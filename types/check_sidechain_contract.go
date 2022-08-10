package types

import (
	"github.com/thanhxeon2470/beowulf-go/encoding/transaction"
)

//SmartContractOperation represents transfer operation data.
type CheckSidechainOperation struct {
	Committer   string `json:"committer"`
	Csid        string `json:"csid"`
	CsOperation string `json:"cs_operation"`
	Fee         string `json:"fee"`
}

//Type function that defines the type of operation SmartContractOperation.
func (op *CheckSidechainOperation) Type() OpType {
	return TypeCheckSidechain
}

//Data returns the operation data SmartContractOperation.
func (op *CheckSidechainOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type SmtCreateOperation to bytes.
func (op *CheckSidechainOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeCheckSidechain.Code()))
	//enc.Encode(op.RequiredOwners)
	// encode AccountAuths as map[string]uint16
	enc.EncodeString(op.Committer)
	enc.Encode(op.Csid)
	enc.Encode(op.CsOperation)
	enc.EncodeMoney(op.Fee)
	//enc.Encode(op.Extensions)
	//enc.EncodeUVarint(0)
	return enc.Err()
}
