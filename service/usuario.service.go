package service

import (
	usuario_dao "api-produtos/dao/usuario"
	log "api-produtos/internal"
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"context"
	"strconv"
)

type UsuarioService struct {
	usuarioDao       *usuario_dao.UsuarioDao
	transacaoService conexao_mysql.TransacaoService
}

func InstanciarUsuarioService() *UsuarioService {
	return &UsuarioService{}
}

func (s *UsuarioService) Inicializar(usuarioDao *usuario_dao.UsuarioDao, transacaoService conexao_mysql.TransacaoService) {
	s.usuarioDao = usuarioDao
	s.transacaoService = transacaoService
}

func (s *UsuarioService) ObterUsuariosPaginado(context context.Context, pesquisa string, ultimoUsuario *int64) ([]*model.UsuarioResumido, error) {
	log.Start(context, "Iniciando busca de usuários paginada")

	usuarios, err := s.usuarioDao.ObterUsuariosPaginado(context, pesquisa, ultimoUsuario)
	if err != nil {
		log.Exception(context, "Ocorreu um erro ao obter os usuários paginados", &err)
		return nil, err
	}

	log.Success(context, "Busca de usuários paginada realizada com sucesso")
	return usuarios, nil
}

func (s *UsuarioService) AdicionarUsuario(context context.Context, usuario *model.Usuario) error {
	log.Start(context, "Iniciando cadastro de um novo usuário")

	erro := s.usuarioDao.VerificarDisponibilidadeEmail(context, usuario.Email)
	if erro != nil {
		log.Exception(context, "Ocorreu um problema ao verificar a existência do email do usuário", &erro)
		return erro
	}

	erro = s.usuarioDao.AdicionarUsuario(context, usuario)
	if erro != nil {
		log.Exception(context, "Ocorreu um erro ao tentar adicionar o novo usuário no banco de dados", &erro)
		return erro
	}

	log.Success(context, "Usuário cadastrado com sucesso")
	return nil
}

func (s *UsuarioService) ObterUsuarioEspecifico(context context.Context, ultimoUsuario *int64) (*model.Usuario, error) {
	log.Start(context, "Iniciando busca de um usuário específico")

	usuario, erro := s.usuarioDao.ObterUsuarioEspecifico(context, ultimoUsuario)
	if erro != nil {
		log.Exception(context, "Ocorreu um erro ao obter o usuário específico", &erro)
		return nil, erro
	}

	log.Success(context, "Busca de usuário específico realizada com sucesso")
	return usuario, nil
}

func (s *UsuarioService) DesativarUsuarioEspecifico(context context.Context, ultimoUsuario *int64) error {
	log.Start(context, "Iniciando desativação do usuario de id "+strconv.FormatInt(*ultimoUsuario, 10))

	erro := s.usuarioDao.DesativarUsuarioEspecifico(context, ultimoUsuario)
	if erro != nil {
		log.Exception(context, "Ocorreu um erro ao desativar o usuário específico", &erro)
		return erro
	}

	log.Success(context, "Desativação de usuário específico realizada com sucesso")
	return nil
}
