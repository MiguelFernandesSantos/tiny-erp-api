package conexao_mysql

import (
	"database/sql"
	"embed"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqldriver "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type MySQL struct {
	conexao *sql.DB
}

type TransacaoService interface {
	ExecutarTransacao(function func(*Transacao) error) error
}

func NovaConexao(host string, usuario string, senha string, schema string) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", usuario, senha, host, schema)
	db, erro := sql.Open("mysql", dsn)
	if erro != nil {
		return nil, erro
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return &MySQL{
		conexao: db,
	}, nil
}

func (m *MySQL) Conexao() *sql.DB {
	return m.conexao
}

func (m *MySQL) FecharConexao() error {
	return m.conexao.Close()
}

func (m *MySQL) Select(query string, parametros ...any) ([]map[string]any, error) {
	rows, err := m.executarQuery(query, parametros...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	return m.lerResultados(rows, columnNames, columnTypes)
}

func (m *MySQL) executarQuery(query string, parametros ...any) (*sql.Rows, error) {
	statement, err := m.conexao.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	return statement.Query(parametros...)
}

func (m *MySQL) lerResultados(rows *sql.Rows, columnNames []string, columnTypes []*sql.ColumnType) ([]map[string]any, error) {
	var resultados []map[string]any

	for rows.Next() {
		linha, err := m.lerLinha(rows, columnNames, columnTypes)
		if err != nil {
			return nil, err
		}
		resultados = append(resultados, linha)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resultados, nil
}

func (m *MySQL) lerLinha(rows *sql.Rows, columnNames []string, columnTypes []*sql.ColumnType) (map[string]any, error) {
	rawValues := make([]any, len(columnTypes))
	dest := make([]any, len(columnTypes))

	for i := range rawValues {
		dest[i] = &rawValues[i]
	}

	if err := rows.Scan(dest...); err != nil {
		return nil, err
	}

	linha := make(map[string]any, len(columnTypes))
	for i, colType := range columnTypes {
		linha[columnNames[i]] = converterValor(colType, rawValues[i])
	}

	return linha, nil
}

func (m *MySQL) Update(query string, parametros ...any) error {
	var (
		statement *sql.Stmt
		erro      error
	)
	if statement, erro = m.conexao.Prepare(query); erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(parametros...)
	return erro
}

func (m *MySQL) Insert(query string, parametros ...any) (*int64, error) {
	var (
		statement *sql.Stmt
		err       error
	)
	if statement, err = m.conexao.Prepare(query); err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(parametros...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *MySQL) ExecutarTransacao(funcao func(*Transacao) error) (erro error) {
	var transacao *Transacao
	transacao, erro = m.begin()
	if erro != nil {
		return erro
	}

	defer func() {
		if r := recover(); r != nil {
			_ = transacao.rollback()
			if e, ok := r.(error); ok {
				erro = e
			} else {
				erro = fmt.Errorf("Erro na transacao: %v", r)
			}
		} else if erro != nil {
			_ = transacao.rollback()
		} else {
			erro = transacao.commit()
		}
	}()

	erro = funcao(transacao)
	return erro
}

func (m *MySQL) begin() (*Transacao, error) {
	tx, erro := m.conexao.Begin()
	if erro != nil {
		return nil, erro
	}
	return &Transacao{
		tx: tx,
	}, nil
}

func (m *MySQL) ExecutarMigracao(migrationFS embed.FS, nomeTabelaMigration string) error {
	driver, err := mysqldriver.WithInstance(m.conexao, &mysqldriver.Config{
		MigrationsTable: nomeTabelaMigration,
	})
	if err != nil {
		return fmt.Errorf("Falha ao criar instância do driver MySQL: %w", err)
	}

	source, err := iofs.New(migrationFS, "embed/migration/mysql")
	if err != nil {
		return fmt.Errorf("Falha ao criar fonte de migração: %w", err)
	}

	migration, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("Falha ao criar instância de migração: %w", err)
	}
	defer source.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Falha ao executar migração: %w", err)
	}
	return nil
}

func converterValor(colType *sql.ColumnType, valor any) any {
	if valor == nil {
		return nil
	}

	switch colType.DatabaseTypeName() {
	case "INT", "BIGINT", "SMALLINT", "TINYINT":
		return converterInteiro(valor)
	case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "REAL":
		return converterDecimal(valor)
	case "BOOL", "BOOLEAN":
		return converterBooleano(valor)
	default: // VARCHAR, TEXT, DATE, DATETIME, etc.
		return converterTexto(valor)
	}
}

func converterInteiro(valor any) any {
	switch v := valor.(type) {
	case int64:
		return &v
	case []byte:
		n, _ := strconv.ParseInt(string(v), 10, 64)
		return &n
	default:
		return fmt.Sprintf("%v", v)
	}
}

func converterDecimal(valor any) any {
	switch v := valor.(type) {
	case float64:
		return &v
	case []byte:
		n, _ := strconv.ParseFloat(string(v), 64)
		return &n
	default:
		return fmt.Sprintf("%v", v)
	}
}

func converterBooleano(valor any) any {
	switch v := valor.(type) {
	case bool:
		return &v
	case []byte:
		b, _ := strconv.ParseBool(string(v))
		return &b
	default:
		return fmt.Sprintf("%v", v)
	}
}

func converterTexto(valor any) any {
	switch v := valor.(type) {
	case string:
		return &v
	case []byte:
		conteudo := string(v)
		return &conteudo
	default:
		return fmt.Sprintf("%v", v)
	}
}
