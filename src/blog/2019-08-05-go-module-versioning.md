--- title
go module versioning.md
--- description
notes on go using go modules
--- main


package management in Go,
seems like it's _finally_ comming together

_1 module_ per repo,
everything in a module is versioned together

whatever the _module name_ is in `go.mod` is the name that is used when importing it

# _using_

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

# _updating_

or just delete `go.mod` and `go.sum` and have it recalculate all dependencies

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

# _releasing_

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

# _multi module_

realm of confusion,
avoid if possible

## parallel

no conflict in module scope,
versioning may be be recorded s `pseudo versions` due to only a single ambiguous VCS tag

```
root/
  |- package1/
  |   `- gp.mod
  `- package2/
      `- go.mod
```

## subdir

the root module includes everything except stuff in subdir,
subdir has its _own module name_ and versioning,
versions may be recorded as `pseudo versions` created from time and git commits instead of tags

```
root/
  |- go.mod
  `- subdir
      `- go.mod
```
