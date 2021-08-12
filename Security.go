package msgraph

import (
	"net"
	"time"
)

// Alert represents a security alert.
type Alert struct {
	ActivityGroupName    string                    `json:"activityGroupName"`
	AssignedTo           string                    `json:"assignedTo"`
	AzureSubscriptionID  string                    `json:"azureSubscriptionId"`
	AzureTenantID        string                    `json:"azureTenantId"`
	Category             string                    `json:"category"`
	ClosedDateTime       time.Time                 `json:"closedDateTime"`
	CloudAppStates       []CloudAppSecurityState   `json:"cloudAppStates"`
	Comments             []string                  `json:"comments"`
	Confidence           int32                     `json:"confidence"`
	CreatedDateTime      time.Time                 `json:"createdDateTime"`
	Description          string                    `json:"description"`
	DetectionIDs         []string                  `json:"detectionIds"`
	EventDateTime        time.Time                 `json:"eventDateTime"`
	Feedback             string                    `json:"feedback"`
	FileStates           []FileSecurityState       `json:"fileStates"`
	HostStates           []HostSecurityState       `json:"hostStates"`
	ID                   string                    `json:"id"`
	IncidentIDs          []string                  `json:"incidentIds"`
	LastModifiedDateTime time.Time                 `json:"lastModifiedDateTime"`
	MalwareStates        []MalwareState            `json:"malwareStates"`
	NetworkConnections   []NetworkConnection       `json:"networkConnections"`
	Processes            []Process                 `json:"processes"`
	RecommendedActions   []string                  `json:"recommendedActions"`
	RegistryKeyStates    []RegistryKeyState        `json:"registryKeyStates"`
	SecurityResources    []SecurityResource        `json:"securityResources"`
	Severity             string                    `json:"severity"`
	SourceMaterials      []string                  `json:"sourceMaterials"`
	Status               string                    `json:"status"`
	Tags                 []string                  `json:"tags"`
	Title                string                    `json:"title"`
	Triggers             []AlertTrigger            `json:"triggers"`
	UserStates           []UserSecurityState       `json:"userStates"`
	VendorInformation    SecurityVendorInformation `json:"vendorInformation"`
	VulnerabilityStates  []VulnerabilityState      `json:"vulnerabilityStates"`
}

// CloudAppSecurityState contains stateful information about a cloud application related to an alert.
type CloudAppSecurityState struct {
	DestinationServiceIP   net.IP `json:"destinationServiceIp"`
	DestinationServiceName string `json:"destinationServiceName"`
	RiskScore              string `json:"riskScore"`
}

// FileSecurityState contains information about a file (not process) related to an alert.
type FileSecurityState struct {
	FileHash  FileHash `json:"fileHash"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	RiskScore string   `json:"riskScore"`
}

// FileHash contains hash information related to a file.
type FileHash struct {
	HashType  string `json:"hashType"`
	HashValue string `json:"hashValue"`
}

// HostSecurityState contains information about a host (computer, device, etc.) related to an alert.
type HostSecurityState struct {
	FQDN                      string `json:"fqdn"`
	IsAzureAADJoined          bool   `json:"isAzureAadJoined"`
	IsAzurAADRegistered       bool   `json:"isAzureAadRegistered"`
	IsHybridAzureDomainJoined bool   `json:"isHybridAzureDomainJoined"`
	NetBiosName               string `json:"netBiosName"`
	OS                        string `json:"os"`
	PrivateIPAddress          net.IP `json:"privateIpAddress"`
	PublicIPAddress           net.IP `json:"publicIpAddress"`
	RiskScore                 string `json:"riskScore"`
}

// MalwareState contains information about a malware entity.
type MalwareState struct {
	Category   string `json:"category"`
	Family     string `json:"family"`
	Name       string `json:"name"`
	Severity   string `json:"severity"`
	WasRunning bool   `json:"wasRunning"`
}

// NetworkConnection contains stateful information describing a network connection related to an alert.
type NetworkConnection struct {
	ApplicationName          string    `json:"applicationName"`
	DestinationAddress       net.IP    `json:"destinationAddress"`
	DestinationLocation      string    `json:"destinationLocation"`
	DestinationDomain        string    `json:"destinationDomain"`
	DestinationPort          string    `json:"destinationPort"` // spec calls it a string, not a number
	DestinationURL           string    `json:"destinationUrl"`
	Direction                string    `json:"direction"`
	DomainRegisteredDateTime time.Time `json:"domainRegisteredDateTime"`
	LocalDNSName             string    `json:"localDnsName"`
	NATDestinationAddress    net.IP    `json:"natDestinationAddress"`
	NATDestinationPort       string    `json:"natDestinationPort"`
	NATSourceAddress         net.IP    `json:"natSourceAddress"`
	NATSourcePort            string    `json:"natSourcePort"`
	Protocol                 string    `json:"protocol"`
	RiskScore                string    `json:"riskScore"`
	SourceAddress            net.IP    `json:"sourceAddress"`
	SourceLocation           string    `json:"sourceLocation"`
	SourcePort               string    `json:"sourcePort"`
	Status                   string    `json:"status"`
	URLParameters            string    `json:"urlParameters"`
}

// Process describes a process related to an alert.
type Process struct {
	AccountName                  string    `json:"accountName"`
	CommandLine                  string    `json:"commandLine"`
	CreatedDateTime              time.Time `json:"createdDateTime"` // translated
	FileHash                     FileHash  `json:"fileHash"`
	IntegrityLevel               string    `json:"integrityLevel"`
	IsElevated                   bool      `json:"isElevated"`
	Name                         string    `json:"name"`
	ParentProcessCreatedDateTime time.Time `json:"parentProcessCreatedDateTime"` // translated
	ParentProcessID              int32     `json:"parentProcessId"`
	ParentProcessName            string    `json:"parentProcessName"`
	Path                         string    `json:"path"`
	ProcessID                    int32     `json:"processId"`
}

// RegistryKeyState contains information about registry key changes related to an alert, and about the process which changed the keys.
type RegistryKeyState struct {
	Hive         string `json:"hive"`
	Key          string `json:"key"`
	OldKey       string `json:"oldKey"`
	OldValueData string `json:"oldValueData"`
	OldValueName string `json:"oldValueName"`
	Operation    string `json:"operation"`
	ProcessID    int32  `json:"processId"`
	ValueData    string `json:"valueData"`
	ValueName    string `json:"valueName"`
	ValueType    string `json:"valueType"`
}

// SecurityResource represents resources related to an alert.
type SecurityResource struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
}

// AlertTrigger contains information about a property which triggered an alert detection.
type AlertTrigger struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// UserSecurityState contains stateful information about a user account related to an alert.
type UserSecurityState struct {
	AADUserID                    string    `json:"aadUserId"`
	AccountName                  string    `json:"accountName"`
	DomainName                   string    `json:"domainName"`
	EmailRole                    string    `json:"emailRole"`
	IsVPN                        bool      `json:"isVpn"`
	LogonDateTime                time.Time `json:"logonDateTime"`
	LogonID                      string    `json:"logonId"`
	LogonIP                      net.IP    `json:"logonIp"`
	LogonLocation                string    `json:"logonLocation"`
	LogonType                    string    `json:"logonType"`
	OnPremisesSecurityIdentifier string    `json:"onPremisesSecurityIdentifier"`
	RiskScore                    string    `json:"riskScore"`
	UserAccountType              string    `json:"userAccountType"`
	UserPrincipalName            string    `json:"userPrincipalName"`
}

// SecurityVendorInformation contains details about the vendor of a particular security product.
type SecurityVendorInformation struct {
	Provider        string `json:"provider"`
	ProviderVersion string `json:"providerVersion"`
	SubProvider     string `json:"subProvider"`
	Vendor          string `json:"vendor"`
}

// VulnerabilityState contains information about a particular vulnerability.
type VulnerabilityState struct {
	CVE        string `json:"cve"`
	Severity   string `json:"severity"`
	WasRunning bool   `json:"wasRunning"`
}

// ListAlerts returns a slice of Alert objects from MS Graph's security API. Each Alert represents a security event reported by some component.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
func (g *GraphClient) ListAlerts(opts ...ListQueryOption) ([]Alert, error) {
	resource := "/security/alerts"
	var marsh struct {
		Alerts []Alert `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Alerts, err
}

// SecureScore represents the security score of a tenant for a particular day.
type SecureScore struct {
	ID                       string                    `json:"id"`
	AzureTenantID            string                    `json:"azureTenantId"`
	ActiveUserCount          int32                     `json:"activeUserCount"`
	CreatedDateTime          time.Time                 `json:"createdDateTime"`
	CurrentScore             float64                   `json:"currentScore"`
	EnabledServices          []string                  `json:"enabledServices"`
	LicensedUserCount        int32                     `json:"licensedUserCount"`
	MaxScore                 float64                   `json:"maxScore"`
	AverageComparativeScores []AverageComparativeScore `json:"averageComparativeScores"`
	ControlScores            []ControlScore            `json:"controlScores"`
	VendorInformation        SecurityVendorInformation `json:"vendorInformation"`
}

// AverageComparativeScore describes average scores across a variety of different scopes.
// The Basis field may contain the strings "AllTenants", "TotalSeats", or "IndustryTypes".
type AverageComparativeScore struct {
	Basis        string  `json:"basis"`
	AverageScore float64 `json:"averageScore"`
}

// ControlScore contains a score for a single security control.
type ControlScore struct {
	ControlName     string  `json:"controlName"`
	Score           float64 `json:"score"`
	ControlCategory string  `json:"controlCategory"`
	Description     string  `json:"description"`
}

// ListSecureScores returns a slice of SecureScore objects. Each SecureScore represents
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
// a tenant's security score for a particular day.
func (g *GraphClient) ListSecureScores(opts ...ListQueryOption) ([]SecureScore, error) {
	resource := "/security/secureScores"
	var marsh struct {
		Scores []SecureScore `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Scores, err
}

// SecureScoreControlProfile describes in greater detail the parameters of a given security
// score control.
type SecureScoreControlProfile struct {
	ID                    string                          `json:"id"`
	AzureTenantID         string                          `json:"azureTenantId"`
	ActionType            string                          `json:"actionType"`
	ActionURL             string                          `json:"actionUrl"`
	ControlCategory       string                          `json:"controlCategory"`
	Title                 string                          `json:"title"`
	Deprecated            bool                            `json:"deprecated"`
	ImplementationCost    string                          `json:"implementationCost"`
	LastModifiedDateTime  time.Time                       `json:"lastModifiedDateTime"`
	MaxScore              float64                         `json:"maxScore"`
	Rank                  int32                           `json:"rank"`
	Remediation           string                          `json:"remediation"`
	RemediationImpact     string                          `json:"remediationImpact"`
	Service               string                          `json:"service"`
	Threats               []string                        `json:"threats"`
	Tier                  string                          `json:"tier"`
	UserImpact            string                          `json:"userImpact"`
	ComplianceInformation []ComplianceInformation         `json:"complianceInformation"`
	ControlStateUpdates   []SecureScoreControlStateUpdate `json:"controlStateUpdates"`
	VendorInformation     SecurityVendorInformation       `json:"vendorInformation"`
}

// ComplianceInformation contains compliance data associated with a secure score control.
type ComplianceInformation struct {
	CertificationName     string                 `json:"certificationName"`
	CertificationControls []CertificationControl `json:"certificationControls"`
}

// CertificationControl contains compliance certification data associated with a secure score control.
type CertificationControl struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SecureScoreControlStateUpdate records a particular historical state of the control state
// as updated by the user.
type SecureScoreControlStateUpdate struct {
	AssignedTo      string    `json:"assignedTo"`
	Comment         string    `json:"comment"`
	State           string    `json:"state"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
}

// ListSecureScoreControlProfiles returns a slice of SecureScoreControlProfile objects.
// Each object represents a secure score control profile, which is used when calculating
// a tenant's secure score.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
func (g *GraphClient) ListSecureScoreControlProfiles(opts ...ListQueryOption) ([]SecureScoreControlProfile, error) {
	resource := "/security/secureScoreControlProfiles"
	var marsh struct {
		Profiles []SecureScoreControlProfile `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	return marsh.Profiles, err
}
