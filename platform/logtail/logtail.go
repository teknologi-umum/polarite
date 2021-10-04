package logtail

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"net/http"
)

type Meta struct {
	Level string `json:"level"`
	Date string `json:"date,omitempty"`
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

func Log(message interface{}) error {
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