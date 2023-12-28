package interpreter

import (
	"errors"
	"interpreters/src/lex"
)

type Environment struct {
	variables map[string]any
}

func createEnvironment() *Environment {
	return &Environment{variables: make(map[string]any)}
}

func (e *Environment) Assign(name lex.Token, value any) error {
	if _, contained := e.variables[name.Lexeme]; contained {
		e.variables[name.Lexeme] = value
		return nil
	} else {
		return errors.New("Undefined variable: " + name.Lexeme)
	}
}

func (e *Environment) Define(name string, value any) {
	e.variables[name] = value
}

func (e *Environment) Get(name string) (any, error) {
	if value, contained := e.variables[name]; contained {
		return value, nil
	} else {
		return nil, errors.New("Cannot find variable: " + name)
	}
}
