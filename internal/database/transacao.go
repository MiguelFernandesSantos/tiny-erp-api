package conexao_mysql

import "database/sql"

type Transacao struct {
	tx *sql.Tx
}

func (t *Transacao) Select(query string, parametros ...any) ([]map[string]any, error) {
	// Executa a query direto no tx
	rows, err := t.tx.Query(query, parametros...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // fecha apenas depois de ler os resultados

	// Pega os nomes e tipos das colunas
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	return t.lerResultados(rows, columnNames, columnTypes)
}

func (t *Transacao) lerResultados(rows *sql.Rows, columnNames []string, columnTypes []*sql.ColumnType) ([]map[string]any, error) {
	var resultados []map[string]any

	for rows.Next() {
		linha, err := t.lerLinha(rows, columnNames, columnTypes)
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

func (t *Transacao) lerLinha(rows *sql.Rows, columnNames []string, columnTypes []*sql.ColumnType) (map[string]any, error) {
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
func (transacao *Transacao) Update(query string, parametros ...any) error {
	var (
		statement *sql.Stmt
		erro      error
	)
	if statement, erro = transacao.tx.Prepare(query); erro != nil {
		return erro
	}
	defer statement.Close()
	_, erro = statement.Exec(parametros...)
	return erro
}

func (t *Transacao) Insert(query string, parametros ...any) (*int64, error) {
	var (
		statement *sql.Stmt
		err       error
	)
	if statement, err = t.tx.Prepare(query); err != nil {
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

func (transacao *Transacao) Prepare(query string) (*sql.Stmt, error) {
	return transacao.tx.Prepare(query)
}

func (transacao *Transacao) commit() error {
	return transacao.tx.Commit()
}

func (transacao *Transacao) rollback() error {
	return transacao.tx.Rollback()
}
