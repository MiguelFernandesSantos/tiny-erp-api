package service

import (
	"api-produtos/dao/produto"
	log "api-produtos/internal"
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"api-produtos/service/api"
	"context"
	"strconv"
)

type ProdutoService struct {
	produtoDao                 *produto_dao.ProdutoDao
	produtoApiEmbeddingService *api.ProdutoApiEmbeddingService
	departamentoService        *DepartamentoService
	transacaoService           conexao_mysql.TransacaoService
}

func InstanciarProdutoService() *ProdutoService {
	return &ProdutoService{}
}

func (s *ProdutoService) Inicializar(
	produtoDao *produto_dao.ProdutoDao,
	produtoApiEmbeddingService *api.ProdutoApiEmbeddingService,
	departamentoService *DepartamentoService,
	transacaoService conexao_mysql.TransacaoService,
) {
	s.produtoDao = produtoDao
	s.produtoApiEmbeddingService = produtoApiEmbeddingService
	s.transacaoService = transacaoService
	s.departamentoService = departamentoService
}

func (s *ProdutoService) ObterProdutosPaginado(contexto context.Context, pesquisa string, ultimoProduto *int64) ([]*model.ProdutoResumido, error) {
	log.Start(contexto, "Iniciando busca paginada de produtos")
	produtos, erro := s.produtoDao.ObterProdutosPaginado(contexto, pesquisa, ultimoProduto)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao obter os produtos paginado", &erro)
		return nil, erro
	}

	log.Success(contexto, "Busca paginada de produtos realizada com sucesso")
	return produtos, nil
}

func (s *ProdutoService) AdicionarProduto(contexto context.Context, produto *model.Produto) error {
	log.Info(contexto, "Iniciando cadastro de um novo produto")
	if produto.Plu != nil {
		erro := s.verificarDisponibilidadePlu(contexto, produto.Plu)
		if erro != nil {
			log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do PLU", &erro)
			return erro
		}
	}

	log.Info(contexto, "Verificando existencia do departamento e da secao")
	erro := s.departamentoService.VerificarExistenciaDepartamentoESecao(contexto, produto.CodigoDepartamento, produto.CodigoSecao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um problema ao verificar a existencia do departamento e da secao", &erro)
		return erro
	}

	erro = s.transacaoService.ExecutarTransacao(func(transacao *conexao_mysql.Transacao) error {
		erro = s.produtoDao.AdicionarProduto(contexto, produto, transacao)
		if erro != nil {
			log.Exception(contexto, "Ocorreu um erro ao tentar adicionar o novo produto no banco de dados", &erro)
			return erro
		}

		erro = s.produtoApiEmbeddingService.AdicionarProduto(contexto, produto)
		if erro != nil {
			log.Exception(contexto, "Ocorreu um erro ao tentar adicionar o produto via API", &erro)
			return erro
		}

		// TODO imagens

		return nil
	})
	if erro != nil {
		return erro
	}

	log.Success(contexto, "ProdutoResumido adicionado com sucesso")
	return nil

}

func (s *ProdutoService) verificarDisponibilidadePlu(contexto context.Context, plu *int64) error {
	log.Info(contexto, "Verificando disponibilidade do PLU "+strconv.FormatInt(*plu, 10))
	erro := s.produtoDao.VerificarDisponibilidadePLU(contexto, plu)
	if erro != nil {
		return erro
	}
	log.Info(contexto, "Verificacao da disponibilidade do PLU "+strconv.FormatInt(*plu, 10)+" realizada com sucesso")
	return nil
}

func (s *ProdutoService) ObterProdutoEspecifico(contexto context.Context, plu int64) (*model.Produto, error) {
	log.Start(contexto, "Iniciando busca do produto de plu "+strconv.FormatInt(plu, 10))
	produto, erro := s.produtoDao.ObterProdutoEspecifico(contexto, plu)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro na busca do produto de plu "+strconv.FormatInt(plu, 10), &erro)
		return nil, erro
	}

	log.Success(contexto, "Busca do produto de plu "+strconv.FormatInt(plu, 10)+" realizada com sucesso")
	return produto, nil
}

func (s *ProdutoService) RemoverProduto(contexto context.Context, plu int64) error {
	log.Start(contexto, "Iniciando exclusao do produto de PLU "+strconv.FormatInt(plu, 10))

	erro := s.transacaoService.ExecutarTransacao(func(transacao *conexao_mysql.Transacao) error {
		erro := s.produtoDao.RemoverProduto(contexto, plu)
		if erro != nil {
			log.Exception(contexto, "Ocorreu um erro ao tentar remover o produto do banco de dados", &erro)
			return erro
		}

		erro = s.produtoApiEmbeddingService.RemoverProduto(contexto, plu)
		if erro != nil {
			log.Exception(contexto, "Ocorreu um erro ao tentar remover o produto via API", &erro)
			return erro
		}
		return nil
	})
	if erro != nil {
		return erro
	}

	log.Success(contexto, "Remocao do produto de plu "+strconv.FormatInt(plu, 10)+" realizada com sucesso")
	return nil
}

func (s *ProdutoService) ObterSemelhantes(contexto context.Context, plu int64) ([]*model.ProdutoResumido, error) {
	log.Start(contexto, "Iniciando busca de produtos semelhantes ao produto de plu "+strconv.FormatInt(plu, 10))

	plusSemelhantes, erro := s.produtoApiEmbeddingService.BuscarProdutosSemelhantes(contexto, plu)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao buscar na API os produto semelhantes", &erro)
		return nil, erro
	}

	if len(plusSemelhantes) == 0 {
		log.Info(contexto, "Nenhum produto semelhante foi encontrado para o produto de plu "+strconv.FormatInt(plu, 10))
		return []*model.ProdutoResumido{}, nil
	}

	produtos, erro := s.produtoDao.BuscarProdutosEspecificos(contexto, plusSemelhantes)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao buscar aos produtos semelhantes no banco de dados", &erro)
		return nil, erro
	}

	log.Success(contexto, "Busca de produtos semelhantes ao produto de plu "+strconv.FormatInt(plu, 10)+" buscada com sucesso")
	return produtos, nil
}
