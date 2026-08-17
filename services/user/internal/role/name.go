package role

// Name is a role name in the authorization domain. The Role Authority never
// emits a name outside this set (CONTEXT.md: Role).
type Name string

const (
	Student       Name = "student"
	Admin         Name = "admin"
	PlatformAdmin Name = "platform_admin"
)

// ValidName returns the canonical Name for s if s is in the domain.
func ValidName(s string) (Name, bool) {
	switch s {
	case string(Student), string(Admin), string(PlatformAdmin):
		return Name(s), true
	}
	return "", false
}
