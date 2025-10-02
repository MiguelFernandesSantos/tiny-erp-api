package produto_dao

const (
	QUERY_BUSCAR_PRODUTOS_PAGINADO = `
		SELECT
		    produto.plu,
		    produto.descricao,
		    produto.pesavel,
		    departamento.nome AS nome_departamento,
		    secao.nome AS nome_secao,
		    produto.bloqueado_entrada,
		    produto.bloqueado_saida
		FROM produto
		INNER JOIN departamento ON produto.codigo_departamento = departamento.codigo
		INNER JOIN secao ON produto.codigo_secao = secao.codigo
		WHERE (produto.plu, produto.descricao) > (?, ?)
	`

	QUERY_VERIFICAR_DISPONIBILIDADE_PLU = `
		SELECT plu FROM produto
		WHERE plu = ?
		LIMIT 1
	`

	QUERY_INSERIR_NOVO_PRODUTO_COM_PLU_FIXO = `
		INSERT INTO produto (
			plu, 
			descricao, 
			descricao_resumida,
			codigo_departamento, 
			codigo_secao, 
			pesavel,
			bloqueado_entrada,
			bloqueado_saida,
			custo, 
			preco_venda, 
			estoque, 
			estoque_minimo, 
			estoque_maximo, 
			data_cadastro, 
			data_alteracao
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	VERIFICAR_EXISTENCIA_CODIGOS_BARRAS = `
		SELECT plu, codigo_barras FROM codigo_barras
		WHERE codigo_barras IN(:CODIGOS_BARRAS)
		LIMIT 1
	`

	QUERY_INSERIR_CODIGO_BARRA = `
		INSERT INTO codigo_barras (
			plu,
			codigo_barras  
		) VALUES (?, ?)
	`

	QUERY_OBTER_PROXIMO_PLU_DISPONIVEL = `
		SELECT COALESCE(MAX(plu), 0) + 1 AS plu FROM produto
	`

	QUERY_BUSCAR_PRODUTO_ESPECIFICO = `
		SELECT
			plu, 
			descricao, 
			descricao_resumida,
			produto.codigo_departamento, 
			produto.codigo_secao, 
			pesavel,
			bloqueado_entrada,
			bloqueado_saida,
			custo, 
			preco_venda, 
			estoque, 
			estoque_minimo, 
			estoque_maximo, 
			data_cadastro, 
			data_alteracao,
			departamento.nome AS nome_departamento,
		    secao.nome AS nome_secao
		FROM produto
		INNER JOIN departamento ON produto.codigo_departamento = departamento.codigo
		INNER JOIN secao ON produto.codigo_secao = secao.codigo
		WHERE produto.plu = ?
		LIMIT 1
	`

	QUERY_BUSCAR_CODIGOS_BARRAS_POR_PLU = `
		SELECT
		    plu,
			codigo_barras
		FROM codigo_barras
		WHERE plu = ?
	`

	QUERY_BUSCAR_PRODUTOS_ESPECIFICOS = `
		SELECT
		    produto.plu,
		    produto.descricao,
		    produto.pesavel,
		    departamento.nome AS nome_departamento,
		    secao.nome AS nome_secao,
		    produto.bloqueado_entrada,
		    produto.bloqueado_saida
		FROM produto
		INNER JOIN departamento ON produto.codigo_departamento = departamento.codigo
		INNER JOIN secao ON produto.codigo_secao = secao.codigo
		WHERE produto.plu IN(:PLUS_ESPECIFICOS)
	`

	QUERY_REMOVER_CODIGOS_BARRAS_PRODUTO = `
		DELETE FROM codigo_barras
		WHERE plu = ?	
	`

	QUERY_REMOVER_PRODUTO_ESPECIFICO = `
		DELETE FROM produto
		WHERE plu = ?
	`
)
