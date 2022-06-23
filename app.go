package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"time"
)

//go:embed servers.json
var serversBytes []byte
var servers []Server

func printMainMenu() {
	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("         \"No Internet, Secured\" Fixer")
	fmt.Println("==============================================")
	fmt.Println("1. Set NCSI registries automatically")
	fmt.Println("2. Set NCSI registries manually")
	fmt.Println("3. Set NCSI registries manually (show latency)")
	fmt.Println("4. Print system NCSI registry values")
	fmt.Println("5. Test system NCSI server")
	fmt.Println("6. Exit")
	fmt.Println("==============================================")
	fmt.Println()
}

func getLatencyString(duration time.Duration) string {
	if duration == -1 {
		return "Failed"
	} else {
		return fmt.Sprintf("%s", duration)
	}
}

func printServerLatencies() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Server Name", "Web Probe Latency", "Web ProbeV6 Latency", "DNS Probe Latency", "DNS ProbeV6 Latency", "Average Latency"})
	for _, server := range servers {
		table.Append([]string{
			server.Name,
			getLatencyString(server.WebProbeLatency),
			getLatencyString(server.WebProbeV6Latency),
			getLatencyString(server.DnsProbeLatency),
			getLatencyString(server.DnsProbeV6Latency),
			getLatencyString(server.AverageLatency),
		})
	}
	table.Render()
}

func printSystemNCSIRegistry() {
	ncsiReg, err := GetSystemNCSIReg()
	if err != nil {
		log.Fatalln("Failed to get system NCSI registry.")
	}
	ncsiReg.Print()
}

func testServersAndSort() {
	for i := 0; i < len(servers); i++ {
		servers[i].Test()
		fmt.Println()
	}
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].AverageLatency == -1 {
			return false
		} else if servers[j].AverageLatency == -1 {
			return true
		} else {
			return servers[i].AverageLatency < servers[j].AverageLatency
		}
	})
}

func setNCSIRegistriesAutomatically() {
	testServersAndSort()
	printServerLatencies()
	if servers[0].AverageLatency == -1 {
		fmt.Println("No server is available.")
		os.Exit(1)
	}
	server := servers[0]
	log.Println(fmt.Sprintf("Setting system NCSI registry to %s.", server.Name))
	err := server.ToNCSIReg().setSystemNCSIReg()
	if err != nil {
		log.Fatalln("Failed to set system NCSI registry.")
	}
	log.Println(fmt.Sprintf("Successfully set system NCSI registry to %s.", server.Name))
}

func setNCSIRegistriesManuallyWithLatency() {
	testServersAndSort()
	printServerLatencies()
	fmt.Println()
	manualChooseServer()
}

func testSystemNCSIServer() {
	ncsiReg, err := GetSystemNCSIReg()
	if err != nil {
		log.Fatalln("Failed to get system NCSI registry.")
	}
	ncsiReg.Print()
	fmt.Println()
	server := ncsiReg.ToServer()
	server.Name = "System"
	server.Test()
}

func pause() {
	fmt.Println()
	fmt.Print("🔙 Press Enter to continue...")
	fmt.Scanln()
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func mainMenu() {
	printMainMenu()
	choice := getChoice()
	switch choice {
	case 1:
		// Set NCSI registries automatically
		setNCSIRegistriesAutomatically()
		pause()
		break
	case 2:
		// Set NCSI registries manually
		manualChooseServer()
		pause()
		break
	case 3:
		// Set NCSI registries manually (show latency)
		setNCSIRegistriesManuallyWithLatency()
		pause()
		break
	case 4:
		// Print system NCSI registry values
		printSystemNCSIRegistry()
		pause()
		break
	case 5:
		// Test system NCSI server
		testSystemNCSIServer()
		pause()
		break
	case 6:
		// Exit
		fmt.Println("💙 Bye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice.")
		break
	}
}

func getChoice() int {
	var choice int
	fmt.Print("✨ Enter your choice: ")
	fmt.Scanf("%d\n", &choice)
	fmt.Println()
	return choice
}

func printManualChooseServerMenu() {
	fmt.Println("------------------------------------")
	fmt.Println("         Choose a server")
	fmt.Println("------------------------------------")
	for index, server := range servers {
		fmt.Println(fmt.Sprintf("%d. %s", index+1, server.Name))
	}
}

func manualChooseServer() {
	printManualChooseServerMenu()
	choice := getChoice()
	if choice < 1 || choice > len(servers) {
		fmt.Println("Invalid choice.")
		os.Exit(1)
	}
	server := servers[choice-1]
	log.Println(fmt.Sprintf("Setting system NCSI registry to %s.", server.Name))
	err := server.ToNCSIReg().setSystemNCSIReg()
	if err != nil {
		log.Fatalln("Failed to set system NCSI registry.")
	}
	fmt.Println(fmt.Sprintf("Successfully set system NCSI registry to %s.", server.Name))
}

func main() {
	if runtime.GOOS != "windows" {
		log.Fatalln("This program only works on Windows.")
	}

	err := json.Unmarshal(serversBytes, &servers)
	if err != nil {
		log.Fatalln("Failed to load servers.")
	}
	if len(servers) == 0 {
		log.Fatalln("No servers found.")
	}

	for {
		clearScreen()
		mainMenu()
	}
}
