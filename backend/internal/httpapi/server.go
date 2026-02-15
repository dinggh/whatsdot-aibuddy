package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"whatsdot-aibuddy/backend/internal/auth"
	"whatsdot-aibuddy/backend/internal/store"
	"whatsdot-aibuddy/backend/internal/wechat"
)

type Server struct {
	Store          *store.Store
	JWT            *auth.JWT
	WeChat         *wechat.Client
	ForceDevWeChat bool
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
	return s.requestLogger(s.withCORS(mux))
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

	body, _ := io.ReadAll(r.Body)
	var req loginReq
	_ = json.Unmarshal(body, &req)

	// Dev bypass: when ForceDevWeChat or credentials are placeholder, return mock session for local testing
	if s.ForceDevWeChat || isDevWeChatCreds(s.WeChat.AppID, s.WeChat.Secret) {
		openID := strings.TrimSpace(req.Code)
		if openID == "" {
			openID = "dev_local"
		} else {
			openID = "dev_" + openID
		}
		u, err := s.Store.UpsertUserByOpenID(r.Context(), openID, "")
		if err != nil {
			log.Printf("[ERROR] wechat login dev bypass upsert: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
			return
		}
		tk, err := s.JWT.Sign(u.ID)
		if err != nil {
			log.Printf("[ERROR] wechat login dev bypass sign: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token error"})
			return
		}
		log.Printf("[INFO] wechat login dev bypass ok user_id=%d", u.ID)
		writeJSON(w, http.StatusOK, map[string]any{"token": tk, "user": u})
		return
	}

	if s.WeChat.AppID == "" || s.WeChat.Secret == "" {
		log.Printf("[ERROR] wechat login: WECHAT_APP_ID/WECHAT_APP_SECRET not configured")
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "WECHAT_APP_ID/WECHAT_APP_SECRET not configured"})
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		log.Printf("[WARN] wechat login: code required")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "code required"})
		return
	}

	sess, err := s.WeChat.Code2Session(r.Context(), strings.TrimSpace(req.Code))
	if err != nil {
		log.Printf("[ERROR] wechat login code2session: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	u, err := s.Store.UpsertUserByOpenID(r.Context(), sess.OpenID, sess.UnionID)
	if err != nil {
		log.Printf("[ERROR] wechat login upsert: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}

	tk, err := s.JWT.Sign(u.ID)
	if err != nil {
		log.Printf("[ERROR] wechat login sign: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token error"})
		return
	}
	log.Printf("[INFO] wechat login ok user_id=%d openid=%s", u.ID, sess.OpenID)
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
		log.Printf("[WARN] update profile invalid json: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	nick := strings.TrimSpace(req.NickName)
	if nick == "" {
		log.Printf("[WARN] update profile: nickName required user_id=%d", uid)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "nickName required"})
		return
	}

	u, err := s.Store.UpdateUserProfile(r.Context(), uid, nick, strings.TrimSpace(req.AvatarURL))
	if err != nil {
		log.Printf("[ERROR] update profile: %v user_id=%d", err, uid)
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
		log.Printf("[WARN] bind phone invalid json: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	phone, err := s.WeChat.GetPhoneNumberByCode(r.Context(), strings.TrimSpace(req.Code))
	if err != nil {
		log.Printf("[ERROR] bind phone get: %v user_id=%d", err, uid)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	u, err := s.Store.UpdateUserPhone(r.Context(), uid, phone)
	if err != nil {
		log.Printf("[ERROR] bind phone update: %v user_id=%d", err, uid)
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
			log.Printf("[WARN] get me: user not found user_id=%d", uid)
		} else {
			log.Printf("[ERROR] get me: %v user_id=%d", err, uid)
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
		log.Printf("[ERROR] list history: %v user_id=%d", err, uid)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authz := strings.TrimSpace(r.Header.Get("Authorization"))
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			log.Printf("[WARN] auth: missing token path=%s", r.URL.Path)
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing token"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		claims, err := s.JWT.Parse(token)
		if err != nil {
			log.Printf("[WARN] auth: invalid token path=%s err=%v", r.URL.Path, err)
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid token"})
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		dur := time.Since(start)
		log.Printf("[%s] %s %s %d %d %s", r.RemoteAddr, r.Method, r.URL.Path, rec.status, rec.bytes, dur)
		if rec.status >= 400 {
			log.Printf("[ERROR] %s %s status=%d", r.Method, r.URL.Path, rec.status)
		}
	})
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
				log.Printf("[ERROR] panic recovered: %v path=%s", rec, r.URL.Path)
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": fmt.Sprintf("internal error: %v", rec)})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var ErrBadRequest = errors.New("bad request")

func isDevWeChatCreds(appID, secret string) bool {
	return appID == "" || secret == "" ||
		strings.Contains(appID, "your_app") || strings.Contains(secret, "your_wechat") ||
		appID == "wx_your_app_id" || secret == "your_wechat_secret"
}
