package departamento_dao

import (
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"context"
)

type DepartamentoDao struct {
	conexao *conexao_mysql.MySQL
}

func InstanciarDepartamentoDao(conexao *conexao_mysql.MySQL) *DepartamentoDao {
	return &DepartamentoDao{
		conexao: conexao,
	}
}

func (d *DepartamentoDao) ObterDepartamentosPaginado(context context.Context, pesquisa string, ultimoDepartamento *int64) ([]*model.Departamento, error) {
	pesquisaLike := "%" + pesquisa + "%"

	if ultimoDepartamento == nil {
		ultimoDepartamento = new(int64)
		*ultimoDepartamento = 0
	}

	resultados, erro := d.conexao.Select(QUERY_BUSCAR_DEPARTAMENTOS_PAGINADO, pesquisaLike, ultimoDepartamento)
	if erro != nil {
		return nil, erro
	}

	departamentos := make([]*model.Departamento, 0)
	for _, resultado := range resultados {
		departamento := d.mapearDepartamento(resultado)
		departamentos = append(departamentos, departamento)
	}

	return departamentos, erro
}

func (d *DepartamentoDao) VerificarDisponibilidadeNome(ctx context.Context, nome *string) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_DISPONIBILIDADE_NOME_DEPARTAMENTO, nome)
	if erro != nil {
		return erro
	}

	if len(resultado) > 0 {
		return model.DepartamentoJaExiste
	}

	return nil
}

func (d *DepartamentoDao) AdicionarDepartamento(context context.Context, departamento *model.Departamento) error {
	erro := d.conexao.Update(QUERY_INSERIR_NOVO_DEPARTAMENTO, departamento.Nome)
	if erro != nil {
		return erro
	}

	return nil
}

func (d *DepartamentoDao) ObterDepartamentoEspecifico(ctx context.Context, numeroDepartamento int64) (*model.Departamento, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_DEPARTAMENTO_ESPECIFICO, numeroDepartamento)
	if erro != nil {
		return nil, erro
	}

	if len(resultados) == 0 {
		return nil, model.DepartamentoNaoEncontrado
	}

	departamentoAsMap := resultados[0]

	departamento := d.mapearDepartamento(departamentoAsMap)

	secoes, erro := d.buscarSecoesDoDepartamento(ctx, numeroDepartamento)
	if erro != nil {
		return nil, erro
	}

	departamento.Secoes = secoes

	return departamento, nil
}

func (d *DepartamentoDao) buscarSecoesDoDepartamento(context context.Context, numeroDepartamento int64) ([]*model.Secao, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_SECOES_DO_DEPARTAMENTO, numeroDepartamento)
	if erro != nil {
		return nil, erro
	}

	secoes := make([]*model.Secao, 0)
	for _, resultado := range resultados {
		secao := d.mapearSecao(resultado)
		secoes = append(secoes, secao)
	}

	return secoes, nil
}

func (d *DepartamentoDao) VerificarSePossuiProdutosAssociados(contexto context.Context, departamento *int64) error {

	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_SE_DEPARTAMENTO_POSSUI_PRODUTOS_ASSOCIADOS, departamento)
	if erro != nil {
		return erro
	}

	if len(resultado) == 0 {
		return model.DepartamentoNaoEncontrado
	}

	resultadoMap := resultado[0]

	quantidadeProdutos := resultadoMap["quantidade_produtos"].(*int64)
	if quantidadeProdutos == nil || *quantidadeProdutos > 0 {
		return model.PossuiProdutosAssociados

	}

	quantidadeSecoes := resultadoMap["quantidade_secoes"].(*int64)
	if quantidadeSecoes == nil || *quantidadeSecoes > 0 {
		return model.PossuiSecoesAssociadas
	}

	return nil
}

func (d *DepartamentoDao) RemoverDepartamentoEspecifico(contexto context.Context, departamento *int64) error {
	erro := d.conexao.Update(QUERY_REMOVER_DEPARTAMENTO_ESPECIFICO, departamento)
	if erro != nil {
		return erro
	}
	return nil
}

func (d *DepartamentoDao) VerificarExistenciaDepartamento(contexto context.Context, departamento *int64) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_EXISTENCIA_DEPARTAMENTO, departamento)
	if erro != nil {
		return erro
	}

	if len(resultado) == 0 {
		return model.DepartamentoNaoEncontrado
	}

	return nil
}

func (d *DepartamentoDao) VerificarExistenciaSecao(ctx context.Context, secao *int64, departamento *int64) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_EXISTENCIA_SECAO, secao, departamento)
	if erro != nil {
		return erro
	}

	if len(resultado) == 0 {
		return model.SecaoNaoEncontrada
	}

	return nil
}

func (d *DepartamentoDao) VerificarDisponibilidadeNomeSecao(contexto context.Context, secao *model.Secao) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_DISPONIBILIDADE_NOME_SECAO, secao.Nome, secao.NumeroDepartamento)
	if erro != nil {
		return erro
	}

	if len(resultado) > 0 {
		return model.SecaoJaExiste
	}

	return nil
}

func (d *DepartamentoDao) AdicionarSecaoAoDepartamento(contexto context.Context, secao *model.Secao) error {
	erro := d.conexao.Update(QUERY_INSERIR_NOVA_SECAO, secao.Nome, secao.NumeroDepartamento)
	if erro != nil {
		return erro
	}
	return nil
}

func (d *DepartamentoDao) ObterSecaoEspecifica(contexto context.Context, secao int64, departamento int64) (*model.Secao, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_SECAO_ESPECIFICA, secao, departamento)
	if erro != nil {
		return nil, erro
	}

	if len(resultados) == 0 {
		return nil, model.SecaoNaoEncontrada
	}

	secaoAsMap := resultados[0]
	return d.mapearSecao(secaoAsMap), nil
}

func (d *DepartamentoDao) VerificarSePossuiProdutosAssociadosASecao(contexto context.Context, secao int64, departamento int64) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_SE_SECAO_POSSUI_PRODUTOS_ASSOCIADOS, departamento, secao)
	if erro != nil {
		return erro
	}

	if len(resultado) == 0 {
		return model.SecaoNaoEncontrada
	}

	resultadoMap := resultado[0]

	quantidadeProdutos := resultadoMap["quantidade_produtos"].(*int64)
	if quantidadeProdutos == nil || *quantidadeProdutos > 0 {
		return model.PossuiProdutosAssociados
	}

	return nil
}

func (d *DepartamentoDao) RemoverSecaoEspecifica(contexto context.Context, secao int64, departamento int64) error {
	erro := d.conexao.Update(QUERY_REMOVER_SECAO_ESPECIFICA, departamento, secao)
	if erro != nil {
		return erro
	}
	return nil
}
