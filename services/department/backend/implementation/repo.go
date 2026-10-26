package implementation

import (
    "database/sql"
    
    "github.com/nadezhda/documents/services/department/backend/repo"
    "github.com/nadezhda/documents/services/department/domain"
    
    _ "github.com/lib/pq"
)

type repository struct {
    db *sql.DB
}

func NewDepartmentRepo(db *sql.DB) repo.DepartmentRepo {
    return &repository{
        db: db,
    }
}

func (r *repository) CreateDepartment(department domain.Department) (*domain.Department, error) {
    sql := `
        INSERT INTO departments (name, work_conditions_class, work_day_start, work_day_end) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `
    result, err := r.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    var insertedId int
    if err := result.QueryRow(
        department.Name,
        department.WorkConditionsClass,
        department.WorkDayStart,
        department.WorkDayEnd,
    ).Scan(&insertedId); err != nil {
        return nil, err
    }

    return r.DepartmentById(insertedId)
}

func (r *repository) GetDepartments() ([]domain.Department, error) {
    sql := `
        SELECT id, name, work_conditions_class, TO_CHAR(work_day_start, 'HH24:MI') as work_day_start, TO_CHAR(work_day_end, 'HH24:MI') as work_day_end
        FROM departments
        ORDER BY id
    `

    rows, err := r.db.Query(sql)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    departments := []domain.Department{}

    for rows.Next() {
        department := domain.Department{}
        rows.Scan(
            &department.ID,
            &department.Name,
            &department.WorkConditionsClass,
            &department.WorkDayStart,
            &department.WorkDayEnd,
        )
        departments = append(departments, department)
    }

    return departments, nil
}

func (r *repository) UpdateDepartment(department domain.Department) (*domain.Department, error) {
    sql := `
        UPDATE departments 
        SET name = $1, 
            work_conditions_class = $2, 
            work_day_start = $3, 
            work_day_end = $4 
        WHERE id = $5
    `
    _, err := r.db.Exec(sql, 
        department.Name,
        department.WorkConditionsClass,
        department.WorkDayStart,
        department.WorkDayEnd,
        department.ID,
    )
    if err != nil {
        return nil, err
    }

    return r.DepartmentById(department.ID)
}

func (r *repository) DeleteDepartment(id int) error {
    sql := `
        DELETE FROM departments WHERE id = $1
    `
    _, err := r.db.Exec(sql, id)
    return err
}

func (r *repository) DepartmentById(id int) (*domain.Department, error) {
    sql := `
        SELECT id, name, work_conditions_class, TO_CHAR(work_day_start, 'HH24:MI') as work_day_start, TO_CHAR(work_day_end, 'HH24:MI') as work_day_end
        FROM departments 
        WHERE id = $1
    `
    department := domain.Department{}
    if err := r.db.QueryRow(sql, id).Scan(
        &department.ID,
        &department.Name,
        &department.WorkConditionsClass,
        &department.WorkDayStart,
        &department.WorkDayEnd,
    ); err != nil {
        return nil, err
    }
    return &department, nil
}