package node_storage

import (
	"database/sql"
	"fmt"
	"strconv"

	guid "github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	_ "github.com/jackc/pgx/stdlib"
)

func NewPathStorage(dbType string, uri string) (*PathStorage, error) {
	db, err := sql.Open(dbType, uri)
	if err != nil {
		return nil, err
	}
	table := "nodes.Paths"
	createSchema := `CREATE SCHEMA IF NOT EXISTS nodes`
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		node_id uuid REFERENCES nodes.Nodes (id),
		path text PRIMARY KEY,
		size bigint NOT NULL
		)`, table)

	_, err = db.Exec(createSchema)
	if err != nil {
		dbErr := db.Close()
		if dbErr != nil {
			err = multierror.Append(err, dbErr)
		}
		return nil, err
	}
	_, err = db.Exec(createTableQuery)
	if err != nil {
		dbErr := db.Close()
		if dbErr != nil {
			err = multierror.Append(err, dbErr)
		}
		return nil, err
	}

	return &PathStorage{db: db, table: table}, nil
}

type PathStorage struct {
	db    *sql.DB
	table string
}

func (p *PathStorage) Add(nodeId guid.UUID, path string, size int64) error {
	_, err := p.db.Exec(fmt.Sprintf("INSERT INTO %s VALUES ('%s', '%s', %d)", p.table, nodeId.String(), path, size))
	return err
}

func (p *PathStorage) Remove(path string) error {
	_, err := p.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE path = '%s'", p.table, path))
	return err
}

func (p *PathStorage) Replace(path string, nodeId guid.UUID, size int64) error {
	_, err := p.db.Exec(fmt.Sprintf("UPDATE %s SET (nodeId, size) = ('%s', %d) WHERE path = '%s'", p.table, nodeId.String(), size, path))
	return err
}

func (p *PathStorage) GetAll() ([]guid.UUID, []string, []int64, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", p.table))

	if row.Err() != nil {
		return nil, nil, nil, row.Err()
	}
	rowAmountStr := ""
	err := row.Scan(&rowAmountStr)
	if err != nil {
		return nil, nil, nil, err
	}
	rows, err := p.db.Query(fmt.Sprintf("SELECT * FROM %s", p.table))
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	var (
		nodeId guid.UUID
		path   string
		size   int64
	)

	rowAmount, err := strconv.ParseInt(rowAmountStr, 0, 64)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeIds := make([]guid.UUID, rowAmount)
	paths := make([]string, rowAmount)
	sizes := make([]int64, rowAmount)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&nodeId, &path, &size)
		if err != nil {
			return nil, nil, nil, err
		}
		nodeIds[i] = nodeId
		paths[i] = path
		sizes[i] = size
	}

	err = row.Err()
	if err != nil {
		return nil, nil, nil, err
	}
	return nodeIds, paths, sizes, nil
}

func (p *PathStorage) Get(path string) (guid.UUID, int64, error) {
	row := p.db.QueryRow(fmt.Sprintf(`SELECT node_id, size FROM %s WHERE path = '%s'`, p.table, path))
	if row.Err() != nil {
		return guid.Nil, 0, row.Err()
	}

	var (
		nodeId guid.UUID
		size   int64
	)
	err := row.Scan(&nodeId, &size)
	if err != nil {
		return guid.Nil, 0, row.Err()
	}
	return nodeId, size, nil
}

func (p *PathStorage) GetPathsByNodeId(id guid.UUID) ([]string, []int64, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE node_id = '%s'", p.table, id.String()))

	if row.Err() != nil {
		return nil, nil, row.Err()
	}
	rowAmountStr := ""
	err := row.Scan(&rowAmountStr)
	if err != nil {
		return nil, nil, err
	}
	rows, err := p.db.Query(fmt.Sprintf("SELECT path, size FROM %s WHERE node_id = '%s'", p.table, id.String()))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var (
		path string
		size int64
	)

	rowAmount, err := strconv.ParseInt(rowAmountStr, 0, 64)
	if err != nil {
		return nil, nil, err
	}

	paths := make([]string, rowAmount)
	sizes := make([]int64, rowAmount)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&path, &size)
		if err != nil {
			return nil, nil, err
		}
		paths[i] = path
		sizes[i] = size
	}

	err = row.Err()
	if err != nil {
		return nil, nil, err
	}
	return paths, sizes, nil
}
