-- wrk -t2 -d30s -c10 -s ./scripts/wrk/publish.lua http://localhost:8989/logs
-- Running 30s test @ http://localhost:8989/logs
--   2 threads and 10 connections
--   Thread Stats   Avg      Stdev     Max   +/- Stdev
--     Latency     5.02ms    1.87ms  21.28ms   72.16%
--     Req/Sec     1.01k   262.57     1.75k    69.00%
--   60144 requests in 30.02s, 10.27MB read
-- Requests/sec:   2003.64
-- Transfer/sec:    350.25KB

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body =
	'{"timestamp": "2025-01-08T12:00:00Z","client_ip": "192.168.1.1","endpoint": "/api/v1/resource", "method": "GET", "status_code": 401}'
