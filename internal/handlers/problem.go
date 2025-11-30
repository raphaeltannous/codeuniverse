package handlers

import (
	"errors"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
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
		IsPaid      bool   `json:"isPaid"`
		IsPublic    bool   `json:"isPublic"`

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

	response := map[string]string{
		"message": "Problem is created.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted)
}

// GET
// Optional params: offset limit search.
func (h *ProblemHanlder) GetProblems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	offset, ok := ctx.Value("offset").(int)
	if !ok {
		offset = middleware.OffsetDefault
	}

	limit, ok := ctx.Value("limit").(int)
	if !ok {
		limit = middleware.LimitDefault
	}

	problems, err := h.pS.GetAllProblems(ctx, offset, limit)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, problems, http.StatusAccepted)
}

// GET
func (h *ProblemHanlder) GetProblem(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "problemSlug")

	ctx := r.Context()

	problem, err := h.pS.GetBySlug(
		ctx,
		slug,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrProblemNotFound):
			apiError.Code = "PROBLEM_NOT_FOUND"
			apiError.Message = "Problem not found."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	handlersutils.WriteResponseJSON(w, problem, http.StatusAccepted)
}

// PUT
func (h *ProblemHanlder) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// DELETE
func (h *ProblemHanlder) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}
