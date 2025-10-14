package enum

type UserRole int

const (
	RoleReader UserRole = iota // 0
	RoleEditor                 // 1
	RoleAuthor                 // 2
	RoleAdmin                  // 3
)

func IsValidRole(role UserRole) bool {
	switch role {
	case RoleReader, RoleEditor, RoleAuthor, RoleAdmin:
		return true
	default:
		return false
	}
}
