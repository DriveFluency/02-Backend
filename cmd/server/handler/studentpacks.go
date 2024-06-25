package handler

import (
	"net/http"
	"strconv"
	"github.com/DriveFluency/02-Backend/internal/studentpacks"
	"github.com/gin-gonic/gin"
)

type StudentPacksHandler struct {
	service studentPacks.Service
}

func NewStudentPacksHandler(service studentPacks.Service) *StudentPacksHandler {
	return &StudentPacksHandler{service: service}
}

// @Summary Buscar StudentPacks
// @Tag domain.StudentPacks
// @Param dni path int true "Item dni"
// @Success 200
// @Router /StudentPacks/{dni} [get]
func (h *StudentPacksHandler) SearchByDni() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dniParam := ctx.Param("dni")
		dni, err := strconv.Atoi(dniParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "formato id inválido"})
			return
		}
		p, err := h.service.SearchByDni(dni)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "id inexistente"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"StudentPacks": p})

	}

}

// @Summary Buscar StudentPackss
// @Tag domain.StudentPacks
// @Success 200
// @Router /StudentPacks [get]
func (h *StudentPacksHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		StudentsPacks, err := h.service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"StudentsPacks": StudentsPacks})
	}
}
