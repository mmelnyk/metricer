# Example of usage
TODO

## Prometheus
TODO

prometheus.yml
```
...
scrape_configs:
  - job_name: 'metricer'
    metrics_path: /metrics/values
    static_configs:
    - targets: ['localhost:9110']
...
```

## Telegraf
TODO

telegraf.conf
```
```
