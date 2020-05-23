---
description: relational databases
title: db relational
---

### _databases_

> People keep asking me career advice, here is one:
>
> Make state someone else's problem.
> &mdash Jaana Dogan (@rakyll) [2020-05-19](https://twitter.com/rakyll/status/1262599727394074626)

[jepsen.io](http://jepsen.io/analyses) analyses arre imformative

- [2020-06-02 | db relational](../2020-06-01-db-relational/)
- [2020-06-02 | db key value](../2020-06-02-db-key-value/)
- [2020-06-03 | db document](../2020-06-03-db-document/)
- [2020-06-04 | db wide](../2020-06-04-db-wide/)
- [2020-06-05 | db timeseries](../2020-06-05-db-timeseries/)
- [2020-06-06 | db message queue](../2020-06-06-db-message-queue/)

#### _basics_

common-ish databases you might think about using

slight focus on distributed databases

##### _isolation_ and consistency

Read [isolation](http://dbmsmusings.blogspot.com/2019/05/introduction-to-transaction-isolation.html)
and [consistency](http://dbmsmusings.blogspot.com/2019/07/overview-of-consistency-levels-in.html),
and [this](https://fauna.com/blog/a-comparison-of-scalable-database-isolation-levels)
and [problems that may occur](http://dbmsmusings.blogspot.com/2019/06/correctness-anomalies-under.html) (the last table is a simple overview)
tldr: it's confusing, moreso when transactions are involved

Simple way to think about it,
add external consistency and stric serializable above serializable in isolation

- **external consistency**: behave as if a single machine, regardless of replica
- **strict serializable**: add realtime ordering constraint on serial order
- **serializable**: behave as if transactions occured in _some_ serial order
- **snapshot (repeatable read)**: view a snapshot taken when transaction begins
- **read committed**: view things after they have committed,
- **read uncommitted**: view as things change, even if they will get rolled back

##### _CAP_

Choose 2, only google can have _CA_

- Consistency
- Availability
- Partitions

#### _relational_ databases

##### _embedded_ relational

not much choice here

| database           | Language | SQL syntax | Isolation (default - max) | License       |
| ------------------ | -------- | ---------- | ------------------------- | ------------- |
| _[SQLite][sqlite]_ | C        | sqlite     | snapshot                  | Public Domain |

[sqlite]: https://sqlite.org/src/doc/trunk/README.md

notes:

- Sqlite needs WAL (write ahead log) for snapshot
- Sqlite has no isolation in a single connection

##### _single_ node relational

postgres is usually recommended

| database                 | Language | SQL syntax | Isolation (default - max)     | License            |
| ------------------------ | -------- | ---------- | ----------------------------- | ------------------ |
| [MariaDB][mariadb]       | C/C++    | mysql      | snapshot - serializable       | GPLv2              |
| _[PostgreSQL][postgres]_ | C        | postgres   | read committed - serializable | PostgreSQL License |

[mariadb]: https://github.com/MariaDB/server
[postgres]: https://github.com/postgres/postgres

notes:

- PostgreSQL License is similar to MIT / BSD
- some of the multi node ones below can be run in single node mode too

##### _multi_ node relational

| database                   | Language | SQL syntax | Operation modes | CAP  | Isolation (default - max) | Consensus      | License                 |
| -------------------------- | -------- | ---------- | --------------- | ---- | ------------------------- | -------------- | ----------------------- |
| _[CockroachDB][cockroach]_ | Go       | postgres   | cluster         | CP   | serializable              | Raft           | Business Source License |
| [dqlite][dqlite]           | C        | sqlite     | embedded        | CP   |                           | Raft           | LGPLv3                  |
| [Galera][galera]           | C/C++    | mysql      | cluster         | CP   | snapshot - serializable   | master-master  | GPLv2                   |
| [rqlite][rqlite]           | C/Go     | sqlite     | cluster         | CP   |                           | Raft           | MIT                     |
| _[Spanner][spanner]_       | ?        | spanner    | hosted          | _CA_ | _external consistency_    | Paxos          | -                       |
| [TiDB][tidb]               | Go       | mysql      | cluster         | CP   | snapshot                  | Raft           | Apache 2.0              |
| [Vitess][vitess]           | Go       | mysql      | cluster         | -    | read committed            | 2 Phase Commit | Apache 2.0              |
| [YugaByte DB][yugabyte]    | C/C++    | postgres   | cluster         | CP   | snapshot - serializable   | Raft           | Apache 2.0              |

[cockroach]: https://github.com/cockroachdb/cockroach
[dqlite]: https://github.com/canonical/dqlite
[galera]: https://github.com/codership/galera
[rqlite]: https://github.com/rqlite/rqlite
[spanner]: https://cloud.google.com/spanner
[tidb]: https://github.com/pingcap/tidb
[vitess]: https://github.com/vitessio/vitess
[yugabyte]: https://github.com/yugabyte/yugabyte-db

notes:

- CockroachDB and yugabyte db both claim to be modelled after spanner
- CockroachDB Business Source License converts to Apache 2.0 3 years after release
- dqlite is built on sqlite
- Galera is the clustering mode for mysql / mariadb
- Galera is serializable in master-slave, snapshot in multi master
- rqlite is built on sqlite
- Spanner is technically CP
- TiDB is built on TiKV and RocksDB
- TiDB is designed for sharded data
- Vitess is designed for sharded data
- YugaByte universe = primary cluster + read clusters
- YugaByte cluster = multi node auto sharded
- YugaByte is serializable in YSQL (sql), snapshot in YCQL (nosql)
