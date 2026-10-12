# LABS TASK SPECIFICATION: {{NOME_DA_TAREFA}}

**Data:** {{DATA_ISO}}  
**Autor Arquitetural:** Claude Opus (Headless Session)  
**Status do Contrato:** FECHADO / IMUTÁVEL  

---

## 1. Escopo Físico Estrito (Filesystem Fencing)
*Regra Dura:* Qualquer modificação em arquivo fora desta lista resulta em rejeição automática do PR.
- **Arquivos com permissão de escrita:**
  - `{{CAMINHO_ARQUIVO_1}}`
  - `{{CAMINHO_TESTE_1}}`
- **Arquivos de consulta apenas-leitura (LSP/Tipos):**
  - `{{CAMINHO_CONTRATO_PROTO_OU_INTERFACE}}`

---

## 2. Contrato de Assinaturas & Tipagens (LSP Delta)
Substituição cirúrgica exigida. Proibido alterar assinaturas públicas não listadas aqui.

### Assinatura Atual (Antes)
```{{LINGUAGEM}}
// Assinatura extraída via LSP
```

### Assinatura Alvo (Depois)
```{{LINGUAGEM}}
// Assinatura exata esperada pelo design
```

---

## 3. Invariantes Invioláveis de Negócio & Segurança
1. **Invariante 1:** (ex.: Transações financeiras exigem idempotência estrita no PostgreSQL Blackboard).
2. **Invariante 2:** (ex.: Tratamento de erro sem panic; propagação explícita com contexto de erro).
3. **Invariante 3:** (ex.: Proibido adicionar goroutines desgovernadas sem context cancellation).

---

## 4. Non-Goals (O Que o Worker ESTÁ PROIBIDO de Fazer)
- ❌ Proibido refatorar funções vizinhas ou "limpar código" fora do escopo.
- ❌ Proibido adicionar novas dependências no manifesto de pacotes (`go.mod`, `package.json`).
- ❌ Proibido reescrever a suíte de testes existente; apenas estender com o caso novo.

---

## 5. Juiz Compilador (Comando Determinístico de Validação)
O worker deve rodar este comando exato antes de dar a tarefa como encerrada:
```bash
rtk {{COMANDO_DE_TESTE_OU_BUILD}}
```
*Critério de Sucesso:* `exit code == 0` sem supressão de flags de erro.
