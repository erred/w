---
description: notes on reading through the OpenTelemetry specification
title: reading OTel spec
---

### _OpenTelemetry_

The [OpenTelemetry Spec](https://github.com/open-telemetry/opentelemetry-specification)
recently hit [v1.0.1](https://github.com/open-telemetry/opentelemetry-specification/tree/v1.0.1)
(`v1.0.0` with a fix soon after),
so what did I learn?

#### _notes_

The spec itself is reasonably approachable,
if a bit annoyingly spread out over multiple files.

Recommended reading is to start with the
[glossary](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.0.1/specification/glossary.md)
then onto the [overview](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.0.1/specification/overview.md).
[library layout](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.0.1/specification/library-layout.md)
is good for understanding where in general you should be looking for things.
The [semantic conventions](https://github.com/open-telemetry/opentelemetry-specification/tree/v1.0.1/semantic_conventions)
are also worth looking over for the standardized names.

The idea is:
there's a standard library interface that is implemented in roughly the same way across multiple languages.
Client code (apps you write) that is being will always (and only) interact with `api`.
`sdk` will provide a default implementation,
which needs to be instantiated at the start of the program and wired to `api`,
and is responsible for exporting the data.
There are multiple transport formats to support different vendors,
including a new one specifically for OpenTelemetry,
hello [XKCD Standards](https://xkcd.com/927/).
