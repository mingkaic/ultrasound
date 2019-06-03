package data

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type (
	Node struct {
		GraphID   string            `gorm:"graph_id"`
		NodeID    int               `gorm:"node_id"`
		Shape     string            `gorm:"shape"`
		Maxheight int               `gorm:"maxheight"`
		Minheight int               `gorm:"minheight"`
		Tags      map[string]string `sql:"-"`
	}

	Edge struct {
		GraphID  string `gorm:"graph_id"`
		ParentID int    `gorm:"parent_id"`
		ChildID  int    `gorm:"child_id"`
		Label    string `gorm:"label"`
		Shaper   string `gorm:"shaper"`
		Coorder  string `gorm:"coorder"`
	}

	NodeTag struct {
		GraphID string `gorm:"graph_id"`
		NodeID  int    `gorm:"node_id"`
		TagKey  string `gorm:"tag_key"`
		TagVal  string `gorm:"tag_val"`
	}

	NodeData struct {
		GraphID string    `gorm:"graph_id"`
		NodeID  int       `gorm:"node_id"`
		RawData []float64 `sql:"-"`
	}

	GraphData interface {
		ListGraphs() ([]string, error)
		ListNodes(params map[string]interface{}) ([]*Node, error)
		ListEdges(params map[string]interface{}) ([]*Edge, error)
		GetNodeData(graphID string, nodeID int) (*NodeData, error)
		TagNode(node *Node) (*Node, error)

		CreateNodes(nodes []*Node) error
		CreateEdges(edges []*Edge) error
		UpdateData(data *NodeData) error
		TagNodes(tags []*NodeTag) error
	}

	graphData struct {
		db *gorm.DB
	}
)

func NewGraphData(db *gorm.DB) GraphData {
	return &graphData{db: db}
}

func (d *graphData) ListGraphs() ([]string, error) {
	rows, err := d.db.Raw(`
		SELECT distinct graph_id
		FROM nodes
		`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gids := []string{}
	for rows.Next() {
		var gid string
		if rows.Scan(&gid) != nil {
			break
		}
		gids = append(gids, gid)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return gids, nil
}

func (d *graphData) ListNodes(params map[string]interface{}) (out []*Node, err error) {
	err = d.db.Find(&out, params).Error
	return
}

func (d *graphData) ListEdges(params map[string]interface{}) (out []*Edge, err error) {
	err = d.db.Find(&out, params).Error
	return
}

func (d *graphData) GetNodeData(graphID string, nodeID int) (out *NodeData, err error) {
	holder := struct{ Data []float64 }{}
	err = d.db.Raw(`
		SELECT data
		FROM node_data
		WHERE graph_id = ? and node_id = ?
		`, graphID, nodeID).Scan(&holder).Error
	if err != nil {
		return
	}
	out = &NodeData{
		GraphID: graphID,
		NodeID:  nodeID,
		RawData: holder.Data,
	}
	return
}

func (d *graphData) TagNode(node *Node) (*Node, error) {
	var tags []*NodeTag
	if err := d.db.Where(map[string]interface{}{
		"graph_id": node.GraphID,
		"node_id":  node.NodeID,
	}).Find(&tags).Error; err != nil {
		return nil, err
	}

	node.Tags = make(map[string]string)
	for _, tag := range tags {
		node.Tags[tag.TagKey] = tag.TagVal
	}

	return node, nil
}

func (d *graphData) CreateNodes(nodes []*Node) (err error) {
	// todo: use BatchInsert once it's implemented
	for _, node := range nodes {
		if err = d.db.Create(node).Error; err != nil {
			return
		}
	}
	return
}

func (d *graphData) CreateEdges(edges []*Edge) (err error) {
	// todo: use BatchInsert once it's implemented
	for _, edge := range edges {
		if err = d.db.Create(edge).Error; err != nil {
			return
		}
	}
	return
}

func (d *graphData) UpdateData(dentry *NodeData) (err error) {
	if err = d.db.FirstOrCreate(&NodeData{},
		NodeData{
			GraphID: dentry.GraphID,
			NodeID:  dentry.NodeID,
		}).Error; err != nil {
		return
	}
	updateStmt := `
	UPDATE ONLY node_data SET data = $1
	WHERE graph_id = $2 AND node_id = $3
	`
	err = d.db.Exec(updateStmt, pq.Array(dentry.RawData), dentry.GraphID, dentry.NodeID).Error
	return
}

func (d *graphData) TagNodes(tags []*NodeTag) (err error) {
	// todo: use BatchInsert once it's implemented
	for _, tag := range tags {
		if err = d.db.Create(tag).Error; err != nil {
			return
		}
	}
	return
}
