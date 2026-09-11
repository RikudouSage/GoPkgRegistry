package dto

import (
	"go.chrastecky.dev/go-pkg-repository/types"
)

// RepositoryType identifies the version control system hosting a package.
type RepositoryType string

const (
	// RepositoryTypeGit identifies a Git repository.
	RepositoryTypeGit RepositoryType = "git"
	// RepositoryTypeMercurial identifies a Mercurial repository.
	RepositoryTypeMercurial RepositoryType = "hg"
	// RepositoryTypeSubversion identifies a Subversion repository.
	RepositoryTypeSubversion RepositoryType = "svn"
	// RepositoryTypeFossil identifies a Fossil repository.
	RepositoryTypeFossil RepositoryType = "fossil"
)

// Package describes a Go package and the repository where its source is hosted.
type Package struct {
	// ID is the package's database identifier.
	ID int `json:"id"`

	// ImportPath is the canonical path used to import the package.
	ImportPath string `json:"import_path"`
	// VCS is the version control system used by the repository.
	VCS RepositoryType `json:"vcs"`
	// RepositoryURL is the URL understood by the version control system.
	RepositoryURL *types.URI `json:"repository_url"`

	// SourceURL is the URL of a web-based source browser.
	SourceURL *types.URI `json:"source_url"`
	// SourceDirURL is the source-browser URL template for a directory.
	SourceDirURL *types.URI `json:"source_dir_url"`
	// SourceFileURL is the source-browser URL template for a file.
	SourceFileURL *types.URI `json:"source_file_url"`
}

// IsValid reports whether all metadata required by the go-import protocol is set.
func (receiver Package) IsValid() bool {
	return receiver.ImportPath != "" && receiver.VCS != "" && receiver.RepositoryURL != nil
}
