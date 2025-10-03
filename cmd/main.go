package main

import (
	api_produtos "api-produtos"
	internal "api-produtos/internal"
	conexao_mysql "api-produtos/internal/database"
	"api-produtos/internal/echo_config"
	"api-produtos/internal/state"
	"api-produtos/model"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "time/tzdata"
)

func main() {
	internal.Inicializar(internal.ConfiguracaoPadrao())

	contextoBackgroud := context.Background()

	internal.Start(contextoBackgroud, "Iniciando API")
	api := echo.New()

	erro := godotenv.Load()
	if erro != nil {
		internal.Info(contextoBackgroud, "Nenhum arquivo .env encontrado")
	}

	internal.Info(contextoBackgroud, "Configurando CORS")
	echo_config.ConfigurarCORS(api)

	internal.Info(contextoBackgroud, "Configurando validador JSON")
	api.Validator = configurarValidadorJson()

	api.Use(echo_config.RequestIDMiddleware)

	parametrosBanco, erro := obterParametrosBanco()
	if erro != nil {
		internal.Exception(contextoBackgroud, "Ocorreu um erro ao tentar obter os dados do banco de dados", &erro)
		os.Exit(1)
	}

	conexaoBanco, erro := conexao_mysql.NovaConexao(parametrosBanco.Host, parametrosBanco.Usuario, parametrosBanco.Senha, parametrosBanco.NomeBanco)
	if erro != nil {
		internal.Exception(contextoBackgroud, "Ocorreu um erro ao tentar conectar com o banco de dados com os dados fornecedos", &erro)
		os.Exit(1)
	}

	state.URL_API_EMBEDDING, erro = obterVariavelAmbiente("URL_API_EMBEDDING")
	if erro != nil {
		internal.Exception(contextoBackgroud, "Nao foi encontrado a URL da api de embedding", &erro)
		os.Exit(1)
	}

	utilizaEmbedding, erro := obterVariavelAmbiente("UTILIZAR_API_EMBEDDING")
	if erro != nil {
		internal.Exception(contextoBackgroud, "Nao foi encontrado a variavel de ambiente que define se utiliza a api de embedding", &erro)
		os.Exit(1)
	}

	if strings.ToLower(utilizaEmbedding) == "true" {
		state.UTILIZAR_API_EMBEDDING = true
	}

	api_produtos.RegistrarDependencias(contextoBackgroud, api, conexaoBanco)

	porta, erro := obterVariavelAmbiente("PORTA_API")
	if erro != nil {
		internal.Info(contextoBackgroud, "Variavel de ambiente referente a porta não encontrada, subindo como porta padrao 8081")
		porta = "8081"
	}

	api.Logger.Fatal(api.Start(":" + porta))
}

func configurarValidadorJson() *echo_config.JSONValidator {
	validatorJson := validator.New()

	return &echo_config.JSONValidator{Validator: validatorJson}
}

func obterParametrosBanco() (*model.ParametrosConexaoBanco, error) {
	nomeBanco, erro := obterVariavelAmbiente("NOME_BANCO")
	if erro != nil {
		return nil, erro
	}

	urlBanco, erro := obterVariavelAmbiente("URL_BANCO")
	if erro != nil {
		return nil, erro
	}

	senhaBanco, erro := obterVariavelAmbiente("SENHA_BANCO")
	if erro != nil {
		return nil, erro
	}

	usuarioBanco, erro := obterVariavelAmbiente("USUARIO_BANCO")
	if erro != nil {
		return nil, erro
	}

	return &model.ParametrosConexaoBanco{
		NomeBanco: nomeBanco,
		Host:      urlBanco,
		Senha:     senhaBanco,
		Usuario:   usuarioBanco,
	}, nil
}

func obterVariavelAmbiente(chave string) (string, error) {
	valor, existe := os.LookupEnv(chave)
	if !existe {
		return "", fmt.Errorf("Variavel de ambiente obrigatoria não encontrada: %s", chave)
	}
	return valor, nil
}
