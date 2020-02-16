package gopls

import (
	"context"
	"fmt"
	"strings"

	"github.com/kr/pretty"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
)

type clienthandler struct {
}

var _ protocol.Client = (*clienthandler)(nil)

func (ch *clienthandler) ShowMessage(ctxt context.Context, params *protocol.ShowMessageParams) error {
	defer absorbShutdownErr()
	fmt.Printf("ShowMessage callback: %v", params.Message)

	var hl string
	switch params.Type {
	case protocol.Error:
		hl = "ErrorMsg"
	case protocol.Warning:
		hl = "WarningMsg"
	default:
		return nil
	}

	opts := make(map[string]interface{})
	opts["mousemoved"] = "any"
	opts["moved"] = "any"
	opts["padding"] = []int{0, 1, 0, 1}
	opts["wrap"] = true
	opts["border"] = []int{}
	opts["highlight"] = hl
	opts["line"] = 1
	opts["close"] = "click"

	// g.ChannelCall("popup_create", strings.Split(params.Message, "\n"), opts)
	fmt.Printf(">>> Popup create %s %v", strings.Split(params.Message, "-"), opts)
	return nil
}

func (ch *clienthandler) ShowMessageRequest(context.Context, *protocol.ShowMessageRequestParams) (*protocol.MessageActionItem, error) {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) LogMessage(ctxt context.Context, params *protocol.LogMessageParams) error {
	defer absorbShutdownErr()
	fmt.Printf("LogMessage callback: %v", pretty.Sprint(params))
	return nil
}

func (ch *clienthandler) Telemetry(context.Context, interface{}) error {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) RegisterCapability(ctxt context.Context, params *protocol.RegistrationParams) error {
	defer absorbShutdownErr()
	fmt.Printf("RegisterCapability: %v", pretty.Sprint(params))
	return nil
}

func (ch *clienthandler) UnregisterCapability(context.Context, *protocol.UnregistrationParams) error {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) WorkspaceFolders(context.Context) ([]protocol.WorkspaceFolder, error) {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) Configuration(ctxt context.Context, params *protocol.ParamConfiguration) ([]interface{}, error) {
	defer absorbShutdownErr()

	// TODO this is a rather fragile workaround for https://github.com/golang/go/issues/35817
	// It's fragile because we are relying on gopls not handling any requests until the response
	// to Configuration is received and processed. In practice this appears to currently be
	// the case but there is no guarantee of this going forward. Rather we hope that a fix
	// for https://github.com/golang/go/issues/35817 lands sooner rather than later at whic
	// point this workaround can go.
	//
	// We also use a lock here because, despite it appearing that will only be a single
	// Configuration call and that if there were more they would be serial, we can't rely on
	// this.
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) ApplyEdit(context.Context, *protocol.ApplyWorkspaceEditParams) (*protocol.ApplyWorkspaceEditResponse, error) {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) Event(context.Context, *interface{}) error {
	defer absorbShutdownErr()
	panic("not implemented yet")
}

func (ch *clienthandler) PublishDiagnostics(ctxt context.Context, params *protocol.PublishDiagnosticsParams) error {
	defer absorbShutdownErr()
	fmt.Printf("PublishDiagnostics callback: %v", pretty.Sprint(params))
	// TODO: add some temp logging for https://github.com/golang/go/issues/36601
	// Note this only captures situations where a file's version is increasing.
	// Because it's possible to receive new diagnostics for a file without its
	// version increasing. And in that case it's impossible to know if the
	// diagnostics we receive are old or not.
	fmt.Printf("** Received non-new diagnostics for %v; ignoring. Currently have: \n\nGot: \n%v", params.URI, tabIndent(pretty.Sprint(params)))
	return nil
}

func absorbShutdownErr() {
	if r := recover(); r != nil {
		panic(r)
	}
}

func (ch *clienthandler) logGoplsClientf(format string, args ...interface{}) {
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	fmt.Printf("gopls client start =======================\n"+format+"gopls client end =======================\n", args...)
}

func tabIndent(s string) string {
	return "\t" + strings.ReplaceAll(s, "\n", "\n\t")
}
