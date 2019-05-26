package api

import (
	"context"

	pb "github.com/mingkaic/ultrasound/viewer/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type viewerServer struct{}

// NewViewerServer ...
// Creates viewer grpc server in compliance with protobuf definition
func NewViewerServer() pb.ViewerServer {
	return &viewerServer{}
}

func (*viewerServer) ListGraph(ctx context.Context, req *pb.ListGraphRequest) (*pb.ListGraphResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGraph not implemented")
}

func (*viewerServer) GetGraph(ctx context.Context, req *pb.GetGraphRequest) (*pb.GetGraphResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGraph not implemented")
}

func (*viewerServer) GetNodeData(req *pb.GetNodeDataRequest, srv pb.Viewer_GetNodeDataServer) error {
	return status.Errorf(codes.Unimplemented, "method GetNodeData not implemented")
}

func (*viewerServer) DeleteGraph(ctx context.Context, req *pb.DeleteGraphRequest) (*pb.DeleteGraphResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGraph not implemented")
}
