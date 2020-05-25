---
description: databases you might consider using
title: databases
---

### _databases_

> People keep asking me career advice, here is one:
>
> Make state someone else's problem.
>
> — _Jaana Dogan_ (@rakyll) [2020-05-19](https://twitter.com/rakyll/status/1262599727394074626)

[jepsen.io](http://jepsen.io/analyses) analyses arre imformative

#### _basics_

common-ish databases you might think about using

slight focus on distributed databases

#### _isolation_ and consistency

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

#### _CAP_

Choose 2, only google can have _CA_

- Consistency
- Availability
- Partitions

### _relational_

The standard database.

#### _embedded_ relational

not much choice here

| database           | Language | SQL syntax | Isolation (default - max) | License       |
| ------------------ | -------- | ---------- | ------------------------- | ------------- |
| _[SQLite][sqlite]_ | C        | sqlite     | snapshot                  | Public Domain |

[sqlite]: https://sqlite.org/src/doc/trunk/README.md

##### _notes_

- Sqlite needs WAL (write ahead log) for snapshot
- Sqlite has no isolation in a single connection

#### _single_ node relational

postgres is usually recommended

| database                 | Language | SQL syntax | Isolation (default - max)     | License            |
| ------------------------ | -------- | ---------- | ----------------------------- | ------------------ |
| [MariaDB][mariadb]       | C/C++    | mysql      | snapshot - serializable       | GPLv2              |
| _[PostgreSQL][postgres]_ | C        | postgres   | read committed - serializable | PostgreSQL License |

[mariadb]: https://github.com/MariaDB/server
[postgres]: https://github.com/postgres/postgres

##### _notes_

- PostgreSQL License is similar to MIT / BSD
- some of the multi node ones below can be run in single node mode too

#### _multi_ node relational

| database                   | Language | SQL syntax | Operation modes | CAP  | Isolation (default - max) | Consensus | License                 |
| -------------------------- | -------- | ---------- | --------------- | ---- | ------------------------- | --------- | ----------------------- |
| _[CockroachDB][cockroach]_ | Go       | postgres   | cluster         | CP   | serializable              | Raft      | Business Source License |
| [dqlite][dqlite]           | C        | sqlite     | embedded        | CP   | ?                         | Raft      | LGPLv3                  |
| [Galera][galera]           | C/C++    | mysql      | master-master   | -    | snapshot - serializable   | -         | GPLv2                   |
| [rqlite][rqlite]           | C/Go     | sqlite     | cluster         | CP   | ?                         | Raft      | MIT                     |
| _[Spanner][spanner]_       | ?        | spanner    | hosted          | _CA_ | **external consistency**  | Paxos     | -                       |
| [Vitess][vitess]           | Go       | mysql      | master-master   | -    | read committed            | -         | Apache 2.0              |
| [YugaByte DB][yugabyte]    | C/C++    | postgres   | cluster         | CP   | snapshot - serializable   | Raft      | Apache 2.0              |

[cockroach]: https://github.com/cockroachdb/cockroach
[dqlite]: https://github.com/canonical/dqlite
[galera]: https://github.com/codership/galera
[rqlite]: https://github.com/rqlite/rqlite
[spanner]: https://cloud.google.com/spanner
[vitess]: https://github.com/vitessio/vitess
[yugabyte]: https://github.com/yugabyte/yugabyte-db

##### _notes_

- CockroachDB and yugabyte db both claim to be modelled after spanner
- CockroachDB Business Source License converts to Apache 2.0 3 years after release
- dqlite is built on sqlite
- Galera is the clustering mode for mysql / mariadb
- Galera is serializable in master-slave, snapshot in multi master
- rqlite is built on sqlite
- Spanner is technically CP
- Vitess is designed for sharded data
- YugaByte universe = primary cluster + read clusters
- YugaByte cluster = multi node auto sharded
- YugaByte is serializable in YSQL (sql), snapshot in YCQL (nosql)

### _key_ value

Stores a blob of unindexed data.

#### _embedded_ key value

| database               | Language | Isolation (default - max) | License                 |
| ---------------------- | -------- | ------------------------- | ----------------------- |
| [badger][badger]       | Go       | snapshot                  | Apache 2.0              |
| [bbolt][bbolt]         | Go       | serializable              | MIT                     |
| [LevelDB][leveldb]     | C++      | serializable              | BSD 3                   |
| [goleveldb][goleveldb] | Go       | serializable              | BSD 2                   |
| [LMDB][lmdb]           | C        | serializable              | OpenLDAP Public License |
| [RocksDB][rocksdb]     | C++      | read committed - snapshot | GPLv2                   |

[badger]: https://github.com/dgraph-io/badger
[bbolt]: https://github.com/etcd-io/bbolt
[goleveldb]: https://github.com/syndtr/goleveldb
[leveldb]: https://github.com/google/leveldb
[lmdb]: https://github.com/LMDB/lmdb
[rocksdb]: https://github.com/facebook/rocksdb

##### _notes_

- LevelDB isolation is probably serializable since only a single writer is allowed
- LMDB OpenLDAP Public License is similar to MIT / BSD
- RocksDB actually allows concurrent

#### _single_ node key value

| database       | Language | Isolation (default - max) | License |
| -------------- | -------- | ------------------------- | ------- |
| [Redis][redis] | C        | serializable              | BSD 3   |

[redis]: https://github.com/antirez/redis

##### _notes_

- Redis is primarily an in memory cache, durability is not guaranteed
- Redis core is BSD 3, redis modules are Redis Source Available License
- Redis technically has cluster / replication but it's async (useful for sometimes stale caching)
- Redis isolation is probably serializable

#### _multi_ node key value

| database                     | Language | Operation modes | CAP | Isolation (default - max)                      | Consensus | License    |
| ---------------------------- | -------- | --------------- | --- | ---------------------------------------------- | --------- | ---------- |
| [Consul][consul]             | Go       | cluster         | CP  | strict serializable - **external consistency** | Raft      | MPL 2.0    |
| _[etcd][etcd]_               | Go       | cluster         | CP  | strict serializable - **external consistency** | Raft      | Apache 2.0 |
| [FoundationDB][foundationdb] | C/C++    | cluster         | CP  | Strict serializable                            | Paxos     | Apache 2.0 |
| [ZooKeeper][zookeeper]       | Java     | cluster         | CP  | read uncommitted - read committed              | Zab       | Apache 2.0 |

[consul]: https://github.com/hashicorp/consul
[etcd]: https://github.com/etcd-io/etcd
[foundationdb]: https://github.com/apple/foundationdb
[zookeeper]: https://github.com/apache/zookeeper

##### _notes_

- Consul external consistency requires extra roundtrips
- etcd external consistency requires extra roundtrips

### _document_

Stores structured data (can be indexed).
Can be used like a key-value store.

| database                 | Language        | Operation modes | CAP | Isolation (default - max)       | Consensus | License    |
| ------------------------ | --------------- | --------------- | --- | ------------------------------- | --------- | ---------- |
| [ArangoDB][arangodb]     | C++             | cluster         | CP  | read uncommitted?               | Raft      | Apache 2.0 |
| [Couchbase][couchbase]   | C/C++/Erlang/Go | master-master   | -   | read committed                  | -         | Apache 2.0 |
| [CouchDB][couchdb]       | Erlang          | cluster         | PA  | snapshot                        | -         | Apache 2.0 |
| [Fauna DB][faunadb]      | Scala           | cluster         | CP  | **strict serializable**         | Calvin    | -          |
| _[Firestore][firestore]_ | -               | hosted          | CP  | **external consistency**        | -         | -          |
| [MongoDB][mongodb]       | C++             | cluster         | CP  | read uncommitted - snapshot     | homegrown | SSPL       |
| [RethinkDB][rethinkdb]   | C++             | cluster         | CP  | read uncommitted - serializable | Raft      | Apache 2.0 |

[arangodb]: https://github.com/arangodb/arangodb
[couchbase]: https://github.com/couchbase
[couchdb]: https://github.com/apache/couchdb
[faunadb]: https://fauna.com/
[firestore]: https://firebase.google.com/docs/firestore
[mongodb]: https://github.com/mongodb/mongo
[rethinkdb]: https://github.com/rethinkdb/rethinkdb

##### _notes_

- ArangoDB is multimodal but nobody recommends the other modes
- CouchDB is eventually consistent
- MongoDB loses data by default
- RethinkDB is probably serializable

#### _wide_ column

Stores column data contiguously (as opposed to row data contiguously).
Primary usecase is analytics.

| database               | Language | Operation modes | CAP | Isolation (default - max) | Consensus | License    |
| ---------------------- | -------- | --------------- | --- | ------------------------- | --------- | ---------- |
| [BigTable][bigtable]   | -        | hosted          | CP  | row level atomic          | -         | -          |
| [Cassandra][cassandra] | Java     | cluster         | CA  | row level atomic          | Paxos     | Apache 2.0 |
| [HBase][hbase]         | Java     | cluster         | CP  | read committed            | Paxos     | Apache 2.0 |

[bigtable]: https://cloud.google.com/bigtable
[cassandra]: https://github.com/apache/cassandra
[hbase]: https://github.com/apache/hbase

##### _notes_

- BigTable is single node
- Google internally has Percolator with snapshot isolation built on BigTable
- Cassandra can optionally be CP
- Cassandra only has row level isolation
- HBase requires hadoop

#### _graph_

Stores nodes and edges, super specialized.

| database           | Language | Operation modes | CAP | Isolation (default - max) | Consensus | License    |
| ------------------ | -------- | --------------- | --- | ------------------------- | --------- | ---------- |
| _[Dgraph][dgraph]_ | Go       | cluster         | CP  | snapshot                  | Raft      | Apache 2.0 |
| [Neo4j][neo4j]     | Java     | cluster         | CP  | read committed            | Raft      | GPLv3      |

[dgraph]: https://github.com/dgraph-io/dgraph
[neo4j]: https://github.com/neo4j/neo4j

##### _notes_

- Neo4j cluster also has primary-read only

#### _time_ series

Write once, append only style. No transactions, no isolation concerns.

| database                   | Language | Operation modes | CAP | Consensus | License    |
| -------------------------- | -------- | --------------- | --- | --------- | ---------- |
| [Druid][druid]             | Java     | cluster         | -   | -         | Apache 2.0 |
| _[InfluxDB][influxdb]_     | Go       | cluster         | CA  | Raft      | MIT        |
| [OpenTSDB][opentsdb]       | Java     | cluster         | CP  | Zab       | LGPL 2.1   |
| [TimescaleDB][timescaledb] | C        | single node     | -   | ?         | Apache 2.0 |

[druid]: https://github.com/apache/druid/
[influxdb]: https://github.com/influxdata/influxdb
[opentsdb]: https://github.com/OpenTSDB/opentsdb
[timescaledb]: https://github.com/timescale/timescaledb

##### _notes_

- Druid is a column database but doesn't support writes
- Graphite has a TSDB but cannot be operated independently, also not clustered
- InfluxDB clustering is available as enterprise offering
- OpenTSDB requires HBase, Hadoop
- TimescaleDB is based on postgres
- TimescaleDB clustering is private beta
