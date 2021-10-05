package models

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
