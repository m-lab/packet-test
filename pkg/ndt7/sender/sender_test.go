package sender

import (
	"testing"

	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/tcp-info/inetdiag"
	"github.com/m-lab/tcp-info/tcp"
)

func Test_terminateTest(t *testing.T) {
	tests := []struct {
		name string
		p    *Params
		m    model.Measurement
		want bool
	}{
		{
			name: "max_cwnd_gain and max_elapsed_time",
			p: &Params{
				MaxCwndGain:    10,
				MaxElapsedTime: 10,
			},
			m: model.Measurement{
				BBRInfo: &model.BBRInfo{
					BBRInfo: inetdiag.BBRInfo{
						CwndGain: 10,
					},
				},
				TCPInfo: &model.TCPInfo{
					ElapsedTime: 10,
				},
			},
			want: true,
		},
		{
			name: "max_cwnd_gain and early_exit",
			p: &Params{
				MaxCwndGain: 10,
				MaxBytes:    10,
			},
			m: model.Measurement{
				BBRInfo: &model.BBRInfo{
					BBRInfo: inetdiag.BBRInfo{
						CwndGain: 10,
					},
				},
				TCPInfo: &model.TCPInfo{
					LinuxTCPInfo: tcp.LinuxTCPInfo{
						BytesAcked: 10,
					},
				},
			},
			want: true,
		},
		{
			name: "max_cwnd_gain only",
			p: &Params{
				MaxCwndGain: 10,
			},
			m: model.Measurement{
				BBRInfo: &model.BBRInfo{
					BBRInfo: inetdiag.BBRInfo{
						CwndGain: 10,
					},
				},
			},
			want: true,
		},
		{
			name: "max_elapsed_time only",
			p: &Params{
				MaxElapsedTime: 10,
			},
			m: model.Measurement{
				TCPInfo: &model.TCPInfo{
					ElapsedTime: 10,
				},
			},
			want: true,
		},
		{
			name: "early_exit only",
			p: &Params{
				MaxBytes: 10,
			},
			m: model.Measurement{
				TCPInfo: &model.TCPInfo{
					LinuxTCPInfo: tcp.LinuxTCPInfo{
						BytesAcked: 10,
					},
				},
			},
			want: true,
		},
		{
			name: "no limit",
			p:    &Params{},
			m:    model.Measurement{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := terminateTest(tt.p, tt.m); got != tt.want {
				t.Errorf("terminateTest() = %v, want %v", got, tt.want)
			}
		})
	}
}
