CREATE TABLE permissoes_usuario
(
    codigo    INT PRIMARY KEY,
    chave     VARCHAR(50)  NOT NULL,
    nome      VARCHAR(100) NOT NULL,
    descricao TEXT         NOT NULL
);

CREATE INDEX idx_permissoes_chave ON permissoes_usuario (chave);

CREATE TABLE IF NOT EXISTS permissoes_ativas_usuario
(
    usuario_codigo   INT NOT NULL,
    permissao_codigo INT NOT NULL,
    PRIMARY KEY (usuario_codigo, permissao_codigo),
    FOREIGN KEY (usuario_codigo) REFERENCES usuario (codigo) ON DELETE CASCADE,
    FOREIGN KEY (permissao_codigo) REFERENCES permissoes_usuario (codigo) ON DELETE CASCADE
);
