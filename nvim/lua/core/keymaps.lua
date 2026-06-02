vim.g.mapleader = " "
vim.g.maplocalleader = " "

local map = function(mode, lhs, rhs, desc)
    vim.keymap.set(mode, lhs, rhs, { noremap = true, silent = true, desc = desc })
end

-- Navigation: keep cursor centered while jumping
map("n", "<C-d>", "<C-d>zz", "Scroll down (centered)")
map("n", "<C-u>", "<C-u>zz", "Scroll up (centered)")
map("n", "n",     "nzzzv",   "Next search result (centered)")
map("n", "N",     "Nzzzv",   "Prev search result (centered)")

-- Editing: move selected lines up/down
map("v", "J", ":m '>+1<CR>gv=gv", "Move selection down")
map("v", "K", ":m '<-2<CR>gv=gv", "Move selection up")

-- Editing: keep indent after visual shift
map("v", "<", "<gv", "Indent left (stay in visual)")
map("v", ">", ">gv", "Indent right (stay in visual)")

-- Clipboard: delete/paste without clobbering yank register
map("x",            "<leader>p",    [["_dP]],   "Paste over selection (preserve yank)")
map("v",            "p",            '"_dp',     "Paste in visual (preserve yank)")
map({ "n", "v" },   "<leader>d",    [["_d]],    "Delete to void register")
-- x in normal mode goes to void so it doesn't pollute the yank register
map("n", "x", '"_x', "Delete char (void register)")

-- Search: clear highlight
map("i", "<C-c>", "<Esc>",      "Escape insert mode")
map("n", "<C-c>", ":nohl<CR>",  "Clear search highlight")

-- LSP
map("n", "<leader>f", vim.lsp.buf.format, "Format buffer (LSP)")

-- Refactor: rename word under cursor globally
map("n", "<leader>s", [[:%s/\<<C-r><C-w>\>/<C-r><C-w>/gI<Left><Left><Left>]], "Rename word globally")

-- File
map("n", "<leader>x", "<cmd>!chmod +x %<CR>", "Make file executable")
map("n", "<leader>fp", function()
    local path = vim.fn.expand("%:~")
    vim.fn.setreg("+", path)
    print("Copied: " .. path)
end, "Copy file path to clipboard")

-- Tabs
map("n", "<leader>to", "<cmd>tabnew<CR>",   "New tab")
map("n", "<leader>tx", "<cmd>tabclose<CR>", "Close tab")
map("n", "<leader>tn", "<cmd>tabn<CR>",     "Next tab")
map("n", "<leader>tp", "<cmd>tabp<CR>",     "Prev tab")
map("n", "<leader>tf", "<cmd>tabnew %<CR>", "Open buffer in new tab")

-- Splits
map("n", "<leader>sv", "<C-w>v",          "Split vertical")
map("n", "<leader>sh", "<C-w>s",          "Split horizontal")
map("n", "<leader>se", "<C-w>=",          "Equalize splits")
map("n", "<leader>sx", "<cmd>close<CR>",  "Close split")
