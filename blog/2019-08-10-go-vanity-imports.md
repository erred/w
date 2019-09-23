vanity imports are vain

---

So how do you get `go get your.domain/package` to work?

# _go get_

`go get` queries `your.domain/path/to/module/path/to/package?go-get=1`
and expects a meta tag of the form: `<meta name="go-import" content="module_name vcs repo root">`,
it then expects the path after module name to be the dir structure in the repo,

```
ex:
<meta name="go-import" content="seankhliao.com/rsssubsbot git https://github.com/seankhliao/rsssubsbot">

ex:
# go get GETs:
your.domain/path/to/module/path/to/package?go-get=1

# server responds with:
<meta name="go-import" content="your.domain/path/to/module git https://github.com/your/repo">

# go downloads (git clones):
https://github.com/your/repo

# go expects to find the package in
repo/path/to/package
```

**side note:
the `go` cmd as of go1.12 does not use tls1.3 by default
(need `GODEBUG=tls13=1`)
so make sure your server works with at least tls1.2
if you want stuff like [godoc.org](godoc.org) to work**

# _go source_

Documentation is important,
[godoc.org](godoc.org) is standard for _Go_

support it with the extra meta tag of the form: `<meta name="go-source" content="module_name link_to_repo link_to_repo_ref_dir link_to_repo_ref_line"`

```
ex:
<meta name="go-source" content="seankhliao.com/rsssubsbot
                                https://github.com/seankhliao/rsssubsbot
                                https://github.com/seankhliao/rsssubsbot/tree/master{/dir}
                                https://github.com/seankhliao/rsssubsbot/tree/master{/dir}/{file}#L{line}">
```

# other

so all that's for the `go` cmd,
what about humans?

Another meta tag!
With the form `<meta http-equiv="Refresh" content="delay_in_seconds; url='url.to/redirect/to'">`,
another way is with js: `window.location='url.to/redirect/to'`,
or http redirects if you can parse the `go-get=1` query param
Redirect to `godoc.org`:

```
<meta http-equiv="Refresh" content="0; url='https://godoc.org/seankhliao.com/rsssubsbot'">
```

# everything _together_

```
<meta name="go-import" content="seankhliao.com/rsssubsbot git https://github.com/seankhliao/rsssubsbot">
<meta name="go-source" content="seankhliao.com/rsssubsbot
                                https://github.com/seankhliao/rsssubsbot
                                https://github.com/seankhliao/rsssubsbot/tree/master{/dir}
                                https://github.com/seankhliao/rsssubsbot/tree/master{/dir}/{file}#L{line}">
<meta http-equiv="Refresh" content="0; url='https://godoc.org/seankhliao.com/rsssubsbot'">
```
