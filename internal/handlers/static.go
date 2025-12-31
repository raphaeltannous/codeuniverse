package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

type StaticHandler struct {
	staticService services.StaticService
}

func NewStaticHandler(
	staticService services.StaticService,
) *StaticHandler {
	return &StaticHandler{
		staticService: staticService,
	}
}

// GET
func (s *StaticHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	ctx := r.Context()

	avatarPath, err := s.staticService.GetAvatar(ctx, filename)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"AVATAR_NOT_FOUND",
			"Avatar not found.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, avatarPath)
}

// GET
func (s *StaticHandler) GetCourseThumbnail(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	ctx := r.Context()

	thumbnailPath, err := s.staticService.GetCourseThumbnail(ctx, filename)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"THUMBNAIL_NOT_FOUND",
			"Thumbnail not found.",
		)
		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, thumbnailPath)
}

func (h *StaticHandler) StreamVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lesson, ok := ctx.Value(middleware.LessonCtxKey).(*models.Lesson)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	expiresStr := r.URL.Query().Get("expires")
	signature := r.URL.Query().Get("signature")

	expiresAt, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INVALID_PARAMETER",
			"Invalid parameter.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	if time.Now().Unix() > expiresAt || !utils.ValidateSignedUrl(lesson.ID.String(), signature, int(expiresAt)) {
		apiError := handlersutils.NewAPIError(
			"LINK_EXPIRED",
			"Link expired.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		return
	}

	videoPath, err := h.staticService.GetLessonVideo(
		ctx,
		lesson.VideoURL,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrNotFound):
			apiError.Code = "NO_VIDEO"
			apiError.Message = "No video for this lesson."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, videoPath)
}
