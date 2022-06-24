package main

import (
	"fmt"
	"os"
	"strings"
)

var tags []string = []string{
	"a", "abbr", "acronym", "address", "applet", "area", "article", "aside", "audio", "b",
	"base", "basefont", "bdi", "bdo", "big", "blockquote", "body", "br", "button", "canvas",
	"caption", "center", "cite", "code", "col", "colgroup", "datalist", "dd", "del", "details",
	"dfn", "dir", "div", "dl", "dt", "em", "embed", "fieldset", "figcaption",
	"figure", "font", "footer", "form", "frame", "frameset", "head", "header", "hgroup",
	"h1", "h2", "h3", "h4", "h5", "h6", "hr", "html", "i", "iframe", "img",
	"input", "ins", "kbd", "keygen", "label", "legend", "li", "link", "map", "mark",
	"menuitem", "menu", "meta", "meter", "nav", "noframes", "noscript", "object",
	"ol", "optgroup", "option", "output", "p", "param", "pre", "progress",
	"q", "rp", "rt", "ruby", "s", "samp", "script", "section", "select", "small",
	"source", "span", "strike", "strong", "style", "sub", "summary", "sup", "table",
	"tbody", "td", "textarea", "tfoot", "th", "thead", "time", "title", "tr", "track",
	"tt", "u", "ul", "var", "video", "wbr"}

func main() {
	fmt.Println("Generating HTML tag functions")

	if f, err := os.Create("functions.go.txt"); err == nil {
		defer f.Close()
		for i := range tags {
			s := "func %funcName(attr interface{}, children ...interface{}) *HTMLElement {\n   return E(\"%tagName\", attr, children...)\n}\n"
			s = strings.ReplaceAll(s, "%funcName", strings.ToUpper(tags[i]))
			s = strings.ReplaceAll(s, "%tagName", strings.ToUpper(tags[i]))
			_, _ = f.WriteString(s)
		}
		fmt.Println("Done")
	} else {
		fmt.Println(err.Error())
	}
}
