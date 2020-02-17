package gopls

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/jsonrpc2"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
	"gopkg.in/tomb.v2"
)

// Client ..
type Client struct {
	Server protocol.Server
	Point  *types.Point
	B      *types.Buffer
	tomb   tomb.Tomb
}

// NewGoPlsClient creates a GoPls client from the local running gopls server
func NewGoPlsClient(errChan chan error) *Client {
	goplsClient := &Client{}
	// Server
	goplsArgs := []string{"-rpc.trace", "-logfile", "log"}
	gopls := exec.Command("/Users/nmartin/go/bin/gopls", goplsArgs...)

	stdout, err := gopls.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to create stdout pipe for gopls: %v", err)
	}
	stdin, err := gopls.StdinPipe()
	if err != nil {
		log.Fatalf("failed to create stdin pipe for gopls: %v", err)
	}
	stderr, err := gopls.StderrPipe()
	if err != nil {
		log.Fatalf("failed to create stderr pipe for gopls: %v", err)
	}
	goplsClient.tomb.Go(func() error {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("gopls stderr: %v", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("reading standard input: %v", err)
		}
		return nil
	})
	if err := gopls.Start(); err != nil {
		log.Fatalf("failed to start gopls: %v", err)
	}

	goplsClient.tomb.Go(func() (err error) {
		if err = gopls.Wait(); err != nil {
			err = fmt.Errorf("got error running gopls: %v", err)
		}
		select {
		// case <-g.inShutdown:
		// 	return nil
		default:
			if err != nil {
				errChan <- err
			}
			return
		}
	})
	stream := jsonrpc2.NewStream(stdout, stdin)
	conn := jsonrpc2.NewConn(stream)
	server := protocol.ServerDispatcher(conn)

	// Client
	ch := &clienthandler{}
	conn.AddHandler(protocol.ClientHandler(ch))
	conn.AddHandler(protocol.Canceller{})
	ctxt := context.Background()
	ctxt = protocol.WithClient(ctxt, ch)
	goplsClient.tomb.Go(func() error {
		return conn.Run(ctxt)
	})

	s := loggingGoplsServer{
		u: server,
	}

	goplsClient.Server = s

	if _, err := goplsClient.Server.Initialize(context.Background(), nil); err != nil {
		log.Fatalf("failed to initialise gopls: %v", err)
	}

	if err := goplsClient.Server.Initialized(context.Background(), &protocol.InitializedParams{}); err != nil {
		log.Fatalf("failed to call gopls.Initialized: %v", err)
	}

	return goplsClient
}

// Complete sends a complete request
func (c *Client) Complete(args ...json.RawMessage) (interface{}, error) {
	params := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				// NOTE: assume 1 file only
				URI: c.B.Name,
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
		pos, err := types.PointFromPosition(c.B, res.Items[0].TextEdit.Range.Start)
		if err != nil {
			return nil, fmt.Errorf("failed to derive completion start: %v", err)
		}
		start = pos.Col - 1 // see help complete-functions
	}
	return start, nil
}
