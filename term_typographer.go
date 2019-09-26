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
	"unicode/utf8"
)

// fixTermTypo 修正 node 下文本节点中出现的术语拼写问题。
func (t *Tree) fixTermTypo(node *Node) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			child.tokens = fixTermTypo0(child.tokens)
		} else {
			t.fixTermTypo(child) // 递归处理子节点
		}
		child = next
	}
}

func fixTermTypo0(tokens items) items {
	length := len(tokens)
	var token byte
	var i, j, k, l int
	var before, after byte
	var originalTerm items
	for ; i < length; i++ {
		token = tokens[i]
		if isNotTerm(token) {
			continue
		}
		if 1 <= i {
			before = tokens[i-1]
			if !isNotTerm(before) {
				// 前一个字节必须是非术语，否则无法分隔
				continue
			}
		}
		if isASCIIPunct(before) {
			// 比如术语前面如果是 . 则不进行修正，因为可能是链接
			// 比如 test.html 虽然不能识别为自动链接，但是也不能进行修正
			continue
		}

		for j = i; j < length; j++ {
			after = tokens[j]
			if isNotTerm(after) || itemDot == after {
				break
			}
		}
		if isASCIIPunct(after) {
			// 比如术语后面如果是 . 则不进行修正，因为可能是链接
			// 比如 github.com 虽然不能识别为自动链接，但是也不能进行修正
			continue
		}

		originalTerm = bytes.ToLower(tokens[i:j])
		if to, ok := terms[fromItems(originalTerm)]; ok {
			l = 0
			for k = i; k < j; k++ {
				tokens[k] = to[l]
				l++
			}
		}
	}

	return tokens
}

func isNotTerm(token byte) bool {
	return token >= utf8.RuneSelf || isWhitespace(token) || isASCIIPunct(token)
}

func replaceAtIndex(str string, r rune, i int) string {
	out := []rune(str)
	out[i] = r
	return string(out)
}

// terms 定义了术语字典，用于术语拼写修正。Key 必须是全小写的。
// TODO: 考虑提供接口支持开发者添加
var terms = map[string]string{
	"jetty":         "Jetty",
	"tomcat":        "Tomcat",
	"jdbc":          "JDBC",
	"mariadb":       "MariaDB",
	"ipfs":          "IPFS",
	"saas":          "SaaS",
	"paas":          "PaaS",
	"iaas":          "IaaS",
	"ioc":           "IoC",
	"freemarker":    "FreeMarker",
	"ruby":          "Ruby",
	"mri":           "MRI",
	"rails":         "Rails",
	"mina":          "Mina",
	"puppet":        "Puppet",
	"vagrant":       "Vagrant",
	"chef":          "Chef",
	"npm":           "NPM",
	"beego":         "Beego",
	"gin":           "Gin",
	"iris":          "Iris",
	"php":           "PHP",
	"ssh":           "SSH",
	"web":           "Web",
	"api":           "API",
	"css":           "CSS",
	"html":          "HTML",
	"json":          "JSON",
	"jsonp":         "JSONP",
	"xml":           "XML",
	"yaml":          "YAML",
	"yml":           "YAML",
	"ini":           "INI",
	"csv":           "CSV",
	"soap":          "SOAP",
	"ajax":          "AJAX",
	"messagepack":   "MessagePack",
	"javascript":    "JavaScript",
	"java":          "Java",
	"jsp":           "JSP",
	"restful":       "RESTFul",
	"gorm":          "GORM",
	"orm":           "ORM",
	"oauth":         "OAuth",
	"facebook":      "Facebook",
	"github":        "GitHub",
	"gist":          "Gist",
	"heroku":        "Heroku",
	"stackoverflow": "Stack Overflow",
	"stackexchange": "Stack Exchange",
	"twitter":       "Twitter",
	"youtube":       "YouTube",
	"dynamodb":      "DynamoDB",
	"mysql":         "MySQL",
	"postgresql":    "PostgreSQL",
	"sqlite":        "SQLite",
	"memcached":     "Memcached",
	"mongodb":       "MongoDB",
	"redis":         "Redis",
	"elasticsearch": "Elasticsearch",
	"solr":          "Solr",
	"solo":          "Solo",
	"sym":           "Sym",
	"b3log":         "B3log",
	"hacpai":        "HacPai",
	"lute":          "Lute",
	"sphinx":        "Sphinx",
	"linux":         "Linux",
	"mac":           "Mac",
	"osx":           "OS X",
	"ubuntu":        "Ubuntu",
	"centos":        "CentOS",
	"centos7":       "CentOS7",
	"redhat":        "RedHat",
	"gitlab":        "GitLab",
	"jquery":        "jQuery",
	"angularjs":     "AngularJS",
	"ffmpeg":        "FFMPEG",
	"git":           "Git",
	"svn":           "SVN",
	"vim":           "VIM",
	"emacs":         "Emacs",
	"sublime":       "Sublime",
	"virtualbox":    "VirtualBox",
	"safari":        "Safari",
	"chrome":        "Chrome",
	"ie":            "IE",
	"firefox":       "Firefox",
	"iterm":         "iTerm",
	"iterm2":        "iTerm2",
	"iwork":         "iWork",
	"itunes":        "iTunes",
	"iphoto":        "iPhoto",
	"ibook":         "iBook",
	"imessage":      "iMessage",
	"photoshop":     "Photoshop",
	"excel":         "Excel",
	"powerpoint":    "PowerPoint",
	"ios":           "iOS",
	"iphone":        "iPhone",
	"ipad":          "iPad",
	"android":       "Android",
	"imac":          "iMac",
	"macbook":       "MacBook",
	"vps":           "VPS",
	"vpn":           "VPN",
	"arm":           "ARM",
	"cpu":           "CPU",
	"spring":        "Spring",
	"springboot":    "SpringBoot",
	"springcloud":   "SpringCloud",
	"sprintmvc":     "SpringMVC",
	"mybatis":       "MyBatis",
	"qq":            "QQ",
}
