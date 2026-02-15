package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"whatsdot-aibuddy/backend/internal/openai"
	"whatsdot-aibuddy/backend/internal/store"
)

type Server struct {
	Store       *store.Store
	OpenAI      *openai.Client
	UploadDir   string
	AnalyzeMock bool
	Limiter     *DeviceLimiter
}

type apiResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type homeworkResp struct {
	ID             int64                `json:"id"`
	Mode           string               `json:"mode"`
	SourceImage    string               `json:"sourceImageUrl"`
	QuestionText   string               `json:"questionText"`
	SuggestedGrade string               `json:"suggestedGrade"`
	Result         openai.AnalyzeResult `json:"result"`
	SolvedAt       time.Time            `json:"solvedAt"`
}

var errOpenAIConfigMissing = errors.New("openai not configured: set OPENAI_API_KEY or enable ANALYZE_MOCK=true")

func (s *Server) Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(s.requestLogger())
	r.Use(s.cors())

	_ = os.MkdirAll(s.UploadDir, 0o755)
	r.Static("/uploads", s.UploadDir)

	r.GET("/health", func(c *gin.Context) {
		s.success(c, gin.H{"ok": true, "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")
	api.Use(s.withRateLimit())
	{
		api.POST("/homework/analyze", s.handleAnalyze)
		api.POST("/homework/:id/regenerate", s.handleRegenerate)
		api.GET("/history", s.handleHistory)
		api.GET("/history/:id", s.handleHistoryDetail)
	}
	return r
}

func (s *Server) handleAnalyze(c *gin.Context) {
	deviceID := deviceIDFromRequest(c)
	if deviceID == "" {
		s.fail(c, http.StatusBadRequest, 40001, "device_id required")
		return
	}

	mode := normalizeMode(c.PostForm("mode"))
	fileHeader, err := c.FormFile("image")
	if err != nil {
		s.fail(c, http.StatusBadRequest, 40002, "image file required")
		return
	}
	bytes, contentType, imageURL, err := s.readAndSaveUpload(fileHeader)
	if err != nil {
		s.fail(c, http.StatusBadRequest, 40003, err.Error())
		return
	}

	result, err := s.analyze(c, bytes, contentType, mode)
	if err != nil {
		log.Printf("[ERROR] analyze: %v", err)
		if errors.Is(err, errOpenAIConfigMissing) {
			s.fail(c, http.StatusInternalServerError, 50007, err.Error())
			return
		}
		s.fail(c, http.StatusBadGateway, 50001, "analyze failed")
		return
	}

	rec, err := s.Store.CreateHomework(c.Request.Context(), deviceID, mode, imageURL, result.QuestionText, result.SuggestedGrade, result)
	if err != nil {
		log.Printf("[ERROR] create homework: %v", err)
		s.fail(c, http.StatusInternalServerError, 50002, "save record failed")
		return
	}

	s.success(c, gin.H{"record": toHomeworkResp(rec)})
}

func (s *Server) handleRegenerate(c *gin.Context) {
	deviceID := deviceIDFromRequest(c)
	if deviceID == "" {
		s.fail(c, http.StatusBadRequest, 40001, "device_id required")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		s.fail(c, http.StatusBadRequest, 40004, "invalid id")
		return
	}
	mode := normalizeMode(c.PostForm("mode"))
	if mode == "guided" {
		if jmode := normalizeMode(c.Query("mode")); jmode != "guided" {
			mode = jmode
		}
	}

	rec, err := s.Store.GetHomeworkByIDAndDevice(c.Request.Context(), id, deviceID)
	if err != nil {
		if store.IsNotFound(err) {
			s.fail(c, http.StatusNotFound, 40401, "record not found")
			return
		}
		log.Printf("[ERROR] get homework: %v", err)
		s.fail(c, http.StatusInternalServerError, 50003, "query record failed")
		return
	}

	imgPath := s.localPathFromURL(rec.SourceImage)
	b, err := os.ReadFile(imgPath)
	if err != nil {
		s.fail(c, http.StatusBadRequest, 40005, "image source missing")
		return
	}

	result, err := s.analyze(c, b, "image/jpeg", mode)
	if err != nil {
		log.Printf("[ERROR] regenerate analyze: %v", err)
		if errors.Is(err, errOpenAIConfigMissing) {
			s.fail(c, http.StatusInternalServerError, 50007, err.Error())
			return
		}
		s.fail(c, http.StatusBadGateway, 50001, "analyze failed")
		return
	}

	updated, err := s.Store.UpdateHomeworkResult(c.Request.Context(), id, deviceID, mode, result.QuestionText, result.SuggestedGrade, result)
	if err != nil {
		if store.IsNotFound(err) {
			s.fail(c, http.StatusNotFound, 40401, "record not found")
			return
		}
		log.Printf("[ERROR] update homework: %v", err)
		s.fail(c, http.StatusInternalServerError, 50004, "update record failed")
		return
	}

	s.success(c, gin.H{"record": toHomeworkResp(updated)})
}

func (s *Server) handleHistory(c *gin.Context) {
	deviceID := deviceIDFromRequest(c)
	if deviceID == "" {
		s.fail(c, http.StatusBadRequest, 40001, "device_id required")
		return
	}
	items, err := s.Store.ListHistoryByDevice(c.Request.Context(), deviceID, 100)
	if err != nil {
		log.Printf("[ERROR] list history: %v", err)
		s.fail(c, http.StatusInternalServerError, 50005, "query history failed")
		return
	}
	s.success(c, gin.H{"items": items})
}

func (s *Server) handleHistoryDetail(c *gin.Context) {
	deviceID := deviceIDFromRequest(c)
	if deviceID == "" {
		s.fail(c, http.StatusBadRequest, 40001, "device_id required")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		s.fail(c, http.StatusBadRequest, 40004, "invalid id")
		return
	}

	rec, err := s.Store.GetHomeworkByIDAndDevice(c.Request.Context(), id, deviceID)
	if err != nil {
		if store.IsNotFound(err) {
			s.fail(c, http.StatusNotFound, 40401, "record not found")
			return
		}
		log.Printf("[ERROR] history detail: %v", err)
		s.fail(c, http.StatusInternalServerError, 50006, "query detail failed")
		return
	}

	s.success(c, gin.H{"record": toHomeworkResp(rec)})
}

func (s *Server) analyze(c *gin.Context, imageBytes []byte, contentType string, mode string) (openai.AnalyzeResult, error) {
	if s.AnalyzeMock {
		return mockResult(mode), nil
	}
	if s.OpenAI == nil || strings.TrimSpace(s.OpenAI.APIKey) == "" {
		return openai.AnalyzeResult{}, errOpenAIConfigMissing
	}
	return s.OpenAI.AnalyzeHomework(c.Request.Context(), imageBytes, contentType, mode)
}

func (s *Server) readAndSaveUpload(file *multipart.FileHeader) ([]byte, string, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", "", err
	}
	defer src.Close()

	b, err := io.ReadAll(io.LimitReader(src, 8*1024*1024))
	if err != nil {
		return nil, "", "", err
	}
	if len(b) == 0 {
		return nil, "", "", errors.New("empty file")
	}
	contentType := http.DetectContentType(b)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", "", errors.New("only image is supported")
	}
	ext := extByContentType(contentType)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(s.UploadDir, name)
	if err := os.WriteFile(fullPath, b, 0o644); err != nil {
		return nil, "", "", err
	}
	return b, contentType, "/uploads/" + name, nil
}

func (s *Server) localPathFromURL(imageURL string) string {
	imageURL = strings.TrimSpace(imageURL)
	if strings.HasPrefix(imageURL, "/uploads/") {
		return filepath.Join(s.UploadDir, filepath.Base(imageURL))
	}
	return filepath.Join(s.UploadDir, filepath.Base(imageURL))
}

func (s *Server) withRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.Limiter == nil {
			c.Next()
			return
		}
		deviceID := deviceIDFromRequest(c)
		if deviceID == "" {
			c.Next()
			return
		}
		if !s.Limiter.Allow(deviceID) {
			s.fail(c, http.StatusTooManyRequests, 42901, "rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%s] %s %s %d %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}

func (s *Server) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Device-Id")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *Server) success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, apiResp{Code: 0, Message: "ok", Data: data})
}

func (s *Server) fail(c *gin.Context, status int, code int, message string) {
	c.JSON(status, apiResp{Code: code, Message: message})
}

func deviceIDFromRequest(c *gin.Context) string {
	deviceID := strings.TrimSpace(c.GetHeader("X-Device-Id"))
	if deviceID == "" {
		deviceID = strings.TrimSpace(c.Query("device_id"))
	}
	if len(deviceID) > 128 {
		deviceID = deviceID[:128]
	}
	return deviceID
}

func extByContentType(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func normalizeMode(mode string) string {
	v := strings.TrimSpace(strings.ToLower(mode))
	switch v {
	case "guided", "detailed", "noanswer", "quick":
		return v
	default:
		return "guided"
	}
}

func toHomeworkResp(rec store.HomeworkRecord) homeworkResp {
	var parsed openai.AnalyzeResult
	if len(rec.ResultJSONRaw) > 0 {
		_ = json.Unmarshal(rec.ResultJSONRaw, &parsed)
	}
	return homeworkResp{
		ID:             rec.ID,
		Mode:           rec.Mode,
		SourceImage:    rec.SourceImage,
		QuestionText:   rec.QuestionText,
		SuggestedGrade: rec.Grade,
		Result:         parsed,
		SolvedAt:       rec.SolvedAt,
	}
}

func mockResult(mode string) openai.AnalyzeResult {
	prefix := "引导思考"
	switch mode {
	case "detailed":
		prefix = "详细讲解"
	case "noanswer":
		prefix = "不给答案"
	case "quick":
		prefix = "快速提示"
	}
	return openai.AnalyzeResult{
		QuestionText:     "24 × 15 = ?",
		SolutionThoughts: prefix + "：把 15 拆成 10 和 5，分别与 24 相乘后相加，过程比答案更重要。",
		ExplainToChild:   "我们先算 24×10，再算 24×5，最后把两个结果加起来。",
		ParentGuidance: []string{
			"你先说说为什么可以把 15 拆成 10 和 5？",
			"如果先算 24×5，你会怎么口算？",
			"两部分结果加起来前，先估一估答案大概是多少？",
		},
		ChildStuckPoints: []string{
			"容易忘记把两部分乘积相加。",
			"对两位数乘法拆分不熟悉。",
		},
		KnowledgePoints: []string{"两位数乘法", "乘法分配律", "口算与估算"},
		SuggestedGrade:  "三年级",
	}
}
