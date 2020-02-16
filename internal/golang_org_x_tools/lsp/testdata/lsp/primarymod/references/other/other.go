package other

import (
	references "github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/references"
)

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
