package controller

import (
	"api-produtos/model"
	"api-produtos/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DepartamentoController struct {
	departamentoService *service.DepartamentoService
}

func InstanciarDepartamentoController(grupoV1 *echo.Group, departamentoService *service.DepartamentoService) *DepartamentoController {
	controller := &DepartamentoController{
		departamentoService: departamentoService,
	}

	grupoDepartamentos := grupoV1.Group("/departamento")

	grupoDepartamentos.GET("", controller.obterDepartamentosPaginado)
	grupoDepartamentos.POST("", controller.adicionarDepartamento)
	grupoDepartamentos.GET("/:numeroDepartamento", controller.obterDepartamentoEspecifico)
	grupoDepartamentos.DELETE("/:numeroDepartamento", controller.removerDepartamentoEspecifico)
	grupoDepartamentos.POST("/:numeroDepartamento/secao", controller.adicionarSecaoAoDepartamento)
	grupoDepartamentos.GET("/:numeroDepartamento/secao/:numeroSecao", controller.obterSecaoEspecifica)
	grupoDepartamentos.DELETE("/:numeroDepartamento/secao/:numeroSecao", controller.removerSecaoEspecifica)

	return controller
}

func (c *DepartamentoController) obterDepartamentosPaginado(contexto echo.Context) error {
	pesquisa := ""
	var ultimoDepartamento *int64 // TODO precisa receber esses dois valores

	departamentos, erro := c.departamentoService.ObterDepartamentosPaginado(contexto.Request().Context(), pesquisa, ultimoDepartamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, departamentos)
}

func (c *DepartamentoController) adicionarDepartamento(contexto echo.Context) error {
	var (
		departamento *model.Departamento
	)

	if erro := contexto.Bind(&departamento); erro != nil {
		return contexto.JSON(http.StatusBadRequest, nil)
	}

	erro := departamento.Validar()
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, erro.Error())
	}

	erro = c.departamentoService.AdicionarDepartamento(contexto.Request().Context(), departamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusCreated, departamento)
}

func (c *DepartamentoController) obterDepartamentoEspecifico(contexto echo.Context) error {
	var (
		numeroDepartamentoAsString = contexto.Param("numeroDepartamento")
	)

	if numeroDepartamentoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento é obrigatório")
	}

	numeroDepartamento, erro := strconv.ParseInt(numeroDepartamentoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento inválido")
	}

	departamento, erro := c.departamentoService.ObterDepartamentoEspecifico(contexto.Request().Context(), numeroDepartamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, departamento)
}

func (c *DepartamentoController) removerDepartamentoEspecifico(contexto echo.Context) error {
	var (
		numeroDepartamentoAsString = contexto.Param("numeroDepartamento")
	)

	if numeroDepartamentoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento é obrigatório")
	}

	numeroDepartamento, erro := strconv.ParseInt(numeroDepartamentoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento inválido")
	}

	erro = c.departamentoService.RemoverDepartamentoEspecifico(contexto.Request().Context(), &numeroDepartamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.NoContent(http.StatusNoContent)
}

func (c *DepartamentoController) adicionarSecaoAoDepartamento(contexto echo.Context) error {
	var (
		secao *model.Secao
	)

	if erro := contexto.Bind(&secao); erro != nil {
		return contexto.JSON(http.StatusBadRequest, nil)
	}

	erro := secao.Validar()
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, erro.Error())
	}

	erro = c.departamentoService.AdicionarSecaoAoDepartamento(contexto.Request().Context(), secao)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusCreated, secao)
}

func (c *DepartamentoController) obterSecaoEspecifica(contexto echo.Context) error {
	var (
		numeroSecaoAsString        = contexto.Param("numeroSecao")
		numeroDepartamentoAsString = contexto.Param("numeroDepartamento")
	)

	if numeroSecaoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número da seção é obrigatório")
	}

	numeroSecao, erro := strconv.ParseInt(numeroSecaoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número da seção inválido")
	}

	if numeroDepartamentoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento é obrigatório")
	}

	numeroDepartamento, erro := strconv.ParseInt(numeroDepartamentoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento inválido")
	}

	secao, erro := c.departamentoService.ObterSecaoEspecifica(contexto.Request().Context(), numeroSecao, numeroDepartamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, secao)
}

func (c *DepartamentoController) removerSecaoEspecifica(contexto echo.Context) error {
	var (
		numeroSecaoAsString        = contexto.Param("numeroSecao")
		numeroDepartamentoAsString = contexto.Param("numeroDepartamento")
	)

	if numeroSecaoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número da seção é obrigatório")
	}

	numeroSecao, erro := strconv.ParseInt(numeroSecaoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número da seção inválido")
	}

	if numeroDepartamentoAsString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento é obrigatório")
	}

	numeroDepartamento, erro := strconv.ParseInt(numeroDepartamentoAsString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do departamento inválido")
	}

	erro = c.departamentoService.RemoverSecaoEspecifica(contexto.Request().Context(), numeroSecao, numeroDepartamento)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.NoContent(http.StatusNoContent)
}

func (c *DepartamentoController) tratarErros(contexto echo.Context, erro error) error {
	switch {
	case errors.Is(erro, model.DepartamentoNaoEncontrado):
		return contexto.JSON(http.StatusNotFound, erro.Error())
	case errors.Is(erro, model.DepartamentoJaExiste):
		return contexto.JSON(http.StatusConflict, erro.Error())
	case errors.Is(erro, model.SecaoNaoEncontrada):
		return contexto.JSON(http.StatusNotFound, erro.Error())
	case errors.Is(erro, model.SecaoJaExiste):
		return contexto.JSON(http.StatusConflict, erro.Error())
	case errors.Is(erro, model.PossuiProdutosAssociados):
		return contexto.JSON(http.StatusConflict, erro.Error())
	case errors.Is(erro, model.PossuiSecoesAssociadas):
		return contexto.JSON(http.StatusConflict, erro.Error())
	default:
		return contexto.JSON(http.StatusInternalServerError, erro.Error())
	}
}
