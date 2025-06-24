package helpers

import "time"

type UnbondingDelegationEntry struct {
	CreationHeight string `json:"creation_height"`
	CompletionTime string `json:"completion_time"`
	InitialBalance string `json:"initial_balance"`
	Balance        string `json:"balance"`
	UnbondingID    string `json:"unbonding_id"`
}

type Unbond struct {
	DelegatorAddress string                     `json:"delegator_address"`
	ValidatorAddress string                     `json:"validator_address"`
	Entries          []UnbondingDelegationEntry `json:"entries"`
}

type UnbondingDelegationResponse struct {
	Unbond Unbond `json:"unbond"`
}

type DelegationQueryResponse struct {
	DelegationResponse struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_response"`
}

// From stakingtypes.Validator
type Vals struct {
	Validators []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation       string `json:"min_self_delegation"`
		UnbondingOnHoldRefCount string `json:"unbonding_on_hold_ref_count"`
		UnbondingIds            []any  `json:"unbonding_ids"`
	} `json:"validators"`
	Pagination struct {
		NextKey any    `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
