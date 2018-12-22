package core

type Player struct {
	platform  int
	handler   func(...interface{})
	params []interface{}
}

func (Player *Player) NewClient(player int) *Player {
	Player.platform = player
	return Player
}
func (Player *Player) RecHandleFunc(hFuc handlerFunc,params ...interface{}) {
	Player.handler = hFuc
}
func (Player *Player) Connect() {
	Player.handler(Player.params)
}

type handlerFunc func(...interface{})

