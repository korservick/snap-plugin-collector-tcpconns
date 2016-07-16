package tcpconns

import (
	"strconv"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "tcpconns"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	// make sure that we actually satisify requierd interface
	_ plugin.CollectorPlugin = (*Tcpconns)(nil)

	metricNames []string
)

type Tcpconns struct {
}

func New() *Tcpconns {
	return &Tcpconns{}
}

// CollectMetrics collects metrics for testing
func (p *Tcpconns) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	var err error

	conf := mts[0].Config().Table()
	var statpath string
	statpathConf, ok := conf["tcp_net_path"]
	if ok {
		statpath = statpathConf.(ctypes.ConfigValueStr).Value
	} else {
		statpath = "/proc/net/tcp"
	}

	metrics, err := tcpconns(statpath, mts)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func tcpconns(tcp_net_path string, mts []plugin.MetricType) ([]plugin.MetricType, error) {
	results, err := GatherTcpconnsInfo(tcp_net_path)
	if err != nil {
		return nil, err
	}
	runTime := time.Now()

	var arrLen int
	listeners := make(map[uint16]bool)
	for p, r := range results.Connections {
		arrLen++
		_, listeners[p] = r.Status["LISTEN"]
		arrLen += len(r.Status)
	}
	metrics := make([]plugin.MetricType, 0, arrLen)
	for port, r := range results.Connections {
		if !listeners[port] {
			continue
		}
		for _, m := range mts {
			stat := m.Namespace()[4].Value
			if value, ok := r.Status[stat]; ok {
				mt := plugin.MetricType{
					Data_:      value,
					Namespace_: core.NewNamespace("raintank", "tcpconns", strconv.FormatUint(uint64(port), 10), "local", stat),
					Timestamp_: runTime,
					Version_:   m.Version(),
				}
				metrics = append(metrics, mt)
			}
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (p *Tcpconns) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	if metricNames == nil || len(metricNames) == 0 {
		// Assemble metricNames
		metricNames = make([]string, len(TcpStates))
		i := 0
		for _, s := range TcpStates {
			metricNames[i] = s
			i++
		}
	}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("raintank", "tcpconns").
				AddDynamicElement("port", "the tcp port").
				AddStaticElement("local").
				AddStaticElement(metricName),
		})
	}
	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (p *Tcpconns) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("tcp_net_path", false, "/proc/net/tcp")
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	c.Add([]string{"raintank", "tcpconns"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(5000),
	)
}
