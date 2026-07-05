package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

type SummaryResult struct {
	TotalLines  int
	LogCounts   map[string]int
	IPCounts    map[string]int
	TimeBuckets map[string]int
}

func RunStaticAnalysis(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		color.Red("ERROR Cannot open file: %v", err)
		return
	}
	defer file.Close()

	stats := SummaryResult{
		LogCounts:   make(map[string]int),
		IPCounts:    make(map[string]int),
		TimeBuckets: make(map[string]int),
	}

	ipRegex := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

	scanner := bufio.NewScanner(file)
	startTime := time.Now()

	fmt.Println(" -> Analyzing file statistics...")

	for scanner.Scan() {
		line := scanner.Text()
		stats.TotalLines++

		upperLine := strings.ToUpper(line)
		if strings.Contains(upperLine, "ERROR") || strings.Contains(upperLine, "FAIL") || strings.Contains(upperLine, "CRITICAL") {
			stats.LogCounts["ERRORS"]++
		} else if strings.Contains(upperLine, "WARN") {
			stats.LogCounts["WARNINGS"]++
		} else {
			stats.LogCounts["INFO/OTHER"]++
		}

		foundIP := ipRegex.FindString(line)
		if foundIP != "" {
			stats.IPCounts[foundIP]++
		}

		timeVal := extractTime(line)
		if timeVal != "" {
			stats.TimeBuckets[timeVal]++
		}
	}

	saveAndPrintReport(filePath, stats, time.Since(startTime))
}

func extractTime(line string) string {
	re := regexp.MustCompile(`\d{2}:\d{2}`)
	return re.FindString(line)
}

func saveAndPrintReport(sourceFile string, stats SummaryResult, duration time.Duration) {

	reportDir := "reports"
	os.MkdirAll(reportDir, 0755)

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	reportFileName := filepath.Join(reportDir, fmt.Sprintf("summary_report_%s.txt", timestamp))

	f, err := os.Create(reportFileName)
	if err != nil {
		color.Red("ERROR Could not create report file: %v", err)
		return
	}
	defer f.Close()

	printBoth := func(title string, content string, colorFunc func(a ...interface{}) string) {

		if colorFunc != nil {
			fmt.Print(colorFunc(content))
		} else {
			fmt.Print(content)
		}

		f.WriteString(content)
	}

	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	lineSep := strings.Repeat("=", 40) + "\n"

	printBoth("", "\n"+lineSep, cyan)
	printBoth("", fmt.Sprintf("   STATIC ANALYSIS REPORT: %s\n", sourceFile), nil)
	printBoth("", lineSep, cyan)

	printBoth("", fmt.Sprintf(" Analysis Time         : %s\n", time.Now().Format(time.RFC1123)), nil)
	printBoth("", fmt.Sprintf(" Total Lines Processed : %d\n", stats.TotalLines), nil)
	printBoth("", fmt.Sprintf(" Duration              : %s\n", duration), nil)

	printBoth("", "\n--- Log Level Distribution ---\n", yellow)
	for k, v := range stats.LogCounts {
		printBoth("", fmt.Sprintf(" %-12s : %d\n", k, v), nil)
	}

	printBoth("", "\n--- Top 5 IP Addresses ---\n", yellow)
	type kv struct {
		Key   string
		Value int
	}
	var sortedIPs []kv
	for k, v := range stats.IPCounts {
		sortedIPs = append(sortedIPs, kv{k, v})
	}
	sort.Slice(sortedIPs, func(i, j int) bool { return sortedIPs[i].Value > sortedIPs[j].Value })

	count := 0
	for _, kv := range sortedIPs {
		if count >= 5 {
			break
		}

		msg := fmt.Sprintf(" %-15s : %d events\n", kv.Key, kv.Value)
		fmt.Print(green(fmt.Sprintf(" %-15s", kv.Key)) + fmt.Sprintf(" : %d events\n", kv.Value)) // Ekrana özel renkli bas
		f.WriteString(msg)                                                                        // Dosyaya düz bas
		count++
	}

	printBoth("", "\n--- Time Activity ---\n", yellow)
	var times []string
	for t := range stats.TimeBuckets {
		times = append(times, t)
	}
	sort.Strings(times)

	maxVal := 0
	for _, c := range stats.TimeBuckets {
		if c > maxVal {
			maxVal = c
		}
	}

	if len(times) > 0 {
		for _, t := range times {
			val := stats.TimeBuckets[t]
			barLen := 0
			if maxVal > 0 {
				barLen = (val * 20) / maxVal
			}
			bar := strings.Repeat("█", barLen)

			fmt.Printf(" [%s] %s %d\n", t, color.RedString(bar), val)

			f.WriteString(fmt.Sprintf(" [%s] %s %d\n", t, bar, val))
		}
	} else {
		printBoth("", " (No time data detected)\n", nil)
	}

	printBoth("", lineSep, cyan)
	color.Green("\nReport saved to: %s", reportFileName)
}
