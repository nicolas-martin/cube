package gopls

import (
	"context"
	"log"
	"os"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/jsonrpc2"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
)

// GoplsClient ..
type GoplsClient struct {
	server protocol.Server
	point  *types.Point
}

// NewGoPlsClient creates a GoPls client from the local running gopls server
func NewGoPlsClient() *GoplsClient {
	stream := jsonrpc2.NewHeaderStream(os.Stdout, os.Stdin)
	conn := jsonrpc2.NewConn(stream)
	server := protocol.ServerDispatcher(conn)
	ch := &clienthandler{}
	conn.AddHandler(protocol.ClientHandler(ch))
	conn.AddHandler(protocol.Canceller{})
	ctxt := context.Background()
	ctxt = protocol.WithClient(ctxt, ch)

	log.Fatal(conn.Run(ctxt))

	s := loggingGoplsServer{
		u: server,
	}

	goplsClient := &GoplsClient{
		server: s,
		// gopls: gopls.Process,
		// in:    os.Stdin,
		// out:   os.Stdout,
	}

	return goplsClient
}

// func (c *goplsClient) complete(args ...json.RawMessage) (interface{}, error) {
// 	params := &protocol.CompletionParams{
// 		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
// 			TextDocument: protocol.TextDocumentIdentifier{
// 				URI: string(b.URI()),
// 			},
// 			Position: c.point,
// 		},
// 	}
// 	res, err := v.server.Completion(context.Background(), params)
// 	if err != nil {
// 		return nil, fmt.Errorf("called to gopls.Completion failed: %v", err)
// 	}

// 	start := pos.Col()
// 	if len(res.Items) > 0 {
// 		pos, err := types.PointFromPosition(b, res.Items[0].TextEdit.Range.Start)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to derive completion start: %v", err)
// 		}
// 		start = pos.Col() - 1 // see help complete-functions
// 	}
// 	v.lastCompleteResults = res
// 	return start, nil
// }
