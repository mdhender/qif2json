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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mdhender/qif2json/buffer"
	"github.com/mdhender/qif2json/reader"
	"github.com/mdhender/qif2json/transformer"
	"github.com/peterbourgon/ff/v3"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	fs := flag.NewFlagSet("my-program", flag.ExitOnError)
	var (
		input = fs.String("input", "", "QIF file to translate")
		accts = fs.String("accounts", "", "file to write accounts to")
		cats  = fs.String("categories", "", "file to write categories to")
		trans = fs.String("transactions", "", "file to write transactions to")
		_     = fs.String("config", "", "config file (optional)")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("QIFXLAT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser)); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	if *input == "" {
		fmt.Printf("please provide the name of the QIF file to translate\n")
		os.Exit(2)
	}
	fmt.Printf("%-30s == %q\n", "QIFXLAT_INPUT", *input)
	if *accts != "" {
		fmt.Printf("%-30s == %q\n", "QIFXLAT_ACCOUNTS", *accts)
	}
	if *cats != "" {
		fmt.Printf("%-30s == %q\n", "QIFXLAT_CATEGORIES", *cats)
	}
	if *trans != "" {
		fmt.Printf("%-30s == %q\n", "QIFXLAT_TRANSACTIONS", *trans)
	}

	if err := run(*input, *accts, *cats, *trans); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

func run(name, accounts, categories, transactions string) error {
	started := time.Now()

	input, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	buf, err := buffer.NewBuffer(input)
	if err != nil {
		return err
	}
	r, err := reader.Read(buf)
	if err != nil {
		return err
	}
	if accounts != "" {
		type Account struct {
			Type                 string `json:"type"`
			Name                 string `json:"name"`
			CreditLimit          string `json:"credit_limit,omitempty"`
			Description          string `json:"descr,omitempty"`
			StatementBalance     string `json:"balance,omitempty"`
			StatementBalanceDate string `json:"statement_date,omitempty"`
		}
		var data struct {
			Accounts []Account `json:"accounts"`
		}
		for _, account := range r.Accounts.Records {
			var typ string
			switch account.Type {
			case "Bank":
				typ = "bank"
			case "CCard":
				typ = "creditCard"
			case "Cash":
				typ = "cash"
			case "Oth A":
				typ = "asset"
			case "Oth L":
				typ = "liability"
			case "Port":
				typ = "brokerage"
			case "401(k)/403(b)":
				typ = "retirement"
			default:
				panic(fmt.Sprintf("assert(account.type != %q)", account.Type))
			}
			data.Accounts = append(data.Accounts, Account{
				Type:                 typ,
				Name:                 account.Name,
				CreditLimit:          account.CreditLimit,
				Description:          account.Description,
				StatementBalance:     account.StatementBalance,
				StatementBalanceDate: account.StatementBalanceDate,
			})
		}
		if err := writeJSON(accounts, data); err != nil {
			return err
		}
	}

	if categories != "" {
		type Category struct {
			Name        string `json:"name"`
			Description string `json:"descr,omitempty"`
			Income      bool   `json:"income,omitempty"`
			TaxRelated  bool   `json:"tax_related,omitempty"`
			TaxSchedule string `json:"tax_schedule,omitempty"`
		}
		var data struct {
			Categories []Category `json:"categories"`
		}
		for _, category := range r.Categories.Records {
			data.Categories = append(data.Categories, Category{
				Name:        category.Name,
				Description: category.Description,
				Income:      category.IsIncome,
				TaxRelated:  category.IsTaxRelated,
				TaxSchedule: category.TaxSchedule,
			})
		}
		if err := writeJSON(categories, data); err != nil {
			return err
		}
	}

	if transactions != "" {
		normalized := transformer.NormalizeSplits(r.Transactions)

		type Split struct {
			Line     int    `json:"line,omitempty"`
			Account  string `json:"account,omitempty"`
			Amount   string `json:"amount,omitempty"`
			Category string `json:"category,omitempty"`
			Memo     string `json:"memo,omitempty"`
		}
		type Transaction struct {
			Line          int     `json:"line,omitempty"`
			Type          string  `json:"type,omitempty"`
			Date          string  `json:"date,omitempty"`
			Account       string  `json:"account,omitempty"`
			ToAccount     string  `json:"to_account,omitempty"`
			Amount        string  `json:"amount,omitempty"`
			Category      string  `json:"category,omitempty"`
			ClearedStatus string  `json:"cleared_status,omitempty"`
			Memo          string  `json:"memo,omitempty"`
			Payee         string  `json:"payee,omitempty"`
			RefNo         string  `json:"ref_no,omitempty"`
			Split         []Split `json:"lines,omitempty"`
		}
		var data struct {
			Transactions []Transaction `json:"transactions"`
		}
		for _, transaction := range normalized {
			xact := Transaction{
				Line:          transaction.Line,
				Account:       transaction.Account,
				ClearedStatus: transaction.ClearedStatus,
				Date:          transaction.Date,
				Memo:          transaction.Memo,
				Payee:         transaction.Payee,
				RefNo:         transaction.RefNo,
			}
			for _, line := range transaction.Split {
				split := Split{
					Line:     line.Line,
					Account:  line.Account,
					Amount:   line.Amount,
					Category: line.Category,
					Memo:     line.Memo,
				}
				xact.Split = append(xact.Split, split)
			}
			data.Transactions = append(data.Transactions, xact)
		}
		if err := writeJSON(transactions, data); err != nil {
			return err
		}
	}

	var totalRecords int
	if r.Accounts != nil {
		fmt.Printf("processed %8d accounts\n", len(r.Accounts.Records))
		totalRecords += len(r.Accounts.Records)
	}
	if r.Categories != nil {
		fmt.Printf("processed %8d categories\n", len(r.Categories.Records))
		totalRecords += len(r.Categories.Records)
	}
	totalRecords += len(r.Memorized)
	fmt.Printf("processed %8d memorized\n", len(r.Memorized))
	totalRecords += len(r.Prices)
	fmt.Printf("processed %8d prices\n", len(r.Prices))
	if r.Securities != nil {
		fmt.Printf("processed %8d securities\n", len(r.Securities.Records))
		totalRecords += len(r.Securities.Records)
	}
	if r.Tags != nil {
		fmt.Printf("processed %8d tags\n", len(r.Tags.Records))
	}
	totalRecords += len(r.Transactions)
	fmt.Printf("processed %8d transactions\n", len(r.Transactions))

	duration := time.Now().Sub(started)
	fmt.Printf("processed %8d records in %v\n", totalRecords, duration)

	return nil
}

func writeJSON(name string, data interface{}) error {
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(name, buf, 0600); err != nil {
		return err
	}
	fmt.Printf("wrote %8d bytes to %q\n", len(buf), name)
	return nil
}
