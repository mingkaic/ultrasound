package api

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/mingkaic/ultrasound/data"
	pb "github.com/mingkaic/ultrasound/emitter/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type emitterServer struct{}

// NewEmitterServer ...
// Creates emitter grpc server in compliance with protobuf definition
func NewEmitterServer() pb.GraphEmitterServer {
	return &emitterServer{}
}

func (*emitterServer) CreateGraph(ctx context.Context, req *pb.CreateGraphRequest) (*pb.CreateGraphResponse, error) {
	graphInfo := req.Payload
	id := graphInfo.GraphId

	if nil == graphInfo.Nodes {
		return nil, status.Errorf(codes.InvalidArgument, "missing nodes from GraphInfo")
	}
	if nil == graphInfo.Edges {
		return nil, status.Errorf(codes.InvalidArgument, "missing edges from GraphInfo")
	}

	nodes := make([]*data.Node, len(graphInfo.Nodes))
	edges := make([]*data.Edge, len(graphInfo.Edges))
	labels := make([]*data.NodeLabel, 0)
	for i, node := range graphInfo.Nodes {
		nodes[i] = &data.Node{
			GraphID: id,
			NodeID:  node.Id,
			Repr:    node.Repr,
			Shape:   node.Shape,
		}
		for _, label := range node.Labels {
			labels = append(labels, &data.NodeLabel{
				GraphID: id,
				NodeID:  node.Id,
				Label:   label,
			})
		}
	}

	for i, edge := range graphInfo.Edges {
		edges[i] = &data.Edge{
			GraphID: id,
			Parent:  edge.Parent,
			Child:   edge.Child,
			Label:   edge.Label,
		}
	}

	if err := data.Transaction(func(db *gorm.DB) (err error) {
		gData := data.NewGraphData(db)
		if err = gData.CreateNodes(nodes); err != nil {
			return
		}
		if err = gData.LabelNodes(labels); err != nil {
			return
		}
		if err = gData.CreateEdges(edges); err != nil {
			return
		}
		return
	}); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateGraphResponse{
		Status:  pb.Status_OK,
		Message: "Successfully created graph",
	}, nil
}

func (*emitterServer) UpdateNodeMeta(srv pb.GraphEmitter_UpdateNodeMetaServer) error {
	return status.Errorf(codes.Unimplemented, "method UpdateNodeMeta not implemented")
}

func (*emitterServer) UpdateNodeData(srv pb.GraphEmitter_UpdateNodeDataServer) error {
	return status.Errorf(codes.Unimplemented, "method UpdateNodeData not implemented")
}
