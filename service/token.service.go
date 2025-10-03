package service

import (
	log "api-produtos/internal"
	"api-produtos/model"
	"context"
)

type TokenService struct {
	criptografiaService *CriptografiaService
}

func InstanciarTokenService() *TokenService {
	return &TokenService{}
}

func (t *TokenService) Inicializar(criptografiaService *CriptografiaService) {
	t.criptografiaService = criptografiaService
}

func (t *TokenService) GerarTokensAutenticacao(contexto context.Context, usuarioAutenticado *model.UsuarioAutenticado) (*string, *string, error) {
	log.Start(contexto, "Iniciando geração dos tokens de autenticação do usuário")
	token, erro := t.criptografiaService.GerarTokenAcesso(usuarioAutenticado)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar gerar o token de autenticação do usuário", &erro)
		return nil, nil, erro
	}

	refreshToken, erro := t.criptografiaService.GerarRefreshToken(usuarioAutenticado)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar gerar o refresh token do usuário", &erro)
		return nil, nil, erro
	}

	log.Success(contexto, "Geração dos tokens de autenticação do usuário realizada com sucesso")
	return token, refreshToken, nil
}
