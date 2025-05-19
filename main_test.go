//go:build features

package main

import "testing"

func TestFeatures(t *testing.T) {
	if err := app.Run(t.Context()); err != nil {
		t.Fatal(err.Error())
	}
}
