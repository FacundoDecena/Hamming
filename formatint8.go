package main

import (
	"errors"
	"math/bits"
	"strconv"
)

// FormatInt8 returns the string representation of i in the given base,
// for 2 <= base <= 36. The result uses the lower-case letters 'a' to 'z'
// for digit values >= 10.
func FormatInt(i int8, base int) string {
	return formatBits(uint8(i), base)
}

const digits = "0123456789"

// formatBits computes the string representation of u in the given base.
// If neg is set, u is treated as negative int64 value. If append_ is
// set, the string is appended to dst and the resulting byte slice is
// returned as the first result value; otherwise the string is returned
// as the second result value.
//
func formatBits(u uint8, base int) string {

	var a [8 + 1]byte // +1 for sign of 64bit value in base 2
	i := len(a)

	shift := uint(bits.TrailingZeros(uint(base))) & 7
	b := uint8(base)
	m := uint(base) - 1 // == 1<<shift - 1
	for u >= b {
		i--
		a[i] = digits[uint(u)&m]
		u >>= shift
	}
	// u < base
	i--
	a[i] = digits[uint(u)]
	s := string(a[i:])
	return s
}

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string {
	return "strconv." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, str, ErrRange}
}

// ParseUint is like ParseInt but for unsigned numbers.
func ParseUint(s string, base int) (uint8, error) {
	var n uint8
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		}

		n *= uint8(base)

		n1 := n + uint8(d)

		n = n1
	}

	return n, nil
}

func ParseInt(s string, base int, bitSize int) (i int8, err error) {
	const fnParseInt = "ParseInt"

	// Pick off leading sign.
	s0 := s

	// Convert unsigned and check range.
	var un uint8
	un, err = ParseUint(s, base)

	cutoff := uint8(1 << uint(bitSize-1))
	if un >= cutoff {
		return int8(cutoff - 1), rangeError(fnParseInt, s0)
	}
	n := int8(un)

	return n, nil
}
