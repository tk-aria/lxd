//go:build linux && cgo && !agent

package cluster

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var projectObjects = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  ORDER BY projects.name
`)

var projectObjectsByName = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  WHERE ( projects.name = ? )
  ORDER BY projects.name
`)

var projectObjectsByID = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  WHERE ( projects.id = ? )
  ORDER BY projects.name
`)

var projectCreate = RegisterStmt(`
INSERT INTO projects (description, name)
  VALUES (?, ?)
`)

var projectID = RegisterStmt(`
SELECT projects.id FROM projects
  WHERE projects.name = ?
`)

var projectRename = RegisterStmt(`
UPDATE projects SET name = ? WHERE name = ?
`)

var projectUpdate = RegisterStmt(`
UPDATE projects
  SET description = ?
 WHERE id = ?
`)

var projectDeleteByName = RegisterStmt(`
DELETE FROM projects WHERE name = ?
`)

// projectColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Project entity.
func projectColumns() string {
	return "projects.id, projects.description, projects.name"
}

// getProjects can be used to run handwritten sql.Stmts to return a slice of objects.
func getProjects(ctx context.Context, stmt *sql.Stmt, args ...any) ([]Project, error) {
	objects := make([]Project, 0)

	dest := func(scan func(dest ...any) error) error {
		p := Project{}
		err := scan(&p.ID, &p.Description, &p.Name)
		if err != nil {
			return err
		}

		objects = append(objects, p)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	return objects, nil
}

// getProjects can be used to run handwritten query strings to return a slice of objects.
func getProjectsRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]Project, error) {
	objects := make([]Project, 0)

	dest := func(scan func(dest ...any) error) error {
		p := Project{}
		err := scan(&p.ID, &p.Description, &p.Name)
		if err != nil {
			return err
		}

		objects = append(objects, p)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	return objects, nil
}

// GetProjects returns all available projects.
// generator: project GetMany
func GetProjects(ctx context.Context, tx *sql.Tx, filters ...ProjectFilter) ([]Project, error) {
	var err error

	// Result slice.
	objects := make([]Project, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(tx, projectObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"projectObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Name != nil && filter.ID == nil {
			args = append(args, []any{filter.Name}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, projectObjectsByName)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"projectObjectsByName\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(projectObjectsByName)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"projectObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID != nil && filter.Name == nil {
			args = append(args, []any{filter.ID}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, projectObjectsByID)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"projectObjectsByID\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(projectObjectsByID)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"projectObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID == nil && filter.Name == nil {
			return nil, fmt.Errorf("Cannot filter on empty ProjectFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getProjects(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getProjectsRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	return objects, nil
}

// GetProjectConfig returns all available Project Config
// generator: project GetMany
func GetProjectConfig(ctx context.Context, tx *sql.Tx, projectID int, filters ...ConfigFilter) (map[string]string, error) {
	projectConfig, err := GetConfig(ctx, tx, "project", filters...)
	if err != nil {
		return nil, err
	}

	config, ok := projectConfig[projectID]
	if !ok {
		config = map[string]string{}
	}

	return config, nil
}

// GetProject returns the project with the given key.
// generator: project GetOne
func GetProject(ctx context.Context, tx *sql.Tx, name string) (*Project, error) {
	filter := ProjectFilter{}
	filter.Name = &name

	objects, err := GetProjects(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "Project not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"projects\" entry matches")
	}
}

// ProjectExists checks if a project with the given key exists.
// generator: project Exists
func ProjectExists(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	_, err := GetProjectID(ctx, tx, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateProject adds a new project to the database.
// generator: project Create
func CreateProject(ctx context.Context, tx *sql.Tx, object Project) (int64, error) {
	// Check if a project with the same key exists.
	exists, err := ProjectExists(ctx, tx, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"projects\" entry already exists")
	}

	args := make([]any, 2)

	// Populate the statement arguments.
	args[0] = object.Description
	args[1] = object.Name

	// Prepared statement to use.
	stmt, err := Stmt(tx, projectCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"projectCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"projects\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"projects\" entry ID: %w", err)
	}

	return id, nil
}

// CreateProjectConfig adds new project Config to the database.
// generator: project Create
func CreateProjectConfig(ctx context.Context, tx *sql.Tx, projectID int64, config map[string]string) error {
	referenceID := int(projectID)
	for key, value := range config {
		insert := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err := CreateConfig(ctx, tx, "project", insert)
		if err != nil {
			return fmt.Errorf("Insert Config failed for Project: %w", err)
		}

	}

	return nil
}

// GetProjectID return the ID of the project with the given key.
// generator: project ID
func GetProjectID(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	stmt, err := Stmt(tx, projectID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"projectID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, name)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "Project not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"projects\" ID: %w", err)
	}

	return id, nil
}

// RenameProject renames the project matching the given key parameters.
// generator: project Rename
func RenameProject(ctx context.Context, tx *sql.Tx, name string, to string) error {
	stmt, err := Stmt(tx, projectRename)
	if err != nil {
		return fmt.Errorf("Failed to get \"projectRename\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(to, name)
	if err != nil {
		return fmt.Errorf("Rename Project failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}

	return nil
}

// DeleteProject deletes the project matching the given key parameters.
// generator: project DeleteOne-by-Name
func DeleteProject(ctx context.Context, tx *sql.Tx, name string) error {
	stmt, err := Stmt(tx, projectDeleteByName)
	if err != nil {
		return fmt.Errorf("Failed to get \"projectDeleteByName\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("Delete \"projects\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "Project not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d Project rows instead of 1", n)
	}

	return nil
}
