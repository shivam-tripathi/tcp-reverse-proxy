package tcpReverseProxy

import (
	"math/rand"
	"sync"
)

type State struct {
	sync.Mutex
	clientTargetMapping map[string]string
	targets             []string
	validTargets        map[string]bool
}

func NewState() *State {
	ret := new(State)
	ret.clientTargetMapping = make(map[string]string)
	ret.targets = make([]string, 0)
	return ret
}

func (m *State) AddTarget(targetIP string) {
	m.Lock()
	defer m.Unlock()
	if isValid, ok := m.validTargets[targetIP]; !ok || !isValid {
		m.targets = append(m.targets, targetIP)
		m.validTargets[targetIP] = true
	}
}

func (m *State) RemoveTarget(targetIP string) {
	m.Lock()
	targetIdx := -1
	for idx, tIP := range m.targets {
		if tIP == targetIP {
			targetIdx = idx
			break
		}
	}
	if targetIdx != -1 {
		m.validTargets[targetIP] = false
		m.targets = append(m.targets[:targetIdx], m.targets[targetIdx+1:]...)
		for cIP, tIP := range m.clientTargetMapping {
			if tIP == targetIP {
				m.clientTargetMapping[cIP] = m.targets[rand.Intn(len(m.targets))]
			}
		}
	}
	defer m.Unlock()
}
