package easyRequest

import (
	"errors"
	"reflect"
)

func WSBind(name string, f interface{}) error {

	v := reflect.ValueOf(f)

	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}

	// f must return either value and error or just error
	if n := v.Type().NumOut(); n > 0 {
		return errors.New("no return for websocket function")
	}

	if v.Type().NumIn() == 0 {
		return errors.New("function arguments mismatch")
	}

	return wsRegister(name, func(req WsRequest) error {

		args := []reflect.Value{}

		p := reflect.New(reflect.TypeOf(req))
		p.Elem().Set(reflect.ValueOf(req))

		args = append(args, p.Elem())

		errorType := reflect.TypeOf((*error)(nil)).Elem()

		res := v.Call(args)

		switch len(res) {
		case 0:
			// No results from the function, just return nil
			return nil
		case 1:
			// One result may be a value, or an error

			if res[0].Type().Implements(errorType) {
				if res[0].Interface() != nil {
					return res[0].Interface().(error)
				}
				return nil
			}

		default:
			return errors.New("unexpected number of return values")
		}

		return nil

	})
}

func wsRegister(name string, f wsBindingFunc) error {

	_, exists := wsBindings[name]

	if exists {
		return errors.New("")
	}

	wsBindings[name] = wsBindingType{
		f,
		name,
	}

	return nil
}
