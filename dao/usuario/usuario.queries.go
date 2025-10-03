package usuario_dao

const (
	QUERY_BUSCAR_USUARIOS_PAGINADO = `
		SELECT codigo, nome, email, perfil, ativo
		FROM usuario
		WHERE codigo > ? AND (nome LIKE ? OR email LIKE ?)
		ORDER BY codigo
		LIMIT 24
	`

	QUERY_VERIFICAR_DISPONIBILIDADE_EMAIL_USUARIO = `
		SELECT codigo FROM usuario WHERE email = ?
	`

	QUERY_INSERIR_NOVO_USUARIO = `
		INSERT INTO usuario (nome, email, senha, perfil, data_criacao) VALUES (?, ?, ?, ?, ?)
	`

	QUERY_BUSCAR_USUARIO_ESPECIFICO = `
		SELECT codigo, nome, email, perfil, ativo
		FROM usuario
		WHERE codigo = ?
	`

	QUERY_DESATIVAR_USUARIO_ESPECIFICO = `
		UPDATE usuario SET ativo = 0 WHERE codigo = ?
	`

	QUERY_BUSCAR_SENHA_USUARIO_POR_EMAIL = `
		SELECT ativo, senha FROM usuario WHERE email = ?	
	`

	QUERY_BUSCAR_DADOS_AUTENTICACAO_USUARIO_POR_EMAIL = `
		SELECT codigo, nome, email, perfil FROM usuario WHERE email = ?
	`
)
