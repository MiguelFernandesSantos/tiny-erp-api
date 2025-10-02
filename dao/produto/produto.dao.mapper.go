package produto_dao

import "api-produtos/model"

func (d *ProdutoDao) mapearProdutoResumido(resultado map[string]interface{}) *model.ProdutoResumido {
	plu := resultado["plu"].(*int64)

	var descricao *string
	if resultado["descricao"] != nil {
		descricao = resultado["descricao"].(*string)
	}

	var departamento *string
	if resultado["nome_departamento"] != nil {
		departamento = resultado["nome_departamento"].(*string)
	}

	var secao *string
	if resultado["nome_secao"] != nil {
		secao = resultado["nome_secao"].(*string)
	}

	var bloqueadoEntrada bool = false
	if resultado["bloqueado_entrada"] != nil {
		bloqueadoEntrada = *resultado["bloqueado_entrada"].(*int64) == 1
	}

	var bloqueadoSaida bool = false
	if resultado["bloqueado_saida"] != nil {
		bloqueadoSaida = *resultado["bloqueado_saida"].(*int64) == 1
	}

	var pesavel bool = false
	if resultado["pesavel"] != nil {
		pesavel = *resultado["pesavel"].(*int64) == 1
	}

	return &model.ProdutoResumido{
		Plu:              plu,
		Descricao:        descricao,
		Departamento:     departamento,
		Secao:            secao,
		BloqueadoEntrada: &bloqueadoEntrada,
		BloqueadoSaida:   &bloqueadoSaida,
		Pesavel:          &pesavel,
	}
}

func (d *ProdutoDao) mapearProduto(resultado map[string]interface{}) *model.Produto {
	var plu *int64
	if resultado["plu"] != nil {
		plu = resultado["plu"].(*int64)
	}

	var descricao *string
	if resultado["descricao"] != nil {
		descricao = resultado["descricao"].(*string)
	}

	var descricaoResumida *string
	if resultado["descricao_resumida"] != nil {
		descricaoResumida = resultado["descricao_resumida"].(*string)
	}

	var codigoDepartamento *int64
	if resultado["codigo_departamento"] != nil {
		codigoDepartamento = resultado["codigo_departamento"].(*int64)
	}

	var nomeDepartamento *string
	if resultado["nome_departamento"] != nil {
		nomeDepartamento = resultado["nome_departamento"].(*string)
	}

	var codigoSecao *int64
	if resultado["codigo_secao"] != nil {
		codigoSecao = resultado["codigo_secao"].(*int64)
	}

	var nomeSecao *string
	if resultado["nome_secao"] != nil {
		nomeSecao = resultado["nome_secao"].(*string)
	}

	var pesavel bool = false
	if resultado["pesavel"] != nil {
		pesavel = *resultado["pesavel"].(*int64) == 1
	}

	var bloqueadoEntrada bool = false
	if resultado["bloqueado_entrada"] != nil {
		bloqueadoEntrada = *resultado["bloqueado_entrada"].(*int64) == 1
	}

	var bloqueadoSaida bool = false
	if resultado["bloqueado_saida"] != nil {
		bloqueadoSaida = *resultado["bloqueado_saida"].(*int64) == 1
	}

	var custo *float64
	if resultado["custo"] != nil {
		custo = resultado["custo"].(*float64)
	}

	var precoVenda *float64
	if resultado["preco_venda"] != nil {
		precoVenda = resultado["preco_venda"].(*float64)
	}

	var estoque *float64
	if resultado["estoque"] != nil {
		estoque = resultado["estoque"].(*float64)
	}

	var estoqueMinimo *float64
	if resultado["estoque_minimo"] != nil {
		estoqueMinimo = resultado["estoque_minimo"].(*float64)
	}

	var estoqueMaximo *float64
	if resultado["estoque_maximo"] != nil {
		estoqueMaximo = resultado["estoque_maximo"].(*float64)
	}

	var dataCadastro string = ""
	if resultado["data_cadastro"] != nil {
		dataCadastro = *resultado["data_cadastro"].(*string)
	}

	var dataAlteracao string = ""
	if resultado["data_alteracao"] != nil {
		dataAlteracao = *resultado["data_alteracao"].(*string)
	}

	return &model.Produto{
		Plu:                plu,
		Descricao:          descricao,
		DescricaoResumida:  descricaoResumida,
		CodigoDepartamento: codigoDepartamento,
		NomeDepartamento:   nomeDepartamento,
		CodigoSecao:        codigoSecao,
		NomeSecao:          nomeSecao,
		Pesavel:            &pesavel,
		BloqueadoEntrada:   &bloqueadoEntrada,
		BloqueadoSaida:     &bloqueadoSaida,
		Custo:              custo,
		PrecoVenda:         precoVenda,
		Estoque:            estoque,
		EstoqueMinimo:      estoqueMinimo,
		EstoqueMaximo:      estoqueMaximo,
		DataCadastro:       dataCadastro,
		DataAlteracao:      dataAlteracao,
	}
}

func (d *ProdutoDao) mapearCodigosBarras(resultados []map[string]any) []*model.CodigoBarra {
	codigosBarras := make([]*model.CodigoBarra, 0)
	for _, resultado := range resultados {
		codigoBarra := d.mapearCodigoBarras(resultado)
		codigosBarras = append(codigosBarras, codigoBarra)
	}
	return codigosBarras
}

func (d *ProdutoDao) mapearCodigoBarras(resultado map[string]interface{}) *model.CodigoBarra {
	var plu *int64
	if resultado["plu"] != nil {
		plu = resultado["plu"].(*int64)
	}

	codigoBarra := ""
	if resultado["codigo_barras"] != nil {
		codigoBarra = *resultado["codigo_barras"].(*string)
	}

	return &model.CodigoBarra{
		Plu:         plu,
		CodigoBarra: codigoBarra,
	}
}
