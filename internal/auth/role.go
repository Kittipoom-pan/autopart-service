package auth

const (
	RoleSuperAdmin = "super_admin"
	RoleStaff      = "staff"
	RoleCustomer   = "customer"
)

func IsAdminRole(role string) bool {
	return role == RoleSuperAdmin || role == RoleStaff
}
