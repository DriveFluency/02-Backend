package store

import (

"github.com/DriveFluency/02-Backend/internal/domain"
	
)

type PackInterface interface{

	Read(id int)(domain.Pack, error)
	Create(Pack domain.Pack)error
	Update(Pack domain.Pack)error 
	Delete(id int)error
	GetAll()([]domain.Pack, error)

}