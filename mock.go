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

package continusec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type ProxyAndRecordHandler struct {
	Host                  string
	InHeaders, OutHeaders []string
	Dir                   string
	Sequence              int
	FailOnMissing         bool
}

type SavedResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type SavedRequest struct {
	URL     string
	Method  string
	Headers map[string][]string
	Body    []byte
}

func (us *SavedRequest) Equals(them *SavedRequest) bool {
	return reflect.DeepEqual(us, them)
}

type SavedPair struct {
	Request  *SavedRequest
	Response *SavedResponse
}

func FilePathForSeq(path string, seq int) string {
	return filepath.Join(path, fmt.Sprintf("%04d.response", seq))
}

func (self *SavedPair) Write(path string, seq int) error {
	fi, err := os.Create(FilePathForSeq(path, seq))
	if err != nil {
		return err
	}
	err = json.NewEncoder(fi).Encode(self)
	if err != nil {
		return err
	}
	fi.Close()
	return nil
}

func LoadSavedIfThere(path string, seq int) (*SavedPair, error) {
	fi, err := os.Open(FilePathForSeq(path, seq))
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	var rv SavedPair
	err = json.NewDecoder(fi).Decode(&rv)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

func saveRequest(r *http.Request, altHost string, headerFilter []string) (*SavedRequest, error) {
	url := r.URL.String()
	url = altHost + url[strings.Index(url, "/"):]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string][]string)
	for _, h := range headerFilter {
		canon := http.CanonicalHeaderKey(h)
		z, ok := r.Header[canon]
		if ok {
			headers[canon] = z
		}
	}

	return &SavedRequest{
		Method:  r.Method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}, nil
}

func saveResponse(resp *http.Response, headerFilter []string) (*SavedResponse, error) {
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	headers := make(map[string][]string)
	for _, h := range headerFilter {
		canon := http.CanonicalHeaderKey(h)
		z, ok := resp.Header[canon]
		if ok {
			headers[canon] = z
		}
	}

	return &SavedResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       contents,
	}, nil
}

func writeResponse(saved *SavedResponse, w http.ResponseWriter) {
	for k, vs := range saved.Headers {
		w.Header()[k] = vs
	}
	w.WriteHeader(saved.StatusCode)
	w.Write(saved.Body)
}

func sendSavedRequest(savedReq *SavedRequest, headerIn, headerOut []string) (*SavedResponse, error) {
	req, err := http.NewRequest(savedReq.Method, savedReq.URL, bytes.NewReader(savedReq.Body))
	if err != nil {
		return nil, err
	}

	for _, h := range headerIn {
		canon := http.CanonicalHeaderKey(h)
		req.Header[canon] = savedReq.Headers[canon]
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return saveResponse(resp, headerOut)
}

func (self *ProxyAndRecordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	canonReq, err := saveRequest(r, self.Host, self.InHeaders)
	if err != nil {
		fmt.Println(self.Sequence, "Error saving request:", err)
		return
	}
	savedPair, err := LoadSavedIfThere(self.Dir, self.Sequence)
	if err != nil {
		if self.FailOnMissing {
			fmt.Println(self.Sequence, "Error loading response:", err)
			return
		} else {
			fmt.Println(self.Sequence, "Fetching", canonReq.URL)
			sr, err := sendSavedRequest(canonReq, self.InHeaders, self.OutHeaders)
			if err != nil {
				fmt.Println(self.Sequence, "Error receiving response:", err)
				return
			}
			savedPair = &SavedPair{
				Request:  canonReq,
				Response: sr,
			}
			err = savedPair.Write(self.Dir, self.Sequence)
			if err != nil {
				fmt.Println(self.Sequence, "Error saving response:", err)
				return
			}
		}
	} else {
		fmt.Println(self.Sequence, "From cache", canonReq.URL)
	}
	if !savedPair.Request.Equals(canonReq) {
		fmt.Println(self.Sequence, "Bad request, got", canonReq, "wanted", savedPair.Request)

		return
	}

	writeResponse(savedPair.Response, w)
	self.Sequence++
}

func RunMockServer(hostport string, pr *ProxyAndRecordHandler) {
	http.ListenAndServe(hostport, pr)
}

/*
func main() {
	http.ListenAndServe(":8080", &ProxyAndRecordHandler{
		Host:          "https://api.continusec.com",
		InHeaders:     []string{"Authorization"},
		OutHeaders:    []string{"Content-Type", "X-Verified-TreeSize", "X-Verified-Proof"},
		Dir:           "responses",
		FailOnMissing: false,
	})
}
*/
