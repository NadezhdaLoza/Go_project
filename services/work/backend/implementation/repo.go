package implementation

import (
    "database/sql"
    
    "github.com/nadezhda/documents/services/work/backend/repo"
    "github.com/nadezhda/documents/services/work/domain"
    
    _ "github.com/lib/pq"
)

type repository struct {
    db *sql.DB
}

func NewWorkRepo(db *sql.DB) repo.WorkRepo {
    return &repository{
        db: db,
    }
}

func (r *repository) CreateWork(work domain.Work) (*domain.Work, error) {
    sql := `
        INSERT INTO works (name, rate, unit) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `
    result, err := r.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    var insertedId int
    if err := result.QueryRow(work.Name, work.Rate, work.Unit).Scan(&insertedId); err != nil {
        return nil, err
    }

    return r.WorkById(insertedId)
}

func (r *repository) GetWorks() ([]domain.Work, error) {
    sql := `
        SELECT id, name, rate, unit 
        FROM works
        ORDER BY id
    `

    rows, err := r.db.Query(sql)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    works := []domain.Work{}

    for rows.Next() {
        work := domain.Work{}
        rows.Scan(
            &work.ID,
            &work.Name,
            &work.Rate,
            &work.Unit,
        )
        works = append(works, work)
    }

    return works, nil
}

func (r *repository) UpdateWork(work domain.Work) (*domain.Work, error) {
    sql := `
        UPDATE works 
        SET name = $1, 
            rate = $2,
            unit = $3
        WHERE id = $4
    `
    _, err := r.db.Exec(sql, work.Name, work.Rate, work.Unit, work.ID)
    if err != nil {
        return nil, err
    }

    return r.WorkById(work.ID)
}

func (r *repository) DeleteWork(id int) error {
    sql := `
        DELETE FROM works WHERE id = $1
    `
    _, err := r.db.Exec(sql, id)
    return err
}

func (r *repository) WorkById(id int) (*domain.Work, error) {
    sql := `
        SELECT id, name, rate, unit 
        FROM works 
        WHERE id = $1
    `
    work := domain.Work{}
    if err := r.db.QueryRow(sql, id).Scan(
        &work.ID,
        &work.Name,
        &work.Rate,
        &work.Unit,
    ); err != nil {
        return nil, err
    }
    return &work, nil
}