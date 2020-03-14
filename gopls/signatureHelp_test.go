// Package gopls provides ...
package gopls

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/types"
)

//TODO: Make rootURI settable
func TestClient_SignatureHelp(t *testing.T) {
	in := `package abc

import "fmt"

func abc() {

	fmt.Println("test")
}`
	want := &protocol.SignatureHelp{
		Signatures: []protocol.SignatureInformation{
			protocol.SignatureInformation{
				Label:         "Printf(format string, a ...interface{}) (n int, err error)",
				Documentation: "Printf formats according to a format specifier and writes to standard output.",
				Parameters: []protocol.ParameterInformation{
					protocol.ParameterInformation{
						Label: "format string", Documentation: ""},
					protocol.ParameterInformation{
						Label: "a ...interface{}", Documentation: ""},
				},
			},
		},
		ActiveSignature: 0,
		ActiveParameter: 0,
	}
	errChan := make(chan error, 1)
	c := NewGoPlsClient(errChan)
	dir, file := createTmp("signatureHelp_test", "signatureHelp_test-buffer")
	err := file.Chmod(os.FileMode(777))
	if err != nil {
		t.Fatal(err)
	}
	file.Write([]byte(in))

	c.Buffer = &types.Buffer{
		Name:     dir,
		Contents: []byte(in),
	}
	c.Point = &types.Point{
		Line: 7,
		Col:  13,
	}

	got, err := c.SignatureHelp(nil)
	if err != nil {
		t.Errorf("Client.SignatureHelp() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Client.SignatureHelp() = %v, want %v", got, want)
	}

	defer func() {
		os.RemoveAll(dir)
	}()
}

func createTmp(folder, file string) (string, *os.File) {

	tmpFolder, err := ioutil.TempDir("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.TempFile(tmpFolder, file)
	if err != nil {
		log.Fatal(err)
	}

	return tmpFolder, f

}
