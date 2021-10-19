package msgraph

import (
	"strings"
)

// Domains represents multiple Domains-instances and provides funcs to work with them.
type Domains []Domain

func (d Domains) String() string {
	var domains = make([]string, len(d))
	for i, calendar := range d {
		domains[i] = calendar.String()
	}
	return "Domains(" + strings.Join(domains, " | ") + ")"
}

// setGraphClient sets the GraphClient within that particular instance. Hence it's directly created by GraphClient
func (d Domains) setGraphClient(gC *GraphClient) Domains {
	for i := range d {
		d[i].setGraphClient(gC)
	}
	return d
}

// GetById returns the Domain obj of that array whose ID matches
// the given id. Returns an ErrFindDomain if no domain exists that matches the given
// DisplayName.
func (d Domains) GetById(id string) (Domain, error) {
	for _, domain := range d {
		if domain.ID == id {
			return domain, nil
		}
	}
	return Domain{}, ErrFindDomain
}
