package userutil

import (
	"context"
	"strings"
)

const (
	passwdPath = "/etc/passwd"
	groupsPath = "/etc/groups"
)

var filters = []string{"_"}

func users(ctx context.Context, all bool) ([]User, error) {
	users := []User{}

	err := Open(passwdPath, func(line string) error {
		if user := parsePasswdLine(line, all); user != nil {
			users = append(users, *user)
		}

		return nil
	})
	return users, err
}

func parsePasswdLine(line string, all bool) *User {
	line = strings.TrimSpace(line)

	parts := strings.Split(line, ":")

	if len(parts) >= 7 {

		// Pass if username is filtered
		if !all && toFilter(parts[0]) {
			return nil
		}

		// username:password:uid:gid:userinfo:home:shell
		// [0] username
		// [1] password
		// [2] uid
		// [3] gid
		// [4] userinfo
		// [5] home
		// [6] shell
		return &User{
			Username:     parts[0],
			Uid:          parts[2],
			Gid:          parts[3],
			UserInfo:     parts[4],
			HomeDir:      parts[5],
			DefaultShell: parts[6],
		}
	}
	return nil
}

func groups(ctx context.Context, all bool, membership bool) ([]Group, error) {
	groups := make([]Group, 0)

	err := Open(groupsPath, func(line string) error {
		if group := parseGroupLine(line, all, membership); group != nil {
			groups = append(groups, *group)
		}

		return nil
	})
	return groups, err
}

func parseGroupLine(line string, all bool, membership bool) *Group {
	line = strings.TrimSpace(line)

	parts := strings.Split(line, ":")

	if len(parts) >= 3 {
		// Pass if group name is filtered
		if !all && toFilter(parts[0]) {
			return nil
		}

		// name:x:1:user1,user2,user3
		// [0] name
		// [1] password
		// [2] gid
		// [3:] user
		group := &Group{
			Name: parts[0],
			Gid:  parts[2],
		}

		// include user membership if asked
		if membership && len(parts) > 3 {
			group.Users = append(group.Users, parts[3:]...)
		}

		return group
	}
	return nil
}
