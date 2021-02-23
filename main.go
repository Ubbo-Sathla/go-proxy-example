package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.TraceLevel)

}

func main() {
	var listen, backend string
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&listen, "listen", "", "Listen for connections on this address.")
	flags.StringVar(&backend, "backend", "", "The address of the backend to forward to.")
	flags.Parse(os.Args[1:])

	if listen == "" || backend == "" {
		fmt.Fprintln(os.Stderr, "listen and backend options required")
		os.Exit(1)
	}

	p := Proxy{Listen: listen, Backend: backend}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := p.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
