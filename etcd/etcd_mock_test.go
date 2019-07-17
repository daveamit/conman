package etcd_test

import (
	"context"

	etcdCli "go.etcd.io/etcd/clientv3"
)

func newAuthMock() *authMock {
	return &authMock{}
}

type authMock struct {
	authEnable           error
	authDisable          error
	userGet              error
	userAdd              error
	userDelete           error
	userChangePassword   error
	userGrantRole        error
	userList             error
	userRevokeRole       error
	roleAdd              error
	roleGrantPermission  error
	roleGet              error
	roleList             error
	roleRevokePermission error
	roleDelete           error
}

// AuthEnable enables auth of an etcd cluster.
func (a *authMock) AuthEnable(ctx context.Context) (*etcdCli.AuthEnableResponse, error) {
	if a.authEnable != nil {
		return nil, a.authEnable
	}
	return nil, nil
}

// AuthDisable disables auth of an etcd cluster.
func (a *authMock) AuthDisable(ctx context.Context) (*etcdCli.AuthDisableResponse, error) {
	if a.authDisable != nil {
		return nil, a.authDisable
	}
	return nil, nil
}

// UserAdd adds a new user to an etcd cluster.
func (a *authMock) UserAdd(ctx context.Context, name string, password string) (*etcdCli.AuthUserAddResponse, error) {
	if a.userAdd != nil {
		return nil, a.userAdd
	}
	return nil, nil
}

// UserDelete deletes a user from an etcd cluster.
func (a *authMock) UserDelete(ctx context.Context, name string) (*etcdCli.AuthUserDeleteResponse, error) {
	if a.userDelete != nil {
		return nil, a.userDelete
	}
	return nil, nil
}

// UserChangePassword changes a password of a user.
func (a *authMock) UserChangePassword(ctx context.Context, name string, password string) (*etcdCli.AuthUserChangePasswordResponse, error) {
	if a.userChangePassword != nil {
		return nil, a.userChangePassword
	}
	return nil, nil
}

// UserGrantRole grants a role to a user.
func (a *authMock) UserGrantRole(ctx context.Context, user string, role string) (*etcdCli.AuthUserGrantRoleResponse, error) {
	if a.userGrantRole != nil {
		return nil, a.userGrantRole
	}
	return nil, nil
}

// UserGet gets a detailed information of a user.
func (a *authMock) UserGet(ctx context.Context, name string) (*etcdCli.AuthUserGetResponse, error) {
	if a.userGet != nil {
		return nil, a.userGet
	}
	return nil, nil
}

// UserList gets a list of all users.
func (a *authMock) UserList(ctx context.Context) (*etcdCli.AuthUserListResponse, error) {
	if a.userList != nil {
		return nil, a.userList
	}
	return nil, nil
}

// UserRevokeRole revokes a role of a user.
func (a *authMock) UserRevokeRole(ctx context.Context, name string, role string) (*etcdCli.AuthUserRevokeRoleResponse, error) {
	if a.userRevokeRole != nil {
		return nil, a.userRevokeRole
	}
	return nil, nil
}

// RoleAdd adds a new role to an etcd cluster.
func (a *authMock) RoleAdd(ctx context.Context, name string) (*etcdCli.AuthRoleAddResponse, error) {
	if a.roleAdd != nil {
		return nil, a.roleAdd
	}
	return nil, nil
}

// RoleGrantPermission grants a permission to a role.
func (a *authMock) RoleGrantPermission(ctx context.Context, name string, key, rangeEnd string, permType etcdCli.PermissionType) (*etcdCli.AuthRoleGrantPermissionResponse, error) {
	if a.roleGrantPermission != nil {
		return nil, a.roleGrantPermission
	}
	return nil, nil
}

// RoleGet gets a detailed information of a role.
func (a *authMock) RoleGet(ctx context.Context, role string) (*etcdCli.AuthRoleGetResponse, error) {
	if a.roleGet != nil {
		return nil, a.roleGet
	}
	return nil, nil
}

// RoleList gets a list of all roles.
func (a *authMock) RoleList(ctx context.Context) (*etcdCli.AuthRoleListResponse, error) {
	if a.roleList != nil {
		return nil, a.roleList
	}
	return nil, nil
}

// RoleRevokePermission revokes a permission from a role.
func (a *authMock) RoleRevokePermission(ctx context.Context, role string, key, rangeEnd string) (*etcdCli.AuthRoleRevokePermissionResponse, error) {
	if a.roleRevokePermission != nil {
		return nil, a.roleRevokePermission
	}
	return nil, nil
}

// RoleDelete deletes a role.
func (a *authMock) RoleDelete(ctx context.Context, role string) (*etcdCli.AuthRoleDeleteResponse, error) {
	if a.roleDelete != nil {
		return nil, a.roleDelete
	}
	return nil, nil
}

func newMockKV() *mockKV {
	return &mockKV{}
}

type mockKV struct {
	put     error
	get     error
	delete  error
	compact error
	do      error
	txn     error
}

func (m *mockKV) Put(ctx context.Context, key, val string, opts ...etcdCli.OpOption) (*etcdCli.PutResponse, error) {
	return nil, m.put
}

func (m *mockKV) Get(ctx context.Context, key string, opts ...etcdCli.OpOption) (*etcdCli.GetResponse, error) {
	return nil, m.get
}

func (m *mockKV) Delete(ctx context.Context, key string, opts ...etcdCli.OpOption) (*etcdCli.DeleteResponse, error) {
	return nil, m.delete
}

func (m *mockKV) Compact(ctx context.Context, rev int64, opts ...etcdCli.CompactOption) (*etcdCli.CompactResponse, error) {
	return nil, m.compact
}

func (m *mockKV) Do(ctx context.Context, op etcdCli.Op) (etcdCli.OpResponse, error) {
	return etcdCli.OpResponse{}, m.do
}

func (m *mockKV) Txn(ctx context.Context) etcdCli.Txn {
	return nil
}
