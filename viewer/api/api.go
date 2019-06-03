package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mingkaic/ultrasound/data"
	pb "github.com/mingkaic/ultrasound/viewer/proto"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type viewerServer struct{}

// NewViewerServer ...
// Creates viewer grpc server in compliance with protobuf definition
func NewViewerServer() pb.ViewerServer {
	return &viewerServer{}
}

func convertDbShape(shape string) []uint32 {
	slist := strings.Split(shape, ",")
	out := make([]uint32, len(slist))
	for i, dim := range slist {
		val, err := strconv.ParseUint(dim, 10, 32)
		if err != nil {
			log.Errorf("unknown dimension %s: %+v", dim, err)
		}
		out[i] = uint32(val)
	}
	return out
}

func (*viewerServer) ListGraphs(ctx context.Context, req *pb.ListGraphRequest) (*pb.ListGraphResponse, error) {
	var gids []string
	if err := data.Transaction(func(db *gorm.DB) (err error) {
		gData := data.NewGraphData(db)
		gids, err = gData.ListGraphs()
		return
	}); err != nil {
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ListGraphResponse{
		Result:  gids,
		Status:  pb.Status_OK,
		Message: fmt.Sprintf("Got %d graphs", len(gids)),
	}, nil
}

func (*viewerServer) GetGraph(ctx context.Context, req *pb.GetGraphRequest) (*pb.GetGraphResponse, error) {
	gid := req.GraphId
	var (
		nodes []*data.Node
		edges []*data.Edge
		args  = map[string]interface{}{
			"graph_id": gid,
		}
	)
	if err := data.Transaction(func(db *gorm.DB) (err error) {
		gData := data.NewGraphData(db)
		nodes, err = gData.ListNodes(args)
		if err != nil {
			return
		}
		edges, err = gData.ListEdges(args)
		if err != nil {
			return
		}
		for i, node := range nodes {
			node, err = gData.TagNode(node)
			if err != nil {
				return
			}
			nodes[i] = node
		}
		return
	}); err != nil {
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbNodes := make([]*pb.NodeInfo, len(nodes))
	pbEdges := make([]*pb.EdgeInfo, len(edges))
	for i, node := range nodes {
		pbNodes[i] = &pb.NodeInfo{
			Id:    int32(node.NodeID),
			Shape: convertDbShape(node.Shape),
			Tags:  node.Tags,
			Location: &pb.NodeLoc{
				Maxheight: uint32(node.Maxheight),
				Minheight: uint32(node.Minheight),
			},
		}
	}
	for i, edge := range edges {
		pbEdges[i] = &pb.EdgeInfo{
			Parent:  int32(edge.ParentID),
			Child:   int32(edge.ChildID),
			Label:   edge.Label,
			Shaper:  edge.Shaper,
			Coorder: edge.Coorder,
		}
	}
	return &pb.GetGraphResponse{
		Result: &pb.GraphInfo{
			GraphId: gid,
			Nodes:   pbNodes,
			Edges:   pbEdges,
		},
		Status:  pb.Status_OK,
		Message: fmt.Sprintf("Got %s graph", gid),
	}, nil
}

func (*viewerServer) GetNodeData(ctx context.Context, req *pb.GetNodeDataRequest) (*pb.GetNodeDataResponse, error) {
	var node *data.NodeData
	gid := req.GraphId
	nid := req.NodeId

	if err := data.Transaction(func(db *gorm.DB) (err error) {
		gData := data.NewGraphData(db)
		node, err = gData.GetNodeData(gid, int(nid))
		return
	}); err != nil {
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	rdata := make([]float32, len(node.RawData))
	for i, rdatum := range node.RawData {
		rdata[i] = float32(rdatum)
	}
	return &pb.GetNodeDataResponse{
		Result: &pb.NodeData{
			GraphId: gid,
			NodeId:  nid,
			Data:    rdata,
		},
		Status:  pb.Status_OK,
		Message: fmt.Sprintf("Got node %d of graph %s", nid, gid),
	}, nil
}

func (*viewerServer) DeleteGraph(ctx context.Context, req *pb.DeleteGraphRequest) (*pb.DeleteGraphResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGraph not implemented")
}
