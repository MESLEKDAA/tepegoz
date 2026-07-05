package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/MESLEKDAA/tepegoz/pkg/analyzer"
	"github.com/MESLEKDAA/tepegoz/pkg/config"
	"github.com/MESLEKDAA/tepegoz/pkg/model"
	"github.com/fatih/color"
)

func main() {

	filePath := flag.String("file", "", "Path to the log file for static analysis")
	tailMode := flag.String("tail", "", "Path to the log file for live monitoring")
	configFile := flag.String("config", "", "Path to the configuration file (YAML)")

	flag.Parse()

	if *tailMode != "" || *filePath != "" {
		printBanner()
		runCliMode(*filePath, *tailMode, *configFile)
		return
	}

	printBanner()
	runInteractiveMenu()
}

func runCliMode(file, tailFile, conf string) {

	if tailFile != "" {
		fmt.Printf("[*] Starting Live Monitoring Mode.\nTarget: %s\n", tailFile)

		var cfg *model.Config
		var err error

		if conf != "" {
			cfg, err = config.LoadConfig(conf)
			if err != nil {
				color.Red(" Failed to load config file: %v", err)
				return
			}
			color.Green("-> Rule set loaded: %s", conf)
		} else {
			color.Yellow("[!] WARNING: No config file provided. Monitoring raw logs only.")
		}

		analyzer.RunLiveMonitoring(tailFile, cfg)
		return
	}

	if file != "" {
		fmt.Printf("[*] Starting Static Analysis Mode.\nTarget: %s\n", file)
		analyzer.RunStaticAnalysis(file)
		return
	}
}

func runInteractiveMenu() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- MAIN MENU ---")
		fmt.Printf("[%s] Static File Analysis\n", cyan("1"))
		fmt.Printf("[%s] Live Monitoring \n", cyan("2"))
		fmt.Printf("[%s] Exit\n", cyan("3"))

		fmt.Print("\n" + yellow("Select Option > "))

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:

			fmt.Print("Enter Log File Path (e.g., server.log): ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)

			if path == "" {
				fmt.Println(red("Invalid path!"))
				continue
			}
			analyzer.RunStaticAnalysis(path)

		case 2:

			fmt.Print("Enter Log File Path to Monitor: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)

			fmt.Print("Enter Rules File Path (Leave empty for Raw Mode): ")
			rulePath, _ := reader.ReadString('\n')
			rulePath = strings.TrimSpace(rulePath)

			var cfg *model.Config
			var err error

			if rulePath != "" {
				cfg, err = config.LoadConfig(rulePath)
				if err != nil {
					fmt.Printf(red("Error loading rules: %v\n"), err)
					continue
				}
				fmt.Println(green("-> Rules loaded successfully."))
			} else {
				fmt.Println(yellow("-> No rules provided. Running in Raw Mode."))
			}

			if path != "" {
				analyzer.RunLiveMonitoring(path, cfg)
			} else {
				fmt.Println(red("Invalid path!"))
			}

		case 3:
			fmt.Println(green("Exiting Tepegoz. Stay safe!"))
			os.Exit(0)

		default:
			fmt.Println(red("Invalid selection. Please try again."))
		}
	}
}

func printBanner() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	banner := `
        .---.
       /     \    ` + cyan("TEPEGOZ") + ` v1.0.0
      |  ` + red("(@)") + `  |   ` + white("The Eye On Your System") + `
       \  ^  /    ------------------------------------
        \___/     github.com/MESLEKDAA/tepegoz
       _ | | _
      (__| |__)
    `
	fmt.Println(banner)

	quotes := []string{
		"Watching the unseen...",
		"Logs never lie.",
		"Scanning for anomalies...",
		"System integrity check initiated...",
	}
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("[*] %s\n", quotes[rand.Intn(len(quotes))])
}
