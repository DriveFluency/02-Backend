package internal

type User struct {
	Username  string `json:"username" binding:"required"` 
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	LastName string `json:"lastname" binding:"required"`
	DNI       int `json:"dni" binding:"required"`
	Telefono  int `json:"telefono" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
}