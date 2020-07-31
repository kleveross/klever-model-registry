package util

import "testing"

func TestExecOSCommand(t *testing.T) {
	type args struct {
		command string
		args    []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exec ls command",
			args: args{
				command: "ls",
				args:    []string{},
			},
			wantErr: false,
		},
		{
			name: "exec ls xxx failed",
			args: args{
				command: "ls xxx",
				args:    []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecOSCommand(tt.args.command, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("ExecOSCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
