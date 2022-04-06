# intership-task-go

Тестовое задание на вакансию Go-разработчик.

## Запуск и остановка
Запуск:
``` sh
make
```
Остановка:

``` sh
make clean
```


## Benchmarking
Четыре датчика, интервал обновления 5 секунд, таймаут одна секунда:
``` sh
> wrk -t10 -c500 http://localhost:8080
Running 10s test @ http://localhost:8080
  10 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.96ms    4.77ms  81.15ms   81.02%
    Req/Sec     9.17k     1.31k   22.47k    74.40%
  918092 requests in 10.10s, 104.19MB read
Requests/sec:  90920.98
Transfer/sec:     10.32MB
```

Десять датчиков, интервал обновления 1 секунд, таймаут 20 секунд:
``` sh
> wrk -t10 -c5000 -d30s http://localhost:8080
Running 30s test @ http://localhost:8080
  10 threads and 5000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.21ms    8.08ms 191.28ms   82.98%
    Req/Sec    11.03k     7.11k   34.53k    74.35%
  2636466 requests in 30.09s, 304.14MB read
  Socket errors: connect 3989, read 0, write 0, timeout 0
Requests/sec:  87622.83
Transfer/sec:     10.11MB
```
