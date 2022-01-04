package servicecheck

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	//nolint:gosec // This is the well-known path to Kubernetes serviceaccount tokens.
	tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

// doRequest does an http request only to get the http status code
func (c *Checker) doRequest(url string) (string, error) {
	// Read Bearer Token file from ServiceAccount
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return "error", fmt.Errorf("could not load token %s: %s", tokenFile, err)
	}

	req, _ := http.NewRequest("GET", url, http.NoBody)

	// Only add the Bearer for API Server Requests
	if strings.HasSuffix(url, "/version") {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err.Error(), err
	}

	// Body is non-nil if err is nil, so close it
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "ok", nil
	}

	return resp.Status, errors.New(resp.Status)
}
