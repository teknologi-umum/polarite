package controllers

import (
	"errors"
	"net/http"
	"strings"

	h "polarite/platform/highlight"
	"polarite/repository"
	"polarite/resources"

	"github.com/gofiber/fiber/v2"
)

// For parsing query string on GET /:id
// Read documentation on how to parse it here
// https://docs.gofiber.io/api/ctx#queryparser
type QueryString struct {
	// For language syntax highlighting
	Language string `query:"lang"`

	// For syntax colorscheme
	Theme string `query:"theme"`

	// Whether or not to enable line number
	LineNr string `query:"linenr"`
}

// Get route to find a paste by given ID (on path parameters).
func (d *Dependency) Get(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.Context(), "Get")
	defer span.End()

	// Parse the URL param first
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).Send([]byte("ID must not be empty"))
	}

	var qs QueryString
	err := c.QueryParser(&qs)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte("Invalid query string input"))
	}

	// Validate if the ID exists or not
	i, err := d.Paste.GetItemById(ctx, id)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if errors.Is(err, repository.ErrNotFound) {
		return c.Status(http.StatusNotFound).Send([]byte("ID specified was not found"))

	}

	// we need to replace escaped newline to literal newline
	content := strings.Replace(string(i.Paste), `\n`, "\n", -1)

	if qs.Language != "" {
		highlighted, err := h.Highlight(content, qs.Language, qs.Theme, qs.LineNr)
		if err != nil {
			// they should still be able to get the plain text even if the highlighter is b0rked
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(http.StatusOK).Send([]byte(content))
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.Render("template", fiber.Map{
			"Content":     "",
			"Highlighted": highlighted,
			"Truncated":   resources.TruncateString(content, 50),
		})
	}

	if strings.Contains(c.Get(fiber.HeaderAccept), fiber.MIMETextHTML) {
		return c.Render("template", fiber.Map{
			"Content":     content,
			"Highlighted": "",
			"Truncated":   resources.TruncateString(content, 50),
		})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusOK).Send([]byte(content))
}
