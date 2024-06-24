package handler

import (
	"github.com/DriveFluency/02-Backend/internal/domain"
	"github.com/DriveFluency/02-Backend/internal/pack"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type packHandler struct {
	service pack.Service
}

func NewPackHandler(service pack.Service) *packHandler {
	return &packHandler{service: service}
}

// @Summary Buscar pack
// @Tag domain.pack
// @Param id path int true "Item ID"
// @Success 200
// @Router /packs/{id} [get]
func (h *packHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "formato id inválido"})
			return
		}
		p, err := h.service.GetByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "id inexistente"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"pack": p})

	}

}

// faltan validaciones de campos

// @Summary Crear pack
// @Tag domain.pack
// @Param token header string true "TOKEN"
// @Success 201
// @Router /packs [post]
func (h *packHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pack domain.Pack
		err := ctx.ShouldBindJSON(&pack)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		p, err := h.service.Create(pack)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"pack": p})

	}
}

// @Summary Modificar pack
// @Tag domain.pack
// @Param token header string true "TOKEN"
// @Param id path int true "Item ID"
// @Success 200
// @Router /packs/{id} [put]
func (h *packHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "formato id inválido"})
			return
		}
		var pack domain.Pack
		err = ctx.ShouldBindJSON(&pack)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		p, err := h.service.Update(id, pack)
		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"pack": p})

	}
}

// @Summary Modificar un campo del pack
// @Tag domain.pack
// @Param token header string true "TOKEN"
// @Param id path int true "Item ID"
// @Success 200
// @Router /packs/{id} [patch]
func (h *packHandler) Patch() gin.HandlerFunc {
	// estructura
	type Request struct {
		Name            string  `json:"name,omitempty"`
		Description     string  `json:"description,omitempty"`
		NumberClasses   int     `json:"number_classes,omitempty"`
		DurationClasses int     `json:"duration_classes,omitempty"`
		Cost            float64 `json:"cost,omitempty"`
	}

	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "formato id inválido"})
			return
		}
		var request Request
		err = ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		update := domain.Pack{
			Name:            request.Name,
			Description:     request.Description,
			NumberClasses:   request.NumberClasses,
			DurationClasses: request.DurationClasses,
			Cost:            request.Cost,
		}
		p, err := h.service.Update(id, update)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"pack": p})

	}
}

// @Summary Eliminar pack
// @Tag domain.pack
// @Param token header string true "TOKEN"
// @Param id path int true "Item ID"
// @Success 204
// @Router /packs/{id} [delete]
func (h *packHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "formato id inválido"})
			return
		}
		err = h.service.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusNoContent, gin.H{"message": "pack delete"})
	}
}

// @Summary Buscar packs
// @Tag domain.pack
// @Success 200
// @Router /packs [get]
func (h *packHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		packs, err := h.service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"packs": packs})
	}
}
