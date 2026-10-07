package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return db, mock
}

func TestUserRepositoryFindByUsername(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("alice", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "nickname", "role"}).
			AddRow(1, "alice", "a@b.c", "Alice", "user"))
	u, err := repo.FindByUsername("alice")
	if err != nil {
		t.Fatalf("FindByUsername: %v", err)
	}
	if u.ID != 1 || u.Role != "user" {
		t.Errorf("unexpected user: %+v", u)
	}
}

func TestUserRepositoryFindByUsernameNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("nobody", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	if _, err := repo.FindByUsername("nobody"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepositoryCreateDuplicate(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WillReturnError(errors.New("Duplicate entry 'alice' for key 'users.username'"))
	mock.ExpectRollback()
	u := &model.User{Username: "alice", Email: "a@b.c"}
	if err := repo.Create(u); !errors.Is(err, ErrDuplicate) {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}
