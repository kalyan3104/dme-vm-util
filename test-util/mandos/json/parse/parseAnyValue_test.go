package mandosjsonparse

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestBool(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("true")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = p.parseAnyValueAsByteArray("false")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestString(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("``abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)

	result, err = p.parseAnyValueAsByteArray("``")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("```")
	require.Nil(t, err)
	require.Equal(t, []byte("`"), result)

	result, err = p.parseAnyValueAsByteArray("`` ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)
}

func TestUnsignedNumber(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = p.parseAnyValueAsByteArray("0x")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("12")
	require.Nil(t, err)
	require.Equal(t, []byte{12}, result)

	result, err = p.parseAnyValueAsByteArray("256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)

	result, err = p.parseAnyValueAsByteArray("0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = p.parseAnyValueAsByteArray("0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x05}, result)
}

func TestSignedNumber(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = p.parseAnyValueAsByteArray("255")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = p.parseAnyValueAsByteArray("+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = p.parseAnyValueAsByteArray("0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = p.parseAnyValueAsByteArray("+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = p.parseAnyValueAsByteArray("-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0x00}, result)

	result, err = p.parseAnyValueAsByteArray("-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)
}

func TestConcat(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("0x01|5")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = p.parseAnyValueAsByteArray("|||0x01|5||||")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = p.parseAnyValueAsByteArray("|")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("|||||||")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("|0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = p.parseAnyValueAsByteArray("``a|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = p.parseAnyValueAsByteArray("``a|0x62")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = p.parseAnyValueAsByteArray("0x61|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)
}

func TestKeccak256(t *testing.T) {
	p := Parser{}
	result, err := p.parseAnyValueAsByteArray("keccak256:0x01|5")
	require.Nil(t, err)
	expected, _ := keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:|||0x01|5||||")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:|")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:|||||||")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:|0")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:``a|``b")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:``a|0x62")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:0x61|``b")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	// some values from the old ERC20 tests
	result, err = p.parseAnyValueAsByteArray("keccak256:1|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("19efaebcc296cffac396adb4a60d54c05eff43926a6072498a618e943908efe1")
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:1|0x7777777777777777777707777777777777777777777777177777777777771234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("a3da7395b9df9b4a0ad4ce2fd40d2db4c5b231dbc2a19ce9bafcbc2233dc1b0a")
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:1|0x5555555555555555555505555555555555555555555555155555555555551234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("648147902a606bf61e05b8b9d828540be393187d2c12a271b45315628f8b05b9")
	require.Equal(t, expected, result)

	result, err = p.parseAnyValueAsByteArray("keccak256:2|0x7777777777777777777707777777777777777777777777177777777777771234|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("e314ce9b5b28a5927ee30ba28b67ee27ad8779e1101baf4224590c8f1e287891")
	require.Equal(t, expected, result)

}

func TestFile(t *testing.T) {
	p := Parser{
		FileResolver: NewDefaultFileResolver(),
	}
	result, err := p.parseAnyValueAsByteArray("file:../integrationTests/exampleFile.txt")
	require.Nil(t, err)
	require.Equal(t, []byte("hello!"), result)
}
