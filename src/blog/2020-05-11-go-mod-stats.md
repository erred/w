---
description: stats about go modules
title: go mod stats
---

### _go_ module survey

**Biases**:

- 2020-05-09
- uses the default [index][index] / [proxy][proxy]
- modules only in other proxies not included
- private modules not included
- applications never pulled through proxy not included

[index]: https://index.golang.org/
[proxy]: https://proxy.golang.org/

**Process**:

Download and run [`go/scanner.Scan`][scan] on all available versions
took approx 12h on `Intel(R) Xeon(R) CPU E3-1240L v5 @ 2.10GHz`
(4 core / 8 hyperthread) with 16GB ram (max ~10GB used).

[scan]: https://golang.org/pkg/go/scanner/#Scanner.Scan

#### _basic_

- _173955_ unique modules
- _1184798_ total unique versions

Ever wondered how those versions are distributed?

<iframe width="1415.6560457516339" height="585.0395299145299" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1969100263&amp;format=interactive"></iframe>

#### _availability_

by versions

- _1050881_ 88.7% retrieved successfully
- _132066_ 11.1% 410 Gone
- _1478_ 0.1% modfile parsing errors
- _285_ 503 Service Unavailable (should've retried)
- _88_ stream INTERNAL_ERROR (should've retried)

What kind of parsing errors are there?

- invalid go directives
- `require(` note the missing space
- invalid versions (tags, branches)
- mismatch major versions
- git merge conflict (`>>>>>>`, `======`, `<<<<<<`)

#### _hosting_

##### domains

What domains do modules live under?

- _161519_ github.com
- _2655_ gitlab.com
- _1753_ gopkg.in
- _1338_ bitbucket.org
- _940_ gitee.com
- _277_ k8s.io
- _187_ git.sr.ht
- _140_ sigs.k8s.io
- _139_ gitea.com
- _138_ code.cloudfoundry.org
- _69_ code.gitea.io
- _64_ 山山.xyz
- _53_ golang.org
- _51_ modernc.org
- _45_ rsc.io
- _41_ decred.org
- _39_ moul.io
- _35_ launchpad.net
- _34_ go.elastic.co
- _32_ code.aliyun.com

##### vcs extension

do people use vanity domains or `.$vcs` extension for their repos?

<iframe width="605.6603773584906" height="374.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1597446649&amp;format=interactive"></iframe>

### _latest_

the following stats only count the latest version of every module

#### _modfile_

##### go directive

apparently documentation isn't good enough and people try to use invalid versions

<iframe width="606.8206751054852" height="375.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1565143455&amp;format=interactive"></iframe>

##### require directive

who requires so many things?

todo: split by direct / indirect

<iframe width="1410.4082661290322" height="585" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=733939&amp;format=interactive"></iframe>

##### replace directive

todo: split by replace type

<iframe width="1410.5" height="563.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1094968622&amp;format=interactive"></iframe>

##### exclude directive

the _exclude_ directive feels excluded

<iframe width="605.6603773584906" height="374.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=854001003&amp;format=interactive"></iframe>

#### _code scan_

Since `go.mod` doesn't include all imports
and I neglected to record the imports per package / file
I don't think I can do tokens - imports graph
(and find the left-pad equivalent in go)

##### tokens per module

10000 tokens is not too much, not too little

(the labels should be **between the previous label and x**)

<iframe width="610.5121293800539" height="377.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1499721188&amp;format=interactive"></iframe>

##### token popularity

ever wondered which operator was the most popular?

<iframe width="1417.5" height="586.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=1734309664&amp;format=interactive"></iframe>

##### identifiers per module

how many different identifiers do you need? a-z 26?

<iframe width="612.1293800539083" height="378.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=712316250&amp;format=interactive"></iframe>

##### identifier popularity

err??????

i have some doubts about this since who uses `TRUE`, `FALSE`,`iNdEx`, `dAtA`, `autorest`

is this maybe skewed by vendor?

<iframe width="1413" height="585.5" seamless frameborder="0" scrolling="no" src="https://docs.google.com/spreadsheets/d/e/2PACX-1vSKIZYBdvEdDVWnHnrksmrT99ag-MoTFRg6LpddYWqRlHlYsiF9It8Pq-upL-iyrs1MygyIqFx3DRvY/pubchart?oid=403214845&amp;format=interactive"></iframe>
