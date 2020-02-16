package nodisk

import (
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
