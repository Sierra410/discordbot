package main

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	err := cfg.Load()

	fmt.Printf("%#v\n", cfg)

	if err != nil {
		t.Error(err)
	}

	err = cfg.Save()

	fmt.Printf("%#v\n", cfg)

	if err != nil {
		t.Error(err)
	}
}
