## hw2
### Итоги
Обратный индекс дает наибольшую производительность, поскольку оптимизирован под поиск по префиксу
Btree индекс, дает увеличение производительности только в 2 раза
### Как повторить
```shell
make compose-rebuild
```
потом дождаться загрузки тестовых данных(сделить за логами контейнера cli):
```shell
{"level":"info","time":"2023-06-20T18:20:15Z","message":"import users started"}
{"level":"info","time":"2023-06-20T18:20:41Z","message":"imported 100000"}
{"level":"info","time":"2023-06-20T18:21:08Z","message":"imported 200000"}
{"level":"info","time":"2023-06-20T18:21:35Z","message":"imported 300000"}
{"level":"info","time":"2023-06-20T18:22:07Z","message":"imported 400000"}
{"level":"info","time":"2023-06-20T18:22:35Z","message":"imported 500000"}
{"level":"info","time":"2023-06-20T18:23:04Z","message":"imported 600000"}
{"level":"info","time":"2023-06-20T18:23:33Z","message":"imported 700000"}
{"level":"info","time":"2023-06-20T18:24:00Z","message":"imported 800000"}
{"level":"info","time":"2023-06-20T18:24:30Z","message":"imported 900000"}
{"level":"info","time":"2023-06-20T18:24:56Z","message":"imported 1000000"}
{"level":"info","time":"2023-06-20T18:24:56Z","message":"import users finished"}
```
потом:
```shell
make test-hw-2
```
Будет виден результат с индексами, если хотим без, то:
```sql
DROP INDEX fs_name_idx;
DROP INDEX fs_name_gin_idx;
```
### График запросов без индексов
```shell
wrk -t1 -c1 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   103.79ms   95.23ms 687.06ms   90.08%
    Req/Sec    12.72      5.37    20.00     58.17%
  349 requests in 30.03s, 41.98MB read
Requests/sec:     11.62
Transfer/sec:      1.40MB
wrk -t10 -c10 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   700.86ms  206.91ms   1.37s    66.19%
    Req/Sec     0.93      0.46     3.00     78.96%
  423 requests in 30.07s, 25.92MB read
Requests/sec:     14.07
Transfer/sec:      0.86MB
wrk -t100 -c100 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  100 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.06s     3.23s   23.80s    79.25%
    Req/Sec     0.00      0.00     0.00    100.00%
  422 requests in 30.05s, 23.08MB read
Requests/sec:     14.04
Transfer/sec:    786.34KB
```
### Explain запроса без индекса 
```sql
EXPLAIN SELECT id, role, first_name, second_name, birthdate, biography, city
FROM users
WHERE first_name LIKE 'И%' AND second_name LIKE 'Ф%';
```
```sql
Gather  (cost=1000.00..34068.40 rows=314 width=110)
  Workers Planned: 2
  ->  Parallel Seq Scan on users  (cost=0.00..33037.00 rows=131 width=110)
        Filter: (((first_name)::text ~~ 'И%'::text) AND ((second_name)::text ~~ 'Ф%'::text))
```
### График запросов с btree индексом
```shell
wrk -t1 -c1 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    87.65ms  117.90ms 768.95ms   89.92%
    Req/Sec    20.27     14.76    80.00     79.46%
  521 requests in 30.04s, 69.08MB read
Requests/sec:     17.34
Transfer/sec:      2.30MB
wrk -t10 -c10 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   514.04ms  556.20ms   3.25s    84.95%
    Req/Sec     6.86      8.44    40.00     83.47%
  760 requests in 30.07s, 85.90MB read
Requests/sec:     25.27
Transfer/sec:      2.86MB
wrk -t100 -c100 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  100 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.89s     3.17s   17.95s    86.59%
    Req/Sec     1.44      2.50    10.00     85.82%
  1071 requests in 30.07s, 76.43MB read
Requests/sec:     35.62
Transfer/sec:      2.54MB
```
### Explain запроса с btree индексом
```sql
EXPLAIN SELECT id, role, first_name, second_name, birthdate, biography, city
FROM users
WHERE first_name LIKE 'И%' AND second_name LIKE 'Ф%';
```
```sql
Index Scan using fs_name_idx on users  (cost=0.42..28992.47 rows=314 width=110)
  Index Cond: (((second_name)::text ~>=~ 'Ф'::text) AND ((second_name)::text ~<~ 'Х'::text))
  Filter: (((first_name)::text ~~ 'И%'::text) AND ((second_name)::text ~~ 'Ф%'::text))
```
### График запросов с gin индексом
```shell
wrk -t1 -c1 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    87.44ms  154.72ms 950.92ms   88.63%
    Req/Sec    64.71     55.64   250.00     73.37%
  1416 requests in 30.10s, 167.77MB read
Requests/sec:     47.04
Transfer/sec:      5.57MB
wrk -t10 -c10 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   237.52ms  574.85ms   3.87s    93.40%
    Req/Sec    21.16     17.58   130.00     76.48%
  4367 requests in 30.09s, 542.60MB read
Requests/sec:    145.15
Transfer/sec:     18.04MB
wrk -t100 -c100 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
Running 30s test @ http://localhost:8086
  100 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   570.16ms  598.17ms   4.71s    87.95%
    Req/Sec     4.16      4.52    40.00     77.68%
  6288 requests in 30.10s, 508.23MB read
Requests/sec:    208.88
Transfer/sec:     16.88MB
```
### Explain запроса с gin индексом
```sql
EXPLAIN SELECT id, role, first_name, second_name, birthdate, biography, city
FROM users
WHERE first_name LIKE 'И%' AND second_name LIKE 'Ф%';
```
```sql
Bitmap Heap Scan on users  (cost=35.21..1190.42 rows=314 width=110)
  Recheck Cond: (((first_name)::text ~~ 'И%'::text) AND ((second_name)::text ~~ 'Ф%'::text))
  ->  Bitmap Index Scan on fs_name_gin_idx  (cost=0.00..35.14 rows=314 width=0)
        Index Cond: (((first_name)::text ~~ 'И%'::text) AND ((second_name)::text ~~ 'Ф%'::text))

```