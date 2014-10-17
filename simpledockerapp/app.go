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

func init() {
	hostname, _ := os.Hostname()
	log.SetPrefix(hostname + " ")
}

func main() {
	defer log.Println("DONE")

	log.Println("STARTING")

	var workers = flag.Int("workers", 1, "Number of workers")
	var sleepIntervalInSeconds = flag.Int("sleep", 1, "Sleep interval in seconds")
	flag.Parse()

	sleepInterval := time.Duration(int64(*sleepIntervalInSeconds)) * time.Second
	quitChannel := make(chan struct{})
	var waitGroup sync.WaitGroup

	waitGroup.Add(*workers)
	for worker := 1; worker <= *workers; worker++ {
		workerIdentifier := "w" + strconv.Itoa(worker)
		go runWorker(workerIdentifier, sleepInterval, quitChannel, &waitGroup)
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
