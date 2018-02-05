package service

import (
	"net/http"
	"errors"
	"fmt"
)

func checkStatus(resp *http.Response) error {
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	return errors.New(fmt.Sprintf("Bad response type: %s", resp.StatusCode))
}