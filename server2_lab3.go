package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type KeyValue struct {
	Key   int    `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

var s []KeyValue
var index int

type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func GetAllKeys(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {

	sort.Sort(ByKey(s))
	result, _ := json.Marshal(s)
	fmt.Fprintln(rw, string(result))
}

func PutKeys(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
	key, _ := strconv.Atoi(p.ByName("key_id"))
	s = append(s, KeyValue{key, p.ByName("value")})
	index++
}

func GetKey(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {

	ind := index
	key, _ := strconv.Atoi(p.ByName("key_id"))
	for i := 0; i < ind; i++ {
		if s[i].Key == key {
			result, _ := json.Marshal(s[i])
			fmt.Fprintln(rw, string(result))
		}
	}
}

func main() {
	index = 0
	mux := httprouter.New()
	mux.GET("/keys", GetAllKeys)
	mux.GET("/keys/:key_id", GetKey)
	mux.PUT("/keys/:key_id/:value", PutKeys)
	go http.ListenAndServe(":3001", mux)
	select {}
}
