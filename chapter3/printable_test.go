package chapter3

import (
	"iter"

	epl "github.com/panyam/eplgo"
)

// --- Helper to Extract Items from Printable Iterator ---
// Needed for comparing structures
func collectPrintables(p *epl.Printable) []*epl.Printable {
	if p.Iter == nil {
		return nil // Not an iterator or empty
	}
	var items []*epl.Printable
	// Need iter.Pull to consume the sequence
	pull, stop := iter.Pull(p.Iter)
	defer stop() // Ensure cleanup
	for {
		item, ok := pull()
		if !ok {
			break
		}
		items = append(items, item)
	}
	return items
}
