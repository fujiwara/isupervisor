package isupervisor

import (
	"context"
)

const (
	// StatusPending is the status of a job that is waiting to be run.
	StatusPending = 0
	// StatusRunning is the status of a job that is currently running.
	StatusRunning = 1
	// StatusFinished is the status of a job that has finished running.
	StatusFinished = 2
)

type Isupervisor struct {
	Recieve func(context.Context) (*Job, error)
	Report  func(context.Context, *Report) error
}

type Job struct {
	ID          string   `json:"id"`
	TeamID      string   `json:"team_id"`
	Command     []string `json:"command"`
	SoftTimeout int64    `json:"soft_timeout"`
	HardTimeout int64    `json:"hard_timeout"`
}

type Report struct {
	Job      *Job   `json:"job"`
	Status   int    `json:"status"` // 0: pending, 1: running, 2: finished
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Result   Result `json:"result"`
}

type Result struct {
	Finished bool  `json:"finished"`
	Passed   bool  `json:"passed"`
	Score    int64 `json:"score"`
}
