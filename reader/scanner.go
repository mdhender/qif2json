/*
 *  qif2json - a QIF data conversion utility
 *
 *  Copyright (c) 2021 Michael D Henderson
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package reader

//import (
//	"bytes"
//	"unicode"
//	"unicode/utf8"
//)
//
//type Scanner struct {
//	Line   int
//	Col    int
//	Buffer []byte
//}
//
//// New returns a new Scanner.
//// The Scanner will use the supplied buffer.
//func New(line, col int, b []byte) Scanner {
//	return Scanner{
//		Line:   line,
//		Col:    col,
//		Buffer: bdup(b),
//	}
//}
//
//// New returns a new Scanner.
//// The Scanner will use a copy of the supplied buffer.
//func NewWithCopy(line, col int, b []byte) Scanner {
//	return Scanner{
//		Line:   line,
//		Col:    col,
//		Buffer: bdup(b),
//	}
//}
//
//// AccountType will accept all the text up to the end of line.
//func (s Scanner) AccountType() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Amount will accept all the text up to the end of line.
//func (s Scanner) Amount() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Bang will accept '!'.
//func (s Scanner) Bang() (int, Scanner) {
//	if len(s.Buffer) == 0 || s.Buffer[0] != '!' {
//		return 0, s
//	}
//	s.Buffer, s.Col = s.Buffer[1:], s.Col+1
//	return 1, s
//}
//
//// Colon will accept ':'.
//func (s Scanner) Colon() (int, Scanner) {
//	if !bytes.HasPrefix(s.Buffer, []byte{':'}) {
//		return 0, s
//	}
//	s.Buffer, s.Col = s.Buffer[1:], s.Col+1
//	return 1, s
//}
//
//// Date will accept a date string which looks like
////    digit digit? slash (space digit) digit tic digit digit
//func (s Scanner) Date() ([]byte, Scanner) {
//	var length, w int
//	var r rune
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); !unicode.IsDigit(r) { // digit
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); unicode.IsDigit(r) { // digit?
//		length += w
//	}
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); r != '/' { // slash
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); !(r == ' ' || unicode.IsDigit(r)) { // (space digit)
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); !unicode.IsDigit(r) { // digit
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); r != '\'' { // tic
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); !unicode.IsDigit(r) { // digit
//		return nil, s
//	}
//	length += w
//
//	if r, w = utf8.DecodeRune(s.Buffer[length:]); !unicode.IsDigit(r) { // digit
//		return nil, s
//	}
//	length += w
//
//	text := s.Buffer[:length]
//	s.Buffer = s.Buffer[length:]
//	return text, s
//}
//
//// Description will accept all the text up to the end of line.
//func (s Scanner) Description() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// EOL will accept \r\n and \n.
//func (s Scanner) EOL() ([]byte, Scanner) {
//	if len(s.Buffer) == 0 {
//		return nil, s
//	}
//	if s.Buffer[0] == '\n' {
//		s.Buffer, s.Line, s.Col = s.Buffer[1:], s.Line+1, 1
//		return []byte{'\n'}, s
//	}
//	if len(s.Buffer) > 1 && s.Buffer[0] == '\r' && s.Buffer[1] == '\n' {
//		s.Buffer, s.Line, s.Col = s.Buffer[2:], s.Line+1, 1
//		return []byte{'\n'}, s
//	}
//	return nil, s
//}
//
//// EOS will accept '!' or end-of-input.
//func (s Scanner) EOS() ([]byte, Scanner) {
//	if len(s.Buffer) == 0 || s.Buffer[0] == '!' {
//		return []byte{'!'}, s
//	}
//	return nil, s
//}
//
//// Label will accept all the text up to the end of line.
//func (s Scanner) Label() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Literal will accept a literal. The scanner's line and col
//// variables will be hosed if the literal has an embedded newline.
//func (s Scanner) Literal(lit []byte) ([]byte, Scanner) {
//	if !bytes.HasPrefix(s.Buffer, lit) {
//		return nil, s
//	}
//	s.Buffer, s.Col = s.Buffer[len(lit):], s.Col+len(lit)
//	return lit, s
//}
//
//// Name will accept all the text up to the end of line.
//func (s Scanner) Name() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Position returns the current line and column
//func (s Scanner) Position() (line, col int) {
//	return s.Line, s.Col
//}
//
//// Risk will accept all the text up to the end of line.
//func (s Scanner) Risk() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// SecurityType will accept all the text up to the end of line.
//func (s Scanner) SecurityType() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Spaces will accept any run of spaces, including new-lines.
//func (s Scanner) Spaces() ([]byte, Scanner) {
//	var length int
//	for r, w := utf8.DecodeRune(s.Buffer); unicode.IsSpace(r); r, w = utf8.DecodeRune(s.Buffer[length:]) {
//		if r == '\n' {
//			s.Line, s.Col = s.Line+1, 0
//		}
//		s.Col, length = s.Col+1, length+w
//	}
//	spaces := s.Buffer[:length]
//	s.Buffer = s.Buffer[length:]
//	return spaces, s
//}
//
//// Text will accept all the text up to the end of line.
//func (s Scanner) Text() ([]byte, Scanner) {
//	var length int
//	for r, w := utf8.DecodeRune(s.Buffer); r != utf8.RuneError && r != '\n'; r, w = utf8.DecodeRune(s.Buffer[length:]) {
//		s.Col, length = s.Col+1, length+w
//	}
//	if length == 0 {
//		return nil, s
//	}
//	text := s.Buffer[:length]
//	s.Buffer = s.Buffer[length:]
//	return text, s
//}
//
//// Ticker will accept all the text up to the end of line.
//func (s Scanner) Ticker() ([]byte, Scanner) {
//	return s.Text()
//}
//
//// Type will accept 'Type'.
//func (s Scanner) Type() (int, Scanner) {
//	if !bytes.HasPrefix(s.Buffer, []byte{'T', 'y', 'p', 'e'}) {
//		return 0, s
//	}
//	s.Buffer, s.Col = s.Buffer[4:], s.Col+4
//	return 1, s
//}
