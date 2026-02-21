package di

import (
	"reflect"
)

// Injectable marks a constructor Function of a Struct for DI
func Injectable[T any](creator func() T) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	getContainer().injectable(typ, func() any { return creator() })
}

func Replace[T any](creator func() T, identifier ...string) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	getContainer().replace(typ, func() any { return creator() })
}

func ReplaceInstance[T any](instance T, identifier ...string) {
	Replace(func() T { return instance }, identifier...)
}

// Inject gets or create a Instance of the Struct used the Injectable constructor Function
func Inject[T any](identifier ...string) T {
	var result T
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if instance, ok := getContainer().inject(typ, identifier...); ok {
		if result, ok = instance.(T); ok {
			return result
		}
	}
	return result
}

func Destroy[T any](identifier ...string) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	getContainer().destroy(typ, identifier...)
}

func DestroyAllMatching(match func(string) bool) {
	getContainer().destroyAllMatching(match)
}
