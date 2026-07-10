package user

import "testing"

func TestActionText(t *testing.T) {
	cases := []struct {
		action string
		want   string
	}{
		{ActionCreateUser, "创建用户"},
		{ActionUpdateUser, "更新资料"},
		{ActionSuspendUser, "禁用用户"},
		{ActionSuspendWithReason, "禁用用户"},
		{ActionActivateUser, "启用用户"},
		{ActionDeleteUser, "删除用户"},
		{ActionResetPassword, "重置密码"},
		{"UNKNOWN_ACTION", "UNKNOWN_ACTION"},
	}
	for _, c := range cases {
		if got := ActionText(c.action); got != c.want {
			t.Errorf("ActionText(%q) = %q, want %q", c.action, got, c.want)
		}
	}
}
