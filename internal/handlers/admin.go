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
	courseService services.CourseService
	staticService services.StaticService
}

func NewAdminHandler(
	courseService services.CourseService,
	staticService services.StaticService,
) *AdminHandler {
	return &AdminHandler{
		courseService: courseService,
		staticService: staticService,
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

	handlersutils.WriteResponseJSON(w, course, http.StatusCreated)
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
		"message": "Course is updated.",
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

	response := map[string]string{
		"message": "Course deleted.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
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

	response := map[string]string{
		"thumbnailUrl": thumbnailUrl,
		"message":      "Thumbnail uploaded successfully.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
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
