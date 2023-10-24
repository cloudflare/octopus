//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const httpTimeout = time.Minute

func FetchHTTP(url string) ([]byte, error) {
	return fetchHTTPWithHeaders(url, make(http.Header))
}

func FetchHTTPWithHeaders(url string, header map[string]string) ([]byte, error) {
	httpHeader := make(http.Header)
	for k, v := range header {
		httpHeader.Add(k, v)
	}

	return fetchHTTPWithHeaders(url, httpHeader)
}

func fetchHTTPWithHeaders(url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request: %v", err)
	}

	req.Header = header
	req.Header.Set("User-Agent", getUserAgent())

	client := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with status %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading HTTP response failed: %v", err)
	}

	return ret, nil
}

func getUserAgent() string {
	envVar := os.Getenv("OCTOPUS_USER_AGENT")
	if envVar != "" {
		return envVar
	}

	return "octopus"
}
