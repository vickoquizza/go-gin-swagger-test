package usecases

import (
	"context"
	"errors"
	"go-gin-swagger-test/app/db"
)

//In memoryDB ID Counter
var PersistenceCurrentID int

// Persistence errors pool
var (
	NoAccounts = errors.New("no acounts present on the current instance of the persistence")
	NoRow      = errors.New("no acounts present on the current instance of the persistence with the given ID")
)

type Persistence []db.Account

func NewPersistence() *Persistence {
	return &Persistence{}
}

func (p *Persistence) InsertAccount(ctx context.Context, account db.AccountDTO) (db.Account, error) {
	PersistenceCurrentID++

	var newAccount db.Account
	newAccount.ID = PersistenceCurrentID
	newAccount.Name = account.Name

	(*p) = append((*p), newAccount)

	return newAccount, nil
}

func (p *Persistence) GetAllAccounts(ctx context.Context) ([]db.Account, error) {

	if len((*p)) == 0 {
		return nil, NoAccounts
	}

	return *p, nil
}

func (p *Persistence) GetAccountById(ctx context.Context, id int) (db.Account, error) {

	for _, account := range *p {
		if id == account.ID {
			return account, nil
		}
	}

	return db.Account{}, NoRow
}

func (p *Persistence) UpdateAccountById(ctx context.Context, id int, accountUpdate db.AccountDTO) error {
	for index, account := range *p {
		if id == account.ID {
			(*p)[index].Name = accountUpdate.Name
			return nil
		}
	}

	return NoRow
}

func (p *Persistence) DeleteAccountById(ctx context.Context, id int) error {
	for index, account := range *p {
		if id == account.ID {
			(*p) = append((*p)[:index], (*p)[index+1:]...)
			return nil
		}
	}

	return NoRow
}
