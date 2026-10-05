package model

type TransferRequest struct {
    SourceAccountID string 'json: "source_account_id"'
    TargetAccountID string 'json: "target_account_id"'
    AmountCets      int64  'json: "amount_cents"'
    Currency        string 'json: "currency"'

}



