package grpcserver

import (
	"context"
	"net"
	"time"

	pb "orchestrator/proto"
	"orchestrator/registry"
	"orchestrator/router"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Server struct {
	pb.UnimplementedOrchestratorServiceServer
	reg       *registry.AgentRegistry
	msgRouter *router.MessageRouter
}

func NewServer(reg *registry.AgentRegistry, msgRouter *router.MessageRouter) *Server {
	return &Server{reg: reg, msgRouter: msgRouter}
}

func (s *Server) Health(ctx context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Service: "orchestrator", Status: "ok"}, nil
}

func (s *Server) ListAgents(ctx context.Context, _ *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	agents := s.reg.List()
	out := make([]*pb.Agent, 0, len(agents))
	for _, a := range agents {
		out = append(out, &pb.Agent{
			Name:        a.Name,
			BaseUrl:     a.BaseURL,
			Healthy:     a.Health,
			LastChecked: a.LastChecked.Format(time.RFC3339),
			LastError:   a.LastError,
		})
	}
	return &pb.ListAgentsResponse{Agents: out}, nil
}

func (s *Server) RouteMessage(ctx context.Context, req *pb.A2AMessage) (*pb.RouteResponse, error) {
	payload := map[string]interface{}{}
	for k, v := range req.Payload {
		payload[k] = v.AsInterface()
	}

	msg := router.A2AMessage{
		From:      req.From,
		To:        req.To,
		EventType: req.EventType,
		Payload:   payload,
	}

	if err := s.msgRouter.Forward(ctx, msg); err != nil {
		return nil, status.Error(14, err.Error()) // Unavailable
	}

	return &pb.RouteResponse{
		Status: "routed",
		To:     req.To,
	}, nil
}

func Start(port string, reg *registry.AgentRegistry, msgRouter *router.MessageRouter) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	g := grpc.NewServer()
	pb.RegisterOrchestratorServiceServer(g, NewServer(reg, msgRouter))
	return g.Serve(lis)
}

// helper if needed later
func toStructValue(v interface{}) *structpb.Value {
	sv, _ := structpb.NewValue(v)
	return sv
}
