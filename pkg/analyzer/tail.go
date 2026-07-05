package analyzer

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/MESLEKDAA/tepegoz/pkg/model"
	"github.com/MESLEKDAA/tepegoz/pkg/reporter"

	"github.com/fatih/color"
	"github.com/hpcloud/tail"
)

type CompiledRule struct {
	RuleDef model.Rule
	Regexp  *regexp.Regexp
}

func RunLiveMonitoring(filePath string, config *model.Config) {

	var activeRules []CompiledRule

	if config != nil {
		fmt.Println("Engine is starting...")

		for _, r := range config.Rules {
			re, err := regexp.Compile(r.Regex)
			if err != nil {
				color.Red("[!] Invalid Regex Rule: %s Skipping", r.Name)
				continue
			}
			activeRules = append(activeRules, CompiledRule{RuleDef: r, Regexp: re})
		}
		color.Green(" -> Live Mode Active! %d Rule Watching.", len(activeRules))
	} else {
		color.Yellow(" -> Warning: Ruleset not loaded. Only Row Log Stream will be Monitored.")
	}

	t, err := tail.TailFile(filePath, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
	})

	if err != nil {
		color.Red("Live Mode has not started: %v", err)
		return
	}

	fmt.Printf(" -> %s is Monitoring... (PRESS CTRL+C to Stop)\n\n", filePath)

	for line := range t.Lines {
		text := line.Text

		fmt.Printf("[DEBUG]: '%s'\n", text)

		if len(activeRules) == 0 {
			fmt.Println(text)
			continue
		}

		for _, cr := range activeRules {
			if cr.Regexp.MatchString(text) {

				now := time.Now().Format("15:04:05")
				alertMsg := fmt.Sprintf("[%s] [ALERT: %s] %s", now, cr.RuleDef.Name, strings.TrimSpace(text))
				color.Red(alertMsg)

				event := model.LogEvent{
					TimeStamp:    time.Now(),
					Level:        cr.RuleDef.Level,
					RuleID:       cr.RuleDef.ID,
					RuleName:     cr.RuleDef.Name,
					SourceFile:   filePath,
					OriginalLine: text,
				}

				dateStr := time.Now().Format("2006-01-02")
				csvFileName := fmt.Sprintf("alerts_%s.csv", dateStr)

				if err := reporter.WriteToCsv(csvFileName, event); err != nil {
					fmt.Printf("CSV ERROR %v\n", err)
				}
			}
		}
	}
}
