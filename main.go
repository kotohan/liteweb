package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var serverSettings struct {
	FileDir string `json:"directory"`
	Port    int    `json:"port"`
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "s", "./settings.json", "Configurations File Path")
	flag.Parse()
	err := readSettings(fileName)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	routingRequests(serverSettings.Port)
}

func readSettings(fileName string) (err error) {
	jsonDef := []byte(`{"directory":"./data","port":80}`)
	rawData, err := os.ReadFile(fileName)
	if err != nil {
		json.Unmarshal(jsonDef, &serverSettings)
		return fmt.Errorf("file error")
	}
	err = json.Unmarshal(rawData, &serverSettings)
	if err != nil {
		json.Unmarshal(jsonDef, &serverSettings)
		return fmt.Errorf("format error")
	}
	return nil
}

func routingRequests(port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{request}", requestHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)["request"]
	raw, err := os.ReadFile(serverSettings.FileDir + "/" + request)
	if err != nil {
		res, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	}
}
