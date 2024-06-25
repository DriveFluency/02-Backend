package studentPacks

import (
	"errors"
	"github.com/DriveFluency/02-Backend/internal/domain"
	"github.com/DriveFluency/02-Backend/pkg/store"
	
)

type Repository interface {
	SearchByDni(dni int)([]domain.InfoPacks, error)
	Create(StudentPacks domain.StudentPacks)(domain.StudentPacks, error)
	GetAll()([]domain.InfoPacks, error)
}

type repository struct {
	storage store.StudentPacksInterface
}


func NewRepositoryStudentPacks(storage store.StudentPacksInterface) Repository {
	return &repository{storage}
}


func (repo *repository) SearchByDni(dni int) ([]domain.InfoPacks, error) {
	StudentPack, err := repo.storage.SearchByDni(dni)
	if err != nil {
		return []domain.InfoPacks{}, errors.New("No se encontraron packs para ese dni")
	}
	return StudentPack, nil

}

func (repo *repository) Create(p domain.StudentPacks) (domain.StudentPacks, error) {
	err := repo.storage.Create(p)
	if err != nil {
		return domain.StudentPacks{}, err 
	}
	return p, nil

}

func (repo *repository) GetAll()([]domain.InfoPacks,error){
	StudentPacks,err:= repo.storage.GetAll()
	if err != nil {
		return nil,errors.New("No se encontraron packs")
	}
	return StudentPacks,nil

}