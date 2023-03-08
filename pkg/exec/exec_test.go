package exec

import "testing"

func TestExecutor(t *testing.T) {
	type args struct {
		commands []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid cmd", args: args{commands: []string{"ls", "-lht"}}, wantErr: false},
		{name: "invalid cmd", args: args{commands: []string{"ls", "-lhte"}}, wantErr: true},
		{name: "valid cat", args: args{commands: []string{"cat", "/etc/os-release"}}, wantErr: false},
		{name: "invalid cat", args: args{commands: []string{"cat", "/etc/os-releasee"}}, wantErr: true},
		{name: "print env", args: args{commands: []string{"echo", "$HOME"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewExecutor()
			if err := e.Exec(tt.args.commands).Err(); (err != nil) != tt.wantErr {
				t.Errorf("ExecBashCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBashExecutor(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid cmd", args: args{command: "ls -lht"}, wantErr: false},
		{name: "invalid cmd", args: args{command: "ls -lhte"}, wantErr: true},
		{name: "valid pipe", args: args{command: "cat /etc/os-release | grep Ubuntu"}, wantErr: false},
		{name: "invalid pipe", args: args{command: "cataa /etc/os-release | grepa Ubuntu"}, wantErr: true},
		{name: "valid multi cmd", args: args{command: "ls -lht && cat /etc/os-release"}, wantErr: false},
		{name: "invalid multi cmd", args: args{command: "ls -lhte && cat /etc/os-release"}, wantErr: true},
		{name: "print env", args: args{command: "echo $HOME"}, wantErr: false},
		{name: "print custom env", args: args{command: "CUSTOM=custom eval 'echo $CUSTOM'"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewExecutor()
			if err := e.BashExec(tt.args.command).Err(); (err != nil) != tt.wantErr {
				t.Errorf("ExecBashCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
