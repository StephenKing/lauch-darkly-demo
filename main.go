package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"log"
	"os"
	"time"
)

// Set sdkKey to your LaunchDarkly SDK key.
const sdkKey = "sdk-fe492bc2-e3fe-4cfc-9d81-66358f24cdca"

// Set featureFlagKey to the feature flag key you want to evaluate.
const featureFlagKey = "cloudland-feature"

// steffen-test-boolean-flag flips the first user example-user-1 at 12% -> perfect

func showMessage(s string) { fmt.Printf("*** %s\n", s) }

func main() {
	if sdkKey == "" {
		showMessage("Please edit main.go to set sdkKey to your LaunchDarkly SDK key first")
		os.Exit(1)
	}

	ldClient, _ := ld.MakeClient(sdkKey, 5*time.Second)
	if ldClient.Initialized() {
		showMessage("SDK successfully initialized!")
	} else {
		showMessage("SDK failed to initialize")
	}

	count := 50

	for {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetAutoIndex(true)

		userId := 1
		start := time.Now()

		for row := 1; row <= int(count/10); row++ {

			rowData := make([]interface{}, 10)

			for column := 1; column <= 10; column++ {

				// user-1, user-2 etc.
				userKey := fmt.Sprintf("user-%d", userId)

				// a User object for LD
				user := lduser.NewUserBuilder(userKey).
					Name(fmt.Sprintf("User %d", userId)).
					Custom("row", ldvalue.Int(row)).
					Custom("column", ldvalue.Int(column)).
					Build()

				// query the LD SDK
				flagValue, err := ldClient.StringVariation(featureFlagKey, user, "OFF")
				if err != nil {
					showMessage("error: " + err.Error())
				}
				rowData[column-1] = flagValue
				userId++
			}

			t.AppendRow(rowData)
		}

		// clear screen
		fmt.Println("\033[2J")

		log.Printf("- %s", time.Since(start))

		t.SetStyle(table.StyleLight)
		t.Render()

		time.Sleep(1 * time.Second)
	}

	ldClient.Close()
}
