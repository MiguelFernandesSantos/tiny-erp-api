package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Produto struct {
	Plu                *int64           `json:"plu,omitempty"`
	Descricao          *string          `json:"descricao,omitempty"`
	DescricaoResumida  *string          `json:"descricao_resumida,omitempty"`
	CodigoDepartamento *int64           `json:"codigo_departamento,omitempty"`
	NomeDepartamento   *string          `json:"nome_departamento,omitempty"`
	CodigoSecao        *int64           `json:"codigo_secao,omitempty"`
	NomeSecao          *string          `json:"nome_secao,omitempty"`
	Pesavel            *bool            `json:"pesavel,omitempty"`
	BloqueadoEntrada   *bool            `json:"bloqueado_entrada,omitempty"`
	BloqueadoSaida     *bool            `json:"bloqueado_saida,omitempty"`
	Custo              *float64         `json:"custo,omitempty"`
	PrecoVenda         *float64         `json:"preco_venda,omitempty"`
	Estoque            *float64         `json:"estoque,omitempty"`
	EstoqueMinimo      *float64         `json:"estoque_minimo,omitempty"`
	EstoqueMaximo      *float64         `json:"estoque_maximo,omitempty"`
	DataCadastro       string           `json:"data_cadastro"`
	DataAlteracao      string           `json:"data_alteracao"`
	CodigosBarras      []*CodigoBarra   `json:"codigos_barras,omitempty"`
	Imagens            []*ImagemProduto `json:"imagens,omitempty"`
}

func (p *Produto) Validar(pluObrigatorio bool) error {
	if p.Plu == nil && pluObrigatorio {
		return fmt.Errorf("O produto enviado não tem PLU")
	}

	if p.Descricao == nil || *p.Descricao == "" {
		return fmt.Errorf("A descricao do produto nao foi informada")
	}

	if p.CodigoDepartamento == nil || *p.CodigoDepartamento <= 0 {
		return fmt.Errorf("O codigo do departamento nao foi informado")
	}

	if p.CodigoSecao == nil || *p.CodigoSecao <= 0 {
		return fmt.Errorf("O codigo da secao nao foi informado")
	}

	if p.Custo == nil || *p.Custo < 0 {
		return fmt.Errorf("O custo do produto nao foi informado ou é invalido")
	}

	if p.PrecoVenda == nil || *p.PrecoVenda < 0 {
		return fmt.Errorf("O preco de venda do produto nao foi informado ou é invalido")
	}

	if p.Estoque != nil && *p.Estoque < 0 {
		return fmt.Errorf("O estoque do produto é invalido")
	}

	if p.EstoqueMinimo != nil && *p.EstoqueMinimo < 0 {
		return fmt.Errorf("O estoque minimo do produto é invalido")
	}

	if p.EstoqueMaximo != nil && *p.EstoqueMaximo < 0 {
		return fmt.Errorf("O estoque maximo do produto é invalido")
	}

	for _, codigoBarra := range p.CodigosBarras {
		if erro := codigoBarra.Validar(); erro != nil {
			return erro
		}
	}

	for _, imagem := range p.Imagens {
		if erro := imagem.Validar(); erro != nil {
			return erro
		}
	}

	return nil
}

type CodigoBarra struct {
	Plu         *int64 `json:"plu,omitempty"`
	CodigoBarra string `json:"codigo_barra"`
}

func (c *CodigoBarra) Validar() error {
	if c.CodigoBarra == "" {
		return fmt.Errorf("O codigo de barra nao foi informado")
	}

	return nil
}

type ImagemProduto struct {
	Codigo    *int64 `json:"codigo,omitempty"`
	Plu       *int64 `json:"plu,omitempty"`
	UrlImagem string `json:"url_imagem"`
}

func (i *ImagemProduto) Validar() error {
	if i.UrlImagem == "" {
		return fmt.Errorf("A url da imagem nao foi informada")
	}

	return nil
}

type ProdutoResumido struct {
	Plu              *int64  `json:"plu,omitempty"`
	Descricao        *string `json:"descricao,omitempty"`
	Departamento     *string `json:"departamento,omitempty"`
	Secao            *string `json:"secao,omitempty"`
	Pesavel          *bool   `json:"pesavel,omitempty"`
	BloqueadoEntrada *bool   `json:"bloqueado_entrada,omitempty"`
	BloqueadoSaida   *bool   `json:"bloqueado_saida,omitempty"`
}

func (p *ProdutoResumido) Validar(pluObrigatorio bool) error {
	if p.Plu == nil && pluObrigatorio {
		return fmt.Errorf("O produto enviado não tem PLU")
	}

	if p.Descricao == nil || *p.Descricao == "" {
		return fmt.Errorf("A descricao do produto nao foi informada")
	}

	return nil
}

func (p *ProdutoResumido) PluComoString() *string {
	if p.Plu != nil {
		plu := strconv.FormatInt(*p.Plu, 10)
		return &plu
	}
	return nil
}

var (
	PluJaCadastrado         = errors.New("Já existe um PLU com essas informações")
	ProdutoNaoEncontrado    = errors.New("Não foi encontrado um produto com essas informações")
	CodigoBarraJaCadastrado = errors.New("Os seguintes códigos de barras já estão cadastrados: ")
)

func MontarErroCodigosBarrasCadastrado(codigosBarras []*CodigoBarra) error {
	if len(codigosBarras) == 0 {
		return nil
	}

	var mensagensErro []string
	for _, codigoBarra := range codigosBarras {
		mensagemErro := "O código de barras " + codigoBarra.CodigoBarra + " já está cadastrado para o produto de PLU " + strconv.FormatInt(*codigoBarra.Plu, 10)
		mensagensErro = append(mensagensErro, mensagemErro)
	}

	return fmt.Errorf("%w %s", CodigoBarraJaCadastrado, strings.Join(mensagensErro, "\n- "))
}
