package monads

import "fmt"

//MIT License
//
//Copyright (c) 2022 Samuel Berthe
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.

// Result represents a result of an action having one
// of the following output: success or failure.
// An instance of Result is an instance of either Ok or Err.
// It could be compared to `Either[error, T]`.
type Result[T any] struct {
	isErr bool
	value T
	err   error
}

// Ok builds a Result when value is valid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		isErr: false,
	}
}

// Err builds a Result when value is invalid.
// Play: https://go.dev/play/p/PDwADdzNoyZ
func Err[T any](err error) Result[T] {
	return Result[T]{
		err:   err,
		isErr: true,
	}
}

// Errf builds a Result when value is invalid.
// Errf formats according to a format specifier and returns the error as a value that satisfies Result[T].
// Play: https://go.dev/play/p/N43w92SM-Bs
func Errf[T any](format string, a ...any) Result[T] {
	return Err[T](fmt.Errorf(format, a...))
}

// IsOk returns true when value is valid.
// Play: https://go.dev/play/p/sfNvBQyZfgU
func (r Result[T]) IsOk() bool {
	return !r.isErr
}

// IsError returns true when value is invalid.
// Play: https://go.dev/play/p/xkV9d464scV
func (r Result[T]) IsError() bool {
	return r.isErr
}

// Error returns error when value is invalid or nil.
// Play: https://go.dev/play/p/CSkHGTyiXJ5
func (r Result[T]) Error() error {
	return r.err
}

// MustGet returns value when Result is valid or panics.
// Play: https://go.dev/play/p/8LSlndHoTAE
func (r Result[T]) MustGet() T {
	if r.isErr {
		panic(r.err)
	}

	return r.value
}
