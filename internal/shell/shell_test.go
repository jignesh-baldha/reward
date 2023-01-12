package shell

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"reward/internal/util"
)

type ShellTestSuite struct {
	suite.Suite
}

func (suite *ShellTestSuite) SetupTest() {
}

func TestShellTestSuite(t *testing.T) {
	suite.Run(t, new(ShellTestSuite))
}

func (suite *ShellTestSuite) TestLocalShell_Execute() {
	// Cannot run in parallel execution as it uses the os.stdout
	type fields struct {
		CatchStdout bool
	}

	type args struct {
		name string
		arg  []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test bash command",
			fields: fields{
				CatchStdout: true,
			},
			args: args{
				name: "/bin/bash",
				arg:  []string{"-c", "echo test"},
			},
			want:    []byte("test\n"),
			wantErr: false,
		},
		{
			name: "test bash command without writing to stdout",
			fields: fields{
				CatchStdout: true,
			},
			args: args{
				name: "/bin/bash",
				arg:  []string{"-c", "echo test"},
			},
			want:    []byte("test\n"),
			wantErr: false,
		},
		{
			name: "invalid bash command",
			fields: fields{
				CatchStdout: true,
			},
			args: args{
				name: "/bin/false",
				arg:  []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test docker compose version",
			fields: fields{
				CatchStdout: true,
			},
			args: args{
				name: "docker",
				arg:  []string{"compose", "version", "--short"},
			},
			want:    []byte("2.13.0\n"),
			wantErr: false,
		},
		{
			name: "test docker-compose version",
			fields: fields{
				CatchStdout: true,
			},
			args: args{
				name: "docker-compose",
				arg:  []string{"version", "--short"},
			},
			want:    []byte("2.13.0\n"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			c := NewLocalShellWithOpts(WithCatchOutput(tt.fields.CatchStdout))
			got, err := c.Execute(tt.args.name, tt.args.arg...)
			w.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %s, wantErr %t", err, tt.wantErr)

				return
			}

			assert.Equal(t, tt.want, got)

			if !tt.fields.CatchStdout {
				out, _ := io.ReadAll(r)
				// The io.ReadAll() function will return an empty byte slice but the want variable is nil
				// if !reflect.DeepEqual(out, tt.want) {
				if len(out) != len(tt.want) {
					assert.Errorf(t, err, "Execute() got = %s, want %s", out, tt.want)
				}
			}

			os.Stdout = rescueStdout
		})
	}
}

func (suite *ShellTestSuite) TestMockShell_Execute() {
	suite.T().Parallel()

	type fields struct {
		Output      []byte
		Err         error
		LastCommand string
	}

	type args struct {
		name string
		args []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test mock output",
			fields: fields{
				Output:      []byte("test"),
				Err:         nil,
				LastCommand: "test",
			},
			args: args{
				name: "",
				args: nil,
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "test mock error",
			fields: fields{
				Output:      nil,
				Err:         fmt.Errorf("test error"),
				LastCommand: "",
			},
			args: args{
				name: "",
				args: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewMockShell(
				tt.fields.LastCommand,
				tt.fields.Output,
				tt.fields.Err,
			)
			got, err := c.Execute(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "Execute() error = %s, wantErr %s", err, tt.wantErr)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func (suite *ShellTestSuite) TestLocalShell_ExecuteWithOptions() {
	type fields struct {
		CatchStdout    *bool
		SuppressStdout *bool
	}

	type args struct {
		name string
		args []string
		opts []Opt
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test without suppress stdout",
			fields: fields{
				CatchStdout: util.BoolPtr(false),
			},
			args: args{
				name: "/bin/bash",
				args: []string{"-c", "echo test"},
				opts: []Opt{WithCatchOutput(false)},
			},
			want:    []byte(nil),
			wantErr: false,
		},
		{
			name: "test with suppress stdout",
			fields: fields{
				CatchStdout: util.BoolPtr(false),
			},
			args: args{
				name: "/bin/bash",
				args: []string{"-c", "echo test"},
				opts: []Opt{WithCatchOutput(true)},
			},
			want:    []byte("test\n"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			c := &LocalShell{
				CatchStdout: tt.fields.CatchStdout,
			}

			got, err := c.ExecuteWithOptions(tt.args.name, tt.args.args, tt.args.opts...)
			w.Close()

			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "ExecuteWithOptions() error = %s, wantErr %s", err, tt.wantErr)

				return
			}
			assert.Equal(t, tt.want, got)

			if len(tt.args.opts) > 0 && reflect.DeepEqual(tt.args.opts[0], WithCatchOutput(true)) {
				out, _ := io.ReadAll(r)
				// The io.ReadAll() function will return an empty byte slice but the want variable is nil
				// if !reflect.DeepEqual(out, tt.want) {
				if len(out) != len(tt.want) {
					assert.Errorf(t, err, "ExecuteWithOptions() got = %s, want %s", out, tt.want)
				}
			}

			os.Stdout = rescueStdout
		})
	}
}

func (suite *ShellTestSuite) TestMockShell_ExecuteWithOptions() {
	type fields struct {
		Output      []byte
		Err         error
		LastCommand string
	}

	type args struct {
		name string
		args []string
		opts []Opt
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test mock output",
			fields: fields{
				Output:      []byte("test"),
				Err:         nil,
				LastCommand: "test",
			},
			args: args{
				name: "",
				args: nil,
				opts: nil,
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "test mock error",
			fields: fields{
				Output:      nil,
				Err:         fmt.Errorf("test error"),
				LastCommand: "",
			},
			args: args{
				name: "",
				args: nil,
				opts: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			c := &MockShell{
				Output:      tt.fields.Output,
				Err:         tt.fields.Err,
				LastCommand: tt.fields.LastCommand,
			}

			got, err := c.ExecuteWithOptions(tt.args.name, tt.args.args, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "ExecuteWithOptions() error = %s, wantErr %t", err, tt.wantErr)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
