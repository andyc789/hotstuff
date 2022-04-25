package wendy

import (
	//"github.com/relab/hotstuff"
	"github.com/relab/hotstuff/consensus"
	//"github.com/relab/hotstuff/crypto"
	"github.com/relab/hotstuff/crypto/bls12"
	"github.com/relab/hotstuff/modules"
	"github.com/herumi/bls-eth-go-binary/bls"
	"strconv"
	"encoding/binary"
)

func init() {
	modules.RegisterModule("ed25519", func() consensus.CryptoImpl {
		return New()
	})
}

const (
	
)

type KeyStor struct {
	sbit [][]bls.SecretKey
	//sbit1 []bls.SecretKey
	pbit [][]bls.PublicKey
	//pbit1 []bls.PublicKey
}



func KGen(len int) KeyStor {
	s := make([][]bls.SecretKey, len)
	//s1 := make([]bls.SecretKey, len)
	p := make([][]bls.PublicKey, len)
	//p1 := make([]bls.PublicKey, len)

//	bls.

	for x := 0; x < len; x++ {
		/*s0[x], _ = bls.()
		p0[x] = s0[x].()
		s1[x], _ = bls12.GeneratePrivateKey()
		p1[x] = s1[x].Public()*/
		var a bls.SecretKey
		var b bls.SecretKey

		s[0][x] = a.SetByCSPRNG()
		s[1][x] = b.SetByCSPRNG()
		p[0][x] = *a.GetPublicKey()
		p[1][x] = *b.GetPublicKey()
	}

	return KeyStor{
		/*sbit0: s0,
		sbit1: s1,
		pbit0: p0,
		pbit1: p1,*/

		sbit: s,
		pbit: p,
	}
} //PUBLIC KEY RELATIONSHIP WITH PI?

func SignShare(kStor KeyStor, vDiff int, v int) *bls.Sign {
	vDiff64 := int64(vDiff)
	binaryString := strconv.FormatInt(vDiff64, 2)
	//binaryNum, _ := strconv.ParseInt(binaryString, 10, 32)
	var sKey bls.SecretKey
	var vByte []byte //use make function to initialize array, size 8
	binary.BigEndian.PutUint64(vByte[0:8], uint64(v))

	for index, bit := range binaryString {
		if (string(bit)=="1") {
			sKey.Add(&kStor.sbit[1][index])
		} else if (string(bit)=="0") {
			sKey.Add(&kStor.sbit[0][index])
		}
	}
	//FINISH
	return sKey.SignByte(vByte)
}

//FOR REFERENCE
func toByteArray(i int32) (arr [4]byte) {
    binary.BigEndian.PutUint32(arr[0:4], uint32(i))
    return
}

func VerifyShare(kStor KeyStor, vDiff int, v int, sig bls.Sign) bool {
	//call verifybyte on signature
	vDiff64 := int64(vDiff)
	binaryString := strconv.FormatInt(vDiff64, 2)
	//binaryNum, _ := strconv.ParseInt(binaryString, 10, 32)
	var pKey bls.PublicKey
	var vByte []byte //use make function to initialize array, size 8
	binary.BigEndian.PutUint64(vByte[0:8], uint64(v))

	for index, bit := range binaryString {
		if (string(bit)=="1") {
			pKey.Add(&kStor.pbit[1][index])
		} else if (string(bit)=="0") {
			pKey.Add(&kStor.pbit[0][index])
		}
	}

	return sig.VerifyByte(&pKey, vByte)
}

func Agg(sigshares []bls.Sign) bls.Sign {
	var rV bls.Sign
	rV.Aggregate(sigshares) //USE THIS FUNCTION
	return rV
}

func VerifyAgg(kStor KeyStor, vDiffs []int, v int, sig bls.Sign) {
	//fastaggregateverify
	//agg sigs combine sigs on different messages into 1 sig (WHAT IS SPECIFIED HERE)
	//multi sigs combine sigs on same message into 1 sig
	//READ SECTION A, B OF SECTION IV IN PAPER FOR INTUITION, DONT PAY ATTN TO MATH FORMULA AS MUCH
	var pubkeys []bls.PublicKey
	var vByte []byte
	binary.BigEndian.PutUint64(vByte[0:8], uint64(v))
	//ITERATE THROUGH
	for i, vDiff := range vDiffs {
		vDiff64 := int64(vDiff)
		binaryString := strconv.FormatInt(vDiff64, 2)

		var pKey bls.PublicKey
		for index, bit := range binaryString {
			if (string(bit)=="1") {
				pKey.Add(&kStor.pbit[1][index])
			} else if (string(bit)=="0") {
				pKey.Add(&kStor.pbit[0][index])
			}
		}

		pubkeys[i] = pKey
	}


	sig.FastAggregateVerify(pubkeys, vByte)
}




/*
func New() consensus.CryptoImpl {
	initKeys()

}

func initKeys() {
	
}*/
