package constants

type Role int

const (
	ROLE_SUPER_ADMIN Role = 9
	ROLE_ADMIN       Role = 8
	ROLE_USER        Role = 1
	ROLE_GUEST       Role = 0
)
