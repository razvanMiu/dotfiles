# Neovim Config

Neovim 0.12+ config using `vim.pack` (built-in package manager). No lazy.nvim.

## Structure

```
init.lua                  -- entry: require("core") + require("plugins")
lua/
  core/
    init.lua              -- auto-loads all files in core/ except itself
    options.lua           -- vim.opt settings
    keymaps.lua           -- leader key + all core keymaps
    autocmds.lua          -- global autocmds (e.g. highlight yank)
  plugins/
    init.lua              -- auto-loads all files in plugins/ except itself
    colorscheme.lua
    autopairs.lua
    oil.lua
    ...
nvim-pack-lock.json       -- auto-managed lockfile, commit this
```

### Auto-discovery

Both `core/init.lua` and `plugins/init.lua` use `vim.fs.dir` to auto-load every
`.lua` file in their directory. Drop a new file in — it loads automatically. No
manual registration needed.

```lua
for name, type in vim.fs.dir(dir) do
  if type == "file" and name ~= "init.lua" and name:match("%.lua$") then
    require("prefix." .. name:gsub("%.lua$", ""))
  end
end
```

Note: load order is not guaranteed. If order matters (rare in plugins, sometimes
in core), prefix filenames: `01_options.lua`, `02_keymaps.lua`.

## vim.pack

`vim.pack` installs and manages plugins. Available functions:

- `vim.pack.add(specs)` — install plugins and add to runtimepath
- `vim.pack.update()` — update all plugins
- `vim.pack.del(name)` — remove a plugin
- `vim.pack.get()` — list installed plugins

There is **no lazy loading** built into `vim.pack`. Calling `vim.pack.add`
immediately adds the plugin to rtp and Neovim sources its `plugin/` files.

### Basic usage (eager)

```lua
vim.pack.add { { src = "https://github.com/author/plugin", name = "plugin" } }
require("plugin").setup({})
```

### Lazy loading pattern

Defer `vim.pack.add` inside an autocmd or keymap callback. Plugin only loads
when the trigger fires.

**Event-based:**
```lua
vim.api.nvim_create_autocmd("InsertEnter", {
    once = true,
    callback = function()
        vim.pack.add { { src = "https://github.com/author/plugin", name = "plugin" } }
        require("plugin").setup({})
    end,
})
```

**Keymap-based:**
```lua
local loaded = false
vim.keymap.set("n", "<leader>x", function()
    if not loaded then
        vim.pack.add { { src = "https://github.com/author/plugin", name = "plugin" } }
        require("plugin").setup({})
        loaded = true
    end
    -- run actual action
    require("plugin").do_thing()
end)
```

**Advanced: reusable lazy loader via `load` callback**

`vim.pack.add` accepts a second arg `{ load = fn }`. When provided, the callback
fires instead of default loading — use it to defer `packadd` to events/keymaps.
The `data` field on each spec passes arbitrary config through to the callback.

```lua
vim.pack.add({
    {
        src  = "https://github.com/author/plugin",
        name = "plugin",
        data = {
            event  = "BufReadPre",
            config = function() require("plugin").setup({}) end,
        },
    },
}, {
    load = function(plugin)
        local data = plugin.spec.data or {}
        if data.event then
            vim.api.nvim_create_autocmd(data.event, {
                once = true,
                callback = function()
                    vim.cmd.packadd(plugin.spec.name)
                    if data.config then data.config() end
                end,
            })
        end
    end,
})
```

## Keymaps

Core keymaps live in `lua/core/keymaps.lua`. Uses a local `map` helper:

```lua
local map = function(mode, lhs, rhs, desc)
    vim.keymap.set(mode, lhs, rhs, { noremap = true, silent = true, desc = desc })
end
```

Plugin files use `vim.keymap.set` directly — `map` is local to `keymaps.lua`.

## Adding a plugin

1. Create `lua/plugins/myplugin.lua`
2. Call `vim.pack.add` with the GitHub src
3. Call `require("myplugin").setup({})` — or wrap in autocmd for lazy loading
4. Restart nvim — plugin downloads automatically
