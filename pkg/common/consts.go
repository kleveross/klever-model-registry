package common

const (
	// ORMBDomainEnvKey is the domain of ORMB
	ORMBDomainEnvKey = "SERVER_ORMB_DOMAIN"
	// ORMBUsernameEnvkey is the username of ORMB
	ORMBUsernameEnvkey = "SERVER_ORMB_USERNAME"
	// ORMBPasswordEnvKey is the password of ORMB
	ORMBPasswordEnvKey = "SERVER_ORMB_PASSWORD"
)

var (
	ORMBDomain   string
	ORMBUserName string
	ORMBPassword string
)
