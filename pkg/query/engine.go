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
	interpreters[""] = &YamlInterpreter{}
	interpreters["yaml"] = &YamlInterpreter{}

	return &Engine{
		interpreters: interpreters,
	}
}

func (q *Engine) Interpret(search *schema.Search) (*Query, error) {
	interpreter, ok := q.interpreters[search.QueryDefinition.Type]
	if !ok {
		return nil, fmt.Errorf("Don't know how to interpret query of type %s", search.QueryDefinition.Type)
	}

	return interpreter.Interpret(search)
}
