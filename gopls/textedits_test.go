package gopls

import (
	"context"
	"strings"
	"testing"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_applyEasyProtocolTextEdits(t *testing.T) {
	e := `package abc

import "fmt"

func abc() {

	fmt.Println("test")
}
`
	e = strings.ReplaceAll(e, "\n", "\r\n")

	in := `package abc

import "fmt"

func abc() {

	fmt.Println(  "test")
}
`
	c := createMockClient()
	c.Buffer = &types.Buffer{
		Name:     "tmp-wd/test.go",
		Contents: []byte(in),
	}
	c.Server = &ServerMock{
		DidChangeFunc: func(in1 context.Context, in2 *protocol.DidChangeTextDocumentParams) error {
			return nil
		},
	}

	edits := []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 6, Character: 13},
				End:   protocol.Position{Line: 6, Character: 15},
			},
			NewText: "",
		},
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 6, Character: 22},
				End:   protocol.Position{Line: 6, Character: 23},
			},
			NewText: "",
		}}

	type args struct {
		edits []protocol.TextEdit
	}
	tests := []struct {
		name     string
		fields   *Client
		args     args
		expected string
		wantErr  bool
	}{
		{
			name:     "test",
			fields:   c,
			args:     args{edits: edits},
			expected: e,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.applyEasyProtocolTextEdits(tt.args.edits); (err != nil) != tt.wantErr {
				t.Errorf("Client.applyEasyProtocolTextEdits() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, string(c.Buffer.Contents), "Buffers does not match")
		})
	}
}

func createClient() *Client {
	errChan := make(chan error, 1)
	c := NewGoPlsClient(errChan)
	return c
}

func createMockClient() *Client {
	goplsClient := &Client{}
	return goplsClient
}

type inOut struct {
	in    string
	out   string
	edits []protocol.TextEdit
}
