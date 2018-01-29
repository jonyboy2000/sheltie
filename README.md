# Sheltie

This program receives messages from **RabbitMQ** and parses them, and then runs the corresponding **Kubernetes** job and read their results from file (if we consider kubernetes job result's are saved to a file)

***

## Used Libraries :

* [Go RabbitMQ Client (amqp)](https://github.com/streadway/amqp)

* [Go Kubernetes Client (client-go)](https://github.com/kubernetes/client-go)

* ~~[Go HTTP and REST client (resty)](https://github.com/go-resty/resty)~~

* [Go YAML support (yaml)](https://github.com/go-yaml/yaml)


## How to build :


* **Install Golang > v1.8.x** 


* **Install "Go RabbitMQ Client" :** 

```
go get github.com/streadway/amqp
```

* **Install "Go Kubernetes Client" :** 


```
go get k8s.io/client-go/...
```

```
go get -u k8s.io/apimachinery/...
```

* **~~Install "Go HTTP and REST client" :~~** 

```
go get -u gopkg.in/resty.v1
```

* **Install "Go YAML support" :** 

```
go get gopkg.in/yaml.v2
```


* **Go to project root directory :** 

```
go build sheltie.go
```
