package model

import (
	"runtime"

	"github.com/Karmenzind/kd/pkg"
	"go.uber.org/zap"
	"golang.org/x/term"
)

type RunInfo struct {
	StartTime int64
	PID       int
	Port      string
	ExeName   string
	ExePath   string
	Version   string

	OS *pkg.OSInfo

	isServer   bool
	termHeight int
	termWidth  int
}

func (r *RunInfo) IsServer() bool {
	return r.isServer
}

func (r *RunInfo) SetServer(v bool) {
	r.isServer = v
}

func (r *RunInfo) SetPort(v string) {
	r.Port = v
}

func (r *RunInfo) SetOSInfo() {
	var err error
	if r.OS, err = pkg.GetOSInfo(); err != nil {
		zap.S().Warn("Failed to fetch os info: %s. (Current GOOS: %s)", err, runtime.GOOS)
	}
}

func (r *RunInfo) GetOSInfo() *pkg.OSInfo {
	if r.OS == nil {
		r.SetOSInfo()
	}
	return r.OS
}

func (r *RunInfo) GetTermSize() (int, int, error) {
	if r.termHeight > 0 && r.termWidth > 0 {
		return r.termWidth, r.termHeight, nil
	}
	w, h, err := term.GetSize(0)
	if err != nil {
		return 0, 0, err
	}
	r.termHeight = h
	r.termWidth = w
	return w, h, nil
}

func (r *RunInfo) SaveToFile(path string) (err error) {
	err = pkg.SaveJson(path, r)
	if err == nil {
		zap.S().Infof("Recorded running information of daemon %+v", r)
	} else {
		zap.S().Warnf("Failed to record running info of daemon %+v", err)
	}
	return
}
