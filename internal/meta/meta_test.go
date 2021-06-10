package meta_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/maaslalani/slides/internal/meta"
)

func TestMeta_ParseHeader(t *testing.T) {
	type fields struct {
		Theme string
	}
	type args struct {
		header string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Meta
		wantErr bool
	}{
		{
			name: "Parse theme from header",
			fields: fields{
				Theme: "dark",
			},
			args: args{
				header: fmt.Sprintf("---\ntheme: %q\n", "dark"),
			},
			want: &meta.Meta{
				Theme: "dark",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meta.Meta{}
			got, err := m.ParseHeader(tt.args.header)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Nil(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *meta.Meta
	}{
		{name: "Create meta struct", want: &meta.Meta{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, meta.New(), tt.want)
		})
	}
}

func ExampleMeta_ParseHeader() {
	header := `
---
theme: "dark"
---
`
	// Parse the header from the markdown
	// file
	m, err := meta.New().ParseHeader(header)
	if err != nil {
		return
	}

	// Print the return theme
	// meta
	fmt.Println(m.Theme)
}
