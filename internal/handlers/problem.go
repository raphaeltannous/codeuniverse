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
	problemService services.ProblemService
}

func NewProblemsHandlers(pS services.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		problemService: pS,
	}
}

func (h *ProblemHandler) GetProblems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search, ok := ctx.Value(middleware.SearchCtxKey).(string)
	if !ok {
		search = ""
	}

	offset, ok := ctx.Value(middleware.OffsetCtxKey).(int)
	if !ok {
		offset = middleware.OffsetDefault
	}

	limit, ok := ctx.Value(middleware.LimitCtxKey).(int)
	if !ok {
		limit = middleware.LimitDefault
	}

	premium, ok := ctx.Value(middleware.ProblemPremiumFilterCtxKey).(repository.ProblemParam)
	if !ok {
		premium = 0
	}
	difficulty, ok := ctx.Value(middleware.ProblemDifficultyFilterCtxKey).(models.ProblemDifficulty)
	if !ok {
		difficulty = 0
	}
	sortBy, ok := ctx.Value(middleware.ProblemSortByFilterCtxKey).(repository.ProblemParam)
	if !ok {
		sortBy = 0
	}
	sortOrder, ok := ctx.Value(middleware.ProblemSortOrderFilterCtxKey).(repository.ProblemParam)
	if !ok {
		sortOrder = 0
	}

	getParams := &repository.GetProblemsParams{
		Offset:     offset,
		Limit:      limit,
		Search:     search,
		IsPremium:  premium,
		IsPublic:   repository.ProblemPublic,
		Difficulty: difficulty,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	problems, total, err := h.problemService.GetProblems(ctx, getParams)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"problems": problems,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	hints, err := h.problemService.GetProblemHints(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	codeSnippets, err := h.problemService.GetProblemCodes(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	testcases, err := h.problemService.GetProblemTestcases(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := new(models.PublicProblem)
	response.Problem = problem

	response.Hints = make([]string, len(hints))
	for i, hint := range hints {
		response.Hints[i] = hint.Hint
	}

	response.CodeSnippets = make([]*models.ProblemCodeCodeSnippet, 0, len(codeSnippets))
	for _, codeSnippet := range codeSnippets {
		if codeSnippet.IsPublic {
			response.CodeSnippets = append(response.CodeSnippets, &models.ProblemCodeCodeSnippet{
				Code:     codeSnippet.CodeSnippet,
				Language: codeSnippet.Language,
			})
		}
	}

	response.Testcases = make([]*models.ProblemTestcase, 0, 3)
	for _, testcase := range testcases {
		if testcase.IsPublic {
			response.Testcases = append(response.Testcases, testcase)
		}
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted)
}

// PUT
func (h *ProblemHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

// DELETE
func (h *ProblemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	handlersutils.Unimplemented(w, r)
}

func (h *ProblemHandler) Run(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Code string `json:"code"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	language, ok := ctx.Value(middleware.ProblemLanguageCtxKey).(models.ProblemLanguage)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlerChannel := make(chan string)

	go func() {
		h.problemService.Run(
			context.WithoutCancel(ctx),
			user,
			problem,
			language,
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
		Code string `json:"code"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	language, ok := ctx.Value(middleware.ProblemLanguageCtxKey).(models.ProblemLanguage)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlerChannel := make(chan string)

	go func() {
		h.problemService.Submit(
			context.WithoutCancel(ctx),
			user,
			problem,
			language,
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

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	submissions, err := h.problemService.GetSubmissions(
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

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
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

	submission, err := h.problemService.GetSubmission(
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

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
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

	run, err := h.problemService.GetRun(
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

	if len(requestBody.Markdown) > 5000 {
		apiError := handlersutils.NewAPIError(
			"CHARACTERS_LIMIT_EXCEEDED",
			"Characters limit exceeded.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
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
	err := h.problemService.CreateNote(
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

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	note, err := h.problemService.GetNote(
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

	if len(requestBody.Markdown) > 5000 {
		apiError := handlersutils.NewAPIError(
			"CHARACTERS_LIMIT_EXCEEDED",
			"Characters limit exceeded.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	note, err := h.problemService.GetNote(ctx, user, problem)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.problemService.UpdateNote(ctx, note, requestBody.Markdown)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Problem note updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *ProblemHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	solvedProblems, err := h.problemService.GetSolvedProblems(
		ctx,
		user,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, solvedProblems, http.StatusOK)
}
