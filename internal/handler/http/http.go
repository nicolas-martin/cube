package http

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/cube/gopls"
	"github.com/nicolas-martin/cube/internal/repository"
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
	err := h.client.FormatCurrentBuffer(nil)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading body")
		return
	}
	h.client.Buffer.URI

	h.client.Buffer.SetContents(body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error formatting the buffer")
		return
	}

	c.String(http.StatusOK, string(h.client.Buffer.Contents))
	return

}
