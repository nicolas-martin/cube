// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/protocol"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/source"
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/span"
)

func (s *Server) implementation(ctx context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	uri := span.NewURI(params.TextDocument.URI)
	view, err := s.session.ViewOf(uri)
	if err != nil {
		return nil, err
	}
	snapshot := view.Snapshot()
	fh, err := snapshot.GetFile(uri)
	if err != nil {
		return nil, err
	}
	if fh.Identity().Kind != source.Go {
		return nil, nil
	}
	return source.Implementation(ctx, snapshot, fh, params.Position)
}