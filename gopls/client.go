package gopls

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/jsonrpc2"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/diff/myers"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/source"
	"github.com/nicolas-martin/cube/internal/types"
	"gopkg.in/tomb.v2"
)

// Client ..
type Client struct {
	Buffer      *types.Buffer
	Point       *types.Point
	Server      protocol.Server
	goplsCancel context.CancelFunc
	tomb        tomb.Tomb
}

var _ protocol.Client = &clienthandler{}

// NewGoPlsClient creates a GoPls client from the local running gopls server
func NewGoPlsClient(errChan chan error, wdPath string) *Client {
	goplsClient := &Client{}
	// Server
	goplsArgs := []string{"-rpc.trace", "-logfile", "log"}
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 {
		log.Fatal("GOPATH NOT SET")
	}

	gopls := exec.Command(fmt.Sprintf("%s/bin/gopls", goPath), goplsArgs...)

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

	s := loggingGoplsServer{
		u: server,
	}

	goplsClient.Server = s
	params := &protocol.ParamInitialize{}

	//TODO: probably shouldn't hardcode this..
	// params.WorkspaceFolders = []protocol.WorkspaceFolder{
	// 	{Name: "test", URI: string(protocol.URIFromPath(wdPath))},
	// }
	// RootURI is deprecated
	params.RootURI = protocol.URIFromPath(wdPath)

	params.Capabilities.Workspace.Configuration = true
	opts := source.DefaultOptions()
	params.Capabilities.TextDocument.SignatureHelp = protocol.SignatureHelpClientCapabilities{
		DynamicRegistration: false,
		ContextSupport:      true,
	}

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

// DefaultOptions comes from the gopls default, this might be too much for our case.
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
