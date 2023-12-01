package isupervisor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/Songmu/wrapcommander"
)

var sleep = time.Second * 1
var ReportTimeout = time.Second * 5

func (isu *Isupervisor) Run(ctx context.Context) error {
	ticker := time.NewTicker(sleep)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
		if err := isu.runJob(ctx); err != nil {
			log.Println("[error] Error while running job:", err)
		}
	}
}

func (isu *Isupervisor) runJob(ctx context.Context) error {
	job, err := isu.Recieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to recieve job: %w", err)
	}
	log.Printf("[info] Recieved job: %+v", job)
	report := &Report{
		Job:    job,
		Status: StatusRunning,
	}
	// report書きだし用のtmpfileを作る
	tmpfile, err := os.CreateTemp("", "isupervisor-")
	if err != nil {
		return fmt.Errorf("failed to create tmpfile: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	// ctxを継承すると親のctxがキャンセルされたときに子もキャンセルされるので
	// benchmarkerはctxを継承しない
	jobCtx, cancelJob := context.WithTimeout(context.Background(), time.Duration(job.SoftTimeout)*time.Second)
	defer cancelJob()
	reportCtx, cancelReport := context.WithTimeout(context.Background(), time.Duration(job.HardTimeout)*time.Second+ReportTimeout)
	defer cancelReport()

	var cmd *exec.Cmd
	command := job.Command
	switch len(command) {
	case 0:
		return fmt.Errorf("Job command is empty")
	case 1:
		cmd = exec.CommandContext(jobCtx, command[0])
	default:
		cmd = exec.CommandContext(jobCtx, command[0], command[1:]...)
	}
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.Env = append(os.Environ(), "ISUPERVISOR_RESULT="+tmpfile.Name())
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = time.Duration(job.HardTimeout-job.SoftTimeout) * time.Second
	log.Println("[info] Start job:", report.Job.ID)

	if err = cmd.Run(); err != nil {
		report.ExitCode = wrapcommander.ResolveExitCode(err)
	}
	report.Stdout = stdoutBuf.String()
	report.Stderr = stderrBuf.String()

	result, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		log.Println("[warn] Error while reading result file:", err)
	} else {
		if err := json.Unmarshal(result, &report.Result); err != nil {
			log.Println("[error] Error while unmarshaling result file:", err)
		}
	}
	log.Println("[info] Finished job:", report.Job.ID)

	if err := isu.Report(reportCtx, report); err != nil {
		log.Println("[error] Error while sending report:", err)
	}
	return nil
}
