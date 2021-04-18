package userutil

import (
	"reflect"
	"strings"
	"testing"
)

const usersInput = `
www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
sshd:x:109:65534::/run/sshd:/usr/sbin/nologin
_apt:x:104:65534::/nonexistent:/usr/sbin/nologin
vagrant:x:1000:1000:,,,:/home/vagrant:/bin/bash
ubuntu:x:1001:1001:Ubuntu:/home/ubuntu:/bin/bash
`
const groupsInput = `
_admin:x:113:
netdev:x:114:ubuntu
_vagrant:x:1000:
_ubuntu:x:1001:
_dummy:x:1002:
`

func TestUsers(t *testing.T) {
	tcs := []struct {
		all      bool
		expected []User
	}{
		{
			// return all users
			true,
			[]User{
				{Username: "www-data", Uid: "33", Gid: "33", UserInfo: "www-data", DefaultShell: "/usr/sbin/nologin", HomeDir: "/var/www"},
				{Username: "sshd", Uid: "109", Gid: "65534", UserInfo: "", DefaultShell: "/usr/sbin/nologin", HomeDir: "/run/sshd"},
				{Username: "_apt", Uid: "104", Gid: "65534", UserInfo: "", DefaultShell: "/usr/sbin/nologin", HomeDir: "/nonexistent"},
				{Username: "vagrant", Uid: "1000", Gid: "1000", UserInfo: ",,,", DefaultShell: "/bin/bash", HomeDir: "/home/vagrant"},
				{Username: "ubuntu", Uid: "1001", Gid: "1001", UserInfo: "Ubuntu", DefaultShell: "/bin/bash", HomeDir: "/home/ubuntu"},
			},
		},
		{
			// filter util users
			false,
			[]User{
				{Username: "www-data", Uid: "33", Gid: "33", UserInfo: "www-data", DefaultShell: "/usr/sbin/nologin", HomeDir: "/var/www"},
				{Username: "sshd", Uid: "109", Gid: "65534", UserInfo: "", DefaultShell: "/usr/sbin/nologin", HomeDir: "/run/sshd"},
				{Username: "vagrant", Uid: "1000", Gid: "1000", UserInfo: ",,,", DefaultShell: "/bin/bash", HomeDir: "/home/vagrant"},
				{Username: "ubuntu", Uid: "1001", Gid: "1001", UserInfo: "Ubuntu", DefaultShell: "/bin/bash", HomeDir: "/home/ubuntu"},
			},
		},
	}

	for i, tc := range tcs {
		userSlice := make([]User, 0)
		for _, line := range strings.Split(usersInput, "\n") {
			if user := parsePasswdLine(line, tc.all); user != nil {
				userSlice = append(userSlice, *user)
			}
		}

		if !reflect.DeepEqual(userSlice, tc.expected) {
			t.Errorf("[%d] TestUsers failed. Wanted: %v, got %v", i, tc.expected, userSlice)
		}
	}
}

func TestGroups(t *testing.T) {
	tcs := []struct {
		all        bool
		membership bool
		expected   []Group
	}{
		{
			// return all groups without membership
			true,
			false,
			[]Group{
				{Name: "_admin", Gid: "113", Users: nil},
				{Name: "netdev", Gid: "114", Users: nil},
				{Name: "_vagrant", Gid: "1000", Users: nil},
				{Name: "_ubuntu", Gid: "1001", Users: nil},
				{Name: "_dummy", Gid: "1002", Users: nil},
			},
		},
		{
			// filter util users, do not include membership
			false,
			false,
			[]Group{
				{Name: "netdev", Gid: "114", Users: nil},
			},
		},
		{
			// filter util users, include membership
			false,
			true,
			[]Group{
				{Name: "netdev", Gid: "114", Users: []string{"ubuntu"}},
			},
		},
	}

	for i, tc := range tcs {
		groupSlice := make([]Group, 0)
		for _, line := range strings.Split(groupsInput, "\n") {
			if group := parseGroupLine(line, tc.all, tc.membership); group != nil {
				groupSlice = append(groupSlice, *group)
			}
		}

		if !reflect.DeepEqual(groupSlice, tc.expected) {
			t.Errorf("[%d] TestGroups failed. Wanted: %v, got %v", i, tc.expected, groupSlice)
		}
	}
}
