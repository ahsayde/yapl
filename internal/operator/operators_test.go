package operator

import "testing"

func TestEqualOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 1,
			inValue: 1,
			result:  true,
		},
		{
			opValue: 1,
			inValue: 2,
			result:  false,
		},
		{
			opValue: true,
			inValue: true,
			result:  true,
		},
		{
			opValue: true,
			inValue: false,
			result:  false,
		},
		{
			opValue: 1.5,
			inValue: 1.5,
			result:  true,
		},
		{
			opValue: 1.5,
			inValue: 1.6,
			result:  false,
		},
		{
			opValue: "test",
			inValue: "test",
			result:  true,
		},
		{
			opValue: "test",
			inValue: "test2",
			result:  false,
		},
		{
			opValue: []int{1, 2, 3},
			inValue: []int{1, 2, 3},
			result:  true,
		},
		{
			opValue: []int{1, 2, 3},
			inValue: []int{1, 2},
			result:  false,
		},
		{
			opValue: map[string]interface{}{"a": 1, "b": 2},
			inValue: map[string]interface{}{"a": 1, "b": 2},
			result:  true,
		},
		{
			opValue: map[string]interface{}{"a": 1, "b": 2, "c": 3},
			inValue: map[string]interface{}{"a": 1, "b": 2},
			result:  false,
		},
	}

	operator := GetOperator("eq")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestHasSuffixOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: "bbb",
			inValue: "aaa bbb",
			result:  true,
		},
		{
			opValue: "aaa",
			inValue: "aaa bbb",
			result:  false,
		},
		{
			opValue: 1,
			inValue: "aaa bbb",
			result:  false,
			err:     true,
		},
		{
			opValue: "aaa",
			inValue: 1,
			result:  false,
			err:     true,
		},
	}

	operator := GetOperator("hasSuffix")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestHasPrefixOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: "aaa",
			inValue: "aaa bbb",
			result:  true,
		},
		{
			opValue: "bbb",
			inValue: "aaa bbb",
			result:  false,
		},
		{
			opValue: 1,
			inValue: "aaa bbb",
			result:  false,
			err:     true,
		},
		{
			opValue: "aaa",
			inValue: 1,
			result:  false,
			err:     true,
		},
	}

	operator := GetOperator("hasPrefix")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestMaxValueOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 10,
			inValue: 1,
			result:  true,
		},
		{
			opValue: 10,
			inValue: 10,
			result:  true,
		},
		{
			opValue: 10,
			inValue: 11,
			result:  false,
		},
		{
			opValue: "10",
			inValue: 11,
			err:     true,
		},
		{
			opValue: 10,
			inValue: "a",
			err:     true,
		},
	}

	operator := GetOperator("max")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestMinValueOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 1,
			inValue: 1,
			result:  true,
		},
		{
			opValue: 1,
			inValue: 2,
			result:  true,
		},
		{
			opValue: 1,
			inValue: 0,
			result:  false,
		},
		{
			opValue: "10",
			inValue: 1,
			err:     true,
		},
		{
			opValue: 10,
			inValue: "a",
			err:     true,
		},
	}

	operator := GetOperator("min")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestInOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: []interface{}{1, 2, 3},
			inValue: 1,
			result:  true,
		},
		{
			opValue: []interface{}{1.1, 2.2, 3.3},
			inValue: 1.1,
			result:  true,
		},
		{
			opValue: []interface{}{"a", "b", "c"},
			inValue: "a",
			result:  true,
		},
		{
			opValue: []interface{}{true, false},
			inValue: true,
			result:  true,
		},
		{
			opValue: []interface{}{
				[]interface{}{1, 2},
				[]interface{}{3, 4},
				[]interface{}{5, 6},
			},
			inValue: []interface{}{1, 2},
			result:  true,
		},
		{
			opValue: []interface{}{
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"a": 3, "b": 4},
				map[string]interface{}{"a": 5, "b": 6},
			},
			inValue: map[string]interface{}{"a": 1, "b": 2},
			result:  true,
		},
		{
			opValue: []interface{}{1, 2, 3},
			inValue: 4,
			result:  false,
		},
		{
			opValue: 1,
			inValue: 1,
			err:     true,
		},
	}

	operator := GetOperator("in")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestContainsOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 1,
			inValue: []interface{}{1, 2, 3},
			result:  true,
		},
		{
			opValue: 1.1,
			inValue: []interface{}{1.1, 2.2, 3.3},
			result:  true,
		},
		{
			opValue: "a",
			inValue: []interface{}{"a", "b", "c"},
			result:  true,
		},
		{
			opValue: true,
			inValue: []interface{}{true, false},
			result:  true,
		},
		{
			opValue: []interface{}{1, 2},
			inValue: []interface{}{
				[]interface{}{1, 2},
				[]interface{}{3, 4},
				[]interface{}{5, 6},
			},
			result: true,
		},
		{
			opValue: map[string]interface{}{"a": 1, "b": 2},
			inValue: []interface{}{
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"a": 3, "b": 4},
				map[string]interface{}{"a": 5, "b": 6},
			},
			result: true,
		},
		{
			opValue: 4,
			inValue: []interface{}{1, 2, 3},
			result:  false,
		},
		{
			opValue: 1,
			inValue: 1,
			err:     true,
		},
	}

	operator := GetOperator("contains")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestLengthOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 0,
			inValue: []interface{}{},
			result:  true,
		},
		{
			opValue: 1,
			inValue: []interface{}{1},
			result:  true,
		},
		{
			opValue: 2,
			inValue: []interface{}{1.1, 1.2},
			result:  true,
		},
		{
			opValue: 3,
			inValue: []interface{}{"a", "b", "c"},
			result:  true,
		},
		{
			opValue: 4,
			inValue: []interface{}{"a", "b", "c"},
			result:  false,
		},
		{
			inValue: 1,
			opValue: 4,
			err:     true,
		},
		{
			opValue: "aaa",
			inValue: []interface{}{},
			err:     true,
		},
	}

	operator := GetOperator("len")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestMaxLengthOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 1,
			inValue: []interface{}{},
			result:  true,
		},
		{
			opValue: 1,
			inValue: []interface{}{1},
			result:  true,
		},
		{
			opValue: 2,
			inValue: []interface{}{1, 2, 3},
			result:  false,
		},
		{
			inValue: 1,
			opValue: 1,
			err:     true,
		},
		{
			opValue: "aaa",
			inValue: []interface{}{},
			err:     true,
		},
	}

	operator := GetOperator("maxLength")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestMinLengthOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: 0,
			inValue: []interface{}{},
			result:  true,
		},
		{
			opValue: 1,
			inValue: []interface{}{1},
			result:  true,
		},
		{
			opValue: 1,
			inValue: []interface{}{1, 2},
			result:  true,
		},
		{
			opValue: 2,
			inValue: []interface{}{1},
			result:  false,
		},
		{
			inValue: 1,
			opValue: 1,
			err:     true,
		},
		{
			opValue: "aaa",
			inValue: []interface{}{},
			err:     true,
		},
	}

	operator := GetOperator("minLength")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}

func TestRegexOperator(t *testing.T) {
	cases := []struct {
		opValue interface{}
		inValue interface{}
		result  bool
		err     bool
	}{
		{
			opValue: "^test$",
			inValue: "test",
			result:  true,
		},
		{
			opValue: "^test$",
			inValue: "aaa",
			result:  false,
		},
		{
			opValue: "?",
			inValue: "test",
			err:     true,
		},
		{
			opValue: "^test$",
			inValue: 1,
			err:     true,
		},
		{
			opValue: 1,
			inValue: "test",
			err:     true,
		},
	}

	operator := GetOperator("regex")
	for _, test := range cases {
		result, err := operator.Eval(test.inValue, test.opValue)
		if err != nil {
			if !test.err {
				t.Error(err)
			}
		} else {
			if result != test.result {
				t.Errorf("failed")
			}
		}
	}
}
