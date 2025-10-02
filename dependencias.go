package api_produtos

import (
	"api-produtos/controller"
	departamento_dao "api-produtos/dao/departamento"
	"api-produtos/dao/produto"
	log "api-produtos/internal"
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/internal/state"
	"api-produtos/service"
	"api-produtos/service/api"
	"context"
	"embed"
	"os"

	"github.com/labstack/echo/v4"
)

// DAO
var (
	produtoDao      *produto_dao.ProdutoDao
	departamentodao *departamento_dao.DepartamentoDao
)

// SERVICES
var (
	produtoService             *service.ProdutoService
	produtoApiEmbeddingService *api.ProdutoApiEmbeddingService

	departamentoService *service.DepartamentoService
)

// EMBED
var (
	//go:embed embed/migration/mysql
	migrationFS embed.FS
)

func RegistrarDependencias(contextoBackgroud context.Context, apiRoot *echo.Echo, conexao *conexao_mysql.MySQL) {

	erro := conexao.ExecutarMigracao(migrationFS, "schema_migration")
	if erro != nil {
		log.Exception(contextoBackgroud, "Ocorreu um erro ao tentar realizar o migration", &erro)
		os.Exit(1)
	}

	instanciarDAOs(conexao)
	instanciarServices()
	inicializarServices(conexao)

	grupoV1 := apiRoot.Group("/api/v1")

	controller.InstanciarProdutoController(grupoV1, produtoService)
	controller.InstanciarDepartamentoController(grupoV1, departamentoService)
}

func instanciarDAOs(conexao *conexao_mysql.MySQL) {
	produtoDao = produto_dao.InstanciarProdutoDao(conexao)
	departamentodao = departamento_dao.InstanciarDepartamentoDao(conexao)
}

func instanciarServices() {
	produtoApiEmbeddingService = api.InstanciarProdutoApiEmbeddingService()
	produtoService = service.InstanciarProdutoService()
	departamentoService = service.InstanciarDepartamentoService()
}

func inicializarServices(conexao *conexao_mysql.MySQL) {
	produtoApiEmbeddingService.Inicializar(state.URL_API_EMBEDDING, state.UTILIZAR_API_EMBEDDING)
	produtoService.Inicializar(produtoDao, produtoApiEmbeddingService, departamentoService, conexao)
	departamentoService.Inicializar(departamentodao, conexao)
}
