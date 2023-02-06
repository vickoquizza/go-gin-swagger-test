package repository

import (
	"context"
	"go-gin-swagger-test/app/db"
)

type Repository interface {
	InsertAccount(ctx context.Context, account db.AccountDTO) (db.Account, error)
	GetAllAccounts(ctx context.Context) ([]db.Account, error)
	GetAccountById(ctx context.Context, id int) (db.Account, error)
	UpdateAccountById(ctx context.Context, id int, accountUpdate db.AccountDTO) error
	DeleteAccountById(ctx context.Context, id int) error
}
