package deej

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/thoas/go-funk"
)

type SessionOpTyp string

const (
	sessionOpTyp_mute   = "mute"
	sessionOpTyp_volume = "volume"
)

func SessionOpTypFromString(s string) SessionOpTyp {
	switch s {
	case "mute":
		return sessionOpTyp_mute
	case "volume":
		return sessionOpTyp_volume
	default:
		return sessionOpTyp_volume
	}
}

type sliderMap struct {
	m    map[int][]string
	mm   map[int]SessionOpTyp
	lock sync.Locker
}

func newSliderMap() *sliderMap {
	return &sliderMap{
		m:    make(map[int][]string),
		mm:   make(map[int]SessionOpTyp),
		lock: &sync.Mutex{},
	}
}

func sliderMapFromConfigs(userMapping map[string][]string, internalMapping map[string][]string, userModeMapping map[string]string) *sliderMap {
	resultMap := newSliderMap()

	// copy targets from user config, ignoring empty values
	for sliderIdxString, targets := range userMapping {
		sliderIdx, _ := strconv.Atoi(sliderIdxString)

		resultMap.set(sliderIdx, funk.FilterString(targets, func(s string) bool {
			return s != ""
		}))
	}

	// add targets from internal configs, ignoring duplicate or empty values
	for sliderIdxString, targets := range internalMapping {
		sliderIdx, _ := strconv.Atoi(sliderIdxString)

		existingTargets, ok := resultMap.get(sliderIdx)
		if !ok {
			existingTargets = []string{}
		}

		filteredTargets := funk.FilterString(targets, func(s string) bool {
			return (!funk.ContainsString(existingTargets, s)) && s != ""
		})

		existingTargets = append(existingTargets, filteredTargets...)
		resultMap.set(sliderIdx, existingTargets)
	}

	for sliderIdxString, mode := range userModeMapping {
		sliderIdx, _ := strconv.Atoi(sliderIdxString)

		resultMap.setMode(sliderIdx, SessionOpTypFromString(mode))
	}

	return resultMap
}

func (m *sliderMap) iterate(f func(int, []string)) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for key, value := range m.m {
		f(key, value)
	}
}

func (m *sliderMap) get(key int) ([]string, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	value, ok := m.m[key]
	return value, ok
}

func (m *sliderMap) set(key int, value []string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.m[key] = value
}

func (m *sliderMap) getMode(key int) SessionOpTyp {
	m.lock.Lock()
	defer m.lock.Unlock()

	value, ok := m.mm[key]
	if !ok {
		return sessionOpTyp_volume
	}
	return value
}

func (m *sliderMap) setMode(key int, value SessionOpTyp) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.mm[key] = value
}

func (m *sliderMap) String() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	sliderCount := 0
	targetCount := 0

	for _, value := range m.m {
		sliderCount++
		targetCount += len(value)
	}

	return fmt.Sprintf("<%d sliders mapped to %d targets>", sliderCount, targetCount)
}
