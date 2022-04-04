// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.5.0-devel
// 	protoc            v3.19.1
// source: internal/proto/hotstuffpb/hotstuff.proto

package hotstuffpb

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(5 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 5)
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	gorums.Configuration
	qspec QuorumSpec
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
func (c *Configuration) Nodes() []*Node {
	nodes := make([]*Node, 0, c.Size())
	for _, n := range c.Configuration {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// And returns a NodeListOption that can be used to create a new configuration combining c and d.
func (c Configuration) And(d *Configuration) gorums.NodeListOption {
	return c.Configuration.And(d.Configuration)
}

// Except returns a NodeListOption that can be used to create a new configuration
// from c without the nodes in rm.
func (c Configuration) Except(rm *Configuration) gorums.NodeListOption {
	return c.Configuration.Except(rm.Configuration)
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

// Manager maintains a connection pool of nodes on
// which quorum calls can be performed.
type Manager struct {
	*gorums.Manager
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) (mgr *Manager) {
	mgr = &Manager{}
	mgr.Manager = gorums.NewManager(opts...)
	return mgr
}

// NewConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewConfiguration(opts ...gorums.ConfigOption) (c *Configuration, err error) {
	if len(opts) < 1 || len(opts) > 2 {
		return nil, fmt.Errorf("wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.Configuration, err = gorums.NewConfiguration(m.Manager, v)
			if err != nil {
				return nil, err
			}
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("unknown option type: %v", v)
		}
	}
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
		return nil, fmt.Errorf("missing required QuorumSpec")
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.Manager.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

type Node struct {
	*gorums.Node
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Propose is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Propose(ctx context.Context, in *Proposal, opts ...gorums.CallOption) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "hotstuffpb.Hotstuff.Propose",
	}

	c.Configuration.Multicast(ctx, cd, opts...)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Timeout is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Timeout(ctx context.Context, in *TimeoutMsg, opts ...gorums.CallOption) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "hotstuffpb.Hotstuff.Timeout",
	}

	c.Configuration.Multicast(ctx, cd, opts...)
}

// QuorumSpec is the interface of quorum functions for Hotstuff.
type QuorumSpec interface {
	gorums.ConfigOption

	// FetchQF is the quorum function for the Fetch
	// quorum call method. The in parameter is the request object
	// supplied to the Fetch method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *BlockHash'.
	FetchQF(in *BlockHash, replies map[uint32]*Block) (*Block, bool)
}

// Fetch is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Fetch(ctx context.Context, in *BlockHash) (resp *Block, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "hotstuffpb.Hotstuff.Fetch",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Block, len(replies))
		for k, v := range replies {
			r[k] = v.(*Block)
		}
		return c.qspec.FetchQF(req.(*BlockHash), r)
	}

	res, err := c.Configuration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Block), err
}

// Hotstuff is the server-side API for the Hotstuff Service
type Hotstuff interface {
	Propose(ctx gorums.ServerCtx, request *Proposal)
	Vote(ctx gorums.ServerCtx, request *PartialCert)
	Timeout(ctx gorums.ServerCtx, request *TimeoutMsg)
	NewView(ctx gorums.ServerCtx, request *SyncInfo)
	Fetch(ctx gorums.ServerCtx, request *BlockHash) (response *Block, err error)
}

func RegisterHotstuffServer(srv *gorums.Server, impl Hotstuff) {
	srv.RegisterHandler("hotstuffpb.Hotstuff.Propose", func(ctx gorums.ServerCtx, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*Proposal)
		defer ctx.Release()
		impl.Propose(ctx, req)
	})
	srv.RegisterHandler("hotstuffpb.Hotstuff.Vote", func(ctx gorums.ServerCtx, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*PartialCert)
		defer ctx.Release()
		impl.Vote(ctx, req)
	})
	srv.RegisterHandler("hotstuffpb.Hotstuff.Timeout", func(ctx gorums.ServerCtx, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*TimeoutMsg)
		defer ctx.Release()
		impl.Timeout(ctx, req)
	})
	srv.RegisterHandler("hotstuffpb.Hotstuff.NewView", func(ctx gorums.ServerCtx, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*SyncInfo)
		defer ctx.Release()
		impl.NewView(ctx, req)
	})
	srv.RegisterHandler("hotstuffpb.Hotstuff.Fetch", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*BlockHash)
		defer ctx.Release()
		resp, err := impl.Fetch(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
}

type internalBlock struct {
	nid   uint32
	reply *Block
	err   error
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Vote is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Vote(ctx context.Context, in *PartialCert, opts ...gorums.CallOption) {
	cd := gorums.CallData{
		Message: in,
		Method:  "hotstuffpb.Hotstuff.Vote",
	}

	n.Node.Unicast(ctx, cd, opts...)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// NewView is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) NewView(ctx context.Context, in *SyncInfo, opts ...gorums.CallOption) {
	cd := gorums.CallData{
		Message: in,
		Method:  "hotstuffpb.Hotstuff.NewView",
	}

	n.Node.Unicast(ctx, cd, opts...)
}
