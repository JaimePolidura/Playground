package interpreter

import (
	"errors"
	"interpreters/src/lex"
)

type Environment struct {
	parent    *Environment //Used for scope variables
	variables map[string]any
}

func createRootEnvironment() *Environment {
	return &Environment{variables: make(map[string]any)}
}

func createChildEnvironment(parent *Environment) *Environment {
	return &Environment{variables: make(map[string]any), parent: parent}
}

func (e *Environment) CopyInto(otherEnvironment *Environment) *Environment {
	actual := otherEnvironment
	for actual != nil {
		for k, v := range actual.variables {
			if _, contained := e.variables[k]; !contained {
				e.variables[k] = v
			}
		}

		actual = actual.parent
	}

	return e
}

func (e *Environment) Assign(name lex.Token, value any) error {
	if _, contained := e.variables[name.Lexeme]; contained {
		e.variables[name.Lexeme] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	return errors.New("Undefined variable: " + name.Lexeme)
}

func (e *Environment) Define(name string, value any) {
	e.variables[name] = value
}

func (e *Environment) Get(name string) (any, error) {
	if value, contained := e.variables[name]; contained {
		return value, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}

	return nil, errors.New("Cannot find variable: " + name)
}
