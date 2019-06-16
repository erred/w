package main

import "testing"

func TestURL(t *testing.T) {
	type TU struct {
		in           string
		outDst       string
		outRelative  string
		outCanonical string
	}
	cases := []TU{
		TU{
			"src/index.html",
			"dst/index.html",
			"/",
			"https://seankhliao.com/",
		}, TU{
			"src/tos.html",
			"dst/tos.html",
			"/tos",
			"https://seankhliao.com/tos",
		}, TU{
			"src/authed/index.html",
			"dst/authed/index.html",
			"/authed",
			"https://seankhliao.com/authed",
		}, TU{
			"src/blog/post.md",
			"dst/blog/post.html",
			"/blog/post",
			"https://seankhliao.com/blog/post",
		},
	}
	for i, c := range cases {
		u := NewURL(c.in)

		exp, got := c.outDst, u.Dst()
		if exp != got {
			t.Errorf("TestURL Dst case %v: expected %v got %v\n", i, exp, got)
		}

		exp, got = c.outCanonical, u.Canonical()
		if exp != got {
			t.Errorf("TestURL Canonical case %v: expected %v got %v\n", i, exp, got)
		}

		exp, got = c.outRelative, u.Relative()
		if exp != got {
			t.Errorf("TestURL Relative case %v: expected %v got %v\n", i, exp, got)
		}
	}
}
