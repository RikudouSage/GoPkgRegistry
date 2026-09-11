package db

import (
	"database/sql"
	"fmt"

	"go.chrastecky.dev/go-pkg-repository/cfg"
	"go.chrastecky.dev/go-pkg-repository/dto"
	"go.chrastecky.dev/go-pkg-repository/migrations"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// Client stores and retrieves package metadata in SQLite.
type Client struct {
	db *sql.DB
}

// NewClient opens the configured SQLite database and applies pending migrations.
func NewClient(config *cfg.GlobalConfig) (*Client, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", config.DatabasePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	instance := &Client{
		db: db,
	}
	if err := instance.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return instance, nil
}

func (receiver *Client) migrate() error {
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetBaseFS(migrations.Assets)

	if err := goose.Up(receiver.db, "."); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}

// FindPackageByImportPath returns the package with the given import path.
// It returns (nil, nil) when no package matches.
func (receiver *Client) FindPackageByImportPath(path string) (*dto.Package, error) {
	rows, err := receiver.db.Query("SELECT id, import_path, vcs, repository_url, source_url, source_dir_url, source_file_url from packages where import_path=? limit 1", path)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	if pkg, err := receiver.scanPackage(rows); err != nil {
		return nil, fmt.Errorf("failed to scan package: %w", err)
	} else {
		return pkg, nil
	}
}

// GetPackages returns all stored packages.
func (receiver *Client) GetPackages() ([]*dto.Package, error) {
	rows, err := receiver.db.Query("SELECT id, import_path, vcs, repository_url, source_url, source_dir_url, source_file_url from packages")
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	packages := make([]*dto.Package, 0)
	for rows.Next() {
		pkg, err := receiver.scanPackage(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// StorePackage inserts pkg or updates the existing package with the same import
// path. On success, it sets pkg.ID to the stored database identifier.
func (receiver *Client) StorePackage(pkg *dto.Package) error {
	query := `insert into packages (import_path, vcs, repository_url, source_url, source_dir_url, source_file_url)
				values (?, ?, ?, ?, ?, ?)
				on conflict (import_path) do update set vcs = excluded.vcs, repository_url = excluded.repository_url,
					source_url = excluded.source_url, source_dir_url = excluded.source_dir_url,
					source_file_url = excluded.source_file_url
				returning id
			 `
	row := receiver.db.QueryRow(query, pkg.ImportPath, pkg.VCS, pkg.RepositoryURL, pkg.SourceURL, pkg.SourceDirURL, pkg.SourceFileURL)
	if err := row.Scan(&pkg.ID); err != nil {
		return fmt.Errorf("failed to upsert package: %w", err)
	}

	return nil
}

func (receiver *Client) scanPackage(rows *sql.Rows) (*dto.Package, error) {
	result := dto.Package{}
	if err := rows.Scan(&result.ID, &result.ImportPath, &result.VCS, &result.RepositoryURL, &result.SourceURL, &result.SourceDirURL, &result.SourceFileURL); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &result, nil
}
