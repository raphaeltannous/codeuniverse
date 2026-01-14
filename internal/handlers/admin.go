package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	courseService  services.CourseService
	staticService  services.StaticService
	userService    services.UserService
	problemService services.ProblemService
}

func NewAdminHandler(
	courseService services.CourseService,
	staticService services.StaticService,
	userService services.UserService,
	problemService services.ProblemService,
) *AdminHandler {
	return &AdminHandler{
		courseService:  courseService,
		staticService:  staticService,
		userService:    userService,
		problemService: problemService,
	}
}

func (h *AdminHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	courses, err := h.courseService.GetAllCourses(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, courses, http.StatusOK)
}

func (h *AdminHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	course := &models.Course{
		Title:       requestBody.Title,
		Slug:        requestBody.Slug,
		Description: requestBody.Description,
		Difficulty:  requestBody.Difficulty,
	}

	ctx := r.Context()
	course, err := h.courseService.CreateCourse(ctx, course)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrCourseAlreadyExists):
			apiError.Code = "COURSE_ALREADY_EXISTS"
			apiError.Message = "Course slug already exists."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	handlersutils.WriteSuccessMessage(
		w,
		"Course created.",
		http.StatusCreated,
	)
}

func (h *AdminHandler) UpdateCourseInfo(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	patch := map[string]any{
		"title":       requestBody.Title,
		"slug":        requestBody.Slug,
		"description": requestBody.Description,
		"difficulty":  requestBody.Difficulty,
	}

	err := h.courseService.UpdateCourse(ctx, course, patch)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrInvalidPatch):
			apiError.Code = "INVALID_PATCH"
			apiError.Message = "Invalid patch."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	response := map[string]string{
		"message": "Course updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateCoursePublishStatus(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		IsPublished bool `json:"isPublished"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	patch := map[string]any{
		"isPublished": requestBody.IsPublished,
	}

	err := h.courseService.UpdateCourse(
		ctx,
		course,
		patch,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrInvalidPatch):
			apiError.Code = "INVALID_PATCH"
			apiError.Message = "Invalid patch."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	response := map[string]string{
		"message": "isPublished is updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.courseService.DeleteCourse(
		ctx,
		course,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteSuccessMessage(
		w,
		"Course deleted.",
		http.StatusOK,
	)
}

func (h *AdminHandler) UpdateThumbnail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "courseSlug")
	ctx := r.Context()

	course, err := h.courseService.GetCourseBySlug(ctx, slug)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_PARSE_FORM",
			"Failed to parse form.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"NO_THUMBNAIL_FILE_PROVIDED",
			"No thumbnail file provided.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	if !slices.Contains(allowedExts, ext) {
		apiError := handlersutils.NewAPIError(
			"INVALID_FILE_TYPE",
			"Invalid file type.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	if header.Size > 5*1024*1024 {
		apiError := handlersutils.NewAPIError(
			"FILE_TOO_LARGE",
			"File to large. (Max 5MB).",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	thumbnailUrl, err := h.staticService.SaveCourseThumbnail(
		ctx,
		file,
		ext,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.courseService.UpdateCourse(
		ctx,
		course,
		map[string]any{"thumbnailUrl": thumbnailUrl},
	)
	if err != nil {
		h.staticService.DeleteAvatar(ctx, thumbnailUrl)
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	h.staticService.DeleteCourseThumbnail(ctx, course.ThumbnailURL)

	handlersutils.WriteSuccessMessage(
		w,
		"Thumbnail updated.",
		http.StatusOK,
	)
}

func (h *AdminHandler) GetLessons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	lessons, err := h.courseService.GetCourseLessons(
		ctx,
		course,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"lessons":     lessons,
		"courseTitle": course.Title,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		LessonNumber int    `json:"lessonNumber"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	lesson := &models.Lesson{
		Title:        requestBody.Title,
		Description:  requestBody.Description,
		LessonNumber: requestBody.LessonNumber,
	}

	lesson, err := h.courseService.CreateLesson(
		ctx,
		course,
		lesson,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Lesson created.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lesson, ok := ctx.Value(middleware.LessonCtxKey).(*models.Lesson)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return

	}

	err := h.courseService.DeleteLesson(
		ctx,
		lesson,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Lesson deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		LessonNumber int    `json:"lessonNumber"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	patch := map[string]any{
		"title":        requestBody.Title,
		"description":  requestBody.Description,
		"lessonNumber": requestBody.LessonNumber,
	}

	ctx := r.Context()
	lesson, ok := ctx.Value(middleware.LessonCtxKey).(*models.Lesson)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.courseService.UpdateLesson(
		ctx,
		lesson,
		patch,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Lesson updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateLessonVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lesson, ok := ctx.Value(middleware.LessonCtxKey).(*models.Lesson)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := r.ParseMultipartForm(500 << 20)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_PARSE_FORM",
			"Failed to parse form.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"NO_VIDEO_FILE_PROVIDED",
			"No video file provided.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}
	defer file.Close()

	durationSecondsStr := r.FormValue("durationSeconds")
	durationSeconds, err := strconv.Atoi(durationSecondsStr)
	if err != nil || durationSeconds <= 0 {
		apiError := handlersutils.NewAPIError(
			"INVALID_DURATION",
			"Invalid duration.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".mp4"}
	if !slices.Contains(allowedExts, ext) {
		apiError := handlersutils.NewAPIError(
			"INVALID_FILE_TYPE",
			"Invalid file type.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	if header.Size > 500*1024*1024 {
		apiError := handlersutils.NewAPIError(
			"FILE_TOO_LARGE",
			"File to large. (Max 500MB).",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadGateway)
		return
	}

	videoUrl, err := h.staticService.SaveLessonVideo(
		ctx,
		file,
		ext,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.courseService.UpdateLesson(
		ctx,
		lesson,
		map[string]any{
			"videoUrl":        videoUrl,
			"durationSeconds": durationSeconds,
		},
	)
	if err != nil {
		h.staticService.DeleteAvatar(ctx, videoUrl)
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	h.staticService.DeleteCourseThumbnail(ctx, lesson.VideoURL)

	response := map[string]string{
		"videoUrl": videoUrl,
		"message":  "Lesson uploaded successfully.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search, ok := r.Context().Value(middleware.SearchCtxKey).(string)
	if !ok {
		search = ""
	}

	offset, ok := r.Context().Value("offset").(int)
	if !ok {
		offset = middleware.OffsetDefault
	}

	limit, ok := r.Context().Value("limit").(int)
	if !ok {
		limit = middleware.LimitDefault
	}

	role, ok := ctx.Value(middleware.UserRoleFilterCtxKey).(repository.UserParam)
	if !ok {
		role = 0
	}
	status, ok := ctx.Value(middleware.UserStatusFilterCtxKey).(repository.UserParam)
	if !ok {
		status = 0
	}
	verified, ok := ctx.Value(middleware.UserVerificationFilterCtxKey).(repository.UserParam)
	if !ok {
		verified = 0
	}
	sortBy, ok := ctx.Value(middleware.UserSortByFilterCtxKey).(repository.UserParam)
	if !ok {
		sortBy = 0
	}
	sortOrder, ok := ctx.Value(middleware.UserSortOrderFilterCtxKey).(repository.UserParam)
	if !ok {
		sortOrder = 0
	}

	getParams := &repository.GetUsersParams{
		Offset:     offset,
		Limit:      limit,
		Search:     search,
		Role:       role,
		IsActive:   status,
		IsVerified: verified,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	users, total, err := h.userService.GetAllUsers(ctx, getParams)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		Role       string `json:"role"`
		IsActive   bool   `json:"isActive"`
		IsVerified bool   `json:"isVerified"`
		AvatarUrl  string `json:"avatarUrl"`
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

	patch := make(map[string]any)
	addFunc := makePatchAdder(patch)
	addFunc("username", user.Username, requestBody.Username)
	addFunc("email", user.Email, requestBody.Email)
	addFunc("role", user.Role, requestBody.Role)
	addFunc("isActive", user.IsActive, requestBody.IsActive)
	addFunc("isVerified", user.IsVerified, requestBody.IsVerified)
	addFunc("avatarUrl", user.AvatarURL, requestBody.AvatarUrl)

	err := h.userService.UpdateUserPatch(
		ctx,
		user,
		patch,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message": "User is updated.",
		"updated": patch,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.userService.Delete(
		ctx,
		user,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message": "User is deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		IsActive   bool   `json:"isActive"`
		IsVerified bool   `json:"isVerified"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	user := &models.User{
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: requestBody.Password,
		Role:         requestBody.Role,
		IsActive:     requestBody.IsActive,
		IsVerified:   requestBody.IsVerified,
	}

	user, err := h.userService.RegisterUser(
		ctx,
		user,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()

		switch err {
		case repository.ErrUserAlreadyExists:
			apiError.Code = "USER_ALREADY_EXISTS"
			apiError.Message = "User already exists."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusConflict)
		case services.ErrInvalidEmail, services.ErrInvalidSlug, services.ErrWeakPasswordLength:
			apiError.Code = "INVALID_CONSTRAINTS"
			apiError.Message = "Invalid constraints."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	response := map[string]string{
		"message": "User created.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
		IsPremium   bool   `json:"isPremium"`
		IsPublic    bool   `json:"isPublic"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	problem, err := models.NewProblem(
		requestBody.Title,
		requestBody.Slug,
		requestBody.Description,
		requestBody.Difficulty,
		requestBody.IsPremium,
		requestBody.IsPublic,
	)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INVALID_DIFFICULTY",
			"Invalid difficulty.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	problem, err = h.problemService.Create(
		ctx,
		problem,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch err {
		case repository.ErrProblemAlreadyExists:
			apiError.Code = "PROBLEM_ALREADY_EXISTS"
			apiError.Message = "Problem already exists"

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Problem is created.",
	}
	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
		IsPremium   bool   `json:"isPremium"`
		IsPublic    bool   `json:"isPublic"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	difficulty, err := models.NewProblemDifficulty(requestBody.Difficulty)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INVALID_DIFFICULTY",
			"Invalid difficulty.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	patch := make(map[string]any)
	addFunc := makePatchAdder(patch)
	addFunc("title", problem.Title, requestBody.Title)
	addFunc("slug", problem.Slug, requestBody.Slug)
	addFunc("description", problem.Description, requestBody.Description)
	addFunc("difficulty", problem.Difficulty, difficulty)
	addFunc("isPremium", problem.IsPremium, requestBody.IsPremium)
	addFunc("isPublic", problem.IsPublic, requestBody.IsPublic)

	err = h.problemService.UpdateProblem(
		ctx,
		problem,
		patch,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message":     "Problem updated.",
		"updatePatch": patch,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func makePatchAdder(patch map[string]any) func(key string, v1, v2 any) {
	return func(key string, v1, v2 any) {
		switch v1.(type) {
		case string:
			if sV2, ok := v2.(string); ok && v1 != sV2 && sV2 != "" {
				patch[key] = sV2
			}
		case bool:
			if sV2, ok := v2.(bool); ok && v1 != sV2 {
				patch[key] = sV2
			}
		case models.ProblemDifficulty:
			if sV2, ok := v2.(models.ProblemDifficulty); ok && v1 != sV2 {
				patch[key] = sV2
			}
		}
	}
}

func (h *AdminHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, problem, http.StatusOK)
}

func (h *AdminHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.problemService.Delete(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Problem deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetProblems(w http.ResponseWriter, r *http.Request) {
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
	public, ok := ctx.Value(middleware.ProblemPublicFilterCtxKey).(repository.ProblemParam)
	if !ok {
		public = 0
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
		IsPublic:   public,
		Difficulty: difficulty,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}
	problems, total, err := h.problemService.GetProblems(
		ctx,
		getParams,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	type responseProblem struct {
		*models.Problem
		Hints        []string                         `json:"hints"`
		CodeSnippets []*models.ProblemCodeCodeSnippet `json:"codeSnippets"`
	}
	responseProblems := make([]*responseProblem, 0)

	for _, problem := range problems {
		responseProblems = append(
			responseProblems,
			&responseProblem{
				Problem:      problem,
				Hints:        []string{},
				CodeSnippets: []*models.ProblemCodeCodeSnippet{},
			},
		)

	}
	response := map[string]any{
		"problems": responseProblems,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetSupportedLanguages(w http.ResponseWriter, r *http.Request) {
	type language struct {
		LanguageName string `json:"languageName"`
		LanguageSlug string `json:"languageSlug"`
	}

	response := make([]*language, models.LanguageEnd-1)
	for lang := models.LanguageGo; lang < models.LanguageEnd; lang++ {
		response[lang-1] = &language{
			LanguageName: lang.String(),
			LanguageSlug: lang.Slug(),
		}
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetProblemHints(w http.ResponseWriter, r *http.Request) {
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

	handlersutils.WriteResponseJSON(w, hints, http.StatusOK)
}

func (h *AdminHandler) CreateProblemHint(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Hint string `json:"hint"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	hint := &models.ProblemHint{
		Hint: requestBody.Hint,
	}

	err := h.problemService.CreateHint(
		ctx,
		problem,
		hint,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Hint created.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateProblemHint(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Hint string `json:"hint"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	hint, ok := ctx.Value(middleware.ProblemHintCtxKey).(*models.ProblemHint)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.problemService.UpdateHint(
		ctx,
		hint,
		requestBody.Hint,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Hint updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) DeleteProblemHint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hint, ok := ctx.Value(middleware.ProblemHintCtxKey).(*models.ProblemHint)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.problemService.DeleteHint(
		ctx,
		hint,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Hint deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetProblemCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problemCodes, err := h.problemService.GetProblemCodes(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, problemCodes, http.StatusOK)
}

func (h *AdminHandler) UpdateProblemCodeSnippet(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		CodeSnippet  string `json:"codeSnippet"`
		Checker      string `json:"checker"`
		Driver       string `json:"driver"`
		IsPublic     bool   `json:"isPublic"`
		LanguageName string `json:"languageName"`
		LanguageSlug string `json:"languageSlug"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
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

	problemCode := &models.ProblemCode{
		CodeSnippet: requestBody.CodeSnippet,
		Checker:     requestBody.Checker,
		Driver:      requestBody.Driver,
		IsPublic:    requestBody.IsPublic,
		Language:    language,
	}

	err := h.problemService.SaveProblemCode(
		ctx,
		problem,
		problemCode,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Code snippet updated.",
	}
	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetProblemTestcases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	problemTestcases, err := h.problemService.GetProblemTestcases(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, problemTestcases, http.StatusOK)
}

func (h *AdminHandler) CreateProblemTestcase(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Input    any  `json:"input"`
		Expected any  `json:"expected"`
		IsPublic bool `json:"isPublic"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	testcase := &models.ProblemTestcase{
		Input:    requestBody.Input,
		Expected: requestBody.Expected,
		IsPublic: requestBody.IsPublic,
	}

	err := h.problemService.AddProblemTestcase(
		ctx,
		problem,
		testcase,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Testcase added.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) UpdateProblemTestcase(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Input    any  `json:"input"`
		Expected any  `json:"expected"`
		IsPublic bool `json:"isPublic"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	testcaseId, ok := ctx.Value(middleware.ProblemCodeTestcaseIdCtxKey).(int)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	testcase := &models.ProblemTestcase{
		Input:    requestBody.Input,
		Expected: requestBody.Expected,
		IsPublic: requestBody.IsPublic,
	}

	err := h.problemService.UpdateProblemTestcase(
		ctx,
		problem,
		testcaseId,
		testcase,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Testcase updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) DeleteProblemTestcase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	testcaseId, ok := ctx.Value(middleware.ProblemCodeTestcaseIdCtxKey).(int)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.problemService.DeleteProblemTestcase(
		ctx,
		problem,
		testcaseId,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Testcase deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *AdminHandler) GetProblemConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	config, err := h.problemService.GetProblemCodeConfig(
		ctx,
		problem,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, config, http.StatusOK)
}

func (h *AdminHandler) UpdateProblemConfig(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		TimeLimit   int `json:"timeLimit"`
		MemoryLimit int `json:"memoryLimit"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	problem, ok := ctx.Value(middleware.ProblemCtxKey).(*models.Problem)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	config := &models.ProblemCodeConfig{
		TimeLimit:   requestBody.TimeLimit,
		MemoryLimit: requestBody.MemoryLimit,
	}
	if config.TimeLimit <= 0 || config.MemoryLimit <= 0 {
		apiError := handlersutils.NewAPIError(
			"INVALID_CONFIG",
			"Config value cannot be zero or negative.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	err := h.problemService.UpdateProblemCodeConfig(
		ctx,
		problem,
		config,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Problem config updated.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}
