package grpc

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/you/workflow-runner/internal/pb"
	"go.temporal.io/sdk/client"
)

type GRPCServer struct {
	pb.UnimplementedWorkflowServiceServer
	TemporalClient client.Client
	DB             *pg.DB
	TaskQueue      string
}

func NewGRPCServer(tc client.Client, db *pg.DB, tq string) *GRPCServer {
	return &GRPCServer{TemporalClient: tc, DB: db, TaskQueue: tq}
}

func (s *GRPCServer) StartWorkflow(ctx context.Context, req *pb.StartReq) (*pb.StartResp, error) {
	wid := req.WorkflowId
	// Load nodes from DB via go-pg
	var nodes []Node
	if err := s.DB.Model(&nodes).Where("workflow_id = ?", wid).Order("id ASC").Select(); err != nil {
		return nil, err
	}
	// convert to NodeSpec
	var specs []map[string]interface{}
	for _, n := range nodes {
		specs = append(specs, map[string]interface{}{"id": n.ID.String(), "type": n.Type, "script": n.Script})
	}

	runID := uuid.New().String()
	// Start the Temporal workflow: pass nodes as input (non-deterministic data), runID separate param
	we, err := s.TemporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        runID,
		TaskQueue: s.TaskQueue,
	}, "OrchestrateWorkflow", runID, specs, s.TaskQueue)
	if err != nil {
		return nil, err
	}
	return &pb.StartResp{WorkflowId: we.GetID(), RunId: we.GetRunID()}, nil
}

func (s *GRPCServer) GetStatus(ctx context.Context, req *pb.StatusReq) (*pb.StatusResp, error) {
	// You can query Temporal or your DB; simplified: check workflow run existence
	// Implement as needed
	return &pb.StatusResp{WorkflowId: req.WorkflowId, Status: "running"}, nil
}
