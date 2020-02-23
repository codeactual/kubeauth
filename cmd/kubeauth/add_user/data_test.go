package add_user_test

import "github.com/codeactual/kubeauth/internal/testkit"

func SecretName(prefix string) string {
	return prefix + testkit.SecretNameSuffix
}

func CertData() []byte {
	return []byte(testkit.Prefix + "-cert-data")
}

func TokenData() []byte {
	return []byte(testkit.Prefix + "-token-data")
}
