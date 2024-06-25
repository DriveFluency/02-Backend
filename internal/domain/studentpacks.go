package domain

type StudentPacks struct{
Id int `json:"id"`
StudentDNI int `json:"student_dni" binding:"required" `
PackId int `json:"pack_id" binding:"required" `
PayId int `json:"pay_id" binding:"required" `
}


type InfoPacks struct{
	Id int `json:"id"`
	StudentDNI int `json:"student_dni" binding:"required" `
	Date string 
	Method string
	NamePack string
	DescriptionPack string 
	Cost float64
}

