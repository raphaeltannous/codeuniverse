package handlers

import (
	"context"
	"errors"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

type ProblemHandler struct {
	pS services.ProblemService
}

func NewProblemsHandlers(pS services.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		pS: pS,
	}
}

// POST
func (h *ProblemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
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
func (h *ProblemHandler) GetProblems(w http.ResponseWriter, r *http.Request) {
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
func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
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
func (h *ProblemHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// DELETE
func (h *ProblemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// POST
func (h *ProblemHandler) Run(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ProblemSlug  string `json:"problemSlug"`
		LanguageSlug string `json:"languageSlug"`
		Code         string `json:"code"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, err := h.pS.GetBySlug(
		ctx,
		requestBody.ProblemSlug,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlerChannel := make(chan string)

	go func() {
		h.pS.Run(
			context.WithoutCancel(ctx),
			user,
			problem,
			requestBody.LanguageSlug,
			requestBody.Code,
			handlerChannel,
		)
	}()

	runId := <-handlerChannel
	if runId == repository.ErrInternalServerError.Error() {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"runId": runId,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusCreated)
}

// POST
func (h *ProblemHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ProblemSlug  string `json:"problemSlug"`
		LanguageSlug string `json:"languageSlug"`
		Code         string `json:"code"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, err := h.pS.GetBySlug(
		ctx,
		requestBody.ProblemSlug,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlerChannel := make(chan string)

	go func() {
		h.pS.Submit(
			context.WithoutCancel(ctx),
			user,
			problem,
			requestBody.LanguageSlug,
			requestBody.Code,
			handlerChannel,
		)
	}()

	submissionId := <-handlerChannel
	if submissionId == repository.ErrInternalServerError.Error() {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"submissionId": submissionId,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusCreated)
}

// GET
func (h *ProblemHandler) GetSubmissions(w http.ResponseWriter, r *http.Request) {
	problemSlug := chi.URLParam(r, "problemSlug")
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, err := h.pS.GetBySlug(
		ctx,
		problemSlug,
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

	submissions, err := h.pS.GetSubmissions(
		ctx,
		user,
		problem,
	)

	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_GET_SUBMISSIONS.",
			"Failed to get submissions.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, submissions, http.StatusOK)
}
