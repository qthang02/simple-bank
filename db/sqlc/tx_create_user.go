package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(cxt context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(cxt, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(cxt, arg.CreateUserParams)
		if err != nil {
			return err
		}
		
		return arg.AfterCreate(result.User)
	})

	return result, err
}
	