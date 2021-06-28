package process

import "testing"

func TestExecute(t *testing.T) {
	tt := []struct {
		block Block
		want  string
	}{
		{
			block: Block{
				Command: "cat",
				Input:   "Hello, world!",
			},
			want: "Hello, world!",
		},
		{
			block: Block{
				Command: "sed -e s/Find/Replace/g",
				Input:   "Find",
			},
			want: "Replace",
		},
	}

	for _, tc := range tt {
		t.Run(tc.want, func(t *testing.T) {
			if testing.Short() {
				t.SkipNow()
			}
			tc.block.Execute()
			got := tc.block.Output
			if tc.want != got {
				t.Fatalf("Invalid execution, want %s, got %s", tc.want, got)
			}
		})
	}
}
