package repo

import "github.com/nadezhda/documents/services/proxy/domain"

type ProxyRepo interface {
	CreateProxyHeader(domain.ProxyHeader) (*domain.ProxyHeader, error)
	CreateProxyBodyItem(item domain.ProxyBodyItem) (*domain.Proxy, error)
	GetProxies() ([]domain.ProxyHeader, error)
	ProxyById(id int) (*domain.Proxy, error)
	UpdateProxyHeader(domain.ProxyHeader) (*domain.Proxy, error)
	DeleteProxy(id int) error
	DeleteProxyBodyItem(id int) error
}
