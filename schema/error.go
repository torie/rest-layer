package schema

import (
	"fmt"
	"sort"
	"strings"
)

// ErrorMap contains a map of errors by field name.
type ErrorMap map[string][]interface{}

// Error implements the built-in error interface.
func (e ErrorMap) Error() string {
	errs := make([]string, 0, len(e))
	for key := range e {
		errs = append(errs, key)
	}
	sort.Strings(errs)
	for i, key := range errs {
		errs[i] = fmt.Sprintf("%s is %s", key, e[key])
	}
	return strings.Join(errs, ", ")
}

// Extend copies all errors from err to e.
func (e ErrorMap) Extend(err ErrorMap) {
	for k, v := range err {
		e[k] = append(e[k], v...)
	}
}

// Tidy returns an un-typed nil if e is empty, and an ErrorMap otherwise.
func (e ErrorMap) Tidy() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// ErrorSlice contains a concatenation of several errors.
type ErrorSlice []error

// Add adds an error to e and returns a new slice. If err is ErrorSlice, it is
// extended so that all the elements in err are apppended to e. If err is nil,
// then no error is appended. Note that calling Add on a nil slice is valid.
func (e ErrorSlice) Add(err error) ErrorSlice {
	switch et := err.(type) {
	case nil:
		// don't append nil errors.
	case ErrorSlice:
		// Extend error slices.
		e = append(e, et...)
	default:
		e = append(e, et)
	}
	return e
}

// Error implements the built-in error interface.
func (e ErrorSlice) Error() string {
	sl := make([]string, 0, len(e))
	for _, err := range e {
		sl = append(sl, err.Error())
	}
	return strings.Join(sl, ", ")
}

// Tidy returns an un-typed nil if e is empty, a single error if there is only
// one error, or an ErrorSlice otherwise.
func (e ErrorSlice) Tidy() error {
	switch len(e) {
	case 0:
		return nil
	case 1:
		return e[0]
	default:
		return e
	}
}
