package gopls

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
)

type textEdit struct {
	buffer int
	call   string
	start  int
	end    int
	lines  []string
}

func (c *Client) applyEasyProtocolTextEdits(edits []protocol.TextEdit) error {
	b := c.Buffer
	// prepare the changes to make in Vim
	blines := bytes.Split(b.Contents[:len(b.Contents)-1], []byte("\n"))
	var changes []textEdit
	for ie := len(edits) - 1; ie >= 0; ie-- {
		e := edits[ie]
		start, err := types.PointFromPosition(b, e.Range.Start)
		if err != nil {
			return fmt.Errorf("failed to derive start point from position: %v", err)
		}
		end, err := types.PointFromPosition(b, e.Range.End)
		if err != nil {
			return fmt.Errorf("failed to derive end point from position: %v", err)
		}
		// Skip empty edits
		if start == end && e.NewText == "" {
			continue
		}
		// special case deleting of complete lines
		if start.Col == 1 && end.Col == 1 && e.NewText == "" {
			delstart := min(start.Line, len(blines))
			delend := min(end.Line-1, len(blines))
			changes = append(changes, textEdit{
				call:   "deletebufline",
				buffer: b.Num,
				start:  delstart,
				end:    delend,
			})
			blines = append(blines[:delstart-1], blines[delend:]...)
			continue
		}
		newLines := strings.Split(e.NewText, "\n")
		appendAdjust := 1
		if start.Line-1 < len(blines) {
			appendAdjust = 0
			startLine := blines[start.Line-1]
			newLines[0] = string(startLine[:start.Col-1]) + newLines[0]
			if end.Line-1 < len(blines) {
				endLine := blines[end.Line-1]
				newLines[len(newLines)-1] += string(endLine[end.Col-1:])
			}
			// we only need to update the start line because we can't have
			// overlapping edits
			blines[start.Line-1] = []byte(newLines[0])
			changes = append(changes, textEdit{
				call:   "setbufline",
				buffer: b.Num,
				start:  start.Line,
				lines:  []string{string(blines[start.Line-1])},
			})
		} else {
			blines = append(blines, []byte(newLines[0]))
		}
		if start.Line != end.Line {
			// We can't delete beyond the end of the buffer. So the start end end here are
			// both min() reduced
			delstart := min(start.Line+1, len(blines))
			delend := min(end.Line, len(blines))
			changes = append(changes, textEdit{
				call:   "deletebufline",
				buffer: b.Num,
				start:  delstart,
				end:    delend,
			})
			blines = blines[:delstart-1]
		}
		if len(newLines) > 1 {
			changes = append(changes, textEdit{
				call:   "appendbufline",
				buffer: b.Num,
				start:  start.Line - appendAdjust,
				lines:  newLines[1-appendAdjust : len(newLines)-appendAdjust],
			})
		}
	}

	for _, e := range changes {
		switch e.call {
		case "setbufline":
			// blines[e.start-1] = []byte(e.lines[0])
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.lines[0])
		case "deletebufline":
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.end)
		case "appendbufline":
			//TODO: iterate / append range for e.lines
			// Create a new line
			blines = append(blines, nil)
			// Shift everything
			copy(blines[e.start+1:], blines[e.start:])
			// Insert
			blines[e.start] = []byte(e.lines[0])
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.lines)
		default:
			panic(fmt.Errorf("unknown change type: %v", e.call))
		}
	}

	newContents := ""
	for _, v := range blines {
		newContents += fmt.Sprintf("%s\r\n", v)
	}

	c.Buffer.SetContents([]byte(newContents))
	c.Buffer.Version++
	params := &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: c.Buffer.ToTextDocumentIdentifier(),
			Version:                float64(c.Buffer.Version),
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Text: newContents,
			},
		},
	}
	return c.Server.DidChange(context.Background(), params)
}

func (c *Client) applyProtocolTextEdits(edits []protocol.TextEdit) error {
	// prepare the changes to make in Vim
	blines := bytes.Split(c.Buffer.Contents[:len(c.Buffer.Contents)-1], []byte("\n"))
	var changes []textEdit
	for ie := len(edits) - 1; ie >= 0; ie-- {
		e := edits[ie]
		start, err := types.PointFromPosition(c.Buffer, e.Range.Start)
		if err != nil {
			return fmt.Errorf("failed to derive start point from position: %v", err)
		}
		end, err := types.PointFromPosition(c.Buffer, e.Range.End)
		if err != nil {
			return fmt.Errorf("failed to derive end point from position: %v", err)
		}
		// Skip empty edits
		if start == end && e.NewText == "" {
			continue
		}

		// special case deleting of complete lines
		newLines := strings.Split(e.NewText, "\n")
		appendAdjust := 1
		if start.Line-1 < len(blines) {
			appendAdjust = 0
			startLine := blines[start.Line-1]
			newLines[0] = string(startLine[:start.Col-1]) + newLines[0]
			if end.Line-1 < len(blines) {
				endLine := blines[end.Line-1]
				newLines[len(newLines)-1] += string(endLine[end.Col-1:])
			}
			// we only need to update the start line because we can't have
			// overlapping edits
			blines[start.Line-1] = []byte(newLines[0])
			changes = append(changes, textEdit{
				call:   "setbufline",
				buffer: 1,
				start:  start.Line,
				lines:  []string{string(blines[start.Line-1])},
			})
		} else {
			blines = append(blines, []byte(newLines[0]))
		}

		if start.Line != end.Line {
			// We can't delete beyond the end of the buffer. So the start end end here are
			// both min() reduced
			delstart := min(start.Line+1, len(blines))
			delend := min(end.Line, len(blines))
			changes = append(changes, textEdit{
				call:   "deletebufline",
				buffer: 1,
				start:  delstart,
				end:    delend,
			})
			blines = blines[:delstart-1]
		}

		if len(newLines) > 1 {
			changes = append(changes, textEdit{
				call:   "appendbufline",
				buffer: 1,
				start:  start.Line - appendAdjust,
				lines:  newLines[1-appendAdjust : len(newLines)-appendAdjust],
			})
		}
	}

	for _, e := range changes {
		switch e.call {
		case "setbufline":
			blines[e.start] = []byte(e.lines[0])
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.lines[0])
		case "deletebufline":
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.end)
		case "appendbufline":
			// blines = blines[:e.start]
			//TODO: Should interate through e.lines[]
			// blines[e.start+1] = []byte(e.lines[0])
			fmt.Printf("%s = %v %v \r\n", e.call, e.start, e.lines)
		default:
			panic(fmt.Errorf("unknown change type: %v", e.call))
		}
	}

	// c.Buffer.Contents = []byte(newContents)
	// c.Buffer.Version++
	// params := &protocol.DidChangeTextDocumentParams{
	// 	TextDocument: protocol.VersionedTextDocumentIdentifier{
	// 		TextDocumentIdentifier: c.Buffer.ToTextDocumentIdentifier(),
	// 		Version:                float64(c.Buffer.Version),
	// 	},
	// 	ContentChanges: []protocol.TextDocumentContentChangeEvent{
	// 		{
	// 			Text: newContents,
	// 		},
	// 	},
	// }
	// return c.Server.DidChange(context.Background(), params)
	return nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
