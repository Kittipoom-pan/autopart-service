package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbAdminToAdminRes(dbAdmin db.GetAdminUserRow) *AdminRes {
	return &AdminRes{
		ID:       int32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    dbAdmin.Email.String,
		Role:     string(dbAdmin.Role),
	}
}

func MapDbAdminToAdminEntity(dbAdmin db.GetAdminByUsernameRow) *Admin {
	return &Admin{
		ID:       int32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Role:     string(dbAdmin.Role),
		Password: dbAdmin.Password,
	}
}

func MapAdminToAdminParam(admin *AdminReq, createBy *int) db.CreateAdminParams {
	return db.CreateAdminParams{
		Username:         admin.Username,
		Email:            utils.StringToNullString(admin.Email),
		Role:             db.AdminUserRole(admin.Role),
		Password:         admin.Password,
		CreatedByAdminID: utils.IntNullToNullInt32(createBy),
		CreatedAt:        utils.NowUTC(),
	}
}

func MapDbAdminsToAdminEntity(dbAdmin db.ListAdminUsersRow) *AdminRes {
	return &AdminRes{
		ID:       int32(dbAdmin.AdminUserID),
		Username: dbAdmin.Username,
		Email:    utils.NullStringToString(dbAdmin.Email),
		Role:     string(dbAdmin.Role),
	}
}

func MapUpdateAdminParams(id int, admin *AdminReq, UpdatedByAdminID *int) db.UpdateAdminParams {
	return db.UpdateAdminParams{
		AdminUserID:      int32(id),
		Username:         admin.Username,
		Password:         admin.Password,
		Role:             db.AdminUserRole(admin.Role),
		Email:            utils.StringToNullString(admin.Email),
		UpdatedByAdminID: utils.IntNullToNullInt32(UpdatedByAdminID),
		UpdatedAt:        utils.NullTimeNow(),
	}
}

func MapUpdateAdminIsActiveParams(id int, isActive bool, UpdatedByAdminID *int) db.UpdateAdminIsActiveParams {
	return db.UpdateAdminIsActiveParams{
		IsActive:         isActive,
		UpdatedByAdminID: utils.IntNullToNullInt32(UpdatedByAdminID),
		UpdatedAt:        utils.NullTimeNow(),
		AdminUserID:      int32(id),
	}
}

func MapAdminToLoginRes(admin *Admin, token string, expiresIn int32) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}
}
