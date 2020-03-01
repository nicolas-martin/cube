package types

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/span"
)

// Point represents a position within a Buffer
type Point struct {
	// line is Vim's line number within the buffer, i.e. 1-indexed
	Line int

	// col is the Vim representation of column number, i.e.  1-based byte index
	Col int

	// offset is the 0-index byte-offset
	offset int

	// is the 0-index character offset in line
	utf16Col int
}

// ToPosition converts p to a protocol.Position
func (p Point) ToPosition() protocol.Position {
	return protocol.Position{
		Line:      p.GoplsLine(),
		Character: p.GoplsChar(),
	}
}

// GoplsLine is the 0-index line in the buffer, returned as a float64 value. This
// is how gopls refers to lines.
func (p Point) GoplsLine() float64 {
	return float64(p.Line - 1)
}

// GoplsChar is the 0-index character offset in a buffer.
func (p Point) GoplsChar() float64 {
	return float64(p.utf16Col)
}

// Buffer represents the current state of the page
type Buffer struct {
	Name     string
	Num      int
	Contents []byte
	// cc is lazily set whenever position information is required
	cc      *span.TokenConverter
	Version int
}

func (b *Buffer) tokenConvertor() *span.TokenConverter {
	if b.cc == nil {
		b.cc = span.NewContentConverter(b.Name, b.Contents)
	}
	return b.cc
}

// URI returns the b's Name as a span.URI, assuming it is a file.
// TODO we should panic here is this is not a file-based buffer
func (b *Buffer) URI() span.URI {
	return span.FileURI(b.Name)
}

// ToTextDocumentIdentifier converts b to a protocol.TextDocumentIdentifier
func (b *Buffer) ToTextDocumentIdentifier() protocol.TextDocumentIdentifier {
	return protocol.TextDocumentIdentifier{
		//Name instead of URI
		URI: string(b.Name),
	}
}

// PointFromPosition converts protocol psoition to point position (?!)
func PointFromPosition(b *Buffer, pos protocol.Position) (Point, error) {
	cc := b.tokenConvertor()
	sline := f2int(pos.Line) + 1
	scol := f2int(pos.Character)
	soff, err := cc.ToOffset(sline, 1)
	p := span.NewPoint(sline, 1, soff)
	p, err = span.FromUTF16Column(p, scol+1, b.Contents)
	if err != nil {
		return Point{}, fmt.Errorf("failed to translate char colum for buffer %v", err)
	}
	res := Point{
		Line:     p.Line(),
		Col:      p.Column(),
		offset:   p.Offset(),
		utf16Col: scol,
	}
	return res, nil
}

func f2int(f float64) int {
	return int(math.Round(f))
}

// Parse parses raw json to struct
func Parse(j json.RawMessage, i interface{}) {
	if err := json.Unmarshal(j, i); err != nil {
		log.Fatalf("failed to parse from %q: %v", j, err)
	}
}
