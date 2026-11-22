package feed

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Handle    string    `json:"handle"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	limit := 20

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			http.Error(
				w,
				"limit must be between 1 and 100",
				http.StatusBadRequest,
			)
			return
		}

		limit = parsed
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT
			p.id,
			u.handle,
			p.body,
			p.created_at,
			p.expires_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.expires_at > now()
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Printf("feed query failed: %v", err)
		http.Error(
			w,
			"unable to retrieve feed",
			http.StatusInternalServerError,
		)
		return
	}
	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.ID,
			&post.Handle,
			&post.Body,
			&post.CreatedAt,
			&post.ExpiresAt,
		)
		if err != nil {
			log.Printf("feed row scan failed: %v", err)
			http.Error(
				w,
				"unable to retrieve feed",
				http.StatusInternalServerError,
			)
			return
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("feed iteration failed: %v", err)
		http.Error(
			w,
			"unable to retrieve feed",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("feed response encoding failed: %v", err)
	}
}
