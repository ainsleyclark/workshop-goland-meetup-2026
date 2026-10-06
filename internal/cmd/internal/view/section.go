// Package view renders the CLI's styled output. Commands fetch the
// data and describe what they want shown as a list of sections,
// view decides how each one looks.
package view

import (
	"fmt"
	"workshop/internal/common/printer"
	"workshop/internal/common/printer/styles"
)

// textWidth is how wide a free text paragraph is allowed to grow.
const textWidth = 80

// Section is one part of a command's styled output. A section with
// nothing to show skips itself, so a command can describe its whole
// view up front rather than guarding each part with an if.
type Section interface {
	Render(p *printer.Console)
}

// Render prints each section in order.
func Render(p *printer.Console, sections ...Section) {
	for _, section := range sections {
		section.Render(p)
	}
}

// Card is a bordered detail box. Heading, when set, is printed as a
// section title above the box, otherwise the box is given a blank line
// to breathe. Body is used when there are no key/value Rows.
type Card struct {
	Heading string
	Title   string
	Rows    []printer.KV
	Body    string
}

// Render implements Section.
func (c Card) Render(p *printer.Console) {
	if c.Heading != "" {
		p.Title(c.Heading)
	} else {
		p.LineBreak()
	}

	body := c.Body
	if len(c.Rows) > 0 {
		body = printer.KeyValueString(c.Rows...)
	}

	p.Box(c.Title, body)
}

// Table renders Items as a styled table. Row maps one item to its cells,
// so a command only has to say what a single row looks like. Empty, when
// set, is shown in place of a table with no rows.
type Table[T any] struct {
	Title   string
	Headers []string
	Items   []T
	Row     func(T) []string
	Empty   string
}

// Render implements Section.
func (t Table[T]) Render(p *printer.Console) {
	if t.Title != "" {
		p.Title(t.Title)
	}

	rows := make([][]string, 0, len(t.Items))
	for _, item := range t.Items {
		rows = append(rows, t.Row(item))
	}

	if t.Empty == "" {
		p.Table(t.Headers, rows)
		return
	}

	p.TableOrEmpty(t.Headers, rows, t.Empty)
}

// Text is a titled paragraph, skipped when there's nothing to say.
type Text struct {
	Title string
	Body  string
}

// Render implements Section.
func (t Text) Render(p *printer.Console) {
	if t.Body == "" {
		return
	}
	p.Title(t.Title)
	p.Println(styles.Value.Width(textWidth).Render(t.Body))
}

// List is titled bullet points, skipped when there's nothing to list.
type List struct {
	Title string
	Items []string
}

// Render implements Section.
func (l List) Render(p *printer.Console) {
	if len(l.Items) == 0 {
		return
	}

	items := make([]any, 0, len(l.Items))
	for _, item := range l.Items {
		items = append(items, item)
	}

	p.Title(l.Title)
	p.List(items...)
}

// Lineage is a staircase of nested levels, such as a taxonomy.
type Lineage struct {
	Title string
	Rows  []printer.KV
}

// Render implements Section.
func (l Lineage) Render(p *printer.Console) {
	p.Title(l.Title)
	p.Lineage(l.Rows...)
}

// Details is a titled block of aligned key/value lines, without a box.
type Details struct {
	Title string
	Rows  []printer.KV
}

// Render implements Section.
func (d Details) Render(p *printer.Console) {
	p.Title(d.Title)
	p.KeyValue(d.Rows...)
}

// Heading is a bare section title.
type Heading string

// Render implements Section.
func (h Heading) Render(p *printer.Console) {
	p.Title(string(h))
}

// Muted is a dim standalone line, used for empty states.
type Muted string

// Render implements Section.
func (m Muted) Render(p *printer.Console) {
	p.Muted(string(m))
}

// Hint is a muted footer suggesting what to run next.
type Hint string

// Hintf builds a hint from a format string.
func Hintf(format string, args ...any) Hint {
	return Hint(fmt.Sprintf(format, args...))
}

// Render implements Section.
func (h Hint) Render(p *printer.Console) {
	if h == "" {
		return
	}
	p.LineBreak()
	p.Muted(string(h))
}
