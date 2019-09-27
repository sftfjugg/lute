// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"bytes"
	"errors"
)

// RendererFunc 描述了渲染器函数签名。
type RendererFunc func(n *Node, entering bool) (WalkStatus, error)

// Renderer 描述了渲染器接口。
type Renderer interface {
	// Render 渲染输出。
	Render() (output []byte, err error)
}

// BaseRenderer 描述了渲染器结构。
type BaseRenderer struct {
	writer              *bytes.Buffer             // 输出缓冲
	lastOut             byte                      // 最新输出的一个字节
	rendererFuncs       map[nodeType]RendererFunc // 渲染器
	defaultRendererFunc RendererFunc              // 默认渲染器，在 rendererFuncs 中找不到节点渲染器时会使用该默认渲染器进行渲染
	disableTags         int                       // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况，比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>。
	option              *options                  // 解析渲染选项
	treeRoot            *Node                     // 待渲染的树的根节点
}

// Render 从指定的根节点 root 开始遍历并渲染。
func (r *BaseRenderer) Render() (output []byte, err error) {
	defer recoverPanic(&err)

	r.lastOut = itemNewline
	r.writer = &bytes.Buffer{}
	r.writer.Grow(4096)

	err = Walk(r.treeRoot, func(n *Node, entering bool) (WalkStatus, error) {
		f := r.rendererFuncs[n.typ]
		if nil == f {
			if nil != r.defaultRendererFunc {
				return r.defaultRendererFunc(n, entering)
			} else {
				return r.renderDefault(n, entering)
			}
		}
		return f(n, entering)
	})
	if nil != err {
		return
	}

	output = r.writer.Bytes()
	return
}

func (r *BaseRenderer) renderDefault(n *Node, entering bool) (WalkStatus, error) {
	return WalkStop, errors.New("not found render function for node [type=" + n.typ.String() + ", tokens=" + itemsToStr(n.tokens) + "]")
}

// writeByte 输出一个字节 c。
func (r *BaseRenderer) writeByte(c byte) {
	r.writer.WriteByte(c)
	r.lastOut = c
}

// writeBytes 输出字节数组 bytes。
func (r *BaseRenderer) writeBytes(bytes []byte) {
	if length := len(bytes); 0 < length {
		r.writer.Write(bytes)
		r.lastOut = bytes[length-1]
	}
}

// write 输出指定的 tokens 数组 content。
func (r *BaseRenderer) write(content items) {
	bytes := itemsToBytes(content)
	if length := len(bytes); 0 < length {
		r.writer.Write(bytes)
		r.lastOut = bytes[length-1]
	}
}

// writeString 输出指定的字符串 content。
func (r *BaseRenderer) writeString(content string) {
	if length := len(content); 0 < length {
		r.writer.WriteString(content)
		r.lastOut = content[length-1]
	}
}

// newline 会在最新内容不是换行符 \n 时输出一个换行符。
func (r *BaseRenderer) newline() {
	if itemNewline != r.lastOut {
		r.writer.WriteByte(itemNewline)
		r.lastOut = itemNewline
	}
}
