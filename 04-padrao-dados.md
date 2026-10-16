# PADRÃO DE DADOS LABS (LABS DATA STANDARD)

**Classificação:** Norma Canônica de Bancos de Dados  
**Status:** ATIVO / CANÔNICO  
**Última Revisão:** Outubro de 2026  

---

Como o Labs faz bancos **bem estruturados, eficientes, robustos, seguros e sem ponto único de falha**. Escopo: **PostgreSQL**, o armazenamento padrão da Norma 01 §3.2. Segue os níveis da Norma 01: cada regra é marcada **[0]** (todo projeto) ou **[1]** (mesmo gatilho do Nível 1 da Norma 01); o Nível 2 entra por gatilho registrado (§6).

A norma governa o resultado, não a ferramenta: qualquer migrador serve, desde que cumpra §3. Os testes e fontes que sustentam cada regra estão na PR que introduziu esta norma.

**Projetos existentes.** Regra de schema (§1, §2) vale para tabela e coluna novas; tabela existente converte quando for tocada por outro motivo. Regra de segurança e de operação (§3 item 7, §4, §5) vale já.

---

## 1. Bem Estruturado

### Chaves
- **[0]** PK `id` do tipo nativo `uuid`, padrão **UUIDv7** (`DEFAULT uuidv7()` no PostgreSQL 18+; abaixo disso, gerado na aplicação). Nunca UUID em `text`/`varchar`: índice 2,4× maior.
- **[0]** **UUIDv4** quando o identificador sai para fora e o instante de criação é sensível (token, convite, link de redefinição): v7 carrega o timestamp de criação.
- **[0]** Tabela interna, nunca exposta e de alto volume pode usar `bigint GENERATED ALWAYS AS IDENTITY`. Nunca `serial` nem `int`.
- **[0]** Autorização nunca depende do identificador ser difícil de adivinhar.

### Tipos
- **[0]** Instante: `timestamptz`. Proibido `timestamp`, `timetz` ou instante em `text`.
- **[0]** Dinheiro: `numeric(p,s)`, ou inteiro em unidade mínima + coluna de moeda. Proibido `money` e ponto flutuante.
- **[0]** Texto: `text`, com `CHECK` de tamanho se houver regra. Proibido `varchar(n)`.
- **[0]** Conjunto fechado de valores: `text` + `CHECK (col IN (...))`. Proibido `ENUM`: remover valor exige reescrever a tabela.
- **[0]** `jsonb` só para payload externo ou esparso; campo filtrado, juntado ou validado vira coluna.
- **[0]** `NOT NULL` por padrão. Toda tabela tem `created_at timestamptz NOT NULL DEFAULT now()`.

### Nomes
- **[0]** `snake_case` minúsculo, sem aspas. Projeto novo: tabelas em inglês, no plural. Um projeto, uma convenção.
- **[0]** FK: `<entidade>_id`. Constraints e índices: `<tabela>_<colunas>_{pkey,fkey,key,check,idx}`.

### Integridade
- **[0]** O banco garante o que pode: tipo, `NOT NULL`, `CHECK`, `UNIQUE` e FK no banco, não só na aplicação.
- **[0]** Idempotência é `UNIQUE` no banco, não `SELECT` antes do `INSERT`.

### Vetores e Grafos
- **[0]** Vetor: `vector(n)` (`pgvector`) no mesmo banco do dado de origem, similaridade calculada no banco. Proibido vetor em `jsonb` com cálculo na aplicação.
- **[0]** Toda linha com vetor carrega o modelo e a dimensão que o geraram; vetores de modelos diferentes nunca são comparados.
- **[0]** Grafo: tabela de arestas (`origem_id`, `destino_id`, `tipo`) indexada nas duas pontas; travessia por `WITH RECURSIVE` com limite de profundidade e `CYCLE`.

---

## 2. Eficiente

- **[0]** Toda FK tem índice na coluna que referencia. O PostgreSQL não o cria; sem ele, `DELETE` no pai varre o filho.
- **[0]** Busca vetorial no acervo inteiro usa índice `hnsw`. Com filtro seletivo (documento, tenant), B-tree no filtro basta.
- **[0]** `pg_stat_statements` ligado desde o provisionamento (ligar depois exige reinício) e `log_min_duration_statement` definido.
- **[0]** `idle_in_transaction_session_timeout` definido: transação esquecida impede o vacuum e a tabela incha sem erro.
- **[1]** Revisão trimestral: consultas mais caras, índices sem uso (nunca PK nem `UNIQUE`) e duplicados, inchaço.

---

## 3. Robusto

1. **[0]** Migração só para frente: `NNNN_nome.sql`, sem `down`. Desfazer passo destrutivo é restauração.
2. **[0]** Migração aplicada é imutável; teste reprova alteração em migração existente na `main`.
3. **[0]** Migração destrutiva (`DROP`, `TRUNCATE`, `RENAME`, `DELETE`) começa com `-- destrutiva: <motivo>`; teste reprova a falta.
4. **[0]** CI aplica todas as migrações do zero, na versão maior do PostgreSQL de produção, e roda os testes contra PostgreSQL real. Sem SQLite, sem mock de banco.
5. **[0]** Sem travar o tráfego:
   - `SET lock_timeout = '2s'` no início; falhou, tenta de novo. Migração esperando lock bloqueia todas as consultas que chegam depois.
   - Índice com `CONCURRENTLY`, fora de transação.
   - FK, `CHECK` e `SET NOT NULL` em tabela existente: `NOT VALID`, depois `VALIDATE`.
   - Renomear, trocar tipo ou remover coluna em uso: expandir/contrair, em duas versões.
6. **[0]** Backfill em lotes, fora da migração de DDL.
7. **[0]** Proibido em produção: `db push`, sincronização de schema sem arquivo de migração, comando que dropa schema.
8. **[1]** `squawk` no CI sobre toda migração nova; exceção anotada no arquivo, com motivo.
9. **[1]** Migração ensaiada contra cópia restaurada de produção antes de QA (Norma 01 §4.4).

Implementação de referência: `wibx-maestro/engine · db.TestMigrationsDestrutivasTemMarcador` e `db.TestBaselineIgualAoHistorico`.

---

## 4. Seguro

### Papéis e Credenciais
- **[0]** Dono dos objetos: papel `NOLOGIN`, usado só pelo migrador.
- **[0]** Aplicação: papel próprio por serviço, só DML, sem posse de objeto, sem `SUPERUSER`, sem `BYPASSRLS`.
- **[0]** String de conexão é segredo (Norma 01 §3.6). Proibido valor padrão de credencial no código.
- **[1]** Papel `readonly` para consulta e análise; ninguém consulta produção com o papel da aplicação.

### Multi-Tenant
- **[1]** Serviço com dado de mais de um *whitelabel*: coluna `wlc` em toda tabela de tenant, `ENABLE` **e `FORCE ROW LEVEL SECURITY`**. Sem `FORCE`, o dono da tabela ignora as políticas.
- **[1]** O tenant vive só na transação, nunca na sessão: `set_config('app.wlc', $1, true)`, num único ponto do código. Tenant na sessão vaza para a próxima requisição que reusa a conexão.
- **[1]** Política falha fechada (sem tenant, zero linhas) e é testada como o papel da aplicação, nunca como dono.

Implementação de referência: `wibx-recybacks/recybacks · packages/db.comTenant`.

### Dados Pessoais
- **[0]** Tabela com dado pessoal declara no banco: `COMMENT ON TABLE <t> IS 'pii: <colunas>'`.
- **[0]** Dado pessoal em colunas próprias, em poucas tabelas; evento e log referenciam o titular por id, não copiam o dado.
- **[0]** Remoção é `DELETE`; `deleted_at` só com expurgo agendado.
- **[0]** Proibido copiar dado de produção para DEV.

**Fora do Labs:** prazos de retenção, base legal e atendimento a pedido do titular são do encarregado de dados (LGPD art. 41). A norma garante que o dado é localizável e apagável; quando e por quê apagar, não.

---

## 5. Sem Ponto Único de Falha

Backup protege contra dado destruído; standby protege contra servidor morto. Um não substitui o outro.

- **[0]** Backup agendado, cifrado, com cópia **fora do host do banco**. Falha de backup alerta.
- **[1]** Arquivamento contínuo de WAL (`wal-g` ou `pgBackRest`) fora do host; chave guardada separada do backup.
- **[1]** **Standby em outro host**, por replicação assíncrona, com atraso de replicação monitorado e alertado. Standby no mesmo host não conta: cai junto.
- **[1]** Promoção manual (`pg_promote()`) com runbook escrito; o PostgreSQL não promove sozinho.
- **[1]** RPO e RTO escritos por serviço.
- **[1]** Ensaio com no máximo 7 dias, de **restauração** e de **promoção**, terminando na sonda da Norma 01 §3.8 contra o banco resultante, com tempo medido contra o RTO. A idade do último ensaio verde é conferida no deploy e alerta quando passa de 7 dias.

Implementação de referência (restauração): `wibx-maestro/environments · bin/arquivo-db`.

---

## 6. Nível 2 — Por Gatilho

Cada item entra quando o gatilho for observado e registrado no projeto.

| Item | Gatilho |
| :--- | :--- |
| Failover automático (Patroni ou similar, 3+ nós) | RTO de minutos fora do horário em que há alguém para promover. |
| Replicação síncrona | Perda das últimas transações na promoção é inaceitável para o serviço. Todo commit passa a esperar o standby. |
| Réplica de leitura | CPU ou I/O do primário saturado por leitura, medido. |
| Pooler externo (pgbouncer) | Esgotamento de conexões observado. |
| Particionamento | Expurgo por `DELETE` que demora ou trava, ou tabela que domina a instância. |
| Tabela de histórico por entidade | Exigência de auditoria, ou incidente "quem mudou isto" registrado. |
| Cifragem de coluna na aplicação | Coluna com documento de identificação ou dado pessoal sensível. |
| Banco dedicado (vetorial, grafo) | Latência medida que o PostgreSQL ajustado (`hnsw`, `WITH RECURSIVE` indexado) não cumpre. |

---

## 7. Exceções

- **`wibx-master-console`:** fork do Dokploy; segue as convenções de schema do upstream (§1, §2). §3 item 7, §4 e §5 valem integralmente.
- Repositório que só lê banco alheio não tem schema próprio e não cumpre §1–§3.
