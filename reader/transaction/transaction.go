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

package transaction

import (
	"fmt"
	"github.com/mdhender/qif2json/buffer"
	"github.com/mdhender/qif2json/stdlib"
	"strings"
)

type Section struct {
	Line    int       `json:"-"`
	Col     int       `json:"-"`
	Records []*Record `json:"records,omitempty"`
}

type Record struct {
	Line          int
	Col           int
	Account       string
	Address       []string // Up to five lines (the sixth line is an optional message)
	AmountTCode   string
	AmountUCode   string
	BudgetAmount  []string
	Category      string // Category/Subcategory/Transfer/Class
	ClearedStatus string
	Commission    string
	Date          string
	Interest      string
	Memo          string
	MemorizedFlag string
	Quantity      string
	Payee         string
	Price         string
	RefNo         string // (check or reference number)
	Split         []*Split
	Ticker        string
	ToAccount     string // if category is [xxxx], then ToAccount is 'xxxx'
	Type          string
}

type Split struct {
	Line     int    `json:"-"`
	Col      int    `json:"-"`
	Account  string `json:"account,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Category string `json:"category,omitempty"`
	Memo     string `json:"memo,omitempty"`
}

func ReadSection(buf buffer.Buffer, account, accountType string) (*Section, buffer.Buffer, error) {
	saved, sname, section := buf, "transactions", Section{Line: buf.Line, Col: buf.Col}

	var literal string
	switch accountType {
	case "Bank", "Cash", "CCard", "Invst", "Oth A", "Oth L", "Memorized", "Prices":
		literal = "!Type:" + accountType
	default:
		panic(fmt.Sprintf("assert(account.type != %q)", accountType))
	}
	lit, bb := buf.Literal(literal)
	if lit == nil {
		return nil, saved, nil
	}
	buf = bb

	// read the section detail
	var err error
	for {
		var record *Record
		record, buf, err = ReadRecord(buf, account, accountType)
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
		fmt.Printf("error here %q\n", string(buf.Buffer[:20]))
		return nil, saved, fmt.Errorf("%d: %s: %d:%d: unexpected input", section.Line, sname, buf.Line, buf.Col)
	}
	buf = bb

	return &section, buf, nil
}

func ReadRecord(buf buffer.Buffer, account, accountType string) (*Record, buffer.Buffer, error) {
	saved, sname, record := buf, "transaction", Record{Line: buf.Line, Col: buf.Col, Account: account, Type: accountType}

	var found bool
	var category, cleared, commission, date, interest, memo, memorized, payee, qty, refNo, ticker, tcode, toAccount, ucode []byte
	var split *Split
	for {
		if addrLine, bb := buf.Field("A"); addrLine != nil {
			found, record.Address = true, append(record.Address, string(addrLine))
			buf = bb
			continue
		}
		if cleared == nil {
			if cleared, buf = buf.Field("C"); cleared != nil {
				found, record.ClearedStatus = true, string(cleared)
				continue
			}
		}
		if commission == nil {
			if commission, buf = buf.Field("O"); commission != nil {
				found, record.Commission = true, string(commission)
				continue
			}
		}
		if date == nil {
			if date, buf = buf.Date("D"); date != nil {
				found, record.Date = true, string(date)
				continue
			}
		}
		if interest == nil {
			if interest, buf = buf.Field("I"); interest != nil {
				found, record.Interest = true, string(interest)
				continue
			}
		}
		if memo == nil {
			if memo, buf = buf.Field("M"); memo != nil {
				found, record.Memo = true, string(memo)
				continue
			}
		}
		if memorized == nil {
			if memorized, buf = buf.Field("K"); memorized != nil {
				found, record.MemorizedFlag = true, string(memorized)
				continue
			}
		}
		if payee == nil {
			if payee, buf = buf.Field("P"); payee != nil {
				found, record.Payee = true, string(payee)
				continue
			}
		}
		if price, bb := buf.Field("\""); price != nil {
			lexeme := strings.TrimRight(string(price), "\"")
			if fields := strings.Split(lexeme, "\""); len(fields) == 3 {
				found, record.Ticker = true, fields[0]
				record.Price = strings.ReplaceAll(fields[1], ",", "")
				date = []byte(fields[2])
				record.Date = stdlib.Date(date)
				buf = bb
				continue
			}
		}
		if qty == nil {
			if qty, buf = buf.Field("Q"); qty != nil {
				found, record.Quantity = true, string(qty)
				continue
			}
		}
		if refNo == nil {
			if refNo, buf = buf.Field("N"); refNo != nil {
				found, record.RefNo = true, string(refNo)
				continue
			}
		}
		if splitAmount, bb := buf.Field("$"); splitAmount != nil {
			if split == nil {
				split = &Split{Line: bb.Line}
				record.Split = append(record.Split, split)
			}
			found, split.Amount = true, string(splitAmount)
			buf = bb
			continue
		}
		if splitCategory, bb := buf.Field("S"); splitCategory != nil {
			split = &Split{Line: bb.Line}
			found, record.Split = true, append(record.Split, split)
			split.Category = string(splitCategory)
			if strings.HasPrefix(split.Category, "[") {
				split.Account = strings.Trim(split.Category, "[]")
				split.Category = ""
			}
			buf = bb
			continue
		}
		if splitMemo, bb := buf.Field("E"); splitMemo != nil {
			if split == nil {
				split = &Split{Line: bb.Line}
				record.Split = append(record.Split, split)
			}
			found, split.Memo = true, string(splitMemo)
			buf = bb
			continue
		}
		if tcode == nil {
			if tcode, buf = buf.Field("T"); tcode != nil {
				found, record.AmountTCode = true, string(tcode)
				continue
			}
		}
		if ticker == nil {
			if ticker, buf = buf.Field("Y"); ticker != nil {
				found, record.Ticker = true, string(ticker)
				continue
			}
		}
		if ucode == nil {
			if ucode, buf = buf.Field("U"); ucode != nil {
				found, record.AmountUCode = true, string(ucode)
				continue
			}
		}

		// category must follow toAccount since they share a common prefix
		if toAccount == nil {
			if toAccount, buf = buf.Field("L["); toAccount != nil {
				found, record.ToAccount = true, strings.TrimRight(string(toAccount), "]")
				continue
			}
		}
		if category == nil {
			if category, buf = buf.Field("L"); category != nil {
				found, record.Category = true, string(category)
				continue
			}
		}

		// add budget amount logic here
		if record.Type == "Memorized" {
			var budgetAmount []byte
			for _, flag := range []string{"1", "2", "3", "4", "5", "6", "7"} {
				if budgetAmount == nil {
					budgetAmount, buf = buf.Field(flag)
				}
			}
			if budgetAmount != nil {
				found, record.BudgetAmount = true, append(record.BudgetAmount, string(budgetAmount))
				continue
			}
		}

		break
	}

	if !found { // no fields found
		return nil, saved, nil
	}

	// check for required fields
	switch record.Type {
	case "Memorized":
		if memorized == nil {
			return nil, saved, fmt.Errorf("%d: %s: missing field %q", record.Line, sname, "memorized")
		}
	default:
		if date == nil {
			return nil, saved, fmt.Errorf("%d: %s: missing field %q", record.Line, sname, "date")
		}
	}

	eor, bb := buf.EndOfRecord()
	if eor == nil {
		return nil, saved, fmt.Errorf("%d: %s: missing record terminator", buf.Line, sname)
	}
	buf = bb

	return &record, buf, nil
}
