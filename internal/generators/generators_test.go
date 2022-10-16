package generators_test

import (
	"fmt"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
)

func ExampleGenerateString() {
	// Generate the string consisting of 10 symbols.
	val, _ := generators.GenerateString(10)
	fmt.Println(len(val))

	// The zero-valued size leads to the error.
	_, err := generators.GenerateString(0)
	fmt.Println(err.Error())

	// Output:
	// 10
	// random string length is missing
}

func ExampleGenerateID() {
	db := storage.NewMemoryRepo()

	// Generate the ID represented by a string consisting of 7 symbols.
	id, _ := generators.GenerateID(db, 7)
	fmt.Println(len(id))

	// Output:
	// 7
}
