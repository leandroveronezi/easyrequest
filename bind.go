package easyRequest

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func register(name string, f bindingFunc, in []interface{}, out interface{}, comment string) error {

	_, exists := bindings[name]

	if exists {
		return errors.New("")
	}

	comment = strings.Trim(comment, " ")
	comment = strings.TrimRight(comment, "\r\n")

	bindings[name] = bindingType{
		f,
		name,
		comment,
		in,
		out,
	}

	return nil
}

func Bind(name string, f interface{}) error {

	v := reflect.ValueOf(f)

	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}

	// f must return either value and error or just error
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}

	if v.Type().NumIn() == 0 {
		return errors.New("function arguments mismatch")
	}

	in := make([]interface{}, 0)

	for i := 1; i < v.Type().NumIn(); i++ {

		in = append(in, getName(v.Type().In(i)))

	}

	var out interface{}

	for i := 0; i < v.Type().NumOut(); i++ {

		if v.Type().Out(i).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			continue
		}

		out = getName(v.Type().Out(i))

	}

	return register(name, func(req Request, raw []json.RawMessage) (interface{}, error) {

		if len(raw) != v.Type().NumIn()-1 {
			return nil, errors.New("function arguments mismatch")
		}

		args := []reflect.Value{}

		p := reflect.New(reflect.TypeOf(req))
		p.Elem().Set(reflect.ValueOf(req))

		args = append(args, p.Elem())

		for i := range raw {

			arg := reflect.New(v.Type().In(i + 1))

			if err := json.Unmarshal(raw[i], arg.Interface()); err != nil {
				return nil, err
			}

			args = append(args, arg.Elem())
		}

		errorType := reflect.TypeOf((*error)(nil)).Elem()

		res := v.Call(args)

		switch len(res) {
		case 0:
			// No results from the function, just return nil
			return nil, nil
		case 1:
			// One result may be a value, or an error
			if res[0].Type().Implements(errorType) {
				if res[0].Interface() != nil {
					return nil, res[0].Interface().(error)
				}
				return nil, nil
			}
			return res[0].Interface(), nil
		case 2:
			// Two results: first one is value, second is error
			if !res[1].Type().Implements(errorType) {
				return nil, errors.New("second return value must be an error")
			}
			if res[1].Interface() == nil {
				return res[0].Interface(), nil
			}
			return res[0].Interface(), res[1].Interface().(error)
		default:
			return nil, errors.New("unexpected number of return values")
		}
	}, in, out, "" /*funcDescription(f)*/)
}

func getName(t reflect.Type) interface{} {

	if t == reflect.TypeOf(time.Time{}) {

		return t.Name()
	}

	if t.Kind() == reflect.Slice {

		r := getName(t.Elem())

		valueSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(r)), 0, 0)

		//mv := reflect.New(reflect.TypeOf(r))

		valueSlice = reflect.Append(valueSlice, reflect.ValueOf(r))

		return valueSlice.Interface()

	} else if t.Kind() == reflect.Struct {

		b := make([]reflect.StructField, 0)

		for i := 0; i < t.NumField(); i++ {

			aux := t.Field(i)

			if t.Field(i).Type.Kind() == reflect.Slice {

				aux.Type = reflect.TypeOf(getName(t.Field(i).Type))

			} else if t.Field(i).Type.Kind() == reflect.Struct {

				aux.Type = reflect.TypeOf(getName(t.Field(i).Type))

			} else {
				aux.Type = reflect.TypeOf(string(""))
			}

			aux.Name = t.Field(i).Name

			b = append(b, aux)

		}

		v := reflect.New(reflect.StructOf(b)).Elem()

		for i := 0; i < t.NumField(); i++ {

			if !v.Field(i).IsValid() {
				continue
			}

			if t.Field(i).Type == reflect.TypeOf(time.Time{}) {
				v.Field(i).SetString("time.Time")
				continue
			}

			if t.Field(i).Type.Kind() == reflect.Struct {

				v.Field(i).Set(reflect.ValueOf(getName(t.Field(i).Type)))

				continue
			}

			if t.Field(i).Type.Kind() == reflect.Slice {

				v.Field(i).Set(reflect.ValueOf(getName(t.Field(i).Type)))
				continue
			}

			if t.Field(i).Type.Kind() == reflect.Struct {

				v.Field(i).Set(reflect.ValueOf(getName(t.Field(i).Type)))

				continue
			}

			v.Field(i).SetString(t.Field(i).Type.Name())

		}

		return v.Addr().Interface()

	}

	return t.Name()

}

// Get the name and path of a func
func funcPathAndName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// Get the name of a func (with package path)
func funcName(f interface{}) string {
	splitFuncName := strings.Split(funcPathAndName(f), ".")
	return splitFuncName[len(splitFuncName)-1]
}

/*
func funcDescription(f interface{}) string {

	fileName, _ := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).FileLine(0)
	funcName := funcName(f)
	fset := token.NewFileSet()

	// Parse src
	parsedAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	pkg := &ast.Package{
		Name:  "Any",
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = parsedAst

	comment := ""

	ast.Inspect(parsedAst, func(n ast.Node) bool {

		switch n.(type) {

		case *ast.FuncDecl:
			{
				fn, ok := n.(*ast.FuncDecl)
				if ok {

					if fn.Name.String() != funcName {
						return true
					}

					if fn.Recv == nil {

						fmt.Println(fn.Name)

						comment = fn.Doc.Text()
						comment = strings.Trim(comment, " ")
						comment = strings.TrimSuffix(comment, "\n")

						if fn.Type.Params != nil {

							for i := 0; i < len(fn.Type.Params.List); i++ {

								fmt.Println(
									fn.Type.Params.List[i].Tag,
									fn.Type.Params.List[i].Type,
									fn.Type.Params.List[i].Comment,
									fn.Type.Params.List[i].Doc,
									fn.Type.Params.List[i].Names,
								)

								res, ok := fn.Type.Params.List[i].Type.(*ast.Ident)
								if ok {

									fmt.Println(res)

								}

							}

						}

					}

				}

			}

		}

		return true
	})

	return comment

	/*
		importPath, _ := filepath.Abs("/")
		myDoc := doc.New(pkg, importPath, doc.AllDecls)

		for _, theFunc := range myDoc.Funcs {
			if theFunc.Name == funcName {
				return theFunc.Doc
			}
		}
	*/

	return ""
}
*/