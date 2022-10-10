package yapl

import (
	"fmt"

	"github.com/ahsayde/yapl/internal/operator"
	"github.com/ahsayde/yapl/internal/parser"
	"github.com/ahsayde/yapl/internal/renderer"
)

type renderedCondition struct {
	Field    string
	Operator string
	Value    interface{}
	Expr     *renderer.Renderable
}

type ConditionResult struct {
	Field    *parser.Node
	Operator string
	Expr     interface{}
	Value    interface{}
}

type Condition struct {
	Field    *renderer.Renderable `json:"field,omitempty" yaml:"field,omitempty"`
	Expr     *renderer.Renderable `json:"expr,omitempty" yaml:"expr,omitempty"`
	Operator *renderer.Renderable `json:"operator" yaml:"operator"`
	Value    *renderer.Renderable `json:"value" yaml:"value"`
}

func (c *Condition) Eval(ctx *Context, input *parser.Node) ([]ConditionResult, error) {
	condition, err := c.render(ctx)
	if err != nil {
		return nil, err
	}

	operator := operator.GetOperator(condition.Operator)
	if operator == nil {
		return nil, RuntimeError{
			msg: fmt.Sprintf("invalid operator %s", condition.Operator),
		}
	}

	var results []ConditionResult

	if condition.Field == "" {
		value, err := condition.Expr.Render(ctx)
		if err != nil {
			return nil, err
		}

		ok, err := operator.Eval(value, condition.Value)
		if err != nil {
			return nil, err
		}

		if !ok {
			results = append(results, ConditionResult{
				Expr:     value,
				Operator: condition.Operator,
				Value:    condition.Value,
			})
			return results, nil
		}
		return nil, nil

	} else {

		fields, err := input.FindAll(condition.Field)
		if err != nil {
			return nil, err
		}

		if fields == nil {
			return nil, fmt.Errorf("field %s not found", condition.Field)
		}

		for i := range fields {
			ctx.Cond.Field = fields[i]
			value := fields[i].Value

			if condition.Expr != nil {
				value, err = condition.Expr.Render(ctx)
				if err != nil {
					return nil, err
				}
				ctx.Cond.Expr = value
			}

			ok, err := operator.Eval(value, condition.Value)
			if err != nil {
				return nil, err
			}
			if !ok {
				result := ConditionResult{
					Field:    fields[i],
					Operator: condition.Operator,
					Value:    condition.Value,
				}
				if condition.Expr != nil {
					result.Expr = value
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func (c *Condition) render(ctx *Context) (*renderedCondition, error) {
	cond := renderedCondition{Expr: c.Expr}

	if c.Field != nil {
		field, err := c.Field.Render(ctx)
		if err != nil {
			return nil, err
		}
		cond.Field = field.(string)
	}

	operator, err := c.Operator.Render(ctx)
	if err != nil {
		return nil, err
	}

	cond.Operator = operator.(string)

	value, err := c.Value.Render(ctx)
	if err != nil {
		return nil, err
	}

	cond.Value = value

	return &cond, nil
}

func (c *Condition) validate() []parserError {
	var result []parserError

	if c.Operator == nil {
		result = append(result, parserError{
			msg: "missing condition's operator",
		})
	}

	if c.Value == nil {
		result = append(result, parserError{
			msg: "missing condition's value",
		})
	}

	if c.Field == nil && c.Expr == nil {
		result = append(result, parserError{
			msg: "condition must have field or expr defined",
		})
	}

	return result
}

type LogicalCondition struct {
	Not        *LogicalCondition  `json:"not,omitempty" yaml:"not,omitempty"`
	And        []LogicalCondition `json:"and,omitempty" yaml:"and,omitempty"`
	Or         []LogicalCondition `json:"or,omitempty" yaml:"or,omitempty"`
	*Condition `json:",inline,omitempty" yaml:",inline,omitempty"`
}

func (lc *LogicalCondition) Eval(ctx *Context, input *parser.Node) (bool, error) {
	if len(lc.And) > 0 {
		for _, condition := range lc.And {
			ok, err := condition.Eval(ctx, input)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil
	}
	if len(lc.Or) > 0 {
		for _, condition := range lc.Or {
			ok, err := condition.Eval(ctx, input)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
	if lc.Not != nil {
		ok, err := lc.Not.Eval(ctx, input)
		if err != nil {
			return false, err
		}
		return !ok, nil
	}

	failures, err := lc.Condition.Eval(ctx, input)
	if err != nil {
		return false, err
	}

	return failures == nil, nil
}

func (lc *LogicalCondition) validate() []parserError {
	var result []parserError

	if lc.And == nil && lc.Or == nil && lc.Not == nil && lc.Condition == nil {
		result = append(result, parserError{
			msg: "invalid condition",
		})
	}

	for i := range lc.And {
		if errs := lc.And[i].validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	for i := range lc.Or {
		if errs := lc.Or[i].validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if lc.Not != nil {
		if errs := lc.Not.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if lc.Condition != nil {
		if errs := lc.Condition.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	return result
}
