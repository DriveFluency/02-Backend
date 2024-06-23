package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/DriveFluency/02-Backend/internal/domain"
)

type sqlPack struct {
	db *sql.DB
}

func NewSqlPack(db *sql.DB) PackInterface {
	return &sqlPack{
		db: db,
	}
}

func (p *sqlPack) Read(id int) (domain.Pack, error) {

	Pack := domain.Pack{}
	query:="SELECT * FROM Packs WHERE id=? ;"
	row := p.db.QueryRow(query,id)
	err:= row.Scan(&Pack.Id,&Pack.Name,&Pack.Description,&Pack.NumberClasses,&Pack.DurationClasses,&Pack.Cost)
	if err != nil{
		return domain.Pack{},err
	}
	return  Pack,nil
}

func (p *sqlPack) Create(Pack domain.Pack)  error {

	var id int
	queryExists := "SELECT id FROM Packs WHERE name = ?"
	row := p.db.QueryRow(queryExists,Pack.Name)
	err:= row.Scan(&id)
	if err == nil{
		    fmt.Println(id)
			return errors.New("Pack existente")
	}
	query:="INSERT INTO Packs(name,description,number_classes,duration_classes,cost) VALUES (?,?,?,?,?);"
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return err
	}

	res,err := stmt.Exec(Pack.Name,Pack.Description,Pack.NumberClasses,Pack.DurationClasses,Pack.Cost)	
	if err != nil{
		return err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return err
	}
	return nil
}


func (p *sqlPack) Update(Pack domain.Pack)  error {

	query:="UPDATE Packs SET name=?, description=?, number_classes=?, duration_classes=?, cost=? WHERE id=? ;"
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return err
	}
	res,err := stmt.Exec(Pack.Name,Pack.Description,Pack.NumberClasses,Pack.DurationClasses,Pack.Cost,Pack.Id)	
	if err != nil{
		return err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return err
	}

	return nil
}
func (p *sqlPack) Delete(id int) error {

	query:="DELETE FROM Packs WHERE id=? "
	stmt,err := p.db.Prepare(query)
	if err != nil{
		return err
	}
	res,err := stmt.Exec(id)	
	if err != nil{
		return err
	}
	_,err = res.RowsAffected()
	if err != nil{
		return err
	}

	return nil

}


func (p *sqlPack) GetAll() ([]domain.Pack, error) {
	
	Packs := []domain.Pack{}
	query:="SELECT * FROM Packs  ;"
	rows,err := p.db.Query(query)
	if err != nil{
		log.Fatal(err)
	}
	err = rows.Err()
	if err != nil{
		log.Fatal(err)
	}
	for rows.Next(){
		var Pack domain.Pack
		err:= rows.Scan(&Pack.Id,&Pack.Name,&Pack.Description,&Pack.NumberClasses,&Pack.DurationClasses,&Pack.Cost)
		if err != nil{
			log.Fatal(err)
		}
		Packs = append(Packs,Pack )
	}
	return Packs,nil
}
