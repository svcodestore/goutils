package rand

import (
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"
	"time"

	"github.com/thanhpk/randstr"
)

func RandRange(min, max int) (n int) {
	rand.Seed(time.Now().Unix())
	n = rand.Intn(max-min) + min
	return
}

func RandomInt(max int64) (n int64) {
	r, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(max))

	return r.Int64()
}

func GenerateClientId() string {
	return randstr.Hex(10)
}

func GenerateClientSecret() string {
	return randstr.Hex(20)
}
