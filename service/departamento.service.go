package service

import (
	departamento_dao "api-produtos/dao/departamento"
	log "api-produtos/internal"
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"context"
	"strconv"
)

type DepartamentoService struct {
	departamentoDao  *departamento_dao.DepartamentoDao
	transacaoService conexao_mysql.TransacaoService
}

func InstanciarDepartamentoService() *DepartamentoService {
	return &DepartamentoService{}
}

func (s *DepartamentoService) Inicializar(
	departamentoDao *departamento_dao.DepartamentoDao,
	transacaoService conexao_mysql.TransacaoService,
) {
	s.departamentoDao = departamentoDao
	s.transacaoService = transacaoService
}

func (s *DepartamentoService) ObterDepartamentosPaginado(contexto context.Context, pesquisa string, ultimoDepartamento *int64) ([]*model.Departamento, error) {
	log.Start(contexto, "Iniciando busca paginada de departamentos")
	departamentos, erro := s.departamentoDao.ObterDepartamentosPaginado(contexto, pesquisa, ultimoDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao obter os departamentos paginado", &erro)
		return nil, erro
	}

	log.Success(contexto, "Busca paginada de departamentos realizada com sucesso")
	return departamentos, nil
}

func (s *DepartamentoService) AdicionarDepartamento(contexto context.Context, departamento *model.Departamento) error {
	log.Info(contexto, "Iniciando cadastro de um novo departamento")

	erro := s.departamentoDao.VerificarDisponibilidadeNome(contexto, departamento.Nome)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do nome do departamento", &erro)
		return erro
	}

	erro = s.departamentoDao.AdicionarDepartamento(contexto, departamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar adicionar o novo departamento no banco de dados", &erro)
		return erro
	}

	log.Success(contexto, "Departamento cadastrado com sucesso")
	return nil
}

func (s *DepartamentoService) ObterDepartamentoEspecifico(contexto context.Context, numeroDepartamento int64) (*model.Departamento, error) {
	log.Start(contexto, "Iniciando busca de departamento de número "+strconv.FormatInt(numeroDepartamento, 10))
	departamento, erro := s.departamentoDao.ObterDepartamentoEspecifico(contexto, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao obter o departamento", &erro)
		return nil, erro
	}

	log.Success(contexto, "Busca de departamento realizada com sucesso")
	return departamento, nil
}

func (s *DepartamentoService) RemoverDepartamentoEspecifico(contexto context.Context, numeroDepartamento *int64) error {
	log.Info(contexto, "Iniciando remoção do departamento de número "+strconv.FormatInt(*numeroDepartamento, 10))

	erro := s.departamentoDao.VerificarSePossuiProdutosAssociados(contexto, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar se o departamento possui produtos associados", &erro)
		return erro
	}

	erro = s.departamentoDao.RemoverDepartamentoEspecifico(contexto, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar remover o departamento do banco de dados", &erro)
		return erro
	}

	log.Success(contexto, "Departamento removido com sucesso")
	return nil
}

func (s *DepartamentoService) AdicionarSecaoAoDepartamento(contexto context.Context, secao *model.Secao) error {
	log.Info(contexto, "Iniciando cadastro de uma nova seção no departamento "+strconv.FormatInt(*secao.NumeroDepartamento, 10))
	erro := s.departamentoDao.VerificarExistenciaDepartamento(contexto, secao.NumeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do departamento", &erro)
		return erro
	}

	erro = s.departamentoDao.VerificarDisponibilidadeNomeSecao(contexto, secao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do nome da seção", &erro)
		return erro
	}

	erro = s.departamentoDao.AdicionarSecaoAoDepartamento(contexto, secao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar adicionar a nova seção no banco de dados", &erro)
		return erro
	}

	log.Success(contexto, "Seção adicionada com sucesso")
	return nil
}

func (s *DepartamentoService) ObterSecaoEspecifica(contexto context.Context, numeroSecao int64, numeroDepartamento int64) (*model.Secao, error) {
	log.Start(contexto, "Iniciando busca da seção de número "+strconv.FormatInt(numeroSecao, 10)+" do departamento "+strconv.FormatInt(numeroDepartamento, 10))
	secao, erro := s.departamentoDao.ObterSecaoEspecifica(contexto, numeroSecao, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao obter a seção", &erro)
		return nil, erro
	}
	log.Success(contexto, "Busca da seção realizada com sucesso")
	return secao, nil
}

func (s *DepartamentoService) RemoverSecaoEspecifica(contexto context.Context, numeroSecao int64, numeroDepartamento int64) error {
	log.Info(contexto, "Iniciando remoção da seção de número "+strconv.FormatInt(numeroSecao, 10)+" do departamento "+strconv.FormatInt(numeroDepartamento, 10))

	erro := s.departamentoDao.VerificarSePossuiProdutosAssociadosASecao(contexto, numeroSecao, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar se a seção possui produtos associados", &erro)
		return erro
	}

	erro = s.departamentoDao.RemoverSecaoEspecifica(contexto, numeroSecao, numeroDepartamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar remover a seção do banco de dados", &erro)
		return erro
	}

	log.Success(contexto, "Seção removida com sucesso")
	return nil
}

func (s *DepartamentoService) VerificarExistenciaDepartamentoESecao(contexto context.Context, departamento *int64, secao *int64) error {
	log.Start(contexto, "Verificando existencia do departamento de número "+strconv.FormatInt(*departamento, 10)+" e da seção de número "+strconv.FormatInt(*secao, 10))
	erro := s.departamentoDao.VerificarExistenciaDepartamento(contexto, departamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do departamento", &erro)
		return erro
	}

	erro = s.departamentoDao.VerificarExistenciaSecao(contexto, secao, departamento)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia da seção", &erro)
		return erro
	}

	log.Info(contexto, "Verificação de existencia do departamento e da seção realizada com sucesso")
	return nil
}
