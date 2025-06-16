package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

var (
	timeout   int
	threads   int
	verbose   bool
	portRange string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rully-scan",
		Short: "Rully Scan - Multi-platform Network Tool",
		Long: `Rully Scan is a complete network tool for:
- IP and DNS verification
- Local and remote port scanning
- Network connectivity analysis`,
		Run: func(cmd *cobra.Command, args []string) {
			printBanner()
			cmd.Help()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 3, "Timeout in seconds")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "T", 100, "Number of threads")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")

	// Command for IP/DNS verification
	var ipCmd = &cobra.Command{
		Use:   "ip [host]",
		Short: "Check IP and DNS information",
		Args:  cobra.ExactArgs(1),
		Run:   runIPCheck,
	}

	// Command for port scanning
	var scanCmd = &cobra.Command{
		Use:   "scan [host]",
		Short: "Scan host ports",
		Args:  cobra.ExactArgs(1),
		Run:   runPortScan,
	}

	scanCmd.Flags().StringVarP(&portRange, "ports", "p", "1-1000", "Port range (ex: 1-1000, 80,443,8080)")

	// Command for open local ports
	var localCmd = &cobra.Command{
		Use:   "local",
		Short: "List locally open ports",
		Run:   runLocalScan,
	}

	rootCmd.AddCommand(ipCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(localCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printBanner() {
	banner := `
 ____        _ _        ____                  
|  _ \ _   _| | |_   _ / ___|  ___ __ _ _ __  
| |_) | | | | | | | | \___ \ / __/ _` + "`" + ` | '_ \ 
|  _ <| |_| | | | |_| |___) | (_| (_| | | | |
|_| \_\\__,_|_|_|\__, |____/ \___\__,_|_| |_|
                |___/                       
` + ColorCyan + "Network Scanner Tool v1.0" + ColorReset + `
` + ColorYellow + "Developed in Go Lang - Multi-platform" + ColorReset + `
`
	fmt.Println(ColorPurple + banner + ColorReset)
}

func runIPCheck(cmd *cobra.Command, args []string) {
	host := args[0]

	fmt.Printf("\n%s[INFO]%s Checking information for: %s%s%s\n",
		ColorBlue, ColorReset, ColorBold, host, ColorReset)

	// DNS resolution
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("%s[ERROR]%s DNS resolution failed: %v\n", ColorRed, ColorReset, err)
		return
	}

	fmt.Printf("\n%s[DNS]%s IP addresses found:\n", ColorGreen, ColorReset)
	for _, ip := range ips {
		ipType := "IPv4"
		if ip.To4() == nil {
			ipType = "IPv6"
		}
		fmt.Printf("  %s%s%s (%s)\n", ColorCyan, ip.String(), ColorReset, ipType)
	}

	// Reverse DNS for the first IP
	if len(ips) > 0 {
		names, err := net.LookupAddr(ips[0].String())
		if err == nil && len(names) > 0 {
			fmt.Printf("\n%s[REVERSE DNS]%s\n", ColorGreen, ColorReset)
			for _, name := range names {
				fmt.Printf("  %s%s%s\n", ColorCyan, name, ColorReset)
			}
		}
	}

	// CNAME information
	cname, err := net.LookupCNAME(host)
	if err == nil && cname != host+"." {
		fmt.Printf("\n%s[CNAME]%s %s%s%s\n",
			ColorGreen, ColorReset, ColorCyan, cname, ColorReset)
	}

	// MX records
	mxRecords, err := net.LookupMX(host)
	if err == nil && len(mxRecords) > 0 {
		fmt.Printf("\n%s[MX RECORDS]%s\n", ColorGreen, ColorReset)
		for _, mx := range mxRecords {
			fmt.Printf("  %s%s%s (priority: %d)\n",
				ColorCyan, mx.Host, ColorReset, mx.Pref)
		}
	}
}

func runPortScan(cmd *cobra.Command, args []string) {
	host := args[0]

	fmt.Printf("\n%s[INFO]%s Scanning ports for: %s%s%s\n",
		ColorBlue, ColorReset, ColorBold, host, ColorReset)

	ports := parsePortRange(portRange)
	if len(ports) == 0 {
		fmt.Printf("%s[ERROR]%s Invalid port range\n", ColorRed, ColorReset)
		return
	}

	fmt.Printf("%s[INFO]%s Scanning %d ports with %d threads (timeout: %ds)\n",
		ColorBlue, ColorReset, len(ports), threads, timeout)

	start := time.Now()
	openPorts := scanPorts(host, ports)
	duration := time.Since(start)

	if len(openPorts) > 0 {
		fmt.Printf("\n%s[OPEN PORTS]%s Found %d ports:\n",
			ColorGreen, ColorReset, len(openPorts))

		for _, port := range openPorts {
			service := getServiceName(port)
			fmt.Printf("  %s%d/tcp%s\t%s%s%s\n",
				ColorGreen, port, ColorReset, ColorYellow, service, ColorReset)
		}
	} else {
		fmt.Printf("\n%s[INFO]%s No open ports found\n",
			ColorYellow, ColorReset)
	}

	fmt.Printf("\n%s[COMPLETED]%s Scan finished in %v\n",
		ColorBlue, ColorReset, duration.Round(time.Millisecond))
}

func runLocalScan(cmd *cobra.Command, args []string) {
	fmt.Printf("\n%s[INFO]%s Checking locally open ports...\n",
		ColorBlue, ColorReset)

	// Check common TCP ports locally
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 993, 995, 1433, 3306, 3389, 5432, 5900, 8080, 8443}

	openPorts := scanPorts("localhost", commonPorts)

	if len(openPorts) > 0 {
		fmt.Printf("\n%s[OPEN LOCAL PORTS]%s\n", ColorGreen, ColorReset)
		for _, port := range openPorts {
			service := getServiceName(port)
			fmt.Printf("  %s%d/tcp%s\t%s%s%s\n",
				ColorGreen, port, ColorReset, ColorYellow, service, ColorReset)
		}
	} else {
		fmt.Printf("\n%s[INFO]%s No common ports open locally\n",
			ColorYellow, ColorReset)
	}
}

func parsePortRange(portRange string) []int {
	var ports []int

	if strings.Contains(portRange, ",") {
		// Port list: 80,443,8080
		portStrs := strings.Split(portRange, ",")
		for _, portStr := range portStrs {
			port, err := strconv.Atoi(strings.TrimSpace(portStr))
			if err == nil && port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
	} else if strings.Contains(portRange, "-") {
		// Port range: 1-1000
		parts := strings.Split(portRange, "-")
		if len(parts) == 2 {
			start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

			if err1 == nil && err2 == nil && start > 0 && end <= 65535 && start <= end {
				for i := start; i <= end; i++ {
					ports = append(ports, i)
				}
			}
		}
	} else {
		// Single port
		port, err := strconv.Atoi(strings.TrimSpace(portRange))
		if err == nil && port > 0 && port <= 65535 {
			ports = append(ports, port)
		}
	}

	return ports
}

func scanPorts(host string, ports []int) []int {
	var openPorts []int
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Channel to limit threads
	semaphore := make(chan struct{}, threads)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, p),
				time.Duration(timeout)*time.Second)

			if err == nil {
				conn.Close()
				mutex.Lock()
				openPorts = append(openPorts, p)
				mutex.Unlock()

				if verbose {
					fmt.Printf("%s[OPEN]%s %s:%d\n",
						ColorGreen, ColorReset, host, p)
				}
			} else if verbose {
				fmt.Printf("%s[CLOSED]%s %s:%d\n",
					ColorRed, ColorReset, host, p)
			}
		}(port)
	}

	wg.Wait()
	return openPorts
}

func getServiceName(port int) string {
	services := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		135:  "RPC",
		139:  "NetBIOS",
		143:  "IMAP",
		443:  "HTTPS",
		445:  "SMB",
		993:  "IMAPS",
		995:  "POP3S",
		1433: "MSSQL",
		3306: "MySQL",
		3389: "RDP",
		5432: "PostgreSQL",
		5900: "VNC",
		8080: "HTTP-Alt",
		8443: "HTTPS-Alt",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}
