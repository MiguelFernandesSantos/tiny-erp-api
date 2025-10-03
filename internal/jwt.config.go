package internal

import (
	"api-produtos/internal/state"

	"github.com/golang-jwt/jwt/v5"
)

func CarregarChavesJWT(chavePrivada []byte, chavePublica []byte) error {
	chavePrivadaJWT, erro := jwt.ParseRSAPrivateKeyFromPEM(chavePrivada)
	if erro != nil {
		return erro
	}

	chavePublicaJWT, erro := jwt.ParseRSAPublicKeyFromPEM(chavePublica)
	if erro != nil {
		return erro
	}

	state.CHAVE_PRIVADA_JWT = chavePrivadaJWT
	state.CHAVE_PUBLICA_JWT = chavePublicaJWT

	return nil
}
