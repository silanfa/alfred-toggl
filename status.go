package main

import (
	"fmt"
	"time"

	"github.com/jason0x43/go-alfred"
)

// StatusFilter is a command
type StatusFilter struct{}

// About returns information about this command
func (c StatusFilter) About() alfred.CommandDef {
	return alfred.CommandDef{
		Keyword:     "status",
		Description: "Show current status",
		IsEnabled:   config.APIKey != "",
	}
}

// Items returns a list of filter items
func (c StatusFilter) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("status items with arg=%s, data=%s", arg, data)

	if err = refresh(); err != nil {
		items = append(items, alfred.Item{
			Title:    "Error syncing with toggl.com",
			Subtitle: fmt.Sprintf("%v", err),
		})
		return
	}

	if entry, found := getRunningTimer(); found {
		startTime := entry.StartTime().Local()
		seconds := round(time.Now().Sub(startTime).Seconds())
		duration := float64(seconds) / float64(60*60)
		date := toHumanDateString(startTime)
		time := startTime.Format("15:04:05")
		subtitle := fmt.Sprintf("%s, started %s at %s",
			formatDuration(round(duration*100.0)), date, time)

		if project, _, ok := getProjectByID(entry.Pid); ok {
			subtitle = "[" + project.Name + "] " + subtitle
		}

		item := alfred.Item{
			Title:    entry.Description,
			Subtitle: subtitle,
			Arg: &alfred.ItemArg{
				Keyword: "timers",
				Data:    alfred.Stringify(timerCfg{Timer: &entry.ID}),
			},
		}

		item.AddMod(alfred.ModCmd, alfred.ItemMod{
			Subtitle: "Stop this timer",
			Arg: &alfred.ItemArg{
				Keyword: "timers",
				Mode:    alfred.ModeDo,
				Data:    alfred.Stringify(timerCfg{ToToggle: &toggleCfg{entry.ID, config.DurationOnly}}),
			},
		})

		items = append(items, item)
	} else {
		items = append(items, alfred.Item{
			Title: "No timers currently running",
			Icon:  "off.png",
		})
	}

	span1, _ := getSpan("today")
	var report1 *summaryReport
	report1, err = generateReport(span1.Start, span1.End, -1, "")
	for _, date1 := range report1.dates {
		items = append(items, alfred.Item{
			Title: fmt.Sprintf("Total time for today: %s", formatDuration(date1.total)),
		})
		break
	}

	span2, _ := getSpan("week")
	var report2 *summaryReport
	report2, err = generateReport(span2.Start, span2.End, -1, "")
	items = append(items, alfred.Item{
		Title: fmt.Sprintf("Total time for %s: %s",span2.Label, formatDuration(report2.total)),
	})

	return
}
