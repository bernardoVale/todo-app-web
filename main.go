package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

var slaveConnection string
var slavePassword string
var masterConnection string
var masterPassword string

var APP_VERSION string = "latest" // Change this by setting compile flag -ldflags "-X main.APP_VERSION=${CI_GIT_TAG}"

func init() {
	// Metrics have to be registered to be exposed:
	registerMetrics()
}

func main() {
	flag.StringVar(&masterConnection, "master", "redis-master:6379", "The connection string to the Redis master as <hostname/ip>:<port>")
	flag.StringVar(&slaveConnection, "slave", "redis-slave:6379", "The connection string to the Redis slave as <hostname/ip>:<port>")
	flag.StringVar(&masterPassword, "master-password", "", "The password used to connect to the master")
	flag.StringVar(&slavePassword, "slave-password", "", "The password used to connect to the slave")
	flag.Parse()

	// Iniitialize metrics
	getHealthStatus()

	r := mux.NewRouter()
	r.Path("/read/{key}").Methods("GET").HandlerFunc(readTodoHandler)
	r.Path("/insert/{key}/{value}").Methods("GET").HandlerFunc(insertTodoHandler)
	r.Path("/delete/{key}/{value}").Methods("GET").HandlerFunc(deleteTodoHandler)
	r.Path("/health").Methods("GET").HandlerFunc(healthCheckHandler)
	r.Path("/metrics").Methods("GET").Handler(prometheus.Handler())
	r.Path("/whoami").Methods("GET").HandlerFunc(whoAmIHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	http.ListenAndServe(":3000", n)
}
