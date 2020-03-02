package a

import (
	_ "github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/triple/b" //@diag("_ \"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/circular/triple/b\"", "compiler", "import cycle not allowed", "error")
)
