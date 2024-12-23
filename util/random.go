package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ownerLength          = 6
	currencyLength       = 4
	moneyLimit     int64 = 1000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandStr(length int) string {
	var sb strings.Builder
	for length > 0 {
		length--
		sb.WriteByte(alphabet[RandInt(0, int64(len(alphabet)))])
	}
	return sb.String()
}

func RandOwner() string {
	return RandStr(ownerLength)
}

func RandMoney() int64 {
	return RandInt(0, moneyLimit)
}

func RandCurrency() string {
	return RandStr(currencyLength)
}
