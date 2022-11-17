package randoms

import (
	"math/rand"
)

const (
	defaultNumberSet     = "0123456789"
	defaultAlphaLowerSet = "abcdefghijklmnopqrstuvwxyz"
	defaultAlphaUpperSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultSymbolSet     = "!$%^&*()_+{}:@[];'#<>?,./|\\-=?"
	asciiSetConst        = ` !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_\abcdefghijklmnopqrstuvwxyz{|}~` + "`"
)

var (
	lowerAlphaSet   = []rune(defaultAlphaLowerSet)
	upperAlphaSet   = []rune(defaultAlphaUpperSet)
	symbolSet       = []rune(defaultSymbolSet)
	numberSet       = []rune(defaultNumberSet)
	alphaSet        = []rune(defaultAlphaLowerSet + defaultAlphaUpperSet)
	alphaNumericSet = []rune(defaultAlphaLowerSet + defaultAlphaUpperSet + defaultNumberSet)
	allCharSet      = []rune(defaultAlphaLowerSet + defaultAlphaUpperSet + defaultNumberSet + defaultSymbolSet)
	asciiSet        = []rune(asciiSetConst)
)

func RandomStringAscii(maxLength int, minLength int) string {
	if minLength > maxLength {
		minLength = maxLength
	}

	length := maxLength
	if minLength > 0 && minLength != maxLength {
		length = rand.Intn(maxLength-minLength) + minLength
	}

	b := make([]rune, length)
	for i := range b {
		b[i] = asciiSet[rand.Intn(len(asciiSet))]
	}

	return string(b)
}

func RandomStringAlpha(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphaSet[rand.Intn(len(alphaSet))]
	}
	return string(b)
}

func RandomStringAlphaNumericRandomLength(maxLength int, minLength int) string {
	length := maxLength

	if minLength > maxLength {
		minLength = maxLength
	}

	if minLength > 0 && minLength != maxLength {
		length = rand.Intn(maxLength-minLength) + minLength
	}

	b := make([]rune, length)
	for i := range b {
		b[i] = alphaNumericSet[rand.Intn(len(alphaNumericSet))]
	}

	return string(b)
}

func RandomStringAlphaNumeric(length int) string {
	return RandomStringAlphaNumericRandomLength(length, length)
}

func RandomStringPassword(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = allCharSet[rand.Intn(len(allCharSet))]
	}
	return string(b)
}
