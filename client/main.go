// client

package main

import (
	"log"
	"net/rpc"
)

type KVPair struct {
	Key, Value string
}

func getClient() *rpc.Client {

	client, err := rpc.DialHTTP("tcp", "localhost:5001")
	if err != nil {
		log.Fatal("error while connecting to localhost:5001: ", err)
	}
	return client
}

func getValue(key string) string {
	var replyValue string
	log.Println("Sending Get with: ", key)

	client := getClient()

	err := client.Call("RPCObj.Get", key, &replyValue)
	if err != nil {
		log.Fatal("error while calling RPCObj.Get: ", err)
	}

	log.Println("Received: ", replyValue)
	return replyValue
}

func setValue(kvPair KVPair) KVPair {
	var replyValue KVPair
	log.Println("Sending Set with: ", kvPair)

	client := getClient()

	err := client.Call("RPCObj.Set", kvPair, &replyValue)
	if err != nil {
		log.Fatal("error while calling RPCObj.Get: ", err)
	}

	log.Println("Received: ", replyValue)
	return replyValue
}

func updateValue(kvPair KVPair) KVPair {
	var replyValue KVPair
	log.Println("Sending Update with: ", kvPair)

	client := getClient()

	err := client.Call("RPCObj.Update", kvPair, &replyValue)
	if err != nil {
		log.Fatal("error while calling RPCObj.Update: ", err)
	}

	log.Println("Received: ", replyValue)
	return replyValue
}

func deleteValue(key string) KVPair {
	var replyValue KVPair
	log.Println("Sending Delete with: ", key)

	client := getClient()

	err := client.Call("RPCObj.Delete", key, &replyValue)
	if err != nil {
		log.Fatal("error while calling RPCObj.Delete: ", err)
	}

	log.Println("Received: ", replyValue)
	return replyValue
}

func main() {
	getValue("ping")
	updateValue(KVPair{Key: "ping", Value: "Hello, World!"})
	getValue("ping")
	setValue(KVPair{Key: "checkin", Value: "All clear"})
	getValue("checkin")
	deleteValue("ping")
	getValue("ping")
}
