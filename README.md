# URL Checker

Checks URLs and tracks metrics via Prometheus and Grafana.

## Architecture

```
Go app (2112) → scraped by → Prometheus (9090) → queried by → Grafana (3000)
```

- **2112** — Go app runs the URL checker and exposes `/metrics` (raw Prometheus text format)
- **9090** — Prometheus scrapes `:2112/metrics` every 15s and stores time-series data
- **3000** — Grafana connects to Prometheus as a data source and renders dashboards

## Running

```bash
docker-compose up
```

- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Metrics endpoint: http://localhost:2112/metrics
