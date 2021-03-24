package seqt

func byteIsDigit(b byte) bool {
	return ('0' <= b) && (b <= '9')
}

func byteIsUpperLetter(b byte) bool {
	return ('A' <= b) && (b <= 'Z')
}

func byteIsLowerLetter(b byte) bool {
	return ('a' <= b) && (b <= 'z')
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
