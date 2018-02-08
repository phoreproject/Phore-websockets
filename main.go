// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/phoreproject/btcd/rpcclient"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	var connCfg = &rpcclient.ConnConfig{
		Host:                 "rpc.phore.io/rpc",
		HTTPPostMode:         true,
		DisableTLS:           false,
		DisableAutoReconnect: false,
		DisableConnectOnNew:  false,
	}

	client, _ := rpcclient.New(connCfg, nil)

	blockCount, err := client.GetBlockCount()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Block count: %d\n", blockCount)
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
