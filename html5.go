package html5

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

type Attr map[string]string
type ElementCollection []*HTMLElement

var unclosedTags map[string]bool = map[string]bool{"input": true}

type HTMLElement struct {
	tag           string
	attributes    Attr
	children      []*HTMLElement
	text          string
	hasClosingTag bool
}

func Element(tag string, attr interface{}, children ...interface{}) *HTMLElement {
	v := HTMLElement{
		tag:        strings.ToLower(tag),
		attributes: toAttrCollection(attr),
		children:   make(ElementCollection, 0),
		text:       ""}

	for i := range children {
		ne := children[i]
		switch ne.(type) {
		case *HTMLElement:
			v.add(ne.(*HTMLElement))
			break
		case []interface{}:
			ea := ne.([]interface{})
			for i := range ea {
				v.add(ea[i])
			}
			break
		case string:
			v.add(Element("", nil).setText(ne.(string)))
		default:
			panic("Type mismatch. Expecting *HTMLElement")
		}
	}

	_, found := unclosedTags[v.tag]
	v.hasClosingTag = !found
	return &v
}

func (me *HTMLElement) AddClass(className string) *HTMLElement {
	parts := strings.Split(me.GetAttribute("class"), " ")
	classfound := false
	for i := range parts {
		if className == parts[i] {
			classfound = true
		}
	}
	if !classfound {
		parts = append(parts, className)
	}
	me.SetAttribute("class", strings.Join(parts, " "))
	return me
}

func (me *HTMLElement) DelClass(className string) *HTMLElement {
	parts := strings.Split(me.GetAttribute("class"), " ")
	newParts := []string{}
	for i := range parts {
		if parts[i] != className && parts[i] != "" {
			newParts = append(newParts, parts[i])
		}
	}
	me.SetAttribute("class", strings.Join(newParts, " "))
	return me
}

func (me *HTMLElement) SetAttribute(name string, value string) *HTMLElement {
	me.attributes[name] = value
	return me
}

func (me *HTMLElement) GetAttribute(name string) string {
	if a, found := me.attributes[name]; found {
		return a
	} else {
		return ""
	}
}

func (me *HTMLElement) DelAttribute(name string) *HTMLElement {
	if _, found := me.attributes[name]; found {
		delete(me.attributes, name)
	}
	return me
}

func (me *HTMLElement) AssignTo(x **HTMLElement) *HTMLElement {
	*x = me
	return me
}

func (me *HTMLElement) setText(s string) *HTMLElement {
	me.text = s
	return me
}

func (me *HTMLElement) SetInnerText(s string) *HTMLElement {
	me.children = make(ElementCollection, 0)
	me.add(s)
	return me
}

func (me *HTMLElement) Empty() *HTMLElement {
	me.children = make(ElementCollection, 0)
	return me
}

func (me *HTMLElement) SetInnerContents(contents interface{}) *HTMLElement {
	switch contents.(type) {
	case string:
		me.Empty()
		me.add(contents)
	case *HTMLElement:
		me.Empty()
		me.add(contents)
	case []*HTMLElement:
		me.children = contents.([]*HTMLElement)
	default:
		panic("Type mismatch")
	}
	return me
}

func (me *HTMLElement) add(elt interface{}) *HTMLElement {
	switch elt.(type) {
	case *HTMLElement:
		me.children = append(me.children, elt.(*HTMLElement))
	case string:
		me.children = append(me.children, Element("", nil).setText(elt.(string)))
	default:
		panic("Type mismatch")
	}

	return me
}
func toAttrCollection(attr interface{}) Attr {
	if attr == nil {
		return nil
	} else {
		switch attr.(type) {
		case Attr:
			return attr.(Attr)
		case string:
			v := make(Attr)
			me := json.Unmarshal([]byte(attr.(string)), &v)
			if me == nil {
				return v
			} else {
				panic(me)
			}

		default:
			panic("Type mismatch")
		}
	}
}

func Document(children ...interface{}) *HTMLElement {
	return Element("html", nil, children...)
}

func (me *HTMLElement) ToHTML(indent ...bool) string {
	clearHTML := false
	if len(indent) > 0 {
		clearHTML = indent[0]
	}
	if clearHTML {
		return me.toHTMLIndent(0)
	} else {
		return me.toHTMLIndent(-9999)
	}
}

func (me *HTMLElement) toHTMLIndent(indent int) string {
	leftPad := ""
	if indent >= 0 {
		leftPad = "\n" + strings.Repeat(" ", indent*2)
	}
	if me.tag == "" {
		return html.EscapeString(me.text)
	} else {
		sb := strings.Builder{}
		if me.tag == "html" {
			sb.WriteString("<!DOCTYPE html>")
		}

		hasElementchildren := false
		if me.attributes != nil {
			asb := strings.Builder{}
			for k, me := range me.attributes {
				asb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, html.EscapeString(me)))
			}

			sb.WriteString(fmt.Sprintf("%s<%s%s>", leftPad, me.tag, asb.String()))
		} else {
			sb.WriteString(fmt.Sprintf("%s<%s>", leftPad, me.tag))
		}
		for i := range me.children {
			if me.children[i].tag != "" {
				hasElementchildren = true
			}
			sb.WriteString(me.children[i].toHTMLIndent(indent + 1))
		}
		if me.hasClosingTag {
			if hasElementchildren {
				sb.WriteString(fmt.Sprintf("%s</%s>", leftPad, me.tag))
			} else {
				sb.WriteString(fmt.Sprintf("</%s>", me.tag))
			}
		}
		return sb.String()
	}
}
