package genericfunc

import (
	"errors"
	"fmt"
	"reflect"
)

// GenericType represents a any reflect.Type
type GenericType int

var genericType = reflect.TypeOf(new(GenericType)).Elem()

// FunctionCache keeps in cached the GenericFunc reflections objects
type FunctionCache struct {
	FnValue  reflect.Value
	FnType   reflect.Type
	TypesIn  []reflect.Type
	TypesOut []reflect.Type
}

// GenericFunc is a type used to validate and call dynamic functions
type GenericFunc struct {
	Cache *FunctionCache
}

// Call calls a dynamic function
func (d *GenericFunc) Call(params ...interface{}) ([]interface{}, error) {
	paramsIn := make([]reflect.Value, len(params))
	for i, param := range params {
		paramValue := reflect.ValueOf(param)
		paramType := paramValue.Type()
		if d.Cache.TypesIn[i] != genericType {
			if !paramType.ConvertibleTo(d.Cache.TypesIn[i]) {
				return nil, fmt.Errorf("GenericFunc.Call: params[%d] '%s' is not convertible to '%s'", i, paramType, d.Cache.TypesIn[i])
			}
		}
		paramsIn[i] = paramValue
	}
	paramsOut := d.Cache.FnValue.Call(paramsIn)
	interfaceOut := make([]interface{}, len(paramsOut))
	for i, item := range paramsOut {
		interfaceOut[i] = item.Interface()
	}
	return interfaceOut, nil
}

// New instantiates a new GenericFunc pointer
func New(fn interface{}, validateFunc func(*FunctionCache) error) (*GenericFunc, error) {
	cache := &FunctionCache{}
	cache.FnValue = reflect.ValueOf(fn)

	if cache.FnValue.Kind() != reflect.Func {
		return nil, errors.New("GenericFunc.New: fn is not a function")
	}

	cache.FnType = cache.FnValue.Type()
	numTypesIn := cache.FnType.NumIn()
	cache.TypesIn = make([]reflect.Type, numTypesIn)
	for i := 0; i < numTypesIn; i++ {
		cache.TypesIn[i] = cache.FnType.In(i)
	}

	numTypesOut := cache.FnType.NumOut()
	cache.TypesOut = make([]reflect.Type, numTypesOut)
	for i := 0; i < numTypesOut; i++ {
		cache.TypesOut[i] = cache.FnType.Out(i)
	}
	if err := validateFunc(cache); err != nil {
		return nil, err
	}

	return &GenericFunc{Cache: cache}, nil
}

// SimpleParamValidator creates a function to validate GenericFunc based in the In and Out function parameters
func SimpleParamValidator(In []reflect.Type, Out []reflect.Type) func(cache *FunctionCache) error {
	return func(cache *FunctionCache) error {
		if In != nil {
			if len(In) != len(cache.TypesIn) {
				return fmt.Errorf("SimpleParamValidator: Number of parameters In expected: %d, actual: %d", len(In), len(cache.TypesIn))
			}
			for i, paramIn := range In {
				if paramIn != genericType && paramIn != cache.TypesIn[i] {
					return fmt.Errorf("SimpleParamValidator: parameter In[%d] expected type: %s, actual type: %s", i, paramIn, cache.TypesOut[i])
				}
			}
		}
		if Out != nil {
			if len(Out) != len(cache.TypesOut) {
				return fmt.Errorf("SimpleParamValidator: Number of parameters Out expected: %d, actual: %d", len(Out), len(cache.TypesOut))
			}
			for i, paramOut := range Out {
				if paramOut != genericType && paramOut != cache.TypesOut[i] {
					return fmt.Errorf("SimpleParamValidator: parameter Out[%d] expected type: %s, actual type: %s", i, paramOut.String(), cache.TypesOut[i].Kind().String())
				}
			}
		}
		return nil
	}
}

// NewElemTypeSlice creates a slice of items elem types
func NewElemTypeSlice(items ...interface{}) []reflect.Type {
	typeList := make([]reflect.Type, len(items))
	for i, item := range items {
		typeItem := reflect.TypeOf(item)
		if typeItem.Kind() == reflect.Ptr {
			typeList[i] = typeItem.Elem()
		}
	}
	return typeList
}
