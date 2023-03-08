package apt

import (
	"testing"
)

func TestInstall(t *testing.T) {
	type args struct {
		packages []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "multi pkg", args: args{packages: []string{"curl", "vim", "git"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewInstaller()
			if err := a.Install(tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("aptInstaller.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
