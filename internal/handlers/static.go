package handlers

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/services"
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
