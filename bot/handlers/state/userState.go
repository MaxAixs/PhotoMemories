package state

import "sync"

type Manager struct {
	mu     sync.RWMutex
	states map[int64]string
}

func NewStateManager() *Manager {
	return &Manager{
		states: make(map[int64]string),
	}
}

const (
	Default      = "default"
	AwaitPic     = "await_pictures"
	AwaitSaveTag = "await_save_tag"
	AwaitDelTag  = "await_del_tag"
	AwaitGetTag  = "await_get_tag"
	AwaitMyTags  = "await_my_tags"
)

func (b *Manager) SetState(userID int64, state string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.states[userID] = state
}

func (b *Manager) GetUserState(userID int64) string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.states[userID]
}

func (b *Manager) IsUserInState(userID int64, state string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.states[userID] == state
}
