/*
   Copyright 2016 Continusec Pty Ltd

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package continusec provides golang client libraries for interacting with the
// verifiable datastructures provided by Continusec.
//
// Users should start with the NewClient function.
package continusec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is the object that will be used to interact with the Continusec API.
// Call NewClient to construct.
type Client struct {
	account    string
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient returns a new client suitable to interact with the Continusec API.
//
// account, although accepted as a string, is usually the integer shown on your account
// settings page.
//
// apiKey is the string configured on your API Access page, and may be left blank for
// any data that is publicly accessible.
func NewClient(account string, apiKey string) *Client {
	return &Client{
		account:    account,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		baseURL:    "https://api.continusec.com",
	}
}

// WithBaseURL modifies the client to point to a different base URL other than the
// standard. This is typically only used for integration tests where you don't wish to
// make live calls to the Continusec API server.
//
// The same client object is also returned for the convenience of the caller.
func (self *Client) WithBaseURL(baseURL string) *Client {
	self.baseURL = baseURL
	return self
}

// WithHttpClient modifies the client to use a different http.Client than the default
// (which is to use http.DefaultClient). This is useful for applications hosted on
// Google App Engine that may wish to call: client.WithHttpClient(urlfetch.Client(ctx))
//
// The same client object is also returned for the convenience of the caller.
func (self *Client) WithHttpClient(httpClient *http.Client) *Client {
	self.httpClient = httpClient
	return self
}

// VerifiableMap returns an object representing a Verifiable Map. This function simply
// returns a pointer to an object that can be used to interact with the Map, and won't
// by itself cause any API calls to be generated.
func (self *Client) VerifiableMap(name string) *VerifiableMap {
	return &VerifiableMap{
		client: self,
		path:   "/map/" + name,
	}
}

// VerifiableLog returns an object representing a Verifiable Log. This function simply
// returns a pointer to an object that can be used to interact with the Log, and won't
// by itself cause any API calls to be generated.
func (self *Client) VerifiableLog(name string) *VerifiableLog {
	return &VerifiableLog{
		client: self,
		path:   "/log/" + name,
	}
}

func (self *Client) makeRequest(method, path string, data []byte) ([]byte, http.Header, error) {
	url := fmt.Sprintf("%s/v1/account/%s%s", self.baseURL, self.account, path)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Key "+self.apiKey)
	resp, err := self.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	switch resp.StatusCode {
	case 200:
		return contents, resp.Header, nil
	case 403:
		return nil, nil, ErrNotAuthorized
	case 400:
		return nil, nil, ErrInvalidRange
	case 404:
		return nil, nil, ErrNotFound
	case 409:
		return nil, nil, ErrObjectConflict
	default:
		return nil, nil, ErrInternalError
	}
}
