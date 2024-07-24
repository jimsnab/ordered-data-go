package ordereddata

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	orderedJson struct {
		content []byte
	}

	parserPos struct {
		pos  int
		line int
	}
)

var ErrInvaldJsonValue = errors.New("invalid json value")

func newOrderedJson(content []byte) *orderedJson {
	return &orderedJson{
		content: content,
	}
}

func (oj *orderedJson) preview(pp parserPos) string {
	by := oj.content[pp.pos:]
	if len(by) == 0 {
		return "[EOF]"
	}

	str := string(by)

	lines := strings.Split(str, "\n")
	if len(lines) > 3 {
		lines = lines[:3]
	}
	str = strings.Join(lines, "\n")

	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\t", "\\t")

	return fmt.Sprintf("line %d pos %d: %s", pp.line, pp.pos, str)
}

func (oj *orderedJson) peekNextRune(pp parserPos) (r rune, nextPp parserPos) {
	if pp.pos >= len(oj.content) {
		return
	}

	r, size := utf8.DecodeRune(oj.content[pp.pos:])
	nextPp.pos = pp.pos + size
	if r == '\n' {
		nextPp.line = pp.line + 1
	} else {
		nextPp.line = pp.line
	}
	return
}

func (oj *orderedJson) peekString(pp parserPos) (s string, nextPp parserPos) {
	r, pp := oj.peekNextRune(pp)
	if r != '"' {
		return
	}

	var sb strings.Builder
	for {
		r, pp = oj.peekNextRune(pp)
		if pp.line == 0 {
			return
		}

		if r == '"' {
			nextPp = pp
			s = sb.String()
			return
		}

		if r == '\\' {
			r, pp = oj.peekNextRune(pp)
			if pp.line == 0 {
				return
			}
			if r == 'u' {
				var digits strings.Builder
				for i := 0; i < 4; i++ {
					r, pp = oj.peekNextRune(pp)
					if pp.line == 0 {
						return
					}
					digits.WriteRune(r)
				}
				ch, err := strconv.ParseInt(digits.String(), 16, 32)
				if err != nil {
					return
				}
				r = rune(ch)
			}
		}
		sb.WriteRune(r)
	}
}

func (oj *orderedJson) peekFloat(pp parserPos) (f float64, nextPp parserPos) {
	var sb strings.Builder
	var r rune
	var pp2 parserPos

	for {
		r, pp = oj.peekNextRune(pp)
		if pp.line == 0 {
			break
		}

		if !unicode.IsDigit(r) && r != '.' && r != 'e' && r != 'E' && r != '-' && r != '+' {
			break
		}
		sb.WriteRune(r)
		pp2 = pp
	}

	if pp2.line == 0 {
		return
	}

	f, err := strconv.ParseFloat(sb.String(), 64)
	if err != nil {
		return
	}

	nextPp = pp2
	return
}

func (oj *orderedJson) peekValueKeyword(pp parserPos) (v any, nextPp parserPos) {
	var r rune
	var sb strings.Builder
	for {
		r, pp = oj.peekNextRune(pp)
		if pp.line == 0 {
			break
		}

		sb.WriteRune(r)
		if sb.String() == "true" {
			v = true
			break
		}

		if sb.String() == "false" {
			v = false
			break
		}

		if sb.String() == "null" {
			v = nil
			break
		}

		if sb.Len() >= 5 {
			return
		}
	}

	nextPp = pp
	return
}

func (oj *orderedJson) skipWhitespace(pp parserPos) (nextPp parserPos) {
	var r rune
	for {
		nextPp = pp
		r, pp = oj.peekNextRune(pp)
		if pp.line == 0 || !unicode.IsSpace(r) {
			break
		}
	}
	return
}

func (oj *orderedJson) peekJsonMap(pp parserPos) (m StringMap, nextPp parserPos, err error) {
	var r rune

	previewPp := pp
	pp = oj.skipWhitespace(pp)
	r, pp = oj.peekNextRune(pp)
	if r != '{' {
		err = fmt.Errorf("expected start of object at %s", oj.preview(previewPp))
		return
	}

	om := NewStringMap()
	for {
		pp = oj.skipWhitespace(pp)

		previewPp = pp
		r, pp = oj.peekNextRune(pp)
		if r == '}' {
			break
		}

		if om.Len() != 0 {
			if r != ',' {
				err = fmt.Errorf("missing comma in table at %s", oj.preview(previewPp))
				return
			}
			previewPp = pp
		}

		pp = oj.skipWhitespace(previewPp) // goes back a character if no } or , delimeter

		var key string
		key, pp = oj.peekString(pp)
		if pp.line == 0 {
			err = fmt.Errorf("missing map key in table at %s", oj.preview(previewPp))
			return
		}

		pp = oj.skipWhitespace(pp)
		r, pp = oj.peekNextRune(pp)
		if r != ':' {
			err = fmt.Errorf("expected colon after key at %s", oj.preview(previewPp))
		}

		var value any
		var terr error
		value, pp, terr = oj.peekJsonValue(pp)
		if terr != nil {
			err = terr
			return
		}
		if pp.line == 0 {
			err = fmt.Errorf("expected value at %s", oj.preview(previewPp))
			return
		}

		om.Set(key, value)
	}

	m = om
	nextPp = pp
	return
}

func (oj *orderedJson) peekJsonArray(pp parserPos) (a []any, nextPp parserPos, err error) {
	var r rune

	previewPp := pp
	pp = oj.skipWhitespace(pp)
	r, pp = oj.peekNextRune(pp)
	if r != '[' {
		err = fmt.Errorf("expected start of array at %s", oj.preview(previewPp))
		return
	}

	oa := []any{}
	for {
		pp = oj.skipWhitespace(pp)

		previewPp = pp
		r, pp = oj.peekNextRune(pp)
		if r == ']' {
			break
		}

		if len(oa) != 0 {
			if r != ',' {
				err = fmt.Errorf("missing comma in array at %s", oj.preview(previewPp))
				return
			}
			previewPp = pp
		}

		pp = oj.skipWhitespace(previewPp) // goes back a character if no ] or , delimeter

		var value any
		var terr error
		value, pp, terr = oj.peekJsonValue(pp)
		if terr != nil {
			err = terr
			return
		}
		if pp.line == 0 {
			err = fmt.Errorf("expected value at %s", oj.preview(previewPp))
			return
		}

		oa = append(oa, value)
	}

	a = oa
	nextPp = pp
	return
}

func (oj *orderedJson) peekJsonValue(pp parserPos) (val any, nextPp parserPos, err error) {
	pp = oj.skipWhitespace(pp)

	var s string
	s, nextPp = oj.peekString(pp)
	if nextPp.line != 0 {
		val = s
		return
	}

	var v any
	v, nextPp = oj.peekValueKeyword(pp)
	if nextPp.line != 0 {
		val = v
		return
	}

	var f float64
	f, nextPp = oj.peekFloat(pp)
	if nextPp.line != 0 {
		val = f
		return
	}

	om, nextPp, terr := oj.peekJsonMap(pp)
	if terr == nil {
		if nextPp.line != 0 {
			val = om
			return
		}
	}

	oa, nextPp, terr := oj.peekJsonArray(pp)
	if terr == nil {
		if nextPp.line != 0 {
			val = oa
			return
		}
	}

	err = ErrInvaldJsonValue
	return
}
