package msgraph

import (
	"fmt"
	"testing"
)

func GetTestDomain(t *testing.T) Domain {
	t.Helper()
	domains, err := graphClient.ListDomains()
	if err != nil {
		t.Fatalf("Cannot GraphClient.ListDomains(): %v", err)
	}
	domainTest, err := domains.GetById(msGraphExistingDomainId)
	if err != nil {
		t.Fatalf("Cannot domains.GetById(%v): %v", msGraphExistingDomainId, err)
	}
	return domainTest
}

func TestDomain_String(t *testing.T) {
	testDomain := GetTestDomain(t)

	tests := []struct {
		name string
		g    Domain
		want string
	}{
		{
			name: "Test All Domains",
			g:    testDomain,
			want: fmt.Sprintf("Domain(ID: \"%v\", AuthenticationType: \"%v\" IsAdminManaged: \"%v\", IsDefault: \"%v\", IsInitial: \"%v\", IsRoot: \"%v\", IsVerified: \"%v\", SupportedServices: \"%v\", PasswordValidityPeriodInDays: \"%v\", PasswordNotificationWindowInDays: \"%v\", DirectAPIConnection: %v)",
				testDomain.ID, testDomain.AuthenticationType, testDomain.IsAdminManaged, testDomain.IsDefault, testDomain.IsInitial, testDomain.IsRoot, testDomain.IsVerified, testDomain.SupportedServices, testDomain.PasswordValidityPeriodInDays, testDomain.PasswordNotificationWindowInDays, testDomain.graphClient != nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("Domain.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
