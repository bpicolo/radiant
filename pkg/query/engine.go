package query

import (
	"fmt"

	"github.com/bpicolo/radiant/pkg/schema"
)

type Engine struct {
	interpreters map[string]Interpreter
}

func NewEngine() *Engine {
	interpreters := make(map[string]Interpreter)
	interpreters[""] = &JsonInterpreter{}
	interpreters["json"] = &JsonInterpreter{}

	return &Engine{
		interpreters: interpreters,
	}
}

func (q *Engine) Interpret(search schema.Search) (*Query, error) {
	interpreter, ok := q.interpreters[search.Query.Type]
	if !ok {
		return nil, fmt.Errorf("Don't know how to interpret query of type %s", search.Query.Type)
	}

	return interpreter.Interpret(search)
}
