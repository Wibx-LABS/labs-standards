# PADRÃO DE PROJETO LABS (LABS ARCHITECTURAL STANDARD)

**Classificação:** Norma Canônica de Engenharia, DevOps e Arquitetura de Software  
**Status:** ATIVO / CANÔNICO  
**Última Revisão:** Outubro de 2026  

---

Como o Labs estrutura, empacota, promove e mantém software. A norma é dividida em **níveis**: todo projeto cumpre o Nível 0; os níveis seguintes entram por gatilho, nunca por expectativa. É a regra da Carta §4.7 aplicada à infraestrutura — prever uma necessidade é atenção, não construção. Um artefato entregue a um departamento precisa ser mantido pelo dono sem o Labs (Carta §6.7); cada nível acima do necessário é custo que o dono herda.

---

## 1. Princípios

1. **Imutabilidade.** Nenhum ambiente roda tag móvel (`:latest`, `:main`, `:v1`). Toda imagem em manifesto de execução aponta para `image@sha256:...`. O que foi testado é o que vai para o ar.
2. **Nível mínimo que serve.** Padrões locais e transacionais (PostgreSQL, Compose, Caddy) antes de orquestradores distribuídos e brokers externos. Complexidade entra com evidência registrada (§5).
3. **Hermetismo.** Builds reproduzíveis, contêineres não-root, segredos cifrados no repositório e decifrados só em memória.
4. **Ferramenta nativa antes de script próprio.** Se GitHub, Docker ou a plataforma já fazem, o Labs não reescreve (Carta §10.4).

---

## 2. Níveis

| Nível | Aplica-se a | Seção |
| :--- | :--- | :--- |
| **0** | Todo projeto, inclusive os entregues a departamentos. | §3 |
| **1** | Artefato que TI vai manter como código de produção (Carta §7.5) ou que atende usuários finais. | §4 |
| **2** | Item a item, quando o gatilho da tabela §5 for observado e registrado. | §5 |

---

## 3. Nível 0 — Todo Projeto

### 3.1. Runtime
- Linguagens permitidas: **Go** e **TypeScript (Node.js)**. Nenhuma outra.
- **Um runtime por projeto.** UI, integração com LLM e jobs leves: TypeScript. Engine, CLI ou daemon de domínio: Go.
- Go não renderiza UI nem faz SSR. Segundo runtime no mesmo projeto só por gatilho de §5.

### 3.2. Estado e Coordenação
- Armazenamento padrão: **PostgreSQL**.
- Trabalho assíncrono: uma tabela de fila consumida com `FOR UPDATE SKIP LOCKED`, gravada na mesma transação do evento (Transactional Outbox). Broker externo só por gatilho de §5.

### 3.3. Contêineres
- Partir dos templates: [`templates/Dockerfile.go`](./templates/Dockerfile.go) e [`templates/Dockerfile.node`](./templates/Dockerfile.node). Eles já trazem multi-stage, cache mounts, `COPY --link`, binário Go stripped (`-trimpath -ldflags="-s -w"`) e usuário não-root (Go: `engine` UID 10001 sobre `alpine:3`; Node: `65532` sobre `gcr.io/distroless/nodejs24:nonroot`).
- Credenciais de build (deploy keys, tokens de módulos privados) entram só via BuildKit secret (`RUN --mount=type=secret,...`), nunca via `ARG`, `ENV` ou `COPY`. Download de módulos privados por SSH usa host keys pinadas em `/etc/ssh/ssh_known_hosts`.
- **Padrão Toolbox:** daemons do mesmo domínio compartilham uma imagem sem `CMD` padrão; o compose define o papel (`command: orchestrator`, `command: watchdog`, `command: migrate`). Esquemas SQL e contratos compartilhados vão em `/usr/share/<app>/`.

### 3.4. Compose
Todo serviço em `compose.yml` declara:
- imagem por digest (`image@sha256:...`);
- `healthcheck`, `restart: unless-stopped` e `deploy.resources.limits.memory`;
- `init: true` (reaper de zumbis e repasse de sinais como PID 1, sem alterar a imagem);
- limite de heap em ~75% do limite de memória: `GOMEMLIMIT=<MB>MiB` (Go) ou `NODE_OPTIONS=--max-old-space-size=<MB>` (Node).

### 3.5. Borda de Rede
- Só o **Caddy** publica portas no host. Bancos, gateways internos e workers ficam em `networks: [internal]`; quem precisa sair para a internet entra também em `egress`.
- Mesma origem, sem CORS: `/healthz` sem autenticação para sondas de deploy; `/api/*` para o gateway com remoção de prefixo; demais rotas para os estáticos da SPA.
- Cache (RFC 9111): assets com hash → `Cache-Control: public, max-age=31536000, immutable`; `index.html` → `Cache-Control: no-cache`.

### 3.6. Segredos
1. Nenhum `.env`, credencial ou chave de API commitado em texto claro.
2. Segredos de cada ambiente moram cifrados com **SOPS + Age** em `<env>/segredos.sops.env`.
3. No deploy, decifrados para arquivo `0600` em `tmpfs` e apagados ao fim (`trap 'rm -rf "$T"' EXIT`).

### 3.7. CI
- Build e testes com `exit 0` antes de qualquer merge.
- Workflows reutilizados e actions de terceiros pinados por commit SHA de 40 hexadecimais.

### 3.8. Verificação em Produção
- Todo serviço expõe `/healthz`.
- Cada fluxo crítico tem **uma** sonda sintética agendada (jornada de caixa-preta com orçamento de latência).
- Sondas em produção são somente-leitura: abortam se emitirem requisição não-GET fora do fluxo de login. Mutação de teste em produção exige flag explícita (`--commit`) e saldo sintético isolado.

### 3.9. Ambientes
**DEV** (banco efêmero, provedores externos mockados, segredos sintéticos) e **PROD**. Ambientes intermediários são Nível 1.

---

### 3.10. Idioma
- Nome em **inglês**: código, serviço, contêiner, binário, variável de ambiente, chave de configuração, tabela e coluna, campo de log.
- Texto em **pt-br**: documentação, README, CLAUDE.md, PR e comentário.
- Vale para nome novo. Nome antigo em português muda só quando a peça já estiver sendo mexida por outro motivo; sem renomeação em massa.

## 4. Nível 1 — Produção Mantida por TI ou Voltada a Usuários Finais

### 4.1. Matriz de Ambientes

| Ambiente | Papel | Imagens | Segredos | Banco |
| :--- | :--- | :--- | :--- | :--- |
| **DEV** | Ciclo rápido do desenvolvedor. | Bundle pinado + override permissivo; hot-reload só do código local. | Sintéticos; LLM, e-mail e alertas desligados. | Efêmero, zero dado real. |
| **ENSAIO** | Pré-voo descartável (`ensaio-<run>`). | Exatamente os digests candidatos. | Sintéticos; sem portas no host. | Limpo; roda todas as migrações e o boot completo sob `trap` de destruição. |
| **QA** | Homologação estável. | `qa/release.yaml`; muda só por PR aprovado. | Reais de homologação (SOPS/Age); externos em modo sandbox. | Persistente; backup antes de toda migração. |
| **PROD** | Usuários finais. | Cópia do manifesto validado em QA; nenhum digest recalculado. | Produção, só em `tmpfs 0600`. | HA, arquivamento contínuo de WAL, rollback declarativo. |

### 4.2. Trio de Manifestos
```sh
docker compose -f compose.yml -f <env>/override.yaml -f <env>/release.yaml config
```
- `compose.yml`: base neutra — serviços, redes, limites; sem imagens, sem portas públicas.
- `<env>/override.yaml`: topologia local — portas do Caddy, volumes, flags do ambiente.
- `<env>/release.yaml`: identidade imutável — CalVer (`AAAA.MM.N`) e matriz de digests `@sha256`.

### 4.3. Proveniência
- CI publica no GHCR e gera atestação com `actions/attest-build-provenance` (`push-to-registry: true`).
- O deploy verifica antes de aplicar: `gh attestation verify oci://<imagem>@sha256:<hash> --owner <org>`. Falha na verificação aborta o deploy.

### 4.4. Promoção
Nenhum workflow tem escrita direta na `main` de repositórios de ambiente.

Implementação de referência: [`wibx-maestro/environments`](https://github.com/wibx-maestro/environments) (`bin/resolve`, `checa`, `gate`, `aplica`, `vigia-runner`). Ela ainda usa `bin/prova` e `bin/autoriza`; a troca pelos recursos nativos de §4.3 e §4.4 está pendente.
1. **Resolução (`bin/resolve`):** imagens candidatas viram digests e CalVer em `qa/release.yaml`.
2. **PR humano:** revisão e checks (`bin/checa`). `CODEOWNERS` cobre os manifestos de ambiente.
3. **Pré-voo (`para=ensaio`):** `bin/gate` sobe o cluster efêmero, roda migrações e testa `/healthz`. Falha aborta sem tocar QA.
4. **Aplicação (`para=qa`):** autorização por **GitHub Environment `qa` com revisores obrigatórios**; runner dedicado (`bin/vigia-runner`); backup antes de aplicar (`bin/aplica`); `docker compose up --wait` com rollback automático em falha de boot.
5. **Produção:** recebe cópia fiel do `release.yaml` homologado em QA, pelo GitHub Environment `prod` com revisores obrigatórios. Não resolve imagem nem recompila.

### 4.5. Correção de Falhas
Falha de sonda é isolada antes de culpar o código (cliente de teste, rede/gateway, permissão de conta ou regressão real). O patch roda no ENSAIO e precisa passar no **cenário exato que falhou**; o caso só fecha após 3 ciclos consecutivos verdes da sonda.

---

## 5. Nível 2 — Por Gatilho

Cada item entra quando o gatilho for observado e registrado no projeto. Sem registro, não entra.

| Item | Gatilho |
| :--- | :--- |
| Segundo runtime (Go em projeto TS, ou o inverso) | Carga de CPU ou concorrência medida que o runtime atual não sustenta. |
| Contratos Protobuf via **Buf** + gRPC | Dois runtimes no projeto com chamada **síncrona** entre eles. |
| Segunda tabela de fila (separar eventos de tarefas) | Segundo formato de consumidor que não cabe na tabela única. |
| Build multi-arch (`--platform=$BUILDPLATFORM` + `TARGETOS`/`TARGETARCH`, sem QEMU) | Arquitetura de deploy diferente da de build. |
| Jornadas sintéticas além da sonda mínima | Regressão em produção que a sonda mínima não pegou. |
| Invariant watchers (invariantes relacionais, backlog do outbox) | Incidente registrado de violação de invariante ou fila parada. |

**Fora do Labs:** varredura contínua de superfície de ataque (EASM), auditoria de credenciais vazadas e CVEs de imagens são responsabilidade de TI/Segurança (Carta §10.4).

---

## 6. Checklist de Conformidade (Definition of Done de Infraestrutura)

**Nível 0**
- [ ] Um runtime (Go ou TypeScript); segundo só com gatilho registrado.
- [ ] Dockerfile derivado dos templates; usuário não-root; credenciais de build só via BuildKit secret.
- [ ] Toda imagem em manifesto por `@sha256`.
- [ ] Todo serviço com `healthcheck`, `limits.memory`, `init: true` e limite de heap.
- [ ] Só o Caddy publica portas; mesma origem; cabeçalhos de cache RFC 9111.
- [ ] Nenhum segredo em texto claro; SOPS + Age.
- [ ] CI com build e testes `exit 0`; actions pinadas por SHA de 40 hex.
- [ ] `/healthz` e uma sonda sintética por fluxo crítico.
- [ ] `LICENSE` proprietário na raiz; manifestos com `"private": true` e `"license": "UNLICENSED"`.
- [ ] Nome novo em inglês; documentação, PR e comentário em pt-br.

**Nível 1** — além do Nível 0:
- [ ] Trio de manifestos e ambientes ENSAIO e QA.
- [ ] Atestação `attest-build-provenance` verificada no deploy.
- [ ] Deploy aprovado no ENSAIO (`bin/gate`) antes de QA.
- [ ] GitHub Environments `qa` e `prod` com revisores obrigatórios; `CODEOWNERS` nos manifestos.
- [ ] Backup antes de migração em QA e PROD.

---

## 7. Autoria e Propriedade Intelectual

Todo repositório Labs pertence exclusivamente à **WIBX Company**.

1. **`LICENSE` na raiz:** proprietário, declara a titularidade da WIBX Company e veda cópia, redistribuição, engenharia reversa e uso não autorizado.
2. **Pacotes privados:** `"private": true` e `"license": "UNLICENSED"` (ou equivalente do ecossistema).
3. **Sem cabeçalho de copyright por arquivo.** O `LICENSE` cobre o repositório.
4. **Entidade única:** o código pertence à `WIBX Company`. Razões sociais específicas ficam restritas a contratos, Termos de Uso e documentos fiscais.
5. **Sem autoria individual:** proibido `@author`, nome ou e-mail de desenvolvedor no código. O histórico do Git é a fonte de autoria.
6. **Comentários de domínio:** TSDoc/Godoc explicam o **porquê** — regra de negócio, invariante de segurança, contrato de API. Sem narrativa de histórico ("alterado por X").
