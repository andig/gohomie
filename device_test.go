package homie

import (
	"strings"
	"testing"
)

func TestDevice(t *testing.T) {
	if d := NewDevice("id"); d.ID != "id" || d.Name != "id" ||
		d.State != StateInit || d.Version != Version {
		t.Fail()
	}
}
func TestDeviceAdd(t *testing.T) {
	d := NewDevice("id")
	n := NewNode("node")

	if err := d.Add(n); err != nil || d.Nodes["node"] != n {
		t.Fail()
	}

	if err := d.Add(n); err == nil {
		t.Fail()
	}
}

func TestDevicePublish(t *testing.T) {
	d := NewDevice("id")
	d.Name = "name"
	d.State = StateReady
	d.Extensions = append(d.Extensions, "foo", "bar")
	d.Add(NewNode("n1"))
	d.Add(NewNode("n2"))

	exp := []struct {
		t, m string
		r    bool
	}{
		{">/id/$name", "name", true},
		{">/id/$version", Version, true},
		{">/id/$state", string(StateReady), true},
		{">/id/$extensions", "foo,bar", true},
		{">/id/$nodes", "n1,n2", true},
	}

	idx := 0
	d.Publish(func(topic string, retained bool, message string) {
		// filter node properties
		if strings.Contains(topic, "/n1/") || strings.Contains(topic, "/n2/") {
			return
		}

		if idx >= len(exp) {
			t.Errorf("unexpected index %d", idx)
			return
		}

		e := exp[idx]
		if e.t != topic || e.m != message || e.r != retained {
			t.Errorf("expected %s %s %v", e.t, e.m, e.r)
			t.Errorf("got %s %s %v", topic, message, retained)
		}
		idx++
	}, ">")

	if idx != len(exp) {
		t.Errorf("unexpected number of matches %d", idx)
	}
}