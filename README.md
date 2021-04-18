# Userutil

This is a Go library for listing OS user & group info without the use of cgo, so the library can be used in cross-compilation environments.

What is wrong with `os/user`? The built-in os/user package requires cgo on Darwin systems. This means that any Go code that uses that package cannot cross compile. This library provides user & group information without touching the cgo subsytem, hence enabling cross-compilation.

## Doc
```golang
func Users(all bool) ([]User, error)
func UsersWithContext(ctx context.Context, all bool) ([]User, error)
func (u *User) Groups(all bool) ([]Group, error) 
func (u *User) GroupsWithContext(ctx context.Context, all bool) ([]Group, error) 
func Groups(all bool) ([]Group, error) 
func GroupsWithContext(ctx context.Context, all, membership bool) ([]Group, error)
func GroupsWithMembership(all bool) ([]Group, error) 
func GroupsWithMembershipAndContext(ctx context.Context, all bool) ([]Group, error)
```

## avaibility matrix
| os                | Linux | FreeBSD | Plan9 | MacOSX | Windows |
|-------------------|-------|---------|-------|--------|---------|
| implemented       |     x |       x |    x  |    x   |         |
