package userutil

import (
	"context"
)

// User has the same fields as the built-in User type
type User struct {
	Uid          string
	Gid          string
	UserInfo     string
	Username     string
	HomeDir      string
	DefaultShell string
}

// Group has the same fields as the built-in Group type,
// plus an additional string slice to represent user membership
type Group struct {
	Gid   string
	Name  string
	Users []string
}

/*
Users lists OS users

Set all=true to return all users

Set all=false to return users without _ and com.apple prefix
*/
func Users(all bool) ([]User, error) {
	return UsersWithContext(context.Background(), all)
}

/*
Users lists OS users. The function expects context.

Set all=true to return all users

Set all=false to return users without _ and com.apple prefix
*/
func UsersWithContext(ctx context.Context, all bool) ([]User, error) {
	return users(ctx, all)
}

/*
Groups lists all groups the user is a member of. The function expects context.

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func (u *User) Groups(all bool) ([]Group, error) {
	return u.GroupsWithContext(context.Background(), all)
}

/*
Groups lists all groups the user is a member of.

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func (u *User) GroupsWithContext(ctx context.Context, all bool) ([]Group, error) {
	return u.membership(ctx, all)
}

/*
Groups lists OS groups

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func Groups(all bool) ([]Group, error) {
	return GroupsWithContext(context.Background(), all, false)
}

/*
Groups lists OS groups. The function expects context.

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func GroupsWithContext(ctx context.Context, all, membership bool) ([]Group, error) {
	return groups(ctx, all, false)
}

/*
GroupsWithMembership lists all groups with users belonging to them.

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func GroupsWithMembership(all bool) ([]Group, error) {
	return groups(context.Background(), all, true)
}

/*
GroupsWithMembership lists all groups with users belonging to them. The function expects context.

Set all=true to return all groups

Set all=false to return groups without _ and com.apple prefix
*/
func GroupsWithMembershipAndContext(ctx context.Context, all bool) ([]Group, error) {
	return groups(ctx, all, true)
}
