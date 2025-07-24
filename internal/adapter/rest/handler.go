package rest

import (
	"fmt"
	"net/http"

	"computerInventory/internal/domain"
	"computerInventory/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *usecase.Service
}

func NewHandler(service *usecase.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/computers", h.addComputer)
	router.GET("/computers", h.getAll)
	router.GET("/computers/:mac", h.getByMac)
	router.DELETE("/computers/:mac", h.deleteComputer)
	router.PUT("/computers/:mac", h.updateComputer)
	router.GET("/employee/:abbr/computers", h.getByEmployee)
}

func (h *Handler) addComputer(c *gin.Context) {
	var comp domain.Computer
	if err := c.ShouldBindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddComputer(&comp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comp)
}

func (h *Handler) getAll(c *gin.Context) {
	comps, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comps)
}

func (h *Handler) getByMac(c *gin.Context) {
	mac := c.Param("mac")
	comp, err := h.service.Get(mac)
	if err != nil {
		fmt.Println("getByMac: not found on " + mac)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	fmt.Println("getByMac: found on " + mac)
	c.JSON(http.StatusOK, comp)
}

func (h *Handler) deleteComputer(c *gin.Context) {
	mac := c.Param("mac")
	err := h.service.Delete(mac)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) updateComputer(c *gin.Context) {
	mac := c.Param("mac")
	var comp domain.Computer
	if err := c.ShouldBindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comp.MACAddress = mac
	if err := h.service.Update(&comp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

func (h *Handler) getByEmployee(c *gin.Context) {
	abbr := c.Param("abbr")
	comps, err := h.service.GetByEmployee(abbr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comps)
}
