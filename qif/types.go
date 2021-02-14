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

// Package qif defines the types to be imported from the QIF file
package qif

// FILE is a structure that owns all of the data imported from the file.
type FILE struct {
	Accounts     []*ACCOUNT
	Transactions []*TRANSACTION
}

// ACCOUNT is information on a single account.
type ACCOUNT struct {
	Name        string
	Type        string
	CreditLimit int
	Descr       string
}

// TRANSACTION is information about a single transaction.
type TRANSACTION struct {
	ID      int
	Date    string
	Account string
	Payee   string `json:"payee,omitempty"`
	Memo    string `json:"memo,omitempty"`
	Lines   []*ENTRY
}

// ENTRY is information about a single entry in a transaction
type ENTRY struct {
	From        string   `json:"from,omitempty"`
	To          string   `json:"to,omitempty"`
	Category    string   `json:"category,omitempty"`
	Amount      int      `json:"amount,omitempty"`
	Memo        string   `json:"memo,omitempty"`
	FromAccount *ACCOUNT `json:"-"`
	ToAccount   *ACCOUNT `json:"-"`
}

// File contains the data imported from the QIF file.
type File struct {
	Account               AccountDetail                 `json:"-"`
	Accounts              []*AccountDetail              `json:",omitempty"`
	Banks                 []*BankDetail                 `json:"-"`
	Budget                []*BudgetDetail               `json:"-"`
	Cash                  []*CashDetail                 `json:"-"`
	Categories            []*CategoryDetail             `json:",omitempty"`
	CreditCards           []*CreditCardDetail           `json:"-"`
	Investments           []*InvestmentDetail           `json:"-"`
	MemorizedTransactions []*MemorizedTransactionDetail `json:"-"`
	OtherAssets           []*OtherAssetDetail           `json:"-"`
	OtherLiabilities      []*OtherLiabilityDetail       `json:"-"`
	Prices                []*PriceDetail                `json:"-"`
	Securities            []*SecurityDetail             `json:"-"`
	Tags                  []*TagDetail                  `json:",omitempty"`
}

// AccountDetail is
type AccountDetail struct {
	Name                 string
	Type                 string
	CreditLimit          int
	Descr                string
	StatementBalance     int
	StatementBalanceDate string
	Banks                []*BankDetail           `json:",omitempty"`
	Budget               []*BudgetDetail         `json:",omitempty"`
	Cash                 []*CashDetail           `json:",omitempty"`
	CreditCards          []*CreditCardDetail     `json:",omitempty"`
	Investments          []*InvestmentDetail     `json:",omitempty"`
	OtherAssets          []*OtherAssetDetail     `json:",omitempty"`
	OtherLiabilities     []*OtherLiabilityDetail `json:",omitempty"`
}

// BankDetail is
type BankDetail struct {
	Address       []string // Up to five lines (the sixth line is an optional message)
	AmountTCode   int
	AmountUCode   int
	Category      string // Category/Subcategory/Transfer/Class
	ClearedStatus string
	Date          string
	Memo          string
	Num           string // (check or reference number)
	Payee         string
	Split         []*Split
}

// BudgetDetail is
type BudgetDetail struct {
	Raw []string
}

// CashDetail is
type CashDetail struct {
	Raw []string
}

// CategoryDetail is
type CategoryDetail struct {
	Name        string // Category/Subcategory/Transfer/Class
	Descr       string
	IsExpense   bool
	TaxRelated  bool
	TaxSchedule string
}

// CreditCardDetail is
type CreditCardDetail struct {
	Address       []string // Up to five lines (the sixth line is an optional message)
	AmountTCode   int
	AmountUCode   int
	Category      string // Category/Subcategory/Transfer/Class
	ClearedStatus string
	Date          string
	Memo          string
	Num           string // (check or reference number)
	Payee         string
	Split         []*Split
}

// InvestmentDetail is
type InvestmentDetail struct {
	Raw []string
}

// MemorizedTransactionDetail is
type MemorizedTransactionDetail struct {
	Type                     string
	Address                  []string // Up to five lines (the sixth line is an optional message)
	AmountTCode              int
	AmountUCode              int
	Category                 string // Category/Subcategory/Transfer/Class
	ClearedStatus            string
	Date                     string
	Memo                     string
	MemorizedTransactionType string
	Num                      string // (check or reference number)
	Payee                    string
	Split                    []*Split
}

// OtherAssetDetail is
type OtherAssetDetail struct {
	Raw []string
}

// OtherLiabilityDetail is
type OtherLiabilityDetail struct {
	Raw []string
}

// PriceDetail is
type PriceDetail struct {
	Raw    []string
	Symbol string
	Price  string
	Date   string
}

// Split allows a detail line to be split into multiple transfers
type Split struct {
	Amount   int    // Dollar amount of split
	Category string // Category in split (Category/Transfer/Class)
	Memo     string // in split
}

// SecurityDetail is
type SecurityDetail struct {
	Description string
	Label       string
	Risk        string
	Symbol      string
	Type        string
}

// TagDetail is
type TagDetail struct {
	Description string
	Label       string
}
