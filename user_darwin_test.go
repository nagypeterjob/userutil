package userutil

import (
	"reflect"
	"strings"
	"testing"
)

const dsclUsers = `name: root
password: *
uid: 0
gid: 0
dir: /var/root
shell: /bin/sh
gecos: System Administrator
		
name: nagypeter
password: ********
uid: 501
gid: 20
dir: /Users/nagypeter
shell: /bin/zsh
gecos: Nagy Peter

name: _fud
password: *
uid: 278
gid: 278
dir: /var/db/fud
shell: /usr/bin/false
gecos: Firmware Update Daemon

name: _knowledgegraphd
password: *
uid: 279
gid: 279
dir: /var/db/knowledgegraphd
shell: /usr/bin/false
gecos: Knowledge Graph Daemon

name: _coreml
password: *
uid: 280
gid: 280
dir: /var/empty
shell: /usr/bin/false
gecos: CoreML Services
`

const dsclGroups = `
name: admin
password: *
gid: 80
users: root 

name: com.apple.access_ssh
password: *
gid: 399

name: com.apple.access_remote_ae
password: *
gid: 400

name: _oahd
password: *
gid: 441
users: _oahd
`

func TestParseUser(t *testing.T) {
	tcs := []struct {
		all      bool
		expected []User
	}{
		{
			// return all users
			true,
			[]User{
				{Username: "root", Uid: "0", Gid: "0", UserInfo: "System Administrator", DefaultShell: "/bin/sh", HomeDir: "/var/root"},
				{Username: "nagypeter", Uid: "501", Gid: "20", UserInfo: "Nagy Peter", DefaultShell: "/bin/zsh", HomeDir: "/Users/nagypeter"},
				{Username: "_fud", Uid: "278", Gid: "278", UserInfo: "Firmware Update Daemon", DefaultShell: "/usr/bin/false", HomeDir: "/var/db/fud"},
				{Username: "_knowledgegraphd", Uid: "279", Gid: "279", UserInfo: "Knowledge Graph Daemon", DefaultShell: "/usr/bin/false", HomeDir: "/var/db/knowledgegraphd"},
				{Username: "_coreml", Uid: "280", Gid: "280", UserInfo: "CoreML Services", DefaultShell: "/usr/bin/false", HomeDir: "/var/empty"},
			},
		},
		{
			// filter util users
			false,
			[]User{
				{Username: "root", Uid: "0", Gid: "0", UserInfo: "System Administrator", DefaultShell: "/bin/sh", HomeDir: "/var/root"},
				{Username: "nagypeter", Uid: "501", Gid: "20", UserInfo: "Nagy Peter", DefaultShell: "/bin/zsh", HomeDir: "/Users/nagypeter"},
			},
		},
	}

	for i, tc := range tcs {
		reader := strings.NewReader(dsclUsers)
		users, err := parseUsers(reader, tc.all)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(users, tc.expected) {
			t.Errorf("[%d] TestParseUser failed. Wanted: %v, got %v", i, tc.expected, users)
		}
	}
}

func TestParseGroup(t *testing.T) {
	tcs := []struct {
		all        bool
		membership bool
		expected   []Group
	}{
		{
			// return all groups, do not include membership
			true,
			false,
			[]Group{
				{Name: "admin", Gid: "80", Users: nil},
				{Name: "com.apple.access_ssh", Gid: "399", Users: nil},
				{Name: "com.apple.access_remote_ae", Gid: "400", Users: nil},
				{Name: "_oahd", Gid: "441", Users: nil},
			},
		},
		{
			// filter util users, do not include membership
			false,
			false,
			[]Group{
				{Name: "admin", Gid: "80", Users: nil},
			},
		},
		{
			// filter util users, include membership
			false,
			true,
			[]Group{
				{Name: "admin", Gid: "80", Users: []string{"root"}},
			},
		},
	}

	for i, tc := range tcs {
		reader := strings.NewReader(dsclGroups)
		groups, err := parseGroups(reader, tc.all, tc.membership)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(groups, tc.expected) {
			t.Errorf("[%d] TestParseGroup failed. Wanted: %v, got %v", i, tc.expected, groups)
		}
	}
}
