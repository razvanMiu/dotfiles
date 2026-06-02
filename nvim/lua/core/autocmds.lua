-- Flash yanked region briefly so you know what was copied
vim.api.nvim_create_autocmd("TextYankPost", {
  group = vim.api.nvim_create_augroup("highlight-yank", { clear = true }),
  callback = function()
    vim.hl.on_yank({ higroup = "IncSearch", timeout = 150 })
  end,
})
