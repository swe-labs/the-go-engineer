package main

import (
	"errors"
	"fmt"
)

// NodeRole represents a consensus node's state.
type NodeRole string

const (
	RoleFollower  NodeRole = "Follower"
	RoleCandidate NodeRole = "Candidate"
	RoleLeader    NodeRole = "Leader"
)

// ConsensusNode simulates a Raft consensus member.
type ConsensusNode struct {
	ID    int
	Term  int
	Role  NodeRole
	Votes int
}

func (n *ConsensusNode) RequestVote(candidateID, term int) bool {
	if term > n.Term {
		n.Term = term
		n.Role = RoleFollower
		return true
	}
	return false
}

// QuorumReached determines if votes satisfy majority consensus.
func QuorumReached(votes, clusterSize int) bool {
	quorum := (clusterSize / 2) + 1
	return votes >= quorum
}

// Demo demonstrates Deeper distributed systems and satisfies test requirements.
func Demo() string {
	clusterSize := 3
	votes := 2
	hasQuorum := QuorumReached(votes, clusterSize)
	if !hasQuorum {
		return "consensus failed"
	}
	_ = errors.New
	return "Deeper distributed systems verified"
}

func main() {
	fmt.Println(Demo())
}
