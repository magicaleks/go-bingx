package bingx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type websocketServiceTestSuite struct {
	baseTestSuite
	origWsServe func([]byte, *WsConfig, WsHandler, ErrHandler) (chan struct{}, chan struct{}, error)
	serveCount  int
}

func TestWebsocketService(t *testing.T) {
	suite.Run(t, new(websocketServiceTestSuite))
}

func (s *websocketServiceTestSuite) SetupTest() {
	s.origWsServe = wsServe
}

func (s *websocketServiceTestSuite) TearDownTest() {
	wsServe = s.origWsServe
	s.serveCount = 0
}

func (s *websocketServiceTestSuite) mockWsServe(data [][]byte, err error) {
	wsServe = func(initMessage []byte, config *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
		s.serveCount++
		doneC = make(chan struct{})
		stopC = make(chan struct{})
		go func() {
			<-stopC
			close(doneC)
		}()
		for _, d := range data {
			handler(d)
		}
		if err != nil {
			errHandler(err)
		}
		return doneC, stopC, nil
	}
}

func (s *websocketServiceTestSuite) assertWsServe(count ...int) {
	e := 1
	if len(count) > 0 {
		e = count[0]
	}
	s.r().Equal(e, s.serveCount)
}

func (s *websocketServiceTestSuite) TestKlineServe() {
	data := [][]byte{
		[]byte(`{
			"s": "ETHBTC",
			"o": "0.10278577",
			"c": "0.10278645",
			"h": "0.10278712",
			"l": "0.10278518",
			"v": "17.47929838",
			"t": 1499404860000
		}`),
		[]byte(`{
			"s": "ETHBTC",
			"o": "0.10278575",
			"c": "0.10278648",
			"h": "0.10278718",
			"l": "0.10278513",
			"v": "17.47929834",
			"t": 1499404860000
		}`),
		[]byte(`{
			"s": "ETHBTC",
			"o": "0.10278579",
			"c": "0.10278641",
			"h": "0.10278711",
			"l": "0.10278512",
			"v": "17.47929830",
			"t": 1499404860060
		}`),
	}
	fakeErrMsg := "fake error"
	s.mockWsServe(data, errors.New(fakeErrMsg))
	defer s.assertWsServe()

	steps := 0

	doneC, stopC, err := WsKlineServe("ETHBTC", Interval1, func(event *WsKlineEvent) {
		var e *WsKlineEvent
		switch steps {
		case 0:
			e = &WsKlineEvent{
				Symbol:    "ETHBTC",
				Open:      0.10278577,
				Close:     0.10278645,
				High:      0.10278712,
				Low:       0.10278513,
				Volume:    17.47929834,
				Time:      1499404860000,
				Completed: false,
			}
			steps++
		case 1:
			e = &WsKlineEvent{
				Symbol:    "ETHBTC",
				Open:      0.10278575,
				Close:     0.10278648,
				High:      0.10278718,
				Low:       0.10278518,
				Volume:    17.47929838,
				Time:      1499404860000,
				Completed: true,
			}
		}

		s.assertWsKlineEventEqual(e, event)
	}, func(err error) {
		s.r().EqualError(err, fakeErrMsg)
	})
	s.r().NoError(err)
	stopC <- struct{}{}
	<-doneC
}

func (s *websocketServiceTestSuite) assertWsKlineEventEqual(e, a *WsKlineEvent) {
	r := s.r()
	r.Equal(e.Symbol, a.Symbol, "Symbol")
	r.Equal(e.Open, a.Open, "Open")
	r.Equal(e.Close, a.Close, "Close")
	r.Equal(e.High, a.High, "High")
	r.Equal(e.Low, a.Low, "Low")
	r.Equal(e.Volume, a.Volume, "Volume")
	r.Equal(e.Time, a.Time, "Time")
	r.Equal(e.Completed, a.Completed, "Completed")
}
