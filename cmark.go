package cmark

// #cgo LDFLAGS: -lcmark
// #include <string.h>
// #include <stdlib.h>
// #include <cmark.h>
import "C"
import (
	"errors"
	"unsafe"
)

// Parser is a parser for CommonMark
type Parser struct {
	parser *C.cmark_parser
}

// Opt CommonMark options
type Opt C.int

const (
	OptDefault      Opt = C.CMARK_OPT_DEFAULT
	OptSourcePos        = C.CMARK_OPT_SOURCEPOS
	OptHardBreaks       = C.CMARK_OPT_HARDBREAKS
	OptSafe             = C.CMARK_OPT_SAFE
	OptNoBreaks         = C.CMARK_OPT_NOBREAKS
	OptValidateUtf8     = C.CMARK_OPT_VALIDATE_UTF8
	OptSmart            = C.CMARK_OPT_SMART
)

// NewParser builds a parser with the given options
// when finished call Close
func NewParser(options Opt) Parser {
	return Parser{parser: C.cmark_parser_new(C.int(options))}
}

// Write bytes to the parser using the streaming interface
func (p Parser) Write(b []byte) (n int, err error) {
	buf := C.CBytes(b)
	sz := len(b)
	C.cmark_parser_feed(p.parser, (*C.char)(buf), C.size_t(sz))
	C.free(buf)
	return sz, nil
}

// Tree returns the root node for the generated document
// Call this method only once, and then call Close
func (p Parser) Tree() Node {
	return Node{node: C.cmark_parser_finish(p.parser)}
}

// Close frees the wrapped CommonMark Parser
func (p Parser) Close() {
	C.cmark_parser_free(p.parser)
}

type Node struct {
	node *C.cmark_node
}

// NodeType contains the type of a CommonMark AST node
type NodeType C.cmark_node_type

const (
	NodeNone NodeType = C.CMARK_NODE_NONE

	// Block types

	NodeDocument      = C.CMARK_NODE_DOCUMENT
	NodeBlockQuote    = C.CMARK_NODE_BLOCK_QUOTE
	NodeList          = C.CMARK_NODE_LIST
	NodeItem          = C.CMARK_NODE_ITEM
	NodeCodeBlock     = C.CMARK_NODE_CODE_BLOCK
	NodeHTMLBlock     = C.CMARK_NODE_HTML_BLOCK
	NodeCustomBlock   = C.CMARK_NODE_CUSTOM_BLOCK
	NodeParagraph     = C.CMARK_NODE_PARAGRAPH
	NodeHeading       = C.CMARK_NODE_HEADING
	NodeThematicBreak = C.CMARK_NODE_THEMATIC_BREAK

	NodeFirstBlock = C.CMARK_NODE_FIRST_BLOCK
	NodeLastBlock  = C.CMARK_NODE_LAST_BLOCK

	// Inline types

	NodeText         = C.CMARK_NODE_TEXT
	NodeSoftBreak    = C.CMARK_NODE_SOFTBREAK
	NodeLineBreak    = C.CMARK_NODE_LINEBREAK
	NodeCode         = C.CMARK_NODE_CODE
	NodeHTMLInline   = C.CMARK_NODE_HTML_INLINE
	NodeCustomInline = C.CMARK_NODE_CUSTOM_INLINE
	NodeEmph         = C.CMARK_NODE_EMPH
	NodeStrong       = C.CMARK_NODE_STRONG
	NodeLink         = C.CMARK_NODE_LINK
	NodeImage        = C.CMARK_NODE_IMAGE

	NodeFirstInline = C.CMARK_NODE_FIRST_INLINE
	NodeLastInline  = C.CMARK_NODE_LAST_INLINE
)

func (n Node) Next() Node {
	return Node{node: C.cmark_node_next(n.node)}
}

func (n Node) Previous() Node {
	return Node{node: C.cmark_node_previous(n.node)}
}

func (n Node) Parent() Node {
	return Node{node: C.cmark_node_parent(n.node)}
}

func (n Node) FirstChild() Node {
	return Node{node: C.cmark_node_first_child(n.node)}
}

func (n Node) LastChild() Node {
	return Node{node: C.cmark_node_last_child(n.node)}
}

// UserData returns the UserData associated with a node
func (n Node) UserData() unsafe.Pointer {
	return C.cmark_node_get_user_data(n.node)
}

// SetUserData sets the UserData associated with a node
func (n Node) SetUserData(u unsafe.Pointer) {
	C.cmark_node_set_user_data(n.node, u)
}

func (n Node) Type() (NodeType, error) {
	typ := NodeType(C.cmark_node_get_type(n.node))
	if typ == NodeNone {
		return typ, errors.New("Node type could not be determined")
	}
	return typ, nil
}

// TypeString returns a string for a node's type or "<unknown>" on error
func (n Node) TypeString() string {
	str := C.cmark_node_get_type_string(n.node)
	gstr := C.GoString(str)
	C.free(unsafe.Pointer(str))
	return gstr
}

// Literal returns the content of the node
func (n Node) Literal() string {
	return C.GoString(C.cmark_node_get_literal(n.node))
}

// SetLiteral overwrites the literal with a string
// the old string, if any, is not freed
func (n Node) SetLiteral(lit string) {
	C.cmark_node_set_literal(n.node, C.CString(lit))
}

// HeadingLevel returns the heading level of a node
// e.g. 1 for an h1, etc., or 0 if this node is not a heading
func (n Node) HeadingLevel() (int, error) {
	level := int(C.cmark_node_get_heading_level(n.node))
	if level == 0 {
		return level, errors.New("Node is not a heading")
	}
	return level, nil
}

type ListType C.cmark_list_type

const (
	_NoList     ListType = C.CMARK_NO_LIST
	BulletList           = C.CMARK_BULLET_LIST
	OrderedList          = C.CMARK_ORDERED_LIST
)

type ListDelim C.cmark_delim_type

const (
	_NoDelim    ListDelim = C.CMARK_NO_DELIM
	PeriodDelim           = C.CMARK_PERIOD_DELIM
	ParenDelim            = C.CMARK_PAREN_DELIM
)

func (n Node) ListType() (ListType, error) {
	typ := ListType(C.cmark_node_get_list_type(n.node))
	if typ == _NoList {
		return typ, errors.New("Node is not a list")
	}
	return typ, nil
}

func (n Node) SetListType(typ ListType) error {
	if C.cmark_node_set_list_type(n.node, C.cmark_list_type(typ)) == 0 {
		return errors.New("List type could not be set")
	}
	return nil
}

func (n Node) ListDelim() (ListDelim, error) {
	typ := ListDelim(C.cmark_node_get_list_delim(n.node))
	if typ == _NoDelim {
		return typ, errors.New("Node is not a list")
	}
	return typ, nil
}

func (n Node) SetListDelim(typ ListDelim) error {
	if C.cmark_node_set_list_delim(n.node, C.cmark_delim_type(typ)) == 0 {
		return errors.New("List type could not be set")
	}
	return nil
}

func (n Node) ListStart() (int, error) {
	start := int(C.cmark_node_get_list_start(n.node))
	if start == 0 {
		return start, errors.New("ListStart can only be called on ordered lists")
	}
	return start, nil
}

// SetListStart sets the list start number for an ordered list
func (n Node) SetListStart(start int) error {
	if C.cmark_node_set_list_start(n.node, C.int(start)) == 0 {
		return errors.New("SetListStart failed")
	}
	return nil
}

// TightList returns true if the node is a list and list is "tight"
func (n Node) TightList() bool {
	return C.cmark_node_get_list_tight(n.node) == 1
}

func (n Node) SetTightList(tight bool) error {
	t := 0
	if tight {
		t = 1
	}
	if C.cmark_node_set_list_tight(n.node, C.int(t)) == 0 {
		return errors.New("SetTightList failed")
	}
	return nil
}

// FenceInfo returns the info string at a code block fence
// (e.g. ""```ruby" would return "ruby")
func (n Node) FenceInfo() string {
	return C.GoString(C.cmark_node_get_fence_info(n.node))
}

func (n Node) SetFenceInfo(fence string) error {
	if C.cmark_node_set_fence_info(n.node, C.CString(fence)) == 0 {
		return errors.New("SetFenceInfo failed")
	}
	return nil
}

func (n Node) URL() string {
	return C.GoString(C.cmark_node_get_url(n.node))
}

func (n Node) SetURL(url string) error {
	if C.cmark_node_set_url(n.node, C.CString(url)) == 0 {
		return errors.New("SetURL failed")
	}
	return nil
}

func (n Node) Title() string {
	return C.GoString(C.cmark_node_get_title(n.node))
}

func (n Node) SetTitle(title string) error {
	if C.cmark_node_set_title(n.node, C.CString(title)) == 0 {
		return errors.New("SetTitle failed")
	}
	return nil
}

func (n Node) OnEnter() string {
	return C.GoString(C.cmark_node_get_on_enter(n.node))
}

func (n Node) SetOnEnter(onEnter string) error {
	if C.cmark_node_set_on_enter(n.node, C.CString(onEnter)) == 0 {
		return errors.New("SetOnEnter failed")
	}
	return nil
}

func (n Node) OnExit() string {
	return C.GoString(C.cmark_node_get_on_exit(n.node))
}

func (n Node) SetOnExit(onExit string) error {
	if C.cmark_node_set_on_exit(n.node, C.CString(onExit)) == 0 {
		return errors.New("SetOnExit failed")
	}
	return nil
}

func (n Node) StartLine() int {
	return int(C.cmark_node_get_start_line(n.node))
}

func (n Node) StartColumn() int {
	return int(C.cmark_node_get_start_column(n.node))
}

func (n Node) EndtLine() int {
	return int(C.cmark_node_get_end_line(n.node))
}

func (n Node) EndColumn() int {
	return int(C.cmark_node_get_end_column(n.node))
}

// Unlink unlinks node but does not Close it,
// call Close if it is no longer needed
func (n Node) Unlink() {
	C.cmark_node_unlink(n.node)
}

func (n Node) InsertBefore(s Node) error {
	if C.cmark_node_insert_before(n.node, s.node) == 0 {
		return errors.New("InsertBefore failed")
	}
	return nil
}

func (n Node) InsertAfter(s Node) error {
	if C.cmark_node_insert_after(n.node, s.node) == 0 {
		return errors.New("InsertAfter failed")
	}
	return nil
}

// Replaces replaces this node with another,
// call Close on the old node if no longer needed
func (o Node) Replace(n Node) error {
	if C.cmark_node_replace(o.node, n.node) == 0 {
		return errors.New("Replace failed")
	}
	return nil
}

func (n Node) PrependChild(c Node) error {
	if C.cmark_node_prepend_child(n.node, c.node) == 0 {
		return errors.New("PrependChild failed")
	}
	return nil
}

func (n Node) AppendChild(c Node) error {
	if C.cmark_node_append_child(n.node, c.node) == 0 {
		return errors.New("AppendChild failed")
	}
	return nil
}

// ConsolidateTextNodes consolidates adjacent text nodes into one text node
// for the sub-tree of this node
func (n Node) ConsolidateTextNodes() {
	C.cmark_consolidate_text_nodes(n.node)
}

// SetHeadingLevel sets heading level to value (1 for h1, etc.)
func (n Node) SetHeadingLevel(level int) error {
	if C.cmark_node_set_heading_level(n.node, C.int(level)) == 0 {
		return errors.New("Heading could not be set")
	}
	return nil
}

// Close frees the wrapped CommonMark Node
func (n Node) Close() {
	C.cmark_node_free(n.node)
}

// RenderHTML renders html from the document
func (n Node) RenderHTML(options Opt) string {
	html := C.cmark_render_html(n.node, C.int(options))
	gstr := C.GoString(html)
	C.free(unsafe.Pointer(html))
	return gstr
}

// RenderXML renders xml from the document
// This rendering is basically a serialization of the AST
func (n Node) RenderXML(options Opt) string {
	xml := C.cmark_render_xml(n.node, C.int(options))
	gstr := C.GoString(xml)
	C.free(unsafe.Pointer(xml))
	return gstr
}

// RenderMan renders a manual page from the document using troff
// wrapWidth is the wrap width (0 indicates no wrapping)
func (n Node) RenderMan(options Opt, wrapWidth int) string {
	man := C.cmark_render_man(n.node, C.int(options), C.int(wrapWidth))
	gstr := C.GoString(man)
	C.free(unsafe.Pointer(man))
	return gstr
}

// RenderLaTeX renders LaTeX from the document
// wrapWidth is the wrap width (0 indicates no wrapping)
func (n Node) RenderLaTeX(options Opt, wrapWidth int) string {
	latex := C.cmark_render_latex(n.node, C.int(options), C.int(wrapWidth))
	gstr := C.GoString(latex)
	C.free(unsafe.Pointer(latex))
	return gstr
}

// RenderCommonMark renders CommonMark Markdown from the document
// wrapWidth is the wrap width (0 indicates no wrapping)
//
// This method is especially useful for formatting markdown, as it produces
// a canonical CommonMark output
func (n Node) RenderCommonMark(options Opt, wrapWidth int) string {
	markdown := C.cmark_render_commonmark(n.node, C.int(options), C.int(wrapWidth))
	gstr := C.GoString(markdown)
	C.free(unsafe.Pointer(markdown))
	return gstr
}

type Event C.cmark_event_type

const (
	EventNone  Event = C.CMARK_EVENT_NONE
	EventDone        = C.CMARK_EVENT_DONE
	EventEnter       = C.CMARK_EVENT_ENTER
	EventExit        = C.CMARK_EVENT_EXIT
)

type Iter struct {
	iter *C.cmark_iter
}

func (n Node) Iter() Iter {
	return Iter{iter: C.cmark_iter_new(n.node)}
}

// Next advances the iterator and returns the event that has occurred,
// which may be EventEnter, EventExit, or EventDone
func (i Iter) Next() Event {
	return Event(C.cmark_iter_next(i.iter))
}

// Node returns the current node the iterator is pointing to
//
// It is not necessary to Close this node as it is in the tree
// but it is necessary to Close the root node when done iterating
func (i Iter) Node() Node {
	return Node{node: C.cmark_iter_get_node(i.iter)}
}

// Event returns the event which the last advance emitted
func (i Iter) Event() Event {
	return Event(C.cmark_iter_get_event_type(i.iter))
}

// Root returns the root node of the tree this iterator is
// iterating over
func (i Iter) Root() Node {
	return Node{node: C.cmark_iter_get_root(i.iter)}
}

// Reset resets the iterator to a node and event
// Node must be a child of the root
func (i Iter) Reset(n Node, e Event) {
	C.cmark_iter_reset(i.iter, n.node, C.cmark_event_type(e))
}

func (i Iter) Close() {
	C.cmark_iter_free(i.iter)
}
