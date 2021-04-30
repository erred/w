otel hostmetrics

- Used
- Available

prometheus/client_golang

- rss
- vsize

google/cadvisor

- cache == total_cache || file
- container_memory_rss == total_rss || anon
- container_memory_swap
- container_memory_usage_bytes
- container_memory_working_set_bytes	== usage - inactive_file

----

otel hostmetrics

user
system

nice + iowait + irq + softirq + steal + guest + guestnice

cadvisor

container_cpu_system_seconds_total
container_cpu_usage_seconds_total
container_cpu_user_seconds_total

### _memory_

Usually when I have to think about this, I have 2 questions:
- how much memory is my process / container / system using
- how much memory does it actually need

Memory can be thought of in several tiers:
- virtual memory: everything you could use, allocated on write, doesn't have to be mapped to anything
- resident set: everything that is currently mapped
- working set: everything that you need right now


Important to note,
there's virtual memory, which is free ~~infinite~~ space that can be used,
backed by physical memory, the kind you have to buy.

#### _collecting_ data
