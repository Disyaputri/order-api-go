// Package delivery berisi lapisan HTTP (Gin): router, handler, dan format response JSON.
// Bentuk response mengikuti kontrak API (OpenAPI) dari soal.
package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lc1-order-api/domain"
)

// OrderHandler menghubungkan request HTTP ke usecase.
type OrderHandler struct {
	uc domain.OrderUsecase
}

// NewRouter mendaftarkan semua endpoint Gin.
func NewRouter(uc domain.OrderUsecase) *gin.Engine {
	h := &OrderHandler{uc: uc}
	r := gin.Default()
	r.GET("/orders", h.list)
	r.POST("/orders", h.create)
	r.GET("/orders/:id", h.detail)
	r.PUT("/orders/:id", h.update)
	r.DELETE("/orders/:id", h.remove)
	return r
}

// writeError memilih status code dan isi pesan sesuai kontrak.
// id dipakai untuk pesan 404; isi 0 bila endpoint tidak memakai ID.
func writeError(c *gin.Context, err error, id int64) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("order with ID %d not found", id)})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "detail": err.Error()})
	}
}

// parseID membaca :id dari URL. ID yang bukan angka dianggap tidak ditemukan (404),
// karena kontrak untuk endpoint ber-ID hanya mendefinisikan 200, 404, dan 500.
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("order with ID %s not found", c.Param("id"))})
		return 0, false
	}
	return id, true
}

// parseBody membaca JSON request body ke OrderInput.
func parseBody(c *gin.Context) (domain.OrderInput, bool) {
	var in domain.OrderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input: request body must be valid JSON with correct field types"})
		return in, false
	}
	return in, true
}

func (h *OrderHandler) list(c *gin.Context) {
	orders, err := h.uc.GetAll()
	if err != nil {
		writeError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) create(c *gin.Context) {
	in, ok := parseBody(c)
	if !ok {
		return
	}
	order, err := h.uc.Create(in)
	if err != nil {
		writeError(c, err, 0)
		return
	}
	// Kontrak memakai nama key "shipment" pada response POST.
	c.JSON(http.StatusCreated, gin.H{"message": "success create", "shipment": order})
}

func (h *OrderHandler) detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	order, err := h.uc.GetByID(id)
	if err != nil {
		writeError(c, err, id)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	in, ok := parseBody(c)
	if !ok {
		return
	}
	order, err := h.uc.Update(id, in)
	if err != nil {
		writeError(c, err, id)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success update", "shipment": order})
}

func (h *OrderHandler) remove(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.uc.Delete(id); err != nil {
		writeError(c, err, id)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success delete"})
}
