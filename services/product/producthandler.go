package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductRepo interface {
	AddProduct(product Product) error
	SearchProductByName(query string) (*[]Product, error)
}
type Handler struct {
	repo ProductRepo
}

func NewHandler(repo ProductRepo) *Handler {
	return &Handler{
		repo,
	}
}
func (h *Handler) AddProduct(c *gin.Context) {
	var prod Product
	name := c.GetString("username")
	if name == "" || name != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "only admin can add product"})
		return
	}
	if err := c.BindJSON(&prod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err := h.repo.AddProduct(prod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "successfully added product"})
}
func (h *Handler) SearchProduct(c *gin.Context) {
	product_name := c.Query("name")
	if product_name == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}
	products, err := h.repo.SearchProductByName(product_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": products})
}
