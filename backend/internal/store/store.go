package store

import (
	"context"
	"errors"
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
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Grade    string    `json:"grade"`
	ThumbURL string    `json:"thumbUrl"`
	SolvedAt time.Time `json:"solvedAt"`
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

	if err := s.seedHistoryIfEmpty(ctx, u.ID); err != nil {
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

func (s *Store) ListHistory(ctx context.Context, userID int64, limit int) ([]HistoryItem, error) {
	const q = `
SELECT id, title, grade, COALESCE(thumb_url, ''), solved_at
FROM homework_records
WHERE user_id = $1
ORDER BY solved_at DESC
LIMIT $2`

	rows, err := s.DB.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]HistoryItem, 0, limit)
	for rows.Next() {
		var it HistoryItem
		if err := rows.Scan(&it.ID, &it.Title, &it.Grade, &it.ThumbURL, &it.SolvedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (s *Store) seedHistoryIfEmpty(ctx context.Context, userID int64) error {
	const countQ = `SELECT COUNT(1) FROM homework_records WHERE user_id = $1`
	var n int
	if err := s.DB.QueryRow(ctx, countQ, userID).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	const insertQ = `
INSERT INTO homework_records (user_id, title, grade, thumb_url, solved_at)
VALUES
($1, '24 x 15 = ?', '三年级', '/images/generated-1771139016204.png', now() - interval '15 minutes'),
($1, '阅读理解：小蝌蚪找妈妈', '四年级', '/images/generated-1771138856711.png', now() - interval '90 minutes'),
($1, '长方形面积计算', '三年级', '/images/generated-1771138893602.png', now() - interval '1 day')`
	_, err := s.DB.Exec(ctx, insertQ, userID)
	return err
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
