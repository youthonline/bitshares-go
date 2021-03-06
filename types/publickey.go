package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"fmt"
	//"gxclient-go"
	//"gxclient-go/util"
    "github.com/youthonline/bitshares-go/util"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	//sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type PublicKeys []PublicKey

type PublicKey struct {
	key      *btcec.PublicKey
	prefix   string
	checksum []byte
}

type Address struct {
	prefix   string
	data     []byte
	checksum []byte
}
func (p Address) String() string {
	b := append(p.data, p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p Address) Bytes() []byte {
	return p.data
}

func NewAddressFromString(add string) (*Address, error) {
	prefixChain := "BTS"

	prefix := add[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, errors.New("ErrAddressChainPrefixMismatch")
	}

	b58 := base58.Decode(add[len(prefixChain):])
	if len(b58) < 5 {
		return nil, errors.New("ErrInvalidAddress")
	}

	chk1 := b58[len(b58)-4:]
	data := b58[:len(b58)-4]

	chk2, err := util.Ripemd160Checksum(data)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, errors.New("ErrInvalidAddress")
	}

	a := Address{
		data:     data,
		prefix:   prefix,
		checksum: chk1,
	}

	return &a, nil
}

func (p PublicKey) String() string {
	b := append(p.Bytes(), p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var key string

	if err := ffjson.Unmarshal(data, &key); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	pub, err := NewPublicKeyFromString(key)
	if err != nil {
		return errors.Annotate(err, "NewPublicKeyFromString")
	}

	p.key = pub.key
	p.prefix = pub.prefix
	p.checksum = pub.checksum
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(p.String())
}

//func (p PublicKey) MarshalTransaction(enc *transaction.Encoder) error {
//	return enc.Encode(p.Bytes())
//}

func (p *PublicKey) ToAddress() (*Address, error) {
	return NewAddress(p)
}

func (p PublicKey) Bytes() []byte {
	return p.key.SerializeCompressed()
}

func (p PublicKey) Equal(pub *PublicKey) bool {
	return p.key.IsEqual(pub.key)
}

func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	return p.key.ToECDSA()
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func (p PublicKey) MaxSharedKeyLength() int {
	return (p.key.ToECDSA().Curve.Params().BitSize + 7) / 8
}


func NewAddress(pub *PublicKey) (*Address, error) {
	buf512 := sha512.Sum512(pub.Bytes())
	data, err := util.Ripemd160(buf512[:])
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160")
	}

	chk1, err := util.Ripemd160Checksum(data)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	ad := Address{
		prefix:   "BTS",
		data:     data,
		checksum: chk1,
	}

	return &ad, nil
}
//NewPublicKey creates a new PublicKey from string
//e.g.("GXC6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR")
func NewPublicKeyFromString(key string) (*PublicKey, error) {
	prefixChain := "BTS"

	prefix := key[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, errors.New("ErrPublicKeyChainPrefixMismatch")
	}

	b58 := base58.Decode(key[len(prefixChain):])
	if len(b58) < 5 {
		return nil, errors.New("ErrInvalidPublicKey")
	}

	chk1 := b58[len(b58)-4:]

	keyBytes := b58[:len(b58)-4]
	chk2, err := util.Ripemd160Checksum(keyBytes)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, errors.New("ErrInvalidPublicKey")
	}

	pub, err := btcec.ParsePubKey(keyBytes, btcec.S256())
	if err != nil {
		return nil, errors.Annotate(err, "ParsePubKey")
	}

	k := PublicKey{
		key:      pub,
		prefix:   prefix,
		checksum: chk1,
	}

	return &k, nil
}

func NewPublicKey(pub *btcec.PublicKey) (*PublicKey, error) {
	buf := pub.SerializeCompressed()
	chk, err := util.Ripemd160Checksum(buf)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	k := PublicKey{
		key:      pub,
		prefix:   "BTS",
		checksum: chk,
	}

	return &k, nil
}
//
//func PublicKeyComparator(key1, key2 *PublicKey) (int, error) {
//	addr1, err := key1.ToAddress()
//	if err != nil {
//		return 0, errors.Annotate(err, "ToAddress 1")
//	}
//
//	addr2, err := key2.ToAddress()
//	if err != nil {
//		return 0, errors.Annotate(err, "ToAddress 2")
//	}
//
//	return sort.StringComparator(addr1.String(), addr2.String()), nil
//}
