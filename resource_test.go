// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResourceChangeValid(t *testing.T) {
	Convey("Given a list of valid resource changes", t, func() {
		tests := []ResourceChange{
			{Field: "updatedAt", Kind: ResourceChangeDate, DateFormat: "2006-01-02"},
			{Field: "revision", Kind: ResourceChangeInteger},
		}

		Convey("The validation should not return a error", func() {
			for _, tt := range tests {
				So(tt.Valid(), ShouldBeNil)
			}
		})
	})

	Convey("Given a list of invalid resource changes", t, func() {
		tests := []struct {
			title string
			rc    ResourceChange
		}{
			{
				"Should be missing the field",
				ResourceChange{},
			},
			{
				"Should be missing the kind",
				ResourceChange{Field: "updatedAt"},
			},
			{
				"Should be missing the format",
				ResourceChange{Field: "updatedAt", Kind: ResourceChangeDate},
			},
		}

		for _, tt := range tests {
			Convey(tt.title, func() {
				So(tt.rc.Valid(), ShouldNotBeNil)
			})
		}
	})
}

func TestResourceWildcardReplace(t *testing.T) {
	Convey("Given a list of valid wildcards to be replaced", t, func() {
		tests := []struct {
			resource   Resource
			id         string
			revision   interface{}
			rawContent []string
			expected   []string
			hasErr     bool
		}{
			{
				Resource{Path: "/resource/{id}"},
				"/resource/123",
				nil,
				[]string{"{id}", `{"id":"{id}"}`},
				[]string{"123", `{"id":"123"}`},
				false,
			},
			{
				Resource{Path: "/resource/{id}"},
				"/resource/123",
				1,
				[]string{"{revision}"},
				[]string{"1"},
				false,
			},
			{
				Resource{Path: "/resource/{id}"},
				"/resource/123",
				"sample",
				[]string{"{revision}"},
				[]string{"sample"},
				false,
			},
			{
				Resource{Path: "/resource/{id}"},
				"/resource/123",
				time.Date(2011, time.January, 0, 0, 0, 0, 0, time.UTC),
				[]string{"{revision}"},
				[]string{"2010-12-31T00:00:00Z"},
				false,
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				fn, err := tt.resource.WildcardReplace(tt.id, tt.revision)
				So(err, ShouldBeNil)

				for i, value := range tt.rawContent {
					tt.rawContent[i] = fn(value)
				}

				So(tt.rawContent, ShouldResemble, tt.expected)
			}
		})
	})

	Convey("Given a list of invalid wildcards to be replaced", t, func() {
		tests := []struct {
			resource   Resource
			id         string
			revision   interface{}
			rawContent []string
			expected   []string
			hasErr     bool
		}{
			{
				Resource{},
				"%zzzzz",
				nil,
				nil,
				nil,
				true,
			},
		}

		Convey("It's expected to have a error", func() {
			for _, tt := range tests {
				_, err := tt.resource.WildcardReplace(tt.id, tt.revision)
				So(err, ShouldNotBeNil)
			}
		})
	})
}
