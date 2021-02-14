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

package security

import (
	"fmt"
	"github.com/mdhender/qif2json/buffer"
)

type Section struct {
	Line    int       `json:"-"`
	Col     int       `json:"-"`
	Records []*Record `json:"records,omitempty"`
}

type Record struct {
	Line        int    `json:"-"`
	Col         int    `json:"-"`
	Description string `json:"descr,omitempty"`
	Name        string `json:"name"`
	Risk        string `json:"risk,omitempty"`
	Ticker      string `json:"ticker"`
	Type        string `json:"type"`
}

func ReadSection(buf buffer.Buffer) (*Section, buffer.Buffer, error) {
	saved, sname, section := buf, "securities", Section{Line: buf.Line, Col: buf.Col}

	lit, bb := buf.Literal("!Type:Security")
	if lit == nil {
		return nil, saved, nil
	}
	buf = bb

	// read the section detail
	var err error
	for {
		var record *Record
		record, buf, err = ReadRecord(buf)
		if err != nil {
			return nil, buf, fmt.Errorf("%d: %s: %w", section.Line, sname, err)
		} else if record == nil {
			break
		}
		section.Records = append(section.Records, record)
	}

	// read the end of section marker
	eos, bb := buf.EndOfSection()
	if eos == nil {
		return nil, saved, fmt.Errorf("%d: %s: %d:%d: unexpected input", section.Line, sname, buf.Line, buf.Col)
	}
	buf = bb

	return &section, buf, nil
}

func ReadRecord(buf buffer.Buffer) (*Record, buffer.Buffer, error) {
	saved, sname, record := buf, "security", Record{Line: buf.Line, Col: buf.Col}

	var found bool
	var descr, name, risk, ticker, typ []byte
	for {
		if descr == nil {
			if descr, buf = buf.Field("D"); descr != nil {
				found, record.Description = true, string(descr)
				continue
			}
		}
		if name == nil {
			if name, buf = buf.Field("N"); name != nil {
				found, record.Name = true, string(name)
				continue
			}
		}
		if risk == nil {
			if risk, buf = buf.Field("G"); risk != nil {
				found, record.Risk = true, string(risk)
				continue
			}
		}
		if ticker == nil {
			if ticker, buf = buf.Field("S"); ticker != nil {
				found, record.Ticker = true, string(ticker)
				continue
			}
		}
		if typ == nil {
			if typ, buf = buf.Field("T"); typ != nil {
				found, record.Type = true, string(typ)
				continue
			}
		}

		break
	}

	if !found { // no fields found
		return nil, saved, nil
	}

	// check for required fields
	if name == nil {
		return nil, saved, fmt.Errorf("%d: %s: missing field %q", record.Line, sname, "name")
	}

	eor, bb := buf.EndOfRecord()
	if eor == nil {
		return nil, saved, fmt.Errorf("%d: %s: missing record terminator", buf.Line, sname)
	}
	buf = bb

	return &record, buf, nil
}
