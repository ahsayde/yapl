package yapl

import "encoding/json"

type Result struct {
	Passed  bool          `json:"passed" yaml:"passed"`
	Ignored bool          `json:"ignored" yaml:"ignored"`
	Failed  bool          `json:"failed" yaml:"failed"`
	Results []interface{} `json:"errors" yaml:"errors"`
}

func (r *Result) JSON() string {
	raw, _ := json.MarshalIndent(r, "", "  ")
	return string(raw)
}
