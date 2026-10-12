# PADRÃO DE ENGENHARIA DE CONTEXTO E USO DE IA (LABS AI & CONTEXT STANDARD)

**Classificação:** Norma Canônica de Engenharia de Contexto e Automação Agêntica  
**Status:** ATIVO / CANÔNICO  
**Última Revisão:** Outubro de 2026  

---

Regras para todo desenvolvimento assistido ou conduzido por LLMs e agentes no Labs. Objetivo: código fiel ao design em sessões longas e cota de tokens que dura a semana. O porquê de cada regra mora na documentação interna de contexto, não aqui.

---

## 1. Princípios

1. **Planejar e executar em contextos separados.** A sessão que define o plano não escreve o código. O plano vai para disco; a implementação roda em subagentes de contexto limpo (§2).
2. **Juiz Compilador.** O gate de aceitação é o código de saída de compilador, linter e testes locais (`go test`, `tsc --noEmit`). `exit 0` aceita; o erro estrito dita a correção. Revisão por LLM (`/ponytail-review`, `/code-review`) é permitida como sinal adicional — nunca substitui nem dispensa o gate.
3. **Prefixo estável.** Regras do projeto (`CLAUDE.md`, `AGENTS.md`) e o plano ficam no topo do contexto e não são editados no meio da tarefa. Prefixo estável mantém o cache de prompt válido.
4. **Contexto contaminado não se conserta conversando.** Sessão em loop de correções remendadas é descartada e reiniciada a partir do plano em disco.
5. **Escopo cercado.** Em monorepo ou base legada, o agente trabalha em `git worktree` com `git sparse-checkout` restrito aos diretórios da tarefa. Navegação por símbolo via LSP antes de leitura de arquivo inteiro.
6. **Menor diff funcional.** Rege toda geração de código. Fonte única: plugin Ponytail (escada de decisão, regras e comentários `ponytail:`). Esta norma não duplica o conteúdo.

---

## 2. Fluxo de Trabalho: Coordenador e Subagentes (Superpowers)

Tarefa de um passo: sessão única, sem cerimônia. Tarefa de múltiplos passos:

1. **Sessão coordenadora.** Planeja com Superpowers (`brainstorming`, `writing-plans`) e grava o plano em disco. Não edita código; registra apenas plano e status.
2. **Execução por `subagent-driven-development`.** Cada tarefa do plano vai para um subagente isolado, que implementa, valida no Juiz Compilador (§1.2), commita e encerra. O coordenador recebe o recibo, marca a tarefa e despacha a seguinte.
3. **Trava de modelo.** Coordenador pode rodar no Opus. Subagentes travados no Sonnet via `CLAUDE_CODE_SUBAGENT_MODEL=sonnet` (alias, não ID de versão).

---

## 3. Grafos Semânticos e Mapas de Arquitetura

Ferramentas de grafo são pré-processamento local, não interlocutores. **Só o recorte da consulta entra no prompt; dump bruto nunca.**

- **Joern (Code Property Graph):** consulta CPGQL local para taint analysis, callers profundos e impacto de refatoração. Entra no prompt apenas o slice resultante.
- **Graphify:** `graphify-out/graph.json` é consultado via `jq` ou script local para extrair os acoplamentos de 1º grau do módulo. O JSON inteiro nunca entra no prompt.
- **Mapas de arquitetura** (diagramas, fronteiras de domínio) moram em `docs/` e são lidos sob demanda, nunca carregados no boot de toda sessão.

Ferramenta nova entra nesta seção quando um projeto a usar e o resultado for registrado.

---

## 4. Ferramentas Canônicas

| Ferramenta | Papel | Configuração |
| :--- | :--- | :--- |
| **RTK** | Compressão de saída de terminal e leitura de arquivos. | Hook `PreToolUse` no Claude Code (`rtk hook claude`); plugin global no Antigravity (`~/.gemini/config/plugins/rtk/`). |
| **Ponytail** | Menor diff funcional (§1.6). | Plugin no Claude Code; regras no Antigravity/Gemini. |
| **Superpowers** | Fluxo coordenador-subagente (§2). | Plugin no Claude Code. |
| **gopls-lsp** | Navegação por símbolo (§1.5). | Plugin no Claude Code. |

**Métrica única:** `rtk gain`, alvo ≥ 70% de economia média. Abaixo disso, `rtk discover` aponta os comandos sem filtro.

---

## 5. Setup da Estação e `bin/labs-ai-doctor`

1. **RTK:**
   ```bash
   cargo install rtk || brew install rtk-ai/tap/rtk
   rtk --version
   ```
2. **Claude Code (`~/.claude/settings.json`):**
   * Hook `PreToolUse` apontando para `rtk hook claude`.
   * `env.CLAUDE_CODE_SUBAGENT_MODEL: "sonnet"`.
   * Plugins habilitados: `superpowers`, `ponytail`, `gopls-lsp`.
   * `env.CLAUDE_CODE_FORCE_SESSION_PERSISTENCE` ausente ou `"0"`.
3. **Antigravity / Gemini:** plugin RTK em `~/.gemini/config/plugins/rtk/` com hook `PreToolUse` de `run_command`.

Verificação, por estação: [`bin/labs-ai-doctor`](./bin/labs-ai-doctor). Sai com `exit 1` só em falha crítica.

```bash
~/LABS/labs-standards/bin/labs-ai-doctor
```

---

## 6. Checklist de Conformidade (Definition of Done de Contexto)

Antes de mergear PR construído com agentes:

- [ ] **Juiz Compilador:** build e testes locais com `exit 0`, não presunção do modelo.
- [ ] **Ponytail:** menor diff funcional, sem dependência ou abstração não solicitada.
- [ ] **Separação:** tarefa de múltiplos passos executada a partir de plano em disco, fora da sessão que o escreveu.
- [ ] **Esterilidade:** nenhuma credencial, segredo ou dump massivo de dados passou pelo contexto da IA.
