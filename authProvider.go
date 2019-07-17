package conman

// User defines user attribs
type User interface {
	Username() string
	Name() string
	Email() string
	Mobile() string
	Meta() map[string]string
	Roles() ([]Role, error)
	AddRole(name string) error
	RemoveRole(name string) error
}

// Role defines a bucket of permissions
type Role interface {
	Name() string
	Permissions() string
	GrantPermission(setting string, exact string)
	RevokePermission(setting string)
}

// AuthProvider is an abstraction over authentication / authorization interface
type AuthProvider interface {
	// Users retrives all users.
	Users() ([]User, error)
	// User retrives user with provided username
	User(username string) (User, error)
	// AddUser adds new user
	AddUser(user User) error
	// UpdateUser updates exsisting user
	UpdateUser(user User) error
	// RemoveUser removes user
	RemoveUser(username string) error

	// Roles retrives all roles
	Roles() ([]Role, error)
	// Role retrives role with given name
	Role(name string) (Role, error)
	// AddRole adds new role
	AddRole(role Role) error
	// UpdateRole updates exsisting role
	UpdateRole(role Role) error
	// RemoveRole Removes a role, if its unbound
	RemoveRole(name string) error

	// GrantRole binds role to user
	GrantRole(username string, role string) error
	// RevokeRole unbinds role from the user
	RevokeRole(username, role string) error
	// UsersWithRole retrives all users with give role
	UsersWithRole(role string) ([]User, error)
	// RolesWithUser retrives all roles for a given user
	RolesWithUser(username string) ([]Role, error)

	// Grant will allow the users with given role to read the setting
	Grant(roleName string, setting string, exact bool) error
	// Remove will remove read access for the setting
	Revoke(roleName string, setting string) error

	// MakeAdmin allows user to add more users
	MakeAdmin(username string) error
	// MakeNormal removes the admins access
	MakeNormal(username string) error
}
