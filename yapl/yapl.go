package yapl

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Policy, error) {
	var policy Policy

	if err := yaml.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to marshal policy, error: %w", err)
	}

	if errs := policy.validate(); errs != nil {
		return nil, ParserError(errs)
	}

	return &policy, nil
}
