package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserFind(t *testing.T) {
	cases := []struct {
		name     string
		input    interface{}
		key      string
		expected interface{}
	}{
		{
			name:     "scalar-string",
			input:    map[string]interface{}{"my-key": "my-value"},
			key:      "my-key",
			expected: "my-value",
		},
		{
			name:     "scalar-int",
			input:    map[string]interface{}{"my-key": 1},
			key:      "my-key",
			expected: 1,
		},
		{
			name:     "scalar-float",
			input:    map[string]interface{}{"my-key": 1.5},
			key:      "my-key",
			expected: 1.5,
		},
		{
			name:     "scalar-bool",
			input:    map[string]interface{}{"my-key": true},
			key:      "my-key",
			expected: true,
		},
		{
			name:     "scalar-not-found",
			input:    map[string]interface{}{"my-key": true},
			key:      "my-key-2",
			expected: (*Node)(nil),
		},
		{
			name: "map-scalar-string",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": "my-value",
				},
			},
			key:      "parent.child",
			expected: "my-value",
		},
		{
			name: "map-scalar-int",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": 1,
				},
			},
			key:      "parent.child",
			expected: 1,
		},
		{
			name: "map-scalar-float",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": 1.5,
				},
			},
			key:      "parent.child",
			expected: 1.5,
		},
		{
			name: "map-scalar-bool",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": true,
				},
			},
			key:      "parent.child",
			expected: true,
		},
		{
			name: "map-map",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": map[string]interface{}{
						"child-2": "my-value",
					},
				},
			},
			key:      "parent.child.child-2",
			expected: "my-value",
		},
		{
			name: "map-array",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": []interface{}{
						"my-value-1",
						"my-value-2",
					},
				},
			},
			key:      "parent.child",
			expected: []interface{}{"my-value-1", "my-value-2"},
		},
		{
			name: "map-not-found",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": map[string]interface{}{},
				},
			},
			key:      "parent.child.not-found",
			expected: (*Node)(nil),
		},
	}

	for i, test := range cases {
		node := Parse(test.input)
		field, err := node.Find(test.key)
		if err != nil {
			t.Error(err)
		}
		if field != nil {
			assert.Equal(t, test.expected, field.Value, "testcase #%d failed", i)
		} else {
			assert.Equal(t, test.expected, field, "testcase #%d failed", i)
		}
	}
}

func TestParserFindAll(t *testing.T) {
	cases := []struct {
		name     string
		input    interface{}
		key      string
		expected []interface{}
		err      bool
	}{
		{
			name: "with-index-0",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "user1",
					},
					map[string]interface{}{
						"name": "user2",
					},
				},
			},
			key:      "users[0].name",
			expected: []interface{}{"user1"},
		},
		{
			name: "with-index-1",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "user1",
					},
					map[string]interface{}{
						"name": "user2",
					},
				},
			},
			key:      "users[1].name",
			expected: []interface{}{"user2"},
		},
		{
			name: "with-asterisks",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "user1",
					},
					map[string]interface{}{
						"name": "user2",
					},
				},
			},
			key:      "users[*].name",
			expected: []interface{}{"user1", "user2"},
		},
		{
			name: "with-invalid-key",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "user1",
					},
					map[string]interface{}{
						"name": "user2",
					},
				},
			},
			key: "users[invalid].name",
			err: true,
		},
		{
			name: "with-out-of-index",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "user1",
					},
					map[string]interface{}{
						"name": "user2",
					},
				},
			},
			key: "users[3].name",
			err: true,
		},
	}

	for i, test := range cases {
		node := Parse(test.input)
		fields, err := node.FindAll(test.key)
		if err != nil {
			if test.err {
				continue
			}
			t.Error(err)
		}

		assert.Equal(t, len(test.expected), len(fields), "testcase #%d failed", i)

		for j := range fields {
			assert.Equal(t, test.expected[j], fields[j].Value, "testcase #%d failed", i)
		}
	}
}
