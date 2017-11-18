# sym-bidder

для запуска:
```$: docker-compose up --build```

тестировал с помощью [vegeta](https://github.com/tsenart/vegeta)

```$: vegeta attack -targets=target.txt -duration=10s -rate=10000 > results.bin```

```$: vegetreport -inputs results.bin 
Requests      [total, rate]            99904, 9990.50
Duration      [total, attack, wait]    10.004399384s, 9.999899799s, 4.499585ms
Latencies     [mean, 50, 95, 99, max]  4.136696ms, 434.219µs, 16.766372ms, 72.088714ms, 217.81598ms
Bytes In      [total, mean]            999040, 10.00
Bytes Out     [total, mean]            325437285, 3257.50
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:99904 
```
