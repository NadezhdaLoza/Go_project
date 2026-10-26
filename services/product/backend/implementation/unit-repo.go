package implementation

import (
	"database/sql"
	"errors"
	"sort"

	"github.com/nadezhda/documents/services/product/backend/repo"
	"github.com/nadezhda/documents/services/product/domain"
)

type unitRepository struct {
	db *sql.DB
}

func NewUnitRepository(db *sql.DB) repo.UnitRepo {
	return &unitRepository{
		db: db,
	}
}

func (r *unitRepository) CreateUnit(unit domain.Unit) (*domain.Unit, error) {
	sql := `
                INSERT INTO units (name) VALUES ($1) RETURNING id
        `
	result, err := r.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var insertedId int
	if err := result.QueryRow(unit.Name).Scan(&insertedId); err != nil {
		return nil, err
	}

	return r.UnitById(insertedId)
}

func (r *unitRepository) GetUnits() ([]domain.Unit, error) {
	sql := `SELECT * FROM units`
	rows, err := r.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := []domain.Unit{}

	for rows.Next() {
		unit := domain.Unit{}
		rows.Scan(&unit.Id, &unit.Name)
		units = append(units, unit)
	}

	// Сортировка по id
	sort.Slice(units, func(i, j int) bool {
		return units[i].Id < units[j].Id
	})

	return units, nil
}

func (r *unitRepository) UpdateUnit(unit domain.Unit) (*domain.Unit, error) {
	sql := `
                UPDATE units SET name = $1 WHERE id = $2
        `
	rows, err := r.db.Query(sql, unit.Name, unit.Id)
	if err != nil {
		return nil, err
	}

	rows.Close()

	return r.UnitById(unit.Id)
}

func (r *unitRepository) DeleteUnit(id int) error {
	if id == 0 {
		return errors.New("can't delete default value")
	}
	sql := `
                DELETE FROM units WHERE id = $1        
        `
	_, err := r.db.Exec(sql, id)
	return err
}

func (r *unitRepository) UnitById(id int) (*domain.Unit, error) {
	sql := `
                SELECT * FROM units WHERE id = $1    
        `
	unit := domain.Unit{}
	if err := r.db.QueryRow(sql, id).Scan(&unit.Id, &unit.Name); err != nil {
		return nil, err
	}

	return &unit, nil
}