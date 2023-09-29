package mandosjsonparse

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	twos "github.com/kalyan3104/dme-components-big-int/twos-complement"
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	oj "github.com/kalyan3104/dme-vm-util/test-util/orderedjson"
)

const filePrefix = "file:"
const keccak256Prefix = "keccak256:"

func (p *Parser) parseCheckBytes(obj oj.OJsonObject) (mj.JSONCheckBytes, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return mj.JSONCheckBytes{
			Value:    nil,
			IsStar:   true,
			Original: "*"}, nil
	}

	jb, err := p.processAnyValueAsByteArray(obj)
	if err != nil {
		return mj.JSONCheckBytes{}, err
	}
	return mj.JSONCheckBytes{
		Value:    jb.Value,
		IsStar:   false,
		Original: jb.Original,
	}, nil
}

func (p *Parser) processAnyValueAsByteArray(obj oj.OJsonObject) (mj.JSONBytes, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return mj.JSONBytes{}, err
	}
	result, err := p.parseAnyValueAsByteArray(strVal)
	return mj.JSONBytes{
		Value:    result,
		Original: strVal,
	}, err
}

func (p *Parser) parseAnyValueAsByteArray(strRaw string) ([]byte, error) {
	if len(strRaw) == 0 {
		return []byte{}, nil
	}

	// file contents
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, filePrefix) {
		if p.FileResolver == nil {
			return []byte{}, errors.New("parser FileResolver not provided")
		}
		fileContents, err := p.FileResolver.ResolveFileValue(strRaw[len(filePrefix):])
		if err != nil {
			return []byte{}, err
		}
		return fileContents, nil
	}

	// keccak256
	// TODO: make this part of a proper parser
	if strings.HasPrefix(strRaw, keccak256Prefix) {
		arg, err := p.parseAnyValueAsByteArray(strRaw[len(keccak256Prefix):])
		if err != nil {
			return []byte{}, fmt.Errorf("cannot parse keccak256 argument: %w", err)
		}
		hash, err := keccak256(arg)
		if err != nil {
			return []byte{}, fmt.Errorf("error computing keccak256: %w", err)
		}
		return hash, nil
	}

	// concatenate values of different formats
	// TODO: make this part of a proper parser
	parts := strings.Split(strRaw, "|")
	if len(parts) > 1 {
		concat := make([]byte, 0)
		for _, part := range parts {
			eval, err := p.parseAnyValueAsByteArray(part)
			if err != nil {
				return []byte{}, err
			}
			concat = append(concat, eval...)
		}
		return concat, nil
	}

	if strRaw == "false" {
		return []byte{}, nil
	}

	if strRaw == "true" {
		return []byte{0x01}, nil
	}

	// allow ascii strings, for readability
	if strings.HasPrefix(strRaw, "``") || strings.HasPrefix(strRaw, "''") {
		str := strRaw[2:]
		return []byte(str), nil
	}

	// signed numbers
	if strRaw[0] == '-' || strRaw[0] == '+' {
		numberBytes, err := p.parseUnsignedNumberAsByteArray(strRaw[1:])
		if err != nil {
			return []byte{}, err
		}
		number := big.NewInt(0).SetBytes(numberBytes)
		if strRaw[0] == '-' {
			number = number.Neg(number)
		}
		return twos.ToBytes(number), nil
	}

	// unsigned numbers
	return p.parseUnsignedNumberAsByteArray(strRaw)
}

func (p *Parser) parseUnsignedNumberAsByteArray(strRaw string) ([]byte, error) {
	str := strings.ReplaceAll(strRaw, "_", "") // allow underscores, to group digits
	str = strings.ReplaceAll(str, ",", "")     // also allow commas to group digits

	// hex, the usual representation
	if strings.HasPrefix(strRaw, "0x") || strings.HasPrefix(strRaw, "0X") {
		str := strRaw[2:]
		if len(str)%2 == 1 {
			str = "0" + str
		}
		return hex.DecodeString(str)
	}

	// binary representation
	if strings.HasPrefix(strRaw, "0b") || strings.HasPrefix(strRaw, "0B") {
		result := new(big.Int)
		var parseOk bool
		result, parseOk = result.SetString(str[2:], 2)
		if !parseOk {
			return []byte{}, fmt.Errorf("could not parse binary value: %s", strRaw)
		}

		return result.Bytes(), nil
	}

	// default: parse as BigInt, base 10
	result := new(big.Int)
	var parseOk bool
	result, parseOk = result.SetString(str, 10)
	if !parseOk {
		return []byte{}, fmt.Errorf("could not parse base 10 value: %s", strRaw)
	}

	return result.Bytes(), nil
}
