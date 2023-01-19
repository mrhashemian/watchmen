package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"watchmen/config"
)

type UserRepository interface {
	CellphoneExists(ctx context.Context, cellphone string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	GetAlerts(ctx context.Context, userID uint) ([]*Link, error)
}

// User is a model in base API database
type User struct {
	ID        uint      `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FullName  string    `db:"fullname"`
	Cellphone string    `db:"cellphone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), config.C.User.PasswordHashCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)

	return nil
}

func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err
}

type userRepository struct {
	baseAPIDB *sqlx.DB
}

func NewUserRepository(baseAPIDB *sqlx.DB) UserRepository {
	repo := new(userRepository)
	repo.baseAPIDB = baseAPIDB

	return repo
}

func (r *userRepository) CellphoneExists(ctx context.Context, cellphone string) (bool, error) {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	var exists bool
	err := r.baseAPIDB.GetContext(ctx, &exists, "SELECT EXISTS (Select id FROM users WHERE cellphone = ?) AS `exists`;", cellphone)
	if err != nil {
		return true, err
	}

	return exists, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	var exists bool
	err := r.baseAPIDB.GetContext(ctx, &exists, "SELECT EXISTS (Select id FROM users WHERE email = ?) AS `exists`;", email)
	if err != nil {
		return true, err
	}

	return exists, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.WriteTimeout)
	defer done()

	result, err := r.baseAPIDB.NamedExecContext(ctx, "INSERT INTO users(email, password, fullname, cellphone, created_at, updated_at) VALUES (:email, :password, :fullname, :cellphone, NOW(), NOW())",
		user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("no rows affected")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	p := new(User)
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	err := r.baseAPIDB.GetContext(ctx, p, "SELECT * FROM users WHERE email = ?;", email)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *userRepository) GetAlerts(ctx context.Context, userID uint) ([]*Link, error) {
	var p []*Link
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	err := r.baseAPIDB.SelectContext(ctx, &p, "WITH cte AS (SELECT l.id id, count(*) c FROM links l "+
		"JOIN link_report lr ON lr.link_id = l.id "+
		"WHERE l.user_id = ? "+
		"AND lr.status = FALSE "+
		"GROUP BY l.id) "+
		"SELECT l.* FROM links l "+
		"JOIN cte ON l.id = cte.id "+
		"AND cte.c >= l.error_threshold", userID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
