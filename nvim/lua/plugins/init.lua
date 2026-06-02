local plugin_dir = vim.fn.stdpath("config") .. "/lua/plugins"

for name, type in vim.fs.dir(plugin_dir) do
    if type == "file" and name ~= "init.lua" and name:match("%.lua$") then
        require("plugins." .. name:gsub("%.lua$", ""))
    end
end
