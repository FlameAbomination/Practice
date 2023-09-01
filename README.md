# L0-NatsStreaming

Сервис, реализующий кэш для базы данных с помощью встроенных методов Go.
Код использует библиотеки "pgx" и "nats-io/stan.go". 
Чтение происходит только из кэша, база данных используется только для наполнения кэша при запуске программы.

# Запуск
"L0.exe" - запуск программы в режиме получения данных из NATS и http.
"L0.exe --publisher путь-к-данным" - запуск программы в режиме отправки данных в NATS.

# Vegeta
echo "GET http://172.16.68.21:8000/b563feb7b2b84b7test" | vegeta attack -duration=5s --output results.bin | vegeta report results.bin  
Requests      [total, rate, throughput]  250, 50.20, 50.20  
Duration      [total, attack, wait]      4.9798511s, 4.9797053s, 145.8µs  
Latencies     [mean, 50, 95, 99, max]    64.238µs, 0s, 248µs, 1.1803ms, 1.6312ms  
Bytes In      [total, mean]              210000, 840.00  
Bytes Out     [total, mean]              0, 0.00  
Success       [ratio]                    100.00%  
Status Codes  [code:count]               200:250  
Error Set:  

# WRK  
go-wrk -c 80 -d 5 http://172.16.68.21:8000/b563feb7b2b84b6test  
Running 5s test @ http://172.16.68.21:8000/b563feb7b2b84b6test  
  80 goroutine(s) running concurrently  
219775 requests in 4.848323943s, 195.34MB read  
Requests/sec:           45330.10  
Transfer/sec:           40.29MB  
Avg Req Time:           1.764831ms  
Fastest Request:        0s  
Slowest Request:        49.3048ms  
Number of Errors:       0   
