# Akevitt Quick Start Guide

Get your MUD running in minutes.

## Prerequisites

- Go 1.22+
- SSH client (built into Linux/macOS, use [PuTTY](https://www.putty.org/) on Windows)

## Step 1: Create Your Project

Use the CLI tool to scaffold a new project:

```bash
go run github.com/IvanKorchmit/akevitt/cmd/akevitt@latest init mygame
cd mygame
```

This creates:
```
mygame/
├── main.go          # Your game server
├── go.mod           # Go module
├── scripts/
│   └── look.lua     # Example command
└── keys/
    └── id_rsa       # Auto-generated SSH private key
```

## Step 2: Install Dependencies

```bash
go mod tidy
```

## Step 3: Start the Server

```bash
go run main.go
```

You should see:
```
Listening on :2222
```

## Step 4: Connect

Open a **new terminal** and connect via SSH:

```bash
ssh -p 2222 localhost
```

First-time connections will prompt to accept the server's SSH key—type `yes`.

## Step 5: Play

When you connect, you'll see a welcome screen. Try these commands:

| Command | What it does |
|---------|--------------|
| `look` | Examine your surroundings |
| `help` | List available commands |

Type commands and press Enter.

## Project Structure

```
mygame/
├── main.go          # Server entry point - customize this!
├── go.mod            # Go dependencies
├── scripts/          # Lua command scripts
│   └── look.lua      # The 'look' command implementation
└── keys/             # SSH keys (keep private!)
    ├── id_rsa        # Private key (server uses)
    └── id_rsa.pub    # Public key
```

## Customizing Your MUD

### Add Rooms

Edit `main.go` to create your world:

```go
lobby := akevitt.NewRoom("Lobby")
lobby.Description = "A quiet lobby with soft lighting."

tavern := akevitt.NewRoom("Tavern")
tavern.Description = "A cozy tavern filled with chatter."

// Connect rooms
akevitt.BindRoomsBidirectional(lobby, akevitt.Exit{}, tavern)

// Set spawn point
engine.UseSpawnRoom(lobby)
```

### Add Commands

Create a new `.lua` file in `scripts/`. Example `go.lua`:

```lua
local function go(session, args)
    local direction = args[1]
    if not direction then
        akevitt.sendMessage(session, "Go where? Usage: go <direction>")
        return
    end
    
    local current = akevitt.getPlayerRoom(session)
    local exits = akevitt.getRoomExits(current.name)
    
    for _, exit in ipairs(exits) do
        if exit.direction == direction then
            akevitt.setPlayerRoom(session, exit.target)
            akevitt.sendMessage(session, "You go " .. direction .. ".")
            return
        end
    end
    
    akevitt.sendMessage(session, "You can't go that way.")
end

return go
```

### Lua API Reference

| Function | Description |
|----------|-------------|
| `akevitt.sendMessage(session, text)` | Send text to player |
| `akevitt.getPlayerRoom(session)` | Get player's current room |
| `akevitt.setPlayerRoom(session, roomName)` | Move player to room |
| `akevitt.getRoomExits(roomName)` | List exits from a room |
| `akevitt.addRoom({name, description})` | Create a new room |
| `akevitt.addNPC({name, description, room})` | Create an NPC |
| `akevitt.getSessions()` | Get all connected players |
| `akevitt.setPlayerData(session, key, value)` | Store data on player |

## Server Configuration

In `main.go`:

```go
engine := akevitt.NewEngine().
    // Change port (default: 2222)
    UseBind(":2222").
    // Point to your SSH key
    UseKeyPath("keys/id_rsa").
    // Set your spawn room
    UseSpawnRoom(lobby).
    Finish()
```

## Rebuilding

After editing `main.go` or adding Lua scripts, restart the server.

## Troubleshooting

**Connection refused**: Server might not be running. Start it with `go run main.go`.

**SSH key error**: Ensure `keys/id_rsa` exists and has correct permissions (600).

**Commands not working**: Check that `scripts/` contains `.lua` files and your `main.go` includes the Lua plugin:

```go
AddPlugin(plugins.NewLuaCommandPlugin("scripts"))
```

## Next Steps

- Read `AGENTS.md` for architecture details
- Explore `example/main.go` for more patterns
- Check `scripts/` for example Lua commands

## Getting Help

SSH to your server and type `help` to see available commands.