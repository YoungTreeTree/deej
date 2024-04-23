package deej

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type MockInput struct {
	deej  *Deej
	input chan SliderMoveEvent
}

func NewMockInput(deej *Deej, logger *zap.SugaredLogger) (*MockInput, error) {
	logger = logger.Named("MockInput")

	mi := &MockInput{
		deej:  deej,
		input: make(chan SliderMoveEvent, 10),
	}

	logger.Debug("Created mockInput instance")

	go func() {
		s := ""
		for {
			fmt.Scanln(&s)
			logger.Debugf("Mock Input ============ %s ==========", s)
			ss := strings.Split(s, "|")
			if len(ss) == 2 {
				sliderId, err := strconv.ParseInt(ss[0], 10, 32)
				if err != nil {
					logger.Debugf("Mock Input parse %s error", ss[0])
					continue
				}
				percentValue, err := strconv.ParseFloat(ss[1], 32)
				if err != nil {
					logger.Debugf("Mock Input parse %s error", ss[1])
					continue
				}
				mi.input <- SliderMoveEvent{
					SliderID:     int(sliderId),
					PercentValue: float32(percentValue),
				}
			}
		}
	}()

	return mi, nil
}

func (mi *MockInput) SubscribeToSliderMoveEvents() chan SliderMoveEvent {
	return mi.input
}
