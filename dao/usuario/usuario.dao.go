package usuario_dao

import (
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/model"
	"context"
)

type UsuarioDao struct {
	conexao *conexao_mysql.MySQL
}

func InstanciarUsuarioDao(conexao *conexao_mysql.MySQL) *UsuarioDao {
	return &UsuarioDao{
		conexao: conexao,
	}
}

func (d *UsuarioDao) ObterUsuariosPaginado(contexto context.Context, pesquisa string, ultimoUsuario *int64) ([]*model.UsuarioResumido, error) {
	pesquisaLike := "%" + pesquisa + "%"

	if ultimoUsuario == nil {
		ultimoUsuario = new(int64)
		*ultimoUsuario = 0
	}

	resultados, erro := d.conexao.Select(QUERY_BUSCAR_USUARIOS_PAGINADO, ultimoUsuario, pesquisaLike, pesquisaLike)
	if erro != nil {
		return nil, erro
	}

	usuarios := make([]*model.UsuarioResumido, 0)
	for _, resultado := range resultados {
		usuario := d.mapearUsuarioResumido(resultado)
		usuarios = append(usuarios, usuario)
	}

	return usuarios, erro
}

func (d *UsuarioDao) VerificarDisponibilidadeEmail(contexto context.Context, email *string) error {
	resultado, erro := d.conexao.Select(QUERY_VERIFICAR_DISPONIBILIDADE_EMAIL_USUARIO, email)
	if erro != nil {
		return erro
	}

	if len(resultado) > 0 {
		return model.UsuarioJaExiste
	}

	return nil
}

func (d *UsuarioDao) AdicionarUsuario(contexto context.Context, usuario *model.Usuario, senhaCriptografada *string) error {
	dataAtual := model.DataHoraAtual().FormatarISO8601()
	erro := d.conexao.Update(QUERY_INSERIR_NOVO_USUARIO, usuario.Nome, usuario.Email, senhaCriptografada, usuario.Perfil, dataAtual)
	if erro != nil {
		return erro
	}
	return nil
}

func (d *UsuarioDao) ObterUsuarioEspecifico(contexto context.Context, usuario *int64) (*model.Usuario, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_USUARIO_ESPECIFICO, usuario)
	if erro != nil {
		return nil, erro
	}

	if len(resultados) == 0 {
		return nil, model.UsuarioNaoEncontrado
	}

	usuarioAsMap := resultados[0]

	usuarioResumido := d.mapearUsuario(usuarioAsMap)
	return usuarioResumido, nil
}

func (d *UsuarioDao) DesativarUsuarioEspecifico(contexto context.Context, usuario *int64) error {
	erro := d.conexao.Update(QUERY_DESATIVAR_USUARIO_ESPECIFICO, usuario)
	if erro != nil {
		return erro
	}
	return nil
}

func (d *UsuarioDao) ObterSenhaUsuario(contexto context.Context, email *string) (*string, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_SENHA_USUARIO_POR_EMAIL, email)
	if erro != nil {
		return nil, erro
	}

	if len(resultados) == 0 {
		return nil, model.UsuarioNaoEncontrado
	}

	usuarioAsMap := resultados[0]

	ativo := *usuarioAsMap["ativo"].(*int64) == 1
	if !ativo {
		return nil, model.UsuarioInativo
	}

	return usuarioAsMap["senha"].(*string), nil
}

func (d *UsuarioDao) ObterDadosAutenticacaoUsuario(contexto context.Context, email *string) (*model.UsuarioAutenticado, error) {
	resultados, erro := d.conexao.Select(QUERY_BUSCAR_DADOS_AUTENTICACAO_USUARIO_POR_EMAIL, email)
	if erro != nil {
		return nil, erro
	}

	if len(resultados) == 0 {
		return nil, model.UsuarioNaoEncontrado
	}

	usuarioAsMap := resultados[0]

	return d.mapearUsuarioAutenticado(usuarioAsMap), nil
}
