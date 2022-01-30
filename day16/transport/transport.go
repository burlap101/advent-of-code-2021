package transport

import (
	"fmt"
	"strconv"
)

type Packet struct {
	HexString string
	Version uint8
	TypeID uint8
	BitCount uint32
	Parent *Packet
}

const VersionMask uint8 = 0xE0
const LengthTypeMask uint8 = 0x02

var VersionSum int = 0


func (p Packet) Children() ([]Packet, int, error) {
	if len(p.HexString) < 2 {
		return []Packet{}, 0, nil
	}
	if p.TypeID == 4 {
		return nil, 0, fmt.Errorf("no children for packet typeid = 4")
	} 
	lengthType, err := ExtractLengthType(p.HexString)
	if err != nil  {
		return nil, 0, err
	}

	if lengthType == 0 {
		lengthInBits, err := ExtractLength(p.HexString)
		if err != nil {
			return nil, 0, err
		}
		p.BitCount = 22 + lengthInBits
		packets, err := p.subPacketsFromBits(uint16(lengthInBits))
		if err != nil {
			return nil, 0, err
		}
		for _, pkt := range packets {
			VersionSum += int(pkt.Version)
		}
		
		return packets, int(p.BitCount), nil
	} else if lengthType == 1 {
		lengthInPackets, err := ExtractLength(p.HexString)
		if err != nil {
			return nil, 0, err
		}
		packets, bitCount, err := p.subPacketsFromCount(uint16(lengthInPackets))
		if err != nil {
			return nil, 0, err
		}
		p.BitCount = uint32(18 + bitCount)
		for _, pkt := range packets {
			VersionSum += int(pkt.Version)
		}
		return packets, int(p.BitCount), nil
	}
	
	return nil, 0, fmt.Errorf("unknown length type %d", lengthType)
}


func (p Packet) subPacketsFromBits(n uint16) ([]Packet, error) {
	packets := make([]Packet, 0) 
	totalNumberOfBitsToMask := int(n) + 22 	// Parent's header size will be 22 bits
	fullSubStringWParentHeader := p.HexString[:totalNumberOfBitsToMask/4]
	// The following will retrieve the ramining bits and convert append to hex string
	if numBits := totalNumberOfBitsToMask % 4; numBits > 0 {
		digit := p.HexString[totalNumberOfBitsToMask/4]
		split, err := SplitPartialBits(digit, uint8(numBits))
		if err != nil {
			return nil, err
		}
		fullSubStringWParentHeader += string(split[0])
	}
	childString := fullSubStringWParentHeader[6:]
	shiftn := uint8(2)
	split, err := SplitPartialBits(fullSubStringWParentHeader[5], shiftn)
	if err != nil {
		return nil, err
	}
	for i, numBitsAllocated := 0, 0; numBitsAllocated < int(n); i++ {
		packets = append(packets, Packet{Parent: &p})
		childString = string(split[1]) + childString
		childString, err = realignBits(childString, shiftn)
		if err != nil {
			return nil, err
		}
		packets[i].Version, err = ExtractVersion(childString)
		if err != nil {
			return nil, err
		}
		packets[i].TypeID, err = ExtractTypeID(childString)
		if err != nil {
			return nil, err
		}
		if packets[i].TypeID == 4 {
			_, bitCount, err := ExtractLiteral(childString)
			if err != nil {
				return nil, err
			}
			numBitsAllocated += bitCount
			packets[i].HexString = childString[:bitCount / 4]
			modbits := uint8(bitCount % 4)
			if modbits > 0 {
				split, err = SplitPartialBits(childString[bitCount / 4], modbits)
				if err != nil {
					return nil, err
				}
				packets[i].HexString += string(split[0])
				shiftn = modbits
			} else {
				split = [2]byte{0, childString[bitCount/4]}
				shiftn = 0
			}
			childString = childString[(bitCount / 4)+1:]
		} else {
			packets[i].HexString = childString
			_, bitCount, err := packets[i].Children()
			if err != nil {
				return nil, err
			}
			
			numBitsAllocated += bitCount
			
			packets[i].HexString = childString[:bitCount / 4]
			modbits := uint8(bitCount % 4)
			if modbits > 0 {
				split, err = SplitPartialBits(childString[bitCount / 4], modbits)
				if err != nil {
					return nil, err
				}
				packets[i].HexString += string(split[0])
				shiftn = modbits
			} else {
				shiftn = 0
				split = [2]byte{0, childString[bitCount/4]}
			}
			childString = childString[(bitCount / 4)+1:]
		}
	}
	return packets, nil
}


func (p Packet) subPacketsFromCount(n uint16) ([]Packet, int, error) {
	numBitsAllocated := 0
	packets := make([]Packet, 0)
	childString := p.HexString[5:]
	split, err := SplitPartialBits(p.HexString[4], 2)
	if err != nil {
		return nil, 0, err
	}
	
	var shiftn uint8 = 2
	for i, numPacketsAllocated := 0, 0; numPacketsAllocated < int(n); i++ {
		childString = string(split[1]) + childString
		childString, err = realignBits(childString, shiftn)
		if err != nil {
			return nil, 0, err
		}
		packets = append(packets, Packet{Parent: &p})
		packets[i].Version, err = ExtractVersion(childString)
		if err != nil {
			return nil, 0, err
		}
		packets[i].TypeID, err = ExtractTypeID(childString)
		if err != nil {
			return nil, 0, err
		}
		
		if packets[i].TypeID == 4 {
			_, bitCount, err := ExtractLiteral(childString)
			if err != nil {
				return nil, 0, err
			}
			numBitsAllocated += bitCount
			packets[i].HexString = childString[:bitCount / 4]
			modbits := uint8(bitCount % 4)
			packets[i].BitCount = uint32(bitCount)
			if modbits > 0 {
				split, err = SplitPartialBits(childString[bitCount / 4], modbits)
				if err != nil {
					return nil, 0, err
				}
				packets[i].HexString += string(split[0])
				shiftn = modbits
			} else {
				shiftn = 0
				split = [2]byte{0, childString[bitCount/4]}
			}
			childString = childString[(bitCount / 4)+1:]
			numPacketsAllocated += 1
		} else {
			//TODO: Update childstring, removing the packet returned via call to Children()
			packets[i].HexString = childString
			_, bitCount, err := packets[i].Children()
			if err != nil {
				return nil, 0, err
			}
			numBitsAllocated += int(bitCount)
			packets[i].HexString = childString[:bitCount / 4]
			modbits := uint8(bitCount % 4)
			packets[i].BitCount = uint32(bitCount)
			if modbits > 0 {
				split, err = SplitPartialBits(childString[bitCount / 4], modbits)
				if err != nil {
					return nil, 0, err
				}
				packets[i].HexString += string(split[0])
				shiftn = modbits
			} else {
				shiftn = 0
				split = [2]byte{0, childString[bitCount/4]}
			}
			childString = childString[(bitCount / 4) + 1:]
			numPacketsAllocated += 1
		}
	}
	return packets, numBitsAllocated, nil
}


// returns the header length in bits for the packet
func (p Packet) HeaderLength() (int, error) {
	if p.TypeID == 4 {
		return 6, nil
	}
	lt, err := ExtractLengthType(p.HexString)
	if err != nil {
		return 0, err
	}
	if lt == 0 {
		return 22, nil
	}
	if lt == 1 {
		return 18, nil
	}
	return 0, fmt.Errorf("unable to determine header length for packet %+v", p)
}


//This will refactor a hex string so that there are no unused bits at the beginning
func realignBits(hexString string, shiftn uint8) (string, error) {
	originalLength := len(hexString)
	overFlowMask := ^(^uint8(0) >> shiftn)
	newHexString := ""
	if len(hexString) % 2 > 0 {
		hexString += "0"
	}
	nextOverFlow := uint8(0)
	for i := len(hexString)-2; i >= 0; i -= 2 {
		window := hexString[i:i+2]
		windowUint64, err := strconv.ParseUint(window, 16, 8)
		if err != nil {
			return "", err
		}
		windowUint8 := uint8(windowUint64)
		overFlow := nextOverFlow
		nextOverFlow = (windowUint8 & overFlowMask) >> ( 8 - shiftn )
		newHexByte := strconv.FormatUint(uint64((windowUint8 << shiftn) | overFlow), 16)
		if len(newHexByte) == 1 {
			newHexByte = "0" + newHexByte
		}
		newHexString =  newHexByte + newHexString
	}

	return newHexString[:originalLength], nil
}


func SplitPartialBits(digit byte, leftBitCount uint8) ([2]byte, error) {
	if leftBitCount > 4 {
		return [2]byte{}, fmt.Errorf("invalid bit count %d; needs to be 4 or less", leftBitCount)
	}
	dNum, err := strconv.ParseUint(string(digit), 16, 8)
	result := [2]byte{}
	if err != nil {
		return [2]byte{}, err
	}
	leftMask := ^(^uint64(0) << uint64(leftBitCount))
	result[0] = strconv.FormatUint(((dNum >> uint64(4-leftBitCount)) & leftMask) << (4-leftBitCount), 16)[0]
	rightMask := ^(^uint64(0) << uint64(4-leftBitCount))
	result[1] = strconv.FormatUint(dNum & rightMask, 16)[0]

	return result, nil
}


func ExtractVersion(hexString string) (uint8, error) {
	extractedNum, err := strconv.ParseUint(hexString[:2], 16, 8)
	if err != nil {
		return 0, err
	}
	firstByte := uint8(extractedNum)
	version := firstByte & VersionMask >> 5
	
	return version, nil
}

func ExtractTypeID(hexString string) (uint8, error) {
	mask := uint8(0x1C)
	extractedNum, err := strconv.ParseUint(hexString[:2], 16, 8)
	if err != nil {
		return 0, err
	}
	firstByte := uint8(extractedNum)

	return (firstByte & mask) >> 2, nil
}


func ExtractLengthType(hexString string) (uint8, error) {
	extractedNum, err := strconv.ParseUint(hexString[:2], 16, 8)
	if err != nil {
		return 0, err
	}
	firstByte := uint8(extractedNum)
	
	return firstByte & LengthTypeMask >> 1, nil
}


func ExtractLength(hexString string) (uint32, error) {
	const LengthMask11 uint32 = 0x01FFC000
	const LengthMask15 uint32 = 0x01FFFC00

	lengthType, err := ExtractLengthType(hexString)
	if err != nil {
		return 0, err
	}
	extractedNum, err := strconv.ParseUint(hexString[:8], 16, 32)
	if err != nil {
		return 0, err
	}
	first32 := uint32(extractedNum)

	if lengthType == 0 {
		return first32 & LengthMask15 >> 10, nil
	} else if lengthType == 1 {
		return first32 & LengthMask11 >> 14, nil
	}
	return 0, fmt.Errorf("length type = %d; expected either 0 or 1", lengthType)
}

// Returns the literal, count of bits used in hte packet and/or an error
// The bit count includes all digits extracted plus their last indicator flag plus the packet's header. 
// Trailing 0's are not included in bitCount
func ExtractLiteral(hexString string) (literal uint64, bitCount int, err error) {
	typeMask := uint8(0x10)
	digitMask := uint8(0x0F)
	extractedHexString := ""
	// hsNumberOfBits := len(hexString) * 4
	bitCount = 6
	
	literalString := hexString[1:] // remove the first digit
	literalString, err = realignBits(literalString, 2)
	if err != nil {
		return 0,0, err
	}

	// extractedNum, err := strconv.ParseUint(hexString, 16, 64)
	// if err != nil {
	// 	return 0, 0, err
	// }
	isFinalDigit := false
	for i := 0; !isFinalDigit; i++ {
		digitwTypeString := literalString[0:2]
		split, err := SplitPartialBits(digitwTypeString[1], 1)
		if err != nil {
			return 0,0,err
		}
		digitwTypeString = string(digitwTypeString[0]) + string(split[0])
		literalString = string(split[1]) + literalString[2:]
		literalString, err = realignBits(literalString, 1)
		if err != nil {
			return 0,0,err
		}
		extractedNum, err := strconv.ParseUint(digitwTypeString, 16, 8)
		if err != nil {
			return 0,0,err
		}
		shiftedDigitwType := uint8(extractedNum) >> 3
		isFinalDigit = (shiftedDigitwType & typeMask) == 0
		extractedHexString += strconv.FormatUint(uint64(shiftedDigitwType & digitMask), 16)
		bitCount += 5
	}
	
	result, err := strconv.ParseUint(extractedHexString, 16, 64)
	if err != nil {
		return 0, 0, err
	}

	return result, bitCount, nil
}
