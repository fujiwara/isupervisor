package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/fujiwara/isupervisor"
)

func main() {
	ctx := context.TODO()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	app := isupervisor.Isupervisor{
		Recieve: func(ctx context.Context) (*isupervisor.Job, error) {
			return &isupervisor.Job{
				ID:          "1",
				TeamID:      "1",
				Command:     []string{"./bench.sh"},
				SoftTimeout: 60,
				HardTimeout: 120,
			}, nil
		},
		Report: func(ctx context.Context, report *isupervisor.Report) error {
			json.NewEncoder(log.Writer()).Encode(report)
			return nil
		},
	}
	return app.Run(ctx)
}
