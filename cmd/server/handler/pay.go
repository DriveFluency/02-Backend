package handler

import (
//	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DriveFluency/02-Backend/internal/domain"
	"github.com/DriveFluency/02-Backend/internal/pay"
	"github.com/DriveFluency/02-Backend/internal/studentpacks"
	"github.com/gin-gonic/gin"
)

type payHandler struct {
	service pay.Service
	studentPacksService studentPacks.Service
}

func NewPayHandler(service pay.Service, studentService studentPacks.Service) *payHandler {
	return &payHandler{service: service, studentPacksService:studentService  }
}

// @Summary Buscar pay
// @Tag domain.pay
// @Param id path int true "Item ID"
// @Success 200
// @Router /pay/{id} [get]
func (h *payHandler) GetByID() gin.HandlerFunc {
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

		ctx.JSON(http.StatusOK, gin.H{"pay": p})

	}

}

// faltan validaciones de campos

// @Summary Crear pay
// @Tag domain.pay
// @Param token header string true "TOKEN"
// @Success 201
// @Router /pay [post]
func (h *payHandler) Post() gin.HandlerFunc {
/**/
type Data struct{

	Dni int `json:"dni" binding:"required"`
	PackId int `json:"pack_id" binding:"required"`
	Method  string    `json:"method" binding:"required"`
	Amount  float64   `json:"amount" binding:"required"`
	Receipt string    `json:"receipt" binding:"required"`
}

	return func(ctx *gin.Context) {	
		var data Data 
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var pay = domain.Pay{
			Date: time.Now().Format("2006-01-02 15:04:05.00"),
			Method: data.Method,
			Amount: data.Amount,
			Receipt: data.Receipt,
		}
	
		p, err := h.service.Create(pay)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// llamamos al service y le pasamos el dni y los ids... 

		var studentPack= domain.StudentPacks{
			StudentDNI:data.Dni ,
			PackId: data.PackId,
			PayId:p.Id , 
		}
		studentPacks, err := h.studentPacksService.Create(studentPack)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"pay": p, "student packs": studentPacks })

	}
}

// @Summary Buscar pays
// @Tag domain.pay
// @Success 200
// @Router /pay [get]
func (h *payHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pays, err := h.service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"pays": pays})
	}
}
