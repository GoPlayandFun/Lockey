package node

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SystemBuilders/LocKey/internal/lockservice"
	"github.com/SystemBuilders/LocKey/internal/lockservice/routing"
	"github.com/gorilla/mux"
)

// Start begins the node's operation as a http server.
func Start(ls *lockservice.SimpleLockService, scfg lockservice.SimpleConfig) error {

	IP := scfg.IP()
	IP = strings.TrimPrefix(IP, "http://")
	port := scfg.Port()

	if err := checkValidPort(port); err != nil {
		return err
	}

	router := mux.NewRouter()

	router = routing.SetupRouting(ls, router)

	server := &http.Server{
		Handler: router,
		Addr:    IP + ":" + port,
	}

	go gracefulShutdown(server)

	log.Println("Starting Server on " + IP + ":" + port)
	return server.ListenAndServe()
}

// gracefulShutdown shuts down the server on getting a ^C signal
func gracefulShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for currently serving items.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func checkValidPort(port string) error {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if portInt > 65535 {
		return errors.New("Port number exceeds limit of 65535")
	}
	return nil
}
