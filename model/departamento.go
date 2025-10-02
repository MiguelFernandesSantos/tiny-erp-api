package model

import "errors"

var (
	DepartamentoNaoEncontrado = errors.New("departamento nao encontrado")
	DepartamentoJaExiste      = errors.New("departamento ja existe")
	SecaoNaoEncontrada        = errors.New("secao nao encontrada")
	SecaoJaExiste             = errors.New("secao ja existe")
	PossuiProdutosAssociados  = errors.New("departamento ou secao possui produtos associados")
	PossuiSecoesAssociadas    = errors.New("departamento possui secoes associadas")
)

type Departamento struct {
	NumeroDepartamento *int64   `json:"numeroDepartamento,omitempty"`
	Nome               *string  `json:"nome,omitempty"`
	Secoes             []*Secao `json:"secoes,omitempty"`
}

func (d *Departamento) Validar() error {
	if d.Nome == nil || *d.Nome == "" {
		return errors.New("o nome do departamento e obrigatorio")
	}
	return nil
}

type Secao struct {
	NumeroSecao        *int64  `json:"numeroSecao"`
	NumeroDepartamento *int64  `json:"numeroDepartamento"`
	Nome               *string `json:"nome"`
}

func (s *Secao) Validar() error {
	if s.Nome == nil || *s.Nome == "" {
		return errors.New("o nome da secao e obrigatorio")
	}
	return nil
}
