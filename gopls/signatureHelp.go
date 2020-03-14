// Package gopls provides ...
package gopls

import (
	"context"
	"encoding/json"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
)

func (c *Client) SignatureHelp(args ...json.RawMessage) (*protocol.SignatureHelp, error) {
	p, err := types.PointFromVim(c.Buffer, c.Point.Line, c.Point.Col)
	if err != nil {
		return nil, err
	}

	shp := &protocol.SignatureHelpParams{
		Context: protocol.SignatureHelpContext{
			TriggerKind:      0.0,
			TriggerCharacter: ",",
			IsRetrigger:      false,
			ActiveSignatureHelp: protocol.SignatureHelp{
				Signatures:      nil,
				ActiveSignature: 0.0,
				ActiveParameter: 0.0,
			},
		},
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: c.Buffer.ToTextDocumentIdentifier().URI,
			},
			Position: p.ToPosition(),
		},
		WorkDoneProgressParams: protocol.WorkDoneProgressParams{
			WorkDoneToken: nil,
		},
	}

	resp, err := c.Server.SignatureHelp(context.Background(), shp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
