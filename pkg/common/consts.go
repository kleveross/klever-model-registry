package common

import (
	"github.com/spf13/viper"
)

const (
	// DefaultModelJobNamespace is the default namespace for modeljob
	DefaultModelJobNamespace = "default"

	// ORMBDomainEnvKey is the domain of ORMB
	ORMBDomainEnvKey = "SERVER_ORMB_DOMAIN"
	// ORMBUsernameEnvkey is the username of ORMB
	ORMBUsernameEnvkey = "SERVER_ORMB_USERNAME"
	// ORMBPasswordEnvKey is the password of ORMB
	ORMBPasswordEnvKey = "SERVER_ORMB_PASSWORD"

	// ResourceNameLabelKey is resource name in labels
	ResourceNameLabelKey = "resource_name"
)

var (
	ORMBDomain   string
	ORMBUserName string
	ORMBPassword string
)

func init() {
	viper.AutomaticEnv()
	ORMBDomain = viper.GetString(ORMBDomainEnvKey)
	ORMBUserName = viper.GetString(ORMBUsernameEnvkey)
	ORMBPassword = viper.GetString(ORMBPasswordEnvKey)
}
