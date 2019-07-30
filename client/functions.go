package client

import (
	"beowulf-go/api"
	"beowulf-go/types"
	"errors"
)

const fdt = `"20060102t150405"`
const fee = "0.01000 W"
const smtCreationFee = "1000.00000 W"
const accCreationFee = "1.00000 W"

func (client *Client) GetBlock(blockNum uint32) (*api.Block, error) {
	return client.API.GetBlock(blockNum)
}

func (client *Client) GetTransaction(trx string) (*api.TransactionResponse, error) {
	return client.API.GetTransaction(trx)
}

func (client *Client) GetAccount(account string) (*api.AccountInfo, error) {
	accounts, err := client.API.GetAccounts(account)
	if err != nil {
		return nil, err
	}
	var list []api.AccountInfo
	list = *accounts
	if len(list) == 0 {
		return nil, errors.New("Unknown account")
	}
	return &list[0], nil
}

func (client *Client) GetSupernodeByAccount(account string) (*api.SupernodeInfo, error) {
	return client.API.GetSupernodeByAccount(account)
}

func (client *Client) GetSupernodeVoted(account string) (*api.SupernodeVoteList, error) {
	return client.API.GetSupernodeVoted(account)
}

func (client *Client) GetKeyReferences(publicKey string) (*[][]string, error) {
	return client.API.GetKeyReferences(publicKey)
}

func (client *Client) ListAccounts(lowerBound string, limit uint32) (*[]string, error) {
	return client.API.ListAccounts(lowerBound, limit)
}

func (client *Client) GetActiveSupernodes() (*[]string, error) {
	return client.API.GetActiveSupernodes()
}

func (client *Client) ListSupernodes(lowerBound string, limit uint32) (*[]string, error) {
	return client.API.ListSupernodes(lowerBound, limit)
}

func (client *Client) ListTokens() (*api.TokenList, error) {
	return client.API.ListTokens()
}

func (client *Client) GetToken(name string) (*api.TokenInfo, error) {
	tokens, err := client.API.GetTokens(name)
	if err != nil {
		return nil, err
	}
	var list []api.TokenInfo
	list = *tokens
	if len(list) == 0 {
		return nil, errors.New("Unknown token")
	}
	return &list[0], nil
}

func (client *Client) GetBalance(account, tokenName string, decimals uint8) (*string, error) {
	return client.API.GetBalance(account, tokenName, decimals)
}

//Transfer of funds to any user.
func (client *Client) Transfer(fromName, toName, memo, amount, fee string) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.TransferOperation{
		From:   fromName,
		To:     toName,
		Amount: amount,
		Fee:    fee,
		Memo:   memo,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "Transfer", Bresp: resp}, err
}

func (client *Client) MultiOp(trx []types.Operation) (*OperResp, error) {
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "Multi", Bresp: resp}, err
}

func (client *Client) CreateToken(creator, controlAcc, tokenName string, decimals uint8, maxSuplly uint64) (*OperResp, error) {
	//config, err := client.API.GetConfig()
	//if err != nil{
	//	return nil, err
	//}
	//feeAmt := config.TokenCreationFee

	var trx []types.Operation
	tx := &types.SmtCreateOperation{
		ControlAccount: controlAcc,
		Symbol:         &types.AssetSymbol{Decimals: decimals, AssetName: tokenName},
		Creator:        creator,
		SmtCreationFee: smtCreationFee,
		Precision:      decimals,
		Extensions:     [][]interface{}{},
		MaxSupply:      maxSuplly,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "SmtCreate", Bresp: resp}, err
}

//AccountSupernodeVote of voting for the delegate.
func (client *Client) AccountSupernodeVote(username, witnessName, fee string, approv bool, votes int64) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.AccountSupernodeVoteOperation{
		Account:   username,
		Supernode: witnessName,
		Approve:   approv,
		Votes:     votes,
		Fee:       fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "AccountSupernodeVote", Bresp: resp}, err
}

//TransferToVesting transfer to POWER
func (client *Client) TransferToVesting(from, to, amount, fee string) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.TransferToVestingOperation{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "TransferToVesting", Bresp: resp}, err
}

//WithdrawVesting down POWER
func (client *Client) WithdrawVesting(account, vshares, fee string) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.WithdrawVestingOperation{
		Account:       account,
		VestingShares: vshares,
		Fee:           fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "WithdrawVesting", Bresp: resp}, err
}

//SupernodeUpdate updating delegate data
func (client *Client) SupernodeUpdate(owner, blocksigningkey, fee string) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.SupernodeUpdateOperation{
		Owner:           owner,
		BlockSigningKey: blocksigningkey,
		Fee:             fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "SupernodeUpdate", Bresp: resp}, err
}

//AccountCreate creating a user in systems
func (client *Client) GenKeys(newAccountName string) (*WalletData, error) {
	role := "owner"
	password := randStringBytes(16)
	priv := CreatePrivateKey(newAccountName, role, password)
	pub := CreatePublicKey(ADDRESS_PREFIX, priv)

	return &WalletData{Name: newAccountName, PrivateKey: priv, PublicKey: pub}, nil
}

func (client *Client) AccountCreate(creator, newAccountName, publicKey, fee string) (*OperResp, error) {
	var trx []types.Operation
	empty := map[string]int64{}

	owner := types.Authority{
		WeightThreshold: 1,
		AccountAuths:    empty,
		KeyAuths:        map[string]int64{publicKey: 1},
	}

	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "AccountCreate", Bresp: resp}, err
}

func (client *Client) AccountCreateWS(creator, newAccountName, password, fee string) (*OperResp, error) {
	type Keys struct {
		Private string
		Public  string
	}

	var trx []types.Operation
	var listKeys = make(map[string]Keys)
	empty := map[string]int64{}
	roles := [1]string{"owner"}

	for _, val := range roles {
		priv := CreatePrivateKey(newAccountName, val, password)
		pub := CreatePublicKey(ADDRESS_PREFIX, priv)
		listKeys[val] = Keys{Private: priv, Public: pub}
	}

	owner := types.Authority{
		WeightThreshold: 1,
		AccountAuths:    empty,
		KeyAuths:        map[string]int64{listKeys["owner"].Public: 1},
	}

	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "AccountCreateWS", Bresp: resp}, err
}

//CreateMultiSigAccount creating an account shared among many users in systems
func (client *Client) CreateMultiSigAccount(creator, newAccountName, fee string, owners []string) (*OperResp, error) {
	var trx []types.Operation
	var listKeys = make(map[string]int64)
	empty := map[string]int64{}
	for _, k := range owners {
		listKeys[k] = 1
	}

	owner := types.Authority{
		WeightThreshold: uint32(len(owners)),
		AccountAuths:    empty,
		KeyAuths:        listKeys,
	}
	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "AccountCreate", Bresp: resp}, err
}

//AccountUpdate update owner keys for account
//TODO: every key has different weight on account
func (client *Client) AccountUpdate(account, fee string, owners []string) (*OperResp, error) {
	var trx []types.Operation
	var listKeys = make(map[string]int64)
	empty := map[string]int64{}
	for _, k := range owners {
		listKeys[k] = 1
	}

	owner := types.Authority{
		WeightThreshold: uint32(len(owners)),
		AccountAuths:    empty,
		KeyAuths:        listKeys,
	}
	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountUpdateOperation{
		Account:      account,
		Owner:        &owner,
		JSONMetadata: jsonMeta,
		Fee:          fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx)
	return &OperResp{NameOper: "AccountUpdate", Bresp: resp}, err
}
