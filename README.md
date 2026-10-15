# WibX Labs — Normas e Padrões Canônicos de Engenharia

Este repositório é a fonte única da verdade para arquitetura de software, DevOps, segurança e uso de inteligência artificial no **WibX Labs**.

---

## Documentos Canônicos

| Documento | Escopo & Descrição |
| :--- | :--- |
| **[00-charter-labs.md](./00-charter-labs.md)** | **Carta de Princípios:** Missão, governança e valores fundamentais do departamento interno de Labs. |
| **[01-padrao-engenharia-devops.md](./01-padrao-engenharia-devops.md)** | **Norma Geral de Engenharia & DevOps:** Níveis por gatilho (0 todo projeto · 1 produção TI · 2 por evidência), imutabilidade via digests `@sha256`, um runtime Go ou TypeScript por projeto, fila PostgreSQL e templates de Dockerfile. |
| **[02-padrao-engenharia-ia.md](./02-padrao-engenharia-ia.md)** | **Norma de Contexto & IA:** Planejamento e execução separados (Superpowers), Juiz Compilador, compressão RTK, Doutrina Ponytail (menor diff funcional) e uso de grafos semânticos. |
| **[03-padrao-documentacao.md](./03-padrao-documentacao.md)** | **Norma de Documentação:** Diátaxis, docs-as-code, citação por símbolo e, no Nível 1, catálogo de funcionalidades de backend com status em dois eixos (backend · borda), impacto por funcionalidade e verificação contra o código. |
| **[04-padrao-dados.md](./04-padrao-dados.md)** | **Norma de Dados (rascunho):** PostgreSQL — PK `uuid` v7 (v4 em identificador externo sensível), tipos e nomes, migrações só para frente e imutáveis sem travar tráfego, papéis mínimos, RLS com `FORCE` no Nível 1, backup com restauração ensaiada e dado pessoal declarado e expurgável. |

---

## Utilitários e Templates

* **[`templates/`](./templates/):** Dockerfiles de referência para Go e Node/TypeScript (Norma 01, §3.3).

* **[`bin/labs-ai-doctor`](./bin/labs-ai-doctor):** Healthcheck automático para validar se a estação local de trabalho está em conformidade com os hooks do RTK, plugins do Antigravity e configurações seguras de sessão.

```bash
# Diagnosticar ambiente local:
./bin/labs-ai-doctor
```

---

## Titularidade e Licenciamento
Propriedade intelectual exclusiva da **WIBX Company**. Todo o código e documentação são de uso interno restrito.
