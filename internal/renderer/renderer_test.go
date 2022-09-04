package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		name  string
		expr  interface{}
		ctx   map[string]interface{}
		value interface{}
	}{
		{
			name: "string",
			expr: "${ .param1 }",
			ctx: map[string]interface{}{
				"param1": "test1",
			},
			value: "test1",
		},
		{
			name: "boolean",
			expr: "${ .param1 }",
			ctx: map[string]interface{}{
				"param1": true,
			},
			value: true,
		},
		{
			name: "integer",
			expr: "${ .param1 }",
			ctx: map[string]interface{}{
				"param1": 1,
			},
			value: 1,
		},
		{
			name: "float",
			expr: "${ .param1 }",
			ctx: map[string]interface{}{
				"param1": 1.5,
			},
			value: 1.5,
		},
		{
			name: "array",
			expr: []interface{}{
				"${.param1}",
				"${.param2}",
			},
			ctx: map[string]interface{}{
				"param1": "test1",
				"param2": "test2",
			},
			value: []interface{}{"test1", "test2"},
		},
		{
			name: "object",
			expr: map[string]interface{}{
				"key1": "${.param1}",
				"key2": "${.param2}",
			},
			ctx: map[string]interface{}{
				"param1": "test1",
				"param2": "test2",
			},
			value: map[string]interface{}{
				"key1": "test1",
				"key2": "test2",
			},
		},
	}

	for i, test := range cases {
		r, err := newRenderable(test.expr)
		if err != nil {
			t.Fatal(err)
		}

		val, err := r.Render(test.ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, test.value, val, "testcase #%d failed", i)
	}
}
