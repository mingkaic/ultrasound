package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type (
	Node struct {
		GraphID   string            `gorm:"graph_id"`
		NodeID    int               `gorm:"node_id"`
		Shape     []int64           `gorm:"shape"`
		Maxheight int               `gorm:"maxheight"`
		Minheight int               `gorm:"minheight"`
		Tags      map[string]string `sql:"-"`
		CreatedAt time.Time         `gorm:"created_at"`
		UpdatedAt time.Time         `gorm:"updated_at"`
	}

	Edge struct {
		GraphID   string    `gorm:"graph_id"`
		ParentID  int       `gorm:"parent_id"`
		ChildID   int       `gorm:"child_id"`
		Label     string    `gorm:"label"`
		Shaper    string    `gorm:"shaper"`
		Coorder   string    `gorm:"coorder"`
		CreatedAt time.Time `gorm:"created_at"`
		UpdatedAt time.Time `gorm:"updated_at"`
	}

	NodeTag struct {
		GraphID string `gorm:"graph_id"`
		NodeID  int    `gorm:"node_id"`
		TagKey  string `gorm:"tag_key"`
		TagVal  string `gorm:"tag_val"`
	}

	NodeData struct {
		GraphID   string    `gorm:"graph_id"`
		NodeID    int       `gorm:"node_id"`
		RawData   []float64 `gorm:"data"`
		UpdatedAt time.Time `gorm:"updated_at"`
	}

	Graph struct {
		GraphID   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	GraphData interface {
		ListGraphs() ([]*Graph, error)
		ListNodes(params map[string]interface{}) ([]*Node, error)
		ListEdges(params map[string]interface{}) ([]*Edge, error)
		GetNodeData(graphID string, nodeID int) (*NodeData, error)
		TagNode(node *Node) (*Node, error)

		CreateNodes(nodes []*Node) error
		CreateEdges(edges []*Edge) error
		UpsertData(data *NodeData) error
		UpsertNodeTags(tags []*NodeTag) error

		DeleteNodes(graphID string) error
	}

	graphData struct {
		db *sql.Tx
	}
)

const batch_limit = 1000

func NewGraphData(db *sql.Tx) GraphData {
	return &graphData{db: db}
}

func (d *graphData) ListGraphs() ([]*Graph, error) {
	rows, err := (&queryStmt{from: "nodes"}).query(d.db, "distinct graph_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		gid  string
		gids = []string{}
	)
	for rows.Next() {
		if err := rows.Scan(&gid); err != nil {
			return nil, err
		}
		gids = append(gids, gid)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	stmt, args := (&queryStmt{from: "nodes"}).
		where(map[string]interface{}{
			"graph_id": gids,
		}, "and").
		generate("graph_id, max(created_at), max(updated_at)")
	stmt += " group by graph_id"
	sqlLog(stmt, args)
	rowsDate, err := d.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rowsDate.Close()
	graphs := make([]*Graph, 0, len(gids))
	for rowsDate.Next() {
		entry := Graph{}
		if err := rowsDate.Scan(
			&entry.GraphID,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		graphs = append(graphs, &entry)
	}
	if err = rowsDate.Err(); err != nil {
		return nil, err
	}

	return graphs, nil
}

func (d *graphData) ListNodes(params map[string]interface{}) ([]*Node, error) {
	out := make([]*Node, 0)
	rows, err := (&queryStmt{from: "nodes"}).
		where(params, "and").
		query(d.db, "*")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entry := Node{}
		if err := rows.Scan(
			&entry.GraphID,
			&entry.NodeID,
			pq.Array(&entry.Shape),
			&entry.Maxheight,
			&entry.Minheight,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, &entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (d *graphData) ListEdges(params map[string]interface{}) ([]*Edge, error) {
	out := make([]*Edge, 0)
	rows, err := (&queryStmt{from: "edges"}).
		where(params, "and").
		query(d.db, "*")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entry := Edge{}
		if err := rows.Scan(
			&entry.GraphID,
			&entry.ParentID,
			&entry.ChildID,
			&entry.Label,
			&entry.Shaper,
			&entry.Coorder,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, &entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (d *graphData) GetNodeData(graphID string, nodeID int) (*NodeData, error) {
	out := NodeData{
		GraphID: graphID,
		NodeID:  nodeID,
	}
	stmt, args := (&queryStmt{from: "node_data"}).
		where(map[string]interface{}{
			"graph_id": graphID,
			"node_id":  nodeID,
		}, "and").
		generate("data")
	sqlLog(stmt, args)
	row := d.db.QueryRow(stmt, args...)
	if err := row.Scan(pq.Array(&out.RawData)); err != nil {
		return nil, err
	}
	return &out, nil
}

func (d *graphData) TagNode(node *Node) (*Node, error) {
	var (
		key  string
		val  string
		tags = make(map[string]string)
	)
	rows, err := (&queryStmt{from: "node_tags"}).
		where(map[string]interface{}{
			"graph_id": node.GraphID,
			"node_id":  node.NodeID,
		}, "and").
		query(d.db, "tag_key, tag_val")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&key, &val); err != nil {
			return nil, err
		}
		tags[key] = val
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	node.Tags = tags

	return node, nil
}

func splitBatches(arr []interface{}) [][]interface{} {
	nbatches := 1 + len(arr)/batch_limit
	finalBatchsize := len(arr) % batch_limit
	batches := make([][]interface{}, nbatches)
	for i := range batches[:len(batches)-1] {
		batches[i] = arr[i*batch_limit : (i+1)*batch_limit]
	}
	batches[len(batches)-1] = arr[len(arr)-finalBatchsize:]
	return batches
}

func (d *graphData) CreateNodes(nodes []*Node) error {
	log.Errorf("creating nodes %d...%d", nodes[0].NodeID, nodes[len(nodes)-1].NodeID)
	nodeArgs := make([]interface{}, len(nodes))
	for i, node := range nodes {
		nodeArgs[i] = node
	}
	batches := splitBatches(nodeArgs)
	for i, batch := range batches {
		log.Infof("Processing %d nodes batch[%d/%d]", len(batch), i+1, len(batches))
		results, err := (&createStmt{
			into: "nodes",
			fields: []string{
				"graph_id",
				"node_id",
				"shape",
				"maxheight",
				"minheight",
				"created_at",
				"updated_at",
			},
		}).modify(d.db, batch)
		if err != nil {
			return err
		}
		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return err
		}
		log.Infof("Rows affected: %d", rowsAffected)
	}
	return nil
}

func (d *graphData) CreateEdges(edges []*Edge) (err error) {
	edgeArgs := make([]interface{}, len(edges))
	for i, edge := range edges {
		edgeArgs[i] = edge
	}
	batches := splitBatches(edgeArgs)
	for i, batch := range batches {
		log.Infof("Processing %d edges batch[%d/%d]", len(batch), i+1, len(batches)+1)
		results, err := (&createStmt{
			into: "edges",
			fields: []string{
				"graph_id",
				"parent_id",
				"child_id",
				"label",
				"shaper",
				"coorder",
				"created_at",
				"updated_at",
			},
		}).modify(d.db, batch)
		if err != nil {
			return err
		}
		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return err
		}
		log.Infof("Rows affected: %d", rowsAffected)
	}
	return nil
}

func (d *graphData) UpsertData(entry *NodeData) (err error) {
	results, err := (&upsertStmt{
		into: "node_data",
		keyFields: []string{
			"graph_id", "node_id",
		},
		updateFields: []string{
			"data", "updated_at",
		},
	}).modify(d.db, entry)
	if err != nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected: %d", rowsAffected)
	return nil
}

func (d *graphData) UpsertNodeTags(tags []*NodeTag) (err error) {
	tagArgs := make([]interface{}, len(tags))
	for i, tag := range tags {
		tagArgs[i] = tag
	}
	batches := splitBatches(tagArgs)
	for i, batch := range batches {
		log.Infof("Processing %d tags batch[%d/%d]", len(batch), i+1, len(batches)+1)
		results, err := (&upsertStmt{
			into: "node_tags",
			keyFields: []string{
				"graph_id", "node_id", "tag_key",
			},
			updateFields: []string{
				"tag_val",
			},
		}).modify(d.db, batch)
		if err != nil {
			return err
		}
		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return err
		}
		log.Infof("Rows affected: %d", rowsAffected)
	}
	return nil
}

func (d *graphData) DeleteData(params map[string]interface{}) error {
	results, err := d.db.Exec("delete from node_data where graph_id=($1)", params["graph_id"])
	if err != nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected: %d", rowsAffected)
	return nil
}

func (d *graphData) DeleteTags(params map[string]interface{}) error {
	results, err := d.db.Exec("delete from node_tags where graph_id=($1)", params["graph_id"])
	if err != nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected: %d", rowsAffected)
	return nil
}

func (d *graphData) DeleteEdges(params map[string]interface{}) error {
	results, err := d.db.Exec("delete from edges where graph_id=($1)", params["graph_id"])
	if err != nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected: %d", rowsAffected)
	return nil
}

func (d *graphData) DeleteNodes(graphID string) error {
	if err := d.DeleteData(map[string]interface{}{
		"graph_id": graphID,
	}); err != nil {
		return err
	}
	if err := d.DeleteTags(map[string]interface{}{
		"graph_id": graphID,
	}); err != nil {
		return err
	}
	if err := d.DeleteEdges(map[string]interface{}{
		"graph_id": graphID,
	}); err != nil {
		return err
	}
	// todo: support delete stmt in gorm.go
	results, err := d.db.Exec("delete from nodes where graph_id=($1)", graphID)
	if err != nil {
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected: %d", rowsAffected)
	return nil
}
