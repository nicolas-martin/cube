package lsp

import (
	"context"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/source"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/span"
	errors "golang.org/x/xerrors"
)

func (s *Server) executeCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	switch params.Command {
	case "tidy":
		if len(params.Arguments) == 0 || len(params.Arguments) > 1 {
			return nil, errors.Errorf("expected one file URI for call to `go mod tidy`, got %v", params.Arguments)
		}
		// Confirm that this action is being taken on a go.mod file.
		uri := span.NewURI(params.Arguments[0].(string))
		view, err := s.session.ViewOf(uri)
		if err != nil {
			return nil, err
		}
		snapshot := view.Snapshot()
		fh, err := snapshot.GetFile(uri)
		if err != nil {
			return nil, err
		}
		if fh.Identity().Kind != source.Mod {
			return nil, errors.Errorf("%s is not a mod file", uri)
		}
		// Run go.mod tidy on the view.
		// TODO: This should go through the ModTidyHandle on the view.
		// That will also allow us to move source.InvokeGo into internal/lsp/cache.
		if _, err := source.InvokeGo(ctx, view.Folder().Filename(), snapshot.Config(ctx).Env, "mod", "tidy"); err != nil {
			return nil, err
		}
	}
	return nil, nil
}