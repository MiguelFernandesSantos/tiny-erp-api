package service

import (
	"api-produtos/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CriptografiaService struct {
}

type AcessoTokenClaims struct {
	CodigoUsuario int64  `json:"codigo_usuario"`
	Email         string `json:"email"`
	Nome          string `json:"nome"`
	Perfil        int64  `json:"perfil"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	CodigoUsuario int64  `json:"codigo_usuario"`
	Email         string `json:"email"`
	jwt.RegisteredClaims
}

func InstanciarCriptografiaService() *CriptografiaService {
	return &CriptografiaService{}
}

func (c *CriptografiaService) CriptografarSenha(senha string) (*string, error) {
	hash, erro := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if erro != nil {
		return nil, erro
	}
	senhaCriptografada := string(hash)
	return &senhaCriptografada, nil
}

func (c *CriptografiaService) CompararHash(senha *string, segundoHash *string) error {
	erro := bcrypt.CompareHashAndPassword([]byte(*segundoHash), []byte(*senha))
	if erro != nil {
		return erro
	}
	return nil
}

func (c *CriptografiaService) GerarTokenAcesso(usuarioAutenticado *model.UsuarioAutenticado) (*string, error) {
	claims := AcessoTokenClaims{
		CodigoUsuario: *usuarioAutenticado.Codigo,
		Email:         *usuarioAutenticado.Email,
		Nome:          *usuarioAutenticado.Nome,
		Perfil:        *usuarioAutenticado.Perfil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tiny-erp",
			Subject:   "acess-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenAssinado, erro := c.AssinarToken(token)
	if erro != nil {
		return nil, erro
	}

	return &tokenAssinado, nil
}

func (c *CriptografiaService) GerarRefreshToken(usuarioAutenticado *model.UsuarioAutenticado) (*string, error) {
	claims := RefreshTokenClaims{
		CodigoUsuario: *usuarioAutenticado.Codigo,
		Email:         *usuarioAutenticado.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tiny-erp",
			Subject:   "refresh-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenAssinado, erro := c.AssinarToken(token)
	if erro != nil {
		return nil, erro
	}

	return &tokenAssinado, nil
}

func (c *CriptografiaService) AssinarToken(token *jwt.Token) (string, error) {
	tokenAssinado, erro := token.SignedString([]byte("sua_chave_secreta_aqui"))
	if erro != nil {
		return "", erro
	}
	return tokenAssinado, nil
}
