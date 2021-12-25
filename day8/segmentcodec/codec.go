package segmentcodec

import (
	"strings"
	"strconv"
	"fmt"
)

var DIGIT_STRINGS = [...]string{"abcefg", "cf", "acdeg", "acdfg", "bcdf", "abdfg", "abdefg", "acf", "abcdefg", "abcdfg"}

type Converter struct {
	decodeSegmentMap map[string]string
	encodeSegmentMap map[string]string
	EncodedDigits [10]string
}

func (c *Converter) ConvertToDigit(str string) (rune, error) {
	converted := ""
	for _, r := range strings.Split(str, "") {
		converted += c.decodeSegmentMap[r]
	}
	
	for i, digitStr := range DIGIT_STRINGS {
		if DigitStringMatch(digitStr, converted) {
			return rune(strconv.Itoa(i)[0]), nil
		}
	}
	return 0, fmt.Errorf("failed to find corresponding digit for string %s", str)
}

func (c *Converter) PopulateConverter (uniqueDigits []string) error {
	inputOccurrences := SegmentSectionOccurrences(uniqueDigits)
	c.encodeSegmentMap = make(map[string]string)
	c.decodeSegmentMap = make(map[string]string)
	
	for inputr, count := range inputOccurrences {
		switch count {
		case 4:
			c.encodeSegmentMap["e"] = inputr
		case 6:
			c.encodeSegmentMap["b"] = inputr
		case 9:
			c.encodeSegmentMap["f"] = inputr
		}
	}

	for _, ds := range uniqueDigits {
		if len(ds) == 2 {
			if string(ds[0]) != c.encodeSegmentMap["f"] {
				c.encodeSegmentMap["c"] = string(ds[0])
			} else {
				c.encodeSegmentMap["c"] = string(ds[1])
			}
			break
		}
	}

	for _, ds := range uniqueDigits {
		if len(ds) == 3 {
			for _, r := range strings.Split(ds, "") {
				if r != c.encodeSegmentMap["c"] && r != c.encodeSegmentMap["f"] {
					c.encodeSegmentMap["a"] = r
					break
				}
			}
			break
		}
	}

	for _, ds := range uniqueDigits {
		if len(ds) == 4 {
			for _, r := range strings.Split(ds, "") {
				if r != c.encodeSegmentMap["b"] && r != c.encodeSegmentMap["c"] && r != c.encodeSegmentMap["f"] {
					c.encodeSegmentMap["d"] = r
					break
				}
			}
			break
		}
	}

	for _, ds := range uniqueDigits {
		if len(ds) == 7 {
			for _, r := range strings.Split(ds, "") {
				if r != c.encodeSegmentMap["a"] && 
				r != c.encodeSegmentMap["b"] && 
				r != c.encodeSegmentMap["c"] && 
				r != c.encodeSegmentMap["d"] && 
				r != c.encodeSegmentMap["e"] &&
				r != c.encodeSegmentMap["f"] {
					c.encodeSegmentMap["g"] = r
					break
				}
			}
			break
		}
	}

	for dec, enc := range c.encodeSegmentMap {
		c.decodeSegmentMap[enc] = dec
	}

	for i, digitStr := range DIGIT_STRINGS {
		newStr, err := c.EncodeDigitStr(digitStr)
		if err != nil {
			panic(err)
		}
		c.EncodedDigits[i] = newStr
	}

	return nil
}

func (c *Converter) EncodeDigitStr(digitStr string) (string, error) {
	result := ""
	for _, r := range strings.Split(digitStr, "") {
		result += c.encodeSegmentMap[r]
	}
	return result, nil
}

func (c *Converter) DecodeDigitStr(digitStr string) (string, error) {
	result := ""
	for _, r := range strings.Split(digitStr, "") {
		result += c.decodeSegmentMap[r]
	}
	return result, nil
}

func SegmentSectionOccurrences(digitStrs []string) map[string]int {
	occurrences := make(map[string]int)
	for _, digitStr := range digitStrs {
		for _, r := range strings.Split(digitStr, "") {
			occurrences[r] += 1
		}
	}
	return occurrences
}

/* Test to see if all characters are in both strings */
func DigitStringMatch(str1 string, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}
	return containsAll(str1, str2)
}

/* Test to see if all chars of substr are in str */
func containsAll(str string, substr string) bool {
	for _, r := range substr {
		if !strings.ContainsRune(str, r) {
			return false
		}
	}
	return true
}