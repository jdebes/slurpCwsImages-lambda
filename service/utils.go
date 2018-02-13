package service

import (
	"errors"
	"fmt"
	"net/http"
)

func checkStatus(resp *http.Response) error {
	if resp == nil || resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	return errors.New(fmt.Sprintf("Bad response type: %s", resp.StatusCode))
}
