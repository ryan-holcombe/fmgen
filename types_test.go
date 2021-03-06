package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenStruct(t *testing.T) {
	t.Run("skip struct based on comment fmgen:-", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct fmgen:-",
			},
		}

		assert.True(t, s.Skip())
	})

	t.Run("skip struct based on comment fmgen:exclude", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct fmgen:exclude",
			},
		}

		assert.True(t, s.Skip())
	})

	t.Run("skip struct based on comment fmgen:skip", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct fmgen:skip",
			},
		}

		assert.True(t, s.Skip())
	})

	t.Run("comment with no skip", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: &genComment{
				lineNum: 8,
				value:   "Simple struct",
			},
		}

		assert.False(t, s.Skip())
	})

	t.Run("no comment", func(t *testing.T) {
		s := genStruct{
			lineNum: 9,
			comment: nil,
		}

		assert.False(t, s.Skip())
	})

	t.Run("not in struct includes", func(t *testing.T) {
		*flagStructs = "struct2,struct3"
		s := genStruct{
			name: "struct1",
		}

		assert.True(t, s.Skip())
	})

	t.Run("in struct includes", func(t *testing.T) {
		*flagStructs = "struct1,struct2,struct3"
		s := genStruct{
			name: "struct1",
		}

		assert.False(t, s.Skip())
	})
}
