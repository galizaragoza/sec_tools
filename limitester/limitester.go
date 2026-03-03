package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/jpillora/opts"
)

type Config struct {
	URL      string `opts:"short=u, help=Set the URL to test"`
	Method   string `opts:"short=m, help=Use GET, POST, or both (e.g. -m both, -m get)"`
	Packets  int    `opts:"short=p, help=Amount of packets to be sent (defaults to 1000)"`
	Parallel int    `opts:"short=w, help=Amount of parellel workers"`
	Threads  int    `opts:"short=t, help=Amount of system threads"`
}

func checkOpts(c Config) (Config, error) {
	if c.URL == "" {
		return c, fmt.Errorf("limitester needs a URL to send the packets to, URL is set to: %#v", c.URL)
	}
	if c.Packets <= 0 {
		c.Packets = 1000
	}
	if c.Parallel <= 10 {
		c.Parallel = 10
	}
	return c, nil
}

func getLimits(c Config) error {
	var worker sync.WaitGroup
	jobs := make(chan int, c.Packets)

	for w := 1; w <= c.Parallel; w++ {
		worker.Add(1)
		go func(workerID int) {
			defer worker.Done()
			for j := range jobs {
				req, err := http.Get(c.URL)
				if err != nil {
					fmt.Printf("Worker %d: Error testing limits: %d: %v\n", workerID, j, err)
					continue
				}
				fmt.Printf("Worker %d | Request %d of %d: %s\n", workerID, j, c.Packets, req.Status)
				req.Body.Close()
			}
		}(w)
	}

	count := 0
	for i := c.Packets; i > 0; i-- {
		req, err := http.Get(c.URL)
		if err != nil {
			return fmt.Errorf("Error testing limits:", err)
		}
		count += 1
		fmt.Printf("Request %d of %d: %#v\n", count, c.Packets, req.Status)
	}
	return nil
}

func postLimits(c Config) error {
	var worker sync.WaitGroup
	jobs := make(chan int, c.Packets)

	for w := 1; w <= c.Parallel; w++ {
		worker.Add(1)
		go func(workerID int) {
			defer worker.Done()
			for j := range jobs {
				req, err := http.Post("POST", c.URL, &io.LimitedReader{})
				if err != nil {
					fmt.Printf("Worker %d: Error testing limits: %d: %v\n", workerID, j, err)
					continue
				}
				fmt.Printf("Worker %d | Request %d of %d: %s\n", workerID, j, c.Packets, req.Status)
				req.Body.Close()
			}
		}(w)
	}

	count := 0
	for i := c.Packets; i > 0; i-- {
		req, err := http.Get(c.URL)
		if err != nil {
			return fmt.Errorf("Error testing limits:", err)
		}
		count += 1
		fmt.Printf("Request %d of %d: %#v\n", count, c.Packets, req.Status)
	}
	return nil
}

func main() {
	start := time.Now()
	fmt.Printf("Starting at %#v", start)
	c := Config{}

	opts.Parse(&c)

	runtime.GOMAXPROCS(c.Threads)

	c, err := checkOpts(c)
	if err != nil {
		fmt.Println("Error validating config:", err)
	}

	switch c.Method {
	case "get":
		getLimits(c)
	case "post":
		postLimits(c)
	case "both":
		getLimits(c)
		postLimits(c)
	}

	elapsed := time.Since(start)
	fmt.Printf("Did %d requests in %#v minutes", c.Packets, elapsed)
}
