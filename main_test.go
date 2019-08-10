package main

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

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

func TestParseGomod(t *testing.T) {
	s := `
[[ github ]]
        mod = "test1"
[[ github ]]
        user = "seankhliao"
        mod = "test2"
[[ github ]]
        user = "erred"
        mod  = "test3"
`
	gm := Gomod{
		Github: []GithubModule{
			{
				Mod: "test1",
			},
			{
				User: "seankhliao",
				Mod:  "test2",
			}, {
				User: "erred",
				Mod:  "test3",
			},
		},
	}
	g := &Gomod{}
	_, err := toml.Decode(s, g)
	if err != nil {
		t.Errorf("Parse Gomod.toml %v", err)
		return
	}
	if !reflect.DeepEqual(g, gm) {
		t.Errorf("Parse Gomod.toml not DeepEqual: expected \n%#v\n, got \n%#v\n", gm, g)
	}
}

func TestGithubModule(t *testing.T) {
	type TC struct {
		m        GithubModule
		module   string
		repo     string
		goimport string
		gosource string
	}
	cases := []TC{
		{
			GithubModule{
				User: "seankhliao",
				Mod:  "test1",
			},
			"seankhliao.com/test1",
			"https://github.com/seankhliao/test1",
			"seankhliao.com/test1 git https://github.com/seankhliao/test1",
			`seankhliao.com/test1
https://github.com/seankhliao/test1
https://github.com/seankhliao/test1/tree/master{/dir}
https://github.com/seankhliao/test1/tree/master{/dir}/{file}#L{line}`,
		},
		{
			GithubModule{
				Mod: "test2",
			},
			"seankhliao.com/test2",
			"https://github.com/seankhliao/test2",
			"seankhliao.com/test2 git https://github.com/seankhliao/test2",
			`seankhliao.com/test2
https://github.com/seankhliao/test2
https://github.com/seankhliao/test2/tree/master{/dir}
https://github.com/seankhliao/test2/tree/master{/dir}/{file}#L{line}`,
		},
		{
			GithubModule{
				User: "erred",
				Mod:  "test3",
			},
			"seankhliao.com/test3",
			"https://github.com/erred/test3",
			"seankhliao.com/test3 git https://github.com/erred/test3",
			`seankhliao.com/test3
https://github.com/erred/test3
https://github.com/erred/test3/tree/master{/dir}
https://github.com/erred/test3/tree/master{/dir}/{file}#L{line}`,
		},
		{
			GithubModule{
				Mod: "sub/module",
			},
			"seankhliao.com/sub/module",
			"https://github.com/seankhliao/sub",
			"seankhliao.com/sub/module git https://github.com/seankhliao/sub",
			`seankhliao.com/sub/module
https://github.com/seankhliao/sub
https://github.com/seankhliao/sub/tree/master{/dir}
https://github.com/seankhliao/sub/tree/master{/dir}/{file}#L{line}`,
		},
		{
			GithubModule{
				Mod: "test5/v5",
			},
			"seankhliao.com/test5/v5",
			"https://github.com/seankhliao/test5",
			"seankhliao.com/test5/v5 git https://github.com/seankhliao/test5",
			`seankhliao.com/test5/v5
https://github.com/seankhliao/test5
https://github.com/seankhliao/test5/tree/master{/dir}
https://github.com/seankhliao/test5/tree/master{/dir}/{file}#L{line}`,
		},
		{
			GithubModule{
				User: "erred",
				Mod:  "test6/v6",
			},
			"seankhliao.com/test6/v6",
			"https://github.com/erred/test6",
			"seankhliao.com/test6/v6 git https://github.com/erred/test6",
			`seankhliao.com/test6/v6
https://github.com/erred/test6
https://github.com/erred/test6/tree/master{/dir}
https://github.com/erred/test6/tree/master{/dir}/{file}#L{line}`,
		},
	}
	for i, c := range cases {
		c.m.Parse()
		if c.m.Module != c.module {
			t.Errorf("GithubModule.Module case %v: expected %v got %v\n", i, c.module, c.m.Module)
		}
		if c.m.Repo != c.repo {
			t.Errorf("GithubModule.Repo case %v: expected %v got %v\n", i, c.repo, c.m.Repo)
		}
		if c.m.GoImport != c.goimport {
			t.Errorf("GithubModule.GoImport case %v: expected \n%v got \n%v\n", i, c.goimport, c.m.GoImport)
		}
		if c.m.GoSource != c.gosource {
			t.Errorf("GithubModule.GoSource case %v: expected \n%v got \n%v\n", i, c.gosource, c.m.GoSource)
		}

	}
}
