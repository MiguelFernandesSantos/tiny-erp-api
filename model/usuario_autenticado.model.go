package model

type LoginUsuario struct {
	Email *string `json:"email,omitempty"`
	Senha *string `json:"senha,omitempty"`
}

type UsuarioAutenticado struct {
	Codigo       *int64  `json:"codigo,omitempty"`
	Nome         *string `json:"nome,omitempty"`
	Email        *string `json:"email,omitempty"`
	Perfil       *int64  `json:"perfil,omitempty"`
	Token        *string `json:"token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}
