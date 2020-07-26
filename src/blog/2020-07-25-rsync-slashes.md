---
description: rsync and significant traiing slashes
title: rsync slashes
---

### _rsync_ and slashes

the trailing slash for directories is important for rsync,
get it wrong and your directory becomes nested 10 levels deep at your description

#### _basics_

- `dir`
  - source: include everything in this directory, including the directory
  - destination: place everything in this directory, create if it doesn't exist
- `dir/`
  - source: include everything in this directory, excluding the directory
  - destination: place everything in this directory, create if it doesn't exist

#### _corollary_

- sync directory to destination: `rsync -r src parent`
  - result: `src` == `parent/src`
- sync directory to destination with new name: `rsync -r src/ dst/new-src`
  - result: `src` == `dst/new-src`
- sync directory to existing destination: `rsync -r src/ dst/src`
  - result: `src` == `dst/src`
