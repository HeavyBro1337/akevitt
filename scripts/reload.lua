-- reload.lua
-- Reloads all Lua scripts (admin command)

local function reload(session, engine)
    -- Check if user is admin (simplified - would need proper admin check)
    local is_admin = session.account and session.account.username == "admin"
    
    if not is_admin then
        akevitt.sendMessage(session, "Only administrators can reload scripts.")
        return
    end
    
    -- Note: The actual reload is triggered via the akevitt.ReloadLua() Go function
    -- This script would need to signal the engine to reload
    -- For now, this is a placeholder
    akevitt.sendMessage(session, "Use 'akevitt reload' from the server console to reload scripts.")
end

return reload
