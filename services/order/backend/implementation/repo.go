package implementation

import (
    "database/sql"
    
    //"github.com/nadezhda/documents/services/base"
    "github.com/nadezhda/documents/services/order/backend/repo"
    "github.com/nadezhda/documents/services/order/domain"
    
    _ "github.com/lib/pq"
)

type orderRepo struct {
    db *sql.DB
}

func NewOrderRepo(db *sql.DB) repo.OrderRepo {
    return &orderRepo{
        db: db,
    }
}

// Заголовок наряда
func (r *orderRepo) CreateOrder(order domain.Order) (*domain.Order, error) {
    sql := `
        INSERT INTO work_orders 
        (order_number, order_date, department_id, foreman_id, customer_id, status) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id, total_amount
    `
    result, err := r.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    var insertedId int
    var totalAmount float64
    if err := result.QueryRow(
        order.OrderHeader.OrderNumber,
        order.OrderHeader.OrderDate,
        order.OrderHeader.DepartmentID,
        order.OrderHeader.ForemanID,
        order.OrderHeader.CustomerID,
        "draft",
    ).Scan(&insertedId, &totalAmount); err != nil {
        return nil, err
    }

    order.OrderHeader.ID = insertedId
    order.OrderHeader.TotalAmount = totalAmount
    order.OrderHeader.Status = "draft"
    return &order, nil
}

func (r *orderRepo) GetOrders() ([]domain.Order, error) {
    sql := `
        SELECT id, order_number, order_date, department_id, foreman_id, 
               customer_id, status, total_amount
        FROM work_orders
        ORDER BY order_date DESC, id DESC
    `
    rows, err := r.db.Query(sql)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    orders := []domain.Order{}
    for rows.Next() {
        header := domain.OrderHeader{}
        rows.Scan(
            &header.ID,
            &header.OrderNumber,
            &header.OrderDate,
            &header.DepartmentID,
            &header.ForemanID,
            &header.CustomerID,
            &header.Status,
            &header.TotalAmount,
        )
        
        order := domain.Order{
            OrderHeader: header,
        }
        orders = append(orders, order)
    }

    return orders, nil
}

func (r *orderRepo) OrderById(id int) (*domain.Order, error) {
    order := domain.Order{}
    
    // Загружаем заголовок
    headerSql := `
        SELECT id, order_number, order_date, department_id, foreman_id, 
               customer_id, status, total_amount
        FROM work_orders 
        WHERE id = $1
    `
    header := domain.OrderHeader{}
    if err := r.db.QueryRow(headerSql, id).Scan(
        &header.ID,
        &header.OrderNumber,
        &header.OrderDate,
        &header.DepartmentID,
        &header.ForemanID,
        &header.CustomerID,
        &header.Status,
        &header.TotalAmount,
    ); err != nil {
        return nil, err
    }
    
    order.OrderHeader = header

    // Загружаем работы (work_order_body1)
    itemsSql := `
        SELECT id, work_order_id, work_id, description, unit, quantity, rate, amount
        FROM work_order_body1 
        WHERE work_order_id = $1
    `
    rows, err := r.db.Query(itemsSql, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        item := domain.OrderItem{}
        rows.Scan(
            &item.ID,
            &item.OrderID,
            &item.WorkID,
            &item.Description,
            &item.Unit,
            &item.Quantity,
            &item.Rate,
            &item.Amount,
        )
        order.OrderBodyItems = append(order.OrderBodyItems, item)
    }

    // Загружаем исполнителей (work_order_body2)
    performersSql := `
        SELECT id, work_order_id, employee_id, ktu, amount
        FROM work_order_body2 
        WHERE work_order_id = $1
    `
    performersRows, err := r.db.Query(performersSql, id)
    if err != nil {
        return nil, err
    }
    defer performersRows.Close()

    for performersRows.Next() {
        performer := domain.Performer{}
        performersRows.Scan(
            &performer.ID,
            &performer.OrderID,
            &performer.EmployeeID,
            &performer.KTU,
            &performer.Amount,
        )
        order.Performers = append(order.Performers, performer)
    }

    return &order, nil
}

func (r *orderRepo) UpdateOrder(order domain.Order) (*domain.Order, error) {
    sql := `
        UPDATE work_orders 
        SET order_number = $1, 
            order_date = $2,
            department_id = $3,
            foreman_id = $4,
            customer_id = $5,
            status = $6
        WHERE id = $7
    `
    _, err := r.db.Exec(sql,
        order.OrderHeader.OrderNumber,
        order.OrderHeader.OrderDate,
        order.OrderHeader.DepartmentID,
        order.OrderHeader.ForemanID,
        order.OrderHeader.CustomerID,
        order.OrderHeader.Status,
        order.OrderHeader.ID,
    )
    if err != nil {
        return nil, err
    }

    return r.OrderById(order.OrderHeader.ID)
}

func (r *orderRepo) DeleteOrder(id int) error {
    sql := `DELETE FROM work_orders WHERE id = $1`
    _, err := r.db.Exec(sql, id)
    return err
}

// Работы в наряде (1-я таблица)
func (r *orderRepo) CreateOrderItem(item domain.OrderItem) (*domain.Order, error) {
    sql := `
        INSERT INTO work_order_body1 
        (work_order_id, work_id, description, unit, quantity, rate) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id, amount
    `
    result, err := r.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    defer result.Close()

    var insertedId int
    var amount float64
    if err := result.QueryRow(
        item.OrderID,
        item.WorkID,
        item.Description,
        item.Unit,
        item.Quantity,
        item.Rate,
    ).Scan(&insertedId, &amount); err != nil {
        return nil, err
    }

    item.ID = insertedId
    item.Amount = amount
    return r.OrderById(item.OrderID)
}

func (r *orderRepo) DeleteOrderItem(itemID int) error {
    // Получаем work_order_id для обновления
    var workOrderID int
    sql := `SELECT work_order_id FROM work_order_body1 WHERE id = $1`
    if err := r.db.QueryRow(sql, itemID).Scan(&workOrderID); err != nil {
        return err
    }

    // Удаляем элемент
    deleteSql := `DELETE FROM work_order_body1 WHERE id = $1`
    _, err := r.db.Exec(deleteSql, itemID)
    return err
}

// Исполнители (2-я таблица)
func (r *orderRepo) CreatePerformer(performer domain.Performer) (*domain.Order, error) {
    sql := `
        INSERT INTO work_order_body2 
        (work_order_id, employee_id, ktu, amount) 
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
        performer.OrderID,
        performer.EmployeeID,
        performer.KTU,
        performer.Amount,
    ).Scan(&insertedId); err != nil {
        return nil, err
    }

    performer.ID = insertedId
    return r.OrderById(performer.OrderID)
}

func (r *orderRepo) DeletePerformer(performerID int) error {
    // Получаем work_order_id для обновления
    var workOrderID int
    sql := `SELECT work_order_id FROM work_order_body2 WHERE id = $1`
    if err := r.db.QueryRow(sql, performerID).Scan(&workOrderID); err != nil {
        return err
    }

    // Удаляем исполнителя
    deleteSql := `DELETE FROM work_order_body2 WHERE id = $1`
    _, err := r.db.Exec(deleteSql, performerID)
    return err
}