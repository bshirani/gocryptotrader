package main

import (
	"errors"
	"strings"

	"gocryptotrader/exchange/asset"
)

var (
	errInvalidPair  = errors.New("invalid currency pair supplied")
	errInvalidAsset = errors.New("invalid asset supplied")
)

func validPair(pair string) bool {
	return strings.Contains(pair, pairDelimiter)
}

func validAsset(i string) bool {
	_, err := asset.New(i)
	return err == nil
}
