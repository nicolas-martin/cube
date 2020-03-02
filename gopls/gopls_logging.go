package gopls

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
)

type loggingGoplsServer struct {
	u protocol.Server
}

var _ protocol.Server = loggingGoplsServer{}

func (l loggingGoplsServer) Logf(format string, args ...interface{}) {
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	fmt.Printf("gopls server start =======================\n"+format+"gopls server end =======================\n", args...)
}

func (l loggingGoplsServer) Initialize(ctxt context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	fmt.Printf("gopls.Initialize() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Initialize(ctxt, params)
	fmt.Printf("gopls.Initialize() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Initialized(ctxt context.Context, params *protocol.InitializedParams) error {
	fmt.Printf("gopls.Initialized() call; params:\n%v", pretty.Sprint(params))
	err := l.u.Initialized(ctxt, params)
	fmt.Printf("gopls.Initialized() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) PrepareRename(ctxt context.Context, params *protocol.PrepareRenameParams) (*protocol.Range, error) {
	fmt.Printf("gopls.PrepareRename() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.PrepareRename(ctxt, params)
	fmt.Printf("gopls.PrepareRename() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}
func (l loggingGoplsServer) Shutdown(ctxt context.Context) error {
	fmt.Printf("gopls.Shutdown() call")
	err := l.u.Shutdown(ctxt)
	fmt.Printf("gopls.Shutdown() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) Exit(ctxt context.Context) error {
	fmt.Printf("gopls.Exit() call")
	err := l.u.Exit(ctxt)
	fmt.Printf("gopls.Exit() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) DidChangeWorkspaceFolders(ctxt context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	fmt.Printf("gopls.DidChangeWorkspaceFolders() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidChangeWorkspaceFolders(ctxt, params)
	fmt.Printf("gopls.DidChangeWorkspaceFolders() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) DidChangeConfiguration(ctxt context.Context, params *protocol.DidChangeConfigurationParams) error {
	fmt.Printf("gopls.DidChangeConfiguration() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidChangeConfiguration(ctxt, params)
	fmt.Printf("gopls.DidChangeConfiguration() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) DidChangeWatchedFiles(ctxt context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	fmt.Printf("gopls.DidChangeWatchedFiles() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidChangeWatchedFiles(ctxt, params)
	fmt.Printf("gopls.DidChangeWatchedFiles() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) Symbol(ctxt context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	fmt.Printf("gopls.Symbol() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Symbol(ctxt, params)
	fmt.Printf("gopls.Symbol() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) ExecuteCommand(ctxt context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	fmt.Printf("gopls.ExecuteCommand() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.ExecuteCommand(ctxt, params)
	fmt.Printf("gopls.ExecuteCommand() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DidOpen(ctxt context.Context, params *protocol.DidOpenTextDocumentParams) error {
	fmt.Printf("gopls.DidOpen() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidOpen(ctxt, params)
	fmt.Printf("gopls.DidOpen() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) DidChange(ctxt context.Context, params *protocol.DidChangeTextDocumentParams) error {
	fmt.Printf("gopls.DidChange() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidChange(ctxt, params)
	fmt.Printf("gopls.DidChange() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) WillSave(ctxt context.Context, params *protocol.WillSaveTextDocumentParams) error {
	fmt.Printf("gopls.WillSave() call; params:\n%v", pretty.Sprint(params))
	err := l.u.WillSave(ctxt, params)
	fmt.Printf("gopls.WillSave() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) WillSaveWaitUntil(ctxt context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	fmt.Printf("gopls.WillSaveWaitUntil() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.WillSaveWaitUntil(ctxt, params)
	fmt.Printf("gopls.WillSaveWaitUntil() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DidSave(ctxt context.Context, params *protocol.DidSaveTextDocumentParams) error {
	fmt.Printf("gopls.DidSave() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidSave(ctxt, params)
	fmt.Printf("gopls.DidSave() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) DidClose(ctxt context.Context, params *protocol.DidCloseTextDocumentParams) error {
	fmt.Printf("gopls.DidClose() call; params:\n%v", pretty.Sprint(params))
	err := l.u.DidClose(ctxt, params)
	fmt.Printf("gopls.DidClose() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) Completion(ctxt context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	fmt.Println("???????????")
	fmt.Printf("gopls.Completion() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Completion(ctxt, params)
	fmt.Fprintln(os.Stderr, "hello")
	log.Fatalf("gopls.Completion() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Resolve(ctxt context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	fmt.Printf("gopls.Resolve() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Resolve(ctxt, params)
	fmt.Printf("gopls.Resolve() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Hover(ctxt context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	fmt.Printf("gopls.Hover() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Hover(ctxt, params)
	fmt.Printf("gopls.Hover() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) SignatureHelp(ctxt context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	fmt.Printf("gopls.SignatureHelp() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.SignatureHelp(ctxt, params)
	fmt.Printf("gopls.SignatureHelp() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Definition(ctxt context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	fmt.Printf("gopls.Definition() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Definition(ctxt, params)
	fmt.Printf("gopls.Definition() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) TypeDefinition(ctxt context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	fmt.Printf("gopls.TypeDefinition() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.TypeDefinition(ctxt, params)
	fmt.Printf("gopls.TypeDefinition() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Implementation(ctxt context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	fmt.Printf("gopls.Implementation() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Implementation(ctxt, params)
	fmt.Printf("gopls.Implementation() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) References(ctxt context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	fmt.Printf("gopls.References() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.References(ctxt, params)
	fmt.Printf("gopls.References() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DocumentHighlight(ctxt context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	fmt.Printf("gopls.DocumentHighlight() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.DocumentHighlight(ctxt, params)
	fmt.Printf("gopls.DocumentHighlight() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DocumentSymbol(ctxt context.Context, params *protocol.DocumentSymbolParams) ([]protocol.DocumentSymbol, error) {
	fmt.Printf("gopls.DocumentSymbol() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.DocumentSymbol(ctxt, params)
	fmt.Printf("gopls.DocumentSymbol() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) CodeAction(ctxt context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	fmt.Printf("gopls.CodeAction() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.CodeAction(ctxt, params)
	fmt.Printf("gopls.CodeAction() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) CodeLens(ctxt context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	fmt.Printf("gopls.CodeLens() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.CodeLens(ctxt, params)
	fmt.Printf("gopls.CodeLens() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) ResolveCodeLens(ctxt context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) {
	fmt.Printf("gopls.ResolveCodeLens() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.ResolveCodeLens(ctxt, params)
	fmt.Printf("gopls.ResolveCodeLens() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DocumentLink(ctxt context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	fmt.Printf("gopls.DocumentLink() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.DocumentLink(ctxt, params)
	fmt.Printf("gopls.DocumentLink() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) ResolveDocumentLink(ctxt context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	fmt.Printf("gopls.ResolveDocumentLink() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.ResolveDocumentLink(ctxt, params)
	fmt.Printf("gopls.ResolveDocumentLink() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) DocumentColor(ctxt context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	fmt.Printf("gopls.DocumentColor() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.DocumentColor(ctxt, params)
	fmt.Printf("gopls.DocumentColor() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) ColorPresentation(ctxt context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	fmt.Printf("gopls.ColorPresentation() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.ColorPresentation(ctxt, params)
	fmt.Printf("gopls.ColorPresentation() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Formatting(ctxt context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	fmt.Printf("gopls.Formatting() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Formatting(ctxt, params)
	fmt.Printf("gopls.Formatting() return; err: %v; res:\n%v\n", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) RangeFormatting(ctxt context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	fmt.Printf("gopls.RangeFormatting() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.RangeFormatting(ctxt, params)
	fmt.Printf("gopls.RangeFormatting() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) OnTypeFormatting(ctxt context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	fmt.Printf("gopls.OnTypeFormatting() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.OnTypeFormatting(ctxt, params)
	fmt.Printf("gopls.OnTypeFormatting() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Rename(ctxt context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	fmt.Printf("gopls.Rename() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Rename(ctxt, params)
	fmt.Printf("gopls.Rename() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) FoldingRange(ctxt context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	fmt.Printf("gopls.FoldingRange() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.FoldingRange(ctxt, params)
	fmt.Printf("gopls.FoldingRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Declaration(ctxt context.Context, params *protocol.DeclarationParams) (protocol.Declaration, error) {
	fmt.Printf("gopls.Declaration() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.Declaration(ctxt, params)
	fmt.Printf("gopls.Declaration() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) LogTraceNotification(ctxt context.Context, params *protocol.LogTraceParams) error {
	fmt.Printf("gopls.LogTraceNotification() call; params:\n%v", pretty.Sprint(params))
	err := l.u.LogTraceNotification(ctxt, params)
	fmt.Printf("gopls.LogTraceNotification() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) SetTraceNotification(ctxt context.Context, params *protocol.SetTraceParams) error {
	fmt.Printf("gopls.SetTraceNotification() call; params:\n%v", pretty.Sprint(params))
	err := l.u.SetTraceNotification(ctxt, params)
	fmt.Printf("gopls.SetTraceNotification() return; err: %v", err)
	return err
}

func (l loggingGoplsServer) SelectionRange(ctxt context.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	fmt.Printf("gopls.SelectionRange() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.SelectionRange(ctxt, params)
	fmt.Printf("gopls.SelectionRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) Progress(ctxt context.Context, params *protocol.ProgressParams) error {
	fmt.Printf("gopls.Progress() call; params:\n%v", pretty.Sprint(params))
	err := l.u.Progress(ctxt, params)
	fmt.Printf("gopls.Progress() return; err: %v\n", err)
	return err
}

func (l loggingGoplsServer) NonstandardRequest(ctxt context.Context, method string, params interface{}) (interface{}, error) {
	fmt.Printf("gopls.NonstandardRequest() call; method: %v, params:\n%v", method, pretty.Sprint(params))
	res, err := l.u.NonstandardRequest(ctxt, method, params)
	fmt.Printf("gopls.NonstandardRequest() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) IncomingCalls(ctxt context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	fmt.Printf("gopls.IncomingCalls() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.IncomingCalls(ctxt, params)
	fmt.Printf("gopls.IncomingCalls() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) OutgoingCalls(ctxt context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	fmt.Printf("gopls.OutgoingCalls() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.OutgoingCalls(ctxt, params)
	fmt.Printf("gopls.OutgoingCalls() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) PrepareCallHierarchy(ctxt context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	fmt.Printf("gopls.PrepareCallHierarchy() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.PrepareCallHierarchy(ctxt, params)
	fmt.Printf("gopls.PrepareCallHierarchy() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) SemanticTokens(ctxt context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	fmt.Printf("gopls.SemanticTokens() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.SemanticTokens(ctxt, params)
	fmt.Printf("gopls.SemanticTokens() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) SemanticTokensEdits(ctxt context.Context, params *protocol.SemanticTokensEditsParams) (interface{}, error) {
	fmt.Printf("gopls.SemanticTokensEdits() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.SemanticTokensEdits(ctxt, params)
	fmt.Printf("gopls.SemanticTokensEdits() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) SemanticTokensRange(ctxt context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	fmt.Printf("gopls.SemanticTokensRange() call; params:\n%v", pretty.Sprint(params))
	res, err := l.u.SemanticTokensRange(ctxt, params)
	fmt.Printf("gopls.SemanticTokensRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (l loggingGoplsServer) WorkDoneProgressCancel(ctxt context.Context, params *protocol.WorkDoneProgressCancelParams) error {
	fmt.Printf("gopls.WorkDoneProgressCancel() call; params:\n%v", pretty.Sprint(params))
	err := l.u.WorkDoneProgressCancel(ctxt, params)
	fmt.Printf("gopls.WorkDoneProgressCancel() return; err: %v\n", err)
	return err
}

func (l loggingGoplsServer) WorkDoneProgressCreate(ctxt context.Context, params *protocol.WorkDoneProgressCreateParams) error {
	fmt.Printf("gopls.WorkDoneProgressCreate() call; params:\n%v", pretty.Sprint(params))
	err := l.u.WorkDoneProgressCreate(ctxt, params)
	fmt.Printf("gopls.WorkDoneProgressCreate() return; err: %v\n", err)
	return err
}
