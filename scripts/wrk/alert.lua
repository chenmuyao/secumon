-- wrk -t2 -d30s -c10 -s ./scripts/wrk/alert.lua http://localhost:8989/alerts
-- Running 30s test @ http://localhost:8989/alerts
--   2 threads and 10 connections
--   Thread Stats   Avg      Stdev     Max   +/- Stdev
--     Latency     1.46ms    2.80ms  39.28ms   94.72%
--     Req/Sec     5.58k     1.98k    9.41k    62.50%
--   333203 requests in 30.01s, 472.20MB read
-- Requests/sec:  11102.83
-- Transfer/sec:     15.73MB

wrk.method = "GET"
