// Package common contains functionality not critical to the core project but still essential.
package common

// TODO refactor out the common package

// IClient interface is used with the api clients to hold the repo and org specific info.
type IClient interface {
	GetUserOrganization(login string) (*Owner, error)
	GetRepositoriesFromOwner(target Owner) ([]*Repository, error)
	GetOrganizationMembers(target Owner) ([]*Owner, error)
}
