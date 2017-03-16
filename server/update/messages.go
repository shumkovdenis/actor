package update

import "github.com/AsynkronIT/protoactor-go/actor"

type Message interface {
	UpdateMessage()
}

type Check struct{}

func (*Check) UpdateMessage() {}

func (*Check) Command() string {
	return "command.update.check"
}

type Checking struct{}

func (*Checking) UpdateMessage() {}

func (*Checking) Event() string {
	return "event.update.checking"
}

type No struct{}

func (*No) UpdateMessage() {}

func (*No) Event() string {
	return "event.update.no"
}

type Available struct{}

func (*Available) UpdateMessage() {}

func (*Available) Event() string {
	return "event.update.available"
}

type CheckFailed struct{}

func (*CheckFailed) UpdateMessage() {}

func (*CheckFailed) Event() string {
	return "event.update.check.failed"
}

type Download struct{}

func (*Download) UpdateMessage() {}

func (*Download) Command() string {
	return "command.update.download"
}

type Downloading struct{}

func (*Downloading) UpdateMessage() {}

func (*Downloading) Event() string {
	return "event.update.downloading"
}

type Progress struct {
	Progress float64 `json:"progress"`
}

func (*Progress) UpdateMessage() {}

func (*Progress) Event() string {
	return "event.update.download.progress"
}

type Complete struct{}

func (*Complete) UpdateMessage() {}

func (*Complete) Event() string {
	return "event.update.download.complete"
}

type DownloadFailed struct{}

func (*DownloadFailed) UpdateMessage() {}

func (*DownloadFailed) Event() string {
	return "event.update.download.failed"
}

type Install struct{}

func (*Install) UpdateMessage() {}

func (*Install) Command() string {
	return "command.update.install"
}

type Installing struct{}

func (*Installing) UpdateMessage() {}

func (*Installing) Event() string {
	return "event.update.installing"
}

type Ready struct{}

func (*Ready) UpdateMessage() {}

func (*Ready) Event() string {
	return "event.update.install.ready"
}

type Wait struct{}

func (*Wait) UpdateMessage() {}

func (*Wait) Event() string {
	return "event.update.install.wait"
}

type InstallFailed struct{}

func (*InstallFailed) UpdateMessage() {}

func (*InstallFailed) Event() string {
	return "event.update.install.failed"
}

type Restart struct{}

func (*Restart) UpdateMessage() {}

func (*Restart) Command() string {
	return "command.update.restart"
}

type Relaunch struct {
	UpdateDataFile string `json:"update_data_file"`
	NewServerFile  string `json:"new_server_file"`
}

func (*Relaunch) UpdateMessage() {}

func (*Relaunch) Event() string {
	return "event.update.relaunch"
}

type Join struct {
	SessionPID *actor.PID
}

func (*Join) UpdateMessage() {}

type Leave struct {
	SessionPID *actor.PID
}

func (*Leave) UpdateMessage() {}
