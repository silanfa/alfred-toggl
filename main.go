package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl] ", log.LstdFlags)

var cacheFile string
var configFile string
var config struct {
	APIKey           string `desc:"Toggl API key"`
	DurationOnly     bool   `desc:"Extend time entries instead of creating new ones."`
	Rounding         int    `desc:"Minutes to round to, 0 to disable rounding."`
	DefaultProjectID int    `desc:"Optional default project ID; set to 0 to clear"`
	TestMode         bool   `desc:"If true, disable auto refresh"`
}
var cache struct {
	Workspace int
	Account   Account
	Time      time.Time
}
var workflow alfred.Workflow

func main() {
	if !alfred.IsDebugging() {
		dlog.SetOutput(ioutil.Discard)
		dlog.SetFlags(0)
	}

	var err error
	if workflow, err = alfred.OpenWorkflow(".", true); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	configFile = path.Join(workflow.DataDir(), "config.json")
	cacheFile = path.Join(workflow.CacheDir(), "cache.json")

	dlog.Printf("Using config file: %s", configFile)
	dlog.Printf("Using cache file: %s", cacheFile)

	if err := alfred.LoadJSON(configFile, &config); err != nil {
		dlog.Println("Error loading config:", err)
	}

	if err := alfred.LoadJSON(cacheFile, &cache); err != nil {
		dlog.Println("Error loading config:", err)
	}

	workflow.Run([]alfred.Command{
		StatusFilter{},
		LoginCommand{},
		TokenCommand{},
		TimeEntryCommand{},
		ProjectCommand{},
		TagCommand{},
		ReportFilter{},
		OptionsCommand{},
		LogoutCommand{},
		ResetCommand{},
	})
}
