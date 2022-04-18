### Requirements

·  A service written in python or golang that queries 2 urls (https://httpstat.us/503 & https://httpstat.us/200)

·  The service will check the external urls (https://httpstat.us/503 & https://httpstat.us/200 ) are up (based on http status code 200) and response time in milliseconds

·  The service will run a simple http service that produces metrics using appropriate Prometheus libraries and outputs on /metrics

o Expected response format:

§ sample_external_url_up{url="https://httpstat.us/503 "}  = 0

§ sample_external_url_response_ms{url="https://httpstat.us/503 "}  = [value]

§ sample_external_url_up{url="https://httpstat.us/200 "}  = 1

§ sample_external_url_response_ms{url="https://httpstat.us/200 "}  = [value]


### Looking for:

o Code in python or golang

o Dockerfile to build image

o Kubernetes Deployment Specification to deploy Image to Kubernetes Cluster

o Unit Tests

o Good readme providing instructions for use

o Screen shot of deployment in Prometheus and accompanying Grafana Dashboard 

## go-mon
Go Application to expose custom metrics that checks the status of external urls and capture the response time in miliseconds with output in prometheus format e.g. sample_external_url_up{url="https://httpstat.us/503"} 0

### Prerequisite
- Go 1.15 or above workspace with GO PATH


## Build and Deploy on desktop

`$ git clone https://github.com/lambdafunc/vmw.git`
from root directory run: `$ go mod tidy`
```
$ go test .
ok  	go_mon/go_mon	1.182s
```
`$ cd go_mon && go run main.go collector.go`
Service should start and listening on port 8080
image-1


### Build container image for go_mon

`Dockerfile` has been provided which can be used to generate container image
Run from the root directory
`$ docker build -t go-mon .`

### Run Prometheus and Grafana

- Promotheus
`$ cp prometheus.yml /tmp`
`$ docker run -d -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus`

Prometheus will be up and listening on port 9090.

- Grafana
`$ docker run -d -p 3000:3000 grafana/grafana grafana`
Grafana will be up and running on localhost:3000

- Run Go-Mon
`$ docker run -d -p 8080:8080 go_mon`

Below is sample run for `metrics` endpoint to verify if metrics are being populated.
```
$ curl localhost:8080/metrics

# TYPE sample_external_url_response_ms counter
sample_external_url_response_ms{url="https://httpstat.us/200"} 435
sample_external_url_response_ms{url="https://httpstat.us/503"} 435
# HELP sample_external_url_up shows if url is up
# TYPE sample_external_url_up counter
sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

URL's
```
Go-Mon Application: http://localhost:8080/metrics
Prometheus: http://localhost:9090/graph
Grafana: http://localhost:3000/
```

## Generate some test data
`$ chmod +x generateResponseData.sh`
`./generateResponseData` to generate some metrics


## Deployment on Kubernetes:
`$ cd kubernetes-deployment`

Deployment specs for Grafana, Prometheus and service exist under this directory. Each deployment spec consist of all necessary resources such `deployment`, `containers`, `namespace` and `service` definitions. Container image can be pushed to any container registry(dockerhub, ECR etc) and service spec needs to be updated accordingly. (line-26 kubernetes-deployment/service-monitoring-deployment.yaml).

#Deployment-spec for prometheus
`$ kubectl apply -f prometheus-deployment.yaml`

#Deployment-spec for prometheus for grafana
`$ kubectl apply -f grafana-deployment.yaml`

Screenshot:
Image