// +build !windows

package userutil

import (
	"bytes"
	"context"
	"os/exec"
)

func (u User) membership(ctx context.Context, all bool) ([]Group, error) {
	stdout, err := exec.CommandContext(ctx, "id", u.Username).Output()
	if err != nil {
		return nil, err
	}

	expression := "groups="
	// Finding the firs occurrence of "groups="
	firstPos := bytes.Index(stdout, []byte(expression))
	// only keep stuff after the expression (gid and group names)
	cropped := stdout[firstPos+len(expression):]

	return parseMembership(cropped, all), nil
}

func parseMembership(content []byte, all bool) []Group {
	group := make([]Group, 0)

	// gid(name),gid(name),gid(name)...
	for _, tuple := range bytes.Split(content, []byte(",")) {
		len := len(tuple)
		idx := bytes.Index(tuple, []byte("("))

		// number without group name
		// groups=1,2(group1),3(group2)...
		if idx < 0 {
			continue
		}

		name := string(tuple[idx+1 : len-1])
		if !all && toFilter(name) {
			continue
		}

		group = append(group, Group{
			Gid:  string(tuple[:idx]),
			Name: name,
		})
	}
	return group
}
