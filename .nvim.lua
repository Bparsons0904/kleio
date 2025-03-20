-- client/.nvim.lua
local util = require("lspconfig.util")

-- Tell Neovim to treat this directory as a project root for ESLint
vim.g.root_dir = util.path.dirname(vim.fn.expand("%:p"))
vim.g.eslint_root_dir = vim.g.root_dir

-- Restart ESLint server to pick up the new root directory
local clients = vim.lsp.get_active_clients({ name = "eslint" })
if #clients > 0 then
  vim.cmd("LspRestart eslint")
end
