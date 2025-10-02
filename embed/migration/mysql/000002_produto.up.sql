CREATE TABLE IF NOT EXISTS produto
(
    plu                 BIGINT         NOT NULL PRIMARY KEY,
    descricao           VARCHAR(75)    NOT NULL DEFAULT '',
    descricao_resumida  VARCHAR(50)    NOT NULL DEFAULT '',
    codigo_departamento BIGINT         NOT NULL,
    codigo_secao        BIGINT         NOT NULL,
    pesavel             BOOLEAN                 DEFAULT FALSE,
    bloqueado_entrada   BOOLEAN                 DEFAULT FALSE,
    bloqueado_saida     BOOLEAN                 DEFAULT FALSE,
    custo               DECIMAL(12, 3) NOT NULL DEFAULT 0.0,
    preco_venda         DECIMAL(12, 3) NOT NULL DEFAULT 0.0,
    estoque             DECIMAL(12, 3) NOT NULL DEFAULT 0.0,
    estoque_minimo      DECIMAL(12, 3),
    estoque_maximo      DECIMAL(12, 3),
    data_cadastro       DATETIME       NOT NULL,
    data_alteracao      DATETIME       NOT NULL,
    CONSTRAINT codigo_departamento_produto FOREIGN KEY (codigo_departamento) REFERENCES departamento (codigo),
    CONSTRAINT codigo_secao_produto FOREIGN KEY (codigo_secao) REFERENCES secao (codigo)
);

CREATE TABLE IF NOT EXISTS codigo_barras
(
    plu           BIGINT      NOT NULL,
    codigo_barras VARCHAR(32) NOT NULL,
    CONSTRAINT plu_codigo_barra FOREIGN KEY (plu) REFERENCES produto (plu),
    CONSTRAINT codigo_barras_unico UNIQUE (codigo_barras)
);

CREATE TABLE IF NOT EXISTS imagem_produto
(
    codigo     BIGINT PRIMARY KEY AUTO_INCREMENT,
    plu        BIGINT NOT NULL,
    url_imagem TEXT   NOT NULL,
    CONSTRAINT plu_imagem_produto FOREIGN KEY (plu) REFERENCES produto (plu)
);
