package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"personal-finance-app/controllers"
	"personal-finance-app/database"
	"personal-finance-app/models"
)

func main() {
	r := gin.Default()

	// HTMLテンプレートの場所をロード
	r.LoadHTMLGlob("templates/*")

	// データベース接続
	database.ConnectDatabase()

	// Webページ表示用のルート
	r.GET("/", func(c *gin.Context) { // トップページ：円グラフと登録フォーム
		// 直近1ヶ月間のカテゴリ別支出を計算
		type CategorySummary struct {
			Category string
			Total    int
		}
		var expenseSummary []CategorySummary
		oneMonthAgo := time.Now().AddDate(0, -1, 0).UTC()
		database.DB.Model(&models.Transaction{}).
			Select("category, SUM(amount) as total").
			Where("type = ? AND date >= ? AND category IN ?", models.Expense, oneMonthAgo, models.AllExpenseCategories).
			Group("category").
			Find(&expenseSummary)

		// 直近1ヶ月間のカテゴリ別収入を計算
		var incomeSummary []CategorySummary
		database.DB.Model(&models.Transaction{}).
			Select("category, SUM(amount) as total").
			Where("type = ? AND date >= ? AND category IN ?", models.Income, oneMonthAgo, models.AllIncomeCategories).
			Group("category").
			Find(&incomeSummary)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"expenseSummary": expenseSummary,
			"incomeSummary":  incomeSummary,
		})
	})

	r.GET("/list", func(c *gin.Context) { // 一覧ページ：棒グラフと取引履歴
		var transactions []models.Transaction
		// 日付の降順（新しいものが上）で全件取得
		database.DB.Order("date desc").Find(&transactions)

		// --- 棒グラフ用のデータ集計 (過去12ヶ月) ---
		type MonthlyExpense struct {
			Month    string `gorm:"column:month"`
			Category string
			Total    int
		}
		var monthlyExpenses []MonthlyExpense
		now := time.Now().UTC() // タイムゾーン問題を避けるため、現在時刻をUTCで統一
		twelveMonthsAgo := now.AddDate(0, -11, 0)
		firstDayOfPeriod := time.Date(twelveMonthsAgo.Year(), twelveMonthsAgo.Month(), 1, 0, 0, 0, 0, time.UTC)

		database.DB.Model(&models.Transaction{}).
			Select("strftime('%Y-%m', date) as month, category, sum(amount) as total").
			Where("type = ? AND date >= ? AND category IN ?", models.Expense, firstDayOfPeriod, models.AllExpenseCategories).
			Group("month, category").
			Order("month").
			Scan(&monthlyExpenses)

		// Chart.js用のデータ構造に変換
		monthLabels := make([]string, 12)
		monthMap := make(map[string]int) // "YYYY-MM" -> index
		for i := 0; i < 12; i++ {
			month := now.AddDate(0, -(11 - i), 0)
			label := month.Format("2006-01")
			monthLabels[i] = label
			monthMap[label] = i
		}

		categoryData := make(map[string][]int)
		for _, me := range monthlyExpenses {
			if _, ok := categoryData[me.Category]; !ok {
				categoryData[me.Category] = make([]int, 12)
			}
			if monthIndex, ok := monthMap[me.Month]; ok {
				categoryData[me.Category][monthIndex] = me.Total
			}
		}

		type BarChartDataset struct {
			Label           string `json:"label"`
			Data            []int  `json:"data"`
			BackgroundColor string `json:"backgroundColor"`
		}
		chartColors := []string{"#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0", "#9966FF", "#FF9F40", "#E7E9ED", "#8DDF3C", "#F178B4", "#6A2E35", "#C4D7F2", "#A2D4AB"}
		barChartDatasets := []BarChartDataset{}
		colorIndex := 0
		for category, data := range categoryData {
			barChartDatasets = append(barChartDatasets, BarChartDataset{Label: category, Data: data, BackgroundColor: chartColors[colorIndex%len(chartColors)]})
			colorIndex++
		}

		barChartDataBytes, err := json.Marshal(gin.H{"labels": monthLabels, "datasets": barChartDatasets})
		if err != nil {
			// エラーが発生した場合は、空のJSONを渡す
			barChartDataBytes = []byte("{}")
		}

		c.HTML(http.StatusOK, "list.html", gin.H{
			"transactions":     transactions,
			"barChartDataJSON": template.JS(barChartDataBytes),
		})
	})

	// API ルーティング
	r.POST("/transactions", controllers.CreateTransaction)
	r.GET("/transactions", controllers.FindTransactions)
	r.GET("/transactions/:id", controllers.FindTransaction)
	r.PUT("/transactions/:id", controllers.UpdateTransaction)
	r.DELETE("/transactions/:id", controllers.DeleteTransaction)

	r.Run() // サーバーを起動 (デフォルトは :8080)
}
