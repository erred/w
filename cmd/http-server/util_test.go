package main

import "testing"

func TestCanonical(t *testing.T) {
	tcs := []struct {
		in, out string
	}{
		{
			"/", "/",
		}, {
			"/x", "/x/",
		}, {
			"/x/", "/x/",
		}, {
			"/x/y", "/x/y/",
		}, {
			"/index", "/",
		}, {
			"/index.html", "/",
		}, {
			"/x/index", "/x/",
		}, {
			"/x/index.html", "/x/",
		},
	}

	for i, tc := range tcs {
		if got := canonical(tc.in); got != tc.out {
			t.Errorf("TestCanonical %d: expected %q got %q\n", i, tc.out, got)
		}
	}
}
