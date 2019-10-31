package client

import (
	"beowulf-go/api"
	"beowulf-go/config"
	"beowulf-go/transactions"
	"beowulf-go/types"
	"time"
)

var RefBlockMap = make(map[time.Time]uint32)

//BResp of response when sending a transaction.
type BResp struct {
	ID string
	//BlockNum int32
	//TrxNum   int32
	//Expired  bool
	//CreatedTime int64
	JSONTrx string
}

//OperResp type is returned when the operation is performed.
type OperResp struct {
	NameOper string
	Bresp    *BResp
}

//Get HeadBlockNumber from mem before getting from Blockchain
func (client *Client) GetHeadBlockNum() (uint32, error) {
	if len(RefBlockMap) > 0 {
		for k := range RefBlockMap {
			old := k.Add(config.GET_HEAD_BLOCK_NUM_TIMEOUT_IN_MIN * time.Minute)
			now := time.Now().UTC()
			if old.Before(now) {
				delete(RefBlockMap, k)
				props, err := client.API.GetDynamicGlobalProperties()
				if err != nil {
					return 0, err
				}
				refBlockNum := props.HeadBlockNumber
				if refBlockNum > config.HEAD_BLOCK_NUM_SPAN {
					refBlockNum -= config.HEAD_BLOCK_NUM_SPAN
				}
				RefBlockMap[now] = refBlockNum
				return refBlockNum, nil
			}
			return RefBlockMap[k], nil
		}
	}
	props, err := client.API.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	refBlockNum := props.HeadBlockNumber
	if refBlockNum > config.HEAD_BLOCK_NUM_SPAN {
		refBlockNum -= config.HEAD_BLOCK_NUM_SPAN
	}
	now := time.Now().UTC()
	RefBlockMap[now] = refBlockNum
	return refBlockNum, nil
}

//SendTrx generates and sends an array of transactions to BEOWULF.
func (client *Client) SendTrx(strx []types.Operation, extension string) (*BResp, error) {
	var bresp BResp

	// Getting the necessary parameters
	refBlockNum, err := client.GetHeadBlockNum()
	if err != nil {
		return nil, err
	}
	block, err := client.API.GetBlock(refBlockNum)
	if err != nil {
		return nil, err
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return nil, err
	}

	ex := make([]interface{}, 1)
	as := types.ExtensionJsonType{extension}
	tas := types.ExtensionType{uint8(types.ExtJsonType.Code()), as}
	ex[0] = &tas

	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(refBlockNum),
		RefBlockPrefix: refBlockPrefix,
		Extensions:     ex, //[]interface{}{},
	})

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	//var br BResp
	//br.ID = "1"
	//br.JSONTrx = "{\"name\":\"thu\"}"
	////t := []string{"{\"name\":\"thu\"}"}
	//var t []BResp
	//var t1 testExt
	//t1.tp = "string"
	//t1.value = "test extension"

	//var trx []*types.TAssetSymbol
	//var tas types.TAssetSymbol
	//var t1 types.AssetSymbol
	//t1.Decimals = 5
	//t1.AssetName = "ABC"
	//tas.Type = "asset_symbol_type"
	//tas.Value = t1
	//t := &tas
	//
	//trx = append(trx, t)
	//s := make([]interface{}, len(trx))
	//for i, v := range trx {
	//	s[i] = v
	//}
	//tx.Extensions = s //append(tx.Extensions, s)

	// Obtain the key required for signing
	privKeys, err := client.SigningKeys(strx[0])
	if err != nil {
		return nil, err
	}

	// Sign the transaction
	tx.Transaction.Signatures = []string{}
	txId, err := tx.Sign(privKeys, client.chainID)
	if err != nil || txId == "" {
		return nil, err
	}

	// Sending a transaction
	//var resp *api.AsyncBroadcastResponse
	//resp, err = client.API.BroadcastTransaction(tx.Transaction)
	var errb error
	if client.AsyncProtocol {
		var resp *api.AsyncBroadcastResponse
		resp, errb = client.API.BroadcastTransaction(tx.Transaction)
		if resp != nil {
			//if txId != resp.ID {
			//	return nil, errors.New("TransactionID is not mapped")
			//}
			bresp.ID = resp.ID
		}
	} else {
		var resp *api.BroadcastResponse
		resp, errb = client.API.BroadcastTransactionSynchronous(tx.Transaction)
		if resp != nil {
			//if txId != resp.ID {
			//	return nil, errors.New("TransactionID is not mapped")
			//}
			bresp.ID = resp.ID
		}
	}
	if errb != nil {
		err = errb
	}

	bresp.JSONTrx, _ = JSONTrxString(tx)

	if err != nil {
		return &bresp, err
	}

	return &bresp, nil
}

func (client *Client) GetTrx(strx []types.Operation, extension string) (*types.Transaction, error) {
	// Getting the necessary parameters
	refBlockNum, err := client.GetHeadBlockNum()
	if err != nil {
		return nil, err
	}
	block, err := client.API.GetBlock(refBlockNum)
	if err != nil {
		return nil, err
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return nil, err
	}
	ex := make([]interface{}, 1)
	as := types.ExtensionJsonType{extension}
	tas := types.ExtensionType{uint8(types.ExtJsonType.Code()), as}
	ex[0] = &tas
	tx := &types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(refBlockNum),
		RefBlockPrefix: refBlockPrefix,
		Extensions:     ex, //[]interface{}{},
	}

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
	//expTime := time.Now().Add(59 * time.Minute).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	return tx, nil
}
