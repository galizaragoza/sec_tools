package main

import (
	"fmt"
	"os/exec"

	"github.com/jpillora/opts"
)

type Config struct {
	// Info
	Stealth int    `opts:"short=s, help=Level of stealth wanted in the recon (1 is less agressive, 5 is most agressive), defaults to 1"`
	IP      string `opts:"short=i, help=IP of the target, required for Nmap"`
	URL     string `opts:"short=u, help=Target URL"`
	Domain  string `opts:"short=d, help=Target domain, required for dnsrecon"`
	// Tool selection
	Full     bool `opts:"short=f, help=Uses Nikto, DNSrecon and Nmap (NSE has to be selected separately)"`
	Nikto    bool `opts:"help=Use nikto in the recon"`
	Dnsrecon bool `opts:"help=Use dnsrecon in the recon"`
	Nmap     bool `opts:"help=Use nmap in the recon"`
	// NSE
	NSE    bool   `opts:"help=Use NSE (nmap scripting engine), the name of a script or category has to be provided"`
	Script string `opts:"help=Specify script or category to use with NSE"`
}

func getPath(isSet bool, tool string) (string, error) {
	if isSet {
		path, err := exec.LookPath(tool)
		if err != nil {
			return "", fmt.Errorf("[X] Error: %s is not installed\n\n", tool)
		}
		return path, nil
	}
	return "", nil
}

func goNmap(c Config) error {
	nmapPath, err := getPath(c.Nmap, "nmap")
	if err != nil {
		return err
	}

	var cmd []string

	if c.IP == "" {
		return fmt.Errorf("[X] Error: Nmap is being called but no IP was provided. Value of IP %#v\n\n", c.IP)
	}
	switch c.Stealth {
	case 1:
		cmd = []string{
			nmapPath, "--top-ports", "100", "-sS", "--stats-every=10", "-T1", c.IP, "-oX", "autorecon_nmapScan.xml",
		}
	case 2:
		cmd = []string{
			nmapPath, "--top-ports", "1000", "-sS", "-sV", "--stats-every=10", "-T2", c.IP, "-oX", "autorecon_nmapScan.xml",
		}
	case 3:
		cmd = []string{
			nmapPath, "--top-ports", "10000", "-sT", "-sV", "-sC", "--stats-every=10", "-T3", c.IP, "-oX", "autorecon_nmapScan.xml",
		}
	case 4:
		cmd = []string{
			nmapPath, "-p", "0-65535", "--open", "-sT", "-A", "--stats-every=10", "-T4", c.IP, "-oX", "autorecon_nmapScan.xml",
		}
	case 5:
		cmd = []string{
			nmapPath, "-p", "0-65535", "--open", "-sT", "-A", "--stats-every=10", "-T5", c.IP, "-oX", "autorecon_nmapScan.xml",
		}
	default:
		return fmt.Errorf("[X] Error: The stealth value (1-5) provided was incorrect. Value of Stealth %#v", c.Stealth)
	}
	fmt.Printf("Scanning with nmap...\n%#v", cmd)
	// exec.Command(130H15|, arg ...string)
	return nil
}

func main() {
	c := Config{
		Stealth: 1,
		// SOLO PARA TESTING DE AQUÍ PARA ABAJO
		IP: "127.0.0.1",
	}

	opts.Parse(&c)

	if c.Full {
		c.Nikto, c.Dnsrecon, c.Nmap = true, true, true
	}

	fmt.Printf("Current configuration: %+v\n\n", c)
}
