package departamento_dao

const (
	QUERY_BUSCAR_DEPARTAMENTOS_PAGINADO = ` 
		SELECT
		    codigo,
		    nome
		FROM departamento
		WHERE nome LIKE ?
		AND codigo > ?
		ORDER BY codigo
		LIMIT 24
	`

	QUERY_VERIFICAR_DISPONIBILIDADE_NOME_DEPARTAMENTO = `
		SELECT nome FROM departamento
		WHERE nome = ?
		LIMIT 1
	`

	QUERY_INSERIR_NOVO_DEPARTAMENTO = `
		INSERT INTO departamento (nome)
		VALUES (?)
	`

	QUERY_BUSCAR_DEPARTAMENTO_ESPECIFICO = `
		SELECT
		    codigo,
		    nome
		FROM departamento
		WHERE codigo = ?
	`

	QUERY_BUSCAR_SECOES_DO_DEPARTAMENTO = `
		SELECT
		    codigo,
		    codigo_departamento,
		    nome
		FROM secao
		WHERE codigo_departamento = ?
	`

	QUERY_VERIFICAR_SE_DEPARTAMENTO_POSSUI_PRODUTOS_ASSOCIADOS = `
		SELECT 
			(SELECT COUNT(*) FROM produto WHERE codigo_departamento = departamento.codigo) AS quantidade_produtos,
			(SELECT COUNT(*) FROM secao WHERE codigo_departamento = departamento.codigo) AS quantidade_secoes
		FROM departamento
		WHERE codigo = ?
		LIMIT 1
	`

	QUERY_REMOVER_DEPARTAMENTO_ESPECIFICO = `
		DELETE FROM departamento
		WHERE codigo = ?
	`

	QUERY_VERIFICAR_EXISTENCIA_DEPARTAMENTO = `
		SELECT codigo FROM departamento
		WHERE codigo = ?
		LIMIT 1
	`

	QUERY_VERIFICAR_EXISTENCIA_SECAO = `
		SELECT codigo FROM secao
		WHERE codigo = ? AND codigo_departamento = ?
		LIMIT 1
	`

	QUERY_VERIFICAR_DISPONIBILIDADE_NOME_SECAO = `
		SELECT nome FROM secao
		WHERE nome = ? AND codigo_departamento = ?
		LIMIT 1
	`

	QUERY_INSERIR_NOVA_SECAO = `
		INSERT INTO secao (nome, codigo_departamento)
		VALUES (?, ?)
	`

	QUERY_BUSCAR_SECAO_ESPECIFICA = `
		SELECT
		    codigo,
		    codigo_departamento,
		    nome
		FROM secao
		WHERE codigo = ? AND codigo_departamento = ?
	`

	QUERY_VERIFICAR_SE_SECAO_POSSUI_PRODUTOS_ASSOCIADOS = `
		SELECT (SELECT COUNT(*) FROM produto WHERE codigo_departamento = secao.codigo_departamento AND codigo_secao = secao.codigo) AS quantidade_produtos
		FROM secao
		WHERE codigo = ? AND codigo_departamento = ?
		LIMIT 1
	`

	QUERY_REMOVER_SECAO_ESPECIFICA = `
		DELETE FROM secao
		WHERE codigo = ? AND codigo_departamento = ?
	`
)
