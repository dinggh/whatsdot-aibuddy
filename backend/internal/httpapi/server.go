package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"whatsdot-aibuddy/backend/internal/auth"
	"whatsdot-aibuddy/backend/internal/store"
	"whatsdot-aibuddy/backend/internal/wechat"
)

type Server struct {
	Store  *store.Store
	JWT    *auth.JWT
	WeChat *wechat.Client
}

type ctxKey string

const userIDKey ctxKey = "uid"

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/auth/wechat/login", s.handleWeChatLogin)
	mux.HandleFunc("/api/v1/auth/wechat/profile", s.withAuth(s.handleUpdateProfile))
	mux.HandleFunc("/api/v1/auth/wechat/phone", s.withAuth(s.handleBindPhone))
	mux.HandleFunc("/api/v1/me", s.withAuth(s.handleMe))
	mux.HandleFunc("/api/v1/history", s.withAuth(s.handleHistory))
	return s.withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().Format(time.RFC3339)})
}

type loginReq struct {
	Code string `json:"code"`
}

func (s *Server) handleWeChatLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.WeChat.AppID == "" || s.WeChat.Secret == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "WECHAT_APP_ID/WECHAT_APP_SECRET not configured"})
		return
	}

	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	sess, err := s.WeChat.Code2Session(r.Context(), strings.TrimSpace(req.Code))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	u, err := s.Store.UpsertUserByOpenID(r.Context(), sess.OpenID, sess.UnionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}

	tk, err := s.JWT.Sign(u.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": tk, "user": u})
}

type profileReq struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
}

func (s *Server) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := userIDFromContext(r.Context())

	var req profileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	nick := strings.TrimSpace(req.NickName)
	if nick == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "nickName required"})
		return
	}

	u, err := s.Store.UpdateUserProfile(r.Context(), uid, nick, strings.TrimSpace(req.AvatarURL))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}

type bindPhoneReq struct {
	Code string `json:"code"`
}

func (s *Server) handleBindPhone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := userIDFromContext(r.Context())

	var req bindPhoneReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	phone, err := s.WeChat.GetPhoneNumberByCode(r.Context(), strings.TrimSpace(req.Code))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	u, err := s.Store.UpdateUserPhone(r.Context(), uid, phone)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := userIDFromContext(r.Context())
	u, err := s.Store.GetUserByID(r.Context(), uid)
	if err != nil {
		status := http.StatusInternalServerError
		if store.IsNotFound(err) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]any{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := userIDFromContext(r.Context())
	items, err := s.Store.ListHistory(r.Context(), uid, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authz := strings.TrimSpace(r.Header.Get("Authorization"))
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing token"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		claims, err := s.JWT.Parse(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid token"})
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func userIDFromContext(ctx context.Context) int64 {
	v := ctx.Value(userIDKey)
	if id, ok := v.(int64); ok {
		return id
	}
	panic("missing user id in context")
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": fmt.Sprintf("internal error: %v", rec)})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var ErrBadRequest = errors.New("bad request")
