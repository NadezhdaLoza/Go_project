package implementation

import (
	"database/sql"

	"github.com/nadezhda/documents/services/product/backend/repo"
	"github.com/nadezhda/documents/services/product/domain"
)

// Реализация репозитория товаров
type repository struct {
	db *sql.DB
}

// Возвращает репозиторий товаров
func NewProductRepo(db *sql.DB) repo.ProductRepo {
	return &repository{
		db: db,
	}
}

// Добавляет новый товар
func (r *repository) CreateProduct(p domain.Product) (*domain.Product, error) {
	// SQL-запрос для добавления товара с возвратом id
	sql := `
                INSERT INTO products (name, price, unit_id)
                VALUES ($1, $2, $3) RETURNING id
                `

	// Подготовка запроса
	result, err := r.db.Prepare(sql)
	if err != nil {
		// Возвратnil и ошибки, если запрос не может быть выполнен
		return nil, err
	}
	// Освобождение ресурсов после выполнения функции
	defer result.Close()

	// Выполнение запроса и сканирование результата (Id созданного товара)
	var insertedId int
	if err := result.QueryRow(p.Name, p.Price, p.UnitId).Scan(&insertedId); err != nil {
		// объявление и проверка ошибки внутри блока if
		return nil, err
	}

	// Возврат нового товара
	return r.ProductById(insertedId)
}

// GetProducts выполняет запрос в базу данных для получения списка всех товаров.
// Возвращает массив типа domain.Product и ошибку, если она произошла.
func (r *repository) GetProducts() ([]domain.Product, error) {
	// SQL-запрос для получения информации о всех товарах
	sql := `
        SELECT * FROM products ORDER BY products.name ASC
        `

	// Выполнение запроса
	rows, err := r.db.Query(sql)
	if err != nil {
		// Возврат пустого списка и ошибки, если запрос не может быть выполнен
		return nil, err
	}

	// Освобождение ресурсов после выполнения функции
	defer rows.Close()

	// Создание пустого списка товаров
	products := []domain.Product{}

	// Перебор всех строк в ответе на запрос
	for rows.Next() {
		// Создание пустой структуры товара
		product := domain.Product{}
		// Заполнение структуры данными из текущей строки
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.UnitId); err != nil {
			return nil, err
		}
		// Добавление товара в список товаров
		products = append(products, product)
	}

	// Возврат товаров
	return products, nil
}

// Обновляет информацию о товаре в базе данных.
// Принимает параметр p типа domain.Product, который содержит новые данные для товара.
func (r *repository) UpdateProduct(p domain.Product) (*domain.Product, error) {
	// SQL-запрос для обновления цены, единицы измерения и имени товара
	sql := `
                UPDATE products SET price = $1, unit_id = $2, name = $3 WHERE id = $4
        `
	// Выполнение запроса
	rows, err := r.db.Query(sql, p.Price, p.UnitId, p.Name, p.Id)
	if err != nil {
		return nil, err
	}

	// Закрыть соединение, чтобы осуществить следующий запрос
	rows.Close()

	// Вернуть обновленный товар
	return r.ProductById(p.Id)
}

// DeleteProduct удаляеттовар из базы данных по его ID.
// Idтовара передается в качестве параметра.
func (r *repository) DeleteProduct(id int) error {
	sql := `
                DELETE FROM products WHERE id = $1
        `
	// через _ обозначается переменная, которая не будет
	// использоваться в коде
	_, err := r.db.Exec(sql, id)
	return err
}

func (r *repository) ProductById(id int) (*domain.Product, error) {
	// SQL-запрос для получения данных товара
	getProductSql := `
                SELECT * FROM products WHERE products.id = $1    
        `

	// Инициализация переменной
	product := domain.Product{}

	// Выполнение запроса и сканирование результата (данные товара)
	if err := r.db.QueryRow(getProductSql, id).Scan(
		&product.Id,
		&product.Name,
		&product.Price,
		&product.UnitId); err != nil {

		return nil, err
	}
	// Возврат указателя на товар
	return &product, nil
}