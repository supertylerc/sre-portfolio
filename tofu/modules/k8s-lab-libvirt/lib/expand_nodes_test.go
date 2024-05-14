package lib_test

import (
	"maps"
	"testing"

	"github.com/supertylerc/sre-portfolio/tofu/modules/k8s-lab-libvirt/lib"
)

func TestExpandNodes(t *testing.T) {
	type testCase struct {
		input []lib.Node
		want  map[string]lib.Node
	}
	tt := []testCase{
		testCase{
			input: []lib.Node{
				lib.Node{"control-plane", 1, 1, 1, 1},
			},
			want: map[string]lib.Node{
				"control-plane-0": lib.Node{"control-plane", 0, 1, 1, 1},
			},
		},
		testCase{
			input: []lib.Node{
				lib.Node{"control-plane", 1, 2, 8, 40},
				lib.Node{"node", 6, 2, 8, 40},
			},
			want: map[string]lib.Node{
				"control-plane-0": lib.Node{"control-plane", 0, 2, 8, 40},
				"node-0":          lib.Node{"node", 0, 2, 8, 40},
				"node-1":          lib.Node{"node", 1, 2, 8, 40},
				"node-2":          lib.Node{"node", 2, 2, 8, 40},
				"node-3":          lib.Node{"node", 3, 2, 8, 40},
				"node-4":          lib.Node{"node", 4, 2, 8, 40},
				"node-5":          lib.Node{"node", 5, 2, 8, 40},
			},
		},
	}
	for _, test := range tt {
		got := lib.ExpandNodes(test.input)
		if !maps.Equal(test.want, got) {
			t.Errorf("got: %q\nwant: %q", got, test.want)
		}
	}
}
