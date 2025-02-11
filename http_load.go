package main

import (
	"os"
	"time"

	"net/http"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/lib"
)

func customTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"

		tgt.URL = "http://localhost:8080/request"
		merchantTransactionID := uuid.New().String()
		payload := `{
  "id": "` + merchantTransactionID + `",
  "github_username": "vegeta",
  "commit_hash": "vegeta#hash"
}`

		tgt.Body = []byte(payload)

		header := http.Header{}
		header.Add("Accept", "application/json")
		header.Add("Content-Type", "application/json")
		tgt.Header = header

		return nil
	}
}

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 1 * time.Minute

	targeter := customTargeter()
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Whatever name") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	reporter(os.Stdout)
}
