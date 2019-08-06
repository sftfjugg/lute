// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"strings"
)

func (context *Context) parseLinkRefDef(tokens items) items {
	_, tokens = tokens.trimLeft()
	if 1 > len(tokens) {
		return nil
	}

	linkLabel, remains, label := context.parseLinkLabel(tokens)
	if nil == linkLabel {
		return nil
	}

	if 1 > len(remains) || itemColon != remains[0] {
		return nil
	}

	remains = remains[1:]
	whitespaces, remains := remains.trimLeft()
	newlines, _, _ := whitespaces.statWhitespace()
	if 1 < newlines {
		return nil
	}

	tokens = remains
	linkDest, remains, destination := context.parseLinkDest(tokens)
	if nil == linkDest {
		return nil
	}

	whitespaces, remains = remains.trimLeft()
	if nil == whitespaces && 0 < len(remains) {
		return nil
	}
	newlines, spaces1, tabs1 := whitespaces.statWhitespace()
	if 1 < newlines {
		return nil
	}

	_, tokens = remains.trimLeft()
	validTitle, _, remains, title := context.parseLinkTitle(tokens)
	if !validTitle && 1 > newlines {
		return nil
	}
	if 0 < spaces1+tabs1 && !remains.isBlankLine() && itemNewline != remains[0] {
		return nil
	}

	titleLine := tokens
	whitespaces, tokens = remains.trimLeft()
	_, spaces2, tabs2 := whitespaces.statWhitespace()
	if !tokens.isBlankLine() && 0 < spaces2+tabs2 {
		title = ""
		remains = titleLine
	} else {
		remains = tokens
	}

	link := &Link{&BaseNode{typ: NodeLink}, destination, ""}
	lowerCaseLabel := strings.ToLower(label)
	link.Title = title
	if _, ok := context.linkRefDef[lowerCaseLabel]; !ok {
		context.linkRefDef[lowerCaseLabel] = link
	}

	return remains
}

func (context *Context) parseLinkTitle(tokens items) (validTitle bool, passed, remains items, title string) {
	if 1 > len(tokens) {
		return true, nil, tokens, ""
	}
	if itemOpenBracket == tokens[0] {
		return true, nil, tokens, ""
	}

	validTitle, passed, remains, title = context.parseLinkTitleMatch(itemDoublequote, itemDoublequote, tokens)
	if !validTitle {
		validTitle, passed, remains, title = context.parseLinkTitleMatch(itemSinglequote, itemSinglequote, tokens)
		if !validTitle {
			validTitle, passed, remains, title = context.parseLinkTitleMatch(itemOpenParen, itemCloseParen, tokens)
		}
	}
	if "" != title {
		title = unescapeString(title)
	}

	return
}

func (context *Context) parseLinkTitleMatch(opener, closer item, tokens items) (validTitle bool, passed, remains items, title string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if opener != tokens[0] {
		return
	}

	line := tokens
	length = len(line)
	closed := false
	i := 1
	size := 0
	var r rune
	for ; i < length; i += size {
		token := line[i]
		passed = append(passed, token)
		r, size = decodeRune(line[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		title += string(r)
		if closer == token && !tokens.isBackslashEscape(i) {
			closed = true
			title = title[:len(title)-1]
			break
		}
	}

	if !closed {
		title = ""
		passed = nil
		return
	}

	validTitle = true
	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkDest(tokens items) (ret, remains items, destination string) {
	ret, remains, destination = context.parseLinkDest1(tokens) // <autolink>
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens) // [label](/url)
	}
	if nil != ret {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}

func (context *Context) parseLinkDest2(tokens items) (ret, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	var openParens int
	i := 0
	size := 0
	var r rune
	for ; i < length; {
		token := tokens[i]
		ret = append(ret, token)
		r, size = decodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			ret = append(ret, tokens[i+j])
		}
		destination += string(r)
		if token.isWhitespace() || token.isControl() {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
			break
		}

		if itemOpenParen == token && !tokens.isBackslashEscape(i) {
			openParens++
		}
		if itemCloseParen == token && !tokens.isBackslashEscape(i) {
			openParens--
			if 1 > openParens {
				i++
				break
			}
		}

		i += size
	}

	remains = tokens[i:]
	if length > i && !tokens[i].isWhitespace() {
		ret = nil
		destination = ""
		return
	}

	return
}

func (context *Context) parseLinkDest1(tokens items) (ret, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemLess != tokens[0] {
		return
	}

	closed := false
	i := 0
	size := 0
	var r rune
	for ; i < length; i += size {
		token := tokens[i]
		ret = append(ret, token)
		size = 1
		if 0 < i {
			r, size = decodeRune(tokens[i:])
			for j := 1; j < size; j++ {
				ret = append(ret, tokens[i+j])
			}
			destination += string(r)
			if itemLess == token && !tokens.isBackslashEscape(i) {
				ret = nil
				destination = ""
				return
			}
		}

		if itemGreater == token && !tokens.isBackslashEscape(i) {
			closed = true
			destination = destination[0 : len(destination)-1]
			break
		}
	}

	if !closed {
		ret = nil
		destination = ""

		return
	}

	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkLabel(tokens items) (passed, remains items, label string) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenBracket != tokens[0] {
		return
	}

	line := tokens
	closed := false
	i := 1
	for {
		token := line[i]
		passed = append(passed, token)
		r, size := decodeRune(line[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		label += string(r)
		if itemCloseBracket == token && !tokens.isBackslashEscape(i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = line[i+1:]
			break
		}
		if itemOpenBracket == token && !tokens.isBackslashEscape(i) {
			passed = nil
			label = ""
			return
		}
		i += size
	}

	if !closed || "" == strings.TrimSpace(label) || 999 < len(label) {
		passed = nil
	}

	label = strings.TrimSpace(label)
	label = strings.ReplaceAll(label, "\n", " ")
	for 0 <= strings.Index(label, "  ") {
		label = strings.ReplaceAll(label, "  ", " ")
	}

	return
}
