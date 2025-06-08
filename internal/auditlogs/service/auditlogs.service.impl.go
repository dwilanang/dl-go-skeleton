package service

import (
	"fmt"

	"github.com/dwilanang/psp/internal/auditlogs/dto"
	"github.com/dwilanang/psp/internal/auditlogs/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) *service {
	return &service{repo: r}
}

// GetAuditLog implements the Service interface.
func (s *service) GetAuditLog(page, limit int) (*dto.AuditLogResponse, error) {
	var result dto.AuditLogResponse
	totalRecords, err := s.repo.CountAuditLog()
	if err != nil {
		fmt.Println("s.repo.CountAuditLog() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.AuditLogData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.AuditLogData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	results, err := s.repo.FetchAuditLog(limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchAuditLog() error: ", err)
		return &result, err
	}

	if results == nil {
		result.Data = []dto.AuditLogData{}
		return &result, nil
	}

	for _, v := range results {
		result.Data = append(result.Data, dto.AuditLogData{
			ID:        v.ID,
			RequestID: v.RequestID,
			FullName:  v.FullName,
			Action:    v.Action,
			Tablename: v.TableName,
			RecordID:  v.RecordID,
			IpAddress: v.IpAddress,
		})
	}

	return &result, nil
}
