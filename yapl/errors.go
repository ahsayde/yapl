package yapl

import (
	"fmt"
	"strings"
)

type parserError struct {
	msg      string
	location []string
}

func (pe *parserError) Error() string {
	return pe.msg
}

type ParserError []parserError

func (perr ParserError) Error() string {
	var msgs []string
	for i := range perr {
		msgs = append(msgs, perr[i].Error())
	}
	return fmt.Sprintf("invalid policy\n%s", strings.Join(msgs, "\n"))
}

type RuntimeError struct {
	msg      string
	location []string
}

func (rte RuntimeError) Error() string {
	return rte.msg
}
