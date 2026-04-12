package bot

const (
	StateAwaitingName  = "awaiting_name"
	StateAwaitingPhone = "awaiting_phone"
)

type StateManager struct {
	states    map[int64]string
	tempNames map[int64]string
}

func NewStateManager() *StateManager {
	return &StateManager{
		states:    make(map[int64]string),
		tempNames: make(map[int64]string),
	}
}

func (s *StateManager) Get(telegramID int64) string {
	return s.states[telegramID]
}

func (s *StateManager) Set(telegramID int64, state string) {
	s.states[telegramID] = state
}

func (s *StateManager) Clear(telegramID int64) {
	delete(s.states, telegramID)
}

func (s *StateManager) SetTempName(telegramID int64, name string) {
	s.tempNames[telegramID] = name
}

func (s *StateManager) GetTempName(telegramID int64) string {
	name := s.tempNames[telegramID]
	delete(s.tempNames, telegramID)
	return name
}
