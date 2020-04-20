package gopls

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/span"
)

// FormatCurrentBuffer formats the buffer
func (c *Client) FormatCurrentBuffer(args ...json.RawMessage) (err error) {
	//TODO: Call organizeImports
	// err := c.organizeImports(nil)

	var edits []protocol.TextEdit
	formatParams := &protocol.DocumentFormattingParams{
		TextDocument: c.Buffer.ToTextDocumentIdentifier(),
	}
	edits, err = c.Server.Formatting(context.Background(), formatParams)
	if err != nil {
		log.Fatalf("gopls.Formatting returned an error; nothing to do: %s", err.Error())
		return nil
	}
	if len(edits) != 0 {
		return c.applyEasyProtocolTextEdits(edits)
	}
	return nil
}

func (c *Client) organizeImports(args ...json.RawMessage) (err error) {
	//TODO: This returned an error, do it need this?1
	params := &protocol.CodeActionParams{
		TextDocument: c.Buffer.ToTextDocumentIdentifier(),
	}
	actions, err := c.Server.CodeAction(context.Background(), params)
	if err != nil {
		log.Fatalf("gopls.CodeAction returned an error; nothing to do %v", err)
		return nil
	}
	var organizeImports []protocol.CodeAction
	// We might get other kinds in the response, like QuickFix for example.
	// They will be handled via issue #510 (add/enable support for suggested fixes)
	for _, action := range actions {
		if action.Kind == protocol.SourceOrganizeImports {
			organizeImports = append(organizeImports, action)
		}
	}

	switch len(organizeImports) {
	case 0:
	case 1:
		// there should just be a single file
		dcs := organizeImports[0].Edit.DocumentChanges
		switch len(dcs) {
		case 1:
			dc := dcs[0]
			// verify that the URI and version of the edits match the buffer
			euri := span.URI(dc.TextDocument.TextDocumentIdentifier.URI)
			buri := c.Buffer.URI()
			if euri != buri {
				return fmt.Errorf("got edits for file %v, but buffer is %v", euri, buri)
			}
			if ev := int(math.Round(dc.TextDocument.Version)); ev > 0 && ev != c.Buffer.Version {
				return fmt.Errorf("got edits for version %v, but current buffer version is %v", ev, c.Buffer.Version)
			}
			edits := dc.Edits
			if len(edits) != 0 {
				//NOTE: ApplyEdit here
				log.Printf("^^^^^^^^^^^^ Organize import edits: %v", edits)
				// if err := c.applyProtocolTextEdits(edits); err != nil {
				// 	return err
				// }
			}
		default:
			return fmt.Errorf("expected single file, saw %v", len(dcs))
		}
	default:
		return fmt.Errorf("don't know how to handle %v actions", len(organizeImports))
	}

	return err
}
