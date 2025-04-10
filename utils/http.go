// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// PrepareURL return an url.URL from it's parameters
func PrepareURL(schemeHost string, path string, params map[string]string) (*url.URL, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	u, err := url.Parse(schemeHost)
	if err != nil {
		return nil, err
	}

	u.ForceQuery = true
	u.Path = path
	u.RawQuery = values.Encode()

	return u, nil
}

// FetchURL return body of a fetched url.URL
func FetchURL(ctx context.Context, logger log.Logger, u *url.URL) ([]byte, error) {
	_ = level.Debug(logger).Log("url", u, "context", ctx, "msg", "Fetching URL")

	hresp, err := ctxhttp.Get(ctx, http.DefaultClient, u.String())
	if err != nil {
		return nil, err
	}
	defer hresp.Body.Close()

	body, err := io.ReadAll(hresp.Body)
	_ = level.Debug(logger).Log("len(body)", len(body), "err", err, "msg", "Reading HTTP response body")
	if err != nil {
		return nil, err
	}

	if hresp.StatusCode >= 400 {
		return body, errors.New(hresp.Status)
	}

	return body, nil
}
