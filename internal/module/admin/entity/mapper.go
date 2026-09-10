package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbAdminToAdminRes(dbAdmin db.GetAdminUserRow) *AdminRes {
	return &AdminRes{
		ID:       uint32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    dbAdmin.Email.String,
		Role:     string(dbAdmin.Role),
	}
}

func MapDbAdminToAdminEntity(dbAdmin db.GetAdminByUsernameRow) *Admin {
	return &Admin{
		ID:       uint32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Role:     string(dbAdmin.Role),
		Password: dbAdmin.Password,
	}
}

func MapAdminToAdminParam(admin *AdminReq, createBy *int) db.CreateAdminParams {
	return db.CreateAdminParams{
		Username:  admin.Username,
		Email:     utils.StringToNullString(admin.Email),
		Role:      db.AdminUserRole(admin.Role),
		Password:  admin.Password,
		CreatedBy: utils.IntNullToNullInt32(createBy),
		CreatedAt: utils.NullTimeNow(),
	}
}

func MapDbAdminsToAdminEntity(dbAdmin db.ListAdminUsersRow) *AdminRes {
	return &AdminRes{
		ID:       uint32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    utils.NullStringToString(dbAdmin.Email),
		Role:     string(dbAdmin.Role),
	}
}

func MapUpdateAdminParams(id int, admin *AdminReq, updatedBy *int) db.UpdateAdminParams {
	return db.UpdateAdminParams{
		AdminUserID: int32(id),
		Username:    admin.Username,
		Password:    admin.Password,
		Role:        db.AdminUserRole(admin.Role),
		Email:       utils.StringToNullString(admin.Email),
		UpdatedBy:   utils.IntNullToNullInt32(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdateAdminIsActiveParams(id int, isActive bool, updatedBy *int) db.UpdateAdminIsActiveParams {
	return db.UpdateAdminIsActiveParams{
		IsActive:    isActive,
		UpdatedBy:   utils.IntNullToNullInt32(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
		AdminUserID: int32(id),
	}
}

func MapAdminToLoginRes(admin *Admin, token string, expiresIn int32) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}
}
