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

func (codeBlock *Node) codeBlockContinue(context *Context) int {
	var ln = context.currentLine
	var indent = context.indent
	if codeBlock.isFencedCodeBlock {
		if indent <= 3 && codeBlock.isFencedCodeClose(ln[context.nextNonspace:], codeBlock.codeBlockFenceChar, codeBlock.codeBlockFenceLen) {
			context.finalize(codeBlock, context.lineNum)
			return 2
		} else {
			// 跳过围栏标记符之前可能存在的空格
			var i = codeBlock.codeBlockFenceOffset
			var token byte
			for i > 0 {
				token = ln.peek(context.offset)
				if itemSpace != token && itemTab != token {
					break
				}
				context.advanceOffset(1, true)
				i--
			}
		}
	} else { // 缩进代码块
		if indent >= 4 {
			context.advanceOffset(4, true)
		} else if context.blank {
			context.advanceNextNonspace()
		} else {
			return 1
		}
	}
	return 0
}

func (codeBlock *Node) codeBlockFinalize(context *Context) {
	if codeBlock.isFencedCodeBlock {
		content := codeBlock.tokens
		length := len(content)
		if 1 > length {
			return
		}

		var i int
		for ; i < length; i++ {
			if itemNewline == content[i].term {
				break
			}
		}
		firstLine := content[:i]
		rest := content[i+1:]

		codeBlock.codeBlockInfo = unescapeString(trimWhitespace(firstLine))
		codeBlock.tokens = rest
	} else { // 缩进代码块
		codeBlock.tokens = replaceNewlineSpace(codeBlock.tokens)
	}
}

var codeBlockBacktick = strToItems("`")

func (t *Tree) parseFencedCode() (ok bool, codeBlockFenceChar byte, codeBlockFenceLen int, codeBlockFenceOffset int, codeBlockInfo items) {
	marker := t.context.currentLine[t.context.nextNonspace]
	if itemBacktick != marker.term && itemTilde != marker.term {
		return
	}

	fenceChar := marker
	fenceLength := 0
	for i := t.context.nextNonspace; i < t.context.currentLineLen && fenceChar == t.context.currentLine[i]; i++ {
		fenceLength++
	}

	if 3 > fenceLength {
		return
	}

	var info items
	infoTokens := t.context.currentLine[t.context.nextNonspace+fenceLength:]
	if itemBacktick == marker.term && contains(infoTokens, codeBlockBacktick) {
		// info 部分不能包含 `
		return
	}
	info = trimWhitespace(infoTokens)
	info = unescapeString(info)
	return true, fenceChar.term, fenceLength, t.context.indent, info
}

func (codeBlock *Node) isFencedCodeClose(tokens items, openMarker byte, num int) bool {
	closeMarker := tokens[0]
	if closeMarker.term != openMarker {
		return false
	}
	if num > tokens.accept(closeMarker.term) {
		return false
	}
	tokens = trimWhitespace(tokens)
	for _, token := range tokens {
		if token.term != openMarker {
			return false
		}
	}
	return true
}
