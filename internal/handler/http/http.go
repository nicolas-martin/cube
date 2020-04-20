package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/repository"
	"github.com/nicolas-martin/cube/internal/types"
	"github.com/nicolas-martin/cube/util"
)

// Handler ..
type Handler struct {
	repo   *repository.Repository
	client *gopls.Client
}

// NewWebHandler returns a Webhandler
func NewWebHandler(r *repository.Repository, c *gopls.Client) *Handler {
	return &Handler{repo: r, client: c}
}

// Ping pong
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Format formats the current buffer
func (h *Handler) Format(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
		c.String(http.StatusBadRequest, "Error reading body")
		return
	}

	folderName := "tmp-rest-folder"

	_, file := util.CreateTmp(folderName, fmt.Sprintf("%s-buffer", folderName))
	h.client.Buffer = &types.Buffer{
		Name:     file.Name(),
		Contents: body,
	}
	_, err = file.Write(body)
	if err != nil {
		log.Panic(err)
		c.String(http.StatusBadRequest, "Error formatting the buffer")
		return
	}

	fmt.Println("^^^^^^^^^^^^^")
	fmt.Println(file.Name())
	fmt.Println(string(body))

	err = h.client.FormatCurrentBuffer(nil)
	if err != nil {
		log.Panic(err)
		c.String(http.StatusBadRequest, "Error formatting the buffer")
		return
	}

	c.String(http.StatusOK, string(h.client.Buffer.Contents))
	return

}
