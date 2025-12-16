package xslices

import (
	"context"
	"strconv"
)

func Trans[A, B any](a []A, fn func(A) B) []B {
	results := make([]B, 0, len(a))
	for _, v := range a {
		results = append(results, fn(v))
	}
	return results
}

func TransError[A, B any](a []A, fn func(A) (B, error)) ([]B, error) {
	results := make([]B, 0, len(a))
	for _, v := range a {
		b, err := fn(v)
		if err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, nil
}

func TransContext[A, B any](ctx context.Context, a []A, fn func(context.Context, A) B) []B {
	results := make([]B, 0, len(a))
	for _, v := range a {
		results = append(results, fn(ctx, v))
	}
	return results
}

func TransErrorContext[A, B any](ctx context.Context, a []A, fn func(context.Context, A) (B, error)) ([]B, error) {
	results := make([]B, 0, len(a))
	for _, v := range a {
		b, err := fn(ctx, v)
		if err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, nil
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Number interface {
	Integer | ~float32 | ~float64
}

func Ints[A, B Integer](a []A) []B {
	results := make([]B, 0, len(a))
	for _, v := range a {
		results = append(results, B(v))
	}
	return results
}

func Numbers[A, B Number](a []A) []B {
	results := make([]B, 0, len(a))
	for _, v := range a {
		results = append(results, B(v))
	}
	return results
}

func Atoi[I Integer](s string) (I, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return I(i), err
}

func Atois[I Integer](s []string) ([]I, error) {
	return TransError[string, I](s, Atoi)
}

func Itoa[I Integer](i I) string {
	return strconv.FormatInt(int64(i), 10)
}

func Itoas[I Integer](a []I) []string {
	return Trans[I, string](a, Itoa)
}
