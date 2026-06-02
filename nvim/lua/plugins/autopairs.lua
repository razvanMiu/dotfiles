vim.api.nvim_create_autocmd("InsertEnter", {
    group = vim.api.nvim_create_augroup("lazy-autopairs", { clear = true }),
    once = true,
    callback = function()
        vim.pack.add {
            { src = "https://github.com/hrsh7th/nvim-cmp",      name = "nvim-cmp" },
            { src = "https://github.com/windwp/nvim-autopairs", name = "nvim-autopairs" },
        }

        require("nvim-autopairs").setup({
            check_ts = true,
            ts_config = {
                lua  = { "string" },
                java = false,
            },
            disable_filetype = { "TelescopePrompt", "spectre_panel", "snacks_picker_input" },
            fast_wrap = {},
        })

        local cmp = require("cmp")
        local cmp_autopairs = require("nvim-autopairs.completion.cmp")
        cmp.event:on("confirm_done", cmp_autopairs.on_confirm_done())
    end,
})
