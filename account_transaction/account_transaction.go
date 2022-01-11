package account_transaction

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/common"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/validators"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type AccountTransaction struct {
	Id         uint               `json:"id" gorm:"autoIncrement;unique;notNull"`
	Title      string             `json:"title" binding:"required" gorm:"primaryKey"`
	Date       time.Time          `json:"date" binding:"required" gorm:"primaryKey"`
	Amount     int                `json:"amount" binding:"required" gorm:"notNull"`
	CategoryId *uint              `json:"-"`
	Category   *category.Category `json:"category"`
}

func (t AccountTransaction) GetUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"amount":      t.Amount,
		"category_id": t.CategoryId,
	}
}

func (t *AccountTransaction) SetDate(date string) {
	time, _ := time.Parse(validators.ValidDateLayout, date)
	t.Date = time
}

func (t *AccountTransaction) SetAmount(amount float64) {
	t.Amount = int(amount * 100)
}

func (t *AccountTransaction) SetCategory(categories []category.Category) {
	for _, category := range categories {
		if category.Matches(t.Title) {
			t.Category = &category
			t.CategoryId = &category.Id
			return
		}
	}
	t.Category = nil
	t.CategoryId = nil
}

type AccountTransactionRequest struct {
	Date   string  `json:"date" binding:"required,validdate"`
	Title  string  `json:"title" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (r AccountTransactionRequest) GetAccountTransaction(categories []category.Category) AccountTransaction {
	var transaction AccountTransaction
	transaction.Title = r.Title
	transaction.SetAmount(r.Amount)
	transaction.SetDate(r.Date)
	transaction.SetCategory(categories)
	return transaction
}

type ForceSetCategoryRequest struct {
	TransactionId uint `json:"transactionId" binding:"required"`
	CategoryId    uint `json:"categoryId" binding:"required"`
}

func ImportTransactions(transactionsToImport []AccountTransactionRequest) ([]uint, error) {
	var transactionsWithoutCategory []uint
	err := global.Database.Transaction(func(tx *gorm.DB) error {
		var categories []category.Category
		if err := category.GetAllCategories(&categories); err != nil {
			return err
		}

		for _, transactionToImport := range transactionsToImport {
			var transactionFromDb AccountTransaction
			transaction := transactionToImport.GetAccountTransaction(categories)
			if err := tx.Where("title = ? AND date = ?", transaction.Title, transaction.Date).First(&transactionFromDb).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if transactionFromDb.Id == 0 {
				if err := tx.Create(&transaction).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&transactionFromDb).Updates(transaction.GetUpdateMap()).Error; err != nil {
					return err
				}
			}

			if transaction.CategoryId == nil {
				if transactionFromDb.Id == 0 {
					transactionsWithoutCategory = append(transactionsWithoutCategory, transaction.Id)
				} else {
					transactionsWithoutCategory = append(transactionsWithoutCategory, transactionFromDb.Id)
				}
			}
		}
		return nil
	})
	return transactionsWithoutCategory, err
}

func GetAccountTransactionById(transaction *AccountTransaction, id uint) error {
	if err := global.Database.Preload("Category.Regexps").Where("id = ?", id).First(transaction).Error; err != nil {
		return err
	}
	return nil
}

func ForceSetCategory(transaction *AccountTransaction, category *category.Category) error {
	if err := global.Database.Model(transaction).Update("category_id", category.Id).Error; err != nil {
		return err
	}
	return nil
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "account-transaction"
	rest.POST(controllerName, importTransactions)
	rest.GET(controllerName+"/:id", getOne)
	rest.POST(controllerName+"/force-set-category", forceSetCategory)
}

func importTransactions(context *gin.Context) {
	var transactionsToImport []AccountTransactionRequest
	if err := context.BindJSON(&transactionsToImport); err != nil {
		return
	}
	if transactionsToImport == nil || len(transactionsToImport) == 0 {
		context.Status(http.StatusBadRequest)
		return
	}
	noCategoryTransactionsIdList, err := ImportTransactions(transactionsToImport)
	if err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, gin.H{"noCategoryTransactions": noCategoryTransactionsIdList})
}

func getOne(context *gin.Context) {
	var requestedTransaction AccountTransaction
	var idUri common.IdUri

	if err := context.ShouldBindUri(&idUri); err != nil {
		return
	}
	if err := GetAccountTransactionById(&requestedTransaction, idUri.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Cannot find entity with id %d", idUri.Id)})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, requestedTransaction)
}

func forceSetCategory(context *gin.Context) {
	var forceSetCategoryRequest ForceSetCategoryRequest
	var transaction AccountTransaction
	var requestedCategory category.Category

	if err := context.ShouldBindJSON(&forceSetCategoryRequest); err != nil {
		return
	}
	if err := GetAccountTransactionById(&transaction, forceSetCategoryRequest.TransactionId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Transaction with id %d does not exist", forceSetCategoryRequest.TransactionId)})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := category.GetCategoryById(&requestedCategory, forceSetCategoryRequest.CategoryId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Category with id %d does not exist", forceSetCategoryRequest.CategoryId)})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := ForceSetCategory(&transaction, &requestedCategory); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := GetAccountTransactionById(&transaction, transaction.Id); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, transaction)
}
