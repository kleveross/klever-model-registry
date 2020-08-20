package common

const (
	// DefaultModelJobNamespace is the default namespace for modeljob
	DefaultModelJobNamespace = "default"

	// ORMBDomainEnvKey is the domain of ORMB
	ORMBDomainEnvKey = "ORMB_DOMAIN"
	// ORMBUsernameEnvkey is the username of ORMB
	ORMBUsernameEnvkey = "ORMB_USERNAME"
	// ORMBPasswordEnvKey is the password of ORMB
	ORMBPasswordEnvKey = "ORMB_PASSWORD"

	// ResourceNameLabelKey is resource name in labels
	ResourceNameLabelKey = "resource_name"
)

var (
	ORMBDomain   string
	ORMBUserName string
	ORMBPassword string
)
