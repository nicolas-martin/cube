package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"github.com/nicolas-martin/cube/clog"
	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/handler/http"
	"github.com/nicolas-martin/cube/internal/repository"
	"github.com/nicolas-martin/cube/internal/types"
	log "github.com/sirupsen/logrus"
)

func main() {
	initWebAPI()
	// cmd()
}
func cmd() {
	fn := "tmp-wd/test.go"
	t, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}
	// exeFp, _ := os.Getwd()
	exeFp := "tmp-wd"

	errChan := make(chan error, 1)
	c := gopls.NewGoPlsClient(errChan, exeFp)
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
	l := log.New()
	f, _ := os.Create("gin.log")
	l.SetOutput(f)
	tmpWd, _ := ioutil.TempDir("", "tmp-wd")
	errChan := make(chan error, 1)

	c := gopls.NewGoPlsClient(errChan, tmpWd)

	repo := repository.NewRepository()
	wh := http.NewWebHandler(repo, c)

	r := gin.New()
	gin.DisableConsoleColor()

	// Logging to a file.
	// f, _ := os.Create("gin.log")
	gin.DefaultWriter = l.Out
	// gin.DefaultWriter = io.MultiWriter(f)

	r.Use(clog.Logger(l))
	r.GET("/ping", wh.Ping)
	r.POST("/format", wh.Format)
	r.Run() // listen and serve on 0.0.0.0:8080

	//NOTE: does this even do anything?
	msg := <-errChan
	log.Println(msg)

}
