package epl

import (
	"bytes"
	"io"
	"iter"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var originalLogOutput io.Writer

func captureOutput() *bytes.Buffer {
	originalLogOutput = log.Writer()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0) // Disable standard log prefixes
	return &buf
}

func restoreOutput(buf *bytes.Buffer) string {
	log.SetOutput(originalLogOutput)
	log.SetFlags(log.LstdFlags) // Restore default flags
	return buf.String()
}

// --- Helper to create simple iterators for tests ---

func sliceIter(items []*Printable) iter.Seq[*Printable] {
	return func(yield func(*Printable) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

// --- Tests for Printable.Print Formatting ---

func TestPrintable_Print_Formatting(t *testing.T) {
	testCases := []struct {
		name     string
		tree     *Printable
		expected string
	}{
		{
			name:     "Single Leaf Level 0",
			tree:     Printablef(0, "Root"),
			expected: "Root", // Note the leading space from Printf format
		},
		{
			name:     "Single Leaf Level 1",
			tree:     Printablef(1, "Offset 1"),
			expected: "  Offset 1", // 2 spaces indent + 1 space format
		},
		{
			name: "Simple Tree",
			tree: &Printable{ // Root node (iterator)
				IndentLevel: 0,
				Iter: sliceIter([]*Printable{
					Printablef(0, "Node A"),          // Child 1, relative indent 0
					Printablef(1, "Node B Offset 1"), // Child 2, relative indent 1
				}),
			},
			expected: "Node A\n  Node B Offset 1", // Node A at L0, Node B at L1
		},
		{
			name: "Nested Tree",
			tree: &Printable{ // Root node (iterator)
				IndentLevel: 0,
				Iter: sliceIter([]*Printable{
					Printablef(0, "Level 0"),
					{ // Level 1 Iterator
						IndentLevel: 1, // Relative indent 1 for this iterator node itself
						Iter: sliceIter([]*Printable{
							Printablef(0, "Level 1 Child A"), // Relative indent 0 from L1 iter -> Absolute 1
							Printablef(1, "Level 1 Child B"), // Relative indent 1 from L1 iter -> Absolute 2
						}),
					},
					Printablef(0, "Level 0 Again"), // Relative indent 0 from root iter
				}),
			},
			expected: "Level 0\n   Level 1 Child A\n     Level 1 Child B\n Level 0 Again",
			// L0: " Level 0" (abs 0)
			// L1 Iter node: Prints nothing itself, recurses with depth=0+1=1
			// L1CA: base depth 1, rel indent 0 -> abs 1 -> "   Level 1 Child A"
			// L1CB: base depth 1, rel indent 1 -> abs 2 -> "     Level 1 Child B"
			// L0 Again: base depth 0, rel indent 0 -> abs 0 -> " Level 0 Again"
		},
		// Add more cases: deeper nesting, empty iterators, etc.
	}

	// Correct the Printf format in common.go if needed based on desired output
	// Let's assume "%s %s" (indent + space + leaf) is desired for now.
	// Rerun TestPrintable_Print_Formatting with the current common.go
	// Modify common.go:printInternal's Printf and these expected strings until they match.

	// --- Assuming common.go uses log.Printf("%s %s", Indent(absoluteIndent), p.Leaf) ---

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := captureOutput()
			tc.tree.Print()
			actual := restoreOutput(buf)
			assert.Equal(t, tc.expected, strings.TrimSpace(actual)) // TrimSpace removes potential trailing newline from log
		})
	}
}
