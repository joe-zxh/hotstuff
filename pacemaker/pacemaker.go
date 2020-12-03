package pacemaker

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/joe-zxh/hotstuff"
	"github.com/joe-zxh/hotstuff/config"
	"github.com/joe-zxh/hotstuff/internal/logging"
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// RoundRobin change leader in a RR fashion. The amount of commands to be executed before it changes leader can be customized.
type RoundRobin struct {
	*hotstuff.HotStuff

	termLength int
	schedule   []config.ReplicaID
}

// NewRoundRobin returns a new round robin pacemaker
func NewRoundRobin(termLength int, schedule []config.ReplicaID, timeout time.Duration) *RoundRobin {
	return &RoundRobin{
		termLength: termLength,
		schedule:   schedule,
	}
}

// GetLeader returns the fixed ID of the leader for the view height
func (p *RoundRobin) GetLeader(view int) config.ReplicaID {
	term := int(math.Ceil(float64(view) / float64(p.termLength)))
	return p.schedule[term%len(p.schedule)]
}

// Run runs the pacemaker which will beat when the previous QC is completed
func (p *RoundRobin) Run(ctx context.Context) {

}
