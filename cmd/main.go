package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/handler/http"
	"github.com/nicolas-martin/cube/internal/repository"
	"github.com/nicolas-martin/cube/internal/types"
)

func main() {
	initWebAPI()
	// cmd()
}
func cmd() {
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

	err = c.FormatCurrentBuffer(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\r\n= RESULT =")
	fmt.Println(pretty.Print(string(c.Buffer.Contents)))
	fmt.Println("= END RESULT =")
	msg := <-errChan
	fmt.Println(msg)
}

func initWebAPI() {
	errChan := make(chan error, 1)
	c := gopls.NewGoPlsClient(errChan)

	repo := repository.NewRepository()
	wh := http.NewWebHandler(repo, c)

	r := gin.Default()
	r.GET("/ping", wh.Ping)
	r.POST("/format", wh.Format)
	r.Run() // listen and serve on 0.0.0.0:8080

	//NOTE: does this even do anything?
	msg := <-errChan
	fmt.Println(msg)

}
