package chainedhotstuff

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/relab/hotstuff/consensus"
	"github.com/relab/hotstuff/internal/mocks"
	"github.com/relab/hotstuff/internal/testutil"
	"github.com/relab/hotstuff/synchronizer"
)

// TestPropose checks that a leader broadcasts a new proposal, and then sends a vote to the next leader
/*func TestPropose(t *testing.T) {
	// Setup mocks
	ctrl := gomock.NewController(t)
	hs := New()
	builder := testutil.TestModules(t, ctrl, 1, testutil.GenerateECDSAKey(t))
	synchronizer := synchronizer.New(testutil.FixedTimeout(1000))
	cfg, replicas := testutil.CreateMockConfigurationWithReplicas(t, ctrl, 2)
	builder.Register(hs, cfg, testutil.NewLeaderRotation(t, 1, 2), synchronizer)
	builder.Build()

	// RULES:

	// leader should propose to other replicas.
	cfg.EXPECT().Propose(gomock.AssignableToTypeOf(consensus.ProposeMsg{}))

	// leader should send its own vote to the next leader.
	replicas[1].EXPECT().Vote(gomock.Any())

	hs.Propose(consensus.NewSyncInfo().WithQC(synchronizer.HighQC()))

	if hs.lastVote != 1 {
		t.Errorf("Wrong view: got: %d, want: %d", hs.lastVote, 1)
	}
}*/

// TestCommit checks that a replica commits and executes a valid branch
func TestCommit(t *testing.T) {
	/*
		const n = 4
		ctrl := gomock.NewController(t)
		bl := testutil.CreateBuilders(t, ctrl, n)
		cs := mocks.NewMockConsensus(ctrl)
		bl[0].Register(synchronizer.New(testutil.FixedTimeout(1000)), cs)
		hl := bl.Build()
		hs := hl[0]

		cs.EXPECT().Propose(gomock.AssignableToTypeOf(consensus.NewSyncInfo()))

		ok := false
		ctx, cancel := context.WithCancel(context.Background())
		hs.EventLoop().RegisterObserver(consensus.NewViewMsg{}, func(event interface{}) {
			ok = true
			cancel()
		})

		b := testutil.NewProposeMsg(
			consensus.GetGenesis().Hash(),
			consensus.NewQuorumCert(nil, 1, consensus.GetGenesis().Hash()),
			"test", 1, 1,
		)
		hs.BlockChain().Store(b.Block)

		for i, signer := range hl.Signers() {
			pc, err := signer.CreatePartialCert(b.Block)
			if err != nil {
				t.Fatalf("Failed to create partial certificate: %v", err)
			}
			hs.EventLoop().AddEvent(consensus.VoteMsg{ID: hotstuff.ID(i + 1), PartialCert: pc})
		}

		hs.Run(ctx)

		if !ok {
			t.Error("No new view event happened")
		}
	*/
	const n = 4
	ctrl := gomock.NewController(t)
	acceptor := mocks.NewMockAcceptor(ctrl)
	executor := mocks.NewMockExecutor(ctrl)
	bl := testutil.CreateBuilders(t, ctrl, n)
	cs := mocks.NewMockConsensus(ctrl)
	hl := bl.Build()
	hs := hl[0]
	hotstuff := New()
	cb := consensus.NewConsensusBase(hotstuff)
	cb.InitConsensusModule(hs, nil)
	//signers := hl.Signers()
	bl[0].Register(hs, cs, acceptor, executor, synchronizer.New(testutil.FixedTimeout(1000)))

	ok := false
	ctx, cancel := context.WithCancel(context.Background())
	hs.EventLoop().RegisterObserver(consensus.VoteMsg{}, func(event interface{}) {
		ok = true
		cancel()
	})

	/*const n = 4
	ctrl := gomock.NewController(t)
	hs := New();
	consensusbase := consensus.New(hs);
	keys := testutil.GenerateKeys(t, n, testutil.GenerateECDSAKey)
	bl := testutil.CreateBuilders(t, ctrl, n, keys...)
	acceptor := mocks.NewMockAcceptor(ctrl)
	executor := mocks.NewMockExecutor(ctrl)
	synchronizer := synchronizer.New(testutil.FixedTimeout(1000))
	cfg, replicas := testutil.CreateMockConfigurationWithReplicas(t, ctrl, n, keys...)
	bl[0].Register(hs, cfg, acceptor, executor, synchronizer, leaderrotation.NewFixed(2))
	hl := bl.Build()
	signers := hl.Signers()*/

	// create the needed blocks and QCs
	genesisQC := consensus.NewQuorumCert(nil, 0, consensus.GetGenesis().Hash())
	b1 := testutil.NewProposeMsg(consensus.GetGenesis().Hash(), genesisQC, "1", 1, 2)
	hs.BlockChain().Store(b1.Block)
	/*b1QC := testutil.CreateQC(t, b1.Block, signers)
	b2 := testutil.NewProposeMsg(b1.Block.Hash(), b1QC, "2", 2, 2)
	b2QC := testutil.CreateQC(t, b2.Block, signers)
	b3 := testutil.NewProposeMsg(b2.Block.Hash(), b2QC, "3", 3, 2)
	b3QC := testutil.CreateQC(t, b3.Block, signers)
	b4 := testutil.NewProposeMsg(b3.Block.Hash(), b3QC, "4", 4, 2)*/

	// the second replica will be the leader, so we expect it to receive votes
	/*replicas[1].EXPECT().Vote(gomock.Any()).AnyTimes()
	replicas[1].EXPECT().NewView(gomock.Any()).AnyTimes()
	*/

	// executor will check that the correct command is executed
	/*executor.EXPECT().Exec(gomock.Any()).Do(func(arg interface{}) {
		if arg.(consensus.Command) != b1.Block.Command() {
			t.Errorf("Wrong command executed: got: %s, want: %s", arg, b1.Block.Command())
		}
	})*/

	cb.OnPropose(b1)

	fmt.Print("before")

	hs.Run(ctx)

	if !ok {
		t.Error("No new view event happened")
	}

	/*hs.EventLoop().AddEvent(consensus.ProposeMsg{ID: b1.ID, Block: b1.Block})
	hs.EventLoop().AddEvent(consensus.ProposeMsg{ID: b2.ID, Block: b2.Block})
	hs.EventLoop().AddEvent(consensus.ProposeMsg{ID: b3.ID, Block: b3.Block})
	hs.EventLoop().AddEvent(consensus.ProposeMsg{ID: b4.ID, Block: b4.Block})

	hs.Run(ctx)

	// acceptor expects to receive the commands in order
	gomock.InOrder(
		acceptor.EXPECT().Proposed(gomock.Any()),
		acceptor.EXPECT().Accept(consensus.Command("1")).Return(true),
		acceptor.EXPECT().Proposed(consensus.Command("1")),
		acceptor.EXPECT().Accept(consensus.Command("2")).Return(true),
		acceptor.EXPECT().Proposed(consensus.Command("2")),
		acceptor.EXPECT().Accept(consensus.Command("3")).Return(true),
		acceptor.EXPECT().Proposed(consensus.Command("3")),
		acceptor.EXPECT().Accept(consensus.Command("4")).Return(true),
	)

	if !ok {
		t.Error("No new propose event happened")
	}*/

	/*fmt.Print("a")
	cb.OnPropose(b1)
	fmt.Print("b")
	cb.OnPropose(b2)
	cb.OnPropose(b3)
	cb.OnPropose(b4)*/
}

// TestForkingAttack shows that it is possible to execute a forking attack against HotStuff.
// A forking attack is when a proposal creates a fork in the block chain, leading to some commands never being executed.
// Such as scenario is illustrated in the diagram below.
// Let the arrows from the sides of the blocks represent parent links,
// while the arrows from the corners of the blocks represent QC links:
//          __________________________________
//         /                                  \
//        /                                    +---+
//       /       +-----------------------------| E |
//      / ___    |  ___       ___              +---+
//     / /   \   v /   \     /   \
//  +---+     +---+     +---+     +---+
//  | A |<----| B |<----| C |<----| D |
//  +---+     +---+     +---+     +---+
//
// Here, block E creates a new fork which means that blocks C and D will not be executed.
/*func TestForkingAttack(t *testing.T) {
	const n = 4
	ctrl := gomock.NewController(t)
	hs := New()
	keys := testutil.GenerateKeys(t, n, testutil.GenerateECDSAKey)
	bl := testutil.CreateBuilders(t, ctrl, n, keys...)
	cfg, replicas := testutil.CreateMockConfigurationWithReplicas(t, ctrl, n, keys...)
	executor := mocks.NewMockExecutor(ctrl)
	synchronizer := synchronizer.New(testutil.FixedTimeout(1000))
	bl[0].Register(hs, cfg, executor, synchronizer, leaderrotation.NewFixed(2))
	hl := bl.Build()
	signers := hl.Signers()

	// configure mocks
	replicas[1].EXPECT().Vote(gomock.Any()).AnyTimes()
	replicas[1].EXPECT().NewView(gomock.Any()).AnyTimes()

	genesisQC := consensus.NewQuorumCert(nil, 0, consensus.GetGenesis().Hash())
	a := testutil.NewProposeMsg(consensus.GetGenesis().Hash(), genesisQC, "A", 1, 2)
	aQC := testutil.CreateQC(t, a.Block, signers)
	b := testutil.NewProposeMsg(a.Block.Hash(), aQC, "B", 2, 2)
	bQC := testutil.CreateQC(t, b.Block, signers)
	c := testutil.NewProposeMsg(b.Block.Hash(), bQC, "C", 3, 2)
	cQC := testutil.CreateQC(t, c.Block, signers)
	d := testutil.NewProposeMsg(c.Block.Hash(), cQC, "D", 4, 2)
	e := testutil.NewProposeMsg(b.Block.Hash(), aQC, "E", 5, 2)

	// expected order of execution
	gomock.InOrder(
		executor.EXPECT().Exec(a.Block.Command()),
		executor.EXPECT().Exec(b.Block.Command()),
		executor.EXPECT().Exec(e.Block.Command()),
	)

	hs.OnPropose(a)
	hs.OnPropose(b)
	hs.OnPropose(c)
	hs.OnPropose(d)

	// sanity check
	if hs.bLock != b.Block {
		t.Fatalf("Not locked on B!")
	}

	hs.OnPropose(e)

	// advance views until E is executed
	block := advanceView(t, hs, e.Block, signers)
	block = advanceView(t, hs, block, signers)
	_ = advanceView(t, hs, block, signers)
}

func advanceView(t *testing.T, hs *ChainedHotStuff, lastProposal *consensus.Block, signers []consensus.Crypto) *consensus.Block {
	t.Helper()

	qc := testutil.CreateQC(t, lastProposal, signers)
	b := consensus.NewBlock(lastProposal.Hash(), qc, "foo", hs.lastVote+1, 2)
	hs.OnPropose(consensus.ProposeMsg{ID: b.Proposer(), Block: b})
	return b
}

// TestChainedHotstuff runs chained hotstuff with the gorums backend and expects each replica to execute 10 times.
func TestChainedHotstuff(t *testing.T) {
	const n = 4
	ctrl := gomock.NewController(t)

	baseCfg := config.NewConfig(0, nil, nil)

	listeners := make([]net.Listener, n)
	keys := make([]consensus.PrivateKey, n)
	for i := 0; i < n; i++ {
		listeners[i] = testutil.CreateTCPListener(t)
		key := testutil.GenerateECDSAKey(t)
		keys[i] = key
		id := consensus.ID(i + 1)
		baseCfg.Replicas[id] = &config.ReplicaInfo{
			ID:      id,
			Address: listeners[i].Addr().String(),
			PubKey:  key.Public(),
		}
	}

	builders := testutil.CreateBuilders(t, ctrl, n, keys...)
	configs := make([]*backend.Config, n)
	servers := make([]*backend.Server, n)
	synchronizers := make([]consensus.Synchronizer, n)
	for i := 0; i < n; i++ {
		c := *baseCfg
		c.ID = consensus.ID(i + 1)
		c.PrivateKey = keys[i].(*ecdsa.PrivateKey)
		configs[i] = backend.NewConfig(c.ID, nil, gorums.WithDialTimeout(time.Second))
		servers[i] = backend.NewServer()
		synchronizers[i] = synchronizer.New(
			synchronizer.NewViewDuration(1000, 100, 2),
		)
		builders[i].Register(New(), configs[i], servers[i], synchronizers[i])
	}

	executors := make([]*mocks.MockExecutor, n)
	counters := make([]uint, n)
	c := make(chan struct{}, n)
	errChan := make(chan error, n)
	for i := 0; i < n; i++ {
		counter := &counters[i]
		executors[i] = mocks.NewMockExecutor(ctrl)
		executors[i].EXPECT().Exec(gomock.Any()).AnyTimes().Do(func(arg consensus.Command) {
			if arg != consensus.Command("foo") {
				errChan <- fmt.Errorf("unknown command executed: got %s, want: %s", arg, "foo")
			}
			*counter++
			if *counter >= 100 {
				c <- struct{}{}
			}
		})
		builders[i].Register(executors[i])
	}

	hl := builders.Build()

	ctx, cancel := context.WithCancel(context.Background())

	for i, server := range servers {
		server.StartOnListener(listeners[i])
		defer server.Stop()
	}

	for _, cfg := range configs {
		err := cfg.Connect(baseCfg)
		if err != nil {
			t.Fatal(err)
		}
		defer cfg.Close()
	}

	for _, hs := range hl {
		go func(hs *consensus.Modules) {
			hs.Synchronizer().Start(ctx)
			hs.EventLoop().Run(ctx)
		}(hs)
	}

	for i := 0; i < n; i++ {
		select {
		case <-c:
		case err := <-errChan:
			t.Fatal(err)
		}
	}
	cancel()
}*/
