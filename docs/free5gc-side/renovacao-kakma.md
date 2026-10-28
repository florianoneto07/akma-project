# Renovação da K_AKMA pelo AAnF (`Nudm_UECM_AuthTrigger`): proposta

**MVP:** MVP-12 · **Data:** 23 de outubro de 2026 · **Estado:** **aprovada** pelo responsável em 28/10/2026, com os pontos K1 a K12 da seção 8 como propostos; a implementar. As mudanças no free5GC entram em `correcoes-free5gc.md` como F16 a F19.

Pedido do responsável em 22/10/2026: implementar a decisão do defeito D-B do registro de defeitos, a operação `trigger-auth` com o método POST, num MVP próprio, depois da fase (d) do MVP-8. Este documento define o que muda no AAnF, no UDM, no AMF e nos modelos do `openapi` para que o AAnF peça uma nova autenticação primária e, com ela, uma nova K_AKMA.

## 1. Situação normativa

`[3GPP]` TS 33.535:

- §6.4.1: *"KAKMA shall be re-keyed by running a successful primary authentication as described in clause 6.1."*
- §6.4.4: *"the AAnF may decide to refresh the KAKMA based on the operator's local authentication policy by sending the Nudm_UECM_AuthTrigger Request message to the UDM. The UDM may further decide whether to trigger the primary authentication as defined in clause 6.1.5 of TS 33.501"*. É opcional (*may*), e a política é local.
- §6.4.2: a K_AF não muda sozinha quando muda a K_AKMA; o AF segue com a K_AF que tem até ela vencer (o que o R3 do `e2e.sh` já verifica).

`[3GPP]` TS 33.501 §6.1.5, a autenticação primária disparada pela rede doméstica, opcional para a rede doméstica e a servidora:

- passo 0b: na `Nudm_UECM_Registration`, o AMF *"shall provide a callback URI within the AMF registration for the UDM to create an implicit subscription"* para a `Nudm_UECM_Re-AuthenticationNotification`;
- passo 1: o AAnF manda `Nudm_UECM_AuthTrigger` com o SUPI; NOTA A: o UDM pode recusar pela própria política. *"as AAnF sets the AF key expiry based on operator's local authentication policy, no frequent AF key expiry can happen and there is no risk of signalling overload"*;
- passo 2: o UDM manda a `Nudm_UECM_Re-AuthenticationNotification` ao AMF com o SUPI;
- passo 3: o AMF decide pela política e pelo estado do UE (*"if the UE is under handover, or if the UE is already under authentication"*); recusa com uma causa ou aceita. Se aceita e não alcança o UE, marca a autenticação como pendente e a faz quando o UE voltar;
- passo 4: o AMF *"starts the primary authentication procedure as defined in clause 6.1.2"*.

§14.2.6 e §14.2.7: as duas operações têm só o SUPI como entrada. §6.9.4.2: depois de uma autenticação primária, as chaves NAS novas vêm do novo K_AMF e entram em uso por um NAS SMC. §6.9.4.4: a troca das chaves AS depois de uma autenticação *"may be used"*.

`[3GPP]` TS 29.503:

- §5.3.2.16 e §6.2.3.13: `trigger-auth`, operação customizada em `{apiRoot}/nudm-uecm/v1/{ueId}/registrations/trigger-auth`, com `AuthTriggerInfo` (o SUPI) no corpo. Respostas: `204`; `403 REAUTHENTICATION_NOT_ALLOWED`; `404 USER_NOT_FOUND` ou `CONTEXT_NOT_FOUND`;
- §5.3.2.15 e §6.2.5.6: o UDM faz `POST` no `reauthNotifyCallbackUri` com `ReauthNotificationInfo` (o SUPI; opcionais `oldGuami` e `ausfRediscoveryInd`). O AMF responde `204`, ou `202` quando o UE não está alcançável e a autenticação fica pendente; recusas `403 REAUTHENTICATION_NOT_ALLOWED`, `404 CONTEXT_NOT_FOUND`, `409`;
- `reauthNotifyCallbackUri` é um atributo opcional de `Amf3GppAccessRegistration` e de `AmfNon3GppAccessRegistration`.

`[DIVERGÊNCIA]` Defeito D-B: a tabela 6.2.3.13.4.1-1 mapeia `trigger-auth` para POST, e o `TS29503_Nudm_UECM.yaml` v1.4.2 declara `get:` com corpo obrigatório. O Anexo A.1 diz que o YAML *"takes precedence ... with respect to the encoding of information elements and methods"*. Seguir a tabela contraria essa regra. Ponto K2.

## 2. Situação no free5GC v4.2.3

`[FREE5GC]` Nenhuma parte da cadeia existe:

- **`openapi`** (v1.2.4 e v1.2.5): sem `AuthTriggerInfo`, sem `ReauthNotificationInfo`, sem o atributo `reauthNotifyCallbackUri` e sem a operação no cliente `Nudm_UEContextManagement`;
- **UDM**: sem a rota `trigger-auth`. Ele guarda na memória o registro do AMF (`UdmUeContext.Amf3GppAccessRegistration`) e já notifica o AMF pelo `deregCallbackUri` (`SendOnDeregistrationNotification`), o modelo a seguir. O perfil do UDM no NRF não tem `allowedNfTypes`: o NRF dá ao AAnF um *token* `nudm-uecm`;
- **AMF**: registra-se no UDM só com o `deregCallbackUri` e não tem a rota da notificação. A máquina de estados do GMM já tem a transição de `Registered` para `Authentication`, mas só a usa num `Registration Request` de atualização. `AuthenticationProcedure` pula a autenticação quando o contexto de segurança vale e, na falha, manda `Registration Reject`. Depois do SMC, o estado `ContextSetup` espera um `Registration Request` para responder;
- **AAnF**: só registra-se no NRF; não descobre outras NF, não pede *token* e não guarda quando cada contexto foi registrado.

**UE (UERANSIM).** Sem mudança prevista. O `nr-ue` aceita o `Authentication Request` em qualquer estado (`receiveAuthenticationRequest`). O `Security Mode Command` põe o contexto novo em uso, e o *patch* 0001 deriva ali a nova K_AKMA e o novo A-KID (`deriveAkmaKeys`). O AUSF, pelo *patch* do MVP-2, registra o novo contexto no AAnF, que substitui o anterior (TS 33.535 §6.1 NOTA 1a). Conferir na implementação.

## 3. Fluxo proposto

1. Um AF pede a K_AF de um A-KID (`Naanf_AKMA_ApplicationKey_Get`, direto ou pelo NEF). O AAnF responde como hoje. Se o contexto tem mais tempo que `kakmaRefresh` (K3), ele também pede ao UDM, em segundo plano, a renovação: descobre o UDM no NRF, obtém o *token* `nudm-uecm` e manda `POST .../{supi}/registrations/trigger-auth` com `{"supi": …}`.
2. O UDM aplica a sua política (K4), acha o registro do AMF do acesso 3GPP e faz `POST` no `reauthNotifyCallbackUri` com `ReauthNotificationInfo`. Responde ao AAnF com o resultado do AMF.
3. O AMF (K6), com o UE em `Registered` e conectado, responde `204` e começa a autenticação primária: `Authentication Request`, AUSF, `Authentication Response`, confirmação no AUSF. O AUSF deriva a nova K_AKMA e a registra no AAnF.
4. O AMF faz o NAS SMC com o novo K_AMF (TS 33.501 §6.9.4.2) e volta a `Registered`, sem o `ContextSetup` de um registro. A sessão PDU continua.
5. No `Security Mode Command`, o ME deriva a nova K_AKMA e o novo A-KID. Os pedidos de K_AF seguintes vêm com o A-KID novo; o AF segue aceitando a K_AF antiga até ela vencer (§6.4.2).

## 4. Componentes

| Componente | Mudança |
|---|---|
| `free5gc/patches/openapi` | modelos `AuthTriggerInfo` e `ReauthNotificationInfo`; `reauthNotifyCallbackUri` em `Amf3GppAccessRegistration` e `AmfNon3GppAccessRegistration`; operação `TriggerAuth` (POST) no cliente `Nudm_UEContextManagement` |
| `free5gc/patches/udm` (novo *patch*) | rota `POST /{ueId}/registrations/trigger-auth`; política (K4); notificação ao AMF |
| `free5gc/patches/amf` (novo diretório) | `reauthNotifyCallbackUri` no registro no UDM; rota `POST /namf-callback/v1/reauth/{supi}`; autenticação disparada pela rede: estado do UE, *flag* de pendência, NAS SMC e volta a `Registered` |
| `free5gc/aanf` | instante de registro de cada contexto; `kakmaRefresh` na configuração; descoberta do UDM e *token* no NRF; cliente `trigger-auth` |
| `free5gc/aanf/tools/seed` | comando `trigger-auth --supi`, no papel do AAnF, para os testes do UDM |
| `free5gc/scripts/patch-nf.sh` | alvo `amf` (já cobre qualquer NF em `NFs/`) |
| `af/scripts/e2e.sh`, `e2e-remote.sh` | passo novo (K11) |

## 5. Logs

Na categoria AKMA (D16), no nível *info*, sem chaves:

- AAnF: `K_AKMA of <SUPI> registered <t> ago: refresh requested from the UDM (TS 33.535 cl. 6.4.4)` e o resultado;
- UDM: `AuthTrigger for <SUPI> from <tipo de NF>: re-authentication notified to the AMF <id>` ou a recusa e a causa;
- AMF: `Re-authentication of <SUPI> requested by the UDM: started` ou `pending (UE not connected)` ou a recusa;
- AUSF e UE: os logs de hoje, com o A-KID novo.

## 6. Verificação prevista

| Caso | Esperado |
|---|---|
| K_AKMA mais velha que `kakmaRefresh`, pedido de K_AF | o AF recebe a K_AF; o AAnF pede o `trigger-auth`; o UDM notifica o AMF; o UE se autentica de novo e tem um A-KID novo, igual ao do AUSF e ao do AAnF |
| sessão PDU depois da renovação | continua: o Ua\* do UE com o A-KID novo funciona pelo plano de usuário |
| AF com a K_AF anterior | aceita até vencer (§6.4.2) |
| segundo pedido do AAnF para o mesmo contexto | nenhum `trigger-auth` novo |
| `trigger-auth` repetido dentro do intervalo do UDM | `403 REAUTHENTICATION_NOT_ALLOWED` |
| SUPI desconhecido; UE sem registro de AMF; corpo sem `supi` | `404 USER_NOT_FOUND`; `404 CONTEXT_NOT_FOUND`; `400` |
| `GET` em `trigger-auth` | `405` (K2) |
| UE não conectado | `202` e autenticação no retorno do UE; nos testes de unidade do AMF, porque o UERANSIM não fica em CM-IDLE no laboratório |
| `kakmaRefresh` desligado (padrão) | nada muda: o `e2e.sh` atual passa igual |
| logs | nenhuma chave |

## 7. Limitações

- Só acesso 3GPP: o laboratório não tem não-3GPP.
- Sem a troca das chaves AS (TS 33.501 §6.9.4.4, opcional): o UERANSIM não tem segurança AS. O gNB segue com o K_gNB do K_AMF anterior até o próximo estabelecimento do contexto AS.
- `oldGuami` e `ausfRediscoveryInd` não são usados: o laboratório tem um AMF e um AUSF.

## 8. Pontos para aprovação

| Nº | Ponto | Proposta | Base |
|---|---|---|---|
| K1 | Escopo | a cadeia inteira: AAnF, UDM, AMF e modelos do `openapi`, só no acesso 3GPP; o UE sem mudança | TS 33.535 §6.4.4; TS 33.501 §6.1.5 |
| K2 | Método do `trigger-auth` | POST, como a tabela 6.2.3.13.4.1-1 e o defeito D-B decidiram; `GET` recebe `405`. Como contraria a precedência do Anexo A.1, registrar como DV-12 em `desvios.md` | `[DIVERGÊNCIA]` TS 29.503 tabela 6.2.3.13.4.1-1 e §5.3.2.16.2 (*"POST request"*) contra o YAML e o Anexo A.1; TS 29.501 §4.6.1.2; RFC 9110 |
| K3 | Política do AAnF | `akma.kakmaRefresh` (duração; `0`, o padrão, desliga). Num pedido de K_AF de um contexto mais velho que isso, o AAnF responde e pede o `trigger-auth` em segundo plano, uma vez por contexto; se o contexto não mudar em 5 minutos, pede de novo. Alternativas: (b) um temporizador por contexto, que renova também quem não usa AKMA; (c) só por comando de operação | `[LACUNA]` §6.4.4 deixa a política ao operador. A NOTA A de TS 33.501 §6.1.5 liga a renovação ao uso das chaves dos AF, como em (a) |
| K4 | Política do UDM | aceita e notifica o AMF; `403 REAUTHENTICATION_NOT_ALLOWED` para outro pedido do mesmo SUPI em menos de 60 s; `404 USER_NOT_FOUND` sem assinatura; `404 CONTEXT_NOT_FOUND` sem AMF registrado; `400` sem `supi`. Responde ao AAnF depois da resposta do AMF: `204` se o AMF aceitou (`204` ou `202`), a recusa do AMF repassada | `[LACUNA]` a norma não diz se o UDM espera o AMF; NOTA A (recusa pela política); D-B, observação secundária (`AuthTriggerInfo` sem `required`) |
| K5 | *Callback* do AMF | `reauthNotifyCallbackUri` = `<AMF>/namf-callback/v1/reauth/<SUPI>`, no registro do acesso 3GPP, como o `deregCallbackUri` | TS 33.501 §6.1.5 passo 0b; TS 29.503 §6.2.6.2.2 |
| K6 | AMF ao receber a notificação | UE em `Registered` e conectado: `204` e autenticação já; UE registrado e não conectado: `202`, autenticação pendente, feita depois do próximo `Service Request` ou registro aceito; UE em outro procedimento (autenticação, SMC, registro, desregistro): `403 REAUTHENTICATION_NOT_ALLOWED`; UE desconhecido: `404 CONTEXT_NOT_FOUND`. A falha da autenticação segue o tratamento atual do AMF | TS 33.501 §6.1.5 passo 3; TS 29.503 tabela 6.2.5.6-3 |
| K7 | Depois da autenticação | NAS SMC com o novo K_AMF e volta direta a `Registered`; sem troca das chaves AS | TS 33.501 §6.9.4.2 (*shall*) e §6.9.4.4 (*may*); `[MVP]` o UERANSIM não tem segurança AS |
| K8 | Modelos | no *patch* do `openapi` do projeto, como os do Naanf_AKMA; `oldGuami` e `ausfRediscoveryInd` lidos e ignorados | TS 29.503 §6.2.6.2.35, §6.2.6.2.36 |
| K9 | AAnF como consumidor | descobre o UDM no NRF e pede o *token* `nudm-uecm` quando o OAuth2 está ligado, como o AUSF faz | TS 29.510; TS 33.501 §13.4.1 |
| K10 | Ferramenta de teste | `seed trigger-auth --supi`, no papel do AAnF, para os casos do UDM sem esperar a política | como o `seed register` do MVP-1 |
| K11 | Verificação | passo novo do `e2e.sh`, depois do R6: o AAnF reiniciado com `kakmaRefresh: 20s`, um Ua\* do UE depois disso e as conferências da seção 6; a configuração do AAnF volta ao fim | — |
| K12 | Defeitos novos | registrar como D-O: TS 29.503 §5.3.2.16.2 passo 2a diz *"the NF Service Consumer responds with 204"*, quando quem responde é o UDM; e *"manadatory"* na nota da tabela 6.2.3.13.4.2.2-2, como em TS 29.535 (D-I) | editorial |

## 9. Arquivos previstos

- `free5gc/patches/openapi`, `free5gc/patches/udm`, `free5gc/patches/amf` (novo);
- `free5gc/aanf` (contexto, configuração, consumidor) e `free5gc/aanf/tools/seed`;
- `free5gc/config/aanfcfg.yaml` (`kakmaRefresh: 0`, comentado);
- `af/scripts/e2e.sh`, `e2e-remote.sh`;
- `docs/free5gc-side/` (este documento, `aanf.md`, `udm.md`, um `amf.md` novo), `desvios.md` (DV-12), registro de defeitos (D-O), plano, matriz, tutorial e guia público.
