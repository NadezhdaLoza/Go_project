package repo

import "github.com/nadezhda/documents/services/organization/domain"

type OrganizationRepo interface {
	GetOrganizations() ([]domain.Organization, error)
	CreateOrganization(domain.Organization) (*domain.Organization, error)
	UpdateOrganization(domain.Organization) (*domain.Organization, error)
	DeleteOrganization(int) error
	OrganizationById(int) (*domain.Organization, error)
}
