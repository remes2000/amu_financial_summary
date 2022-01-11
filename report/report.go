package report

import (
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/account_transaction"
	"github.com/remes2000/amu_financial_summary/category"
	"log"
	"net/http"
)

type Report struct {
	Year         uint     `json:"year"`
	Month        uint     `json:"month"`
	TotalIncome  float32  `json:"totalIncome"`
	TotalOutcome float32  `json:"totalOutcome"`
	Total        float32  `json:"total"`
	Details      []Detail `json:"details"`
}

type Detail struct {
	Amount   float32           `json:"amount"`
	Category category.Category `json:"category"`
}

type GenerateReportUri struct {
	Year  uint `uri:"year" binding:"required"`
	Month uint `uri:"month" binding:"required"`
}

func GenerateReport(year uint, month uint) (Report, error) {
	var report Report
	report.Year = year
	report.Month = month
	categoryMap := make(map[uint]*category.Category)
	detailsMap := make(map[uint]float32)
	var withoutCategorySum float32
	var transactions []account_transaction.AccountTransaction
	if err := account_transaction.GetAccountTransactionsByYearAndMonth(year, month, &transactions); err != nil {
		return report, err
	}
	for _, transaction := range transactions {
		amount := float32(transaction.Amount) / 100
		report.Total += amount
		if amount < 0 {
			report.TotalOutcome += amount
		} else if amount > 0 {
			report.TotalIncome += amount
		}
		if transaction.Category != nil {
			categoryMap[*transaction.CategoryId] = transaction.Category
			detailsMap[*transaction.CategoryId] += amount
		} else {
			withoutCategorySum += amount
		}
	}
	for categoryId, sum := range detailsMap {
		report.Details = append(report.Details, Detail{Amount: sum, Category: *categoryMap[categoryId]})
	}
	report.Details = append(report.Details, Detail{Amount: withoutCategorySum, Category: category.Category{Name: "No category"}})
	return report, nil
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "report"
	rest.GET(controllerName+"/:month/:year", generateReport)
}

func generateReport(context *gin.Context) {
	var generateReportUri GenerateReportUri

	if err := context.ShouldBindUri(&generateReportUri); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	report, err := GenerateReport(generateReportUri.Year, generateReportUri.Month)
	if err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, report)
}
