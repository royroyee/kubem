# kubem
A Simple Kubernetes Monitoring System

## Description
With Kubem, you can effortlessly monitor the resources and performance of your Kubernetes cluster.


## Architecture
<img width="600" alt="kubem-아키텍처" src="https://github.com/royroyee/kubem/assets/88774925/17d3bac3-476d-43e2-a0e2-54a65b51684c">

## Features
- Kubernetes Monitoring
	- Controllers
	- Resources
	- Metrics
	- Events

- Kubernetes API client for retrieving data from the cluster (using [client-go](https://github.com/kubernetes/client-go))
- DB handler for storing K8s information and data
- Providing REST API for frontend

## Installation
- Dockerfile 

### MongoDB
Kubem uses MongoDB in order to store and retrieve data. Therefore there must be an MongoDB instance (a containered one or just the native one) that shall be running for Kubem




## Contributors
- [Younghwan Kim](https://github.com/royroyee)

## License
[MIT License](https://github.com/royroyee/kubem/blob/main/LICENSE)
