package node_storage

import (
	"database/sql"
	"fmt"
	guid "github.com/google/uuid"
	"github.com/uroborosq-go-dfs/server/connector"
	"strconv"

	"github.com/hashicorp/go-multierror"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/uroborosq-go-dfs/server/node"
)

func New(dbType string, uri string) (*NodeStorage, error) {
	db, err := sql.Open(dbType, uri)
	if err != nil {
		return nil, err
	}
	nodeTable := "nodes.Nodes"
	createSchema := `CREATE SCHEMA IF NOT EXISTS nodes`
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id uuid PRIMARY KEY,
		ip text NOT NULL,
		port text NOT NULL,
		current_size bigint NOT NULL,
        maximum_size bigint NOT NULL,
        connection_type integer NOT NULL
		)`, nodeTable)

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

	return &NodeStorage{db: db, nodeTable: nodeTable}, nil
}

type NodeStorage struct {
	db        *sql.DB
	nodeTable string
}

func (n *NodeStorage) Add(id guid.UUID, node node.INode) error {
	_, err := n.db.Exec(fmt.Sprintf("INSERT INTO %s VALUES ('%s', '%s', '%s', %d, %d, %d)", n.nodeTable, id.String(), node.GetIp(), node.GetPort(), node.GetCurrentSize(), node.GetMaxSize(), node.GetConnectorType()))
	return err
}

func (n *NodeStorage) Remove(id guid.UUID) error {
	_, err := n.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = '%s'`, n.nodeTable, id.String()))
	return err
}

func (n *NodeStorage) Replace(id guid.UUID, node node.INode) error {
	_, err := n.db.Exec(fmt.Sprintf(`UPDATE %s SET ip = '%s', port = '%s', current_size = %d, maximum_size = %d, connection_type = %d WHERE id = '%s'`, n.nodeTable, node.GetIp(), node.GetPort(), node.GetCurrentSize(), node.GetMaxSize(), node.GetConnectorType(), id.String()))
	return err
}

func (n *NodeStorage) GetAll() ([]guid.UUID, []node.INode, error) {
	row := n.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s`, n.nodeTable))

	rowAmountStr := ""
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	err := row.Scan(&rowAmountStr)
	if err != nil {
		return nil, nil, err
	}

	rows, err := n.db.Query(fmt.Sprintf(`SELECT * FROM %s`, n.nodeTable))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var (
		id             guid.UUID
		ip, port       string
		curSize        int64
		maxSize        int64
		connectionType int
	)
	rowAmount, err := strconv.ParseInt(rowAmountStr, 0, 64)
	if err != nil {
		return nil, nil, err
	}
	nodes := make([]node.INode, rowAmount)
	ids := make([]guid.UUID, rowAmount)
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&id, &ip, &port, &curSize, &maxSize, &connectionType)
		if err != nil {
			return nil, nil, err
		}
		nodes[i] = node.CreateNode(ip, port, maxSize, connector.NetConnectorType(connectionType))
		_ = nodes[i].UpdateCurrentSize(curSize)
		ids[i] = id
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	return ids, nodes, nil
}

func (n *NodeStorage) Get(id guid.UUID) (node.INode, error) {
	row := n.db.QueryRow(fmt.Sprintf(`SELECT ip, port, current_size, maximum_size, connection_type FROM %s WHERE id = '%s'`, n.nodeTable, id.String()))
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		ip, port       string
		curSize        int64
		maxSize        int64
		connectionType int
	)
	err := row.Scan(&ip, &port, &curSize, &maxSize, &connectionType)
	if err != nil {
		return nil, err
	}
	newNode := node.CreateNode(ip, port, maxSize, connector.NetConnectorType(connectionType))
	_ = newNode.UpdateCurrentSize(curSize)
	return newNode, nil
}

func (n *NodeStorage) Close() error {
	return n.db.Close()
}
