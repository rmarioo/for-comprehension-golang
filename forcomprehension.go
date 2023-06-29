package errm

import (
	"errors"
	"fmt"
)

func panicToPair[T any](f func() T) (r T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error happened " + toString(r))
		}
	}()

	i := f()

	return i, nil
}

func toString(r any) string {
	return fmt.Sprintf("%v", r)
}

func wrapToPanic[A any, B any](f1 FuncErr[A, B]) func(a A) B {
	return func(a A) B {
		b, err := f1(a)
		if err != nil {
			panic(err)
		} else {
			return b
		}
	}

}

func panicIfError[T any](v T, e error) T {
	if e != nil {
		panic(e)
	} else {
		return v
	}
}

type FuncErr[A, B any] func(A) (B, error)

func compose2[A any, B any, C any](f1 FuncErr[A, B], f2 FuncErr[B, C]) FuncErr[A, C] {
	return func(a A) (C, error) {
		b, err := f1(a)
		if err != nil {
			var result C
			return result, err
		} else {
			return f2(b)
		}

	}
}

func compose3[A any, B any, C any, D any](f1 FuncErr[A, B], f2 FuncErr[B, C], f3 FuncErr[C, D]) FuncErr[A, D] {
	two := compose2(f1, f2)
	return compose2(two, f3)
}

func compose4[A any, B any, C any, D any, E any](
	f1 FuncErr[A, B],
	f2 FuncErr[B, C],
	f3 FuncErr[C, D],
	f4 FuncErr[D, E]) FuncErr[A, E] {
	two := compose3(f1, f2, f3)
	return compose2(two, f4)
}

func compose5[A any, B any, C any, D any, E any, F any](
	f1 FuncErr[A, B],
	f2 FuncErr[B, C],
	f3 FuncErr[C, D],
	f4 FuncErr[D, E],
	f5 FuncErr[E, F]) FuncErr[A, F] {
	two := compose4(f1, f2, f3, f4)
	return compose2(two, f5)
}
