package departamento_dao

import "api-produtos/model"

func (d *DepartamentoDao) mapearDepartamento(resultado map[string]any) *model.Departamento {
	departamento := &model.Departamento{
		NumeroDepartamento: resultado["codigo"].(*int64),
		Nome:               resultado["nome"].(*string),
	}
	return departamento
}

func (d *DepartamentoDao) mapearSecao(resultado map[string]any) *model.Secao {
	secao := &model.Secao{
		NumeroSecao:        resultado["codigo"].(*int64),
		NumeroDepartamento: resultado["codigo_departamento"].(*int64),
		Nome:               resultado["nome"].(*string),
	}
	return secao
}
