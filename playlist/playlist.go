package playlist

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type PlayListI interface {
	Play()
	Pause()
	AddSong(newSong Song)
	Next()
	Prev()
}

type PlayList struct {
	list           List
	mu             sync.Mutex
	currNode       *Node
	playedDuration time.Duration
	RemoteChan     chan int
	paused         bool
	playing        bool
}

func (p *PlayList) StartToPlay() {
	for {
		if p.currNode == nil && p.list.start == nil {
			p.playing = false
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if p.currNode == nil {
			p.currNode = p.list.start
			p.playedDuration = time.Duration(0)
		}

		for p.playedDuration < p.currNode.Song.Duration {
			switch <-p.RemoteChan {
			case PLAY:
				log.Println(">>> play")
				p.paused = false
			case PAUSE:
				log.Println(">>> pause")
				p.paused = true
			case NEXT:
				log.Println(">>> next")
				p.next()
			case PREV:
				log.Println(">>> prev")
			}
			if p.paused {
				continue
			}
			currTime := parseDuration(p.playedDuration)
			maxTime := parseDuration(p.currNode.Song.Duration)
			fmt.Printf("(%s/%s): %s (play)\n", currTime, maxTime, p.currNode.Song.Name)
			time.Sleep(50 * time.Millisecond)
			p.playedDuration += 50 * time.Millisecond
		}

		if p.next() != nil {
			p.Pause()
		}
	}
}

func (p *PlayList) Play() {
	p.RemoteChan <- PLAY
}

func (p *PlayList) Pause() {
	p.RemoteChan <- PAUSE
}

func (p *PlayList) AddSong(newSong Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	newNode := Node{Song: newSong}
	p.list.Append(&newNode)
	fmt.Printf("song %s added to playlist\n", newSong.Name)
}

func (p *PlayList) Next() {
	p.RemoteChan <- NEXT
}

func (p *PlayList) next() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playedDuration = 0

	if p.currNode.next == nil {
		fmt.Println("end of playlist")
		p.currNode = nil
		return errors.New("EOP")
	}
	p.currNode = p.currNode.next
	return nil
}

func (p *PlayList) Prev() {
	p.RemoteChan <- PREV
}

func (p *PlayList) prev() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playedDuration = 0

	if p.currNode.prev == nil {
		fmt.Println("it is first song, restarted")
		return
	}
	p.currNode = p.currNode.prev
}
