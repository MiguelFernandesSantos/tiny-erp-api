package model

import "time"

type DataHora struct {
	Valor time.Time
}

type Data struct {
	Valor time.Time
}

func DataHoraAtual() *DataHora {
	localidade, _ := time.LoadLocation("America/Sao_Paulo")
	return &DataHora{Valor: time.Now().In(localidade)}
}

func DataAtual() *Data {
	localidade, _ := time.LoadLocation("America/Sao_Paulo")
	agora := time.Now().In(localidade)
	data := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, localidade)
	return &Data{Valor: data}
}

func (d *DataHora) FormatarPadraoBrasil() string {
	return d.Valor.Format("02/01/2006 15:04:05")
}

func (d *DataHora) FormatarISO8601() string {
	return d.Valor.Format("2006-01-02 15:04:05")
}

func (d *Data) FormatarPadraoBrasil() string {
	return d.Valor.Format("02/01/2006")
}

func (d *Data) FormatarISO8601() string {
	return d.Valor.Format("2006-01-02")
}
