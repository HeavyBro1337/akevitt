package basic

import (
	"github.com/IvanKorchmit/akevitt"
)

// Out-of-character chat command
func OocCmd(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	return engine.Message("ooc", command, session.Account.Username, session)
}
