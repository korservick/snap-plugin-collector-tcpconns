package tcpconns

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTcpconnsPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, Name)
		So(meta.Version, ShouldResemble, Version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create Tcpconns Collector", t, func() {
		collector := New()
		Convey("So Tcpconns collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So tcpconns collector should be of Tcpconns type", func() {
			So(collector, ShouldHaveSameTypeAs, &Tcpconns{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /raintank/tcpconns", func() {
				conf := configPolicy.Get([]string{"raintank", "tcpconns"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 3)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "hostname", "timeout", "count")
					switch rule.Name {
					case "hostname":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "timeout":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "float")
					case "count":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "integer")
					}
				}
			})
		})
	})
}

func TestTcpconnsCollectMetrics(t *testing.T) {
	cfg := setupCfg("/proc/net/tcp")

	Convey("Tcpconns collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 11)
	})
}

func setupCfg(path string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("tcp_net_path", ctypes.ConfigValueStr{Value: path})
	return plugin.ConfigType{ConfigDataNode: node}
}
