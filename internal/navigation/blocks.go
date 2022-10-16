package navigation

type StateMachine struct {
	states map[string]func(Model) Model
	transitions map[string]map[string]string

	currentStateName string
}

var IdentityEnterFunc = func(model Model) Model { return model }

func NewStateMachine() *StateMachine {
	sm := new(StateMachine)
	return sm
}

func (sm *StateMachine) SetState(name string, onEnter func(Model) Model) {
	if sm.states == nil {
		sm.states = map[string]func(Model)Model{}
		sm.currentStateName = name
	}
	sm.states[name] = onEnter
}

func (sm *StateMachine) SetInitialState(name string) {
	sm.SetState(name, IdentityEnterFunc)
	sm.currentStateName = name
}

func (sm *StateMachine) SetTransition(startStateName string, action string, endStateName string) {
	if sm.transitions[startStateName] == nil {
		sm.transitions[startStateName] = map[string]string{}
	}

	sm.transitions[startStateName][action] = endStateName
}

func (sm *StateMachine) Transition(action string, m Model) Model {
	transitions, exists := sm.transitions[sm.currentStateName]
	if !exists {
		return m
	}

	targetState, exists := transitions[action]
	if !exists {
		return m
	}

	sm.currentStateName = targetState
	return sm.states[sm.currentStateName](m)
}

func (sm *StateMachine) CurrentStateName() string {
	return sm.currentStateName
}