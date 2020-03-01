package gopls

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/govim/cmd/govim/config"
)

const (
	goplsConfigNoDocsOnHover  = "noDocsOnHover"
	goplsConfigHoverKind      = "hoverKind"
	goplsDeepCompletion       = "deepCompletion"
	goplsCompletionMatcher    = "matcher"
	goplsStaticcheck          = "staticcheck"
	goplsCompleteUnimported   = "completeUnimported"
	goplsGoImportsLocalPrefix = "local"
	goplsCompletionBudget     = "completionBudget"
	goplsTempModfile          = "tempModfile"
	goplsVerboseOutput        = "verboseOutput"
	goplsEnv                  = "env"
)

type clienthandler struct {
	config config.Config
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

	ch.logGoplsClientf("Configuration: %v", pretty.Sprint(params))

	conf := ch.config

	// gopls now sends params.Items for each of the configured
	// workspaces. For now, we assume that the first item will be
	// for the section "gopls" and only configure that. We will
	// configure further workspaces when we add support for them.
	if len(params.Items) == 0 || params.Items[0].Section != "gopls" {
		return nil, fmt.Errorf("govim gopls client: expected at least one item, with the first section \"gopls\"")
	}
	res := make([]interface{}, len(params.Items))
	goplsConfig := make(map[string]interface{})
	goplsConfig[goplsConfigHoverKind] = "FullDocumentation"
	if conf.CompletionDeepCompletions != nil {
		goplsConfig[goplsDeepCompletion] = *conf.CompletionDeepCompletions
	}
	if conf.CompletionMatcher != nil {
		goplsConfig[goplsCompletionMatcher] = *conf.CompletionMatcher
	}
	if conf.Staticcheck != nil {
		goplsConfig[goplsStaticcheck] = *conf.Staticcheck
	}
	if conf.CompleteUnimported != nil {
		goplsConfig[goplsCompleteUnimported] = *conf.CompleteUnimported
	}
	if conf.GoImportsLocalPrefix != nil {
		goplsConfig[goplsGoImportsLocalPrefix] = *conf.GoImportsLocalPrefix
	}
	if conf.CompletionBudget != nil {
		goplsConfig[goplsCompletionBudget] = *conf.CompletionBudget
	}
	if os.Getenv(string(config.EnvVarGoplsVerbose)) == "true" {
		goplsConfig[goplsVerboseOutput] = true
	}
	res[0] = goplsConfig

	ch.logGoplsClientf("Configuration response: %v", pretty.Sprint(res))
	return res, nil
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
