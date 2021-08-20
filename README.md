# ETCD Replicator
A ETCD key-value replicator.

## Features

- Replicate from one node to another.
- Save key value pairs from a node to a CSV file.
- Load and add key abd values to ETCD node.

## Install
````
go get github.com/anjulapaulus/etcd-replicator
````

## Usage

Replicate one node to another

````
package main

import (
	"github.com/anjulapaulus/etcd-replicator"
	"github.com/pickme-go/log"
)

func main(){
	client, err := etcdReplicator.NewReplicator(etcdReplicator.Client{
		Endpoints:   []string{"127.0.0.1:2379"},
		Username:    "root",
		Password:    "root",
		DialTimeout: 2,
	})
	if err != nil{
		log.Error(err)
	}

	err = client.Replicate(etcdReplicator.Client{
		Endpoints:   []string{"127.0.0.1:2379"},
		Username:"",
		Password:"",
		DialTimeout:2,
	}, "/system")
	
	if err != nil {
		return
	}

}

````
Save key value pairs to a CSV file.

````
import (
	"github.com/anjulapaulus/etcd-replicator"
	"github.com/pickme-go/log"
)

func main(){
	client, err := etcdReplicator.NewReplicator(etcdReplicator.Client{
		Endpoints:   []string{"127.0.0.1:2379"},
		Username:    "root",
		Password:    "root",
		DialTimeout: 2,
	})
	if err != nil{
		log.Error(err)
	}
    err = client.Save("data", "/path")
	if err != nil {
		return
	}
}

````

Load and add key abd values to ETCD node.

````
func main(){
    err = client.LoadAndReplicate("data.csv", etcdReplicator.Client{
			Endpoints:   []string{"127.0.0.21:2379"},
			Username:"",
			Password:"",
			DialTimeout:2,
	})
	if err != nil{
		log.Error(err)
	}
}
````
