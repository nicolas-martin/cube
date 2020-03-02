package danglingstmt

import "github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/foo"

func _() {
	foo. //@rank(" //", Foo)
	var _ = []string{foo.} //@rank("}", Foo)
}
