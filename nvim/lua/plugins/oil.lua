vim.pack.add {
    { src = "https://github.com/stevearc/oil.nvim",             name = "oil" },
    { src = "https://github.com/nvim-tree/nvim-web-devicons",   name = "nvim-web-devicons" },
}

require("oil").setup({
    default_file_explorer = true,
    columns = {},
    keymaps = {
        ["<C-h>"] = false,
        ["<C-l>"] = false,
        ["<C-c>"] = false, -- <C-c> is esc, don't let oil capture it
        ["<C-r>"] = "actions.refresh",
        ["<M-h>"] = "actions.select_split",
        ["q"]     = "actions.close",
    },
    delete_to_trash = true,
    view_options = { show_hidden = true },
    skip_confirm_for_simple_edits = true,
})

vim.keymap.set("n", "-",         "<CMD>Oil<CR>",                { desc = "Open parent directory" })
vim.keymap.set("n", "<leader>-", require("oil").toggle_float,   { desc = "Toggle oil float" })

vim.api.nvim_create_autocmd("FileType", {
    pattern = "oil",
    callback = function() vim.opt_local.cursorline = true end,
})
