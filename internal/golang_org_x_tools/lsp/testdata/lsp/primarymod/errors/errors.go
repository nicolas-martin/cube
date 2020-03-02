package errors

import (
	"github.com/nicolas-martin/cube/internal/golang_org_x_tools/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
