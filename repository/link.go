package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"watchmen/config"
)

type LinkRepository interface {
	LinkExists(ctx context.Context, link string) (bool, error)
	CreateLink(ctx context.Context, link *Link) error
	GetLink(ctx context.Context, userID uint) ([]*Link, error)
	GetAllLinks(ctx context.Context) ([]*Link, error)
	RetrieveLinkData(ctx context.Context, userID, linkID uint) (*LinkReportStatus, error)
	CreateLinkReport(ctx context.Context, linkReport *LinkReport) error
}

// Link is a model in base API database
type Link struct {
	ID             uint      `db:"id"`
	UserID         uint      `db:"user_id"`
	ErrorThreshold uint      `db:"error_threshold"`
	URL            string    `db:"url"`
	Method         string    `db:"method"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// LinkReport is a model in base API database
type LinkReport struct {
	ID        uint      `db:"id"`
	LinkID    uint      `db:"link_id"`
	Status    bool      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

type LinkReportStatus struct {
	OK  int `db:"ok"`
	ERR int `db:"error"`
}

type linkRepository struct {
	baseAPIDB *sqlx.DB
}

func NewLinkRepository(baseAPIDB *sqlx.DB) LinkRepository {
	repo := new(linkRepository)
	repo.baseAPIDB = baseAPIDB

	return repo
}

func (r *linkRepository) LinkExists(ctx context.Context, link string) (bool, error) {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	var exists bool
	err := r.baseAPIDB.GetContext(ctx, &exists, "SELECT EXISTS (Select id FROM links WHERE cellphone = ?) AS `exists`;", link)
	if err != nil {
		return true, err
	}

	return exists, nil
}

func (r *linkRepository) CreateLink(ctx context.Context, link *Link) error {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.WriteTimeout)
	defer done()

	result, err := r.baseAPIDB.NamedExecContext(ctx, "INSERT INTO links(url, `method`, user_id, error_threshold, created_at, updated_at) VALUES (:url, :method, :user_id, :error_threshold, NOW(), NOW())",
		link)
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

	link.ID = uint(id)

	return nil
}

func (r *linkRepository) GetLink(ctx context.Context, userID uint) ([]*Link, error) {
	var links []*Link
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	err := r.baseAPIDB.SelectContext(ctx, &links, "SELECT * FROM links WHERE user_id = ?", userID)
	if err != nil {
		return links, err
	}

	return links, nil
}

func (r *linkRepository) GetAllLinks(ctx context.Context) ([]*Link, error) {
	var links []*Link
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	err := r.baseAPIDB.SelectContext(ctx, &links, "SELECT * FROM links")
	if err != nil {
		return links, err
	}

	return links, nil
}

func (r *linkRepository) RetrieveLinkData(ctx context.Context, userID, linkID uint) (*LinkReportStatus, error) {
	status := new(LinkReportStatus)
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.ReadTimeout)
	defer done()

	err := r.baseAPIDB.GetContext(ctx, status, "SELECT COUNT(case lr.status when TRUE then 1 else null end) AS ok, "+
		"COUNT(case lr.status when FALSE then 1 else null end) AS error FROM link_report lr "+
		"JOIN links l ON lr.link_id = l.id "+
		"WHERE l.user_id = ? "+
		"AND l.id = ? AND lr.created_at = CURDATE()", userID, linkID)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (r *linkRepository) CreateLinkReport(ctx context.Context, lr *LinkReport) error {
	ctx, done := context.WithTimeout(ctx, config.C.BaseAPIDatabase.WriteTimeout)
	defer done()

	result, err := r.baseAPIDB.NamedExecContext(ctx, "INSERT INTO link_report(link_id, status, created_at) VALUES (:link_id, :status, NOW())",
		lr)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
