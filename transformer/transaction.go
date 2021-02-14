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

package transformer

import "github.com/mdhender/qif2json/reader/transaction"

type Transaction struct {
	Line          int
	Type          string
	Account       string
	Address       []string // Up to five lines (the sixth line is an optional message)
	Category      string
	ClearedStatus string
	Commission    string
	Date          string
	Interest      string
	Memo          string
	MemorizedFlag string
	Quantity      string
	Payee         string
	Price         string
	RefNo         string
	Split         []*Split
	Ticker        string
}

type Split struct {
	Line     int
	Account  string
	Amount   string
	Category string
	Memo     string
}

func NormalizeSplits(transactions []*transaction.Record) []*Transaction {
	var normalized []*Transaction
	for _, t := range transactions {
		xact := Transaction{
			Line:          t.Line,
			Type:          t.Type,
			Date:          t.Date,
			Account:       t.Account,
			ClearedStatus: t.ClearedStatus,
			Memo:          t.Memo,
			Payee:         t.Payee,
			RefNo:         t.RefNo,
		}
		if len(t.Split) == 0 {
			xact.Memo = ""
			split := Split{
				Line:     t.Line,
				Account:  t.ToAccount,
				Amount:   t.AmountTCode,
				Category: t.Category,
				Memo:     t.Memo,
			}
			xact.Split = append(xact.Split, &split)
		} else {
			for i, line := range t.Split {
				split := Split{
					Line:     line.Line,
					Account:  line.Account,
					Amount:   line.Amount,
					Category: line.Category,
					Memo:     line.Memo,
				}
				if i == 0 && split.Account == "" {
					split.Account = t.ToAccount
				}
				xact.Split = append(xact.Split, &split)
			}
		}
		normalized = append(normalized, &xact)
	}
	return normalized
}
