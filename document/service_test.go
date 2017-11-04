// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package document

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/diegobernardes/flare"
	infraHTTP "github.com/diegobernardes/flare/infra/http/test"
	"github.com/diegobernardes/flare/repository/test"
)

func TestServiceHandleShow(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		status     int
		header     http.Header
		body       []byte
		repository flare.DocumentRepositorier
	}{
		{
			"Not found",
			httptest.NewRequest(http.MethodGet, "http://documents/123", nil),
			http.StatusNotFound,
			http.Header{"Content-Type": []string{"application/json"}},
			load("handleShow.notFound.json"),
			test.NewDocument(),
		},
		{
			"Error at repository",
			httptest.NewRequest(http.MethodGet, "http://documents/123", nil),
			http.StatusInternalServerError,
			http.Header{"Content-Type": []string{"application/json"}},
			load("handleShow.errorRepository.json"),
			test.NewDocument(test.DocumentError(errors.New("error at repository"))),
		},
		{
			"Found",
			httptest.NewRequest(http.MethodGet, "http://documents/456", nil),
			http.StatusOK,
			http.Header{"Content-Type": []string{"application/json"}},
			load("handleShow.found.json"),
			test.NewDocument(
				test.DocumentLoadSliceByteDocument(load("handleShow.foundInput.json")),
				test.DocumentDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
			),
		},
	}

	for _, tt := range tests {
		service, err := NewService(
			ServiceLogger(log.NewNopLogger()),
			ServiceDocumentRepository(tt.repository),
			ServiceResourceRepository(test.NewResource()),
			ServiceGetDocumentId(func(r *http.Request) string {
				return strings.Replace(r.URL.Path, "/", "", -1)
			}),
			ServiceWorker(&Worker{}),
		)
		if err != nil {
			t.Error(errors.Wrap(err, "error during service initialization"))
			t.FailNow()
		}

		t.Run(tt.name, infraHTTP.Handler(tt.status, tt.header, service.HandleShow, tt.req, tt.body))
	}
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name     string
		options  []func(*Service)
		hasError bool
	}{
		{
			"Mising logger",
			[]func(*Service){},
			true,
		},
		{
			"Mising subscription trigger",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
			},
			true,
		},
		{
			"Mising document repository",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
			},
			true,
		},
		{
			"Mising resource repository",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
				ServiceDocumentRepository(test.NewDocument()),
			},
			true,
		},
		{
			"Mising subscription repository",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
				ServiceDocumentRepository(test.NewDocument()),
				ServiceResourceRepository(test.NewResource()),
			},
			true,
		},
		{
			"Mising getDocumentId repository",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
				ServiceDocumentRepository(test.NewDocument()),
				ServiceResourceRepository(test.NewResource()),
			},
			true,
		},
		{
			"Mising getDocumentURI repository",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
				ServiceDocumentRepository(test.NewDocument()),
				ServiceResourceRepository(test.NewResource()),
				ServiceGetDocumentId(func(*http.Request) string { return "" }),
			},
			true,
		},
		{
			"Success",
			[]func(*Service){
				ServiceLogger(log.NewNopLogger()),
				ServiceDocumentRepository(test.NewDocument()),
				ServiceResourceRepository(test.NewResource()),
				ServiceGetDocumentId(func(*http.Request) string { return "" }),
				ServiceWorker(&Worker{}),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.options...)
			if tt.hasError != (err != nil) {
				t.Errorf("NewService invalid result, want '%v', got '%v'", tt.hasError, err)
			}
		})
	}
}
