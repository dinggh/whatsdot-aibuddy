package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	DB *pgxpool.Pool
}

type User struct {
	ID          int64     `json:"id"`
	OpenID      string    `json:"-"`
	UnionID     string    `json:"-"`
	NickName    string    `json:"nickName"`
	AvatarURL   string    `json:"avatarUrl"`
	PhoneNumber string    `json:"phoneNumber"`
	UsedCount   int       `json:"usedCount"`
	RemainCount int       `json:"remainingCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type HistoryItem struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Grade        string    `json:"grade"`
	ThumbURL     string    `json:"thumbUrl"`
	Summary      string    `json:"summary"`
	Mode         string    `json:"mode"`
	SolvedAt     time.Time `json:"solvedAt"`
	QuestionText string    `json:"questionText"`
}

type HomeworkRecord struct {
	ID            int64           `json:"id"`
	DeviceID      string          `json:"deviceId"`
	Mode          string          `json:"mode"`
	Title         string          `json:"title"`
	Grade         string          `json:"grade"`
	ThumbURL      string          `json:"thumbUrl"`
	SourceImage   string          `json:"sourceImageUrl"`
	Summary       string          `json:"summary"`
	QuestionText  string          `json:"questionText"`
	ResultJSONRaw json.RawMessage `json:"result"`
	SolvedAt      time.Time       `json:"solvedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func (s *Store) UpsertUserByOpenID(ctx context.Context, openID, unionID string) (User, error) {
	const q = `
INSERT INTO users (openid, unionid)
VALUES ($1, $2)
ON CONFLICT (openid)
DO UPDATE SET unionid = COALESCE(EXCLUDED.unionid, users.unionid), updated_at = now()
RETURNING id, openid, COALESCE(unionid, ''), nick_name, avatar_url, COALESCE(phone_number, ''), used_count, remaining_count, created_at, updated_at`

	var u User
	err := s.DB.QueryRow(ctx, q, openID, nullable(unionID)).Scan(
		&u.ID, &u.OpenID, &u.UnionID, &u.NickName, &u.AvatarURL,
		&u.PhoneNumber, &u.UsedCount, &u.RemainCount, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (User, error) {
	const q = `
SELECT id, openid, COALESCE(unionid, ''), nick_name, avatar_url, COALESCE(phone_number, ''), used_count, remaining_count, created_at, updated_at
FROM users WHERE id = $1`
	var u User
	err := s.DB.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.OpenID, &u.UnionID, &u.NickName, &u.AvatarURL,
		&u.PhoneNumber, &u.UsedCount, &u.RemainCount, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Store) UpdateUserProfile(ctx context.Context, id int64, nickName, avatarURL string) (User, error) {
	const q = `
UPDATE users
SET nick_name = $2, avatar_url = $3, updated_at = now()
WHERE id = $1
RETURNING id, openid, COALESCE(unionid, ''), nick_name, avatar_url, COALESCE(phone_number, ''), used_count, remaining_count, created_at, updated_at`
	var u User
	err := s.DB.QueryRow(ctx, q, id, nickName, avatarURL).Scan(
		&u.ID, &u.OpenID, &u.UnionID, &u.NickName, &u.AvatarURL,
		&u.PhoneNumber, &u.UsedCount, &u.RemainCount, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Store) UpdateUserPhone(ctx context.Context, id int64, phone string) (User, error) {
	const q = `
UPDATE users
SET phone_number = $2, updated_at = now()
WHERE id = $1
RETURNING id, openid, COALESCE(unionid, ''), nick_name, avatar_url, COALESCE(phone_number, ''), used_count, remaining_count, created_at, updated_at`
	var u User
	err := s.DB.QueryRow(ctx, q, id, phone).Scan(
		&u.ID, &u.OpenID, &u.UnionID, &u.NickName, &u.AvatarURL,
		&u.PhoneNumber, &u.UsedCount, &u.RemainCount, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Store) ListHistoryByDevice(ctx context.Context, deviceID string, limit int) ([]HistoryItem, error) {
	const q = `
SELECT id, title, grade, COALESCE(thumb_url, ''), COALESCE(summary, ''), mode, solved_at, COALESCE(question_text, '')
FROM homework_records
WHERE device_id = $1
ORDER BY solved_at DESC
LIMIT $2`

	rows, err := s.DB.Query(ctx, q, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]HistoryItem, 0, limit)
	for rows.Next() {
		var it HistoryItem
		if err := rows.Scan(&it.ID, &it.Title, &it.Grade, &it.ThumbURL, &it.Summary, &it.Mode, &it.SolvedAt, &it.QuestionText); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (s *Store) GetHomeworkByIDAndDevice(ctx context.Context, id int64, deviceID string) (HomeworkRecord, error) {
	const q = `
SELECT id, device_id, mode, title, grade, COALESCE(thumb_url, ''), COALESCE(source_image_url, ''), COALESCE(summary, ''), COALESCE(question_text, ''), result_json, solved_at, created_at, updated_at
FROM homework_records
WHERE id = $1 AND device_id = $2`

	var rec HomeworkRecord
	err := s.DB.QueryRow(ctx, q, id, deviceID).Scan(
		&rec.ID, &rec.DeviceID, &rec.Mode, &rec.Title, &rec.Grade, &rec.ThumbURL, &rec.SourceImage,
		&rec.Summary, &rec.QuestionText, &rec.ResultJSONRaw, &rec.SolvedAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return HomeworkRecord{}, err
	}
	return rec, nil
}

func (s *Store) CreateHomework(ctx context.Context, deviceID string, mode string, imageURL string, questionText string, grade string, resultJSON any) (HomeworkRecord, error) {
	resultBytes, err := json.Marshal(resultJSON)
	if err != nil {
		return HomeworkRecord{}, err
	}
	title := buildTitle(questionText)
	summary := buildSummary(questionText)

	const q = `
INSERT INTO homework_records (device_id, mode, title, grade, thumb_url, source_image_url, summary, question_text, result_json, solved_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,now())
RETURNING id, device_id, mode, title, grade, COALESCE(thumb_url, ''), COALESCE(source_image_url, ''), COALESCE(summary, ''), COALESCE(question_text, ''), result_json, solved_at, created_at, updated_at`

	var rec HomeworkRecord
	err = s.DB.QueryRow(ctx, q, deviceID, mode, title, grade, imageURL, imageURL, summary, questionText, resultBytes).Scan(
		&rec.ID, &rec.DeviceID, &rec.Mode, &rec.Title, &rec.Grade, &rec.ThumbURL, &rec.SourceImage,
		&rec.Summary, &rec.QuestionText, &rec.ResultJSONRaw, &rec.SolvedAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return HomeworkRecord{}, err
	}
	return rec, nil
}

func (s *Store) UpdateHomeworkResult(ctx context.Context, id int64, deviceID string, mode string, questionText string, grade string, resultJSON any) (HomeworkRecord, error) {
	resultBytes, err := json.Marshal(resultJSON)
	if err != nil {
		return HomeworkRecord{}, err
	}
	title := buildTitle(questionText)
	summary := buildSummary(questionText)

	const q = `
UPDATE homework_records
SET mode=$3, title=$4, grade=$5, summary=$6, question_text=$7, result_json=$8, solved_at=now(), updated_at=now()
WHERE id = $1 AND device_id = $2
RETURNING id, device_id, mode, title, grade, COALESCE(thumb_url, ''), COALESCE(source_image_url, ''), COALESCE(summary, ''), COALESCE(question_text, ''), result_json, solved_at, created_at, updated_at`

	var rec HomeworkRecord
	err = s.DB.QueryRow(ctx, q, id, deviceID, mode, title, grade, summary, questionText, resultBytes).Scan(
		&rec.ID, &rec.DeviceID, &rec.Mode, &rec.Title, &rec.Grade, &rec.ThumbURL, &rec.SourceImage,
		&rec.Summary, &rec.QuestionText, &rec.ResultJSONRaw, &rec.SolvedAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return HomeworkRecord{}, err
	}
	return rec, nil
}

func buildTitle(questionText string) string {
	v := strings.TrimSpace(questionText)
	if v == "" {
		return "未识别题目"
	}
	if len([]rune(v)) <= 24 {
		return v
	}
	rs := []rune(v)
	return string(rs[:24]) + "..."
}

func buildSummary(questionText string) string {
	v := strings.TrimSpace(strings.ReplaceAll(questionText, "\n", " "))
	if v == "" {
		return "暂无题干摘要"
	}
	rs := []rune(v)
	if len(rs) <= 50 {
		return v
	}
	return string(rs[:50]) + "..."
}

func (s *Store) ListHistory(ctx context.Context, userID int64, limit int) ([]HistoryItem, error) {
	// Backward compatibility for old API; user_id history no longer used in mini-program flow.
	const q = `
SELECT id, title, grade, COALESCE(thumb_url, ''), COALESCE(summary, ''), mode, solved_at, COALESCE(question_text, '')
FROM homework_records
WHERE user_id = $1
ORDER BY solved_at DESC
LIMIT $2`
	rows, err := s.DB.Query(ctx, q, userID, limit)
	if err != nil {
		if strings.Contains(err.Error(), "column \"user_id\" does not exist") {
			return []HistoryItem{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	items := make([]HistoryItem, 0, limit)
	for rows.Next() {
		var it HistoryItem
		if err := rows.Scan(&it.ID, &it.Title, &it.Grade, &it.ThumbURL, &it.Summary, &it.Mode, &it.SolvedAt, &it.QuestionText); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func WrapNotFound(entity string, err error) error {
	if IsNotFound(err) {
		return fmt.Errorf("%s not found: %w", entity, err)
	}
	return err
}
