package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	leftDelim  = "${"
	rightDelim = "}"
)

var (
	gtmpl = template.New("").
		Funcs(processors).
		Option("missingkey=error").
		Delims(leftDelim, rightDelim)

	replacer = strings.NewReplacer(
		leftDelim,
		fmt.Sprintf("%s toYaml (", leftDelim),
		rightDelim,
		fmt.Sprintf(") %s", rightDelim),
	)
)

type Renderable struct {
	value interface{}
	tmpl  *template.Template
}

func (r *Renderable) init(value interface{}) error {
	var err error
	r.value = value

	switch v := value.(type) {
	case string:
		r.tmpl, err = gtmpl.New("").Parse(v)
		if err != nil {
			return err
		}
	case map[string]interface{}, []interface{}, interface{}:
		raw, err := yaml.Marshal(r.value)
		if err != nil {
			return err
		}
		r.tmpl, err = gtmpl.New("").Parse(string(raw))
		if err != nil {
			return err
		}
	}
	return nil
}

func newRenderable(value interface{}) (*Renderable, error) {
	r := &Renderable{}
	err := r.init(value)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderable) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		if strings.HasPrefix(n.Value, leftDelim) && strings.HasSuffix(n.Value, rightDelim) {
			n.Value = replacer.Replace(n.Value)
		}
	}
	var value interface{}
	err := n.Decode(&value)
	if err != nil {
		return err
	}
	return r.init(value)
}

func (r *Renderable) Render(ctx interface{}) (interface{}, error) {
	if r.tmpl == nil {
		return r.value, nil
	}

	var buff bytes.Buffer
	err := r.tmpl.Execute(&buff, ctx)
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = yaml.Unmarshal(buff.Bytes(), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
