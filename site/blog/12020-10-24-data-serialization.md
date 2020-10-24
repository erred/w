---
description: data for a wire
title: data serialization
---

### _serialization_

The 2 main axes of serialization: text/binary and schema/self-describing

#### _widely_ supported

- json: text, self-describing
  - json schema, jsonrpc, openapi, ... for schemas
- protobuf: binary, schema
  - prototext for text representation
  - protojson for mapping to json

#### _seriously_ consider

- cbor: binary, self-describing
- flatbuffers: like protobuf, but more optimized for storage than transmission
