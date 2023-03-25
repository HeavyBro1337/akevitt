package akevitt

// In-game object that you can interact within the game.
type GameObject interface {
	*Object
	Create(engine *Akevitt, session *ActiveSession) error
}
