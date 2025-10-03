package service

import (
	log "api-produtos/internal"
	"api-produtos/model"
	"context"
)

type AutenticacaoService struct {
	usuarioService      *UsuarioService
	tokenService        *TokenService
	criptografiaService *CriptografiaService
}

func InstanciarAutenticacaoService() *AutenticacaoService {
	return &AutenticacaoService{}
}

func (s *AutenticacaoService) Inicializar(usuarioService *UsuarioService, tokenService *TokenService, criptografiaService *CriptografiaService) {
	s.usuarioService = usuarioService
	s.tokenService = tokenService
	s.criptografiaService = criptografiaService
}

func (s *AutenticacaoService) Logar(contexto context.Context, email *string, senha *string) (*model.UsuarioAutenticado, error) {
	log.Start(contexto, "Iniciando autenticação do usuário de email: "+*email)
	senhaBanco, erro := s.usuarioService.ObterSenhaUsuario(contexto, email)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar obter a senha do usuário", &erro)
		return nil, erro
	}

	erro = s.criptografiaService.CompararHash(senha, senhaBanco)
	if erro != nil {
		log.Error(contexto, "Senha incorreta para o usuário de email: "+*email)
		return nil, model.SenhaIncorreta
	}

	dadosAutenticacao, erro := s.usuarioService.ObterDadosAutenticacaoUsuario(contexto, email)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar obter os dados de autenticação do usuário", &erro)
		return nil, erro
	}

	token, refreshToken, erro := s.tokenService.GerarTokensAutenticacao(contexto, dadosAutenticacao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar gerar o token de autenticação do usuário", &erro)
		return nil, erro
	}

	dadosAutenticacao.Token = token
	dadosAutenticacao.RefreshToken = refreshToken

	log.Success(contexto, "Autenticação do usuário realizada com sucesso")
	return dadosAutenticacao, nil
}
