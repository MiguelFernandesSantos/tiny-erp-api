package api_produtos

import (
	"api-produtos/controller"
	departamento_dao "api-produtos/dao/departamento"
	"api-produtos/dao/produto"
	usuario_dao "api-produtos/dao/usuario"
	internal "api-produtos/internal"
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
	departamentoDao *departamento_dao.DepartamentoDao
	usuarioDao      *usuario_dao.UsuarioDao
)

// SERVICES
var (
	produtoService             *service.ProdutoService
	produtoApiEmbeddingService *api.ProdutoApiEmbeddingService

	departamentoService *service.DepartamentoService

	usuarioService *service.UsuarioService

	autenticacaoService *service.AutenticacaoService

	tokenService *service.TokenService

	criptografiaService *service.CriptografiaService
)

// EMBED
var (
	//go:embed embed/migration/mysql
	migrationFS embed.FS

	//go:embed embed/keys/private.pem
	bytesChavePrivada []byte

	//go:embed embed/keys/public.pem
	bytesChavePublica []byte
)

func RegistrarDependencias(contextoBackgroud context.Context, apiRoot *echo.Echo, conexao *conexao_mysql.MySQL) {

	erro := conexao.ExecutarMigracao(migrationFS, "schema_migration")
	if erro != nil {
		internal.Exception(contextoBackgroud, "Ocorreu um erro ao tentar realizar o migration", &erro)
		os.Exit(1)
	}

	configurarAPI(contextoBackgroud, erro)

	instanciarDAOs(conexao)
	instanciarServices()
	inicializarServices(conexao)

	grupoV1 := apiRoot.Group("/api/v1")

	inicializarController(grupoV1)
}

func configurarAPI(contextoBackgroud context.Context, erro error) {
	if len(bytesChavePrivada) == 0 || len(bytesChavePublica) == 0 {
		internal.Exception(contextoBackgroud, "As chaves JWT não foram encontradas", nil)
		os.Exit(1)
	}

	erro = internal.CarregarChavesJWT(bytesChavePrivada, bytesChavePublica)
	if erro != nil {
		internal.Exception(contextoBackgroud, "Ocorreu um erro ao tentar carregar as chaves JWT", &erro)
		os.Exit(1)
	}
}

func instanciarDAOs(conexao *conexao_mysql.MySQL) {
	produtoDao = produto_dao.InstanciarProdutoDao(conexao)
	departamentoDao = departamento_dao.InstanciarDepartamentoDao(conexao)
	usuarioDao = usuario_dao.InstanciarUsuarioDao(conexao)
}

func instanciarServices() {
	produtoApiEmbeddingService = api.InstanciarProdutoApiEmbeddingService()
	produtoService = service.InstanciarProdutoService()
	departamentoService = service.InstanciarDepartamentoService()
	usuarioService = service.InstanciarUsuarioService()
	autenticacaoService = service.InstanciarAutenticacaoService()
	tokenService = service.InstanciarTokenService()
	criptografiaService = service.InstanciarCriptografiaService()
}

func inicializarServices(conexao *conexao_mysql.MySQL) {
	produtoApiEmbeddingService.Inicializar(state.URL_API_EMBEDDING, state.UTILIZAR_API_EMBEDDING)
	produtoService.Inicializar(produtoDao, produtoApiEmbeddingService, departamentoService, conexao)

	departamentoService.Inicializar(departamentoDao, conexao)

	usuarioService.Inicializar(usuarioDao, criptografiaService, conexao)

	autenticacaoService.Inicializar(usuarioService, tokenService, criptografiaService)
}

func inicializarController(grupoV1 *echo.Group) {
	controller.InstanciarProdutoController(grupoV1, produtoService)
	controller.InstanciarDepartamentoController(grupoV1, departamentoService)
	controller.InstanciarUsuarioController(grupoV1, usuarioService)
	controller.InstanciarAutenticacaoController(grupoV1, autenticacaoService)
}
