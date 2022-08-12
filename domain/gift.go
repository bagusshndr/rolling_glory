package domain

import (
	"context"
	"time"
)

// Article ...
type Gift struct {
	ID          int64     `json:"id"`
	NamaGift    string    `json:"nama_gift" validate:"required"`
	Deskripsi   string    `json:"deskripsi" validate:"required"`
	JumlahPoint int       `json:"jumlah_point" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Status      int       `json:"status" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// ArticleUsecase represent the article's usecases
type GiftUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Gift, string, error)
	GetByTitle(ctx context.Context, nama_gift string) (Gift, error)
	GetByID(ctx context.Context, id int64) (Gift, error)
	Store(context.Context, *Gift) error
	Update(ctx context.Context, ar *Gift) error
	Delete(ctx context.Context, id int64) error
}

// ArticleRepository represent the article's repository contract
type GiftRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Gift, nextCursor string, err error)
	GetByTitle(ctx context.Context, nama_gift string) (Gift, error)
	GetByID(ctx context.Context, id int64) (Gift, error)
	Store(ctx context.Context, a *Gift) error
	Update(ctx context.Context, ar *Gift) error
	Delete(ctx context.Context, id int64) error
}
