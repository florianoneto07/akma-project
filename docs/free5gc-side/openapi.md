# `free5gc/openapi`: modelos e clientes do Naanf_AKMA e do CAPIF

**MVP:** MVP-1; *patch* 0003 no MVP-12; *patch* 0004 no MVP-9; *patches* 0005 e 0006 no MVP-11 · **Decisão:** D6 do plano; ponto C4 do MVP-9 · **Data:** 17 de outubro de 2026; 0003 em 28 de outubro de 2026; 0004 em 29 de outubro de 2026; 0005 em 4 de novembro de 2026; 0006 em 5 de novembro de 2026

## 1. O que falta no `openapi` do free5GC

`[FREE5GC]` O `free5gc/openapi` v1.2.4, fixado pelo free5GC v4.2.3, já traz o tipo de NF `AANF`, o nome de serviço `naanf-akma` e o `AanfInfo` (`docs/baseline/ambiente-verificado.md` §3). Não traz os tipos de dados do Naanf_AKMA nem um cliente para ele. Sem os tipos não há produtor (AAnF) nem consumidores (AUSF no MVP-2, AF no MVP-3).

`[3GPP]` Os tipos vêm de dois lugares. `AkmaKeyInfo` e `CtxRemove` estão em TS 29.535 §5.1.6.2. `AkmaAfKeyRequest`, `AkmaAfKeyData` e `ServiceDisableNotif` estão em TS 29.522 §5.14.5.3. Os nomes JSON seguem `TS29535_Naanf_AKMA.yaml` v1.2.0 (TS 29.535 Anexo A.2), que prevalece sobre o corpo do texto na codificação (Anexo A.1): `aKId`, `kAkma`, `afId`, `kaf`, `expiry`, `anonInd`, `notifUri`, `suppFeat`.

**CAPIF (MVP-9).** `[FREE5GC]` O `openapi` v1.2.4 não tem nenhum modelo nem cliente das APIs de TS 29.222. `[3GPP]` Os tipos estão em TS 29.222 V19.8.0, §8.1.4, §8.2.4, §8.4.4, §8.5.4 e §9.1.4, e os nomes JSON seguem os YAML que vêm com essa versão (`TS29222_*.yaml`, API v1.4.0; o `externalDocs` deles cita a V19.5.0, a última em que mudaram). A especificação do CCF, do NEF e do AF está em `capif.md`.

## 2. Arquivos

Os *patches* 0001, 0002 e 0004 só acrescentam arquivos. O 0003 (MVP-12) também altera três arquivos existentes: os dois modelos de registro do AMF e o `client.go` do `Nudm_UEContextManagement`. O 0005 (MVP-11) altera `client.go`, `util.go` e `udr/ImsDataRepository/client.go`.

| *Patch* (repositório) | Arquivos acrescentados ao `openapi` | Conteúdo |
|---|---|---|
| `free5gc/patches/openapi/0001-naanf-akma-models.patch` | `models/model_akma_key_info.go`, `model_ctx_remove.go`, `model_akma_af_key_request.go`, `model_akma_af_key_data.go`, `model_service_disable_notif.go` | os cinco tipos de dados |
| `free5gc/patches/openapi/0002-naanf-akma-client.patch` | `aanf/AKMA/{client,configuration,api_common}.go` e um arquivo por operação | cliente das três operações, no estilo dos clientes gerados do free5GC |
| `free5gc/patches/openapi/0003-home-network-triggered-authentication.patch` | `models/model_auth_trigger_info.go`, `model_reauth_notification_info.go`; `udm/UEContextManagement/api_trigger_the_primary_re_authentication.go`, `api_amf_registration_reauth_notification.go`. Altera `models/model_amf3_gpp_access_registration.go`, `model_amf_non3_gpp_access_registration.go` e `udm/UEContextManagement/client.go` | a autenticação primária pedida pela rede doméstica (MVP-12, correção F16): `AuthTriggerInfo`, `ReauthNotificationInfo`, o atributo `reauthNotifyCallbackUri` nos dois registros do AMF, a operação `AuthTrigger` e a notificação `Call3GppRegistrationReauthNotificationPost` |
| `free5gc/patches/openapi/0004-capif-models-and-clients.patch` | 24 modelos `models/model_capif_*.go`; `capif/{APIInvokerManagement,DiscoverService,PublishService,Security,AEFSecurity}/`, cada um com `{client,configuration,api_common}.go` e as operações | o CAPIF do MVP-9 (correção F13; `capif.md`, ponto C4): os modelos e os clientes das operações que o CCF, o NEF e o AF usam |
| `free5gc/patches/openapi/0005-sbi-tls.patch` | `sbitls/{sbitls,identity,crl}.go` e testes; `oauth/peer.go`. Altera `client.go`, `util.go` (`GetNFServiceUri`) e `udr/ImsDataRepository/client.go` | o TLS mútuo da SBI do MVP-11 (correções F21 e F24; `tls-sbi.md` §6.2 e §3.3): o perfil TLS, a conferência do certificado, a CRL e o URI pelo nome |
| `free5gc/patches/openapi/0006-access-token.patch` | `oauth/{verify,signer,scopes}.go` e testes. Altera `oauth/oauth.go` (sai o `VerifyOAuth` antigo; sai `oauth/peer.go`) e `sbitls/` | a verificação do *token* e o assinador do NRF do MVP-11 (correções F5, F12 e F23; `tls-sbi.md` §7.2 e §7.4) |
| `free5gc/patches/openapi/0007-capif-tls.patch` | altera `sbitls/{sbitls,crl}.go` e os testes, e `client.go` | a idade máxima da CRL de uma CA extra, `sbitls.CrlChecker` e `CloseIdleConnections`, da fase (c) do MVP-11 (`tls-sbi.md` §17) |

**0003.** `[3GPP]` TS 29.503 §6.2.6.2.35 e §6.2.6.2.36 (tipos), §6.2.6.2.2 e §6.2.6.2.3 (atributo), §6.2.3.13 (`trigger-auth`) e §6.2.5.6 (notificação); `TS29503_Nudm_UECM.yaml` v1.4.2. O campo novo vai no fim de cada modelo de registro, para o *patch* não mexer no alinhamento dos outros. `AuthTrigger` usa POST, como a tabela 6.2.3.13.4.1-1, e não o `get` do YAML (DV-12, defeito D-B). A notificação devolve `Pending` para o `202` do AMF. As duas operações tratam os códigos de erro do YAML e o 307/308 como os clientes do 0002.

**0004.** `[3GPP]` TS 29.222 V19.8.0 e os YAML da mesma versão.

- **Nomes.** Todos os modelos levam o prefixo `Capif`, como o free5GC prefixa os das outras APIs do NEF. Sem o prefixo, `AccessTokenReq`, `AccessTokenErr` e `AccessTokenClaims` colidiriam com os do NRF que o `openapi` já tem. `WebsockNotifConfig`, de TS 29.122, é o do `openapi`.
- **Atributos.** Cada modelo tem os atributos da sua tabela, menos os das *features* que o CCF não suporta (`CAPIF_Ext1`, `RNAA`, `SliceBasedAPIExposure`, `EdgeApp_2`, `ApiStatusMonitoring`) e os da interconexão de CCF; o comentário de cada tipo diz quais ficaram de fora. Na leitura, o Go ignora os atributos que o modelo não tem.
- **`CapifAccessTokenClaims`.** As *claims* do *token* do CCF, com os métodos de `jwt.Claims` para o `golang-jwt` assinar e conferir. Traz `iss` e `client_id` com o identificador do *invoker* e `exp` como instante (DV-13). Diferente do modelo de *claims* do NRF (F5), a biblioteca acha o `exp` e recusa o *token* vencido.
- **Clientes.** Um pacote por YAML em `capif/`, com os serviços nomeados pelas *tags* do YAML, como os gerados: `OnboardedAPIInvokersCollectionApi` e `IndividualOnboardedAPIInvokerDocumentApi`; `AllPublishedServiceAPIsCollectionApi`; `APFPublishedAPIsCollectionApi` e `IndividualAPFPublishedAPIDocumentApi`; `IndividualTrustedAPIInvokerStoreApi`; `CheckAuthenticationRequestApi`. São só as operações de `capif.md` §5 e §6.3: a entrada e a saída, a descoberta, a publicação (criar, listar, ler e retirar), o método de segurança, os dados do *invoker*, o *token* e o `check-authentication`.
- **Erros.** Os códigos da tabela 5.2.6-1 de TS 29.122 e o 414 dos GET de TS 29.222 devolvem o `ProblemDetails`; 307 e 308 devolvem o `Location` (TS 29.222 §7.7). No *token*, 400 e 401 devolvem `CapifAccessTokenErr` (tabela 8.5.2.3.4.4.2-3).
- **Credenciais.** A credencial de entrada vai como `Authorization: Bearer`, pelo `openapi.ContextAccessToken` do contexto (TS 29.222 §5.5.2.2.2 NOTA 4). O pedido de *token* vai em `application/x-www-form-urlencoded`, com `client_secret` no corpo ou com `Authorization: Basic` pelo `openapi.ContextBasicAuth` (§5.6.2.3.2). O TLS mútuo é o do `http.Client` que a NF põe na configuração (`SetHTTPClient`).

O cliente do 0002 segue a Tab. 5.1.4.1-1 de TS 29.535: as três operações são `POST` em `{apiRoot}/naanf-akma/v1/` com os sufixos `register-anchorkey`, `retrieve-applicationkey` e `remove-context`. Ele trata:

- os 11 códigos de erro obrigatórios para `POST` (TS 29.500 Tab. 5.2.7.1-1), devolvendo o `ProblemDetails` recebido;
- 307 e 308 com `Location` e `3gpp-Sbi-Target-Nf-Id` (TS 29.535 Tab. 5.1.4.x.2-3/-4);
- o 204 de `retrieve-applicationkey` (campo `NoContent`), embora o AAnF deste projeto nunca o emita (DV-7).

## 3. Como o código usa a cópia com *patch*

O módulo do Go em `~/go/pkg/mod` é somente leitura e compartilhado por todas as NFs; ele não é alterado. `free5gc/scripts/prepare-openapi.sh`:

1. baixa `github.com/free5gc/openapi@v1.2.4` com `go mod download`;
2. copia o módulo para `free5gc/build/openapi`, fora do controle de versão;
3. aplica `free5gc/patches/openapi/*.patch` em ordem.

O `go.mod` do AAnF aponta para essa cópia com `replace github.com/free5gc/openapi => ../build/openapi`. O AUSF corrigido no MVP-2 e o UDM corrigido no MVP-12 (`udm/0002`) usam a mesma cópia, por `../../../akma-free5gc/free5gc/build/openapi`. No MVP-9, o CCF a usa como o AAnF, o NEF pelo mesmo `replace` do *patch* `nef/0001`, e o AF por `../free5gc/build/openapi`.

O banco de testes do free5GC (`~/free5gc/test`) e o AMF fixam outra versão, `v1.2.5-0.20260527003827-02dc71b4d94f`. Os mesmos *patches* aplicam sobre ela sem mudança. O AMF corrigido no MVP-12 (`amf/0001`) usa essa cópia, por `../../../akma-free5gc/free5gc/build/openapi-test`. O segundo argumento do *script* escolhe a versão:

```bash
free5gc/scripts/prepare-openapi.sh free5gc/build/openapi-test v1.2.5-0.20260527003827-02dc71b4d94f
```

`[MVP]` No trabalho *upstream* os arquivos entram no próprio `free5gc/openapi`, e o `replace` desaparece.

## 4. Como reproduzir numa instalação limpa

Na máquina em que o AAnF será compilado, com Go 1.26.2 no `PATH` e este repositório clonado:

```bash
free5gc/scripts/prepare-openapi.sh
```

A saída lista os quatro *patches* aplicados e termina com `patched free5gc/openapi v1.2.4 ready at …/free5gc/build/openapi`. O *script* só substitui um diretório que já seja uma cópia do `free5gc/openapi`.

## 5. Como verificar

```bash
cd free5gc/aanf && go test -run TestOpenapiClient ./internal/sbi/
```

O teste usa o cliente do *patch* contra o servidor do AAnF sobre HTTP/2 sem TLS (h2c). Ele registra o par do vetor `testbed-ue`, recupera a K_AF esperada pelo oráculo e remove o contexto. Depois confere que os erros chegam ao cliente com o `ProblemDetails` do AAnF: 403 `K_AKMA_NOT_PRESENT` e 404 `AKMA_CONTEXT_NOT_FOUND`.

**0004** (29/10/2026). O *patch* aplica nas duas versões do `openapi` (`v1.2.4` e `v1.2.5-0.20260527003827-02dc71b4d94f`), e `go build ./...` e `go vet ./capif/... ./models/` passam nas duas. Um teste fora do repositório usou os cinco clientes contra servidores `httptest`: a credencial em `Authorization: Bearer`, o `Location` dos `201`, os filtros da descoberta, o *token* em formulário, com `client_secret` e com `Basic`, e o `invalid_client` decodificado; e assinou e conferiu as *claims* com `RS512`, com o *token* vencido recusado (`jwt.ErrTokenExpired`). Com as cópias novas, os testes do AAnF (com `-race`) e do AF continuam passando. Os testes que ficam no repositório vêm com o CCF, o NEF e o AF, que usam esses clientes uns contra os outros.
