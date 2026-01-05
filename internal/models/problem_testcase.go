package models

type ProblemTestcases struct {
	TimeLimit   int               `json:"timeLimit"`
	MemoryLimit int               `json:"memoryLimit"`
	Testcases   []ProblemTestcase `json:"testCases"`
}

type ProblemTestcase struct {
	Input    any  `json:"input"`
	Expected any  `json:"expected"`
	IsPublic bool `json:"isPublic"`
}
