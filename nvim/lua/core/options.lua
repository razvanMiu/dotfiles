-- Line numbers
vim.opt.nu = true
vim.opt.relativenumber = true
vim.opt.cursorline = true
vim.opt.colorcolumn = "80"
vim.opt.signcolumn = "yes"  -- always show, prevents layout shift

-- Indentation
vim.opt.tabstop = 4
vim.opt.softtabstop = 4
vim.opt.shiftwidth = 4
vim.opt.expandtab = true
vim.opt.autoindent = true
vim.opt.smartindent = true
vim.opt.breakindent = true  -- wrapped lines preserve indent

-- Wrapping & whitespace
vim.opt.wrap = true
vim.opt.list = true
vim.opt.listchars = { tab = "» ", trail = "·", nbsp = "␣" }

-- Search
vim.opt.incsearch = true
vim.opt.hlsearch = true
vim.opt.ignorecase = true
vim.opt.smartcase = true    -- case-sensitive when query has uppercase
vim.opt.inccommand = "split"  -- live preview of :s substitutions

-- Files & undo
vim.opt.swapfile = false
vim.opt.backup = false
vim.opt.undofile = true     -- persist undo history across sessions

-- Splits
vim.opt.splitright = true
vim.opt.splitbelow = true
vim.opt.splitkeep = "cursor"

-- Appearance
vim.opt.termguicolors = true
vim.opt.background = "dark"
vim.opt.scrolloff = 8
vim.opt.smoothscroll = true

-- Misc
vim.opt.backspace = { "start", "eol", "indent" }
vim.opt.isfname:append("@-@")
vim.opt.updatetime = 50     -- faster CursorHold events (used by LSP)
vim.opt.clipboard:append("unnamedplus")
vim.opt.mouse = "a"
vim.g.editorconfig = true
vim.g.netrw_banner = 0
