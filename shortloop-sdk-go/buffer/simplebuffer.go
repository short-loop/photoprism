package buffer

import (
	"github.com/short-loop/shortloop-common-go/models/data"
)

type SimpleBuffer struct {
	apiSamples chan *data.APISample
}

func NewSimpleBuffer(size int) SimpleBuffer {
	return SimpleBuffer{apiSamples: make(chan *data.APISample, size)}
}

func (sb *SimpleBuffer) Offer(e *data.APISample) bool {
	if sb.CanOffer() {
		// select used to convert blocking call to unblocking
		select {
		case sb.apiSamples <- e:
			return true
		default:
			return false
		}
	}
	return false
}

func (sb *SimpleBuffer) CanOffer() bool {
	return len(sb.apiSamples) < cap(sb.apiSamples)
}

func (sb *SimpleBuffer) Poll() *data.APISample {
	if sb.GetContentCount() > 0 {
		// select used to convert blocking call to unblocking
		select {
		case sample := <-sb.apiSamples:
			return sample
		default:
			return nil
		}
	}
	return nil
}

func (sb *SimpleBuffer) Clear() bool {
	if sb.GetContentCount() > 0 {
	L:
		for {
			select {
			case _, ok := <-sb.apiSamples:
				if !ok { //if channel is closed while emptying
					break L
				}
			default:
				break L
			}
		}
		return true
	}
	return false
}

func (sb *SimpleBuffer) GetContentCount() int {
	return len(sb.apiSamples)
}
