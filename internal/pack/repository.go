package pack

import (
	"errors"
	"github.com/DriveFluency/02-Backend/internal/domain"
	"github.com/DriveFluency/02-Backend/pkg/store"
	
)

type Repository interface {
	GetByID(id int) (domain.Pack, error)
	Create(p domain.Pack) (domain.Pack, error)
	Update(id int, p domain.Pack) (domain.Pack, error)
	Delete(id int) error
	GetAll()([]domain.Pack,error)
}

type repository struct {
	storage store.PackInterface
}


func NewRepositoryPack(storage store.PackInterface) Repository {
	return &repository{storage}
}


func (repo *repository) GetByID(id int) (domain.Pack, error) {
	Pack, err := repo.storage.Read(id)
	if err != nil {
		return domain.Pack{}, errors.New("Pack no encontrado")
	}
	return Pack, nil

}

func (repo *repository) Create(p domain.Pack) (domain.Pack, error) {
	err := repo.storage.Create(p)
	if err != nil {
		return domain.Pack{}, err 
	}
	return p, nil

}

func (repo *repository) Update(id int, PackModificado domain.Pack) (domain.Pack, error) {
	err := repo.storage.Update(PackModificado)
	if err != nil {
		return domain.Pack{},errors.New("error al modificar el Pack")
	}
	return PackModificado, nil

}

func (repo *repository) Delete(id int) error {
	err := repo.storage.Delete(id)
	if err != nil {
		return errors.New("El Pack no pudo ser eliminado")
	}
	return  nil

}

func (repo *repository) GetAll()([]domain.Pack,error){
	Packs,err:= repo.storage.GetAll()
	if err != nil {
		return nil,errors.New("No se encontraron Packs")
	}
	return Packs,nil

}