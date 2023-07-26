package str

import (
	"bytes"
	"encoding/json"
)

func Pad(str string, length int, repStr string) string {
	l := len(str)
	if l >= length {
		return str
	} else {
		repStrLen := len(repStr)
		buf := bytes.NewBufferString(str)
		needRepL := length - l
		for i := 0; i < needRepL/repStrLen; i++ {
			buf.Write([]byte(repStr))
		}
		restL := needRepL % repStrLen
		if restL != 0 {
			buf.Write([]byte(repStr[0:restL]))
		}
		return buf.String()
	}

}

func Reverse(s string) string {
	a := []rune(s)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}

func SplitByPairSymbol(str, symbolL, symbolR string) (inner []string) {
	a := []rune(str)
	l := len(a)
	sL := []rune(symbolL)[0]
	sR := []rune(symbolR)[0]
	n := 0
	var s []rune

	for i := 0; i < l; i++ {
		if n == 0 {
			n++
		}
		if a[i] == sL {
			n++
		}
		if a[i] == sR {
			if n > 1 {
				inner = append(inner, string(s))
				s = s[0:0]
			}
			n = 1
		}
		if n > 1 && a[i] != sL {
			s = append(s, a[i])
		}
	}

	return
}

func ToFormattedJson(v any) string {
	b, _ := json.MarshalIndent(v, "", "    ")

	return string(b)
}
