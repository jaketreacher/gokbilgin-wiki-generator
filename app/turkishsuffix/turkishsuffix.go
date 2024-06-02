package turkishsuffix

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	voicelessConsonants = []rune{'ç', 'f', 'h', 'k', 'p', 's', 'ş', 't'}
	aVowels             = []rune{'a', 'ı', 'o', 'u'}
	eVowels             = []rune{'e', 'i', 'ö', 'ü', 'ī'}
	vowels              = append(aVowels, eVowels...)
)

type vowelGroupType rune

var vowelGroup = struct {
	A vowelGroupType
	E vowelGroupType
}{
	A: 'a',
	E: 'e',
}

func Ablative(text string) (string, error) {
	var ablativeBuilder []rune

	if isProperNoun(text) {
		ablativeBuilder = append(ablativeBuilder, '\'')
	}

	lower := strings.ToLowerSpecial(unicode.TurkishCase, text)
	if endsWithVoicelessConsonant(lower) {
		ablativeBuilder = append(ablativeBuilder, 't')
	} else {
		ablativeBuilder = append(ablativeBuilder, 'd')
	}

	vowels := getVowels(lower)
	lastVowel := vowels[len(vowels)-1]
	vowelGroup, err := getVowelGroup(lastVowel)

	if err != nil {
		return "", err
	}

	ablativeBuilder = append(ablativeBuilder, rune(vowelGroup))
	ablativeBuilder = append(ablativeBuilder, 'n')

	return text + string(ablativeBuilder), nil
}

func getVowelGroup(r rune) (vowelGroupType, error) {
	for _, v := range aVowels {
		if v == r {
			return vowelGroup.A, nil
		}
	}

	for _, v := range eVowels {
		if v == r {
			return vowelGroup.E, nil
		}
	}

	return vowelGroup.A, fmt.Errorf("invalid vowel: %s", string(r))
}

func getVowels(text string) []rune {
	var vowels []rune
	for _, r := range text {
		if isVowel(r) {
			vowels = append(vowels, r)
		}
	}

	return vowels
}

func isVowel(r rune) bool {
	for _, v := range vowels {
		if r == v {
			return true
		}
	}

	return false
}

func endsWithVoicelessConsonant(text string) bool {
	lastChar := getLastChar(text)

	for _, r := range voicelessConsonants {
		if lastChar == r {
			return true
		}
	}

	return false
}

func getLastChar(text string) rune {
	runes := []rune(text)
	return runes[len(runes)-1]
}

func isProperNoun(text string) bool {
	runes := []rune(text)
	return unicode.IsUpper(runes[0])
}
