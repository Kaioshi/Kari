package timer

import "time"
import "Kari/lib/logger"

var Tickers map[int]*Ticker = make(map[int]*Ticker)

type Ticker struct {
	duration *int
	ticker   *time.Ticker
	quit     chan struct{}
	events   []*TickerEvent
}

type TickerEvent struct {
	handle   string
	callback func()
}

func AddEvent(handle string, ticker int, callback func()) {
	if !tickerExists(&ticker) {
		StartTicker(ticker)
	}
	Tickers[ticker].events = append(Tickers[ticker].events, &TickerEvent{handle, callback})

}

func RemoveEvent(handle string, ticker int) {
	for i, event := range Tickers[ticker].events {
		if event.handle == handle {
			events := Tickers[ticker].events
			events = append(events[:i], events[i+1:]...)
			Tickers[ticker].events = events
		}
	}
}

func eventExists(ticker *int, handle *string) bool {
	if tickerExists(ticker) {
		for _, event := range Tickers[*ticker].events {
			if event.handle == *handle {
				return true
			}
		}
	}
	return false
}

func tickerExists(ticker *int) bool {
	_, ok := Tickers[*ticker]
	return ok
}

func StartTicker(seconds int) {
	if tickerExists(&seconds) {
		logger.Debug("Not starting another " + string(seconds) + "s ticker.")
		return
	}
	Tickers[seconds] = &Ticker{
		duration: &seconds,
		ticker:   time.NewTicker(time.Duration(seconds) * time.Second),
		quit:     make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-Tickers[seconds].ticker.C:
				for _, event := range Tickers[seconds].events {
					event.callback()
				}
			case <-Tickers[seconds].quit:
				Tickers[seconds].ticker.Stop()
				logger.Debug("Stopped " + string(seconds) + "s ticker.")
				return
			}
		}
	}()
}

func StopTicker(ticker int) {
	if !tickerExists(&ticker) {
		logger.Debug("Tried to stop non-existant " + string(ticker) + "s ticker.")
		return
	}
	close(Tickers[ticker].quit)
	logger.Debug("Stopped " + string(ticker) + "s ticker.")
}
