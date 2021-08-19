package etcdReplicator

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/pickme-go/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	Endpoints   []string
	Username    string
	Password    string
	DialTimeout time.Duration
}
type replicator struct {
	FromNodeClient *clientv3.Client
}

func NewReplicator(fromNode Client) (*replicator, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   fromNode.Endpoints,
		DialTimeout: fromNode.DialTimeout * time.Second,
		Username:    fromNode.Username,
		Password:    fromNode.Password,
	})

	if err != nil {
		return nil, err
	}

	return &replicator{
		FromNodeClient: client,
	}, nil
}

func (r *replicator) Replicate(toNode Client, keyPath string) error {
	defer func(FromNodeClient *clientv3.Client) {
		err := FromNodeClient.Close()
		if err != nil {
			log.Error(err)
		}
	}(r.FromNodeClient)

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   toNode.Endpoints,
		DialTimeout: toNode.DialTimeout * time.Second,
		Username:    toNode.Username,
		Password:    toNode.Password,
	})

	if err != nil {
		return err
	}
	defer func(client *clientv3.Client) {
		err := client.Close()
		if err != nil {
			log.Error(err)
		}
	}(client)

	ctx := context.Background()
	resp, err := r.FromNodeClient.Get(ctx, keyPath, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))

	if err != nil {
		return err
	}

	for i := range resp.Kvs {
		_, err = client.Put(ctx, string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// Save function saves all keys and values in text file
func (r *replicator) Save(filename string, keysPath string) error {
	defer func(FromNodeClient *clientv3.Client) {
		err := FromNodeClient.Close()
		if err != nil {
			log.Error(err)
		}
	}(r.FromNodeClient)

	ctx := context.Background()

	resp, err := r.FromNodeClient.Get(ctx, keysPath, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return err
	}
	f, err := os.Create(filename + ".csv")
	if err != nil {
		log.Error(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for i := range resp.Kvs {
		r := make([]string, 0, 2)
		r = append(r, string(resp.Kvs[i].Key))
		r = append(r, string(resp.Kvs[i].Value))
		err := w.Write(r)
		if err != nil {
			log.Error(err)
		}
	}
	return nil

}

// LoadAndReplicate function loads a txt and puts the keys
func (r *replicator) LoadAndReplicate(filename string, toNode Client) error {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   toNode.Endpoints,
		DialTimeout: toNode.DialTimeout * time.Second,
		Username:    toNode.Username,
		Password:    toNode.Password,
	})

	if err != nil {
		return err
	}

	defer func(client *clientv3.Client) {
		err := client.Close()
		if err != nil {
			log.Error(err)
		}
	}(client)

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	records, err := csvr.ReadAll()
	if err != nil {
		return err
	}

	ctx := context.Background()
	for i := range records {
		fmt.Println(records[i][0] + records[i][1])
		_, err = client.Put(ctx, records[i][0], records[i][1])
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}
