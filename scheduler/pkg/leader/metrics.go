package leader

import "github.com/prometheus/client_golang/prometheus"

// nolint: gochecknoglobals
var (
	totalLeadershipClaimed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "scheduler",
		Subsystem: "leader",
		Name:      "claimed_total",
		Help:      "Total number of times this candidate has claimed leadership",
	}, []string{"uuid"})
	totalLeadershipRefreshed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "scheduler",
		Subsystem: "leader",
		Name:      "refreshed_total",
		Help:      "Total number of times this leader has refreshed leadership",
	}, []string{"uuid"})
	totalLeadershipRead = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "scheduler",
		Subsystem: "leader",
		Name:      "read_total",
		Help:      "Total number of times this candidate or leader has checked leadership",
	}, []string{"uuid"})
	currentLeader = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "scheduler",
		Subsystem: "leader",
		Name:      "current_leader",
		Help:      "1 if this is the leader; 0 if not",
	}, []string{"uuid"})
)

func RegisterMetrics(reg prometheus.Registerer) {
	reg.MustRegister(
		totalLeadershipClaimed,
		totalLeadershipRead,
		totalLeadershipRefreshed,
		currentLeader,
	)
}
