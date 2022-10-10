package yapl

import (
	"github.com/ahsayde/yapl/internal/parser"
	"github.com/ahsayde/yapl/internal/renderer"
)

type Policy struct {
	Metadata map[string]interface{}         `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Match    *LogicalCondition              `json:"match,omitempty" yaml:"match,omitempty"`
	Exclude  *LogicalCondition              `json:"exclude,omitempty" yaml:"exclude,omitempty"`
	Globals  map[string]renderer.Renderable `json:"globals,omitempty" yaml:"globals,omitempty"`
	Rules    []Rule                         `json:"rules" yaml:"rules"`
}

func (p *Policy) Eval(input, params map[string]interface{}) (*Result, error) {
	ctx := newContext(input, params)
	node := parser.Parse(ctx.Input)

	if p.Globals != nil {
		globals := make(map[string]interface{})
		for k, v := range p.Globals {
			value, err := v.Render(ctx)
			if err != nil {
				return nil, err
			}
			globals[k] = value
		}
		ctx.Globals = globals
	}

	if p.Match != nil {
		ok, err := p.Match.Eval(ctx, node)
		if err != nil {
			return nil, err
		}
		if !ok {
			return &Result{Ignored: true}, nil
		}
	}

	if p.Exclude != nil {
		ok, err := p.Exclude.Eval(ctx, node)
		if err != nil {
			return nil, err
		}
		if ok {
			return &Result{Ignored: true}, nil
		}
	}

	result := Result{}
	for _, rule := range p.Rules {
		results, err := rule.Eval(ctx, node)
		if err != nil {
			return nil, err
		}
		result.Results = append(result.Results, results...)
	}

	if len(result.Results) > 0 {
		result.Failed = true
	} else {
		result.Passed = true
	}

	return &result, nil
}

func (p *Policy) validate() []parserError {
	var result []parserError

	if p.Match != nil {
		if errs := p.Match.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if p.Exclude != nil {
		if errs := p.Exclude.validate(); errs != nil {
			result = append(result, errs...)
		}
	}

	if p.Rules == nil {
		result = append(result, parserError{
			msg: "policy must contains at least one rule",
		})
	}

	for i := range p.Rules {
		if errs := p.Rules[i].validate(); errs != nil {
			result = append(result, errs...)
		}
	}
	return result
}
