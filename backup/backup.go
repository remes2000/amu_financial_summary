package backup

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/account_transaction"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Backup struct {
	Categories   []category.CategoryBackup                      `json:"categories" binding:"dive"`
	Regexps      []regexp.RegexpBackup                          `json:"regexps" binding:"dive"`
	Transactions []account_transaction.AccountTransactionBackup `json:"transactions" binding:"dive"`
}

type Summary struct {
	ImportedCategories   int64 `json:"importedCategories"`
	ImportedRegexps      int64 `json:"importedRegexps"`
	ImportedTransactions int64 `json:"importedTransactions"`
}

func ExportData(backup *Backup) error {
	var categories []category.Category
	var regexps []regexp.Regexp
	var transactions []account_transaction.AccountTransaction
	if err := global.Database.Find(&categories).Error; err != nil {
		return err
	}
	if err := global.Database.Find(&regexps).Error; err != nil {
		return err
	}
	if err := global.Database.Find(&transactions).Error; err != nil {
		return err
	}
	for _, categoryToExport := range categories {
		categoryBackup := category.CategoryBackup{}
		categoryBackup.FromCategory(categoryToExport)
		backup.Categories = append(backup.Categories, categoryBackup)
	}
	for _, regexpToExport := range regexps {
		regexpBackup := regexp.RegexpBackup{}
		regexpBackup.FromRegexp(regexpToExport)
		backup.Regexps = append(backup.Regexps, regexpBackup)
	}
	for _, transactionToExport := range transactions {
		transactionBackup := account_transaction.AccountTransactionBackup{}
		transactionBackup.FromAccountTransaction(transactionToExport)
		backup.Transactions = append(backup.Transactions, transactionBackup)
	}
	return nil
}

func ImportData(backup *Backup) (Summary, error) {
	var summary Summary
	var categoriesToImport []category.Category
	var regexpsToImport []regexp.Regexp
	var transactionsToImport []account_transaction.AccountTransaction
	for _, categoryToImport := range backup.Categories {
		categoriesToImport = append(categoriesToImport, categoryToImport.ToCategory())
	}
	for _, regexpToImport := range backup.Regexps {
		regexpsToImport = append(regexpsToImport, regexpToImport.ToRegexp())
	}
	for _, transactionToImport := range backup.Transactions {
		transactionsToImport = append(transactionsToImport, transactionToImport.ToAccountTransaction())
	}
	err := global.Database.Transaction(func(tx *gorm.DB) error {
		categoriesImportResult := tx.Create(&categoriesToImport)
		if categoriesImportResult.Error != nil {
			return categoriesImportResult.Error
		}
		regexpsImportResult := tx.Create(&regexpsToImport)
		if regexpsImportResult.Error != nil {
			return regexpsImportResult.Error
		}
		transactionsImportResult := tx.Create(&transactionsToImport)
		if transactionsImportResult.Error != nil {
			return transactionsImportResult.Error
		}
		summary.ImportedCategories = categoriesImportResult.RowsAffected
		summary.ImportedRegexps = regexpsImportResult.RowsAffected
		summary.ImportedTransactions = transactionsImportResult.RowsAffected
		return nil
	})
	return summary, err
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "backup"
	rest.GET(controllerName+"/export", exportData)
	rest.POST(controllerName+"/import", importData)
}

func exportData(context *gin.Context) {
	var backup Backup
	if err := ExportData(&backup); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, backup)
}

func importData(context *gin.Context) {
	var backup Backup
	if err := context.BindJSON(&backup); err != nil {
		return
	}
	importResult, err := ImportData(&backup)
	if err != nil {
		log.Print(err)
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Data import failed"})
			return
		}
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, importResult)
}
