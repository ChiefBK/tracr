package actions

type ActionQueue struct {
	Queue []*ExternalAction
}

func NewActionQueue() *ActionQueue {
	var queue []*ExternalAction
	return &ActionQueue{queue}
}

func (aq *ActionQueue) Push(action *ExternalAction) {
	aq.Queue = append(aq.Queue, action)
}

func (aq *ActionQueue) Dequeue() *ExternalAction {
	if len(aq.Queue) < 1 {
		return nil
	}

	action := aq.Queue[0]
	aq.Queue = aq.Queue[1:]

	return action
}

func (self ActionQueue) Length() int {
	return len(self.Queue)
}
