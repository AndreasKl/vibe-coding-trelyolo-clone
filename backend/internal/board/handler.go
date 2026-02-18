package board

import (
	"net/http"
	"trello-clone/internal/auth"
	"trello-clone/internal/httputil"
)

type Handler struct {
	Store *Store
}

// Boards

func (h *Handler) ListBoards(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	boards, err := h.Store.ListBoards(r.Context(), u.ID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to list boards")
		return
	}
	if boards == nil {
		boards = []Board{}
	}
	httputil.JSON(w, http.StatusOK, boards)
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	var req struct {
		Name string `json:"name"`
	}
	if err := httputil.Decode(r, &req); err != nil || req.Name == "" {
		httputil.Error(w, http.StatusBadRequest, "name required")
		return
	}
	b, err := h.Store.CreateBoard(r.Context(), u.ID, req.Name)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to create board")
		return
	}
	httputil.JSON(w, http.StatusCreated, b)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")
	b, err := h.Store.GetBoard(r.Context(), id, u.ID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "board not found")
		return
	}
	httputil.JSON(w, http.StatusOK, b)
}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")
	if err := h.Store.DeleteBoard(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "board not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Columns

func (h *Handler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	boardID := r.PathValue("boardID")

	// Verify ownership
	_, err := h.Store.GetBoard(r.Context(), boardID, u.ID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "board not found")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := httputil.Decode(r, &req); err != nil || req.Name == "" {
		httputil.Error(w, http.StatusBadRequest, "name required")
		return
	}

	c, err := h.Store.CreateColumn(r.Context(), boardID, req.Name)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to create column")
		return
	}
	httputil.JSON(w, http.StatusCreated, c)
}

func (h *Handler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if _, err := h.Store.ColumnBoardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "column not found")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Position *int    `json:"position"`
	}
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.Store.UpdateColumn(r.Context(), id, req.Name, req.Position)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to update column")
		return
	}
	httputil.JSON(w, http.StatusOK, c)
}

func (h *Handler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if _, err := h.Store.ColumnBoardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "column not found")
		return
	}

	if err := h.Store.DeleteColumn(r.Context(), id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to delete column")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Cards

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	columnID := r.PathValue("columnID")

	if _, err := h.Store.ColumnBoardOwner(r.Context(), columnID, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "column not found")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := httputil.Decode(r, &req); err != nil || req.Title == "" {
		httputil.Error(w, http.StatusBadRequest, "title required")
		return
	}

	c, err := h.Store.CreateCard(r.Context(), columnID, req.Title, req.Description)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to create card")
		return
	}
	httputil.JSON(w, http.StatusCreated, c)
}

func (h *Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if err := h.Store.CardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "card not found")
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
	}
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.Store.UpdateCard(r.Context(), id, req.Title, req.Description)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to update card")
		return
	}
	httputil.JSON(w, http.StatusOK, c)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if err := h.Store.CardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "card not found")
		return
	}

	if err := h.Store.DeleteCard(r.Context(), id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to delete card")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MoveColumn(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if _, err := h.Store.ColumnBoardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "column not found")
		return
	}

	var req struct {
		Position int `json:"position"`
	}
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.Store.MoveColumn(r.Context(), id, req.Position)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to move column")
		return
	}
	httputil.JSON(w, http.StatusOK, c)
}

func (h *Handler) MoveCard(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromContext(r.Context())
	id := r.PathValue("id")

	if err := h.Store.CardOwner(r.Context(), id, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "card not found")
		return
	}

	var req struct {
		ColumnID string `json:"column_id"`
		Position int    `json:"position"`
	}
	if err := httputil.Decode(r, &req); err != nil || req.ColumnID == "" {
		httputil.Error(w, http.StatusBadRequest, "column_id required")
		return
	}

	// Verify target column ownership
	if _, err := h.Store.ColumnBoardOwner(r.Context(), req.ColumnID, u.ID); err != nil {
		httputil.Error(w, http.StatusNotFound, "target column not found")
		return
	}

	c, err := h.Store.MoveCard(r.Context(), id, req.ColumnID, req.Position)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to move card")
		return
	}
	httputil.JSON(w, http.StatusOK, c)
}
