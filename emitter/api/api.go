package api

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/mingkaic/ultrasound/data"
	pb "github.com/mingkaic/ultrasound/emitter/proto"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type emitterServer struct{}

// NewEmitterServer ...
// Creates emitter grpc server in compliance with protobuf definition
func NewEmitterServer() pb.GraphEmitterServer {
	return &emitterServer{}
}

func convertPbShape(shape []uint32) []int64 {
	slist := make([]int64, len(shape))
	for i, dim := range shape {
		slist[i] = int64(dim)
	}
	return slist
}

func (*emitterServer) HealthCheck(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (*emitterServer) CreateGraph(ctx context.Context, req *pb.CreateGraphRequest) (*pb.CreateGraphResponse, error) {
	graphInfo := req.Payload
	gid := graphInfo.GraphId

	if nil == graphInfo.Nodes {
		return nil, status.Errorf(codes.InvalidArgument, "missing nodes from GraphInfo")
	}
	if nil == graphInfo.Edges {
		return nil, status.Errorf(codes.InvalidArgument, "missing edges from GraphInfo")
	}

	nodes := make([]*data.Node, len(graphInfo.Nodes))
	edges := make([]*data.Edge, len(graphInfo.Edges))
	tags := make([]*data.NodeTag, 0, len(graphInfo.Nodes))
	for i, node := range graphInfo.Nodes {
		var (
			maxheight int
			minheight int
			loc       = node.Location
		)
		if loc != nil {
			maxheight = int(loc.Maxheight)
		}
		if loc != nil {
			minheight = int(loc.Minheight)
		}
		nodes[i] = &data.Node{
			GraphID:   gid,
			NodeID:    int(node.Id),
			Shape:     convertPbShape(node.Shape),
			Maxheight: maxheight,
			Minheight: minheight,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		for key, value := range node.Tags {
			tags = append(tags, &data.NodeTag{
				GraphID: gid,
				NodeID:  int(node.Id),
				TagKey:  key,
				TagVal:  value,
			})
		}
	}

	for i, edge := range graphInfo.Edges {
		edges[i] = &data.Edge{
			GraphID:   gid,
			ParentID:  int(edge.Parent),
			ChildID:   int(edge.Child),
			Label:     edge.Label,
			Shaper:    edge.Shaper,
			Coorder:   edge.Coorder,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	if err := data.Transaction(func(db *sql.Tx) (err error) {
		gData := data.NewGraphData(db)
		if err = gData.CreateNodes(nodes); err != nil {
			return
		}
		if err = gData.UpsertNodeTags(tags); err != nil {
			return
		}
		err = gData.CreateEdges(edges)
		return
	}); err != nil {
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateGraphResponse{
		Status:  pb.Status_OK,
		Message: "Successfully created graph",
	}, nil
}

func (*emitterServer) UpdateNodeData(stream pb.GraphEmitter_UpdateNodeDataServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UpdateNodeDataResponse{
				Status:  pb.Status_OK,
				Message: "Successfully updated node data",
			})
		}
		if err != nil {
			return err
		}

		dataInfo := req.Payload
		datarr := make([]float64, len(dataInfo.Data))
		for i, datum := range dataInfo.Data {
			datarr[i] = float64(datum)
		}
		dentry := &data.NodeData{
			GraphID:   dataInfo.GraphId,
			NodeID:    int(dataInfo.NodeId),
			RawData:   datarr,
			UpdatedAt: time.Now(),
		}

		if err := data.Transaction(func(db *sql.Tx) (err error) {
			gData := data.NewGraphData(db)
			err = gData.UpsertData(dentry)
			return
		}); err != nil {
			log.Error(err.Error())
			return status.Errorf(codes.Internal, err.Error())
		}
	}
}
