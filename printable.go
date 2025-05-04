package epl

import (
	"fmt"
	"iter"
	"log"
)

type Printable struct {
	IndentLevel int
	// one of:
	Leaf string
	Iter iter.Seq[*Printable]
}

func Indent(level int) string {
	out := ""
	for range level {
		out += "  "
	}
	return out
}

func (p *Printable) Print() {
	p.print(0)
}

func (p *Printable) print(depth int) {
	if p.Leaf != "" {
		log.Printf("%s%s", Indent(depth+p.IndentLevel), p.Leaf)
	} else {
		for printable := range p.Iter {
			printable.print(depth + p.IndentLevel)
		}
	}
}

func PrintableIter(i iter.Seq[*Printable]) (out *Printable) {
	out = &Printable{Iter: i}
	return out
}

func Printablef(level int, format string, args ...any) (out *Printable) {
	out = &Printable{
		IndentLevel: level,
		Leaf:        fmt.Sprintf(format, args...),
	}
	return
}
