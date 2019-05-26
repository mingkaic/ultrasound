package data

import (
	"github.com/jinzhu/gorm"
)

type (
	Node struct {
		ID      int32    `gorm:"id"`
		GraphID int32    `gorm:"graph_id"`
		NodeID  int32    `gorm:"node_id"`
		Repr    string   `gorm:"representation"`
		Shape   []uint32 `gorm:"shape"`
	}

	Edge struct {
		GraphID int32  `gorm:"graph_id"`
		Parent  int32  `gorm:"parent_id"`
		Child   int32  `gorm:"child_id"`
		Label   string `gorm:"label"`
	}

	NodeLabel struct {
		GraphID int32  `gorm:"graph_id"`
		NodeID  int32  `gorm:"node_id"`
		Label   string `gorm:"label"`
	}

	GraphData interface {
		ListNodes(params map[string]interface{}) ([]*Node, error)
		CreateNodes(nodes []*Node) error
		CreateEdges(edges []*Edge) error
		LabelNodes(labels []*NodeLabel) error
	}

	graphData struct {
		db *gorm.DB
	}
)

func NewGraphData(db *gorm.DB) GraphData {
	return &graphData{db: db}
}

func (d *graphData) ListNodes(params map[string]interface{}) ([]*Node, error) {
	return nil, nil
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

func (d *graphData) LabelNodes(labels []*NodeLabel) (err error) {
	// todo: use BatchInsert once it's implemented
	for _, label := range labels {
		if err = d.db.Create(label).Error; err != nil {
			return
		}
	}
	return
}
