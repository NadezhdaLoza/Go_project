package repo

import(
"github.com/nadezhda/documents/services/product/domain"
)

    type ProductRepo interface{
CreateProduct(domain.Product)(*domain.Product, error)
GetProducts()([]domain.Product, error)
UpdateProduct(domain.Product)(*domain.Product, error)
DeleteProduct(int) error
ProductById(int)(*domain.Product, error)
}
