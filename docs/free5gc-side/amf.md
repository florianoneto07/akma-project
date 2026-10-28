# AMF: autenticação primária pedida pela rede doméstica

**MVP:** MVP-12 · **Data:** 28 de outubro de 2026 · **Aplicado e verificado na VM do free5GC** (192.168.16.130) em 28/10/2026

## 1. O que o AMF precisa fazer

`[3GPP]` TS 33.501 §6.1.5, a autenticação primária disparada pela rede doméstica. É opcional para a rede doméstica e a servidora (§6.1.5.1). Quando as duas a suportam:

- passo 0b: no `Nudm_UECM_Registration`, o AMF *"shall provide a callback URI within the AMF registration for the UDM to create an implicit subscription"* para a `Nudm_UECM_Re-AuthenticationNotification`;
- passo 3: ao receber a notificação, o AMF decide pela sua política e pelo estado do UE, e responde ao UDM com uma causa de recusa ou aceita. Se aceita e não alcança o UE, *"shall set the authentication pending flag"* e faz a reautenticação *"When UE re-attaches to the same AMF or becomes reachable"*;
- passo 4: *"The AMF/SEAF starts the primary authentication procedure as defined in clause 6.1.2"*.

Depois da autenticação, §6.1.3.2: *"the AMF shall initiate NAS security mode command procedure (see clause 6.7.2) with the UE, to take the newly generated partial native 5G NAS security context into use"*. §6.9.4.2: as chaves NAS novas vêm do novo K_AMF. §6.9.4.4: a troca das chaves AS *"may be used"*.

`[3GPP]` TS 29.503 §5.3.2.15.2 e §6.2.5.6: o UDM faz `POST` no `reauthNotifyCallbackUri` com `ReauthNotificationInfo`. O AMF responde `204`, ou `202` com o *flag* de pendência; recusas na tabela 6.2.5.6-3.

`[FREE5GC]` O AMF v1.4.5 não tem nada disso (correções F18 e F19 de `correcoes-free5gc.md`):

- registra-se no UDM só com o `deregCallbackUri`, e não tem a rota da notificação;
- a máquina de estados do GMM só vai de `Registered` a `Authentication` num `Registration Request`;
- `AuthenticationProcedure` pula a autenticação quando o contexto de segurança vale, e o estado `SecurityMode` pula o NAS SMC pelo mesmo motivo;
- depois do SMC, o estado `ContextSetup` espera um `Registration Request`.

Especificação aprovada: `renovacao-kakma.md`, pontos K5 a K7, com as decisões de 28/10/2026 na implementação (seção 10 daquele documento).

## 2. O que o *patch* faz

*Patch*: `free5gc/patches/amf/0001-home-network-triggered-authentication.patch`.

1. **Registro no UDM** (`internal/sbi/consumer/udm_service.go`). No acesso 3GPP, o registro leva `reauthNotifyCallbackUri` = `<AMF>/namf-callback/v1/reauth/<SUPI>`, como o `deregCallbackUri` (K5).
2. **Rota da notificação** (`internal/sbi/api_reauth.go`, `internal/sbi/server.go`). `POST /namf-callback/v1/reauth/:supi`, num grupo de rotas sem a verificação de *token*: OAuth 2.0 não se aplica a notificações (TS 29.500 §6.7.3; TS 33.501 §13.4.1.0). As outras rotas de *callback* continuam com a verificação do free5GC (F20). `oldGuami` e `ausfRediscoveryInd` são lidos e ignorados (K8).
3. **Decisão do AMF** (`internal/gmm/hn_auth.go`, `ReauthNotification`), pelo estado do UE no acesso 3GPP (tabela 6.2.5.6-3):

   | Estado do UE | Resposta |
   |---|---|
   | `Registered`, conectado, sem procedimento em curso | `204`; a autenticação começa na hora |
   | `Registered`, em CM-IDLE (ou sendo chamado por *paging*) | `202`; *flag* de pendência |
   | `Authentication`, `SecurityMode` ou `ContextSetup` de um registro, ou atualização de registro em curso | `409 TEMPORARY_REJECT_REGISTRATION_ONGOING` |
   | *handover* N2 em curso | `409 TEMPORARY_REJECT_HANDOVER_ONGOING` |
   | autenticação pedida pelo UDM em curso, ou desregistro em curso | `403 REAUTHENTICATION_NOT_ALLOWED` |
   | UE desconhecido ou desregistrado | `404 CONTEXT_NOT_FOUND` |
   | a autenticação não começou (AUSF não encontrado ou sem resposta) | `500 SYSTEM_FAILURE`; o UE segue em `Registered` |

4. **Autenticação fora do registro** (`internal/gmm/hn_auth.go`, `handler.go`, `sm.go`). O AMF escolhe um ngKSI diferente do corrente para o contexto parcial que a autenticação cria (TS 33.501 §6.1.3.2 passo 6) e vai de `Registered` a `Authentication`. Com `HnAuthRunning`:
   - `AuthenticationProcedure` não pula a autenticação; o `Authentication Request` sai protegido pelo contexto corrente;
   - uma falha antes de o UE entrar (AUSF) não manda `Registration Reject`: o novo evento `ReauthAbortEvent` leva de volta a `Registered`, com o ngKSI anterior;
   - o estado `SecurityMode` faz o NAS SMC com o novo K_AMF, embora o contexto corrente seja válido (§6.1.3.2, §6.9.4.2);
   - o SMC não pede a mensagem inicial (RINMR = 0) nem a derivação horizontal do K_AMF (HDP = 0), que não cabem fora de um registro (`internal/gmm/message/build.go`);
   - o `Security Mode Complete` leva, pelo novo evento `ReauthSuccessEvent`, direto a `Registered`, sem o `ContextSetup` de um registro (K7). Um `NAS message container` nessa mensagem é ignorado, com um *warning*: fora de um registro ou de um `Service Request` não há mensagem inicial, e o SMC não a pediu. `[FREE5GC]` O UERANSIM põe ali o último `Registration Request` sempre que esse registro foi mandado sem contexto de segurança (`m_lastRegWithoutNsc` fica ligado depois do registro). Sem essa regra, o AMF tratava o conteúdo como um registro novo e ficava em `ContextSetup`, como se viu na segunda execução do `e2e.sh` de 28/10/2026.

   A falha da autenticação em si (RES\* errado, recusa do AUSF, `Authentication Failure` do UE) segue o tratamento do free5GC (K6): `Authentication Reject` e remoção do contexto.
5. **Pendência** (`sm.go`). O AMF faz a autenticação pendente logo depois do próximo `Service Request` aceito, ou ao voltar a `Registered` depois de um registro (§6.1.5.2 passo 3).
6. **Logs** (`internal/logger/logger.go`): categoria nova `AKMA`, no nível *info* e sem chaves (`logs-akma.md`):
   - `Re-authentication of <SUPI> requested by the UDM: started (ngKSI <n>)`, `pending (UE not connected)` ou `refused, <status> <causa>`;
   - `Re-authentication of <SUPI> requested by the UDM <tempo> ago: started (ngKSI <n>)`, para a pendente;
   - `Re-authentication of <SUPI> requested by the UDM: completed, new KAMF in use (ngKSI <n>)`;
   - `Re-authentication of <SUPI> requested by the UDM: not started` (*warning*).

`[MVP]` Sem a troca das chaves AS (§6.9.4.4, opcional): o UERANSIM não tem segurança AS. O gNB segue com o K_gNB do K_AMF anterior até o próximo estabelecimento do contexto AS. Só o acesso 3GPP: o laboratório não tem acesso não 3GPP.

`[FREE5GC]` O AMF não serializa os eventos do GMM que chegam por NAS e por SBI; o *patch* faz o mesmo que as outras operações SBI do AMF, e as notificações de reautenticação entre si são serializadas.

## 3. Arquivos tocados

Árvore: `~/free5gc/NFs/amf` (`free5gc/amf` v1.4.5). Binário: `~/free5gc/bin/amf`, com o original em `bin/amf.orig`.

| Arquivo | Mudança |
|---|---|
| `go.mod`, `go.sum` | `replace` do `openapi` pela cópia com *patch* da versão que o AMF fixa, `v1.2.5-0.20260527003827-02dc71b4d94f` (`free5gc/build/openapi-test`, a mesma do banco de testes) |
| `internal/context/amf_ue.go` | `HnAuthPending`, `HnAuthPendingSince`, `HnAuthRunning` |
| `internal/gmm/init.go` | eventos `ReauthSuccessEvent` (`SecurityMode` → `Registered`) e `ReauthAbortEvent` (`Authentication` → `Registered`) |
| `internal/gmm/hn_auth.go` (novo) | `ReauthNotification`, início da autenticação, pendência |
| `internal/gmm/handler.go` | `AuthenticationProcedure` e `HandleSecurityModeComplete` |
| `internal/gmm/sm.go` | estados `Registered`, `Authentication`, `SecurityMode` e `DeRegistered` |
| `internal/gmm/message/build.go` | RINMR e HDP do `Security Mode Command` |
| `internal/gmm/hn_auth_test.go` (novo) | recusas, pendência, início com AUSF simulado e gNB falso, falha do AUSF, RINMR e HDP, `Security Mode Complete` com `NAS message container` |
| `internal/sbi/api_reauth.go` (novo), `internal/sbi/server.go` | rota da notificação, sem verificação de *token* |
| `internal/sbi/consumer/udm_service.go` | `reauthNotifyCallbackUri` no registro do acesso 3GPP |
| `internal/logger/logger.go` | categoria `AKMA` |

## 4. Como reproduzir numa instalação limpa

Com o repositório em `~/akma-free5gc`, ao lado de `~/free5gc`, e o free5GC parado. A cópia do `openapi` da versão do AMF, com os *patches* até o `0003`:

```bash
cd ~/akma-free5gc && free5gc/scripts/prepare-openapi.sh free5gc/build/openapi-test v1.2.5-0.20260527003827-02dc71b4d94f
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh amf apply
```

Aplica o *patch*, guarda `bin/amf.orig` e recompila. Para desfazer, use `patch-nf.sh amf revert`.

## 5. Como verificar

```bash
cd ~/free5gc/NFs/amf && go test ./internal/gmm/ ./internal/sbi/...
```

Ponta a ponta, o passo R6r do `e2e.sh` (`docs/af-side/e2e.md`): o AAnF pede a renovação da K_AKMA, o UDM notifica o AMF, o AMF autentica o UE registrado e faz o NAS SMC, e o UE passa a ter um A-KID novo sem perder a sessão PDU. O caso do UE em CM-IDLE (`202` e pendência) só é verificado nos testes de unidade, porque o UERANSIM não fica em CM-IDLE no laboratório.

Resultado em 28/10/2026: os testes do AMF passam, também com `-race`. Na execução `20261028-162316` do `e2e.sh`, o R6r passou inteiro: o AMF recebeu a notificação, autenticou o UE com o ngKSI 1, ignorou o `NAS message container` do `Security Mode Complete` e voltou a `Registered` (`renovacao-kakma.md` §11).
