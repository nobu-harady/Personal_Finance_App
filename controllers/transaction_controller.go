package controllers

import (
	"fmt"
	"net/http"
	"personal-finance-app/database"
	"personal-finance-app/models"

	"github.com/gin-gonic/gin"
)

// CreateTransaction godoc
// @Summary 新しい取引を作成
// @Description 新しい収入または支出の記録を作成します
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param   transaction  body   models.Transaction  true  "取引情報"
// @Success 200 {object} models.Transaction
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// カテゴリのバリデーション
	if !models.IsValidCategory(input.Type, input.Category) {
		errMsg := fmt.Sprintf("invalid category '%s' for type '%s'", input.Category, input.Type)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	transaction := models.Transaction{
		Date:     input.Date,
		Type:     input.Type,
		Category: input.Category,
		Amount:   input.Amount,
		Memo:     input.Memo,
	}
	database.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// FindTransactions godoc
// @Summary 取引一覧を取得
// @Description すべての取引の一覧を取得します
// @Tags transactions
// @Produce  json
// @Success 200 {array} models.Transaction
// @Router /transactions [get]
func FindTransactions(c *gin.Context) {
	var transactions []models.Transaction
	database.DB.Find(&transactions)

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// FindTransaction godoc
// @Summary 特定の取引を取得
// @Description IDを指定して特定の取引を取得します
// @Tags transactions
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [get]
func FindTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// DeleteTransaction godoc
// @Summary 取引を削除
// @Description IDを指定して取引を削除します
// @Tags transactions
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]boolean
// @Router /transactions/{id} [delete]
func DeleteTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	database.DB.Delete(&transaction)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// UpdateTransaction godoc
// @Summary 取引を更新
// @Description IDを指定して取引を更新します
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Param   transaction  body   models.Transaction  true  "更新する取引情報"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [put]
func UpdateTransaction(c *gin.Context) {
	// First, find the existing transaction
	var transaction models.Transaction
	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	// Bind the input to a temporary struct
	var input models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the category for the given type
	if !models.IsValidCategory(input.Type, input.Category) {
		errMsg := fmt.Sprintf("invalid category '%s' for type '%s'", input.Category, input.Type)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Update the transaction in the database
	database.DB.Model(&transaction).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}
