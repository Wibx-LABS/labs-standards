# PADRÃO DE DADOS LABS (LABS DATA STANDARD)

**Classificação:** Norma Canônica de Modelagem, Migração e Operação de Bancos de Dados  
**Status:** RASCUNHO  
**Última Revisão:** Outubro de 2026  

---

Como o Labs constrói, muda e mantém bancos e tabelas. Escopo: **PostgreSQL**, o armazenamento padrão da Norma 01 §3.2. Dados analíticos ficam fora desta norma. Segue os níveis da Norma 01: todo projeto cumpre o Nível 0; o Nível 1 entra com o mesmo gatilho (código de produção mantido por TI ou voltado a usuários finais); o Nível 2 entra item a item, por gatilho registrado.

Base de mercado: documentação do PostgreSQL 18, a página *Don't Do This* do wiki do PostgreSQL, a RFC 9562 (UUID) e as regras do linter `squawk`. Boa parte já é prática do Labs: migração só para frente, `timestamptz`, `CHECK` no lugar de enum, `numeric` para dinheiro, papel por serviço, outbox com `SKIP LOCKED` e ensaio de migração em cópia restaurada. A norma escreve o que já funciona e fecha as divergências.

**A norma governa o resultado, não a ferramenta.** Qualquer migrador serve (runner próprio, Kysely, Prisma, drizzle-kit, golang-migrate), desde que cumpra §3.5.

---

## 1. Princípios

1. **O banco garante o que pode garantir.** Tipo, `NOT NULL`, `CHECK`, `UNIQUE` e FK no banco, não só na aplicação. Validação na aplicação é para mensagem de erro; a constraint é para a verdade.
2. **Migração aplicada é imutável.** Erro numa migração aplicada se corrige com migração nova. O histórico é o que rodou em produção.
3. **Mudança que quebra vai em duas versões.** Expandir, migrar, contrair. Nenhuma migração deixa a versão anterior do código sem funcionar.
4. **Dado pessoal tem endereço.** Toda coluna com dado pessoal é localizável e apagável sem caça ao tesouro.
5. **Backup que nunca foi restaurado não é backup.**

---

## 2. Níveis

| Nível | Aplica-se a | Seção |
| :--- | :--- | :--- |
| **0** | Todo projeto com PostgreSQL. | §3 |
| **1** | Mesmo gatilho da Norma 01 §2. | §4 |
| **2** | Item a item, quando o gatilho da tabela §5 for observado e registrado. | §5 |

**Projetos existentes.** Critério: **segurança e perda de dado valem já; eficiência vale para tabela nova.**
- Valem já, inclusive em tabela existente: §3.6 (credenciais e papéis), §3.8 (backup) e, no Nível 1, `FORCE ROW LEVEL SECURITY` (§4.1).
- O resto vale para **tabela e coluna novas**. Tabela existente não é migrada só para conformidade; converte quando for tocada por outro motivo.

`FORCE` está entre as imediatas porque o risco é de isolamento e o custo é mínimo. Teste em PostgreSQL 17.11 com 2 milhões de linhas em dois tenants: sem `FORCE`, o dono da tabela leu as linhas do outro tenant; com `FORCE`, não. O `ALTER TABLE ... FORCE ROW LEVEL SECURITY` levou 0,4 ms e não reescreveu a tabela. A migração usa `lock_timeout` (§3.5).

---

## 3. Nível 0 — Todo Projeto

### 3.1. Chaves Primárias
- Coluna `id`, tipo nativo **`uuid`**. Nunca `text`/`varchar` guardando UUID: 36 bytes e comparação de texto contra 16 bytes.
- Valor padrão **UUIDv7**: `DEFAULT uuidv7()` no PostgreSQL 18+; abaixo do 18, gerado na aplicação. v7 é ordenado no tempo e mantém o índice da PK compacto; v4 espalha inserções pelo B-tree.
- **UUIDv4** (`gen_random_uuid()`) quando o identificador sai para fora e o instante de criação é sensível: token, link de convite, link de redefinição. v7 carrega o timestamp de criação em milissegundos nos primeiros 48 bits.
- Tabela interna, nunca exposta e de alto volume pode usar `bigint GENERATED ALWAYS AS IDENTITY`. Nunca `serial` nem `int`.
- **Autorização nunca depende do identificador ser impossível de adivinhar** (RFC 9562 §8). ID que sai para fora é público.
- Ordem de v7 não é sequência: não substitui `created_at` nem contador.

Medição de referência (PostgreSQL 18.6, 2 milhões de linhas, inserção em lote; duas execuções em ordem inversa):

| Chave | Índice da PK | Inserção | Fragmentação do índice |
| :--- | :--- | :--- | :--- |
| `bigint` identity | 43 MB | 8 s | 0% |
| `uuid` v7 | 60 MB | 12 s | 0% |
| `uuid` v4 | 76 MB | 20–22 s | 50% |
| `text` com UUID v4 | 145 MB | 20–27 s | 50% |

### 3.2. Tipos
| Uso | Tipo | Proibido |
| :--- | :--- | :--- |
| Instante | `timestamptz` | `timestamp`, `timetz`, instante em `text` |
| Dinheiro | `numeric(p,s)` ou inteiro em unidade mínima + coluna de moeda | `money`, `float`, `real` |
| Texto | `text` (+ `CHECK` de tamanho, se houver regra) | `varchar(n)`, `char(n)` |
| Conjunto fechado de valores | `text` + `CHECK (col IN (...))`; tabela de referência se o conjunto muda com frequência | `CREATE TYPE ... AS ENUM`; valores permitidos só em comentário |
| Payload externo ou esparso | `jsonb` | `json`; `jsonb` para campo que é filtrado, juntado ou validado (vira coluna) |

- Por que não `ENUM` (testado em PostgreSQL 18.6): remover valor não existe (`dropping an enum value is not implemented`); o caminho é recriar o tipo e **reescrever a tabela**. Com `CHECK`, troca-se a constraint por `NOT VALID` + `VALIDATE`, sem reescrita. E valor novo de enum não pode ser usado na mesma transação que o criou (`unsafe use of new value`). Renomear valor funciona nos dois.
- `NOT NULL` por padrão. Nulável só quando a ausência tem significado.
- Toda tabela tem `created_at timestamptz NOT NULL DEFAULT now()`. Tabela com `UPDATE` tem `updated_at`.

### 3.3. Nomes
- `snake_case` minúsculo, sem aspas. Nome entre aspas (`"applicationId"`) é proibido.
- Projeto novo: tabelas em **inglês, no plural** (`orders`, `doc_versions`); colunas no singular. Projeto existente mantém a convenção que já tem (o recybacks segue em português, singular) — um projeto, uma convenção, sem mistura dentro dele. Inglês no plural venceu português no singular: é a convenção da maioria dos schemas do Labs (maestro, support_bot, WBCS) e a dos exemplos do próprio PostgreSQL; a escolha é de convenção, não de evidência técnica.
- FK: `<entidade>_id`. Constraints e índices: `<tabela>_<colunas>_{pkey,fkey,key,check,idx}`, que é o padrão que o próprio PostgreSQL gera.

### 3.4. Integridade
- Toda FK tem índice na coluna que referencia. O PostgreSQL não cria esse índice; sem ele, `DELETE` no pai varre o filho. Testado em PostgreSQL 18.6: apagar 500 pais de uma tabela filha com 2 milhões de linhas levou 33 s sem o índice e 16 ms com ele.
- Idempotência (chave de evento, chave de requisição) é `UNIQUE` no banco, não `SELECT` antes do `INSERT`.

### 3.5. Migrações
1. **Só para frente.** Arquivos `NNNN_nome.sql` numerados, sem migração de descida. Desfazer passo destrutivo é restauração (§3.8), não `down`.
2. **Imutáveis depois de aplicadas.** Um teste reprova alteração em migração já existente na `main` (checksum ou comparação com o histórico).
3. **Destrutiva é marcada.** Migração com `DROP`, `TRUNCATE`, `RENAME` ou `DELETE` começa com `-- destrutiva: <motivo>`; teste reprova a falta do marcador.
4. **CI aplica todas do zero**, na mesma versão maior do PostgreSQL de produção.
5. **Sem travar o tráfego:**
   - `SET lock_timeout = '2s'` no início da migração; falhou, tenta de novo. Migração esperando lock enfileira **todas** as consultas que chegam depois, mesmo as que não conflitam com quem segura o lock. O valor de `lock_timeout` é o tempo máximo de parada de todo o tráfego da tabela por tentativa. Testado em PostgreSQL 18.6: sem `lock_timeout`, um `SELECT` instantâneo esperou 6,1 s atrás de uma migração enfileirada; com `2s`, a migração desistiu e o `SELECT` seguiu em 1,5 s. 2 s venceu 5 s: parada ocasional de 2 s é tolerável, 5 s em tabela quente vira timeout para o usuário; o custo é mais tentativas em tabela movimentada.
   - Índice: `CREATE INDEX CONCURRENTLY` / `DROP INDEX CONCURRENTLY`, em migração fora de transação.
   - FK e `CHECK` em tabela existente: `ADD CONSTRAINT ... NOT VALID`, depois `VALIDATE CONSTRAINT`.
   - `SET NOT NULL` em tabela existente: `CHECK (col IS NOT NULL) NOT VALID` → `VALIDATE` → `SET NOT NULL`.
   - Renomear coluna, trocar tipo ou remover coluna em uso: expandir/contrair (coluna nova, escrita dupla, backfill em lotes, troca de leitura, remoção numa versão seguinte).
6. **Backfill em lotes**, fora da migração de DDL.
7. **Proibido em produção:** `db push`, sincronização de schema sem arquivo de migração, comando que dropa schema.

Implementação de referência: `wibx-maestro/engine · db.TestMigrationsDestrutivasTemMarcador` e `db.TestBaselineIgualAoHistorico`.

### 3.6. Papéis e Credenciais
- **Dono** dos objetos: papel `NOLOGIN`, usado só pelo migrador.
- **Aplicação:** papel próprio por serviço, só DML nas tabelas que usa, sem posse de objeto, sem `SUPERUSER`, sem `BYPASSRLS`.
- Tabela de auditoria ou registro fiscal: aplicação sem `UPDATE`/`DELETE`.
- String de conexão é segredo (Norma 01 §3.6). **Proibido valor padrão de credencial no código**, inclusive "só para desenvolvimento".

### 3.7. Testes
- Teste de integração roda contra PostgreSQL real, efêmero, na versão maior de produção. Sem SQLite, sem mock de banco.
- Banco de desenvolvimento usa dado sintético. **Proibido copiar dado de produção para DEV.**

### 3.8. Operação
- `pg_stat_statements` em `shared_preload_libraries` desde o provisionamento (exige reinício; decidir depois custa uma janela).
- `log_min_duration_statement` ligado; `idle_in_transaction_session_timeout` definido (transação aberta impede o vacuum).
- Testado em PostgreSQL 18.6: sem preload, `CREATE EXTENSION pg_stat_statements` funciona mas a consulta falha (`must be loaded via "shared_preload_libraries"`); recarregar a configuração não basta, só o reinício. Com uma transação aberta, `VACUUM` depois de apagar 1 milhão de linhas deixou as 1 milhão de linhas mortas; fechada a transação, zero.
- Backup agendado (`pg_dump` ou snapshot) com **cópia fora da máquina do banco**, cifrada. Backup só no mesmo host não conta.
- Falha de backup alerta. Backup pulado em silêncio é falha.

### 3.9. Dados Pessoais
> **Pendente de revisão jurídica.** Os artigos citados são da Lei 13.709/2018 (LGPD), conferidos com o texto compilado do Planalto em outubro de 2026; a interpretação não passou por leitura jurídica.
- Tabela com dado pessoal declara isso no banco: `COMMENT ON TABLE <t> IS 'pii: <colunas>'`. A declaração é verificável por script.
- Dado pessoal fica em colunas próprias, em poucas tabelas. Evento e histórico referenciam o titular por `<titular>_id`, não copiam nome, CPF ou e-mail. Dado pessoal fora de texto livre, `jsonb` de log e mensagem de erro.
- Remoção de dado pessoal é `DELETE`. `deleted_at` só quando o negócio exige desfazer — e então com expurgo agendado.

### 3.10. Vetores (Embeddings)
- Vetor mora no PostgreSQL, na mesma transação, RLS e backup do dado de origem. Banco vetorial dedicado é Nível 2 (§5).
- Tipo **`vector(n)`** (`pgvector`) e similaridade calculada **no banco** (`ORDER BY embedding <=> $1 LIMIT k`). Vetor em `jsonb` com cosseno calculado na aplicação é proibido em tabela nova: medido, ocupa ~1,9× o disco, trafega o vetor inteiro por linha candidata e é ~30× mais lento já com 200 candidatos.
- Busca sem filtro prévio (acervo inteiro) usa índice **`hnsw`**. Com filtro seletivo (documento, tenant), índice B-tree na coluna do filtro e cálculo exato basta.
- Toda linha com vetor carrega a **impressão digital do modelo** (modelo + dimensão). Vetores de impressões diferentes nunca são comparados. Troca de modelo é expandir/contrair: coluna nova, reembedding em lotes, troca de leitura, remoção.
- Janela de candidatos limitada (`LIMIT`, intervalo de tempo) é decisão de produto, não de desempenho: o que fica fora da janela não é encontrado. Escrever o limite e o motivo.

Medição de referência (sintética, 10.000 chunks de 1024 dimensões, PostgreSQL 17 + pgvector):

| | `jsonb` + cosseno na aplicação | `vector(1024)` no banco |
| :--- | :--- | :--- |
| Disco por linha | 7.655 B | 4.100 B |
| Busca em 1 documento (200 chunks) | 39 ms | 1,2 ms (exato) |
| Busca no acervo (10.000 chunks) | 1,9–2,3 s | 33–35 ms (`hnsw`) |

### 3.11. Grafos
- Relação em grafo (hierarquia, dependência, linhagem) é **tabela de arestas** (`origem_id`, `destino_id`, `tipo`) com índice nas duas pontas e FK para os nós.
- Travessia por `WITH RECURSIVE` com **limite de profundidade** explícito e proteção contra ciclo (`CYCLE`, PostgreSQL 14+).
- Banco de grafo é Nível 2 (§5).

---

## 4. Nível 1 — Produção TI

### 4.1. Multi-Tenant (RLS)
Vale para todo serviço que guarda dado de mais de um *whitelabel*.
- Coluna de tenant `wlc` em toda tabela de tenant; índices começam por `wlc`.
- `ENABLE ROW LEVEL SECURITY` **e `FORCE ROW LEVEL SECURITY`**. Sem `FORCE`, o dono da tabela ignora as políticas.
- **O tenant vive só dentro da transação, nunca na sessão.** Conexão de pool é reutilizada por outra requisição; tenant gravado na sessão vaza para ela, sem erro. Entra por **um único ponto** de código. No PostgreSQL: `set_config('app.wlc', $1, true)` — com `true`, o valor vale só na transação corrente; com `false` ou `SET`, vale até o fim da sessão (documentação do PostgreSQL 18, `set_config`).
- Política falha fechada: sem `app.wlc` definido, zero linhas ou erro.
- Teste de RLS roda **como o papel da aplicação**. Rodar como dono ou superusuário passa sem testar nada.
- Política lê só colunas da própria linha; subconsulta em outra tabela dentro de política é lenta e sujeita a corrida.

Implementação de referência: `wibx-recybacks/recybacks · packages/db.comTenant`.

### 4.2. Migrações
- `squawk` no CI sobre toda migração nova. Exceção a uma regra é anotada no próprio arquivo, com motivo.
- Toda migração é ensaiada contra cópia restaurada de produção antes de QA (Norma 01 §4.4, passo 3).
- Verificação de deriva: o schema obtido aplicando as migrações do zero é comparado ao de QA/PROD (`pg_dump --schema-only`). Diferença reprova.

### 4.3. Backup e Recuperação
- Arquivamento contínuo de WAL com backup base (`wal-g` ou `pgBackRest`) em armazenamento fora do host, cifrado, com a chave guardada separada do backup.
- **RPO e RTO escritos** por serviço.
- **Ensaio de restauração com no máximo 7 dias**, terminando em verificação da aplicação (sonda de §3.8 da Norma 01 contra o banco restaurado), com tempo medido contra o RTO. A idade do último ensaio verde é conferida no deploy (passou de 7 dias, ensaia antes de aplicar) e **alerta** quando passa de 7 dias mesmo sem deploy. 7 dias venceu trimestral: ensaio automatizado custa o mesmo em qualquer cadência, e ensaio manual trimestral é o que se pula.

Implementação de referência: `wibx-maestro/environments · bin/arquivo-db` (`ensaio` e `idade`). O alerta sem deploy ainda não existe lá.
- Backup antes de migração em QA e PROD (Norma 01 §6).

### 4.4. Dados Pessoais
> **Pendente de revisão jurídica.**
- Toda tabela com `pii:` (§3.9) tem prazo de retenção escrito e job de expurgo. O prazo é fato jurídico ou contratual; a norma não inventa número, exige que exista.
- Runbook de eliminação a pedido do titular: onde o dado está, como apagar, o que fica e por quê. Base: art. 18, IV (dado desnecessário, excessivo ou em desconformidade) e VI (dado tratado com consentimento); o que fica, só nas hipóteses do art. 16 (por exemplo, obrigação legal, inciso I).
- O dado apagado sobrevive no backup até a expiração dele. O prazo de expiração do backup é publicado na política; eliminação instantânea não é prometida. A lei não trata de backup; apoiar isto em "limites técnicos" (art. 16, caput) é interpretação a validar.

### 4.5. Revisão Trimestral
- Consultas mais caras (`pg_stat_statements`).
- Índices sem uso (`pg_stat_user_indexes.idx_scan = 0` após um ciclo de negócio completo; nunca PK nem `UNIQUE`) e duplicados.
- Inchaço de tabela e índice; autovacuum ajustado por tabela, não globalmente.
- Papel `readonly` para consulta e análise; ninguém consulta produção com o papel da aplicação.

Trimestral venceu cadência menor: a revisão é análise humana, não verificação automática, e índice "sem uso" só significa algo depois de um ciclo de negócio completo de estatística; cadência curta gera falso "sem uso".

---

## 5. Nível 2 — Por Gatilho

Cada item entra quando o gatilho for observado e registrado no projeto. Sem registro, não entra.

| Item | Gatilho |
| :--- | :--- |
| Particionamento (por tempo, com expurgo por `DETACH`/`DROP` de partição) | Expurgo por `DELETE` que demora ou trava, ou tabela que domina o uso da instância. |
| Pooler externo (pgbouncer) | Esgotamento de conexões observado. Em modo transação, só `set_config(..., true)`; `SET` de sessão e advisory lock de sessão quebram. |
| Tabela de histórico/auditoria por entidade (append-only, `UPDATE`/`DELETE` bloqueados) | Exigência regulatória ou contratual de auditoria, ou incidente "quem mudou isto" registrado. |
| Cifragem de coluna na aplicação (AES-256-GCM, chave fora do banco) | Coluna com documento de identificação ou dado pessoal sensível (LGPD art. 5º, II: origem racial ou étnica, religião, opinião política, saúde, vida sexual, dado genético ou biométrico). CPF não é dado sensível na lei; cifrá-lo é escolha do Labs. |
| Ferramenta de migração sem parada (`pgroll`) | Serviço que não aceita janela de deploy. |
| Réplica / standby quente | RTO que a restauração (§4.3) comprovadamente não cumpre. |
| Schema ou banco por tenant | Exigência contratual de isolamento físico. |
| Verificação de deriva noturna contra produção | Incidente de deriva registrado. |
| Banco vetorial dedicado (Qdrant ou similar) | Latência de busca medida acima do orçamento com `hnsw` ajustado, ou índice que não cabe na memória da instância. |
| Apache AGE (grafo dentro do PostgreSQL) | Consulta de muitos saltos ou padrão de caminho no caminho quente, medida lenta com `WITH RECURSIVE` indexado. |
| Banco de grafo dedicado (Neo4j) | Apache AGE comprovadamente insuficiente. Serviço novo: backup, segredo e migração próprios. |

---

## 6. Exceções

- **`wibx-master-console`:** fork do Dokploy; segue as convenções de schema do upstream para não divergir dele. §3.6 (credenciais) e §3.8 (backup) valem integralmente.
- Repositório que só lê banco alheio (consultas, análise de legado) não tem schema próprio e não cumpre §3.1–§3.5.

---

## 7. Checklist de Conformidade

**Nível 0**
- [ ] PK `id uuid`, v7 por padrão; v4 em identificador externo sensível; sem `serial`, sem UUID em texto.
- [ ] `timestamptz`, `numeric` para dinheiro, `text` + `CHECK` no lugar de enum e `varchar(n)`; `created_at` em toda tabela.
- [ ] `snake_case` sem aspas; uma convenção de nomes por projeto.
- [ ] FK com índice; idempotência por `UNIQUE`.
- [ ] Migrações só para frente, imutáveis, destrutivas marcadas, aplicadas do zero no CI.
- [ ] `lock_timeout`, `CONCURRENTLY`, `NOT VALID` + `VALIDATE`; expandir/contrair em mudança que quebra.
- [ ] Dono `NOLOGIN`; papel de aplicação sem posse, sem `BYPASSRLS`; nenhuma credencial padrão no código.
- [ ] Testes contra PostgreSQL real; nenhum dado de produção em DEV.
- [ ] `pg_stat_statements` e log de consulta lenta; backup com cópia fora do host e alerta de falha.
- [ ] Tabelas com dado pessoal declaradas (`pii:`).
- [ ] Vetor em `vector(n)` com similaridade no banco e impressão digital do modelo por linha; `hnsw` em busca sem filtro.
- [ ] Grafo em tabela de arestas indexada; `WITH RECURSIVE` com limite de profundidade.

**Nível 1** — além do Nível 0:
- [ ] RLS com `FORCE`, tenant por `set_config(..., true)` em ponto único, testes como papel da aplicação.
- [ ] `squawk` no CI; ensaio de migração em cópia restaurada; verificação de deriva.
- [ ] WAL contínuo fora do host; RPO/RTO escritos; ensaio de restauração com no máximo 7 dias, com alerta de idade.
- [ ] Retenção e expurgo por tabela `pii:`; runbook de eliminação.
- [ ] Revisão trimestral de consultas, índices e inchaço; papel `readonly`.
