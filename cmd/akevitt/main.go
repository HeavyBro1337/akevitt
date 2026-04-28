package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: akevitt <command>")
		fmt.Println("Commands:")
		fmt.Println("  init <directory>  Initialize a new Akevitt project")
		fmt.Println("  generate-keys    Generate SSH keys for the server")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Usage: akevitt init <directory>")
			os.Exit(1)
		}
		initProject(os.Args[2])
	case "generate-keys":
		generateSSHKeys()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func initProject(dir string) {
	fmt.Printf("Initializing Akevitt project in %s...\n", dir)

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating scripts directory: %v\n", err)
		os.Exit(1)
	}

	keysDir := filepath.Join(dir, "keys")
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating keys directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generating SSH keys...")
	if err := generateSSHKeyPair(filepath.Join(keysDir, "id_rsa")); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating SSH keys: %v\n", err)
		os.Exit(1)
	}

	mainGo := filepath.Join(dir, "main.go")
	mainContent := `package main

import (
	"log"

	"github.com/IvanKorchmit/akevitt/engine"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/rivo/tview"
)

func main() {
	room := engine.NewRoom("Lobby")
	room.Description = "A quiet lobby with soft lighting."

	tavern := engine.NewRoom("Tavern")
	tavern.Description = "A cozy tavern filled with warmth and chatter."

	engine.BindRoomsBidirectional(room, engine.Exit{}, tavern)

	eng := engine.NewEngine().
		AddPlugin(plugins.DefaultPlugins()...).
		AddPlugin(plugins.NewAccountPlugin()).
		AddPlugin(plugins.NewBoltPlugin[*engine.Account]("database.db")).
		AddPlugin(plugins.NewLuaCommandPlugin("scripts")).
		UseSpawnRoom(room).
		UseRootUI(RootUI).
		UseBind(":2222").
		UseKeyPath("keys/id_rsa").
		Finish()

	log.Fatal(eng.Run())
}

func RootUI(eng *engine.Akevitt, session *engine.ActiveSession) tview.Primitive {
	form := tview.NewForm()
	form.AddTextView("welcome", "Welcome to Akevitt!\n\nType 'look' to examine your surroundings.\nType 'help' for available commands.", 40, 5, true, false)

	session.RoomID = eng.GetSpawnRoom().GUID

	return form
}
`

	if err := os.WriteFile(mainGo, []byte(mainContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing main.go: %v\n", err)
		os.Exit(1)
	}

	lookLua := `-- look.lua
-- Displays information about the current room

local function look(session)
    local room = akevitt.getPlayerRoom(session)
    
    if not room then
        akevitt.sendMessage(session, "You are not in any room.")
        return
    end
    
    akevitt.sendMessage(session, "=== " .. room.name .. " ===")
    akevitt.sendMessage(session, room.description or "A mysterious place.")
    
    local exits = akevitt.getRoomExits(room.name)
    if exits and #exits > 0 then
        local exit_names = {}
        for _, exit in ipairs(exits) do
            table.insert(exit_names, exit.target)
        end
        akevitt.sendMessage(session, "Exits: " .. table.concat(exit_names, ", "))
    else
        akevitt.sendMessage(session, "There are no obvious exits.")
    end
    
    if room.objects and #room.objects > 0 then
        akevitt.sendMessage(session, "You see:")
        for _, obj in ipairs(room.objects) do
            akevitt.sendMessage(session, "  - " .. obj.name)
        end
    end
end

return look
`

	if err := os.WriteFile(filepath.Join(scriptsDir, "look.lua"), []byte(lookLua), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing look.lua: %v\n", err)
		os.Exit(1)
	}

	goMod := `module mygame

go 1.22

require github.com/IvanKorchmit/akevitt v0.4.0
`

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project initialized successfully!")
	fmt.Println("")
	fmt.Println("To get started:")
	fmt.Printf("  cd %s\n", dir)
	fmt.Println("  GOSUMDB=off go mod tidy")
	fmt.Println("  go run main.go")
	fmt.Println("")
	fmt.Println("Then connect with: ssh -p 2222 localhost")
}

func generateSSHKeys() {
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
		os.Exit(1)
	}

	keyPath := filepath.Join(usr.HomeDir, ".ssh", "id_rsa_akevitt")
	if err := generateSSHKeyPair(keyPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating SSH keys: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SSH keys generated at: %s\n", keyPath)
}

func generateSSHKeyPair(path string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if err := os.WriteFile(path, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to generate public key: %w", err)
	}

	if err := os.WriteFile(path+".pub", ssh.MarshalAuthorizedKey(pubKey), 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}