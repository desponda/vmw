## Service monitoring: Go-Mon

An application written in Golang to expose custom metrics that checks the status of external urls and capture the response time in miliseconds with output in prometheus format e.g. sample_external_url_up{url="https://httpstat.us/503"} 0

* [Requirements](#requirements)
* [Prerequisite](#prerequisite)
* [Build and Deploy on local dev setup](#build-and-deploy-on-local-dev-setup)
    * [Local dev setup and test](#local-dev-setup-and-test)
    * [Build container image for golang app](#build-container-image-for-golang-app)
    * [Run Prometheus and Grafana](#run-prometheus-and-grafana)
    * [Run Golang app](#run-golang-app)
    * [URLs for accessing deployed services](#urls-for-accessing-deployed-services)
* [Deployment on Kubernetes](#deployment-on-kubernetes)
* [Screenshot](#screenshot)

### Requirements

* A service written in python or golang that queries 2 urls (https://httpstat.us/503 & https://httpstat.us/200)
* The service will check the external urls (https://httpstat.us/503 & https://httpstat.us/200 ) are up (based on http status code 200) and response time in milliseconds
* The service will run a simple http service that produces metrics using appropriate Prometheus libraries and outputs on /metrics

* Expected response format:
    ```sh
    sample_external_url_up{url="https://httpstat.us/503 "}  = 0
    sample_external_url_response_ms{url="https://httpstat.us/503 "}  = [value]
    sample_external_url_up{url="https://httpstat.us/200 "}  = 1
    sample_external_url_response_ms{url="https://httpstat.us/200 "}  = [value]
    ```


### Prerequisite
- Go 1.15 or above workspace with GO PATH configured.
- Docker and minikube


## Build and Deploy on local dev setup

### Local dev setup and test
Local development will require `go` binary installed and workspace configured. I have used **docker** community edition and **minikube** for this project.
```sh
$ git clone https://github.com/lambdafunc/vmw.git`

from root directory run:
$ go mod tidy

$ go test .
ok  	go_mon/go_mon	1.182s

$ cd go_mon && go run main.go collector.go
```
Service should start and listening on port 8080

![go-mon](https://github.com/lambdafunc/vmw/blob/main/images/go_run_main.png?raw=true)


### Build container image for golang app


`Dockerfile` has been provided which can be used to generate container image

Run from the root directory

`$ docker build -t go-mon .`

![docker build](https://github.com/lambdafunc/vmw/blob/main/images/docker_build_gomon.png?raw=true)

### Run Prometheus and Grafana

```sh
Promotheus:
$ cp prometheus.yml /tmp

$ docker run -d -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

Prometheus will be up and listening on port 9090.

Grafana:
$ docker run -d -p 3000:3000 grafana/grafana grafana

Grafana will be up and running on localhost:3000
```
![Grafana Build](https://github.com/lambdafunc/vmw/blob/main/images/docker_grafana.png?raw=true)

### Run Golang app

```sh
$ docker run -d -p 8080:8080 go_mon`

Below is sample run for /metrics endpoint to verify if metrics are being populated.

$ curl localhost:8080/metrics

# TYPE sample_external_url_response_ms counter
sample_external_url_response_ms{url="https://httpstat.us/200"} 435
sample_external_url_response_ms{url="https://httpstat.us/503"} 435
# HELP sample_external_url_up shows if url is up
# TYPE sample_external_url_up counter
sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

## URLs for accessing deployed services
```sh
Golang Application: http://localhost:8080/metrics
Prometheus: http://localhost:9090/graph
Grafana: http://localhost:3000/
```

## Generate some test data

`$ chmod +x scripts/generateResponseData.sh`

`./scripts/generateResponseData` to generate some metrics

![testdata](https://github.com/lambdafunc/vmw/blob/main/images/metrics_browser.png?raw=true)


## Deployment on Kubernetes:

`$ cd kubernetes-deployment`

Deployment specs for **Grafana, Prometheus** and **Go-Mon Service** exist under kubernetes-deployments directory. Each deployment spec consist of all necessary resources for k8s deployment such as `deployment`, `containers`, `namespace` and `service` definitions. Container image can be pushed to any container registry(dockerhub, ECR etc) and service spec needs to be updated accordingly. (line-26 `kubernetes-deployment/service-monitoring-deployment.yaml`).

```sh
# Deployment-spec for prometheus
$ kubectl apply -f prometheus-deployment.yaml`
# Deployment-spec for prometheus for grafana
$ kubectl apply -f grafana-deployment.yaml`
# Golang app deployment
$ kubectl apply -f service-monitoring-deployment.yaml
```

## Screenshot:
![Grafana-1](https://github.com/lambdafunc/vmw/blob/main/images/grafana_screenshot.png?raw=true)


![Prometheus-1](https://github.com/lambdafunc/vmw/blob/main/images/prom-1.png?raw=true)


![Prometheus-2](https://github.com/lambdafunc/vmw/blob/main/images/prom-2.png?raw=true)
