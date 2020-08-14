package ssssg

import (
	"errors"
	"fmt"
	"net/http"
)

// Serve from the public directory.
func Serve(dir string) error {
	fmt.Println("Serving", dir, "from http://localhost:2020")
	if err := http.ListenAndServe("localhost:2020", http.FileServer(http.Dir(dir))); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
