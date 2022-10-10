package renderer

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	processors = map[string]interface{}{
		// string processors
		"split":      strings.Split,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      strings.Title,
		"join":       strings.Join,
		"trim":       strings.Trim,
		"trimLeft":   strings.TrimLeft,
		"trimRigh":   strings.TrimRight,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"replace":    strings.Replace,
		"replaceAll": strings.ReplaceAll,
		// math processors
		"round": math.Round,
		"ceil":  math.Ceil,
		"abs":   math.Abs,
		"floor": math.Floor,
		"max":   math.Max,
		"min":   math.Min,
		// date & time processors
		"date":    date,
		"now":     time.Now,
		"year":    year,
		"month":   month,
		"weekday": weekday,
		"day":     day,
		"hour":    hour,
		"minute":  minute,
		"second":  second,
		// types processors
		"bool": boolean,
		// converters
		"toYaml": toYaml,
		// conditions
		"ternary": ternary,
	}
)

func date(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func year(t time.Time) int {
	return t.Year()
}

func month(t time.Time) int {
	return int(t.Month())
}

func weekday(t time.Time) int {
	return int(t.Weekday())
}

func day(t time.Time) int {
	return t.Day()
}

func hour(t time.Time) int {
	return t.Hour()
}

func minute(t time.Time) int {
	return t.Minute()
}

func second(t time.Time) int {
	return t.Second()
}

func boolean(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %s to boolean", s)
	}
}

func toYaml(in interface{}) string {
	raw, err := yaml.Marshal(in)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(raw), "\n")
}

func ternary(v1, v2 interface{}, cond bool) interface{} {
	if cond {
		return v1
	}
	return v2
}
