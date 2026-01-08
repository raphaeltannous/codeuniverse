package models

import "fmt"

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
)

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
	}

	if value, ok := statuses[r]; ok {
		return value
	}

	return "UNKNOWN RESULT STATUS"
}

func (r ResultStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", r), nil
}
