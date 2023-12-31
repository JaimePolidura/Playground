package interpreter

import (
	"errors"
	"interpreters/src/syntax"
)

type LoxClass struct {
	Name    string
	Methods map[string]LoxFunction
}

type LoxInstance struct {
	KClass     LoxClass
	Properties map[string]any
}

func (l LoxInstance) GetProperty(name string) (any, error) {
	if propertyValue, propertyExists := l.Properties[name]; propertyExists {
		return propertyValue, nil
	} else if methodValue, methodExists := l.KClass.Methods[name]; methodExists {
		return methodValue, nil
	} else {
		return nil, errors.New("unknown property " + name + " on class " + l.KClass.Name)
	}
}

func (l LoxClass) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	instance := &LoxInstance{KClass: l, Properties: make(map[string]any)}

	newMethods := make(map[string]LoxFunction)
	for name, method := range instance.KClass.Methods {
		newMethods[name] = method.BindThis(instance)
	}

	instance.KClass.Methods = newMethods

	for name, method := range newMethods {
		if name == "init" {
			if method.Arity() != len(args) {
				return nil, errors.New("invalid constructor args")
			}

			method.Call(interpreter, args)
			break
		}
	}
	
	return syntax.CreateLiteralExpression(instance), nil
}

func (l LoxClass) Arity() int {
	for name, method := range l.Methods {
		if name == "init" {
			//Constructor
			return len(method.FunctionStmt.Params)
		}
	}

	return 0
}
