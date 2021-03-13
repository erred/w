---
description: short circuits instead of if/else
title: short circuit logic
---

### _short_ circuit

Sometimes you are handed the hammer that is `&&` and `||`,
and no proper `if/else`.

#### if...else

```python
if x:
  a
else:
  b
```

```sh
(x && (a || true)) || b
```

#### if...else if...else

```python
if x:
  a
else if y:
  b
else:
  c
```

```sh
(x && (a || true)) || (y && (b || true)) || c
```

#### nested if

```python
if x:
  if y:
    a
  else:
    b
else:
  b
```

```sh
(x && (y && (a || true) || (b || true))) || c
```
