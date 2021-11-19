package get

import (
	"bytes"
	"strings"
	"testing"

	my_testing "github.com/zhigang/arctl/testing"
	"github.com/zhigang/arctl/util"
)

func getTestOptions() (*bytes.Buffer, util.Factory, *getOptions) {
	streams, _, buf, _ := my_testing.NewTestIOStreams()
	tf := my_testing.NewFactory()
	o := newGetOptions(streams)
	return buf, tf, o
}

func TestGetApps(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetApps(o, tf)
	cmd.SetOutput(buf)
	cmd.Flags().Set("all", "true")
	cmd.Run(cmd, []string{""})

	if !strings.Contains(buf.String(), "test1") {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}

func TestGetClusters(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetClusters(o, tf)
	cmd.SetOutput(buf)
	cmd.Flags().Set("type", "-1")
	cmd.Run(cmd, []string{""})

	if !(strings.Contains(buf.String(), "Test集群") && strings.Contains(buf.String(), "PRO集群")) {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}

func TestGetAppConfig(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetAppConfig(o, tf)
	cmd.SetOutput(buf)
	cmd.Flags().Set("type", "-1")
	cmd.Run(cmd, []string{"test1"})

	if !(strings.Contains(buf.String(), "test1") &&
		strings.Contains(buf.String(), "conf-test-1") &&
		strings.Contains(buf.String(), "conf-test-2")) {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}

func TestGetNodes(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetNodes(o, tf)
	cmd.SetOutput(buf)
	cmd.Flags().Set("type", "0")
	cmd.Flags().StringVarP(&o.ID, "id", "i", o.ID, "ID of resource.")
	cmd.Flags().Set("id", my_testing.Cluster1ID)
	cmd.Run(cmd, []string{"test1"})

	if !(strings.Contains(buf.String(), "Test集群") &&
		strings.Contains(buf.String(), "ecs-test-id-1") &&
		strings.Contains(buf.String(), "ecs-test-id-2")) {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}

func TestGetPods(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetPods(o, tf)
	cmd.SetOutput(buf)
	cmd.Flags().Set("type", "-1")
	cmd.Run(cmd, []string{"test1"})

	if !(strings.Contains(buf.String(), "test1") &&
		strings.Contains(buf.String(), "jck-23638-24399-998880-585fb6c46-fzs2l") &&
		strings.Contains(buf.String(), "jck-23638-24399-998880-585fb6c46-dfrec")) {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}

func TestGetUsers(t *testing.T) {
	buf, tf, o := getTestOptions()

	cmd := newCmdGetUsers(o, tf)
	cmd.SetOutput(buf)
	cmd.Run(cmd, []string{""})

	if !(strings.Contains(buf.String(), "user1") &&
		strings.Contains(buf.String(), "user2")) {
		t.Errorf("unexpected output:\n %s", buf.String())
	}
}
