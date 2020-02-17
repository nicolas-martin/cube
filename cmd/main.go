package main

import (
	"fmt"
	"log"

	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/types"
)

func main() {
	testString := `package main

import (
	"os"
)

func main() {
	s, err := os.

}`
	errChan := make(chan error, 1)
	c := gopls.NewGoPlsClient(errChan)
	c.B = &types.Buffer{
		Name:     "-Test-",
		Contents: []byte(testString),
	}
	c.Point = &types.Point{
		Line: 8,
		Col:  18,
	}
	resp, err := c.Complete()
	if err != nil {
		log.Fatal(err)
	}

	msg := <-errChan
	fmt.Println(msg)
	fmt.Println(resp)

}
