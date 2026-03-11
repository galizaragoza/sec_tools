package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/jpillora/opts"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	IP     string `opts:"short=i, help=IP of the target (required for Nmap)"`
	URL    string `opts:"short=u, help=Target URL"`
	Domain string `opts:"short=d, help=Target domain (required for dnsrecon)"`
}

func checkPaths(tools []string) (map[string]string, error) {
	paths := make(map[string]string)
	var missing []string

	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			missing = append(missing, tool)
			continue
		}
		paths[tool] = path
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("critical dependencies missing from system: %v", missing)
	}
	return paths, nil
}

func goNmap(c Config, path map[string]string) error {
	if c.IP == "" {
		return fmt.Errorf("[X] Error: No IP was specified, though it is required for Nmap, value of IP: %s", c.IP)
	}

	log.Println("Nmap scan has started")

	nmapPath := path["nmap"]

	reportName := "autorecon_nmapScan.xml"

	args := []string{
		"--top-ports", "500", "--open", "-sS", "-T2", "-f", c.IP, "-oX", reportName,
	}

	cmd := exec.Command(nmapPath, args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[X] Error: Error running the nmap scan:%w", err)
	}

	log.Println("Nmap scan if finished, parsing report")
	log.Println("Nmap scan is complete")

	return nil
}

func goDNSrecon(c Config, path map[string]string) error {
	if c.Domain == "" {
		return fmt.Errorf("[X] Error: No domain was specified, though it is required for DNSrecon, value of Domain: %s", c.Domain)
	}

	log.Println("DNSrecon scan has started")

	dnsReconPath := path["dnsrecon"]

	reportName := "autorecon_DNSrecon.xml"

	args := []string{
		"-d", c.Domain, "-absykwz", "-x", reportName, "-t", "std",
	}

	cmd := exec.Command(dnsReconPath, args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[X] Error: Error running DNS recon: %w", err)
	}

	log.Println("DNSrecon scan is complete")
	return nil
}

func goWhatWeb(c Config, path map[string]string) error {
	if c.URL == "" {
		return fmt.Errorf("[X] Error: URL was not specified, thogh it is required for WhatWeb, value of URL: %s", c.URL)
	}

	log.Println("WhatWeb scan has started")

	whatWebPath := path["whatweb"]
	reportName := "autorecon_WhatWeb.xml"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	args := []string{ // Add UserAgent value would be nice
		"-a", "3", "-c", "--log-xml", reportName, "--user-agent", userAgent,
	}

	cmd := exec.Command(whatWebPath, args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[X] Error: Something happened while scanning with WhatWeb: %w", err)
	}

	log.Println("WhatWeb scan is complete")

	return nil
}

func goWafW00f(c Config, path map[string]string) error {
	if c.URL == "" {
		return fmt.Errorf("[X] Error: URL was not specified, though it is required to use WafW00f, value of URL is: %s", c.URL)
	}

	wafW00fPath := path["wafw00f"]
	reportName := "autorecon_wafw00f.csv"
	log.Println("WafW00f scan has started")

	args := []string{"-a", "-o", reportName, "-f", "csv", "-l", c.URL}

	cmd := exec.Command(wafW00fPath, args...)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	c := Config{}

	opts.Parse(&c)

	tools := []string{"nmap", "wafw00f", "dnsrecon", "whatweb", "xsltproc"}
	paths, err := checkPaths(tools)
	if err != nil {
		log.Fatalf("[X] Dependency error: %v", err)
	}

	log.Printf("All tools found in the machine(%v), proceeding with all: %v", paths, tools)

	g := new(errgroup.Group)

	g.Go(func() error {
		return goNmap(c, paths)
	})

	g.Go(func() error {
		return goDNSrecon(c, paths)
	})

	g.Go(func() error {
		return goWhatWeb(c, paths)
	})

	g.Go(func() error {
		return goWafW00f(c, paths)
	})

	if err := g.Wait(); err != nil {
		return
	}
}
