package pay

import "github.com/DriveFluency/02-Backend/internal/domain"



type Service interface {
	GetByID(id int) (domain.Pay, error)
	Create(p domain.Pay) (domain.Pay, error)
	GetAll()([]domain.Pay,error)
}

type service struct {
	repo Repository
}

// constr
func NewServicePay( repo Repository) Service {
	return &service{repo}
}


// crud
func (serv *service) GetByID(id int) (domain.Pay, error) {
	Pay, err := serv.repo.GetByID(id)
	if err != nil {
		return domain.Pay{}, err
	}
	return Pay, nil

}

func (serv *service) Create(p domain.Pay) (domain.Pay, error) {
	
	p,err := serv.repo.Create(p)
	if err != nil {
		return domain.Pay{}, err
	}
	return p, nil

}

func (serv *service) GetAll()([]domain.Pay,error){
	Pays,err:= serv.repo.GetAll()
	if err != nil {
		return nil,err
	}
	return Pays,nil

}


