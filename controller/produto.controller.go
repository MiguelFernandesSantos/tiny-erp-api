package controller

import (
	"api-produtos/model"
	"api-produtos/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProdutoController struct {
	produtoService *service.ProdutoService
}

func InstanciarProdutoController(grupoV1 *echo.Group, produtoService *service.ProdutoService) *ProdutoController {
	controller := &ProdutoController{
		produtoService: produtoService,
	}

	grupoProdutos := grupoV1.Group("/produto")

	grupoProdutos.GET("", controller.obterProdutosPaginado)
	grupoProdutos.POST("", controller.adicionarProduto)
	grupoProdutos.DELETE("/:plu", controller.removerProduto)
	grupoProdutos.GET("/:plu", controller.obterProdutoEspecifico)
	grupoProdutos.GET("/:plu/semelhantes", controller.obterSemelhantes)

	return controller
}

func (controller *ProdutoController) obterProdutosPaginado(contexto echo.Context) error {
	pesquisa := ""
	var ultimoProduto *int // TODO precisa receber esses dois valores

	produtos, erro := controller.produtoService.ObterProdutosPaginado(contexto.Request().Context(), pesquisa, ultimoProduto)
	if erro != nil {
		return controller.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, produtos)
}

func (controller *ProdutoController) adicionarProduto(contexto echo.Context) error {
	var (
		produto *model.Produto
	)

	if erro := contexto.Bind(&produto); erro != nil {
		return contexto.JSON(http.StatusBadRequest, nil)
	}

	erro := produto.Validar(false)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, erro.Error())
	}

	erro = controller.produtoService.AdicionarProduto(contexto.Request().Context(), produto)
	if erro != nil {
		return controller.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusCreated, nil)
}

func (controller *ProdutoController) obterProdutoEspecifico(contexto echo.Context) error {
	var (
		pluAsString = contexto.Param("plu")
	)

	if pluAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Nenhum PLU foi informado!")
	}

	plu, erro := strconv.ParseInt(pluAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "O plu informado é invalido!")
	}

	produto, erro := controller.produtoService.ObterProdutoEspecifico(contexto.Request().Context(), plu)
	if erro != nil {
		return controller.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, produto)
}

func (controller *ProdutoController) removerProduto(contexto echo.Context) error {
	var (
		pluAsString = contexto.Param("plu")
	)

	if pluAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Nenhum PLU foi informado!")
	}

	plu, erro := strconv.ParseInt(pluAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "O plu informado é invalido!")
	}

	erro = controller.produtoService.RemoverProduto(contexto.Request().Context(), plu)
	if erro != nil {
		return controller.tratarErros(contexto, erro)
	}

	return contexto.NoContent(http.StatusNoContent)
}

func (controller *ProdutoController) obterSemelhantes(contexto echo.Context) error {
	var (
		pluAsString = contexto.Param("plu")
	)

	if pluAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Nenhum PLU foi informado!")
	}

	plu, erro := strconv.ParseInt(pluAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "O plu informado é invalido!")
	}

	produtosSemelhantes, erro := controller.produtoService.ObterSemelhantes(contexto.Request().Context(), plu)
	if erro != nil {
		return controller.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, produtosSemelhantes)
}

func (controller *ProdutoController) tratarErros(contexto echo.Context, erro error) error {
	switch {
	case errors.Is(erro, model.PluJaCadastrado):
		return contexto.JSON(http.StatusConflict, erro.Error())
	case errors.Is(erro, model.ProdutoNaoEncontrado):
		return contexto.JSON(http.StatusNotFound, erro.Error())
	case errors.Is(erro, model.DepartamentoNaoEncontrado):
		return contexto.JSON(http.StatusBadRequest, erro.Error())
	case errors.Is(erro, model.SecaoNaoEncontrada):
		return contexto.JSON(http.StatusBadRequest, erro.Error())
	case errors.Is(erro, model.CodigoBarraJaCadastrado):
		return contexto.JSON(http.StatusConflict, erro.Error())
	default:
		return contexto.JSON(http.StatusInternalServerError, "Ocorreu um erro interno no servidor!")
	}
}
