---
description: reading another spec OpenMetrics
title: reading OpenMetrics spec
---

### _OpenMetrics_

[OpenMetrics](https://openmetrics.io/)
is the standardization of [prometheus](https://prometheus.io/)'s
data export format, also adding (back) support for protobuf.
It is currently a Standards Track
[Internet Draft](https://datatracker.ietf.org/doc/draft-richih-opsawg-openmetrics/)
(future RFC).

#### _specification_

The [specification](https://github.com/OpenObservability/OpenMetrics/blob/master/specification/OpenMetrics.md)
reads like most RFCs.
The main takeaways should be the naming conventions
especially around data / unit types,
everything else is handled by the libraries you'd usually use.

About the standard itself,
the data produced is a (standalone) snapshot,
which includes all the extra text (`TYPE`, `UNIT`, `HELP`)
which is nice for debugging,
but seems like quite a bit of overhead.
