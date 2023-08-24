# Akevitt

*Akevitt* is a MUD engine designed to be as modular, but powerfully minimal, as possible. This means that there are no
"plugins" or "add-ons" that have narrow functionality. You get the infrastructure, information persistence, accounts,
and other essentials done for you, while the rest is up to you. Strongly inspired by
[Evennia](https://github.com/evennia/evennia). What you can do with Akevitt is only limited by your imagination and
technical skill in Go.

## Features

- Runs on an SSH server, meaning that, just like with Telnet, you can connect to an Akevitt game from almost any
device available in the world.
- Automatic database, meaning any game object (`GObject`) you create is tracked and survives reboots. You do not need
to think about the database, ever.
- Modularity. Use straight-forward API to create your objects, NPCs, rooms, or anything with, and place them all in
one file or organised across as many files as you need.

## How does Akevitt work?

The engine runs on a simple SSH server, which means that anybody with an SSH client installed on their computer can
connect to the game and play in their terminal. Data is then stored in a small database that requires no external
libraries or programs. It is all built into Akevitt. The server is entirely self-sufficient: nothing like SQLite is
required.

## What does Akevitt use to work?

The SSH server is implemented with [ssh](https://github.com/gliderlabs/ssh) by Glider Labs. Database is
[Bolt](https://github.com/boltdb/bolt), a simple non-relational key-value database, which is really all that is needed
for a MUD. UI elements are implemented with [tview](https://github.com/rivo/tview). Basically, whatever challenges
Akevitt meets, they are solved with the latest, mature libraries that the Go ecosystem has to offer. This leads to a
robust infrastructure with the worst bugs eliminated before a single line of Akevitt's source was written.

## Why this over Evennia?

Although Evennia is ambitious, and it certainly delivers in its rich feature set, it is quite bloated for what it is:
over 350 MB of source just to get a simple MUD going? We can do better. Additionally, Python is very difficult to
actually develop large projects in. `pip` is a nightmare. Maintainers of Python frequently introduce breaking changes
in what are supposed to be minor versions of Python 3. Python virtual environments offer mixed results, most often not
really solving the issue they are allegedly addressing. We can go on.

For end-users, it probably will not make a lick of difference whether their MUD runs on Akevitt or Evennia. However,
for developers like us it makes a huge difference, and Akevitt is here to bring joy to developing MUDs, while retaining
the ease of deploying complete games. If functional programming is more your speed, Akevitt is what you are looking
for.

## Licence & Attribution

### Programming and design

Ivan Korchmit (c) 2023

### Design and presentation

Maxwell Jensen (c) 2023

Licensed under European Union Public Licence 1.2. For more information, view LICENCE.
