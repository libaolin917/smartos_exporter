package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/siebenmann/go-kstat"
	"net/http"
	"time"
)

func getNamedUint64Val(ks *kstat.KStat, name string) uint64 {
	n, err := ks.GetNamed(name)
	if err != nil {
		log.Fatalf("getting '%s' from %s: %s", name, ks, err)
	}
	if n.Type != kstat.Uint64 {
		log.Fatalf("Named value is not od Uint64 type: '%s', %v", name, ks)
	}
	return n.UintVal
}

func collectARCstats() {
	log.Debugf("Start collecting ARC stats")
	token, err := kstat.Open()
	if err != nil {
		log.Fatalf("Open failure: %s", err)
	}
	for {
		log.Debugf("Collecting...")
		ks, err := token.Lookup("zfs", 0, "arcstats")
		if err != nil {
			log.Fatalf("lookup failure on %s:0:%s: %s", "zfs", "arcstats", err)
		}
		log.Debugf("Collected: %v", ks)
		log.Debugf("Hits: %d", getNamedUint64Val(ks, "hits"))
		log.Debugf("Misses: %d", getNamedUint64Val(ks, "misses"))
		time.Sleep(10 * time.Second)
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debugf("Starting")
	go collectARCstats()
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe("0.0.0.0:9102", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
