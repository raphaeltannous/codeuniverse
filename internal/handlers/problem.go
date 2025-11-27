package handlers

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

type ProblemHanlder struct {
	pS services.ProblemService
}

func NewProblemsHandlers(pS services.ProblemService) *ProblemHanlder {
	return &ProblemHanlder{
		pS: pS,
	}
}

// POST
func (h *ProblemHanlder) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
		IsPaid      bool   `json:"is_paid"`
		IsPublic    bool   `json:"is_public"`

		Hints []string `json:"hints"`

		CodeSnippets []models.CodeSnippet `json:"codeSnippets"`
		TestCases    []string             `json:"testCases"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	problem, err := models.NewProblem(
		requestBody.Title,
		requestBody.Description,
		requestBody.Difficulty,
		requestBody.IsPaid,
		requestBody.IsPublic,

		requestBody.Hints,

		requestBody.CodeSnippets,
		requestBody.TestCases,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	problem, err = h.pS.Create(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, problem, http.StatusAccepted)
}

// GET
// Optional params: offset limit search.
func (h *ProblemHanlder) GetAllProblems(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// GET
func (h *ProblemHanlder) GetProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// PUT
func (h *ProblemHanlder) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// DELETE
func (h *ProblemHanlder) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}
