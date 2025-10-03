package usuario_dao

import "api-produtos/model"

func (d *UsuarioDao) mapearUsuarioResumido(resultado map[string]any) *model.UsuarioResumido {
	var codigoUsuario *int64
	if resultado["codigo_usuario"] != nil {
		codigoUsuario = resultado["codigo_usuario"].(*int64)
	}

	var nome *string
	if resultado["nome"] != nil {
		nome = resultado["nome"].(*string)
	}

	var email *string
	if resultado["email"] != nil {
		email = resultado["email"].(*string)
	}

	var perfil *model.Perfil
	if resultado["perfil"] != nil {
		perfilConvertido := model.Perfil(*resultado["perfil"].(*int64))
		perfil = &perfilConvertido
	}

	var ativo bool = false
	if resultado["ativo"] != nil {
		ativo = *resultado["ativo"].(*int64) == 1
	}

	return &model.UsuarioResumido{
		Codigo: codigoUsuario,
		Nome:   nome,
		Email:  email,
		Perfil: perfil,
		Ativo:  &ativo,
	}
}

func (d *UsuarioDao) mapearUsuario(resultado map[string]any) *model.Usuario {
	var codigoUsuario *int64
	if resultado["codigo"] != nil {
		codigoUsuario = resultado["codigo"].(*int64)
	}

	var nome *string
	if resultado["nome"] != nil {
		nome = resultado["nome"].(*string)
	}

	var email *string
	if resultado["email"] != nil {
		email = resultado["email"].(*string)
	}

	var perfil *model.Perfil
	if resultado["perfil"] != nil {
		perfilConvertido := model.Perfil(*resultado["perfil"].(*int64))
		perfil = &perfilConvertido
	}

	var ativo bool = false
	if resultado["ativo"] != nil {
		ativo = *resultado["ativo"].(*int64) == 1
	}

	return &model.Usuario{
		Codigo: codigoUsuario,
		Nome:   nome,
		Email:  email,
		Perfil: perfil,
		Ativo:  &ativo,
	}
}
