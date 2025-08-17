package entitie

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbAdminToAdminEntity(dbAdmin db.GetAdminUserRow) *Admin {
	return &Admin{
		ID:       uint32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    dbAdmin.Email.String,
		Role:     string(dbAdmin.Role),
	}
}

func MapAdminToAdminParam(admin *AdminReq, createBy string) db.CreateAdminParams {
	return db.CreateAdminParams{
		Username:  admin.Username,
		Email:     utils.StringToNullString(admin.Email),
		Role:      db.AdminUserRole(admin.Role),
		Password:  admin.Password,
		CreatedBy: utils.StringToNullString(createBy),
		CreatedAt: utils.NullTimeNow(),
	}
}

func MapDbAdminsToAdminEntity(dbAdmin db.ListAdminUsersRow) *Admin {
	return &Admin{
		ID:       uint32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    utils.NullStringToString(dbAdmin.Email),
		Role:     string(dbAdmin.Role),
	}
}

func MapUpdateAdminParams(id int, admin *AdminReq, updatedBy string) db.UpdateAdminParams {
	return db.UpdateAdminParams{
		AdminUserID: int32(id),
		Username:    admin.Username,
		Password:    admin.Password,
		Role:        db.AdminUserRole(admin.Role),
		Email:       utils.StringToNullString(admin.Email),
		UpdatedBy:   utils.StringToNullString(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdateAdminIsActiveParams(id int, isActive bool, updatedBy string) db.UpdateAdminIsActiveParams {
	return db.UpdateAdminIsActiveParams{
		IsActive:    isActive,
		UpdatedBy:   utils.StringToNullString(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
		AdminUserID: int32(id),
	}
}
