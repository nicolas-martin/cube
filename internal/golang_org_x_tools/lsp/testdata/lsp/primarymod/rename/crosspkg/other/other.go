package other

import "github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo() //@rename("Foo", "Flamingo")
}
