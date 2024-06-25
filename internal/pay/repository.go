package pay

import (
	"errors"
	"github.com/DriveFluency/02-Backend/internal/domain"
	"github.com/DriveFluency/02-Backend/pkg/store"
	
)

type Repository interface {
	GetByID(id int) (domain.Pay, error)
	Create(p domain.Pay) (domain.Pay, error)
	GetAll()([]domain.Pay,error)
}

type repository struct {
	storage store.PayInterface
}


func NewRepositoryPay(storage store.PayInterface) Repository {
	return &repository{storage}
}


func (repo *repository) GetByID(id int) (domain.Pay, error) {
	Pay, err := repo.storage.Read(id)
	if err != nil {
		return domain.Pay{}, errors.New("Pago no encontrado")
	}
	return Pay, nil

}

func (repo *repository) Create(p domain.Pay) (domain.Pay, error) {
	err := repo.storage.Create(p)
	if err != nil {
		return domain.Pay{}, err 
	}
	return p, nil

}

func (repo *repository) GetAll()([]domain.Pay,error){
	Pays,err:= repo.storage.GetAll()
	if err != nil {
		return nil,errors.New("No se encontraron pagos")
	}
	return Pays,nil

}