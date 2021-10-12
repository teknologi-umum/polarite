// See https://logtail.com/
package logtail

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"net/http"
)

type Meta struct {
	Level  string `json:"level"`
	Date   string `json:"date,omitempty"`
	Source string `json:"source,omitempty"`
}

type Info struct {
	Meta
	Message string `json:"message,omitempty"`
}

type Error struct {
	Meta
	Error string `json:"message,omitempty"`
}

type Logger interface {
	Log(message interface{}) error
}

// Sends a log to Logtail, a logging service.
//
// Sample usage:
//
//      logtail.Log(logtail.Error{
//        Meta: {
//          Level: "error",
//          Date: time.Now().Format(time.RFC3339),
//          Source: "handlers/get.go",
//        },
//        Error: "Some bird flown out of town!",
//      })
func Log(message interface{}) error {
	// Immediate return if not running on production mode
	if os.Getenv("ENVIRONMENT") != "production" {
		return nil
	}

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	body := bytes.NewReader(j)
	req, err := http.NewRequest(http.MethodPost, "https://in.logtail.com", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("LOGTAIL_TOKEN"))

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		var errMsg string
		if resp.StatusCode == http.StatusForbidden {
			errMsg = "Logtail was provided with invalid source token"
		} else if resp.StatusCode == http.StatusNotAcceptable {
			errMsg = "Logtail request body is not valid JSON"
		}
		return errors.New(errMsg)
	}

	return nil
}
