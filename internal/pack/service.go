package pack

import "github.com/DriveFluency/02-Backend/internal/domain"



type Service interface {
	GetByID(id int) (domain.Pack, error)
	Create(p domain.Pack) (domain.Pack, error)
	Update(id int, p domain.Pack) (domain.Pack, error)
	Delete(id int) error
	GetAll()([]domain.Pack,error)
}

type service struct {
	repo Repository
}

// constr
func NewServicePack( repo Repository) Service {
	return &service{repo}
}


// crud
func (serv *service) GetByID(id int) (domain.Pack, error) {
	Pack, err := serv.repo.GetByID(id)
	if err != nil {
		return domain.Pack{}, err
	}
	return Pack, nil

}

func (serv *service) Create(p domain.Pack) (domain.Pack, error) {
	
	p,err := serv.repo.Create(p)
	if err != nil {
		return domain.Pack{}, err
	}
	return p, nil

}

// mejorar las validaciones... 
func (serv *service) Update(id int, pm domain.Pack) (domain.Pack, error) {
	
    PackBuscado,err:= serv.repo.GetByID(id)
	if err != nil {
		return domain.Pack{}, err
	}
	if pm.Name != ""{
		PackBuscado.Name =pm.Name
	}
	if pm.Description != ""{
		PackBuscado.Description =pm.Description
	}
	if pm.NumberClasses > 0 {
		PackBuscado.NumberClasses =pm.NumberClasses
	}
	if pm.DurationClasses > 0 {
		PackBuscado.DurationClasses =pm.DurationClasses
	}
	if pm.Cost > 0 {
		PackBuscado.Cost =pm.Cost}
	
	p,err := serv.repo.Update(id,PackBuscado) 
	if err != nil {
		return domain.Pack{}, err
	}
	return p, nil

}

func (serv *service) Delete(id int) error {
	err := serv.repo.Delete(id)
	if err != nil {
		return err
	}
	return  nil

}

func (serv *service) GetAll()([]domain.Pack,error){
	Packs,err:= serv.repo.GetAll()
	if err != nil {
		return nil,err
	}
	return Packs,nil

}


