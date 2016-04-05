package process

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ProcessState string

const (
	StateStarted = ProcessState("StateStarted")
	StateStopped = ProcessState("StateStopped")
)

type Proc interface {
	GetProcID() string
	GetSvcName() string
	IsActive() bool
	Active() ProcRun
	ProcRuns() []ProcRun
	NumOfProcRuns() int
	State() ProcessState
	SetState(ProcessState)
	Start() error
	Stop() error
}

type ProcRun interface {
	Start() error
	Signal(syscall.Signal) error
	Kill() error
	WaitingToStopped()
	WaitingToStoppedInMillisecond(time.Duration) bool
}

type Process struct {
	ProcID      string
	SvcName     string
	Executor    []string
	Command     string
	StdoutFile  string
	StderrFile  string
	Environment map[string]string
	Metadata    map[string]string
	Pwd         string
	procRuns    []ProcRun
	active      ProcRun
	state       ProcessState // running state assigned by process manager
	rwMutex     sync.RWMutex // guard of active
}

type ProcessRun struct {
	ID          int
	Cmd         *exec.Cmd
	Error       error
	Started     time.Time
	Stopped     time.Time
	Events      []*Event
	ProcID      string
	SvcName     string
	StdoutFile  string
	StderrFile  string
	StdoutBuf   LogWriter
	StderrBuf   LogWriter
	Environment map[string]string
	WaitStatus  syscall.WaitStatus
	Pwd         string
	Stopc       chan struct{}
}

func (pr *ProcessRun) String() string {
	runInfo := make(map[string]string)
	runInfo["ID"] = strconv.Itoa(pr.ID)
	runInfo["ProcID"] = pr.ProcID
	runInfo["SvcName"] = pr.SvcName
	if pr.Error != nil {
		runInfo["Error"] = pr.Error.Error()
	}
	if pr.Started.IsZero() {
		runInfo["StartedTime"] = pr.Started.String()
	}
	if pr.Stopped.IsZero() {
		runInfo["StoppedTime"] = pr.Stopped.String()
	}
	if pr.Pwd != "" {
		runInfo["PWD"] = pr.Pwd
	}
	if pr.Cmd != nil && pr.Cmd.Process != nil && pr.Cmd.Process.Pid > 0 {
		runInfo["PID"] = strconv.Itoa(pr.Cmd.Process.Pid)
	}
	if pr.WaitStatus != 0 {
		runInfo["WaitStatus"] = strconv.Itoa(int(pr.WaitStatus))
	}

	result := make([]string, 0)
	for k, v := range runInfo {
		str := k + "=" + v
		result = append(result, str)
	}
	return "{ " + strings.Join(result, ", ") + " }"
}

type Event struct {
	Time    time.Time
	Message string
}

func NewProcess(procMgr *ProcessManager, procID string, svcName string, executor []string, command string,
	environment map[string]string, stdoutFile string, stderrFile string, metadata map[string]string,
	pwd string) Proc {

	environment = AddDefaultVars(environment)
	if _, ok := environment["SERVICE"]; !ok {
		environment["SERVICE"] = svcName
	}

	return &Process{
		ProcID:      procID,
		SvcName:     svcName,
		Executor:    executor,
		Command:     command,
		StdoutFile:  ReplaceVars(stdoutFile, environment),
		StderrFile:  ReplaceVars(stderrFile, environment),
		Environment: environment,
		Metadata:    metadata,
		Pwd:         pwd,
		procRuns:    make([]ProcRun, 0),
	}
}

func (p *Process) GetProcID() string {
	return p.ProcID
}

func (p *Process) GetSvcName() string {
	return p.SvcName
}

func (p *Process) GetExecutor() string {
	if len(p.Executor) == 0 {
		return "/bin/sh -c"
	}
	return strings.Join(p.Executor, " ")
}

func (p *Process) SetActive(pr ProcRun) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.active = pr
}

func (p *Process) SetInactive() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.active = nil
}

func (p *Process) Active() ProcRun {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.active
}

func (p *Process) IsActive() bool {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.active != nil
}

func (p *Process) State() ProcessState {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.state
}

func (p *Process) SetState(state ProcessState) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.state = state
}

func (p *Process) Start() error {
	if p.IsActive() {
		return errors.New("Process maybe already started")
	}
	pr := p.NewProcessRun()
	if err := pr.Start(); err != nil {
		return err
	}
	p.SetActive(pr)
	go func() {
		pr.WaitingToStopped()
		p.SetInactive()
	}()
	return nil
}

func (p *Process) Stop() error {
	active := p.Active()
	if active == nil {
		log.Warn("Process is inactive, no need to stop")
		return nil
	}

	if err := active.Signal(syscall.SIGINT); err != nil {
		log.Errorf("Send SIGINT to process unsuccessful, error: %v, procinfo: %v", err, active)
		return err
	} else {
		if isStopped := active.WaitingToStoppedInMillisecond(5000); isStopped {
			log.Debugf("Process terminated after SIGINT, procinfo: %v", active)
			return nil // means process finished
		}
	}
	if err := active.Signal(syscall.SIGTERM); err != nil {
		log.Errorf("Send SIGTERM to process unsuccessful, error: %v, procinfo: %v", err, active)
		return err
	} else {
		if isStopped := active.WaitingToStoppedInMillisecond(5000); isStopped {
			log.Debugf("Process terminated after SIGTERM, procinfo: %v", active)
			return nil // means process finished
		}
	}
	if err := active.Kill(); err != nil {
		log.Errorf("Send SIGKILL to process unsuccessful, error: %v, procinfo: %v", err, active)
		return err
	} else {
		if isStopped := active.WaitingToStoppedInMillisecond(1000); isStopped {
			log.Debugf("Process terminated after SIGKILL, procinfo: %v", active)
			return nil // means process finished
		}
	}
	log.Errorf("Terminate process unsuccessful, process: %v", active)
	return errors.New("Failed to stop process")
}

func (p *Process) ProcRuns() []ProcRun {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.procRuns
}

func (p *Process) NumOfProcRuns() int {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return len(p.procRuns)
}

func (p *Process) HoldProcRun(pr ProcRun) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.procRuns = append(p.procRuns, pr)
}

func (p *Process) NewProcessRun() ProcRun {
	run := p.NumOfProcRuns()
	c := p.Command
	c = ReplaceVars(c, p.Environment)
	var cmd *exec.Cmd
	if len(p.Executor) > 0 {
		cmd = exec.Command(p.Executor[0], append(p.Executor[1:], c)...)
	} else {
		bits := strings.Split(c, " ")
		cmd = exec.Command(bits[0], bits[1:]...)
	}

	vars := map[string]string{
		"PROCID": p.ProcID,
		"RUN":    strconv.Itoa(run),
	}
	if len(p.Pwd) > 0 {
		vars["PWD"] = p.Pwd
	}

	pr := &ProcessRun{
		ID:          run,
		Events:      make([]*Event, 0),
		Cmd:         cmd,
		ProcID:      p.ProcID,
		SvcName:     p.SvcName,
		StdoutFile:  ReplaceVars(p.StdoutFile, vars),
		StderrFile:  ReplaceVars(p.StderrFile, vars),
		Environment: make(map[string]string),
		Pwd:         p.Pwd,
	}
	for k, v := range p.Environment {
		pr.Environment[k] = v
	}

	p.HoldProcRun(pr)
	return pr
}

func (pr *ProcessRun) Start() error {
	pr.Started = time.Now()
	pr.Stopc = make(chan struct{})

	stdout, err := pr.Cmd.StdoutPipe()
	if err != nil {
		pr.Error = err
		close(pr.Stopc)
		return err
	}
	stderr, err := pr.Cmd.StderrPipe()
	if err != nil {
		pr.Error = err
		close(pr.Stopc)
		return err
	}
	if len(pr.StdoutFile) > 0 {
		wr, err := NewFileLogWriter(pr.StdoutFile)
		if err != nil {
			log.Error("Unable to open file %s: %v", pr.StdoutFile, err)
			pr.StdoutBuf = NewInMemoryLogWriter()
		} else {
			pr.StdoutBuf = wr
		}
	} else {
		pr.StdoutBuf = NewInMemoryLogWriter()
	}
	if len(pr.StderrFile) > 0 {
		wr, err := NewFileLogWriter(pr.StderrFile)
		if err != nil {
			log.Error("Unable to open file %s: %v", pr.StderrFile, err)
			pr.StderrBuf = NewInMemoryLogWriter()
		} else {
			pr.StderrBuf = wr
		}
	} else {
		pr.StderrBuf = NewInMemoryLogWriter()
	}
	if len(pr.Pwd) > 0 {
		pr.Cmd.Dir = pr.Pwd
	}
	for k, v := range pr.Environment {
		pr.Cmd.Env = append(pr.Cmd.Env, k+"="+v)
	}

	err = pr.Cmd.Start()
	if err != nil {
		pr.Error = err
		pr.StdoutBuf.Close()
		pr.StderrBuf.Close()
		close(pr.Stopc)
		return err
	}
	if pr.Cmd.Process == nil {
		pr.Error = errors.New("Start process failed")
		pr.StdoutBuf.Close()
		pr.StderrBuf.Close()
		close(pr.Stopc)
		return pr.Error
	}

	ev := &Event{time.Now(), fmt.Sprintf("Process %s[%s] started as PID: %d", pr.SvcName, pr.ProcID, pr.Cmd.Process.Pid)}
	log.Info(ev.Message)
	pr.Events = append(pr.Events, ev)

	go func() {
		go func() {
			io.Copy(pr.StdoutBuf, stdout)
			pr.StdoutBuf.Close()
		}()
		go func() {
			io.Copy(pr.StderrBuf, stderr)
			pr.StderrBuf.Close()
		}()
		pr.Cmd.Wait()
		ps := pr.Cmd.ProcessState
		sy := ps.Sys().(syscall.WaitStatus)
		ev := &Event{time.Now(), fmt.Sprintf("Process %s[%s], PID: %d exited with status: %d", pr.SvcName, pr.ProcID, pr.Cmd.Process.Pid, sy.ExitStatus())}
		log.Info(ev.Message)
		pr.Events = append(pr.Events, ev)
		pr.Stopped = time.Now()
		close(pr.Stopc)
	}()

	return nil
}

func (pr *ProcessRun) Signal(sig syscall.Signal) error {
	if pr.Cmd == nil || pr.Cmd.Process == nil {
		return errors.New("Process not started")
	}
	return pr.Cmd.Process.Signal(sig)
}

func (pr *ProcessRun) Kill() error {
	if pr.Cmd == nil || pr.Cmd.Process == nil {
		return errors.New("Process not started")
	}
	return pr.Cmd.Process.Kill()
}

func (pr *ProcessRun) WaitingToStopped() {
	<-pr.Stopc
}

func (pr *ProcessRun) WaitingToStoppedInMillisecond(timeout time.Duration) bool {
	select {
	case <-pr.Stopc:
		return true
	case <-time.After(timeout * time.Millisecond):
		return false
	}
}
