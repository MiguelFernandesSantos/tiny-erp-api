package model

import (
	"errors"
	"fmt"
)

type Perfil int64

var perfis = map[Perfil]string{
	1: "Comum",
	2: "Fiscal",
	3: "Administrador",
}

type Usuario struct {
	Codigo *int64  `json:"codigo"`
	Nome   *string `json:"nome"`
	Email  *string `json:"email"`
	Ativo  *bool   `json:"ativo"`
	Perfil *Perfil `json:"perfil"`
	Senha  *string `json:"senha,omitempty"`
}

func (u *Usuario) Validar() error {
	if u.Nome == nil || *u.Nome == "" {
		return fmt.Errorf("nome é obrigatório")
	}

	if u.Email == nil || *u.Email == "" {
		return fmt.Errorf("email é obrigatório")
	}

	if u.Senha == nil || *u.Senha == "" {
		return fmt.Errorf("senha é obrigatória")
	}

	if u.Perfil == nil {
		return fmt.Errorf("perfil é obrigatório")
	}

	if _, existe := perfis[*u.Perfil]; !existe {
		return fmt.Errorf("perfil inválido")
	}

	return nil
}

type UsuarioResumido struct {
	Codigo *int64  `json:"codigo"`
	Nome   *string `json:"nome"`
	Email  *string `json:"email"`
	Ativo  *bool   `json:"ativo"`
	Perfil *Perfil `json:"perfil"`
}

var (
	UsuarioNaoEncontrado = errors.New("usuário não encontrado")
	UsuarioInativo       = errors.New("usuário inativo")
	UsuarioJaExiste      = errors.New("usuário com o mesmo email já existe")
	SenhaIncorreta       = errors.New("senha incorreta")
	UsuarioDesativado    = errors.New("usuário desativado")
)
