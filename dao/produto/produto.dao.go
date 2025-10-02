package produto_dao

import (
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"context"
	"strconv"
	"strings"
)

type ProdutoDao struct {
	conexao *conexao_mysql.MySQL
}

func InstanciarProdutoDao(conexao *conexao_mysql.MySQL) *ProdutoDao {
	return &ProdutoDao{
		conexao: conexao,
	}
}

func (d *ProdutoDao) ObterProdutosPaginado(context context.Context, pesquisa string, ultimoProduto *int64) ([]*model.ProdutoResumido, error) {
	pesquisaLike := "%" + pesquisa + "%"

	if ultimoProduto == nil {
		ultimoProduto = new(int64)
		*ultimoProduto = 0
	}

	resultados, erro := d.conexao.Select(QUERY_BUSCAR_PRODUTOS_PAGINADO, pesquisaLike, ultimoProduto)
	if erro != nil {
		return nil, erro
	}

	produtos := make([]*model.ProdutoResumido, 0)
	for _, resultado := range resultados {
		produto := d.mapearProdutoResumido(resultado)
		produtos = append(produtos, produto)
	}

	return produtos, erro
}

func (d *ProdutoDao) VerificarDisponibilidadePLU(contexto context.Context, plu *int64) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_DISPONIBILIDADE_PLU, plu)
	if erro != nil {
		return erro
	}

	if len(resultado) > 0 {
		return model.PluJaCadastrado
	}

	return nil
}

func (d *ProdutoDao) AdicionarProduto(contexto context.Context, produto *model.Produto, transacao *conexao_mysql.Transacao) error {
	erro := d.adicionarProduto(produto, transacao)
	if erro != nil {
		return erro
	}
	erro = d.adicionarCodigosBarras(contexto, produto, transacao)
	if erro != nil {
		return erro
	}
	return erro
}

func (d *ProdutoDao) adicionarProduto(produto *model.Produto, transacao *conexao_mysql.Transacao) error {
	if produto.Plu == nil {
		registros, erro := d.conexao.Select(QUERY_OBTER_PROXIMO_PLU_DISPONIVEL)
		if erro != nil {
			return erro
		}
		produto.Plu = registros[0]["plu"].(*int64)
	}

	return d.adicionarNovoProduto(produto, transacao)
}

func (d *ProdutoDao) adicionarNovoProduto(produto *model.Produto, transacao *conexao_mysql.Transacao) error {
	erro := transacao.Update(
		QUERY_INSERIR_NOVO_PRODUTO_COM_PLU_FIXO,
		produto.Plu,
		produto.Descricao,
		produto.DescricaoResumida,
		produto.CodigoDepartamento,
		produto.CodigoSecao,
		produto.Pesavel,
		produto.BloqueadoEntrada,
		produto.BloqueadoSaida,
		produto.Custo,
		produto.PrecoVenda,
		produto.Estoque,
		produto.EstoqueMinimo,
		produto.EstoqueMaximo,
		produto.DataCadastro,
		produto.DataAlteracao,
	)
	if erro != nil {
		return erro
	}
	return nil
}

func (d *ProdutoDao) adicionarCodigosBarras(contexto context.Context, produto *model.Produto, transacao *conexao_mysql.Transacao) error {
	erro := transacao.Update(QUERY_REMOVER_CODIGOS_BARRAS_PRODUTO, produto.Plu)
	if erro != nil {
		return erro
	}

	if len(produto.CodigosBarras) == 0 {
		return nil
	}

	erro = d.verificarExistenciaCodigoBarras(transacao, produto.CodigosBarras)
	if erro != nil {
		return erro
	}

	for _, codigoBarra := range produto.CodigosBarras {
		erro = transacao.Update(
			QUERY_INSERIR_CODIGO_BARRA,
			produto.Plu,
			codigoBarra.CodigoBarra,
		)
		if erro != nil {
			return erro
		}
	}
	return nil
}

func (d *ProdutoDao) verificarExistenciaCodigoBarras(transacao *conexao_mysql.Transacao, codigosBarras []*model.CodigoBarra) error {
	barrasAsString := make([]string, 0)
	for _, codigoBarra := range codigosBarras {
		barrasAsString = append(barrasAsString, "'"+codigoBarra.CodigoBarra+"'")
	}
	consulta := strings.Replace(VERIFICAR_EXISTENCIA_CODIGOS_BARRAS, ":CODIGOS_BARRAS", strings.Join(barrasAsString, ","), 1)
	registros, erro := transacao.Select(consulta)
	if erro != nil {
		return erro
	}

	if len(registros) == 0 {
		return nil
	}

	cadastrados := d.mapearCodigosBarras(registros)
	return model.MontarErroCodigosBarrasCadastrado(cadastrados)
}

func (d *ProdutoDao) ObterProdutoEspecifico(contexto context.Context, plu int64) (*model.Produto, error) {
	registros, erro := d.conexao.Select(QUERY_BUSCAR_PRODUTO_ESPECIFICO, plu)
	if erro != nil {
		return nil, erro
	}

	if len(registros) == 0 {
		return nil, model.ProdutoNaoEncontrado
	}

	produtoAsMap := registros[0]

	produto := d.mapearProduto(produtoAsMap)

	codigosBarras, erro := d.obterCodigosBarrasDoProduto(contexto, plu)
	if erro != nil {
		return nil, erro
	}

	produto.CodigosBarras = codigosBarras

	return produto, nil
}

func (d *ProdutoDao) obterCodigosBarrasDoProduto(ctx context.Context, plu int64) ([]*model.CodigoBarra, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_CODIGOS_BARRAS_POR_PLU, plu)
	if erro != nil {
		return nil, erro
	}

	codigosBarras := d.mapearCodigosBarras(resultados)

	return codigosBarras, nil
}

func (d *ProdutoDao) BuscarProdutosEspecificos(contexto context.Context, plus []int64) ([]*model.ProdutoResumido, error) {
	plusAsString := d.obterPlusSeparadosPorVirgula(plus)
	query := strings.Replace(QUERY_BUSCAR_PRODUTOS_ESPECIFICOS, ":PLUS_ESPECIFICOS", plusAsString, 1)
	resultados, erro := d.conexao.Select(query)
	if erro != nil {
		return nil, erro
	}

	produtos := make([]*model.ProdutoResumido, 0)
	for _, resultado := range resultados {
		produto := d.mapearProdutoResumido(resultado)
		produtos = append(produtos, produto)
	}

	return produtos, erro
}

func (d *ProdutoDao) obterPlusSeparadosPorVirgula(plus []int64) string {
	var plusString []string
	for _, plu := range plus {
		plusString = append(plusString, strconv.FormatInt(plu, 10))
	}
	return strings.Join(plusString, ",")
}

func (d *ProdutoDao) RemoverProduto(contexto context.Context, plu int64) error {
	erro := d.conexao.Update(QUERY_REMOVER_CODIGOS_BARRAS_PRODUTO, plu)
	if erro != nil {
		return erro
	}

	erro = d.conexao.Update(QUERY_REMOVER_PRODUTO_ESPECIFICO, plu)
	if erro != nil {
		return erro
	}
	return nil
}
