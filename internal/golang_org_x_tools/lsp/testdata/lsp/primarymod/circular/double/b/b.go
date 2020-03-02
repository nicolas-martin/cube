package b

import (
	_ "github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/double/one" //@diag("_ \"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/double/one\"", "compiler", "import cycle not allowed", "error"),diag("\"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/double/one\"", "compiler", "could not import github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/double/one (no package for import github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/double/one)", "error")
)
