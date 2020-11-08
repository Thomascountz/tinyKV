// server

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type KVPair struct {
	Key, Value string
}

type store struct {
	data []KVPair
	mu   sync.Mutex
}

var datastore store

type RPCObj int

func (r *RPCObj) Get(key string, replyValue *string) error {
	log.Println("Received Get with: ", key)

	datastore.mu.Lock()
	defer datastore.mu.Unlock()

	index, recordExists := r.get(key)
	if !recordExists {
		return fmt.Errorf("key %q not found", key)
	}

	*replyValue = datastore.data[index].Value
	return nil
}

func (r *RPCObj) Set(queryKVPair KVPair, replyValue *KVPair) error {
	log.Println("Received Set with: ", queryKVPair)

	datastore.mu.Lock()
	defer datastore.mu.Unlock()

	_, recordExists := r.get(queryKVPair.Key)
	if recordExists {
		return fmt.Errorf("key %q already exists.", queryKVPair.Key)
	}

	datastore.data = append(datastore.data, queryKVPair)

	*replyValue = queryKVPair

	return nil
}

func (r *RPCObj) Update(queryKVPair KVPair, replyValue *KVPair) error {
	log.Println("Received Update with: ", queryKVPair)

	datastore.mu.Lock()
	defer datastore.mu.Unlock()

	_, recordExists := r.get(queryKVPair.Key)
	if !recordExists {
		return fmt.Errorf("key %q does not exist.", queryKVPair.Key)
	}

	for index, kvPair := range datastore.data {
		if kvPair.Key == queryKVPair.Key {
			datastore.data[index] = queryKVPair
		}
	}

	*replyValue = queryKVPair

	return nil
}

func (r *RPCObj) Delete(key string, replyValue *KVPair) error {
	log.Println("Received Delete with: ", key)

	datastore.mu.Lock()
	defer datastore.mu.Unlock()

	index, recordExists := r.get(key)
	if !recordExists {
		return fmt.Errorf("key %q does not exist.", key)
	}

	kvPair := datastore.data[index]
	datastore.data[index] = datastore.data[len(datastore.data)-1]
	datastore.data[len(datastore.data)-1] = KVPair{}
	datastore.data = datastore.data[:len(datastore.data)-1]

	*replyValue = kvPair

	return nil
}

func (r *RPCObj) get(key string) (int, bool) {
	for index, kvPair := range datastore.data {
		if kvPair.Key == key {
			return index, true
		}
	}
	return 0, false
}

func buildDefaultStore() {
	datastore.data = append(datastore.data, KVPair{Key: "ping", Value: "pong"})
}

func main() {
	buildDefaultStore()

	rpcObj := new(RPCObj)

	err := rpc.Register(rpcObj)
	if err != nil {
		log.Fatal("error while registering RPCObj", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":5001")

	if err != nil {
		log.Fatal("error while opening a connection on the network", err)
	}

	log.Println("serving RPCObj on port 5001")

	http.Serve(listener, nil)
	if err != nil {
		log.Fatal("error while serving a connection on the network", err)
	}
}
