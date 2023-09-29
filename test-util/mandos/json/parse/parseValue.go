package mandosjsonparse

import (
	"errors"
	"math/big"

	twos "github.com/kalyan3104/dme-components-big-int/twos-complement"
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	oj "github.com/kalyan3104/dme-vm-util/test-util/orderedjson"
)

type bigIntParseFormat int

const (
	bigIntSignedBytes bigIntParseFormat = iota
	bigIntUnsignedBytes
)

func (p *Parser) processCheckBigInt(obj oj.OJsonObject, format bigIntParseFormat) (mj.JSONCheckBigInt, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckBigInt{
			Value:    nil,
			IsStar:   true,
			Original: "*"}, nil
	}

	jbi, err := p.processBigInt(obj, format)
	if err != nil {
		return mj.JSONCheckBigInt{}, err
	}
	return mj.JSONCheckBigInt{
		Value:    jbi.Value,
		IsStar:   false,
		Original: jbi.Original,
	}, nil
}

func (p *Parser) processBigInt(obj oj.OJsonObject, format bigIntParseFormat) (mj.JSONBigInt, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return mj.JSONBigInt{}, err
	}

	bi, err := p.parseBigInt(strVal, format)
	return mj.JSONBigInt{
		Value:    bi,
		Original: strVal,
	}, err
}

func (p *Parser) parseBigInt(strRaw string, format bigIntParseFormat) (*big.Int, error) {
	bytes, err := p.parseAnyValueAsByteArray(strRaw)
	if err != nil {
		return nil, err
	}
	switch format {
	case bigIntSignedBytes:
		return twos.FromBytes(bytes), nil
	case bigIntUnsignedBytes:
		return big.NewInt(0).SetBytes(bytes), nil
	default:
		return nil, errors.New("unknown format requested")
	}
}

func (p *Parser) processCheckUint64(obj oj.OJsonObject) (mj.JSONCheckUint64, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckUint64{
			Value:    0,
			IsStar:   true,
			Original: "*"}, nil
	}

	ju, err := p.processUint64(obj)
	if err != nil {
		return mj.JSONCheckUint64{}, err
	}
	return mj.JSONCheckUint64{
		Value:    ju.Value,
		IsStar:   false,
		Original: ju.Original}, nil

}

func (p *Parser) processUint64(obj oj.OJsonObject) (mj.JSONUint64, error) {
	bi, err := p.processBigInt(obj, bigIntUnsignedBytes)
	if err != nil {
		return mj.JSONUint64{}, err
	}

	if bi.Value == nil || !bi.Value.IsUint64() {
		return mj.JSONUint64{}, errors.New("value is not uint64")
	}

	return mj.JSONUint64{
		Value:    bi.Value.Uint64(),
		Original: bi.Original}, nil
}

func (p *Parser) parseString(obj oj.OJsonObject) (string, error) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", errors.New("not a string value")
	}
	return str.Value, nil
}

func IsStar(obj oj.OJsonObject) bool {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return false
	}
	return str.Value == "*"
}
