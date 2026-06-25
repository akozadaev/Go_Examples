package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Привет"
	fmt.Println(len(s))                    // 12 (байт)
	fmt.Println(utf8.RuneCountInString(s)) // 6
	for i, r := range s {                  // i - смещение в байтах, r - руна
		fmt.Printf("%d: %c %U\n", i, r, r)
	}

	r, i := firstRune(s)
	fmt.Printf("%d: %c\n", i, r)

	byteLen, runeCount, firstByte, firstRune, _ := AnalyzeString(s)
	fmt.Printf("%v %v %v %U \n", byteLen, runeCount, firstByte, firstRune)
}

func firstRune(s string) (rune, int) {
	return utf8.DecodeRuneInString(s)
}

func AnalyzeString(s string) (byteLen, runeCount, firstByte int, firstRune rune, err error) {
	if s == "" {
		return 0, 0, 0, 0, errors.New("empty string")
	}
	byteLen = len(s)
	runeCount = utf8.RuneCountInString(s)
	firstByte = int(s[0])
	firstRune, _ = utf8.DecodeRuneInString(s)
	return
}
