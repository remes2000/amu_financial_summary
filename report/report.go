package report

import (
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/account_transaction"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/currency"
	"log"
	"net/http"
	"sort"
)

type Report struct {
	Year         uint     `json:"year"`
	Month        uint     `json:"month"`
	TotalIncome  string   `json:"totalIncome"`
	TotalOutcome string   `json:"totalOutcome"`
	Total        string   `json:"total"`
	Details      []Detail `json:"details"`
}

type Detail struct {
	CategoryName   string `json:"category"`
	Amount         string `json:"amount"`
	AmountAsNumber int    `json:"-"`
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
	detailsMap := make(map[uint]int)
	var withoutCategorySum int
	var total, totalOutcome, totalIncome int
	var transactions []account_transaction.AccountTransaction
	if err := account_transaction.GetAccountTransactionsByYearAndMonth(year, month, &transactions); err != nil {
		return report, err
	}
	for _, transaction := range transactions {
		amount := transaction.Amount
		total += amount
		if amount < 0 {
			totalOutcome += amount
		} else if amount > 0 {
			totalIncome += amount
		}
		if transaction.Category != nil {
			categoryMap[*transaction.CategoryId] = transaction.Category
			detailsMap[*transaction.CategoryId] += amount
		} else {
			withoutCategorySum += amount
		}
	}
	for categoryId, sum := range detailsMap {
		report.Details = append(report.Details, Detail{Amount: currency.FormatAsCurrency(sum), AmountAsNumber: sum, CategoryName: (*categoryMap[categoryId]).Name})
	}
	report.Details = append(report.Details, Detail{Amount: currency.FormatAsCurrency(withoutCategorySum), AmountAsNumber: withoutCategorySum, CategoryName: "No category"})
	sort.Slice(report.Details, func(i, j int) bool {
		amount1 := report.Details[i].AmountAsNumber
		if amount1 < 0 {
			amount1 *= -1
		}
		amount2 := report.Details[j].AmountAsNumber
		if amount2 < 0 {
			amount2 *= -1
		}
		return amount1 > amount2
	})
	report.Total = currency.FormatAsCurrency(total)
	report.TotalIncome = currency.FormatAsCurrency(totalIncome)
	report.TotalOutcome = currency.FormatAsCurrency(totalOutcome)
	return report, nil
}

// ---=== REST ===---

func BindRoutes(rest *gin.RouterGroup) {
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
