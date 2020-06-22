# Official Go BEOWULF Library

beowulf-go is the official Beowulf library for Go.  

## Main Functions Supported
1. CHAIN
- get_block
- get_transaction
- get_balance
2. TRANSACTION
- broadcast_transaction
- create transaction transfer
- create account
- create token

## Installation
```go
go get -u github.com/beowulf-foundation/beowulf-go
go import "github.com/beowulf-foundation/beowulf-go"
```

## Configuration
#### Init

```go
// MainNet: https://bw.beowulfchain.com/rpc
// TestNet: https://testnet-bw.beowulfchain.com/rpc
cls, _ := client.NewClient('http://localhost:8376/rpc')
defer cls.Close()
// SetKeys
key := "5Jxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Replace your private key
cls.SetKeys(&client.Keys{OKey: []string{key}})
```

## Example Usage

##### Get config
```go
fmt.Println("========== GetConfig ==========")
config, err := cls.API.GetConfig()
if err != nil {
    fmt.Println(err)
}
json_cfg, _ := json.Marshal(config)
fmt.Println(string(json_cfg))

// Use the last irreversible block number as the initial last block number.
props, err := cls.API.GetDynamicGlobalProperties()
json_props, _ := json.Marshal(props)
fmt.Println(string(json_props))
```

##### Get account
```go
account, err := cls.GetAccount("name-account")
json_acc, _ := json.Marshal(account)
fmt.Println(string(json_acc))
```

##### Get block
```go
lastBlock := props.LastIrreversibleBlockNum
block, err := cls.GetBlock(lastBlock)
json_bk, _ := json.Marshal(block)
fmt.Println(string(json_bk))
```

##### Get transaction
```go
trx, err := cls.API.GetTransaction("673fbd4609d1156bcf6d9e6c36388926f7116acc")
if err != nil {
    fmt.Println(err)
}
json_trx, _ := json.Marshal(trx)
fmt.Println(string(json_trx))
oplist := *trx.Operations
for _, op := range oplist {
    d := op.Data()
    switch d.(type){
    case *types.TransferOperation:
        byteData, _ := json.Marshal(d)
        oop := types.TransferOperation{}
        json.Unmarshal(byteData, &oop)
        fmt.Println(oop)
        fmt.Println("From:", oop.From)
        fmt.Println("To:", oop.To)
        fmt.Println("Amount:", oop.Amount)
        fmt.Println("Fee:", oop.Fee)
        fmt.Println("Memo:", oop.Memo)
    }
}
exlist := trx.Extensions
if len(exlist) > 0 {
    tmp := exlist[0]
    byteex, _ := json.Marshal(tmp)
    var met map[string]interface{}
    json.Unmarshal(byteex, &met)
    et := types.ExtensionType{}
    stype := fmt.Sprintf("%v", met["type"])
    et.Type = uint8(types.GetExtCodes(stype))
    value := met["value"].(map[string]interface{})
    ejt := types.ExtensionJsonType{}
    ejt.Data = fmt.Sprintf("%v", value["data"])
    et.Value = ejt
    fmt.Println(ejt)
    fmt.Println(et)
}
```

##### Transfer native coin
###### Transfer BWF from alice to bob
```go
resp_bwf, err := cls.Transfer("alice", "bob", "", "100.00000 BWF", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rbwf, _ := json.Marshal(resp_bwf)
fmt.Println(string(json_rbwf))
```

###### Transfer W from alice to bob
```go
resp_w, err := cls.Transfer("alice", "bob", "", "10.00000 W", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rw, _ := json.Marshal(resp_w)
fmt.Println(string(json_rw))
```

##### Transfer token
```go
//Transfer token KNOW from alice to bob
resp_tk, err := cls.Transfer("alice", "bob", "", "1000.00000 KNOW", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rtk, _ := json.Marshal(resp_tk)
fmt.Println(string(json_rtk))
```

##### Create account
###### GenKeys
```go
walletData, _ := cls.GenKeys("new-account-name")
json_wd, _ := json.Marshal(walletData)
fmt.Println(string(json_wd))
```

###### AccountCreate
```go
resp_ac, err := cls.AccountCreate("creator", walletData.Name, walletData.PublicKey,"1.00000 W")
if err != nil {
    fmt.Println(err)
}
json_rac, _ := json.Marshal(resp_ac)
fmt.Println(string(json_rac))

###### Write file wallet.
password := "your_password"
err := client.SaveWalletFile("/path/to/folder/save/wallet", "", password, walletData)
if err != nil {
    fmt.Println(err)
}
```

###### Load file wallet.
```go
rs := cls.SetKeysFromFileWallet("/path/to/folder/save/wallet/new-account-name-wallet.json", password)
if rs != nil {
    fmt.Println(rs)
}
// print keys
fmt.Println(cls.GetPrivateKey())
fmt.Println(cls.GetPublicKey())
```
