package dto

import (
	"GoPkgRepository/types"
)

type RepositoryType string

const (
	RepositoryTypeGit        RepositoryType = "git"
	RepositoryTypeMercurial  RepositoryType = "hg"
	RepositoryTypeSubversion RepositoryType = "svn"
	RepositoryTypeFossil     RepositoryType = "fossil"
)

type Package struct {
	ID int `json:"id"`

	ImportPath    string         `json:"import_path"`
	VCS           RepositoryType `json:"vcs"`
	RepositoryURL *types.URI     `json:"repository_url"`

	SourceURL     *types.URI `json:"source_url"`
	SourceDirURL  *types.URI `json:"source_dir_url"`
	SourceFileURL *types.URI `json:"source_file_url"`
}

func (receiver Package) IsValid() bool {
	return receiver.ImportPath != "" && receiver.VCS != "" && receiver.RepositoryURL != nil
}
