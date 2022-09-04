package operator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	floatType = reflect.TypeOf(float64(0))
)

type Operator interface {
	Eval(val, expected interface{}) (bool, error)
}

type EqualOperator struct{}

func (_ *EqualOperator) Eval(val, operatorVal interface{}) (bool, error) {
	return reflect.DeepEqual(val, operatorVal), nil
}

type HasSuffixOperator struct{}

func (_ *HasSuffixOperator) Eval(val, operatorVal interface{}) (bool, error) {
	suffix, ok := operatorVal.(string)
	if !ok {
		return false, fmt.Errorf("suffix must be of type string, found %s", reflect.TypeOf(operatorVal).String())
	}

	value, ok := val.(string)
	if !ok {
		return false, nil
	}

	return strings.HasSuffix(value, suffix), nil
}

type HasPrefixOperator struct{}

func (_ *HasPrefixOperator) Eval(val, operatorVal interface{}) (bool, error) {
	prefix, ok := operatorVal.(string)
	if !ok {
		return false, fmt.Errorf("prefix must be of type string, found %s", reflect.TypeOf(operatorVal).String())
	}

	value, ok := val.(string)
	if !ok {
		return false, nil
	}

	return strings.HasPrefix(value, prefix), nil
}

type RegexOperator struct{}

func (_ *RegexOperator) Eval(val, operatorVal interface{}) (bool, error) {
	expr, ok := operatorVal.(string)
	if !ok {
		return false, fmt.Errorf("pattern must be of type string, found %s", reflect.TypeOf(operatorVal).String())
	}

	regex, err := regexp.Compile(expr)
	if err != nil {
		return false, fmt.Errorf("invalid regex expression")
	}

	value, ok := val.(string)
	if !ok {
		return false, nil
	}

	return regex.MatchString(value), nil
}

type MaxValueOperator struct{}

func (_ *MaxValueOperator) Eval(val, operatorVal interface{}) (bool, error) {
	max, err := convertToNumber(operatorVal)
	if err != nil {
		return false, fmt.Errorf("max value must be of type number, but found %s", reflect.TypeOf(operatorVal).String())
	}

	value, err := convertToNumber(val)
	if err != nil {
		return false, fmt.Errorf("value must be of type number, but found %s", reflect.TypeOf(value).String())
	}

	return value <= max, nil
}

type MinValueOperator struct{}

func (_ *MinValueOperator) Eval(val, operatorVal interface{}) (bool, error) {
	min, err := convertToNumber(operatorVal)
	if err != nil {
		return false, fmt.Errorf("min value must be of type number, but found %s", reflect.TypeOf(operatorVal).String())
	}

	value, err := convertToNumber(val)
	if err != nil {
		return false, fmt.Errorf("value must be of type number, but found %s", reflect.TypeOf(value).String())
	}

	return value >= min, nil
}

type InOperator struct{}

func (_ *InOperator) Eval(val, opval interface{}) (bool, error) {
	if reflect.TypeOf(opval).Kind() != reflect.Slice {
		return false, fmt.Errorf("invalid operator value, expected array found %s", reflect.TypeOf(opval).String())
	}

	arr, _ := opval.([]interface{})

	for i := range arr {
		if reflect.DeepEqual(val, arr[i]) {
			return true, nil
		}
	}
	return false, nil
}

type ContainsOperator struct{}

func (_ *ContainsOperator) Eval(val, opval interface{}) (bool, error) {
	cval, ok := val.([]interface{})
	if !ok {
		return false, fmt.Errorf("value must be of type array, found %s", reflect.TypeOf(cval).String())
	}

	for i := range cval {
		if reflect.DeepEqual(opval, cval[i]) {
			return true, nil
		}
	}

	return false, nil
}

type LengthOperator struct{}

func (_ *LengthOperator) Eval(val, opval interface{}) (bool, error) {
	cval, ok := val.([]interface{})
	if !ok {
		return false, fmt.Errorf("value must be of type array, found %s", reflect.TypeOf(cval).String())
	}

	length, ok := opval.(int)
	if !ok {
		return false, fmt.Errorf("length value must be of type int, found %s", reflect.TypeOf(opval).String())
	}

	return len(cval) == length, nil
}

type MaxLengthOperator struct{}

func (_ *MaxLengthOperator) Eval(val, opval interface{}) (bool, error) {
	cval, ok := val.([]interface{})
	if !ok {
		return false, fmt.Errorf("value must be of type array, found %s", reflect.TypeOf(cval).String())
	}

	maxLength, ok := opval.(int)
	if !ok {
		return false, fmt.Errorf("max length value must be of type int, found %s", reflect.TypeOf(opval).String())
	}

	return len(cval) <= maxLength, nil
}

type MinLengthOperator struct{}

func (_ *MinLengthOperator) Eval(val, opval interface{}) (bool, error) {
	cval, ok := val.([]interface{})
	if !ok {
		return false, fmt.Errorf("value must be of type array, found %s", reflect.TypeOf(cval).String())
	}

	minLength, ok := opval.(int)
	if !ok {
		return false, fmt.Errorf("min length value must be of type int, found %s", reflect.TypeOf(opval).String())
	}

	return len(cval) >= minLength, nil
}

func GetOperator(name string) Operator {
	switch name {
	case "equal", "eq":
		return &EqualOperator{}
	case "hasPrefix":
		return &HasPrefixOperator{}
	case "hasSuffix":
		return &HasSuffixOperator{}
	case "regex":
		return &RegexOperator{}
	case "max", "maxValue":
		return &MaxValueOperator{}
	case "min", "minValue":
		return &MinValueOperator{}
	case "in":
		return &InOperator{}
	case "contains":
		return &ContainsOperator{}
	case "len", "length":
		return &LengthOperator{}
	case "maxlen", "maxLength":
		return &MaxLengthOperator{}
	case "minlen", "minLength":
		return &MinLengthOperator{}
	default:
		return nil
	}
}

func convertToNumber(value interface{}) (float64, error) {
	val := reflect.ValueOf(value)
	if val.Type().ConvertibleTo(floatType) {
		return val.Convert(floatType).Float(), nil
	}
	return 0, fmt.Errorf("value must be of type number, found %s", val.Type())
}
