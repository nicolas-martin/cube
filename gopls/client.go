package gopls

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/jsonrpc2"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/diff/myers"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/source"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/span"
	"github.com/nicolas-martin/cube/internal/types"
	"gopkg.in/tomb.v2"
)

// Client ..
type Client struct {
	Server      protocol.Server
	Point       *types.Point
	B           *types.Buffer
	tomb        tomb.Tomb
	goplsCancel context.CancelFunc
}

var _ protocol.Client = &clienthandler{}

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
			log.Printf("reading standard input: %v", err)
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
	stream := jsonrpc2.NewHeaderStream(stdout, stdin)
	ctxt, cancel := context.WithCancel(context.Background())
	goplsClient.goplsCancel = cancel
	conn := jsonrpc2.NewConn(stream)
	server := protocol.ServerDispatcher(conn)

	// Client
	ch := &clienthandler{}
	conn.AddHandler(protocol.ClientHandler(ch))
	conn.AddHandler(protocol.Canceller{})
	ctxt = protocol.WithClient(ctxt, ch)

	goplsClient.tomb.Go(func() error {
		return conn.Run(ctxt)
	})

	// s := loggingGoplsServer{
	// 	u: server,
	// }

	goplsClient.Server = server
	params := &protocol.ParamInitialize{}
	// params.RootURI = string(span.FileURI(c.Client.app.wd))
	params.RootURI = string(span.FileURI("nm-empty-wd"))
	params.Capabilities.Workspace.Configuration = true
	opts := source.DefaultOptions()
	params.Capabilities.TextDocument.Hover = protocol.HoverClientCapabilities{
		ContentFormat: []protocol.MarkupKind{opts.PreferredContentFormat},
	}

	if _, err := goplsClient.Server.Initialize(context.Background(), params); err != nil {
		log.Fatalf("failed to initialise gopls: %v", err)
	}

	if err := goplsClient.Server.Initialized(context.Background(), &protocol.InitializedParams{}); err != nil {
		log.Fatalf("failed to call gopls.Initialized: %v", err)
	}

	return goplsClient
}
func DefaultOptions() source.Options {
	return source.Options{
		ClientOptions: source.ClientOptions{
			InsertTextFormat:              protocol.PlainTextTextFormat,
			PreferredContentFormat:        protocol.Markdown,
			ConfigurationSupported:        true,
			DynamicConfigurationSupported: true,
			DynamicWatchedFilesSupported:  true,
			LineFoldingOnly:               false,
		},
		ServerOptions: source.ServerOptions{
			SupportedCodeActions: map[source.FileKind]map[protocol.CodeActionKind]bool{
				source.Go: {
					protocol.SourceOrganizeImports: true,
					protocol.QuickFix:              true,
				},
				source.Mod: {
					protocol.SourceOrganizeImports: true,
				},
				source.Sum: {},
			},
			SupportedCommands: []string{
				"tidy", // for go.mod files
			},
		},
		UserOptions: source.UserOptions{
			Env:                     os.Environ(),
			HoverKind:               source.FullDocumentation,
			LinkTarget:              "pkg.go.dev",
			Matcher:                 source.Fuzzy,
			DeepCompletion:          true,
			UnimportedCompletion:    true,
			CompletionDocumentation: true,
		},
		DebuggingOptions: source.DebuggingOptions{
			CompletionBudget: 100 * time.Millisecond,
		},
		ExperimentalOptions: source.ExperimentalOptions{
			TempModfile: false,
		},
		Hooks: source.Hooks{
			ComputeEdits: myers.ComputeEdits,
			// URLRegexp:    regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`),
			URLRegexp: nil,
			// Analyzers:    source.defaultAnalyzers(),
			Analyzers: nil,
			GoDiff:    true,
		},
	}
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

// FormatCurrentBuffer formats the buffer
func (c *Client) FormatCurrentBuffer(args ...json.RawMessage) (err error) {
	// we are an autocmd endpoint so we need to be told the current
	// buffer number via <abuf>
	// _ = types.ParseInt(args[0])

	params := &protocol.CodeActionParams{
		TextDocument: c.B.ToTextDocumentIdentifier(),
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
			buri := c.B.URI()
			if euri != buri {
				return fmt.Errorf("got edits for file %v, but buffer is %v", euri, buri)
			}
			if ev := int(math.Round(dc.TextDocument.Version)); ev > 0 && ev != c.B.Version {
				return fmt.Errorf("got edits for version %v, but current buffer version is %v", ev, c.B.Version)
			}
			edits := dc.Edits
			if len(edits) != 0 {
				//NOTE: ApplyEdit here
				log.Printf("^^^ Organize import edits: %v", edits)
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
	var edits []protocol.TextEdit
	formatParams := &protocol.DocumentFormattingParams{
		TextDocument: c.B.ToTextDocumentIdentifier(),
	}
	edits, err = c.Server.Formatting(context.Background(), formatParams)
	if err != nil {
		log.Fatalf("gopls.Formatting returned an error; nothing to do")
		return nil
	}
	if len(edits) != 0 {
		//NOTE: ApplyEdit here
		// return c.applyProtocolTextEdits(edits)
		log.Printf("$$$$ Formatting edits: %v", edits)
	}
	return nil
}

type textEdit struct {
	buffer int
	call   string
	start  int
	end    int
	lines  []string
}

// func (c *Client) applyProtocolTextEdits(edits []protocol.TextEdit) error {

// 	// prepare the changes to make in Vim
// 	blines := bytes.Split(c.B.Contents[:len(c.B.Contents)-1], []byte("\n"))
// 	var changes []textEdit
// 	for ie := len(edits) - 1; ie >= 0; ie-- {
// 		e := edits[ie]
// 		start, err := types.PointFromPosition(c.B, e.Range.Start)
// 		if err != nil {
// 			return fmt.Errorf("failed to derive start point from position: %v", err)
// 		}
// 		end, err := types.PointFromPosition(c.B, e.Range.End)
// 		if err != nil {
// 			return fmt.Errorf("failed to derive end point from position: %v", err)
// 		}
// 		// Skip empty edits
// 		if start == end && e.NewText == "" {
// 			continue
// 		}
// 		// special case deleting of complete lines
// 		newLines := strings.Split(e.NewText, "\n")
// 		appendAdjust := 1
// 		if start.Line-1 < len(blines) {
// 			appendAdjust = 0
// 			startLine := blines[start.Line-1]
// 			newLines[0] = string(startLine[:start.Col-1]) + newLines[0]
// 			if end.Line-1 < len(blines) {
// 				endLine := blines[end.Line-1]
// 				newLines[len(newLines)-1] += string(endLine[end.Col-1:])
// 			}
// 			// we only need to update the start line because we can't have
// 			// overlapping edits
// 			blines[start.Line-1] = []byte(newLines[0])
// 			changes = append(changes, textEdit{
// 				call:   "setbufline",
// 				buffer: 1,
// 				start:  start.Line,
// 				lines:  []string{string(blines[start.Line-1])},
// 			})
// 		} else {
// 			blines = append(blines, []byte(newLines[0]))
// 		}
// 		if start.Line != end.Line {
// 			// We can't delete beyond the end of the buffer. So the start end end here are
// 			// both min() reduced
// 			delstart := min(start.Line+1, len(blines))
// 			delend := min(end.Line, len(blines))
// 			changes = append(changes, textEdit{
// 				call:   "deletebufline",
// 				buffer: 1,
// 				start:  delstart,
// 				end:    delend,
// 			})
// 			blines = blines[:delstart-1]
// 		}
// 		if len(newLines) > 1 {
// 			changes = append(changes, textEdit{
// 				call:   "appendbufline",
// 				buffer: 1,
// 				start:  start.Line - appendAdjust,
// 				lines:  newLines[1-appendAdjust : len(newLines)-appendAdjust],
// 			})
// 		}
// 	}

// 	var newContents string
// 	types.Parse(newContentsRes(), &newContents)
// 	c.B.Contents = []byte(newContents)
// 	c.B.Version++
// 	params := &protocol.DidChangeTextDocumentParams{
// 		TextDocument: protocol.VersionedTextDocumentIdentifier{
// 			TextDocumentIdentifier: c.B.ToTextDocumentIdentifier(),
// 			Version:                float64(c.B.Version),
// 		},
// 		ContentChanges: []protocol.TextDocumentContentChangeEvent{
// 			{
// 				Text: newContents,
// 			},
// 		},
// 	}
// 	return c.Server.DidChange(context.Background(), params)
// }
// func min(i, j int) int {
// 	if i < j {
// 		return i
// 	}
// 	return j
// }
