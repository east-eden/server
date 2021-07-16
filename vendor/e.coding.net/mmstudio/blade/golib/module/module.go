package module

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"e.coding.net/mmstudio/blade/golib/paniccatcher"
	"e.coding.net/mmstudio/blade/golib/race"
	"e.coding.net/mmstudio/blade/golib/sync2"
	"e.coding.net/mmstudio/blade/golib/version"
	"github.com/rs/zerolog/log"
)

/// module service

type Master struct {
	masterStarted   sync2.AtomicBool
	timeoutDuration time.Duration
	allAgents       []*moduleAgent
	stoppedChan     chan struct{}
	stopNotify      chan struct{}
}

var defaultMaster = NewMaster()

func NewMaster() *Master {
	return &Master{stoppedChan: make(chan struct{}), stopNotify: make(chan struct{})}
}
func (master *Master) registerOne(sif Module) *moduleAgent {
	m := &moduleAgent{mif: sif, closeChan: make(chan struct{}, 1)}
	master.allAgents = append(master.allAgents, m)
	return m
}
func (master *Master) Register(ms ...Module) {
	for i := 0; i < len(ms); i++ {
		master.registerOne(ms[i])
	}
}

func (master *Master) Stop() { close(master.stopNotify) }

func (master *Master) runAll() {
	for i := 0; i < len(master.allAgents); i++ {
		master.allAgents[i].mif.OnInit()
	}

	for i := 0; i < len(master.allAgents); i++ {
		log.Info().Msg(fmt.Sprintf("Module %s starting ...", master.allAgents[i].mif.Name()))
		master.allAgents[i].wg.Add(1)
		go master.allAgents[i].run()
		log.Info().Msg(fmt.Sprintf("Module %s started ...", master.allAgents[i].mif.Name()))
	}
}

func (master *Master) RunModule(sif Module) {
	s := master.registerOne(sif)
	if !master.masterStarted.Get() {
		return
	}
	log.Info().Msg(fmt.Sprintf("Module %s starting ...", s.mif.Name()))
	s.wg.Add(1)
	s.mif.OnInit()
	go s.run()
	log.Info().Msg(fmt.Sprintf("Module %s started", s.mif.Name()))
}

func (master *Master) closeAll() {
	for i := len(master.allAgents) - 1; i >= 0; i-- {
		agent := master.allAgents[i]
		log.Info().Msg(fmt.Sprintf("Module %s closing ...", agent.mif.Name()))
		close(agent.closeChan)
		if master.timeoutDuration == 0 {
			agent.wg.Wait()
		} else {
			if sync2.WaitTimeout(&agent.wg, master.timeoutDuration) {
				log.Info().Dur("timeout", master.timeoutDuration).Msg(fmt.Sprintf("Module %s close with timeout ...", agent.mif.Name()))
			}
		}
		closeModule(agent)
	}
}
func (master *Master) RunWithCloseTimeout(duration time.Duration, module ...Module) {
	master.timeoutDuration = duration
	master.Run(module...)
}
func (master *Master) StoppedChan() chan struct{} { return master.stoppedChan }
func (master *Master) Run(ms ...Module) {
	master.Register(ms...)
	master.runAll()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM /*, syscall.SIGTSTP*/)

	log.Info().Int("pid", os.Getpid()).Str("branch", version.Branch).Str("git_version", version.Version).Bool("race", race.Enabled).Msg("progress started")
	master.masterStarted.Set(true)

	sigGot := "user_defined"
	// Block until a signal is received
	select {
	case <-master.stopNotify:
	case sig := <-c:
		sigGot = sig.String()
	}

	log.Info().Int("pid", os.Getpid()).Str("signal", sigGot).Msg("progress closing down by signal")

	master.masterStarted.Set(false)
	master.closeAll()
	close(master.stoppedChan)
}

type Module interface {
	OnInit()
	OnClose()
	Run(closeChan chan struct{})
	Name() string
}

type moduleAgent struct {
	mif       Module
	wg        sync.WaitGroup
	closeChan chan struct{}
}

func (ma *moduleAgent) run() {
	ma.mif.Run(ma.closeChan)
	ma.wg.Done()
}

func RunWithCloseTimeout(d time.Duration, m ...Module) {
	defaultMaster.RunWithCloseTimeout(d, m...)
}
func Run(module ...Module)       { defaultMaster.Run(module...) }
func Register(ms ...Module)      { defaultMaster.Register(ms...) }
func StoppedChan() chan struct{} { return defaultMaster.StoppedChan() }
func RunModule(m Module)         { defaultMaster.RunModule(m) } //start single module

func closeModule(agent *moduleAgent) {
	paniccatcher.Do(func() {
		agent.mif.OnClose()
		log.Info().Msg(fmt.Sprintf("Module %s closed", agent.mif.Name()))
	}, func(p *paniccatcher.Panic) {
		log.Info().Msg(fmt.Sprintf("Module %s closed with reason:%s", agent.mif.Name(), p.Reason))
	})
}

var _ = StoppedChan()
