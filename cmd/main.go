package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/types"
)

func main() {
	t, err := ioutil.ReadFile("tmp-wd/test.go")
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error, 1)
	c := gopls.NewGoPlsClient(errChan)
	c.Buffer = &types.Buffer{
		Name:     "tmp-wd/test.go",
		Contents: t,
	}
	c.Point = &types.Point{
		Line: 7,
		Col:  13,
	}
	fmt.Println("#############")
	err = c.FormatCurrentBuffer()

	fmt.Println("**************")
	// fmt.Printf("r = %+v\n", r)

	msg := <-errChan
	fmt.Println(msg)

}
