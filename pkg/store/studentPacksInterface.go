package store

import (

"github.com/DriveFluency/02-Backend/internal/domain"
	
)

type StudentPacksInterface interface{
	SearchByDni(dni int)([]domain.InfoPacks, error)
	Create(StudentPacks domain.StudentPacks) error
	GetAll()([]domain.InfoPacks, error)
	/*Update(Pack domain.Pack)error 
	Delete(id int)error*/	
}

