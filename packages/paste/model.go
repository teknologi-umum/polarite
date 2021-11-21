package paste

import "time"

type Item struct {
	ID        string    `db:"id"`
	Paste     []byte    `db:"content"`
	CreatedAt time.Time `db:"created"`
	Hash      string    `db:"hash"`
	IP        string    `db:"ip"`
	User      string    `db:"user"`
}

type Items []Items

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
