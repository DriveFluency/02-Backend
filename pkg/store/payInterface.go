package store

import (

"github.com/DriveFluency/02-Backend/internal/domain"
	
)

type PayInterface interface{

	Read(id int)(domain.Pay, error)
	Create(Pay domain.Pay) ( int, error)
	GetAll()([]domain.Pay, error)
	/*Update(Pack domain.Pack)error 
	Delete(id int)error*/
	

}