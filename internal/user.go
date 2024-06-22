package internal

type User struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	LastName string `json:"lastname" binding:"required"`
	DNI       string `json:"dni" binding:"required"`
	Telefono  string `json:"telefono" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Localidad string `json:"localidad" binding:"required"`
	Ciudad string `json:"ciudad" binding:"required"`

}