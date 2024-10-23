package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type AppError struct {
	message    string `json:"message"`
	statusCode int    `json:"status_code"`
}

func (e AppError) Error() string {
	return e.message
}

func someOperation() error {
	return &AppError{message: "something went wrong", statusCode: 500}
}

func main() {
	// valid use of As
	{
		err := someOperation()
		var appErr *AppError
		if errors.As(err, &appErr) {
			fmt.Println("handled properly")
		} else {
			panic("unexpected")
		}
	}

	// valid error use of Is
	{
		err := os.ErrNotExist
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("handled properly")
		} else {
			panic("unexpected")
		}
	}

	// invalid use of Is: custom error not detected
	{
		err := someOperation()
		if errors.Is(err, &AppError{message: "something went wrong"}) {
			panic("Custom error detected")
		} else {
			fmt.Println("Different error")
		}
	}

	// invalid use of As: invalid error wrapping
	{
		err := fmt.Errorf("some context: %v", io.EOF) // %v does not wrap the error
		if errors.Is(err, io.EOF) {
			fmt.Println("EOF error detected")
		} else {
			fmt.Println("`Different error")
		}
	}
}
