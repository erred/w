---
description: consensus isn't that hard
title: tldr consensus
---

### _consensus_

#### _basic_ paxos

Produces a single value through a 2 phase process.
Higher proposal IDs win leader position,
first value to get accepted wins getting propagated.

Relies on singular leader/proposer to not livelock:
compete repeatedly with higher IDs before accepts are sent out by old leader.

Try to get a majority with your new, higher proposal ID.
If a previous value has already been accepted,
you are now responsible for propagating it.
Else, you can get your value to be accepted.
Acceptors send accepted value to learners,
learners agree on value when a majority is received.

#### _multi_ paxos

Reuse phase 2,
include round number with value,
repeat propagating new values with increasing round numbers.

#### _raft_

Single leader, relies a lot on time:

```txt
broadcastTime << electionTimeout << meanTimeBetweenFailure
```

Leader adds to own log, broadcasts to followers.
On majority ack, broadcast commit.

Leader logs are append only,
followers will rewrite their own if it doesn't match leader's.
