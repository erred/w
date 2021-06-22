---
title: load testing tools
description: surveying the landscape of tools for load testing
---


### _load_ testing

You have a HTTP endpoint,
and you want to know when it might fall over.
Of course someone has written tools for this.

Anyway, roll initial impressions

#### _load_

These tools are based around dumping a lot of requests at your server.
Some are obviously better than others.


##### _hey_

[hey](https://github.com/rakyll/hey): Go cli tool

Importantly does the one thing I want:
send as many requests as it can as fast as it can.

```sh
$ hey -n 1000000 http://localhost:8080

Summary:
  Total:	29.6725 secs
  Slowest:	0.0464 secs
  Fastest:	0.0001 secs
  Average:	0.0015 secs
  Requests/sec:	33701.1903


Response time histogram:
  0.000 [1]	      |
  0.005 [971369]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.009 [22573]	  |■
  0.014 [4660]	  |
  0.019 [1076]	  |
  0.023 [238]	    |
  0.028 [65]	    |
  0.033 [13]	    |
  0.037 [4]	      |
  0.042 [0]	      |
  0.046 [1]	      |


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0006 secs
  50% in 0.0012 secs
  75% in 0.0019 secs
  90% in 0.0028 secs
  95% in 0.0037 secs
  99% in 0.0078 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0001 secs, 0.0464 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0231 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0199 secs
  resp wait:	0.0014 secs, 0.0000 secs, 0.0342 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0300 secs

Status code distribution:
  [200]	1000000 responses
```

##### _wrk2_

[wrk2](https://github.com/giltene/wrk2): C cli tool

Based on _wrk_ but with a new (required) throughput parameter (reqs/seq)
and apparently more accurate reporting.
Equal or better performance than `hey`
but i like the output of `hey` a bit better.

```sh
$ wrk2 -R 100000 -d 30s --latency http://localhost:8080
Running 30s test @ http://localhost:8080
  2 threads and 10 connections
  Thread calibration: mean lat.: 2802.377ms, rate sampling interval: 10002ms
  Thread calibration: mean lat.: 2760.434ms, rate sampling interval: 10059ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.89s     3.47s   18.20s    57.62%
    Req/Sec    19.20k   248.00    19.45k    50.00%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%   12.09s
 75.000%   14.87s
 90.000%   16.61s
 99.000%   17.87s
 99.900%   18.14s
 99.990%   18.20s
 99.999%   18.22s
100.000%   18.22s

  Detailed Percentile spectrum:
       Value   Percentile   TotalCount 1/(1-Percentile)

    5824.511     0.000000           13         1.00
    7069.695     0.100000        77838         1.11
    8216.575     0.200000       155282         1.25
    9437.183     0.300000       233028         1.43
   10747.903     0.400000       310697         1.67
   12091.391     0.500000       388294         2.00
   12664.831     0.550000       427436         2.22
   13205.503     0.600000       465808         2.50
   13680.639     0.650000       504621         2.86
   14278.655     0.700000       544086         3.33
   14868.479     0.750000       582142         4.00
   15171.583     0.775000       601792         4.44
   15450.111     0.800000       621392         5.00
   15736.831     0.825000       640325         5.71
   16023.551     0.850000       660073         6.67
   16310.271     0.875000       679201         8.00
   16465.919     0.887500       689453         8.89
   16605.183     0.900000       698765        10.00
   16752.639     0.912500       708242        11.43
   16908.287     0.925000       719130        13.33
   17055.743     0.937500       728927        16.00
   17121.279     0.943750       733300        17.78
   17186.815     0.950000       738233        20.00
   17268.735     0.956250       742727        22.86
   17367.039     0.962500       747481        26.67
   17481.727     0.968750       752249        32.00
   17530.879     0.971875       754395        35.56
   17596.415     0.975000       756705        40.00
   17661.951     0.978125       759392        45.71
   17711.103     0.981250       761563        53.33
   17776.639     0.984375       764530        64.00
   17793.023     0.985938       765187        71.11
   17825.791     0.987500       766650        80.00
   17858.559     0.989062       768095        91.43
   17874.943     0.990625       768877       106.67
   17907.711     0.992188       770252       128.00
   17924.095     0.992969       770854       142.22
   17940.479     0.993750       771338       160.00
   17973.247     0.994531       772190       182.86
   17989.631     0.995313       772752       213.33
   18006.015     0.996094       773280       256.00
   18022.399     0.996484       773685       284.44
   18022.399     0.996875       773685       320.00
   18038.783     0.997266       774032       365.71
   18055.167     0.997656       774347       426.67
   18087.935     0.998047       774743       512.00
   18087.935     0.998242       774743       568.89
   18104.319     0.998437       774907       640.00
   18120.703     0.998633       775085       731.43
   18137.087     0.998828       775454       853.33
   18137.087     0.999023       775454      1024.00
   18137.087     0.999121       775454      1137.78
   18153.471     0.999219       775553      1280.00
   18169.855     0.999316       775816      1462.86
   18169.855     0.999414       775816      1706.67
   18169.855     0.999512       775816      2048.00
   18169.855     0.999561       775816      2275.56
   18169.855     0.999609       775816      2560.00
   18186.239     0.999658       775901      2925.71
   18186.239     0.999707       775901      3413.33
   18202.623     0.999756       776051      4096.00
   18202.623     0.999780       776051      4551.11
   18202.623     0.999805       776051      5120.00
   18202.623     0.999829       776051      5851.43
   18202.623     0.999854       776051      6826.67
   18202.623     0.999878       776051      8192.00
   18202.623     0.999890       776051      9102.22
   18202.623     0.999902       776051     10240.00
   18202.623     0.999915       776051     11702.86
   18202.623     0.999927       776051     13653.33
   18219.007     0.999939       776100     16384.00
   18219.007     1.000000       776100          inf
#[Mean    =    11892.257, StdDeviation   =     3465.307]
#[Max     =    18202.624, Total count    =       776100]
#[Buckets =           27, SubBuckets     =         2048]
----------------------------------------------------------
  1191527 requests in 30.00s, 9.52GB read
Requests/sec:  39718.16
Transfer/sec:    324.92MB
```

##### _wrk_

[wrk](https://github.com/wg/wrk): C cli tool

Cli tool, with optional lua scripting.
You have to think about threads and connections though

```sh
$ wrk -d 2m -t 2 -c 100 --latency http://localhost:8080
Running 2m test @ http://localhost:8080
  2 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.13ms    2.65ms  47.49ms   87.41%
    Req/Sec    29.16k     4.08k   41.59k    69.10%
  Latency Distribution
     50%    1.00ms
     75%    2.88ms
     90%    5.37ms
     99%   12.54ms
  6963860 requests in 2.00m, 55.63GB read
Requests/sec:  57998.74
Transfer/sec:    474.47MB
```

#### _oha_

[oha](https://github.com/hatoo/oha): Rust cli/tui tool

Has a TUI, but the final output is still text.
Fancy but don't see how the ui is useful.

```
$ oha -n 1000000 -c 1000 http://localhost:8080
Summary:
  Success rate:	1.0000
  Total:	29.6182 secs
  Slowest:	0.2919 secs
  Fastest:	0.0001 secs
  Average:	0.0295 secs
  Requests/sec:	33763.0723

  Total data:	7.86 GiB
  Size/request:	8.24 KiB
  Size/sec:	271.76 MiB

Response time histogram:
  0.008 [157163] |■■■■■■■■■■■■■■■■■■■■■■
  0.016 [226111] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.024 [202547] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.032 [122975] |■■■■■■■■■■■■■■■■■
  0.040 [70000]  |■■■■■■■■■
  0.048 [48070]  |■■■■■■
  0.056 [38563]  |■■■■■
  0.064 [30063]  |■■■■
  0.072 [22626]  |■■■
  0.080 [17009]  |■■
  0.088 [64873]  |■■■■■■■■■

Latency distribution:
  10% in 0.0062 secs
  25% in 0.0113 secs
  50% in 0.0204 secs
  75% in 0.0365 secs
  90% in 0.0658 secs
  95% in 0.0896 secs
  99% in 0.1350 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0175 secs, 0.0001 secs, 0.1563 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0051 secs

Status code distribution:
  [200] 1000000 responses
```

TUI:

```
┌Progress─────────────────────────────────────────────────────────────────────────────────────┐
│                                      329421 / 1000000                                       │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
┌statics for last second──────────────────────┐┌Status code distribution──────────────────────┐
│Requests : 35088                             ││[200] 329421 responses                        │
│Slowest: 0.2831 secs                         ││                                              │
│Fastest: 0.0001 secs                         ││                                              │
│Average: 0.0284 secs                         ││                                              │
│Data: 282.42 MiB                             ││                                              │
│Number of open files: 1010 / 65535           ││                                              │
└─────────────────────────────────────────────┘└──────────────────────────────────────────────┘
┌Error distribution───────────────────────────────────────────────────────────────────────────┐
└─────────────────────────────────────────────────────────────────────────────────────────────┘
┌Requests / past second (auto). press -/+/a to┐┌Response time histogram───────────────────────┐
│                ▅▅▅▅▅▅▅ ▆▆▆▆▆▆▆ ███████      ││███████                                       │
│▃▃▃▃▃▃▃ ▅▅▅▅▅▅▅ ███████ ███████ ███████      ││███████ ▇▇▇▇▇▇▇                               │
│███████ ███████ ███████ ███████ ███████      ││███████ ███████                               │
│███████ ███████ ███████ ███████ ███████      ││███████ ███████                               │
│███████ ███████ ███████ ███████ ███████      ││███████ ███████                               │
│███████ ███████ ███████ ███████ ███████      ││███████ ███████ ▅▅▅▅▅▅▅         ▃▃▃▃▃▃▃       │
│█30797█ █31935█ █37635█ █38692█ █39732█      ││█14355█ █12168█ █3582██ █2152██ █2831██       │
│0s      1s      2s      3s      4s           ││0.0170  0.0340  0.0510  0.0681  0.0851        │
└─────────────────────────────────────────────┘└──────────────────────────────────────────────┘
```

##### _apache_ bench

`ab`, part of [apache](https://github.com/apache/httpd): C cli tool

Performance is meh.

```sh
$ ab -n 1000000 -c 10000 http://localhost:8080/
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100000 requests
Completed 200000 requests
Completed 300000 requests
Completed 400000 requests
Completed 500000 requests
Completed 600000 requests
Completed 700000 requests
Completed 800000 requests
Completed 900000 requests
Completed 1000000 requests
Finished 1000000 requests


Server Software:
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        8440 bytes

Concurrency Level:      10000
Time taken for tests:   70.118 seconds
Complete requests:      1000000
Failed requests:        0
Total transferred:      8537000000 bytes
HTML transferred:       8440000000 bytes
Requests per second:    14261.60 [#/sec] (mean)
Time per request:       701.184 [ms] (mean)
Time per request:       0.070 [ms] (mean, across all concurrent requests)
Transfer rate:          118897.72 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0  251  36.5    249     370
Processing:   113  448  51.1    444     740
Waiting:        0  199  38.9    194     426
Total:        347  699  33.2    693     960

Percentage of the requests served within a certain time (ms)
  50%    693
  66%    702
  75%    711
  80%    720
  90%    745
  95%    758
  98%    771
  99%    782
 100%    960 (longest request)
```

##### webslap

[webslap](https://2ton.com.au/webslap/): C cli tool

Based on apache bench / ab, with multi url support
and a terminal animation! / live display.
Does make copying results out a pita though.

```
code    count      min      avg      max     kbhdrs    kbtotal     kbbody
200  1,000,000        0      347    1,313    122,070  8,376,953  8,242,187

URL: http://localhost:8080
200  1,000,000        0      347    1,313    122,070  8,376,953  8,242,187

Time completed:         Tue, 22 Jun 2021 21:25:25 GMT
Concurrency Level:      10,000
Time taken for tests:   35.213s
Total requests:         1,000,000
Failed requests:        0
Keep-alive requests:    990,000
Non-2xx requests:       0
Total transferred:      8,578,000,000 bytes
Headers transferred:    125,000,000 bytes
Body transferred:       8,440,000,000 bytes
Requests per second:    28,398.60 [#/sec] (mean)
Time per request:       352.130 [ms] (mean)
Time per request:       0.035 [ms] (mean, across all concurrent requests)
Wire Transfer rate:     237,893.76 [Kbytes/sec] received
Body Transfer rate:     234,066.59 [Kbytes/sec] received
                    min      avg      max
Connect Time:           0        1      242
Processing Time:        0        0      193
Waiting Time:           0      347    1,313
Total Time:             0      347    1,313

```


##### _siege_

[siege](https://github.com/JoeDog/siege): C cli tool

config and output are both unintuitive

```sh
$ siege -b -q -c 100 -t 1m -j http://localhost:8080

{
	"transactions":			        530945,
	"availability":			        100.00,
	"elapsed_time":			        59.00,
	"data_transferred":		      4273.59,
	"response_time":		        0.01,
	"transaction_rate":		      8999.07,
	"throughput":			          72.43,
	"concurrency":			        99.37,
	"successful_transactions":	530946,
	"failed_transactions":		  0,
	"longest_transaction":		  0.15,
	"shortest_transaction":		  0.00
}
```

##### _vegeta_

[vegeta](https://github.com/tsenart/vegeta) Go cli tool / library

More of a toolsuite, but the ux is a bit confusing?
Actual performance is meh.

```sh
$ echo "GET http://localhost:8080/" | vegeta attack -duration=30s -rate 0 -max-workers 10000 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         334306, 11140.79, 11140.50
Duration      [total, attack, wait]             30.008s, 30.007s, 766.498µs
Latencies     [min, mean, 50, 90, 95, 99, max]  67.87µs, 3.115ms, 1.701ms, 7.585ms, 9.751ms, 15.574ms, 127.55ms
Bytes In      [total, mean]                     2821542640, 8440.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:334306
Error Set:
```


##### _httperf_

[httperf](https://github.com/httperf/httperf): C cli tool

questionable performance

```sh
$ httperf --server localhost --port 8080 --num-conns 1000 --rate 10000 --num-calls 1000000
httperf --client=0/1 --server=localhost --port=8080 --uri=/ --rate=10000 --send-buffer=4096 --recv-buffer=16384 --ssl-protocol=auto --num-conns=1000 --num-calls=1000000
^CMaximum connect burst length: 106

Total: connections 1000 requests 6791480 replies 6790782 test-duration 215.640 s

Connection rate: 4.6 conn/s (215.6 ms/conn, <=1000 concurrent connections)
Connection time [ms]: min 0.0 avg 0.0 max 0.0 median 0.0 stddev 0.0
Connection time [ms]: connect 5.5
Connection length [replies/conn]: 0.000

Request rate: 31494.5 req/s (0.0 ms/req)
Request size [B]: 62.0

Reply rate [replies/s]: min 22054.7 avg 31484.5 max 41094.8 stddev 4416.3 (43 samples)
Reply time [ms]: response 22.0 transfer 9.8
Reply size [B]: header 125.0 content 8440.0 footer 2.0 (total 8567.0)
Reply status: 1xx=0 2xx=6790782 3xx=0 4xx=0 5xx=0

CPU time [s]: user 36.91 system 136.20 (user 17.1% system 63.2% total 80.3%)
Net I/O: 265307.6 KB/s (2173.4*10^6 bps)

Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
```

#### _user_ simulation

These tools are more about simulating complex user flows.
They usually fall over at actually high loads...
Maybe it's better to just run multiple of the pure load generators hitting multiple urls.
And TBH i wouldn't really trust any of these.


#### _locust_

[locust](https://locust.io/): python module

Write python code to simulate requests.
You also get a web ui.
Being python,
you're in charge of running multiple instances to make use of multiple cores,
at least it has a master-worker distributed mode.

```python
from locust import HttpUser, task, between

class LoadTest(HttpUser):
    wait_time = between(0.5, 1)
    host = "http://localhost:8080"

    @task
    def task1(self):
        self.client.get("/1")
        self.client.get("/2")
```

Example run:

```sh
$ locust
[2021-06-22 19:44:44,960] eevee/INFO/locust.main: Starting web interface at http://0.0.0.0:8089 (accepting connections from all network interfaces)
[2021-06-22 19:44:44,972] eevee/INFO/locust.main: Starting Locust 1.5.3
[2021-06-22 19:44:53,693] eevee/INFO/locust.runners: Spawning 1000 users at the rate 10 users/s (0 users already running)...
[2021-06-22 19:45:24,985] eevee/WARNING/root: CPU usage above 90%! This may constrain your throughput and may even give inconsistent response time measurements! See https://docs.locust.io/en/stable/running-locust-distributed.html for how to distribute the load over multiple CPU cores or machines
[2021-06-22 19:46:46,763] eevee/INFO/locust.runners: All users spawned: LoadTest: 1000 (1000 total running)
[2021-06-22 19:49:03,220] eevee/INFO/locust.runners: Stopping 1000 users
[2021-06-22 19:49:03,535] eevee/INFO/locust.runners: 1000 Users have been stopped, 0 still running
[2021-06-22 19:49:03,535] eevee/WARNING/locust.runners: CPU usage was too high at some point during the test! See https://docs.locust.io/en/stable/running-locust-distributed.html for how to distribute the load over multiple CPU cores or machines
KeyboardInterrupt
2021-06-22T19:49:24Z
[2021-06-22 19:49:24,478] eevee/INFO/locust.main: Running teardowns...
[2021-06-22 19:49:24,479] eevee/INFO/locust.main: Shutting down (exit code 0), bye.
[2021-06-22 19:49:24,479] eevee/INFO/locust.main: Cleaning up runner...
 Name                                                          # reqs      # fails  |     Avg     Min     Max  Median  |   req/s failures/s
--------------------------------------------------------------------------------------------------------------------------------------------
 GET /1                                                         93630     0(0.00%)  |     682       1    2505     890  |  374.78    0.00
 GET /2                                                         93630     0(0.00%)  |     600       1    1472     740  |  374.78    0.00
--------------------------------------------------------------------------------------------------------------------------------------------
 Aggregated                                                    187260     0(0.00%)  |     641       1    2505     770  |  749.56    0.00

Response time percentiles (approximated)
 Type     Name                                                              50%    66%    75%    80%    90%    95%    98%    99%  99.9% 99.99%   100% # reqs
--------|------------------------------------------------------------|---------|------|------|------|------|------|------|------|------|------|------|------|
 GET      /1                                                                890   1100   1100   1100   1200   1200   1200   1300   1400   2100   2500  93630
 GET      /2                                                                740    890    960   1000   1100   1100   1200   1300   1400   1500   1500  93630
--------|------------------------------------------------------------|---------|------|------|------|------|------|------|------|------|------|------|------|
 None     Aggregated                                                        770    990   1100   1100   1100   1200   1200   1300   1400   1900   2500 187260
```

##### _k6_

[k6](https://k6.io/): Go/JavaScript cli tool / SaaS

Write js code to simulate requests.
Not sure about the SaaS part.
I guess it works.

```javascript
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('http://localhost:8080');
}
```

```sh
$ k6 run -u 10000 -i 1000000 script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: script.js
     output: -

  scenarios: (100.00%) 1 scenario, 10000 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1000000 iterations shared among 10000 VUs (maxDuration: 10m0s, gracefulStop: 30s)


running (01m35.9s), 00000/10000 VUs, 1000000 complete and 0 interrupted iterations
default ✓ [======================================] 10000 VUs  01m35.9s/10m0s  1000000/1000000 shared iters

     data_received..................: 8.6 GB  89 MB/s
     data_sent......................: 80 MB   834 kB/s
     http_req_blocked...............: avg=7.94ms   min=810ns    med=2.12µs   max=2.43s p(90)=4.04µs   p(95)=6.66µs
     http_req_connecting............: avg=7.91ms   min=0s       med=0s       max=2.43s p(90)=0s       p(95)=0s
     http_req_duration..............: avg=720.12ms min=125.84µs med=615.92ms max=3.67s p(90)=1.26s    p(95)=1.73s
       { expected_response:true }...: avg=720.12ms min=125.84µs med=615.92ms max=3.67s p(90)=1.26s    p(95)=1.73s
     http_req_failed................: 0.00%   ✓ 0       ✗ 1000000
     http_req_receiving.............: avg=6.48ms   min=14.04µs  med=31.12µs  max=3.09s p(90)=231.81µs p(95)=446.44µs
     http_req_sending...............: avg=2.81ms   min=4.25µs   med=8.96µs   max=2.25s p(90)=63.51µs  p(95)=164.56µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s    p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=710.82ms min=75.91µs  med=614.89ms max=3.59s p(90)=1.24s    p(95)=1.68s
     http_reqs......................: 1000000 10427.261551/s
     iteration_duration.............: avg=861.5ms  min=11.91ms  med=662.06ms max=3.86s p(90)=1.61s    p(95)=2s
     iterations.....................: 1000000 10427.261551/s
     vus............................: 10000   min=9712  max=10000
     vus_max........................: 10000   min=10000 max=10000
```


##### _tsung_

[tsung](https://github.com/processone/tsung): Erlang cli tool

Run as proxy, make request using something else and replay.
Didn't try it.

##### _gatling_

[gatling](https://gatling.io/): Scala library / cli toolset?

Run recorder script and replay, or write scala / java code.
Didn't try it.

#### _artillery_

[artillery](https://github.com/artilleryio/artillery): Javascript cli tool, yaml config

Does it even work?

```yaml
config:
  target: http://localhost:8080
  phases:
    - duration: 30
      arrivalCount: 30000
scenarios:
  - flow:
      - get:
          url: "/"
```

```sh
$ artillery run conf.yaml
Started phase 0, duration: 30s @ 21:12:33(+0000) 2021-06-22
Report @ 21:12:43(+0000) 2021-06-22
Elapsed time: 10 seconds
  Scenarios launched:  9984
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: 1002.41
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Errors:
    EAI_AGAIN: 2

Warning:
CPU usage of Artillery seems to be very high (pids: 1)
which may severely affect its performance.
See https://artillery.io/docs/faq/#high-cpu-warnings for details.

Report @ 21:12:53(+0000) 2021-06-22
Elapsed time: 20 seconds
  Scenarios launched:  10000
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: 1002
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Errors:
    EAI_AGAIN: 1
    ETIMEDOUT: 9981

Warning: High CPU usage warning (pids: 1).
See https://artillery.io/docs/faq/#high-cpu-warnings for details.

Report @ 21:13:03(+0000) 2021-06-22
Elapsed time: 30 seconds
  Scenarios launched:  9976
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: 998.5
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Errors:
    ETIMEDOUT: 10000

Report @ 21:13:13(+0000) 2021-06-22
Elapsed time: 40 seconds
  Scenarios launched:  40
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: 4.1
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Errors:
    ETIMEDOUT: 9975

Report @ 21:13:13(+0000) 2021-06-22
Elapsed time: 40 seconds
  Scenarios launched:  0
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: NaN
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Errors:
    ETIMEDOUT: 41

All virtual users finished
Summary report @ 21:13:13(+0000) 2021-06-22
  Scenarios launched:  30000
  Scenarios completed: 0
  Requests completed:  0
  Mean response/sec: 749.63
  Response time (msec):
    min: NaN
    max: NaN
    median: NaN
    p95: NaN
    p99: NaN
  Scenario counts:
    0: 30000 (100%)
  Errors:
    EAI_AGAIN: 3
    ETIMEDOUT: 29997
```
