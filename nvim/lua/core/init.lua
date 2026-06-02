local core_dir = vim.fn.stdpath("config") .. "/lua/core"

for name, type in vim.fs.dir(core_dir) do
    if type == "file" and name ~= "init.lua" and name:match("%.lua$") then
        require("core." .. name:gsub("%.lua$", ""))
    end
end
