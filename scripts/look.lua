-- look.lua
-- Displays information about the current room

local function look(session)
    local room = akevitt.getPlayerRoom(session)
    
    if not room then
        akevitt.sendMessage(session, "You are not in any room.")
        return
    end
    
    akevitt.sendMessage(session, "=== " .. room.name .. " ===")
    akevitt.sendMessage(session, room.description or "A mysterious place.")
    
    -- Show exits
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
    
    -- Show objects (NPCs, items)
    if room.objects and #room.objects > 0 then
        akevitt.sendMessage(session, "You see:")
        for _, obj in ipairs(room.objects) do
            akevitt.sendMessage(session, "  - " .. obj.name)
        end
    end
end

return look
