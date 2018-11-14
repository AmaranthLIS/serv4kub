package handlers

import (
    "github.com/gorilla/mux"
    "sync/atomic"
    "time"
    "github.com/rs/zerolog/log"
)

// Router register necessary routes and returns an instance of a router.
func Router(buildTime, commit, release string) *mux.Router {
    isReady := &atomic.Value{}
    isReady.Store(false)
    go func() {
        log.Printf("Readyz probe is negative by default...")
        time.Sleep(10 * time.Second)//just for test
        isReady.Store(true)
        log.Printf("Readyz probe is positive.")
    }()


    r := mux.NewRouter()
    //r.HandleFunc("/home", home).Methods("GET")
    r.HandleFunc("/", GetClientIP()).Methods("GET")
    r.HandleFunc("/ip", GetServerIp).Methods("GET")
    r.HandleFunc("/home", home(buildTime, commit, release)).Methods("GET")
    r.HandleFunc("/healthz", healthz)
    r.HandleFunc("/readyz", readyz(isReady))
    return r
}
