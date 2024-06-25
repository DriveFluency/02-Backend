package store

import (
	"database/sql"
	"log"
	"github.com/DriveFluency/02-Backend/internal/domain"
)

type sqlStudentPacks struct {
	db *sql.DB
}

func NewSqlStudentPacks(db *sql.DB) StudentPacksInterface {
	return &sqlStudentPacks{
		db: db,
	}
}


func (p *sqlStudentPacks) SearchByDni(dni int) ([]domain.InfoPacks,error){

	//el precioooo 
	var StudentsPacksDNI = []domain.InfoPacks{}
	query:="SELECT student_packs.id,student_packs.student_dni, pay.date, pay.method, packs.name, packs.description, packs.cost FROM student_packs INNER JOIN pay ON student_packs.pay_id=pay.id INNER JOIN packs ON student_packs.pack_id= packs.id WHERE student_dni=? ;"
	 stmt, err := p.db.Query(query, dni)
	 if err != nil {
		 log.Fatal(err)
	 }
	 for stmt.Next() {
		var StudentPacks domain.InfoPacks
		err:= stmt.Scan(&StudentPacks.Id,&StudentPacks.StudentDNI,&StudentPacks.Date,&StudentPacks.Method,&StudentPacks.NamePack,&StudentPacks.DescriptionPack, &StudentPacks.Cost)
		if err != nil{
			log.Fatal(err)
		}
		StudentsPacksDNI = append(StudentsPacksDNI,StudentPacks )
	}
	return StudentsPacksDNI,nil
	
}


func (p *sqlStudentPacks) Create(StudentPacks domain.StudentPacks)  error {

	query:="INSERT INTO student_packs(student_dni,pack_id,pay_id) VALUES (?,?,?);"
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return err
	}
	res,err := stmt.Exec(StudentPacks.StudentDNI,StudentPacks.PackId,StudentPacks.PayId)	
	if err != nil{
		return err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return err
	}

	// select que traiga el último id .... 

	
	return nil
}


func (p *sqlStudentPacks) GetAll() ([]domain.InfoPacks, error) {	
	StudentsPacks := []domain.InfoPacks{}
	query:="SELECT student_packs.id,student_packs.student_dni, pay.date, pay.method, packs.name, packs.description, packs.cost FROM student_packs INNER JOIN pay ON student_packs.pay_id=pay.id INNER JOIN packs ON student_packs.pack_id= packs.id ;"
	rows,err := p.db.Query(query)
	if err != nil{
		log.Fatal(err)
	}
	//defer rows.Close()	
	err = rows.Err()
	if err != nil{
		log.Fatal(err)
	}
	for rows.Next(){
		var StudentPacks domain.InfoPacks
		err:= rows.Scan(&StudentPacks.Id,&StudentPacks.StudentDNI,&StudentPacks.Date,&StudentPacks.Method,&StudentPacks.NamePack,&StudentPacks.DescriptionPack, &StudentPacks.Cost )
		if err != nil{
			log.Fatal(err)
		}
		StudentsPacks = append(StudentsPacks,StudentPacks )
	}
	return StudentsPacks,nil
}
