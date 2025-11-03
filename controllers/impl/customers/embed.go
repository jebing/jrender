package customers

import "net/http"

type EmbedController struct {
}

func NewEmbedController() *EmbedController {
	return &EmbedController{}
}

func (c *EmbedController) HandleDirectEmbed(w http.ResponseWriter, r *http.Request) {}

func (c *EmbedController) HandleEmbedScript(w http.ResponseWriter, r *http.Request) {}

func (c *EmbedController) HandleFormSubmission(w http.ResponseWriter, r *http.Request) {}
