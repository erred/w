package main

import (
	"fmt"
	"testing"
)

func TestParseTemplates(t *testing.T) {
	o := newOptions()
	if err := o.parseTemplates(); err != nil {
		t.Errorf("TestParseTemplates: %v", err)
	}
}
