package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType は取引の種類（収入 or 支出）を表します。
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// 支出カテゴリ
const (
	// 固定費
	ExpenseCatRent         = "家賃"
	ExpenseCatMedicalLoan  = "医療ローン"
	ExpenseCatInsurance    = "保険"
	ExpenseCatSubscription = "サブスク"
	ExpenseCatInstallment  = "ショッピング分割"
	ExpenseCatUtilities    = "光熱費"
	// 変動費
	ExpenseCatFood             = "食費"
	ExpenseCatDailyNecessities = "日用品"
	ExpenseCatTransport        = "交通費"
	ExpenseCatSkillUp          = "スキルアップ"
	ExpenseCatWorkSupplies     = "仕事用品"
	ExpenseCatMedical          = "医療"
	ExpenseCatBeauty           = "美容"
	ExpenseCatEntertainment    = "娯楽"
	ExpenseCatOther            = "その他（支出）"
)

// 収入カテゴリ
const (
	IncomeCatSalary = "給与"
	IncomeCatBonus  = "賞与"
	IncomeCatSales  = "販売"
	IncomeCatOther  = "その他（収入）"
)

// AllExpenseCategories は支出カテゴリのリストです。
var AllExpenseCategories = []string{
	ExpenseCatRent, ExpenseCatMedicalLoan, ExpenseCatInsurance, ExpenseCatSubscription,
	ExpenseCatInstallment, ExpenseCatUtilities, ExpenseCatFood, ExpenseCatDailyNecessities,
	ExpenseCatTransport, ExpenseCatSkillUp, ExpenseCatWorkSupplies, ExpenseCatMedical,
	ExpenseCatBeauty, ExpenseCatEntertainment, ExpenseCatOther,
}

// AllIncomeCategories は収入カテゴリのリストです。
var AllIncomeCategories = []string{
	IncomeCatSalary, IncomeCatBonus, IncomeCatSales, IncomeCatOther,
}

// IsValidCategory は指定されたカテゴリが取引タイプに対して有効かどうかをチェックします。
func IsValidCategory(transactionType TransactionType, category string) bool {
	var validCategories []string
	if transactionType == Income {
		validCategories = AllIncomeCategories
	} else if transactionType == Expense {
		validCategories = AllExpenseCategories
	} else {
		return false
	}

	for _, validCat := range validCategories {
		if category == validCat {
			return true
		}
	}
	return false
}

// Transaction はデータベースの取引レコードとAPIのJSONボディを表します。
type Transaction struct {
	gorm.Model
	Date     time.Time       `json:"date" binding:"required"`
	Type     TransactionType `json:"type" binding:"required,oneof=income expense"`
	Category string          `json:"category" binding:"required"`
	Amount   int             `json:"amount" binding:"required,gt=0"`
	Memo     string          `json:"memo"`
}
