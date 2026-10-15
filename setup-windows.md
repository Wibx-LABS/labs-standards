# Setup da Estação no Windows

Aplica a [Norma 02 §5](./02-padrao-engenharia-ia.md) numa máquina Windows. Escrito para o Claude Code executar: o colaborador abre uma sessão **Local** na aba Code do Claude Desktop (ou o `claude` no terminal) e pede *"aplique o setup-windows.md de https://github.com/Wibx-LABS/labs-standards"*.

**Aceite:** `bin/labs-ai-doctor`, rodado no Git Bash, termina sem `[FAIL]`, sem `[ALERT]` e sem `[WARN]` do Claude Code. Os dois avisos de Antigravity/Gemini só contam para quem usa Antigravity.

---

## 1. O que muda no Windows

| Ponto | Efeito na norma | Fonte |
| :--- | :--- | :--- |
| Sessão WSL do Desktop não carrega plugins. | Ponytail, Superpowers e gopls-lsp ficam de fora. | Doc Claude Code, *Claude Code Desktop in WSL* |
| Sem Git for Windows, o shell do Claude é o PowerShell. | O hook do RTK (matcher `Bash`) não dispara. | Doc Claude Code, *Set up on Windows* |
| Com Git for Windows, a ferramenta PowerShell segue ligada ao lado do Bash. | Comando que vai pelo PowerShell escapa do RTK. | Idem |
| O Desktop herda as variáveis de ambiente do usuário e do sistema, e não lê o profile do PowerShell. | `rtk` e `gopls` precisam estar no `Path` do Windows. | Doc Claude Code, *Desktop: local sessions* |
| O Git for Windows troca LF por CRLF no checkout. | O doctor quebra na linha 2 (`set: pipefail: invalid option name`). | `.gitattributes` deste repo força LF |

---

## 2. Passos

Faça na ordem. Passo marcado **[humano]** é clique ou janela que o Claude não opera: pare, diga exatamente o que a pessoa deve fazer e espere a confirmação.

1. **Sessão Local, nunca WSL.** Se a sessão atual roda em WSL, pare e peça para abrir uma sessão Local.
2. **Cópia deste repo.** Baixe `https://github.com/Wibx-LABS/labs-standards/archive/refs/heads/main.zip` e extraia em `%USERPROFILE%`; renomeie a pasta `labs-standards-main` para `labs-standards`. O zip não depende de Git e mantém o LF dos scripts.
3. **Git for Windows.** `git --version`. Se faltar: `winget install --id Git.Git -e`.
4. **RTK.** `rtk --version`. Se faltar: `winget install rtk-ai.rtk` (o winget põe no `Path`). Confirme que é o RTK certo: `rtk gain` precisa responder, porque existe outro projeto chamado `rtk`.
5. **[humano] Reiniciar o Claude Desktop** (fechar de vez, inclusive da bandeja) para o Git Bash e o `rtk` entrarem no `Path` da sessão. Retome do passo 6 na sessão nova.
6. **`%USERPROFILE%\.claude\settings.json`.** Copie o arquivo para `settings.json.bak-labs` antes de mexer. Junte as chaves abaixo ao que já existe: nunca sobrescreva o arquivo, nunca apague chave da pessoa. Se `hooks.PreToolUse` já tiver uma entrada com matcher `Bash`, acrescente o hook nela.

   ```json
   {
     "env": {
       "CLAUDE_CODE_SUBAGENT_MODEL": "sonnet",
       "CLAUDE_CODE_USE_POWERSHELL_TOOL": "0"
     },
     "hooks": {
       "PreToolUse": [
         { "matcher": "Bash", "hooks": [{ "type": "command", "command": "rtk hook claude" }] }
       ]
     }
   }
   ```

   `CLAUDE_CODE_USE_POWERSHELL_TOOL=0` desliga a ferramenta PowerShell: todo comando passa pelo Bash, onde o hook do RTK atua. Se `env.CLAUDE_CODE_FORCE_SESSION_PERSISTENCE` estiver em `"1"`, troque por `"0"` e avise a pessoa.
7. **Plugins.**
   - Se `claude --version` responde (CLI instalada), rode:

     ```bash
     claude plugin marketplace add Wibx-LABS/wibx-skills
     claude plugin install superpowers@claude-plugins-official --scope user
     claude plugin install gopls-lsp@claude-plugins-official --scope user
     claude plugin install ponytail@wibx-skills --scope user
     ```

   - Só Desktop, sem CLI: **[humano]** botão **+** ao lado do prompt, **Plugins**, **Add plugin**. Adicione o marketplace `Wibx-LABS/wibx-skills` e instale `superpowers` e `gopls-lsp` (marketplace oficial) e `ponytail` (`wibx-skills`).
8. **Go e gopls** (só quem trabalha em projeto Go). `go version`; se faltar: `winget install --id GoLang.Go -e`. Depois `go install golang.org/x/tools/gopls@latest`. Se `command -v gopls` não achar, **[humano]** adicione `%USERPROFILE%\go\bin` ao `Path` da conta (Iniciar, *Editar as variáveis de ambiente para sua conta*). Não use `setx` em `Path`: ele corta o valor em 1024 caracteres.
9. **[humano] Reiniciar o Claude Desktop** de novo. Hook, `env` e plugins só valem em sessão nova.

---

## 3. Verificação

Na sessão nova, no Git Bash:

```bash
bash ~/labs-standards/bin/labs-ai-doctor
echo '{"tool_name":"Bash","tool_input":{"command":"ls"}}' | rtk hook claude
```

O doctor bate o aceite do topo. O hook devolve `"updatedInput":{"command":"rtk ls"}`: o binário reescreve o que o Claude manda pelo Bash.

---

## 4. Rollback

1. `settings.json.bak-labs` volta a ser `settings.json`.
2. Plugins: `claude plugin uninstall <nome>` ou **+**, **Plugins**, **Manage plugins**.
3. `winget uninstall rtk-ai.rtk`; Git e Go só se foram instalados por este setup.
4. Reiniciar o Claude Desktop.
