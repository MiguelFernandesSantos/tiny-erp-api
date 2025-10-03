package state

import "crypto/rsa"

var (
	URL_API_EMBEDDING      = ""
	UTILIZAR_API_EMBEDDING = false
	CHAVE_PRIVADA_JWT      *rsa.PrivateKey
	CHAVE_PUBLICA_JWT      *rsa.PublicKey
)
