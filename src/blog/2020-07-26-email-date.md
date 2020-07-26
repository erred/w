---
description: dates in email
title: email date
---

### _dates_ and email

Email, having been designed in an era
without always on connects from everywhere to anywhere
can be considered somewhat delay tolerant.

Which led me to think about delivering assignments on time.
We might get directions such as:

> The assignment must be received by email by midnight

but if it was worded slightly differently,
such as "sent by email by midnight"
we can play games.

The `Date` header in email is set by the sender,
so you could, in theory, set it to a time in the past after the due date,
send it, then claim there were "issues out of your control" in delivering the mail
and point to the header as proof.

In practice, this date may already be unreliable due to clients not using NTP

No word on which date the main email clients will display though.
