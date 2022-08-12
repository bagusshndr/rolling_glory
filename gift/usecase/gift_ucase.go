package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type giftUsecase struct {
	giftRepo       domain.GiftRepository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewGiftUsecase(a domain.GiftRepository, timeout time.Duration) domain.GiftUsecase {
	return &giftUsecase{
		giftRepo:       a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (a *giftUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Gift, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.giftRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (a *giftUsecase) GetByTitle(c context.Context, title string) (res domain.Gift, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.giftRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	return
}

func (a *giftUsecase) GetByID(c context.Context, id int64) (res domain.Gift, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.giftRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *giftUsecase) Store(c context.Context, m *domain.Gift) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedGift, _ := a.GetByTitle(ctx, m.NamaGift)
	if existedGift != (domain.Gift{}) {
		return domain.ErrConflict
	}

	err = a.giftRepo.Store(ctx, m)
	return
}

func (a *giftUsecase) Update(c context.Context, ar *domain.Gift) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.giftRepo.Update(ctx, ar)
}

func (a *giftUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.giftRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Gift{}) {
		return domain.ErrNotFound
	}
	return a.giftRepo.Delete(ctx, id)
}
