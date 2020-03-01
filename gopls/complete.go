package gopls

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
)

// Complete sends a complete request
func (c *Client) Complete(args ...json.RawMessage) (interface{}, error) {
	params := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				// NOTE: assume 1 file only
				URI: c.Buffer.Name,
			},
			Position: c.Point.ToPosition(),
		},
	}
	res, err := c.Server.Completion(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("called to gopls.Completion failed: %v", err)
	}

	start := c.Point.Col
	if len(res.Items) > 0 {
		pos, err := types.PointFromPosition(c.Buffer, res.Items[0].TextEdit.Range.Start)
		if err != nil {
			return nil, fmt.Errorf("failed to derive completion start: %v", err)
		}
		start = pos.Col - 1 // see help complete-functions
	}
	return start, nil
}
