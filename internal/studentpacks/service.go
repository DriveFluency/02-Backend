package studentPacks

import "github.com/DriveFluency/02-Backend/internal/domain"



type Service interface {
	SearchByDni(dni int) ([]domain.InfoPacks, error)
	Create(p domain.StudentPacks) (domain.StudentPacks, error)
	GetAll()([]domain.InfoPacks,error)
}

type service struct {
	repo Repository
}

// constr
func NewServiceStudentPacks( repo Repository) Service {
	return &service{repo}
}


// crud
func (serv *service) SearchByDni(dni int) ([]domain.InfoPacks, error) {
	StudentPacks, err := serv.repo.SearchByDni(dni)
	if err != nil {
		return []domain.InfoPacks{}, err
	}
	return StudentPacks, nil

}

func (serv *service) Create(p domain.StudentPacks) (domain.StudentPacks, error) {
	
	p,err := serv.repo.Create(p)
	if err != nil {
		return domain.StudentPacks{}, err
	}
	return p, nil

}

func (serv *service) GetAll()([]domain.InfoPacks,error){
	StudentPackss,err:= serv.repo.GetAll()
	if err != nil {
		return nil,err
	}
	return StudentPackss,nil

}


