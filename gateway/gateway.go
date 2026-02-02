package gateway

type Gateway interface {
    Execute(Message string) error
}

type SessionMessage struct {
    SessionID string
    Message string
  }
type SessionRouter interface {
    Serve(s SessionMessage) error
}

type LaneQueue struct {
    Queue map[string]chan[string]
}

func (lq *LaneQueue) Loop() {
    for id, c := range lq.Queue {
        select {
            case msg := <- c:
                // Agent Runner
            default:
               // do nothing
        }
    }
}
