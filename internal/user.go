package internal 


type User struct{
Username string `json:"username" binding:"required"` // first name .... 
Email string `json:"email" binding:"required"`
Rol string `json:"rol" binding:"required"`
DNI string `json:"dni" binding:"required"`
} 