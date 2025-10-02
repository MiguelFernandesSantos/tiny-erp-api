CREATE TABLE IF NOT EXISTS departamento
(
    codigo BIGINT PRIMARY KEY AUTO_INCREMENT,
    nome   VARCHAR(100) NOT NULL,
    CONSTRAINT nome_departamento_unico UNIQUE (nome)
);

CREATE TABLE IF NOT EXISTS secao
(
    codigo              BIGINT PRIMARY KEY AUTO_INCREMENT,
    codigo_departamento BIGINT       NOT NULL,
    nome                VARCHAR(100) NOT NULL,
    CONSTRAINT codigo_departamento_secao FOREIGN KEY (codigo_departamento) REFERENCES departamento (codigo)
);