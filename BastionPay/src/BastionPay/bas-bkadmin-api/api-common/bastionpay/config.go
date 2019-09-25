package bastionpay

type (
	Credentials struct {
		PrivateKey          string
		PublicKey           string
		BastionPayPublicKey string
	}

	Config struct {
		*Credentials
		ApiKey string
	}
)

func NewCredentials(privateKey string, publicKey string, bastionPayPublicKey string) *Credentials {
	return &Credentials{privateKey, publicKey, bastionPayPublicKey}
}

func NewConfig(apikey string, credentials *Credentials) *Config {
	return &Config{Credentials: credentials, ApiKey: apikey}
}
