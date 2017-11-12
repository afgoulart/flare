// Copyright 2017 Diego Bernardes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flare

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDocumentValid(t *testing.T) {
	Convey("Given a list of valid documents", t, func() {
		tests := []Document{
			{
				Id:               "1",
				ChangeFieldValue: 1,
				Resource: Resource{
					Change: ResourceChange{Field: "revision", Kind: ResourceChangeInteger},
				},
			},
			{
				Id:               "1",
				ChangeFieldValue: "2006-01-02",
				Resource: Resource{
					Change: ResourceChange{
						Field:      "revision",
						Kind:       ResourceChangeDate,
						DateFormat: "2006-01-02",
					},
				},
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				So(tt.Valid(), ShouldBeNil)
			}
		})
	})

	Convey("Given a list of invalid documents", t, func() {
		tests := []struct {
			title string
			doc   Document
		}{
			{
				"Should have a invalid id",
				Document{},
			},
			{
				"Should have a invalid changeFieldValue",
				Document{Id: "1"},
			},
			{
				"Should have a invalid resource.Change",
				Document{
					Id:               "1",
					ChangeFieldValue: 1,
					Resource:         Resource{Change: ResourceChange{}},
				},
			},
			{
				"Shoud have a invalid changeFieldValue 1",
				Document{
					Id:               "1",
					ChangeFieldValue: 1,
					Resource: Resource{
						Change: ResourceChange{
							Field: "revision",
							Kind:  ResourceChangeString,
						},
					},
				},
			},
			{
				"Should have a invalid changeFieldValue 2",
				Document{
					Id:               "1",
					ChangeFieldValue: 1,
					Resource: Resource{
						Change: ResourceChange{
							Field:      "revision",
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02",
						},
					},
				},
			},
			{
				"Should have a invalid date at changeFieldValue",
				Document{
					Id:               "1",
					ChangeFieldValue: "sample",
					Resource: Resource{
						Change: ResourceChange{
							Field:      "revision",
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02",
						},
					},
				},
			},
		}

		for _, tt := range tests {
			Convey(tt.title, func() {
				So(tt.doc.Valid(), ShouldNotBeNil)
			})
		}
	})
}

func TestDocumentNewer(t *testing.T) {
	Convey("Given a list of newer documents", t, func() {
		tests := []struct {
			reference *Document
			target    Document
		}{
			{
				nil,
				Document{},
			},
			{
				&Document{
					ChangeFieldValue: float64(0),
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeInteger,
						},
					},
				},
				Document{
					ChangeFieldValue: float64(1),
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeInteger,
						},
					},
				},
			},
			{
				&Document{
					ChangeFieldValue: "a",
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeString,
						},
					},
				},
				Document{
					ChangeFieldValue: "b",
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeString,
						},
					},
				},
			},
			{
				&Document{
					ChangeFieldValue: "2007-09-07T07:08:08.008Z",
					Resource: Resource{
						Change: ResourceChange{
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02T15:04:05Z07:00",
						},
					},
				},
				Document{
					ChangeFieldValue: "2008-09-07T07:08:08.008Z",
					Resource: Resource{
						Change: ResourceChange{
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.target.Newer(tt.reference)
				So(err, ShouldBeNil)
				So(newer, ShouldBeTrue)
			}
		})
	})

	Convey("Given a list of older documents", t, func() {
		tests := []struct {
			reference *Document
			target    Document
		}{
			{
				&Document{
					ChangeFieldValue: float64(1),
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeInteger,
						},
					},
				},
				Document{
					ChangeFieldValue: float64(0),
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeInteger,
						},
					},
				},
			},
			{
				&Document{
					ChangeFieldValue: "b",
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeString,
						},
					},
				},
				Document{
					ChangeFieldValue: "a",
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeString,
						},
					},
				},
			},
			{
				&Document{
					ChangeFieldValue: "2008-09-07T07:08:08.008Z",
					Resource: Resource{
						Change: ResourceChange{
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02T15:04:05Z07:00",
						},
					},
				},
				Document{
					ChangeFieldValue: "2007-09-07T07:08:08.008Z",
					Resource: Resource{
						Change: ResourceChange{
							Kind:       ResourceChangeDate,
							DateFormat: "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.target.Newer(tt.reference)
				So(err, ShouldBeNil)
				So(newer, ShouldBeFalse)
			}
		})
	})

	Convey("Given a list of invalid documents", t, func() {
		tests := []struct {
			reference *Document
			target    Document
		}{
			{
				&Document{},
				Document{},
			},
			{
				&Document{},
				Document{
					Resource: Resource{
						Change: ResourceChange{
							Kind: "sample",
						},
					},
				},
			},
			{
				&Document{},
				Document{
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeDate,
						},
					},
				},
			},
			{
				&Document{},
				Document{
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeInteger,
						},
					},
				},
			},
			{
				&Document{},
				Document{
					Resource: Resource{
						Change: ResourceChange{
							Kind: ResourceChangeString,
						},
					},
				},
			},
		}
		Convey("The output should be valid", func() {
			for _, tt := range tests {
				_, err := tt.target.Newer(tt.reference)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestDocumentNewerDate(t *testing.T) {
	Convey("Given a list of newer documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{
					ChangeFieldValue: "2007-01-02",
					Resource: Resource{
						Change: ResourceChange{
							DateFormat: "2006-01-02",
						},
					},
				},
				"2006-01-02",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerDate(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeTrue)
			}
		})
	})

	Convey("Given a list of older documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{
					ChangeFieldValue: "2006-01-02",
					Resource: Resource{
						Change: ResourceChange{
							DateFormat: "2006-01-02",
						},
					},
				},
				"2006-01-02",
			},
			{
				Document{
					ChangeFieldValue: "2005-01-02",
					Resource: Resource{
						Change: ResourceChange{
							DateFormat: "2006-01-02",
						},
					},
				},
				"2006-01-02",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerDate(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeFalse)
			}
		})
	})

	Convey("Given a list of invalid documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{},
				nil,
			},
			{
				Document{ChangeFieldValue: "2006-01-02"},
				nil,
			},
			{
				Document{ChangeFieldValue: "2006-01-02"},
				"2006-01-02",
			},
			{
				Document{
					ChangeFieldValue: "2006-01-02",
					Resource: Resource{
						Change: ResourceChange{
							DateFormat: "2006-01-02",
						},
					},
				},
				"",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				_, err := tt.document.newerDate(tt.value)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestDocumentNewerInteger(t *testing.T) {
	Convey("Given a list of newer documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{ChangeFieldValue: float64(2)},
				float64(1),
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerInteger(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeTrue)
			}
		})
	})

	Convey("Given a list of older documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{ChangeFieldValue: float64(1)},
				float64(1),
			},
			{
				Document{ChangeFieldValue: float64(0)},
				float64(1),
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerInteger(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeFalse)
			}
		})
	})

	Convey("Given a list of invalid documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{ChangeFieldValue: 1},
				nil,
			},
			{
				Document{ChangeFieldValue: float64(1)},
				"sample",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				_, err := tt.document.newerInteger(tt.value)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestDocumentNewerString(t *testing.T) {
	Convey("Given a list of newer documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{ChangeFieldValue: "sample1"},
				"sample0",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerString(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeTrue)
			}
		})
	})

	Convey("Given a list of older documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{ChangeFieldValue: "sample"},
				"sample",
			},
			{
				Document{ChangeFieldValue: "sample0"},
				"sample1",
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				newer, err := tt.document.newerString(tt.value)
				So(err, ShouldBeNil)
				So(newer, ShouldBeFalse)
			}
		})
	})

	Convey("Given a list of invalid documents", t, func() {
		tests := []struct {
			document Document
			value    interface{}
		}{
			{
				Document{},
				nil,
			},
			{
				Document{ChangeFieldValue: 1},
				nil,
			},
			{
				Document{ChangeFieldValue: "sample"},
				nil,
			},
			{
				Document{ChangeFieldValue: "sample"},
				1,
			},
		}

		Convey("The output should be valid", func() {
			for _, tt := range tests {
				_, err := tt.document.newerString(tt.value)
				So(err, ShouldNotBeNil)
			}
		})
	})
}
