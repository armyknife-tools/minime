// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

//go:generate go run gen.go
//go:generate go run gen.go -test

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// identifier converts s to a Go exported identifier.
// It converts "div" to "Div" and "accept-charset" to "AcceptCharset".
func identifier(s string) string {
	b := make([]byte, 0, len(s))
	cap := true
	for _, c := range s {
		if c == '-' {
			cap = true
			continue
		}
		if cap && 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		cap = false
		b = append(b, byte(c))
	}
	return string(b)
}

var test = flag.Bool("test", false, "generate table_test.go")

func genFile(name string, buf *bytes.Buffer) {
	b, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(name, b, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	var all []string
	all = append(all, elements...)
	all = append(all, attributes...)
	all = append(all, eventHandlers...)
	all = append(all, extra...)
	sort.Strings(all)

	// uniq - lists have dups
	w := 0
	for _, s := range all {
		if w == 0 || all[w-1] != s {
			all[w] = s
			w++
		}
	}
	all = all[:w]

	if *test {
		var buf bytes.Buffer
		fmt.Fprintln(&buf, "// Code generated by go generate gen.go; DO NOT EDIT.\n")
		fmt.Fprintln(&buf, "//go:generate go run gen.go -test\n")
		fmt.Fprintln(&buf, "package atom\n")
		fmt.Fprintln(&buf, "var testAtomList = []string{")
		for _, s := range all {
			fmt.Fprintf(&buf, "\t%q,\n", s)
		}
		fmt.Fprintln(&buf, "}")

		genFile("table_test.go", &buf)
		return
	}

	// Find hash that minimizes table size.
	var best *table
	for i := 0; i < 1000000; i++ {
		if best != nil && 1<<(best.k-1) < len(all) {
			break
		}
		h := rand.Uint32()
		for k := uint(0); k <= 16; k++ {
			if best != nil && k >= best.k {
				break
			}
			var t table
			if t.init(h, k, all) {
				best = &t
				break
			}
		}
	}
	if best == nil {
		fmt.Fprintf(os.Stderr, "failed to construct string table\n")
		os.Exit(1)
	}

	// Lay out strings, using overlaps when possible.
	layout := append([]string{}, all...)

	// Remove strings that are substrings of other strings
	for changed := true; changed; {
		changed = false
		for i, s := range layout {
			if s == "" {
				continue
			}
			for j, t := range layout {
				if i != j && t != "" && strings.Contains(s, t) {
					changed = true
					layout[j] = ""
				}
			}
		}
	}

	// Join strings where one suffix matches another prefix.
	for {
		// Find best i, j, k such that layout[i][len-k:] == layout[j][:k],
		// maximizing overlap length k.
		besti := -1
		bestj := -1
		bestk := 0
		for i, s := range layout {
			if s == "" {
				continue
			}
			for j, t := range layout {
				if i == j {
					continue
				}
				for k := bestk + 1; k <= len(s) && k <= len(t); k++ {
					if s[len(s)-k:] == t[:k] {
						besti = i
						bestj = j
						bestk = k
					}
				}
			}
		}
		if bestk > 0 {
			layout[besti] += layout[bestj][bestk:]
			layout[bestj] = ""
			continue
		}
		break
	}

	text := strings.Join(layout, "")

	atom := map[string]uint32{}
	for _, s := range all {
		off := strings.Index(text, s)
		if off < 0 {
			panic("lost string " + s)
		}
		atom[s] = uint32(off<<8 | len(s))
	}

	var buf bytes.Buffer
	// Generate the Go code.
	fmt.Fprintln(&buf, "// Code generated by go generate gen.go; DO NOT EDIT.\n")
	fmt.Fprintln(&buf, "//go:generate go run gen.go\n")
	fmt.Fprintln(&buf, "package atom\n\nconst (")

	// compute max len
	maxLen := 0
	for _, s := range all {
		if maxLen < len(s) {
			maxLen = len(s)
		}
		fmt.Fprintf(&buf, "\t%s Atom = %#x\n", identifier(s), atom[s])
	}
	fmt.Fprintln(&buf, ")\n")

	fmt.Fprintf(&buf, "const hash0 = %#x\n\n", best.h0)
	fmt.Fprintf(&buf, "const maxAtomLen = %d\n\n", maxLen)

	fmt.Fprintf(&buf, "var table = [1<<%d]Atom{\n", best.k)
	for i, s := range best.tab {
		if s == "" {
			continue
		}
		fmt.Fprintf(&buf, "\t%#x: %#x, // %s\n", i, atom[s], s)
	}
	fmt.Fprintf(&buf, "}\n")
	datasize := (1 << best.k) * 4

	fmt.Fprintln(&buf, "const atomText =")
	textsize := len(text)
	for len(text) > 60 {
		fmt.Fprintf(&buf, "\t%q +\n", text[:60])
		text = text[60:]
	}
	fmt.Fprintf(&buf, "\t%q\n\n", text)

	genFile("table.go", &buf)

	fmt.Fprintf(os.Stdout, "%d atoms; %d string bytes + %d tables = %d total data\n", len(all), textsize, datasize, textsize+datasize)
}

type byLen []string

func (x byLen) Less(i, j int) bool { return len(x[i]) > len(x[j]) }
func (x byLen) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byLen) Len() int           { return len(x) }

// fnv computes the FNV hash with an arbitrary starting value h.
func fnv(h uint32, s string) uint32 {
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// A table represents an attempt at constructing the lookup table.
// The lookup table uses cuckoo hashing, meaning that each string
// can be found in one of two positions.
type table struct {
	h0   uint32
	k    uint
	mask uint32
	tab  []string
}

// hash returns the two hashes for s.
func (t *table) hash(s string) (h1, h2 uint32) {
	h := fnv(t.h0, s)
	h1 = h & t.mask
	h2 = (h >> 16) & t.mask
	return
}

// init initializes the table with the given parameters.
// h0 is the initial hash value,
// k is the number of bits of hash value to use, and
// x is the list of strings to store in the table.
// init returns false if the table cannot be constructed.
func (t *table) init(h0 uint32, k uint, x []string) bool {
	t.h0 = h0
	t.k = k
	t.tab = make([]string, 1<<k)
	t.mask = 1<<k - 1
	for _, s := range x {
		if !t.insert(s) {
			return false
		}
	}
	return true
}

// insert inserts s in the table.
func (t *table) insert(s string) bool {
	h1, h2 := t.hash(s)
	if t.tab[h1] == "" {
		t.tab[h1] = s
		return true
	}
	if t.tab[h2] == "" {
		t.tab[h2] = s
		return true
	}
	if t.push(h1, 0) {
		t.tab[h1] = s
		return true
	}
	if t.push(h2, 0) {
		t.tab[h2] = s
		return true
	}
	return false
}

// push attempts to push aside the entry in slot i.
func (t *table) push(i uint32, depth int) bool {
	if depth > len(t.tab) {
		return false
	}
	s := t.tab[i]
	h1, h2 := t.hash(s)
	j := h1 + h2 - i
	if t.tab[j] != "" && !t.push(j, depth+1) {
		return false
	}
	t.tab[j] = s
	return true
}

// The lists of element names and attribute keys were taken from
// https://html.spec.whatwg.org/multipage/indices.html#index
// as of the "HTML Living Standard - Last Updated 16 April 2018" version.

// "command", "keygen" and "menuitem" have been removed from the spec,
// but are kept here for backwards compatibility.
var elements = []string{
	"a",
	"abbr",
	"address",
	"area",
	"article",
	"aside",
	"audio",
	"b",
	"base",
	"bdi",
	"bdo",
	"blockquote",
	"body",
	"br",
	"button",
	"canvas",
	"caption",
	"cite",
	"code",
	"col",
	"colgroup",
	"command",
	"data",
	"datalist",
	"dd",
	"del",
	"details",
	"dfn",
	"dialog",
	"div",
	"dl",
	"dt",
	"em",
	"embed",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"head",
	"header",
	"hgroup",
	"hr",
	"html",
	"i",
	"iframe",
	"img",
	"input",
	"ins",
	"kbd",
	"keygen",
	"label",
	"legend",
	"li",
	"link",
	"main",
	"map",
	"mark",
	"menu",
	"menuitem",
	"meta",
	"meter",
	"nav",
	"noscript",
	"object",
	"ol",
	"optgroup",
	"option",
	"output",
	"p",
	"param",
	"picture",
	"pre",
	"progress",
	"q",
	"rp",
	"rt",
	"ruby",
	"s",
	"samp",
	"script",
	"section",
	"select",
	"slot",
	"small",
	"source",
	"span",
	"strong",
	"style",
	"sub",
	"summary",
	"sup",
	"table",
	"tbody",
	"td",
	"template",
	"textarea",
	"tfoot",
	"th",
	"thead",
	"time",
	"title",
	"tr",
	"track",
	"u",
	"ul",
	"var",
	"video",
	"wbr",
}

// https://html.spec.whatwg.org/multipage/indices.html#attributes-3
//
// "challenge", "command", "contextmenu", "dropzone", "icon", "keytype", "mediagroup",
// "radiogroup", "spellcheck", "scoped", "seamless", "sortable" and "sorted" have been removed from the spec,
// but are kept here for backwards compatibility.
var attributes = []string{
	"abbr",
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"allowfullscreen",
	"allowpaymentrequest",
	"allowusermedia",
	"alt",
	"as",
	"async",
	"autocomplete",
	"autofocus",
	"autoplay",
	"challenge",
	"charset",
	"checked",
	"cite",
	"class",
	"color",
	"cols",
	"colspan",
	"command",
	"content",
	"contenteditable",
	"contextmenu",
	"controls",
	"coords",
	"crossorigin",
	"data",
	"datetime",
	"default",
	"defer",
	"dir",
	"dirname",
	"disabled",
	"download",
	"draggable",
	"dropzone",
	"enctype",
	"for",
	"form",
	"formaction",
	"formenctype",
	"formmethod",
	"formnovalidate",
	"formtarget",
	"headers",
	"height",
	"hidden",
	"high",
	"href",
	"hreflang",
	"http-equiv",
	"icon",
	"id",
	"inputmode",
	"integrity",
	"is",
	"ismap",
	"itemid",
	"itemprop",
	"itemref",
	"itemscope",
	"itemtype",
	"keytype",
	"kind",
	"label",
	"lang",
	"list",
	"loop",
	"low",
	"manifest",
	"max",
	"maxlength",
	"media",
	"mediagroup",
	"method",
	"min",
	"minlength",
	"multiple",
	"muted",
	"name",
	"nomodule",
	"nonce",
	"novalidate",
	"open",
	"optimum",
	"pattern",
	"ping",
	"placeholder",
	"playsinline",
	"poster",
	"preload",
	"radiogroup",
	"readonly",
	"referrerpolicy",
	"rel",
	"required",
	"reversed",
	"rows",
	"rowspan",
	"sandbox",
	"spellcheck",
	"scope",
	"scoped",
	"seamless",
	"selected",
	"shape",
	"size",
	"sizes",
	"sortable",
	"sorted",
	"slot",
	"span",
	"spellcheck",
	"src",
	"srcdoc",
	"srclang",
	"srcset",
	"start",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"translate",
	"type",
	"typemustmatch",
	"updateviacache",
	"usemap",
	"value",
	"width",
	"workertype",
	"wrap",
}

// "onautocomplete", "onautocompleteerror", "onmousewheel",
// "onshow" and "onsort" have been removed from the spec,
// but are kept here for backwards compatibility.
var eventHandlers = []string{
	"onabort",
	"onautocomplete",
	"onautocompleteerror",
	"onauxclick",
	"onafterprint",
	"onbeforeprint",
	"onbeforeunload",
	"onblur",
	"oncancel",
	"oncanplay",
	"oncanplaythrough",
	"onchange",
	"onclick",
	"onclose",
	"oncontextmenu",
	"oncopy",
	"oncuechange",
	"oncut",
	"ondblclick",
	"ondrag",
	"ondragend",
	"ondragenter",
	"ondragexit",
	"ondragleave",
	"ondragover",
	"ondragstart",
	"ondrop",
	"ondurationchange",
	"onemptied",
	"onended",
	"onerror",
	"onfocus",
	"onhashchange",
	"oninput",
	"oninvalid",
	"onkeydown",
	"onkeypress",
	"onkeyup",
	"onlanguagechange",
	"onload",
	"onloadeddata",
	"onloadedmetadata",
	"onloadend",
	"onloadstart",
	"onmessage",
	"onmessageerror",
	"onmousedown",
	"onmouseenter",
	"onmouseleave",
	"onmousemove",
	"onmouseout",
	"onmouseover",
	"onmouseup",
	"onmousewheel",
	"onwheel",
	"onoffline",
	"ononline",
	"onpagehide",
	"onpageshow",
	"onpaste",
	"onpause",
	"onplay",
	"onplaying",
	"onpopstate",
	"onprogress",
	"onratechange",
	"onreset",
	"onresize",
	"onrejectionhandled",
	"onscroll",
	"onsecuritypolicyviolation",
	"onseeked",
	"onseeking",
	"onselect",
	"onshow",
	"onsort",
	"onstalled",
	"onstorage",
	"onsubmit",
	"onsuspend",
	"ontimeupdate",
	"ontoggle",
	"onunhandledrejection",
	"onunload",
	"onvolumechange",
	"onwaiting",
}

// extra are ad-hoc values not covered by any of the lists above.
var extra = []string{
	"acronym",
	"align",
	"annotation",
	"annotation-xml",
	"applet",
	"basefont",
	"bgsound",
	"big",
	"blink",
	"center",
	"color",
	"desc",
	"face",
	"font",
	"foreignObject", // HTML is case-insensitive, but SVG-embedded-in-HTML is case-sensitive.
	"foreignobject",
	"frame",
	"frameset",
	"image",
	"isindex",
	"listing",
	"malignmark",
	"marquee",
	"math",
	"mglyph",
	"mi",
	"mn",
	"mo",
	"ms",
	"mtext",
	"nobr",
	"noembed",
	"noframes",
	"plaintext",
	"prompt",
	"public",
	"rb",
	"rtc",
	"spacer",
	"strike",
	"svg",
	"system",
	"tt",
	"xmp",
}
