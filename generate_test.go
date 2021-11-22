package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatStructName(t *testing.T) {
	assert.Equal(t, "NewGenStruct", formatStructName("genStruct"))
	assert.Equal(t, "NewTest", formatStructName("Test"))
}

func TestWritePackageFile(t *testing.T) {
	buf := bytes.Buffer{}
	structs := []genStruct{
		{
			name:    "Sample",
			lineNum: 4,
			fields: []genField{
				{name: "ID", typ: "int64", optional: false, skip: true},
				{name: "Name", typ: "string", optional: false, skip: false},
				{name: "Age", typ: "int64", optional: true, skip: false},
				{name: "LastUpdated", typ: "time.Time", optional: true, skip: false},
			},
			comment: &genComment{
				lineNum: 3,
				value:   "Sample simple struct fmgen:omit\n",
			},
		},
	}
	imports := []string{
		`"time"`,
		`"net/url"`,
	}

	expected := `// Code generated by "fmgen". DO NOT EDIT.
package testdata

import (
	"time"
)

// NewSample generated factory method for Sample
func NewSample(Name string, Age *int64, LastUpdated *time.Time) *Sample {
	result := &Sample{
		Name: Name,
	}
	if Age != nil {
		result.Age = *Age
	}
	if LastUpdated != nil {
		result.LastUpdated = *LastUpdated
	}
	return result
}
`

	writePackageFile(&buf, "testdata", imports, structs)
	assert.Equal(t, expected, buf.String())
}

func TestBuildInputParams(t *testing.T) {
	fields := []genField{
		{
			name:     "A",
			typ:      "string",
			optional: false,
			skip:     false,
			ptr:      false,
			array:    false,
		},
		{
			name:     "B",
			typ:      "time.Time",
			optional: true,
			skip:     false,
			ptr:      false,
			array:    false,
		},
		{
			name:     "C",
			typ:      "int64",
			optional: false,
			skip:     true,
			ptr:      false,
			array:    false,
		},
		{
			name:     "D",
			typ:      "int64",
			optional: false,
			skip:     false,
			ptr:      false,
			array:    true,
		},
	}
	result := buildInputParams(fields)
	expected := `A string,B *time.Time,D []int64`
	assert.Equal(t, expected, result)
}

func TestBuildBody(t *testing.T) {
	t.Run("all required", func(t *testing.T) {
		fields := []genField{
			{
				name:     "A",
				typ:      "string",
				optional: false,
				skip:     false,
				ptr:      false,
			},
			{
				name:     "B",
				typ:      "int64",
				optional: false,
				skip:     false,
				ptr:      true,
			},
			{
				name:     "SKIP",
				typ:      "string",
				optional: false,
				skip:     true,
				ptr:      false,
			},
		}
		result := buildBody("Simple", fields)
		expected := `result := &Simple {
A: A,
B: &B,
}
return result
`
		assert.Equal(t, expected, result)
	})

	t.Run("include optional", func(t *testing.T) {
		fields := []genField{
			{
				name:     "A",
				typ:      "string",
				optional: false,
				skip:     false,
				ptr:      false,
			},
			{
				name:     "B",
				typ:      "time.Time",
				optional: true,
				skip:     false,
				ptr:      true,
			},
			{
				name:     "C",
				typ:      "int64",
				optional: true,
				skip:     false,
				ptr:      false,
			},
			{
				name:     "SKIP",
				typ:      "string",
				optional: false,
				skip:     true,
				ptr:      false,
			},
		}
		result := buildBody("Simple", fields)
		expected := `result := &Simple {
A: A,
}
if B != nil {
result.B = B
}
if C != nil {
result.C = *C
}
return result
`
		assert.Equal(t, expected, result)
	})
}
