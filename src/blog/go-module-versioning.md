title = Go Module Versioning
date = 2019-08-05
desc = notes on go using go modules

---

package management in Go,
seems like it's _finally_ comming together

_1 module_ per repo,
everything in a module is versioned together

# using

## v0, v1

you can't add `v0` or `v1` to the end even if you wanted

```
import "<module name>"
import "<module name>/subpackage"
```

## v2+

think of every version 2+ as a _new module_

```
import "<module name>/v2"
import "<module name>/v2/subpackage"
```

# updating

## patch

```
go get -u=patch
```

## minor

```
go get -u
```

## major

edit all import paths

```
import "<module name>"            => import "<module name>/v2"
import "<module name>/subpackage" => import "<module name>/v2/subpackage"
```

# releasing

or just don't version and make everyone live on `master`

## v0, v1

```
go mod init <module name>

git tag v0.x.x
git tag v1.x.x
git push --tags
```

## v2+

edit `go.mod`

```
module <module name> => module <module name>/v2
```

```
git tag v2.x.x
git push --tags
```
