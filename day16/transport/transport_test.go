package transport

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractVersion(t *testing.T) {
	t.Run("version 1", func(t *testing.T) {
		hs := "38006F45291200"

		result, err := ExtractVersion(hs)
		if err != nil {
			t.Error(err)
		}
		if result != 1 {
			t.Errorf("version = %d; expected 1", result)
		}
	})
	t.Run("version 6", func(t *testing.T) {
		hs := "D2FE28"

		result, err := ExtractVersion(hs)
		if err != nil {
			t.Error(err)
		}
		if result != 6 {
			t.Errorf("version = %d; expected 6", result)
		}
	})
}

func TestExtractTypeID(t *testing.T) {
	t.Run("type 4", func(t *testing.T) {
		hs := "D2FE28"

		result, err := ExtractTypeID(hs)
		if err != nil {
			t.Error(err)
		}
		if result != 4 {
			t.Errorf("version = %d; expected 4", result)
		}
	})
	t.Run("type 6", func(t *testing.T) {
		hs := "38006F45291200"

		result, err := ExtractTypeID(hs)
		if err != nil {
			t.Error(err)
		}
		if result != 6 {
			t.Errorf("version = %d; expected 6", result)
		}
	})
}

func TestExtractLengthType(t *testing.T) {
	t.Run("type 0", func(t *testing.T) {
		hs := "38006F45291200"

		result, err := ExtractLengthType(hs)
		if err != nil {
			t.Error(err)
		}
		if result != 0 {
			t.Errorf("version = %d; expected 0", result)
		}
	})
}

func TestExtractLength(t *testing.T) {
	testCases := map[string]uint32{
		"E058F79802FA00A4": 5693,
		"38006F45291200": 27,
		"EE00D40C823060": 3,
	}
	for hs, expectedLength := range testCases {
		t.Run(hs, func(t *testing.T) {	
			result, err := ExtractLength(hs)
			if err != nil {
				t.Error(err)
			}
			if result != expectedLength {
				t.Errorf("length = %d; expected %d", result, expectedLength)
			}
		})
	}
}

func TestExtractLiteral(t *testing.T) {
	hs := "D2FE28"
	
	result, bitCount, err := ExtractLiteral(hs)

	if err != nil {
		t.Error(err)
	}
	if result != 2021 {
		t.Errorf("literal = %d; expected 2021", result)
	}
	if bitCount != 21 {
		t.Errorf("bitcount = %d; expected 21", bitCount)
	}
}

func TestSplitPartialBits(t *testing.T) {
	expectedResultF := make(map[uint8][2]byte)
	expectedResultF[0] = [2]byte{'0', 'f'}
	expectedResultF[1] = [2]byte{'8', '7'}
	expectedResultF[2] = [2]byte{'c', '3'}
	expectedResultF[3] = [2]byte{'e', '1'}
	expectedResultF[4] = [2]byte{'f', '0'}

	expectedResult5 := make(map[uint8][2]byte)
	expectedResult5[0] = [2]byte{'0', '5'}
	expectedResult5[1] = [2]byte{'0', '5'}
	expectedResult5[2] = [2]byte{'4', '1'}
	expectedResult5[3] = [2]byte{'4', '1'}
	expectedResult5[4] = [2]byte{'5', '0'}

	expectedResultA := make(map[uint8][2]byte)
	expectedResultA[0] = [2]byte{'0', 'a'}
	expectedResultA[1] = [2]byte{'8', '2'}
	expectedResultA[2] = [2]byte{'8', '2'}
	expectedResultA[3] = [2]byte{'a', '0'}
	expectedResultA[4] = [2]byte{'a', '0'}

	for i := uint8(0); i <= 4; i++ {
		t.Run(fmt.Sprintf("F-%d", i), func(t *testing.T) {
			result, err := SplitPartialBits('F', i)
			if err != nil {
				t.Error(err)
			}
			t.Run("Left", func(t *testing.T) {
				if result[0] != expectedResultF[i][0] {
					t.Errorf("left = '%c'; expected '%c'", result[0], expectedResultF[i][0])
				}
			})
			t.Run("Right", func(t *testing.T) {
				if result[1] != expectedResultF[i][1] {
					t.Errorf("right = %c; expected '%c'", result[1], expectedResultF[i][1])
				}
			})
		})
		t.Run(fmt.Sprintf("5-%d", i), func(t *testing.T) {
			result, err := SplitPartialBits('5', i)
			if err != nil {
				t.Error(err)
			}
			t.Run("Left", func(t *testing.T) {
				if result[0] != expectedResult5[i][0] {
					t.Errorf("left = '%c'; expected '%c'", result[0], expectedResult5[i][0])
				}
			})
			t.Run("Right", func(t *testing.T) {
				if result[1] != expectedResult5[i][1] {
					t.Errorf("right = '%c'; expected '%c'", result[1], expectedResult5[i][1])
				}
			})
		})
		t.Run(fmt.Sprintf("A-%d", i), func(t *testing.T) {
			result, err := SplitPartialBits('A', i)
			if err != nil {
				t.Error(err)
			}
			t.Run("Left", func(t *testing.T) {
				if result[0] != expectedResultA[i][0] {
					t.Errorf("left = '%c': expected '%c'", result[0], expectedResultA[i][0])
				}
			})
			t.Run("Right", func(t *testing.T) {
				if result[1] != expectedResultA[i][1] {
					t.Errorf("right = '%c'; expected '%c'", result[1], expectedResultA[i][1])
				}
			})
		})
	}
}

func TestRealignBits(t *testing.T) {
	type testCaseInputs struct {
		hexString string
		shiftn uint8
	}
	testCases := map[testCaseInputs]string{
		{"34d4", 2}: "d350",
		{"68EBA3E0F5F387A", 3}: "475D1F07AF9C3D0",
		{"6BEF8003F5555F1F111", 3}: "5F7C001FAAAAF8F8888",
	}

	for input, expectedOutput := range testCases {
		result, err := realignBits(input.hexString, input.shiftn)
		if err != nil {
			t.Error(err)
		}
		if result != strings.ToLower(expectedOutput) {
			t.Errorf("result = '%s'; expected '%s'", strings.ToUpper(result), strings.ToUpper(expectedOutput))
		}
	}
}

func TestSubPacketsFromBits(t *testing.T) {
	p := Packet{
		HexString: "38006F45291200",
		Parent: nil,
	}
	var err error
	p.Version, err = ExtractVersion("38006F45291200")
	if err != nil {
		t.Error(err)
	}
	p.TypeID, err = ExtractTypeID("38006F45291200")
	if err != nil {
		t.Error(err)
	}

	packets, err := p.subPacketsFromBits(27)
	if err != nil {
		t.Error(err)
	}
	if len(packets) != 2 {
		t.Errorf("sub packet count = %d; expected 2", len(packets))
	}
	expectedResults := []Packet{
		{HexString: "d14"},
		{HexString: "5224"},
	}
	for i, packet := range packets {
		t.Run(fmt.Sprintf("Packet%dTest", i+1), func(t *testing.T) {
			if packet.HexString != expectedResults[i].HexString {
				t.Errorf("hexstring = %s; expected %s", packet.HexString, expectedResults[i].HexString)
			}
		})
	}
}

func TestSubPacketsFromCount(t *testing.T) {
	p := Packet{
		HexString: "EE00D40C823060",
		Parent: nil,
	}

	var err error
	p.Version, err = ExtractVersion(p.HexString)
	if err != nil {
		t.Error(err)
	}
	p.TypeID, err = ExtractTypeID(p.HexString)
	if err != nil {
		t.Error(err)
	}
	packets, bitCount, err := p.subPacketsFromCount(3)
	if err != nil {
		t.Error(err)
	}
	if pl:= len(packets); pl != 3 {
		t.Errorf("packets = %d; expected 3", pl)
	}
	if bitCount != 33 {
		t.Errorf("bitcount = %d; expected 33", bitCount)
	}

	expectedResults := []Packet{
		{HexString: "502"},
		{HexString: "904"},
		{HexString: "306"},
	}
	for i, packet := range packets {
		t.Run(fmt.Sprintf("Packet%dTest", i+1), func(t *testing.T) {
			if packet.HexString != expectedResults[i].HexString {
				t.Errorf("hexstring = %s; expected %s", packet.HexString, expectedResults[i].HexString)
			}
		})
	}
}

func TestChildren(t *testing.T) {
	testCases := map[string]int {
		"8A004A801A8002F478": 16,
		"620080001611562C8802118E34":12,
		"C0015000016115A2E0802F182340":23,
		"A0016C880162017C3686B18A3D4780":31,
	}
	for hs, expectedSum := range testCases {
		VersionSum = 0
		t.Run(hs, func(t *testing.T) {
			p := Packet{
				HexString: hs,
				Parent: nil,
			}
			
			var err error
			p.Version, err = ExtractVersion(p.HexString)
			if err != nil {
				t.Error(err)
			}
			_, _, err = p.Children()
			if err != nil {
				t.Error(err)
			}
			if result := VersionSum + int(p.Version); result != expectedSum {
				t.Errorf("version sum = %d; expected %d", result, expectedSum)
			}
		})
	}
}