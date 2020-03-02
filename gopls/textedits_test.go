package gopls

import (
	"context"
	"fmt"
	"testing"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_applyEasyProtocolTextEdits(t *testing.T) {
	e := `
package abc

import "fmt"

func abc() {

    fmt.Println("test")
}`
	in := `
package abc

import "fmt"

func abc() {

	fmt.Println(  "test")
}`
	c := createClient()
	c.Buffer.SetContents([]byte(in))

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

	type fields struct {
		Buffer      *types.Buffer
		Point       *types.Point
		Server      protocol.Server
		goplsCancel context.CancelFunc
	}
	type args struct {
		edits []protocol.TextEdit
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected string
		wantErr  bool
	}{
		{
			name:     "test",
			fields:   fields{},
			args:     args{edits: edits},
			expected: e,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Buffer:      tt.fields.Buffer,
				Point:       tt.fields.Point,
				Server:      tt.fields.Server,
				goplsCancel: tt.fields.goplsCancel,
			}
			if err := c.applyEasyProtocolTextEdits(tt.args.edits); (err != nil) != tt.wantErr {
				t.Errorf("Client.applyEasyProtocolTextEdits() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, fmt.Sprint(c.Buffer), "Buffers does not match")
		})
	}
}

func createClient() *Client {
	errChan := make(chan error, 1)
	c := NewGoPlsClient(errChan)
	return c
}

type inOut struct {
	in    string
	out   string
	edits []protocol.TextEdit
}
