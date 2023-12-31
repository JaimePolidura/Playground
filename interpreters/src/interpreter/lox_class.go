package interpreter

import (
	"errors"
	"interpreters/src/syntax"
)

type LoxClass struct {
	Name       string
	Methods    map[string]LoxFunction
	SuperClass *LoxClass
}

type LoxInstance struct {
	KClass        LoxClass
	Properties    map[string]any
	superInstance *LoxInstance
}

func (l LoxInstance) GetProperty(name string) (any, error) {
	actualInstace := &l

	for actualInstace != nil {
		if propertyValue, propertyExists := actualInstace.Properties[name]; propertyExists {
			return propertyValue, nil
		} else if methodValue, methodExists := actualInstace.KClass.Methods[name]; methodExists {
			return methodValue, nil
		}

		actualInstace = actualInstace.superInstance
	}

	return nil, errors.New("unknown property " + name + " on class " + l.KClass.Name)
}

func (l LoxClass) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	instance := &LoxInstance{KClass: l, Properties: make(map[string]any)}

	instance.KClass.Methods = instance.bindThisToMethods()

	instance.instanciateSuperClasses()

	instance.callConstructor(interpreter, args)

	return syntax.CreateLiteralExpression(instance), nil
}

func (l *LoxInstance) callConstructor(interpreter *Interpreter, args []any) error {
	for name, method := range l.KClass.Methods {
		if name == "init" {
			if method.Arity() != len(args) {
				return errors.New("invalid constructor args")
			}

			method.Call(interpreter, args)
			break
		}
	}

	return nil
}

func (l *LoxInstance) instanciateSuperClasses() {
	actualSuperClass := l.KClass.SuperClass
	subClass := l

	for actualSuperClass != nil {
		superClassInstance := &LoxInstance{KClass: *actualSuperClass, Properties: make(map[string]any)}
		superClassInstance.bindThisToMethods()
		subClass.superInstance = superClassInstance
		actualSuperClass = actualSuperClass.SuperClass
	}
}

func (l *LoxInstance) bindThisToMethods() map[string]LoxFunction {
	newMethods := make(map[string]LoxFunction)
	for name, method := range l.KClass.Methods {
		newMethods[name] = method.BindThis(l)
	}

	return newMethods
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
