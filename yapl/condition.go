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

type NestedCondition struct {
	Not        *NestedCondition  `json:"not,omitempty" yaml:"not,omitempty"`
	And        []NestedCondition `json:"and,omitempty" yaml:"and,omitempty"`
	Or         []NestedCondition `json:"or,omitempty" yaml:"or,omitempty"`
	*Condition `json:",inline,omitempty" yaml:",inline,omitempty"`
}

func (nc *NestedCondition) Eval(ctx *Context, input *parser.Node) (bool, error) {
	if len(nc.And) > 0 {
		for _, condition := range nc.And {
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
	if len(nc.Or) > 0 {
		for _, condition := range nc.Or {
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
	if nc.Not != nil {
		ok, err := nc.Not.Eval(ctx, input)
		if err != nil {
			return false, err
		}
		return !ok, nil
	}

	failures, err := nc.Condition.Eval(ctx, input)
	if err != nil {
		return false, err
	}

	return failures == nil, nil
}

func (ns *NestedCondition) validate() []parserError {
	var result []parserError

	if ns.And == nil && ns.Or == nil && ns.Not == nil && ns.Condition == nil {
		result = append(result, parserError{
			msg: "invalid condition",
		})
	}

	for i := range ns.And {
		if errs := ns.And[i].validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	for i := range ns.Or {
		if errs := ns.Or[i].validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if ns.Not != nil {
		if errs := ns.Not.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if ns.Condition != nil {
		if errs := ns.Condition.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	return result
}
