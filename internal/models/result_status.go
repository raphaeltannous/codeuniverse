package models

import (
	"errors"
	"fmt"
)

type Results struct {
	Exeption        string            `json:"exception"`
	IsPassed        bool              `json:"isPassed"`
	FailedTestcases []*FailedTestcase `json:"failedTestcases"`

	ExecutionTime float64 `json:"executionTime"`
}

type FailedTestcase struct {
	Id       int `json:"id"`
	Input    any `json:"input"`
	Expected any `json:"expected"`
	Got      any `json:"got"`
	StdOut   any `json:"stdOut"`
}

type ResultStatus int

const (
	StatusPending ResultStatus = iota + 1
	StatusStarted

	StatusAccepted
	StatusFailed

	StatusTimeLimitExceeded
	StatusMemoryLimitExceeded

	StatusCompileError
	StatusRuntimeError

	StatusInternalServerError
)

var (
	ErrInvalidResultStatus = errors.New("invalid result status")
)

func ParseResultStatus(status string) (ResultStatus, error) {
	switch status {
	case "Pending":
		return StatusPending, nil
	case "Started":
		return StatusStarted, nil

	case "Accepted":
		return StatusAccepted, nil
	case "Failed":
		return StatusFailed, nil

	case "Time Limit Exceeded":
		return StatusTimeLimitExceeded, nil
	case "Memory Limit Exceeded":
		return StatusMemoryLimitExceeded, nil

	case "Compile Error":
		return StatusCompileError, nil
	case "Runtime Error":
		return StatusRuntimeError, nil

	case "Internal Server Error":
		return StatusInternalServerError, nil

	default:
		return 0, ErrInvalidResultStatus
	}
}

func (r ResultStatus) String() string {
	statuses := map[ResultStatus]string{
		StatusPending: "Pending",
		StatusStarted: "Started",

		StatusAccepted: "Accepted",
		StatusFailed:   "Failed",

		StatusTimeLimitExceeded:   "Time Limit Exceeded",
		StatusMemoryLimitExceeded: "Memory Limit Exceeded",

		StatusCompileError: "Compile Error",
		StatusRuntimeError: "Runtime Error",

		StatusInternalServerError: "Internal Server Error",
	}

	if value, ok := statuses[r]; ok {
		return value
	}

	return "UNKNOWN RESULT STATUS"
}

func (r ResultStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", r), nil
}
