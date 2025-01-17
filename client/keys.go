package client

import (
	"crypto/sha256"
	"reflect"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/thanhxeon2470/beowulf-go/config"
	"github.com/thanhxeon2470/beowulf-go/encoding/wif"
	"github.com/thanhxeon2470/beowulf-go/types"
	"golang.org/x/crypto/ripemd160"
)

var (
	//OpTypeKey include a description of the operation and the key needed to sign it
	OpTypeKey = make(map[types.OpType][]string)
)

//Keys is used as a keystroke for a specific user.
//Only a few keys can be set.
type Keys struct {
	OKey []string
}

func init() {
	OpTypeKey["transfer"] = []string{"owner"}
	OpTypeKey["transfer_to_vesting"] = []string{"owner"}
	OpTypeKey["withdraw_vesting"] = []string{"owner"}
	OpTypeKey["account_create"] = []string{"owner"}
	OpTypeKey["account_update"] = []string{"owner"}
	OpTypeKey["supernode_update"] = []string{"owner"}
	OpTypeKey["account_supernode_vote"] = []string{"owner"}
	OpTypeKey["smt_create"] = []string{"owner"}
	OpTypeKey["smart_contract"] = []string{"owner"}
	OpTypeKey["check_sidechain"] = []string{"owner"}
}

func HasElem(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

//SetKeys you can specify keys for signing transactions.
func (client *Client) SetKeys(keys *Keys) {
	client.CurrentKeys = keys
}

func (client *Client) GetPrivateKey() string {
	if client.CurrentKeys != nil {
		for _, privKey := range client.CurrentKeys.OKey {
			return privKey
		}
	}
	return ""
}

func (client *Client) GetPublicKey() string {
	if client.CurrentKeys != nil {
		for _, privKey := range client.CurrentKeys.OKey {
			return CreatePublicKey(config.ADDRESS_PREFIX, privKey)
		}
	}
	return ""
}

//SigningKeys returns the key from the CurrentKeys
func (client *Client) SigningKeys(trx types.Operation) ([][]byte, error) {
	var keys [][]byte

	if client.CurrentKeys == nil {
		return nil, errors.New("Client Keys not initialized. Use SetKeys method")
	}

	t := trx.Type()
	opKeys := OpTypeKey[t]
	for _, val := range opKeys {
		switch val {
		case "owner":
			for _, keyStr := range client.CurrentKeys.OKey {
				privKey, err := wif.Decode(keyStr)
				if err != nil {
					return nil, errors.New("error decode Owner Key: " + err.Error())
				}
				keys = append(keys, privKey)
			}
		}
	}

	return keys, nil
}

func (client *Client) GetSigningKeysOwner() ([][]byte, error) {
	var keys [][]byte

	if client.CurrentKeys == nil {
		return nil, errors.New("Client Keys not initialized. Use SetKeys method")
	}

	for _, keyStr := range client.CurrentKeys.OKey {
		privKey, err := wif.Decode(keyStr)
		if err != nil {
			return nil, errors.New("error decode Owner Key: " + err.Error())
		}
		keys = append(keys, privKey)
	}

	return keys, nil
}

//CreatePrivateKey generates a private key based on the specified parameters.
func CreatePrivateKey(user, role, password string) string {
	new_password := password + Wallet_.Salt
	hashSha256 := sha256.Sum256([]byte(user + role + new_password))
	pk := append([]byte{0x80}, hashSha256[:]...)
	chs := sha256.Sum256(pk)
	chs = sha256.Sum256(chs[:])
	b58 := append(pk, chs[:4]...)
	return base58.Encode(b58)
}

//CreatePublicKey generates a public key based on the prefix and the private key.
func CreatePublicKey(prefix, privatekey string) string {
	b58 := base58.Decode(privatekey)
	tpk := b58[:len(b58)-4]
	chs := b58[len(b58)-4:]
	nchs := sha256.Sum256(tpk)
	nchs = sha256.Sum256(nchs[:])
	if string(chs) != string(nchs[:4]) {
		return "Invalid WIF key (checksum miss-match)"
	}
	privKeyBytes := [32]byte{}
	copy(privKeyBytes[:], tpk[1:])
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes[:])
	chHash := ripemd160.New()
	chHash.Write(priv.PubKey().SerializeCompressed())
	nc := chHash.Sum(nil)
	pk := append(priv.PubKey().SerializeCompressed(), nc[:4]...)
	return prefix + base58.Encode(pk)
}
