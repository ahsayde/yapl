package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Type int

const (
	Object Type = iota
	Array
	Scalar
)

var (
	regex = regexp.MustCompile("^([a-zA-Z0-9]+)\\[([^(.)$]+)\\]")
)

type Node struct {
	Type      Type
	Path      string
	Index     int
	Value     interface{}
	Parent    *Node
	mapItems  map[string]*Node
	listItems map[int]*Node
}

func Parse(in interface{}) *Node {
	return lookup(in, nil, "", 0)
}

func (n *Node) Find(key string) (*Node, error) {
	nodes, err := n.find(parseKeyPath(key))
	if err != nil {
		return nil, err
	}
	if nodes == nil {
		return nil, nil
	}
	return nodes[0], nil
}

func (n *Node) FindAll(key string) ([]*Node, error) {
	return n.find(parseKeyPath(key))
}

func parseKeyPath(path string) []string {
	var keys []string
	parts := strings.Split(path, ".")
	for _, part := range parts {
		groups := regex.FindStringSubmatch(part)
		if groups == nil {
			keys = append(keys, part)
		} else {
			keys = append(keys, groups[1:]...)
		}
	}
	return keys
}

func (n *Node) find(path []string) ([]*Node, error) {
	if len(path) == 0 {
		return []*Node{n}, nil
	}
	if n.Type == Object {
		item := n.mapItems[path[0]]
		if item == nil {
			return []*Node{{}}, nil
		}
		return item.find(path[1:])
	} else if n.Type == Array {
		if path[0] == "*" {
			var nodes []*Node
			for i := range n.listItems {
				items, err := n.listItems[i].find(path[1:])
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, items...)
			}
			return nodes, nil
		} else {
			index, err := strconv.Atoi(path[0])
			if err != nil {
				return nil, fmt.Errorf("invalid index '%s'", path[0])
			}
			if int(index) > len(n.listItems)-1 {
				return nil, fmt.Errorf("index '%d' of field '%s' is out of range", index, n.Path)
			}
			return n.listItems[int(index)].find(path[1:])
		}
	}
	return []*Node{{}}, nil
}

func lookup(in interface{}, parent *Node, path string, index int) *Node {
	node := &Node{
		Path:   strings.TrimPrefix(path, "."),
		Index:  index,
		Parent: parent,
		Value:  in,
	}
	switch item := in.(type) {
	case map[string]interface{}:
		node.Type = Object
		node.mapItems = make(map[string]*Node)
		for key, val := range item {
			p := fmt.Sprintf("%s.%s", path, key)
			node.mapItems[key] = lookup(val, node, p, 0)
		}
	case []interface{}:
		node.Type = Array
		node.listItems = make(map[int]*Node)
		for i, val := range item {
			p := fmt.Sprintf("%s[%d]", path, i)
			node.listItems[i] = lookup(val, node, p, index)
		}
	default:
		node.Type = Scalar
		node.Value = in
	}
	return node
}
