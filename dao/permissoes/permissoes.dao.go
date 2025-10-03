package permissoes_dao

import conexao_mysql "api-produtos/internal/database"

type PermissoesDao struct {
	conexao *conexao_mysql.MySQL
}

func InstanciarPermissoesDao(conexao *conexao_mysql.MySQL) *PermissoesDao {
	return &PermissoesDao{
		conexao: conexao,
	}
}
