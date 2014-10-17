package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	SleepInterval time.Duration
	Workers       int
}

func init() {
	hostname, _ := os.Hostname()
	log.SetPrefix(hostname + " ")
}

func getConfig() config {
	var workers = flag.Int("workers", 1, "Number of workers")
	var sleepIntervalInSeconds = flag.Int("sleep", 1, "Sleep interval in seconds")
	flag.Parse()

	return config{
		Workers:       *workers,
		SleepInterval: time.Duration(int64(*sleepIntervalInSeconds)) * time.Second,
	}
}

func main() {
	defer log.Println("DONE")

	log.Println("STARTING")

	config := getConfig()

	quitChannel := make(chan struct{})
	var waitGroup sync.WaitGroup

	waitGroup.Add(config.Workers)
	for worker := 1; worker <= config.Workers; worker++ {
		workerIdentifier := "w" + strconv.Itoa(worker)
		go runWorker(workerIdentifier, config.SleepInterval, quitChannel, &waitGroup)
	}

	var input string
	fmt.Scanln(&input)

	close(quitChannel)
	waitGroup.Wait()
}

func runWorker(identifier string, sleepInterval time.Duration, quitChannel chan struct{}, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	defer log.Printf("%s Exiting\n", identifier)

	log.Printf("%s Starting\n", identifier)

	tickerChannel := time.NewTicker(sleepInterval)

	for {
		select {
		case <-tickerChannel.C:
			log.Printf("%s Working\n", identifier)
		case <-quitChannel:
			log.Printf("%s Quiting\n", identifier)
			tickerChannel.Stop()
			return
		}
	}
}
