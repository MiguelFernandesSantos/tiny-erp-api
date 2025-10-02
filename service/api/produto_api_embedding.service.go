package api

import (
	log "api-produtos/internal"
	"api-produtos/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ProdutoApiEmbeddingService struct {
	urlApi      string
	clienteHttp *http.Client
	estaAtiva   bool
}

func InstanciarProdutoApiEmbeddingService() *ProdutoApiEmbeddingService {
	return &ProdutoApiEmbeddingService{}
}

func (s *ProdutoApiEmbeddingService) Inicializar(urlApi string, estaAtiva bool) {
	s.urlApi = urlApi
	s.clienteHttp = &http.Client{}
	s.estaAtiva = estaAtiva
}

func (s *ProdutoApiEmbeddingService) AdicionarProduto(contexto context.Context, produto *model.Produto) error {
	if !s.estaAtiva {
		return nil
	}

	log.Start(contexto, "Iniciando requisicao para salvar o produto")

	produtoAsBytes, erro := json.Marshal(produto)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar obter o produto como json", &erro)
		return erro
	}

	produtoAsString := string(produtoAsBytes)

	requisicao, erro := http.NewRequest(http.MethodPost, s.obterUrlAPI("produto"), strings.NewReader(produtoAsString))
	if erro != nil {
		return erro
	}

	s.adicionarHeaderJsonType(requisicao)

	resposta, erro := s.clienteHttp.Do(requisicao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar fazer a requisicao", &erro)
		return erro
	}

	defer resposta.Body.Close()

	if s.retornouStatusErro(resposta) {
		return fmt.Errorf("A API retornou o status de erro " + strconv.Itoa(resposta.StatusCode))
	}

	log.Success(contexto, "O produto foi salvo na API com sucesso")
	return nil
}

func (s *ProdutoApiEmbeddingService) RemoverProduto(contexto context.Context, plu int64) error {
	if !s.estaAtiva {
		return nil
	}

	pluAsString := strconv.FormatInt(plu, 10)
	requisicao, erro := http.NewRequest(http.MethodDelete, s.obterUrlAPI("produto/"+pluAsString), nil)
	if erro != nil {
		return erro
	}

	resposta, erro := s.clienteHttp.Do(requisicao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar fazer a requisicao", &erro)
		return erro
	}

	defer resposta.Body.Close()

	if s.retornouStatusErro(resposta) {
		return fmt.Errorf("A API retornou o status de erro " + strconv.Itoa(resposta.StatusCode))
	}

	log.Success(contexto, "O produto foi removido na API com sucesso")
	return nil
}

func (s *ProdutoApiEmbeddingService) BuscarProdutosSemelhantes(contexto context.Context, plu int64) ([]int64, error) {
	if !s.estaAtiva {
		return nil, nil
	}

	pluAsString := strconv.FormatInt(plu, 10)
	requisicao, erro := http.NewRequest(http.MethodGet, s.obterUrlAPI("produto/"+pluAsString+"/semelhantes"), nil)
	if erro != nil {
		return nil, erro
	}

	resposta, erro := s.clienteHttp.Do(requisicao)
	if erro != nil {
		log.Exception(contexto, "Ocorreu um erro ao tentar fazer a requisicao", &erro)
		return nil, erro
	}

	defer resposta.Body.Close()

	if s.retornouStatusErro(resposta) {
		return nil, fmt.Errorf("A API retornou o status de erro " + strconv.Itoa(resposta.StatusCode))
	}

	var respostaJsonBody []int64
	if erro := json.NewDecoder(resposta.Body).Decode(&respostaJsonBody); erro != nil {
		return nil, erro
	}

	return respostaJsonBody, nil
}

func (s *ProdutoApiEmbeddingService) retornouStatusErro(resposta *http.Response) bool {
	return resposta.StatusCode < 200 || resposta.StatusCode > 299
}

func (s *ProdutoApiEmbeddingService) adicionarHeaderJsonType(requisicao *http.Request) {
	s.adicionarHeader(requisicao, "Content-Type", "application/json")
}

func (s *ProdutoApiEmbeddingService) adicionarHeader(requisicao *http.Request, chave string, valor string) {
	requisicao.Header.Set(chave, valor)
}

func (s *ProdutoApiEmbeddingService) obterUrlAPI(complementoURL string) string {
	return s.urlApi + "/api/v1/" + complementoURL
}
