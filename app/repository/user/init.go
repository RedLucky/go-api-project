package user

import (
	"context"
	"database/sql"
	"fmt"
	"go-api-project/domain"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Mysql *sql.DB
}

// NewUserRepository will create an object that represent the user.Repository interface
func NewUserRepository(Conn *sql.DB) domain.UserRepository {
	return &UserRepository{Conn}
}

func (r *UserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := r.Mysql.QueryContext(ctx, query, args...)
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

	result = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Email,
			&t.Name,
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

func (m *UserRepository) Fetch(ctx context.Context) (res []domain.User, err error) {
	query := `SELECT id,username,email, name, updated_at, created_at
  						FROM users ORDER BY created_at LIMIT 10 `

	// decodedCursor, err := repository.DecodeCursor(cursor)
	// if err != nil && cursor != "" {
	// 	return nil, "", domain.ErrBadParamInput
	// }

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	// if len(res) == int(num) {
	// 	nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	// }

	return
}
func (m *UserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
	query := `SELECT id,email,username, name, updated_at, created_at
  						FROM users WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT id,email,username, name, updated_at, created_at
  						FROM users WHERE email = ?`

	list, err := m.fetch(ctx, query, email)
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

func (m *UserRepository) GetByUsername(ctx context.Context, username string) (res domain.User, err error) {
	query := `SELECT id,email,username, name, updated_at, created_at
  						FROM users WHERE username = ?`

	list, err := m.fetch(ctx, query, username)
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

func (m *UserRepository) Store(ctx context.Context, a *domain.User) (err error) {
	query := `INSERT  users SET username=? , email=? , name=?, updated_at=? , created_at=?`
	stmt, err := m.Mysql.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Username, a.Email, a.Name, a.UpdatedAt, a.CreatedAt)
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

func (m *UserRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM users WHERE id = ?"

	stmt, err := m.Mysql.PrepareContext(ctx, query)
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
func (m *UserRepository) Update(ctx context.Context, ar *domain.User) (err error) {
	query := `UPDATE users set username=?, email=?, name=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Mysql.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Username, ar.Email, ar.Name, ar.UpdatedAt, ar.ID)
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
