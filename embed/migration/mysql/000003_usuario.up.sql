CREATE TABLE IF NOT EXISTS usuario
(
    codigo       INT AUTO_INCREMENT PRIMARY KEY,
    nome         VARCHAR(100) NOT NULL,
    email        VARCHAR(100) NOT NULL,
    senha        TEXT         NOT NULL,
    ativo        BOOLEAN      NOT NULL DEFAULT TRUE,
    perfil       INT          NOT NULL,
    data_criacao TIMESTAMP    NOT NULL,
    ultimo_login TIMESTAMP
);

CREATE INDEX idx_usuario_email ON usuario (nome, email);