package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/domain"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/util"
)

type LinkDeps struct {
	Config    config.Config
	Logger    *log.Logger
	LinksRepo storage.LinksRepo
}

type createLinkRequest struct {
	LongURL     string     `json:"originalUrl"`
	CustomAlias *string    `json:"customAlias,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type createLinkResponse struct {
	Key              string     `json:"shortCode"`
	ShortURL         string     `json:"shortUrl"`
	LongURLCanonical string     `json:"originalUrl"`
	IsCustom         bool       `json:"isCustom"`
	CreatedAt        time.Time  `json:"createdAt"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
	Existing         bool       `json:"existing"`
}

func CreateLink(d LinkDeps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var req createLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
			return
		}
		if req.LongURL == "" {
			util.WriteError(w, http.StatusBadRequest, "invalid_url", "long_url is required")
			return
		}
		canon, err := domain.CanonicalizeURL(req.LongURL)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid_url", "must be a valid http/https URL")
			return
		}
		if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
			util.WriteError(w, http.StatusBadRequest, "expiry_in_past", "expires_at must be in the future")
			return
		}

		var link *domain.Link
		var existing bool

		if req.CustomAlias != nil && *req.CustomAlias != "" {
			alias := *req.CustomAlias
			if !domain.ValidateAlias(alias) {
				util.WriteError(w, http.StatusBadRequest, "invalid_alias", "alias must match [A-Za-z0-9_-]{3,32}")
				return
			}
			if domain.IsReserved(alias) {
				util.WriteError(w, http.StatusBadRequest, "reserved_key", "alias is reserved")
				return
			}
			link, err = d.LinksRepo.CreateAlias(r.Context(), alias, canon, req.ExpiresAt)
			if err != nil {
				if errors.Is(err, domain.ErrAliasInUse) {
					util.WriteError(w, http.StatusConflict, "alias_in_use", "alias already taken")
					return
				}
				d.Logger.Printf("create alias error: %v", err)
				util.WriteError(w, http.StatusInternalServerError, "server_error", "could not create link")
				return
			}
		} else {
			// idempotent path
			if l, err := d.LinksRepo.GetSystemByCanonicalURL(r.Context(), canon); err == nil {
				link = l
				existing = true
			} else {
				link, err = d.LinksRepo.CreateSystem(r.Context(), canon, req.ExpiresAt)
				if err != nil {
					d.Logger.Printf("create system error: %v", err)
					util.WriteError(w, http.StatusInternalServerError, "server_error", "could not create link")
					return
				}
			}
		}

		resp := createLinkResponse{
			Key:              link.Key,
			ShortURL:         d.Config.BaseURL + "/" + link.Key,
			LongURLCanonical: link.LongURL,
			IsCustom:         link.IsCustom,
			CreatedAt:        link.CreatedAt,
			ExpiresAt:        link.ExpiresAt,
			Existing:         existing,
		}
		status := http.StatusCreated
		if existing {
			status = http.StatusOK
		}
		util.WriteJSON(w, status, resp)
	})
}
