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

func (p *sqlPay) Create(Pay domain.Pay)  error {

	query:="INSERT INTO pay(date,method,amount,receipt) VALUES (?,?,?,?);"
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return err
	}
	res,err := stmt.Exec(Pay.Date,Pay.Method,Pay.Amount,Pay.Receipt)	
	if err != nil{
		return err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return err
	}
	return nil
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
