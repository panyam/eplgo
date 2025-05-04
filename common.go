package epl

import (
	"fmt"
	"iter"
	"log"
)

func Dict[K comparable, V any](args ...any) (out map[K]V) {
	out = map[K]V{}
	for i := 0; i < len(args); i += 2 {
		k := args[i].(K)
		v := args[i+1].(V)
		out[k] = v
	}
	return out
}

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
		log.Printf("%s %s", Indent(depth+p.IndentLevel), p.Leaf)
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
