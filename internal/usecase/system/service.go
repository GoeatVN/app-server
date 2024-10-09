package system

import (
	"app-server/internal/shared/systemdto"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type SystemServiceInterface interface {
	LoadComboboxData(request systemdto.ComboboxRequest) (*systemdto.ComboboxResponse, error)
}

type systemService struct {
	dbPool *pgxpool.Pool
	db     *gorm.DB
}

func NewComboboxService(dbPool *pgxpool.Pool, db *gorm.DB) SystemServiceInterface {
	return &systemService{dbPool: dbPool, db: db}
}

func (s *systemService) LoadComboboxData(request systemdto.ComboboxRequest) (*systemdto.ComboboxResponse, error) {
	var response systemdto.ComboboxResponse
	for _, req := range request.Data {
		switch req.ComboType {
		case systemdto.AllRole:
			dataResponseWrapper, err := s.GetAllRoleCombo(req)
			if err != nil {
				return nil, err
			}
			response.Data = append(response.Data, *dataResponseWrapper)
		}
	}
	return nil, nil
}

func (s *systemService) GetAllRoleCombo(requestParam systemdto.ComboboxRequestItem) (*systemdto.ComboboxResponseItem, error) {

	var query = "select r.id, r.role_name as name, r.role_code as value from roles r;"
	var dto systemdto.ComboboxDto
	err := s.db.Raw(query).Scan(&dto).Error

	if err != nil {
		return nil, err
	}
	var dataResponseItem systemdto.ComboboxResponseItem
	dataResponseItem.Data = append(dataResponseItem.Data, dto)
	dataResponseItem.ComboType = requestParam.ComboType
	return &dataResponseItem, nil
}
