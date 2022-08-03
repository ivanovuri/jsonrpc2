package inspect

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Returns true in case call accepts positional parameters.
// False if first argument is not struct type.
func NamedCall(method any) bool {
	switch reflect.TypeOf(method).In(0).Kind() {
	case reflect.Struct:
		return true
	default:
		return false
	}
}

func NumInArgs(method any) int {
	return reflect.TypeOf(method).NumIn()
}

func NumOutArgs(method any) int {
	return reflect.TypeOf(method).NumOut()
}

// Compare number of values provided in remote request
// with number of parameters accepted by function
func CompareInNumExpectedWithRequest(inValues []any, rpc any) error {
	inValuesNum := reflect.TypeOf(rpc).NumIn()

	if inValuesNum != len(inValues) {
		return fmt.Errorf("paramers count %d doesn't match rpcCall %d",
			len(inValues),
			inValuesNum)
	}

	return nil
}

func ParseNamedParams(jsParams []byte, rpc any) ([]reflect.Value, error) {
	funcInputType := reflect.TypeOf(rpc)
	structInTheFunction := funcInputType.In(0)
	createdWithReflectionOnlyNew := reflect.New(structInTheFunction)
	createdWithReflection := createdWithReflectionOnlyNew.Interface()

	if err := json.Unmarshal(jsParams, &createdWithReflection); err != nil {
		return nil, fmt.Errorf("paramers marshalling error %s", err)
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.Indirect(createdWithReflectionOnlyNew)

	return in, nil
}

func ParsePositionalParams(jsParams []byte, rpc any) ([]reflect.Value, error) {
	var inValues []interface{}
	if err := json.Unmarshal(jsParams, &inValues); err != nil {
		return nil, fmt.Errorf("paramers marshalling error %s", err)
	}

	if err := CompareInNumExpectedWithRequest(inValues, rpc); err != nil {
		return nil, err
	}

	inValuesNum := reflect.TypeOf(rpc).NumIn()
	inputParams := make([]reflect.Value, inValuesNum)

	for i := 0; i < inValuesNum; i++ {
		pType := reflect.TypeOf(rpc).In(i)

		switch pType.Kind() {

		case reflect.Int:
			if err := TypeEq(inValues[i], reflect.Float64); err != nil {
				return nil, fmt.Errorf("%s", err)
			}

			value := int(reflect.ValueOf(inValues[i]).Interface().(float64))
			inputParams[i] = reflect.ValueOf(value)

		case reflect.Float32:
			if err := TypeEq(inValues[i], reflect.Float64); err != nil {
				return nil, fmt.Errorf("%s", err)
			}

			value := float32(reflect.ValueOf(inValues[i]).Interface().(float64))
			inputParams[i] = reflect.ValueOf(value)

		case reflect.String:
			if err := TypeEq(inValues[i], reflect.String); err != nil {
				return nil, fmt.Errorf("%s", err)
			}

			inputParams[i] = reflect.ValueOf(inValues[i])
		}
	}

	return inputParams, nil
}

func TypeEq(verifiable, reference any) error {
	if reflect.ValueOf(verifiable).Kind() != reference {
		return fmt.Errorf("param caller types doesn't match")
	}
	return nil
}

func ExecuteMethod(method any, parameters []reflect.Value) []reflect.Value {
	return reflect.ValueOf(method).Call(parameters)
}
