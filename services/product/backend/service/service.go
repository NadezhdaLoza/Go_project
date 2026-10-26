package service

import"github.com/nadezhda/documents/services/product/domain"

// Service определяет интерфейс, который предоставляет методы для работы с товарами и единицами измерения.
// - CreateProduct создает новый товар с именем "name", ценой "price" и единицей измерения с id "unitId".
// - GetProducts возвращает список всех товаров.
// - UpdateProduct обновляет информацию о товаре.
// - DeleteProduct удаляет товар из базы данных по id.
// - CreateUnit создает новую единицу измерения с именем "name".
// - GetUnits возвращает список всех единиц измерения.
// - UpdateUnit обновляет информацию о единице измерения.
// - DeleteUnit удаляет единицу измерения из базы данных по id.
    type ProductService interface{
CreateProduct(name string, price int, unitId int)(*domain.Product, error)
GetProducts()([]domain.Product, error)
UpdateProduct(domain.Product)(*domain.Product, error)
DeleteProduct(int) error
ProductById(int)(*domain.Product, error)
CreateUnit(string)(*domain.Unit, error)
GetUnits()([]domain.Unit, error)
UpdateUnit(domain.Unit)(*domain.Unit, error)
DeleteUnit(int) error
UnitById(int)(*domain.Unit, error)
}
