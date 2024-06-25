package store

import (
	"database/sql"
	"log"
	"github.com/DriveFluency/02-Backend/internal/domain"
)

type sqlPay struct {
	db *sql.DB
}

func NewSqlPay(db *sql.DB) PayInterface {
	return &sqlPay{
		db: db,
	}
}

func (p *sqlPay) Read(id int) (domain.Pay, error) {

	Pay := domain.Pay{}
	query:="SELECT * FROM pay WHERE id=? ;"
	row := p.db.QueryRow(query,id)
	err:= row.Scan(&Pay.Id,&Pay.Date,&Pay.Method,&Pay.Amount,&Pay.Receipt)
	if err != nil{
		return domain.Pay{},err
	}
	return  Pay,nil
}

func (p *sqlPay) Create(Pay domain.Pay) (int,error) {

	query:="INSERT INTO pay(date,method,amount,receipt) VALUES (?,?,?,?);"
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return 0 , err
	}
	res,err := stmt.Exec(Pay.Date,Pay.Method,Pay.Amount,Pay.Receipt)	
	if err != nil{
		return 0,err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return 0,err
	}

	ultimoId,err := res.LastInsertId()
	log.Print(ultimoId)

	/*var pay domain.Pay

    query2 :="SELECT id FROM pay order by id desc limit 1 ;"
    row,err := p.db.Query(query2)
	if err != nil{
		log.Fatal(err)
	}
	/*
	for row.Next(){
	
		err:= row.Scan(&pay.Id)
		if err != nil{
			log.Fatal(err)
		}

	}
*/
	return int(ultimoId), nil

}


func (p *sqlPay) GetAll() ([]domain.Pay, error) {	
	Pays := []domain.Pay{}
	query:="SELECT * FROM pay ;"
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
		var Pay domain.Pay
		err:= rows.Scan(&Pay.Id,&Pay.Date,&Pay.Method,&Pay.Amount,&Pay.Receipt)
		if err != nil{
			log.Fatal(err)
		}

		
		Pays = append(Pays,Pay )
	}
	return Pays,nil
}
