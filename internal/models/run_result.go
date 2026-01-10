package models

type RunResult struct {
	Status ResultStatus

	FailedTestcases []*FailedTestcase

	MemoryUsage   float64
	ExecutionTime float64

	StdOut string
	StdErr string
}

func NewRunResult(results *Results) *RunResult {
	runResults := &RunResult{
		FailedTestcases: results.FailedTestcases,
		ExecutionTime:   results.ExecutionTime,
	}

	if results.IsPassed {
		runResults.Status = StatusAccepted
	} else {
		runResults.Status = StatusFailed
	}

	return runResults
}
