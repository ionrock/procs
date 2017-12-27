package procs

import (
	"fmt"
	"sync"
)

type Manager struct {
	Processes map[string]*Process

	lock sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Processes: make(map[string]*Process),
	}

}

func (m *Manager) StdoutHandler(name string) OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

func (m *Manager) StderrHandler(name string) OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

func (m *Manager) Start(name, cmd string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	p := NewProcess(cmd)
	p.OutputHandler = m.StdoutHandler(name)
	p.ErrHandler = m.StderrHandler(name)
	err := p.Start()
	if err != nil {
		return err
	}

	m.Processes[name] = p
	return nil
}

func (m *Manager) StartProcess(name string, p *Process) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := p.Start()
	if err != nil {
		return err
	}

	m.Processes[name] = p
	return nil
}

func (m *Manager) Stop(name string) error {
	p, ok := m.Processes[name]
	// We don't mind stopping a process that doesn't exist.
	if !ok {
		return nil
	}

	return p.Stop()
}

func (m *Manager) Remove(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.Stop(name)
	if err != nil {
		return err
	}

	// Note that if the stop fails we don't remove it from the map of
	// processes to avoid losing the reference.
	delete(m.Processes, name)

	return nil
}

func (m *Manager) Wait() error {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Processes))

	for _, p := range m.Processes {
		go func() {
			defer wg.Done()
			p.Wait()
		}()
	}

	wg.Wait()

	return nil
}
