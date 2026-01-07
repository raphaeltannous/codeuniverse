package models

type ProblemCodeConfig struct {
	TimeLimit   int `json:"timeLimit"`
	MemoryLimit int `json:"memoryLimit"`
}

func DefaultProblemCodeConfig() *ProblemCodeConfig {
	return &ProblemCodeConfig{
		TimeLimit:   10_000,
		MemoryLimit: 250,
	}
}

type ProblemTestcase struct {
	Id       int  `json:"id"`
	Input    any  `json:"input"`
	Expected any  `json:"expected"`
	IsPublic bool `json:"isPublic"`
}
