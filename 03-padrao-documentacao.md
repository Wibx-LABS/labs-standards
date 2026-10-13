# PADRÃO DE DOCUMENTAÇÃO LABS (LABS DOCUMENTATION STANDARD)

**Classificação:** Norma Canônica de Documentação de Sistema e Produto  
**Status:** ATIVO / CANÔNICO  
**Última Revisão:** Outubro de 2026  

---

Como o Labs documenta o que um sistema faz, em que estado está e o que uma mudança toca. Segue os níveis da Norma 01: todo projeto cumpre o Nível 0; o Nível 1 entra com o mesmo gatilho (código de produção mantido por TI ou voltado a usuários finais).

Base de mercado: Diátaxis (quatro tipos de documento), docs-as-code (documento no git, revisado por PR), o modelo de sistema do Backstage (componente, API, recurso e `dependsOn`, para análise de impacto) e os estágios de maturidade de feature gates (status com critério objetivo). A diferença do Labs: documento que descreve código é **conferido contra o código no CI**. Catálogo escrito à mão diverge em semanas; catálogo conferido não.

---

## 1. Princípios

1. **Documento descreve o código como ele é.** Sem histórico de mudança, etapa ou pendência. Histórico é o Git; plano mora no plano.
2. **Um tipo por documento (Diátaxis).** Referência, guia de tarefa (runbook), explicação (por quê, decisões) e tutorial não se misturam. Referência que vira tutorial no meio é documento quebrado.
3. **Fonte de máquina, página de gente.** Quando um documento lista coisas que existem no código (rotas, tabelas, flags, funcionalidades), a fonte é um arquivo estruturado (YAML/JSON) e a página humana é gerada dele. A página gerada nunca é editada à mão.
4. **Conferido, não confiado.** Toda referência a código no catálogo (rota, tabela, símbolo) é verificada por script. Referência que não existe reprova o PR.
5. **Citação por símbolo.** Cita-se `repo · pacote.Simbolo`, nunca `arquivo:linha`, que apodrece no primeiro refactor. Exceção: migração aplicada, que é imutável.
6. **Mesmo PR.** Mudou o comportamento, muda o documento no mesmo PR que muda o código.

---

## 2. Níveis

| Nível | Aplica-se a | Seção |
| :--- | :--- | :--- |
| **0** | Todo projeto. | §3 |
| **1** | Mesmo gatilho do Nível 1 da Norma 01. | §4 |

---

## 3. Nível 0 — Todo Projeto

- **`README.md`**: o que é, como rodar, como testar. Referência, não diário.
- **`CLAUDE.md` / `AGENTS.md`**: regras e armadilhas do repo para agente e humano. Sem histórico.
- **Contrato gerado.** API HTTP publicada tem contrato gerado do código (OpenAPI), com CI que falha se o gerado divergir.
- **Decisão que venceu alternativa** registra o que foi descartado e por quê, uma linha basta. Decisão sem alternativa registrada é relitigada.

---

## 4. Nível 1 — Catálogo de Funcionalidades

Todo sistema de Nível 1 mantém um **catálogo de funcionalidades de backend**: o mapa de baixo para cima do que o sistema faz, em que estado está, por onde passa e do que depende. Front não entra; o catálogo começa na borda de API e desce.

### 4.1. Unidade

Funcionalidade é uma **capacidade** que o produto entrega ("aprovar conteúdo", "gerar dossiê"), não um endpoint. As rotas ficam dentro dela. Toda rota do contrato pertence a alguma funcionalidade, exceto as de infraestrutura listadas no script (`/healthz`).

### 4.2. Campos

| campo | obrigatório | o quê |
| :--- | :--- | :--- |
| `id`, `nome` | sim | identificador kebab-case e nome de produto |
| `dominio` | sim | capacidade de negócio de nível 1 (lista fechada por sistema) |
| `o_que_faz` | sim | visão de produto, uma a três frases |
| `status` | sim | eixo backend (§4.3) |
| `nota` | se não `ativo` | o que falta ou o que arma |
| `rotas` | não | rotas do contrato da borda; vazio = não exposto |
| `disparos` | não | entradas fora da borda: worker, CLI, agente, evento |
| `modulos` | sim | repos/componentes atravessados, em ordem, com símbolos |
| `tabelas` | não | `schema.tabela` lidas ou escritas |
| `config`, `externos`, `testes` | não | flags, serviços de terceiros, testes do fluxo |

### 4.3. Status: dois eixos

Um status só mente: "pronto" no backend e invisível na API é o caso mais comum e o que mais confunde.

- **Backend:** `ativo` (chamador real e ligado na config padrão), `inerte` (pronto e desligado de propósito, com registro do que o arma), `parcial` (falta peça; `nota` diz qual).
- **Borda:** exposto ou não, derivado de `rotas`. Não se escreve à mão.

`ativo` descreve o código, não o ambiente. Saber se roda em QA ou produção é verificação em produção (Norma 01, §3.8), não catálogo.

### 4.4. Impacto

A página gerada mostra, por funcionalidade, **módulos · tabelas · serviços externos**. É a resposta para "vou mexer nesta funcionalidade, o que preciso olhar".

### 4.5. Verificação

Script versionado no repo de documentação, com dois modos:

- **CI do repo de docs:** schema válido e página gerada em dia.
- **Workspace com todos os repos:** rota existe no contrato, toda rota do contrato tem dono, tabela tem `CREATE TABLE` numa migração, símbolo existe no repo citado, `inerte` está no registro de código inerte.

Implementação de referência: `wibx-maestro/docs`, `funcionalidades/` e `tools/funcionalidades.py`.

---

## 5. Checklist de Conformidade

**Nível 0**
- [ ] README e CLAUDE.md descrevem o código atual, sem histórico.
- [ ] Contrato de API gerado do código, com CI de divergência.
- [ ] Citações por símbolo.

**Nível 1** — além do Nível 0:
- [ ] Catálogo de funcionalidades com fonte estruturada e página gerada.
- [ ] Toda rota do contrato pertence a uma funcionalidade.
- [ ] Verificação contra o código passando.
- [ ] Mudança de comportamento muda o catálogo no mesmo PR.
