package yapl

import (
	"github.com/ahsayde/yapl/internal/parser"
	"github.com/ahsayde/yapl/internal/renderer"
)

type Rule struct {
	When      *NestedCondition     `json:"when" yaml:"when"`
	Condition *Condition           `json:"condition" yaml:"condition"`
	Result    *renderer.Renderable `json:"result" yaml:"result"`
}

func (r *Rule) Eval(ctx *Context, input *parser.Node) ([]interface{}, error) {
	var results []interface{}

	if r.When != nil {
		ok, err := r.When.Eval(ctx, input)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, nil
		}
	}

	failures, err := r.Condition.Eval(ctx, input)
	if err != nil {
		return nil, err
	}

	if failures == nil {
		return results, nil
	}

	for i := range failures {
		ctx.Cond = failures[i]
		result, err := r.Result.Render(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (ck *Rule) validate() []parserError {
	var result []parserError

	if ck.Condition == nil {
		result = append(result, parserError{
			msg: "rule must have a condition",
		})
	}

	if ck.Result == nil {
		result = append(result, parserError{
			msg: "rule must have result",
		})
	}

	if ck.When != nil {
		if errs := ck.When.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if errs := ck.Condition.validate(); errs != nil {
		result = append(result, errs...)
	}

	return result
}
