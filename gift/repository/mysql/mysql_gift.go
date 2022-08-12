package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/gift/repository"
	"github.com/sirupsen/logrus"
)

type mysqlGiftRepository struct {
	Conn *sql.DB
}

// NewmysqlGiftRepository will create an object that represent the article.Repository interface
func NewMysqlGiftRepository(Conn *sql.DB) domain.GiftRepository {
	return &mysqlGiftRepository{Conn}
}

func (m *mysqlGiftRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Gift, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Gift, 0)
	for rows.Next() {
		t := domain.Gift{}
		err = rows.Scan(
			&t.ID,
			&t.NamaGift,
			&t.Deskripsi,
			&t.JumlahPoint,
			&t.Stock,
			&t.Status,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlGiftRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Gift, nextCursor string, err error) {
	query := `SELECT id, nama_gift, deskripsi, jumlah_point, stock, status, created_at, updated_at FROM tbl_gift WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *mysqlGiftRepository) Store(ctx context.Context, a *domain.Gift) (err error) {
	query := `INSERT  tbl_gift SET nama_gift=? , deskripsi=? , jumlah_point=?, stock=?, status=?, created_at=? , updated_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.NamaGift, a.Deskripsi, a.JumlahPoint, a.Stock, a.Status, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *mysqlGiftRepository) GetByTitle(ctx context.Context, nama_gift string) (res domain.Gift, err error) {
	query := `SELECT id, nama_gift, deskripsi, jumlah_point, stock, status, created_at, updated_at
  						FROM tbl_gift WHERE nama_gift = ?`

	list, err := m.fetch(ctx, query, nama_gift)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlGiftRepository) GetByID(ctx context.Context, id int64) (res domain.Gift, err error) {
	query := `SELECT id, nama_gift, deskripsi, jumlah_point, stock, status, created_at, updated_at
				FROM tbl_gift WHERE id = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Gift{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlGiftRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM tbl_gift WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *mysqlGiftRepository) Update(ctx context.Context, ar *domain.Gift) (err error) {
	query := `UPDATE tbl_gift set nama_gift=? , deskripsi=? , jumlah_point=?, stock=?, status=?, created_at=? , updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.NamaGift, ar.Deskripsi, ar.JumlahPoint, ar.Stock, ar.Status, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
