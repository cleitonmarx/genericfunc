package genericfunc

import (
	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	tests := []struct {
		function       interface{}
		validationFunc func(*FunctionCache) error
		exception      error
	}{
		{ // A valid function
			func(item int) bool { return item > 10 },
			SimpleParamValidator(NewElemTypeSlice(new(int)), NewElemTypeSlice(new(bool))),
			nil,
		},
		{ //returns error when the function parameter has not the function kind
			"Not a function",
			SimpleParamValidator(nil, []reflect.Type{}),
			errors.New("GenericFunc.New: fn is not a function"),
		},
		{ // Returns error when expected parameters number are not equal
			func(idx, item int) {},
			SimpleParamValidator(NewElemTypeSlice(new(int)), []reflect.Type{}),
			errors.New("SimpleParamValidator: Number of parameters In expected: 1, actual: 2"),
		},
		{ // Returns error when expected parameters types are not equal
			func(items ...int) bool { return false },
			SimpleParamValidator(NewElemTypeSlice(new([]bool)), NewElemTypeSlice(new(bool))),
			errors.New("SimpleParamValidator: parameter In[0] expected type: []bool, actual type: bool"),
		},
		{ // Returns error when expected returns number are not equal
			func(item int) bool { return item > 10 },
			SimpleParamValidator(NewElemTypeSlice(new(int)), []reflect.Type{}),
			errors.New("SimpleParamValidator: Number of parameters Out expected: 0, actual: 1"),
		},
		{ // Returns error when expected return types are not equal
			func(items ...int) bool { return false },
			SimpleParamValidator(NewElemTypeSlice(new([]int)), NewElemTypeSlice(new(int64))),
			errors.New("SimpleParamValidator: parameter Out[0] expected type: int64, actual type: bool"),
		},
	}

	for _, test := range tests {
		_, err := New(test.function, test.validationFunc)
		if !(err == test.exception || err.Error() == test.exception.Error()) {
			t.Errorf("Validate expect error: %s, actual: %s", test.exception, err)
		}
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		function       interface{}
		validationFunc func(*FunctionCache) error
		fnParameter    interface{}
		result         interface{}
		exception      error
	}{
		{ // A valid function and parameters
			func(i int) int { return i * 3 },
			SimpleParamValidator(NewElemTypeSlice(new(GenericType)), NewElemTypeSlice(new(int))),
			3,
			9,
			nil,
		},
		{ // Returns error when the required type doesn't match with the specification
			func(i int) int { return i * 3 },
			SimpleParamValidator(NewElemTypeSlice(new(int)), NewElemTypeSlice(new(int))),
			"not a int",
			9,
			errors.New("GenericFunc.Call: params[0] 'string' is not convertible to 'int'"),
		},
	}

	for _, test := range tests {
		dynaFunc, err := New(test.function, test.validationFunc)
		if err != nil {
			t.Errorf("expect error: nil, actual: %s", err)
		}
		result, err := dynaFunc.Call(test.fnParameter)
		if !(err == test.exception || err.Error() == test.exception.Error()) {
			t.Errorf("expect error: nil, actual: %s", err)
		}
		if result != nil && result[0] != test.result {
			t.Errorf("expect result: %d, actual: %d", test.result, result[0])
		}
	}

}
