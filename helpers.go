// Copyright (c) 2014, David Kitchen <david@buro9.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bluemonday

import (
	"regexp"
)

// A selection of regular expressions that can be used as .Matching() rules on
// HTML attributes.
var (
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-align
	CellAlign = regexp.MustCompile(`(?i)center|justify|left|right|char`)

	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-valign
	CellVerticalAlign = regexp.MustCompile(`(?i)baseline|bottom|middle|top`)

	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdo#attr-dir
	Direction = regexp.MustCompile(`(?i)rtl|ltr`)

	// http://www.w3.org/MarkUp/Test/Img/imgtest.html
	ImageAlign = regexp.MustCompile(
		`(?i)left|right|top|texttop|middle|absmiddle|baseline|bottom|absbottom`,
	)

	// Whole positive integer (including 0) used in places like td.colspan
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-colspan
	Integer = regexp.MustCompile(`[0-9]+`)

	// ISO8601 looks scary but isn't, it's taken from here:
	// http://www.pelagodesign.com/blog/2009/05/20/iso-8601-date-validation-that-doesnt-suck/
	//
	// Minor changes have been made to remove PERL specific syntax that requires
	// regexp backtracking which are not supported in Go
	//
	// Used in places like time.datetime
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/time#attr-datetime
	ISO8601 = regexp.MustCompile(
		`^([\+-]?\d{4})((-?)((0[1-9]|1[0-2])` +
			`([12]\d|0[1-9]|3[01])?|W([0-4]\d|5[0-2])` +
			`(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))` +
			`([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+[^:])?)?` +
			`([0-5]\d([\.,]\d+)?)?([zZ]|([\+-])` +
			`([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`,
	)

	// List types, encapsulates the common value as well as the latest spec
	// values
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ol#attr-type
	ListType = regexp.MustCompile(`(?i)circle|disc|square|a|A|i|I|1`)

	// Textual names/labels separated by spaces, used in places like a.rel and
	// the common attribute 'class'
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a#attr-rel
	NamesAndSpaces = regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)

	// Double value used on HTML5 meter and progress elements
	// http://www.whatwg.org/specs/web-apps/current-work/multipage/the-button-element.html#the-meter-element
	Number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+$`)

	// Whole numbers or %. Used predominantly as units of measurement in width
	// and height attributes
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img#attr-height
	NumberOrPercent = regexp.MustCompile(`[0-9]+%?`)

	// Any block of text in an attribute such as *.'title', img.alt, etc
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes#attr-title
	// Regexp: \p{L} matches unicode letters, \p{N} matches unicode numbers
	// Note that we are not allowing chars that could close tags like '>'
	Paragraph = regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)
)

// AllowStandardURLs is a convenience function that will enable rel="nofollow"
// on "a", "area" and "link" (if you have allowed those elements) and will
// ensure that the URL values are parseable and either relative or belong to the
// "mailto", "http", or "https" schemes
func (p *Policy) AllowStandardURLs() {
	// URLs must be parseable by net/url.Parse()
	p.RequireParseableURLs(true)

	// !url.IsAbs() is permitted
	p.AllowRelativeURLs(true)

	// Most common URL schemes only
	p.AllowURLSchemes("mailto", "http", "https")

	// For all anchors we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)
}

// AllowStandardAttributes will enable "id", "title" and the language specific
// attributes "dir" and "lang" on all elements that are whitelisted
func (p *Policy) AllowStandardAttributes() {
	// "dir" "lang" are permitted as both language attributes affect charsets
	// and direction of text.
	p.AllowAttrs("dir").Matching(Direction).Globally()
	p.AllowAttrs(
		"lang",
	).Matching(regexp.MustCompile(`[a-zA-Z]{2,20}`)).Globally()

	// "id" is permitted. This is pretty much as some HTML elements require this
	// to work well ("dfn" is an example of a "id" being value)
	// This does create a risk that JavaScript and CSS within your web page
	// might identify the wrong elements. Ensure that you select things
	// accurately
	p.AllowAttrs("id").Matching(
		regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`),
	).Globally()

	// "title" is permitted as it improves accessibility.
	p.AllowAttrs("title").Matching(Paragraph).Globally()
}

// AllowImages enables the img element and some popular attributes. It will also
// ensure that URL values are parseable
func (p *Policy) AllowImages() {

	// "img" is permitted
	p.AllowAttrs("align").Matching(ImageAlign).OnElements("img")
	p.AllowAttrs("alt").Matching(Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("img")

	// Standard URLs enabled, which disables data URI images as the data scheme
	// isn't included in the policy
	p.AllowStandardURLs()
	p.AllowAttrs("src").OnElements("img")
}

// AllowLists will enabled ordered and unordered lists, as well as definition
// lists
func (p *Policy) AllowLists() {
	// "ol" "ul" are permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("ol", "ul")

	// "li" is permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("li")
	p.AllowAttrs("value").Matching(Integer).OnElements("li")

	// "dl" "dt" "dd" are permitted
	p.AllowElements("dl", "dt", "dd")
}

// AllowTables will enable a rich set of elements and attributes to describe
// HTML tables
func (p *Policy) AllowTables() {

	// "table" is permitted
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("table")
	p.AllowAttrs("summary").Matching(Paragraph).OnElements("table")

	// "caption" is permitted
	p.AllowElements("caption")

	// "col" "colgroup" are permitted
	p.AllowAttrs("align").Matching(CellAlign).OnElements("col", "colgroup")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("col", "colgroup")
	p.AllowAttrs("span").Matching(Integer).OnElements("colgroup", "col")
	p.AllowAttrs("valign").Matching(
		CellVerticalAlign,
	).OnElements("col", "colgroup")

	// "thead" "tr" are permitted
	p.AllowAttrs("align").Matching(CellAlign).OnElements("thead", "tr")
	p.AllowAttrs("valign").Matching(CellVerticalAlign).OnElements("thead", "tr")

	// "td" "th" are permitted
	p.AllowAttrs("abbr").Matching(Paragraph).OnElements("td", "th")
	p.AllowAttrs("align").Matching(CellAlign).OnElements("td", "th")
	p.AllowAttrs("colspan", "rowspan").Matching(Integer).OnElements("td", "th")
	p.AllowAttrs("headers").Matching(NamesAndSpaces).OnElements("td", "th")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("td", "th")
	p.AllowAttrs(
		"scope",
	).Matching(
		regexp.MustCompile(`(?i)(?:row|col)(?:group)?`),
	).OnElements("td", "th")
	p.AllowAttrs("valign").Matching(CellVerticalAlign).OnElements("td", "th")
	p.AllowAttrs("nowrap").Matching(
		regexp.MustCompile(`(?i)|nowrap`),
	).OnElements("td", "th")

	// "tbody" "tfoot"
	p.AllowAttrs("align").Matching(CellAlign).OnElements("tbody", "tfoot")
	p.AllowAttrs("valign").Matching(
		CellVerticalAlign,
	).OnElements("tbody", "tfoot")
}
