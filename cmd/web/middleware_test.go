package main

import (
	"net/http"
	"testing"
)

// Testing NoSurf and SessionLoad middleware

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
		case http.Handler:
			// ok
		default:
			t.Errorf("Unexpected type %T", v)
	}
}


func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
		case http.Handler:
			// ok
		default:
			t.Errorf("Unexpected type %T", v)
	}
}

 