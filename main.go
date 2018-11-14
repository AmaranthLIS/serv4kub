//package devkub

package main

import (
    "net/http"
    "github.com/rs/zerolog/log"
    "play4j/devkub/handlers"
    "os"
    "play4j/devkub/version"
    "os/signal"
    "syscall"
    "context"
)

func main() {
    log.Printf(
        "Starting the service...\ncommit: %s, build time: %s, release: %s",
        version.Commit, version.BuildTime, version.Release,
    )

    /*http.HandleFunc("/home", func(w http.ResponseWriter, _ *http.Request) {
        fmt.Fprint(w, "Hello! Your request was processed.")
    },)*/
    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal().Msg("Port is not set.")
    }

    router := handlers.Router(version.BuildTime, version.Commit, version.Release)

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: router,
    }
    go func() {
        log.Fatal().Err(srv.ListenAndServe()).Msg("")
    }()

    log.Print("The service is ready to listen and serve.")
    //log.Error().AnErr("httpServer", http.ListenAndServe(":"+port, router)).Msg("terminatedOut")
    killSignal := <-interrupt
    switch killSignal {
    case os.Kill:
        log.Print("Got SIGKILL...")
    case os.Interrupt:
        log.Print("Got SIGINT...")
    case syscall.SIGTERM:
        log.Print("Got SIGTERM...")
    }
    log.Print("The service is shutting down...")
    srv.Shutdown(context.Background())
    log.Print("Done")
}
