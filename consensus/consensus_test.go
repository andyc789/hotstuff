package consensus_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/relab/hotstuff"
	"github.com/relab/hotstuff/consensus"
	"github.com/relab/hotstuff/internal/mocks"
	"github.com/relab/hotstuff/internal/testutil"
	"github.com/relab/hotstuff/synchronizer"
)

// TestVote checks that a leader can collect votes on a proposal to form a QC
func TestVote(t *testing.T) {
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
}

/*
func TestTwoPhase(t *testing.T) {
	const n = 4
	ctrl := gomock.NewController(t)
	bl := testutil.CreateBuilders(t, ctrl, n)
	cs := mocks.NewMockConsensus(ctrl)
	bl[0].Register(synchronizer.New(testutil.FixedTimeout(1000)), cs)
	hl := bl.Build()
	hs := hl[0]

	cs.EXPECT().Propose(gomock.AssignableToTypeOf(consensus.NewSyncInfo()))

	b1 := testutil.NewProposeMsg(
		consensus.GetGenesis().Hash(),
		consensus.NewQuorumCert(nil, 1, consensus.GetGenesis().Hash()),
		"test", 1, 1,
	)

	hs.BlockChain().Store(b1.Block)

	for i, signer := range hl.Signers() {
		pc, err := signer.CreatePartialCert(b1.Block)
		if err != nil {
			t.Fatalf("Failed to create partial certificate: %v", err)
		}
		consensus.VoteMsg{ID: hotstuff.ID(i + 1), PartialCert: pc} //newVoteMsg function?
	}

	b2 := testutil.NewProposeMsg(
		b1.Block.Hash(),
		consensus.NewQuorumCert(nil, 2, b1.Block.Hash()),
		"test", 2, 1,
	)

	hs.BlockChain().Store(b2.Block)

	for i, signer := range hl.Signers() {
		pc, err := signer.CreatePartialCert(b2.Block)
		if err != nil {
			t.Fatalf("Failed to create partial certificate: %v", err)
		}
		consensus.VoteMsg{ID: hotstuff.ID(i + 1), PartialCert: pc} //newVoteMsg function?
	}

	consensus.Consensus.Propose()

	//IMPLEMENT ALTERNATE WAY USING EVENTLOOP AND CATCHING BLOCK2 WITH A HANDLER?

}*/
/*
func TestPropose(t *testing.T) {
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
}

// TestCommit checks that a replica commits and executes a valid branch
/*func TestCommit(t *testing.T) {
	const n = 4
	ctrl := gomock.NewController(t)
	hs := New()
	keys := testutil.GenerateKeys(t, n, testutil.GenerateECDSAKey)
	bl := testutil.CreateBuilders(t, ctrl, n, keys...)
	acceptor := mocks.NewMockAcceptor(ctrl)
	executor := mocks.NewMockExecutor(ctrl)
	synchronizer := synchronizer.New(testutil.FixedTimeout(1000))
	cfg, replicas := testutil.CreateMockConfigurationWithReplicas(t, ctrl, n, keys...)
	bl[0].Register(hs, cfg, acceptor, executor, synchronizer, leaderrotation.NewFixed(2))
	hl := bl.Build()
	signers := hl.Signers()

	// create the needed blocks and QCs
	genesisQC := consensus.NewQuorumCert(nil, 0, consensus.GetGenesis().Hash())
	b1 := testutil.NewProposeMsg(consensus.GetGenesis().Hash(), genesisQC, "1", 1, 2)
	b1QC := testutil.CreateQC(t, b1.Block, signers)
	b2 := testutil.NewProposeMsg(b1.Block.Hash(), b1QC, "2", 2, 2)
	b2QC := testutil.CreateQC(t, b2.Block, signers)
	b3 := testutil.NewProposeMsg(b2.Block.Hash(), b2QC, "3", 3, 2)
	b3QC := testutil.CreateQC(t, b3.Block, signers)
	b4 := testutil.NewProposeMsg(b3.Block.Hash(), b3QC, "4", 4, 2)

	// the second replica will be the leader, so we expect it to receive votes
	replicas[1].EXPECT().Vote(gomock.Any()).AnyTimes()
	replicas[1].EXPECT().NewView(gomock.Any()).AnyTimes()

	// executor will check that the correct command is executed
	executor.EXPECT().Exec(gomock.Any()).Do(func(arg interface{}) {
		if arg.(consensus.Command) != b1.Block.Command() {
			t.Errorf("Wrong command executed: got: %s, want: %s", arg, b1.Block.Command())
		}
	})

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

	hs.OnPropose(b1)
	hs.OnPropose(b2)
	hs.OnPropose(b3)
	hs.OnPropose(b4)
}
*/
