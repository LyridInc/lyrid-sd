# lyrid-sd
Lyrid Prometheus Service Discovery

Full Readme - Coming Soon

File service discovery implementation for the solution. Reads data from the proxy (with HTTP endpoint) and automatically generates a new metrics port endpoint  that represents a metric that is coming from prom2lyrid.

### Build and Run the Container
```
docker build .
docker run --restart always -d --name lyrid_sd -p 8000:8000 -p 8001-9024:8001-9024 -v /mnt/prometheus/config/lyrid:/lyrid-sd/.config  $tag
```

### Prometheus Settings

For prometheus to read into the lyrid
```
  - job_name: 'lyrid-service-discovery'
    scrape_interval:    30s
    file_sd_configs:
    - files:
      - '/etc/prometheus/lyrid/lyrid_sd.json'

```
This file is constantly updated by lyrid-sd and any new service added to the proxy will generate a new endpoint port in the lyrid-sd.

Example of lyrid_sd.json:
```json
[
    {
        "targets": [
            "127.0.0.1:8001"
        ],
        "labels": {
            "__meta_lyrid_id": "4e8687c6-8ab6-4371-a98b-789daaf8b111",
            "__meta_lyrid_port": "8001"
        }
    }
]
```

you can relabel things into the prometheus by adding relabel_configs into the job. Example where we add a machine cluster label and pass it to the prometheus:
```
  - job_name: 'lyrid-service-discovery'
    scrape_interval:    30s
    file_sd_configs:
    - files:
      - '/etc/prometheus/lyrid/lyrid_sd.json'
    relabel_configs:
      - source_labels: [ '__meta_lyrid_cluster' ]
        regex: '(.*)'
        target_label: 'cluster'
        action: replace
        replacement: '${1}'
```