package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util"
)

type Domains struct {
	domains []*backend.Domain
}

func (b *Domains) ListDomains() (domains []*backend.Domain, err error) {
	domains = b.domains
	return
}

func (b *Domains) GetDomainByName(name string) (*backend.Domain, error) {
	for _, d := range b.domains {
		if d.Name == name {
			return d, nil
		}
	}
	return nil, errors.New("No such domain")
}

func (b *Domains) InsertDomain(domain *backend.Domain) (*backend.Domain, error) {
	domain.ID = util.GenerateId()
	b.domains = append(b.domains, domain)
	return domain, nil
}

func NewDomains() backend.DomainsBackend {
	return &Domains{}
}
