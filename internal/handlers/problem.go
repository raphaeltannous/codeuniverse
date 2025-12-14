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
	"github.com/google/uuid"
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
	ctx := r.Context()

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
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

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
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

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
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
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
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

// GET
func (h *ProblemHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	submissionId, err := uuid.Parse(chi.URLParam(r, "submissionId"))
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INVALID_SUBMISSION_ID",
			"Invalid submission id.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	submission, err := h.pS.GetSubmission(
		ctx,
		user,
		submissionId,
	)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_GET_SUBMISSION.",
			"Failed to get submission.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, submission, http.StatusOK)
}

// GET
func (h *ProblemHandler) GetRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	runId, err := uuid.Parse(chi.URLParam(r, "runId"))
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INVALID_RUN_ID",
			"Invalid run id.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	run, err := h.pS.GetRun(
		ctx,
		user,
		runId,
	)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_GET_RUN.",
			"Failed to get run.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, run, http.StatusOK)
}

// POST
func (h *ProblemHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Markdown string `json:"markdown"`
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

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	note := models.NewProblemNote(user.ID, problem.ID, requestBody.Markdown)
	err := h.pS.CreateNote(
		ctx,
		note,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrProblemNoteAlreadyExists):
			apiError.Code = "PROBLEM_NOTE_ALREADY_EXISTS"
			apiError.Message = "Problem note already exists."
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Problem note created.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusCreated)
}

// DELETE
func (h *ProblemHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
}

// GET
func (h *ProblemHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	note, err := h.pS.GetNote(
		ctx,
		user,
		problem,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrProblemNoteNotFound):
			apiError.Code = "PROBLEM_NOTE_NOT_FOUND"
			apiError.Message = "Problem note not found."
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}
		return
	}

	handlersutils.WriteResponseJSON(w, note, http.StatusOK)
}

// PUT
func (h *ProblemHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Markdown string `json:"markdown"`
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

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	note, err := h.pS.GetNote(ctx, user, problem)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.pS.UpdateNote(ctx, note, requestBody.Markdown)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Problem note updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}
