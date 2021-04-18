package userutil

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"
)

var filters = []string{"_", "com.apple"}

func users(ctx context.Context, all bool) ([]User, error) {
	result, err := dscacheutil(ctx, "user")
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(result)
	return parseUsers(reader, all)
}

func groups(ctx context.Context, all bool, membership bool) ([]Group, error) {
	result, err := dscacheutil(ctx, "group")
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(result)
	return parseGroups(reader, all, membership)
}

func parseUsers(reader io.Reader, all bool) ([]User, error) {
	var (
		users = make([]User, 0)

		// cache already added users to prevent
		// duplicates returned by dscacheutil
		cache = map[string]bool{}

		// Current user index
		current = -1

		// Pass parsing user
		pass = false
	)

	err := Parse(reader, func(ln string) error {
		line := strings.TrimSpace(ln)

		if strings.HasPrefix(line, "name") {
			name := line[6:]

			// Pass if username is filtered or already added
			if (!all && toFilter(name)) || cache[name] {
				pass = true
				return nil
			}

			users = append(users, User{Username: name})
			cache[name] = true
			pass = false
			current++
		}

		if pass {
			// Continue scan loop
			return nil
		}

		switch {
		case strings.HasPrefix(line, "uid"):
			users[current].Uid = line[5:]
		case strings.HasPrefix(line, "gid"):
			users[current].Gid = line[5:]
		case strings.HasPrefix(line, "dir"):
			users[current].HomeDir = line[5:]
		case strings.HasPrefix(line, "shell"):
			users[current].DefaultShell = line[7:]
		case strings.HasPrefix(line, "gecos"):
			users[current].UserInfo = line[7:]
		}

		return nil
	})

	return users, err
}

func parseGroups(reader io.Reader, all bool, membership bool) ([]Group, error) {

	var (
		groups = make([]Group, 0)

		// cache already added users to prevent
		// duplicates returned by dscacheutil
		cache = map[string]bool{}

		// Current user index
		current = -1

		// Pass parsing user
		pass = false
	)

	err := Parse(reader, func(ln string) error {
		line := strings.TrimSpace(ln)

		if strings.HasPrefix(line, "name") {
			name := line[6:]

			// Pass if group name is filtered or already added
			if (!all && toFilter(name)) || cache[name] {
				pass = true
				return nil
			}

			groups = append(groups, Group{Name: name})
			cache[name] = true
			pass = false
			current++
		}

		if pass {
			// Continue scan loop
			return nil
		}

		switch {
		case strings.HasPrefix(line, "gid"):
			groups[current].Gid = line[5:]
		case strings.HasPrefix(line, "users"):
			if membership {
				groups[current].Users = strings.Split(line[7:], " ")
			}
		}

		return nil
	})

	return groups, err
}

func dscacheutil(ctx context.Context, mode string) ([]byte, error) {
	return exec.CommandContext(ctx, "dscacheutil", "-q", mode).Output()
}
