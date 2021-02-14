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

import (
	"fmt"
	"github.com/mdhender/qif2json/buffer"
	"github.com/mdhender/qif2json/reader/account"
	"github.com/mdhender/qif2json/reader/category"
	"github.com/mdhender/qif2json/reader/security"
	"github.com/mdhender/qif2json/reader/tag"
	"github.com/mdhender/qif2json/reader/transaction"
)

type Reader struct {
	active struct {
		account     string
		accountType string
	}
	Accounts     *account.Section      `json:"accounts,omitempty"`
	Categories   *category.Section     `json:"categories,omitempty"`
	Securities   *security.Section     `json:"securities,omitempty"`
	Tags         *tag.Section          `json:"tags,omitempty"`
	Transactions []*transaction.Record `json:"transactions,omitempty"`
	Memorized    []*transaction.Record `json:"-"`
	Prices       []*transaction.Record `json:"-"`
}

func Read(buf buffer.Buffer) (*Reader, error) {
	var r Reader
	for len(buf.Buffer) != 0 {
		if literal, bb := buf.Literal("!Clear:AutoSwitch"); literal != nil {
			// ignore
			buf = bb
			continue
		}
		if literal, bb := buf.Literal("!Option:AutoSwitch"); literal != nil {
			// ignore
			buf = bb
			continue
		}
		if section, bb, err := account.ReadSection(buf); err != nil {
			return nil, err
		} else if section != nil {
			if len(section.Records) != 0 {
				if r.Accounts == nil {
					r.Accounts = section
				} else if len(section.Records) == 1 {
					r.active.account = section.Records[0].Name
					r.active.accountType = section.Records[0].Type
				} else {
					panic("!")
				}
			}
			buf = bb
			continue
		}
		if section, bb, err := category.ReadSection(buf); err != nil {
			return nil, err
		} else if section != nil {
			if len(section.Records) != 0 {
				if r.Categories == nil {
					r.Categories = section
				} else {
					panic("!")
				}
			}
			buf = bb
			continue
		}
		if section, bb, err := security.ReadSection(buf); err != nil {
			return nil, err
		} else if section != nil {
			if len(section.Records) != 0 {
				if r.Securities == nil {
					r.Securities = &security.Section{
						Line: buf.Line,
						Col:  buf.Col,
					}
				}
				r.Securities.Records = append(r.Securities.Records, section.Records...)
			}
			buf = bb
			continue
		}
		if section, bb, err := tag.ReadSection(buf); err != nil {
			return nil, err
		} else if section != nil {
			if len(section.Records) != 0 {
				if r.Tags == nil {
					r.Tags = section
				} else {
					panic("!")
				}
			}
			buf = bb
			continue
		}
		if section, bb, err := transaction.ReadSection(buf, r.active.account, r.active.accountType); err != nil {
			return nil, err
		} else if section != nil {
			for _, xact := range section.Records {
				r.Transactions = append(r.Transactions, xact)
			}
			buf = bb
			continue
		}
		if section, bb, err := transaction.ReadSection(buf, "", "Memorized"); err != nil {
			return nil, err
		} else if section != nil {
			for _, xact := range section.Records {
				r.Memorized = append(r.Memorized, xact)
			}
			buf = bb
			continue
		}
		if section, bb, err := transaction.ReadSection(buf, "", "Prices"); err != nil {
			return nil, err
		} else if section != nil {
			for _, xact := range section.Records {
				r.Prices = append(r.Prices, xact)
			}
			buf = bb
			continue
		}
		return nil, fmt.Errorf("%d:%d: unexpected input", buf.Line, buf.Col)
	}
	return &r, nil
}
