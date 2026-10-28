# UDM: indicação AKMA, Routing Indicator e `trigger-auth`

**MVP:** MVP-2; seção 6 no MVP-12 · **Data:** 17 de outubro de 2026; seção 6 em 28 de outubro de 2026 · **Aplicado e verificado na VM do free5GC** (192.168.16.130) em 17/10/2026; seção 6 aplicada em 28/10/2026

## 1. O que o UDM precisa fazer

`[3GPP]` TS 33.501 §6.1.3.1 passo 2 (EAP-AKA') e §6.1.3.2.0 passo 2 (5G AKA): *"If a subscriber has an AKMA subscription, the UDM shall include the AKMA indication and Routing indicator in the Nudm_UEAuthentication_Get Response."* TS 29.503 Tab. 6.3.6.2.3-1, `AuthenticationInfoResult`:

- `akmaInd`: *"true: AKMA keys need to be generated for the UE"*;
- `routingId`: *"shall be present if the akmaInd is set to "true""* e contém *"the Routing Indicator contained in the SUCI"*.

`[FREE5GC]` O UDM v1.4.5 não preenchia nenhum dos dois. O `free5gc/openapi` v1.2.4 já tem os campos, tanto no resultado (`UdmUeauAuthenticationInfoResult.AkmaInd`, `RoutingId`) quanto na assinatura lida do UDR (`AuthenticationSubscription.AkmaAllowed`, `RoutingId`).

## 2. De onde vem o RID (decisão D10)

`[3GPP]` TS 29.505 Tab. 5.4.2.2-1, `AuthenticationSubscription` no UDR:

- `akmaAllowed` indica se o assinante pode usar AKMA;
- `routingId` *"shall be present if the akmaAllowed is set to "true""* e *"shall include the Routing Indicator contained in the SUCI"*. O YAML define o padrão `^[0-9]{1,4}$`, uma cadeia de dígitos.

`[MVP]` Decisão D10 (17/10/2026): o RID enviado ao AUSF é o `routingId` guardado no UDR. Ele existe também quando o UDM recebe SUPI em vez de SUCI, por exemplo numa reautenticação depois de registro por GUTI. Monitoramento pedido pelo responsável: quando o pedido traz SUCI, o UDM compara o RID do SUCI com o do UDR e registra no log (`warning`) qualquer diferença. O valor enviado continua sendo o do UDR.

Com `akmaAllowed: true` e `routingId` ausente ou fora do padrão, o dado viola TS 29.505, e o UDM não envia a indicação AKMA. `[MVP]` Isso fica registrado como `error` no log. A autenticação segue normalmente.

## 3. Arquivos tocados

Árvore: `~/free5gc/NFs/udm` (`free5gc/udm` v1.4.5). Binário: `~/free5gc/bin/udm`, com o original em `bin/udm.orig`. *Patch*: `free5gc/patches/udm/0001-akma-indication-and-routing-id.patch`.

| Arquivo | Mudança |
|---|---|
| `pkg/suci/suci.go` | nova função exportada `RoutingIndicator(suci)`: o RID do SUCI como cadeia, com todos os dígitos (TS 23.003 §2.2B) |
| `pkg/suci/suci_test.go` | RIDs `"0"`, `"012"`, `"0000"` e `"12"` (MNC de 3 dígitos) |
| `internal/sbi/processor/akma.go` (novo) | `setAkmaIndication`: preenche `akmaInd` e `routingId` e faz a comparação de D10 |
| `internal/sbi/processor/generate_auth_data.go` | chama `setAkmaIndication` antes de responder, nos dois métodos de autenticação |
| `internal/sbi/processor/akma_test.go` (novo) | casos da função; procedimento completo com UDR simulado, com e sem `akmaAllowed` |
| `internal/logger/logger.go` | categoria `AKMA` (MVP-3a): as mensagens AKMA do UDM saem como `[UDM][AKMA]`; até o MVP-2 saíam como `[UDM][UEAU]` |

## 4. Como reproduzir numa instalação limpa

Com o repositório em `~/akma-free5gc`, ao lado de `~/free5gc`, e o free5GC parado:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm apply
```

Aplica o *patch*, guarda `bin/udm.orig` e recompila. Para desfazer, use `patch-nf.sh udm revert`.

## 5. Como verificar

```bash
cd ~/free5gc/NFs/udm && go test ./internal/sbi/processor/ ./pkg/suci/
```

Com o banco de testes: `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription` (`verificacao.md`). No log do free5GC, `[UDM][AKMA] AKMA indication sent to the AUSF for imsi-…, Routing Indicator …` confirma o envio (MVP-3a, `logs-akma.md`). A ausência de `Routing Indicator of the SUCI ... differs` confirma que o RID do SUCI e o do UDR coincidem.

Resultado em 17/10/2026: os testes do UDM passam. Com UERANSIM (SUCI `suci-0-208-93-0000-…`), o UDR tinha `routingId: "0000"` e não houve divergência registrada.

## 6. `trigger-auth` e a notificação ao AMF (MVP-12)

`[3GPP]` TS 33.535 §6.4.4 e TS 33.501 §6.1.5.2 passos 1 e 2: o AAnF pode pedir ao UDM, com `Nudm_UECM_AuthTrigger`, uma nova autenticação primária do UE; o UDM decide pela sua política (NOTA A) e notifica o AMF com `Nudm_UECM_Re-AuthenticationNotification`. TS 29.503 §5.3.2.16 e §6.2.3.13 (`trigger-auth`), §5.3.2.15 e §6.2.5.6 (notificação). Especificação aprovada: `renovacao-kakma.md` (K2 e K4). Correção F17 de `correcoes-free5gc.md`.

*Patch*: `free5gc/patches/udm/0002-auth-trigger.patch`, aplicado depois do `0001`. Ele usa os modelos e o cliente do `openapi/0003` (`openapi.md`), pela cópia com *patch* em `../../../akma-free5gc/free5gc/build/openapi`, como o AUSF.

| Arquivo | Mudança |
|---|---|
| `go.mod`, `go.sum` | `replace` do `openapi` pela cópia com *patch* |
| `internal/sbi/api_uecontextmanagement.go` | rotas `POST` e `GET` de `/:ueId/registrations/trigger-auth` |
| `internal/sbi/api_authtrigger.go` (novo) | `HandleAuthTrigger`: confere o SUPI do caminho e o do corpo; `GET` recebe 405 com `Allow: POST` (DV-12) |
| `internal/sbi/processor/auth_trigger.go` (novo) | `AuthTriggerProcedure`: assinatura no UDR, registro do AMF, política, notificação ao AMF e mapeamento da resposta |
| `internal/sbi/api_authtrigger_test.go`, `internal/sbi/processor/auth_trigger_test.go` (novos) | pedidos malformados; respostas do AMF; intervalo de 60 s |

Respostas ao consumidor (tabela 6.2.3.13.4.2.2-2):

| Situação | Resposta |
|---|---|
| o AMF aceitou (`204`) ou marcou a autenticação como pendente (`202`) | `204` |
| SUPI sem assinatura no UDR | `404 USER_NOT_FOUND` |
| nenhum AMF registrado para o acesso 3GPP, ou o AMF sem o UE | `404 CONTEXT_NOT_FOUND` |
| outro `trigger-auth` aceito para o mesmo SUPI há menos de 60 s (política, NOTA A) | `403 REAUTHENTICATION_NOT_ALLOWED` |
| registro do AMF sem `reauthNotifyCallbackUri` (AMF sem o *patch*) | `403 REAUTHENTICATION_NOT_ALLOWED` |
| o AMF recusou (`403`, ou `409` de registro ou *handover* em curso, que o `trigger-auth` não lista) | `403 REAUTHENTICATION_NOT_ALLOWED`, com a causa do AMF no `detail` |
| corpo sem `supi`, ou `supi` diferente do `ueId` | `400 MANDATORY_IE_MISSING` ou `MANDATORY_IE_INCORRECT` |
| outra falha do AMF ou da rede | `500 SYSTEM_FAILURE` |

`[LACUNA]` A norma não diz se o UDM responde ao consumidor antes ou depois do AMF. O UDM do projeto espera o AMF (até 10 s), para que o consumidor saiba se a autenticação vai acontecer (K4). Uma recusa do AMF libera o SUPI para outro pedido.

`[3GPP]` A notificação ao AMF vai sem *token* de acesso: TS 29.500 §6.7.3 e TS 33.501 §13.4.1.0 (decisão do responsável de 28/10/2026; F20).

`[FREE5GC]` O UDM guarda o registro do AMF na memória (`UdmUeContext.Amf3GppAccessRegistration`). O UDR, sem o *patch* do `openapi`, não guarda o `reauthNotifyCallbackUri`. Depois de reiniciar o UDM, o `trigger-auth` responde `404 CONTEXT_NOT_FOUND` até o UE se registrar de novo.

Logs na categoria AKMA (`logs-akma.md`):

- `AuthTrigger for <SUPI> from AANF at <endereço>: re-authentication notified to the AMF <id>`, com `, pending until the UE is reachable` quando o AMF responde `202`. O tipo de NF vem do `User-Agent` (TS 29.500 tabela 5.2.2.2-1);
- `AuthTrigger for <SUPI> from … refused: <status> <causa> (<detalhe>)`.

Reprodução: com o `0001` aplicado por uma versão anterior, desfazer com o diretório antigo e aplicar de novo os dois:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm revert <diretório com só o 0001>
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm apply
```

Numa instalação limpa basta o segundo. O `openapi` com o `0003` tem de estar preparado antes (`free5gc/scripts/prepare-openapi.sh`).

Verificação: `cd ~/free5gc/NFs/udm && go test ./internal/sbi/...`, e o passo R6r do `e2e.sh` (`docs/af-side/e2e.md`). Resultado em 28/10/2026: os testes passam, também com `-race`, e o R6r passou inteiro na execução `20261028-162316` (`renovacao-kakma.md` §11).

