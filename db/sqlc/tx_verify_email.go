package db

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	EmailId int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(cxt context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(cxt, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(cxt, UpdateVerifyEmailParams{
			ID: arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(cxt, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool: true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
	