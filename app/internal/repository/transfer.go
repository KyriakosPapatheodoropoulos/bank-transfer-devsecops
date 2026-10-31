package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

func ExecuteTransfer(
	ctx context.Context,
	pool *pgxpool.Pool,
	sourceAccountID string,
	targetAccountID string,
	amountCents int64,
	currency string,
) error {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sourceBalance int64

	err = tx.QueryRow(
		ctx,
		`
		SELECT balance_cents
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`,
		sourceAccountID,
	).Scan(&sourceBalance)

	if err != nil {
		return err
	}

	var targetCurrency string

	err = tx.QueryRow(
		ctx,
		`
		SELECT currency
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`,
		targetAccountID,
	).Scan(&targetCurrency)

	if err != nil {
		return err
	}

	if sourceBalance < amountCents {
		return ErrInsufficientFunds
	}

	if targetCurrency != currency {
		return errors.New("currency mismatch")
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE accounts
		SET balance_cents = balance_cents - $1
		WHERE id = $2
		`,
		amountCents,
		sourceAccountID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE accounts
		SET balance_cents = balance_cents + $1
		WHERE id = $2
		`,
		amountCents,
		targetAccountID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO transactions (
			id,
			source_account_id,
			target_account_id,
			amount_cents,
			currency
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4
		)
		`,
		sourceAccountID,
		targetAccountID,
		amountCents,
		currency,
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
