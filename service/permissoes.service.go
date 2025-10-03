package service

import permissoes_dao "api-produtos/dao/permissoes"

type PermissoesService struct {
	permissoesDao *permissoes_dao.PermissoesDao
}

func InstanciarPermissoesService() *PermissoesService {
	return &PermissoesService{}
}

func (s *PermissoesService) Inicializar(permissoesDao *permissoes_dao.PermissoesDao) {
	s.permissoesDao = permissoesDao
}
