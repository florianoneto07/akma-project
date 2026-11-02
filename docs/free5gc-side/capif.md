# CAPIF no free5GC — CCF, NEF como APF e AEF, AF como *API invoker*: especificação

**MVP:** MVP-9 · **Data:** 29 de outubro de 2026 · **Estado:** **aprovada** pelo responsável em 29/10/2026, com os pontos C17 a C30 da seção 13 como propostos; o DV-15 e os defeitos D-AD a D-AG registrados na mesma data. O escopo foi aprovado no plano em 29/10/2026, com os pontos C1 a C16 (plano, seção 4).

Este documento detalha o que o plano aprovou: o CAPIF Core Function (CCF), uma NF nova escrita pelo projeto; o NEF do free5GC como APF e AEF; o AF externo como *API invoker*; e os testes. Os pontos ainda abertos continuam a numeração do plano, de C17 em diante. Neste texto, "função de gestão de API" é a *API management function* do CAPIF, e não o AMF do 5GC.

## 1. Situação normativa

O plano (MVP-9, "Situação normativa") resume as normas. As operações que o MVP usa, na ordem em que acontecem:

| Etapa | Operação | Cláusulas | Interface |
|---|---|---|---|
| 1 | o NEF publica `3gpp-akma` | TS 29.222 §5.3.2.2 e §8.2; TS 23.222 §8.3 | CAPIF-4 |
| 2 | o AF entra no CCF | TS 29.222 §5.5.2.2 e §8.4; TS 33.122 §6.1; TS 23.222 §8.1 | CAPIF-1e |
| 3 | o AF descobre a API | TS 29.222 §5.2.2.2 e §8.1; TS 23.222 §8.7; TS 33.535 §6.3, NOTA do passo 1 | CAPIF-1e |
| 4 | o AF negocia o método de segurança | TS 29.222 §5.6.2.2 e §8.5; TS 33.122 §6.3.1.2 | CAPIF-1e |
| 5 | o AF pede ao NEF que confira se tem os dados para autenticá-lo | TS 29.222 §6.2.2.2 e §9.1; TS 29.522 §7.2 | CAPIF-2e |
| 6 | o NEF obtém no CCF os dados do *invoker* | TS 29.222 §5.6.2.4; TS 23.222 §8.16.3, passo 2a | CAPIF-3 |
| 7 | o AF obtém o *token* | TS 29.222 §5.6.2.3.2; TS 33.122 §6.5.2.3 e Anexo C | CAPIF-1e |
| 8 | o AF pede a K_AF ao NEF, com o *token* | TS 29.522 §4.4.23.2 e §7.2; TS 33.535 §6.3 | CAPIF-2e |
| 9 | o AF sai do CCF | TS 29.222 §5.5.2.3; TS 33.122 §6.8 | CAPIF-1e |

`[3GPP]` Regras comuns às APIs do CAPIF (TS 29.222 §7):

- erros conforme TS 29.122 §5.2.6 (§7.7): `ProblemDetails` em `application/problem+json`, com os códigos da tabela 5.2.6-1;
- negociação de *features* conforme TS 29.122 §5.2.7 (§7.8): `supportedFeatures` vai no pedido de criação e na resposta;
- HTTP/1.1 sobre TLS é obrigatório, e HTTP/2 sobre TLS, recomendado (§7.3).

`[DIVERGÊNCIA]` TS 29.222 §7.3: *"A functional entity desiring to use HTTP/2 shall use the HTTP upgrade mechanism to negotiate applicable HTTP version as described in IETF RFC 9113 [10]."* A RFC 9113 §3.1 diz que o uso do `h2c` no *Upgrade* do HTTP *"is deprecated by this document"*, e §3.2 manda negociar o HTTP/2 sobre TLS por ALPN. Igual na V20.1.0. O CCF segue a RFC 9113 (ponto C29).

## 2. Situação no free5GC e no projeto

- `[FREE5GC]` O free5GC não tem CCF, e o `openapi` v1.2.4 não tem nenhum modelo das APIs do CAPIF (F13; plano, C4).
- `[FREE5GC]` No NEF com o *patch* `nef/0001` (MVP-6), o grupo `3gpp-akma` passa pelo *middleware* OAuth2 do free5GC. Ele confere o *token* do NRF (`NefContext.AuthorizationCheck`, `oauth.VerifyOAuth`) e não pede *token* nenhum quando o NRF não exige OAuth2. Depois, a política `akma.authorizedAfs` liga cada AF_ID a endereços de origem (`nef.md` §3.2).
- O AF externo do MVP-6 se registra no NRF pelo túnel SSH e pede o *token* `3gpp-akma` (`af/af/internal/nrf`). O cliente do NEF (`af/af/internal/nef`) recebe o *token* por uma interface, `TokenSource`, que o modo CAPIF troca sem mexer no pedido.
- O AAnF não muda (C14; seção 8).

## 3. Topologia e endereços

```
                 VM do AF (.129)                             VM do free5GC (.130)
  AF externo ──(CAPIF-1e) HTTPS, TLS mútuo ───────────► CCF 192.168.16.130:8020
  AF externo ──(CAPIF-2e) http (DV-9) ────────────────► NEF 192.168.16.130:8000
                                                        NEF ──(CAPIF-3/4) http (DV-1) ──► CCF 127.0.0.1:8020
                                                        NEF ──(SBI, sem mudança) ───────► NRF e AAnF 192.168.16.130:8010
```

| Caminho | Endereço | Segurança | Base |
|---|---|---|---|
| AF → CCF (CAPIF-1e) | `https://192.168.16.130:8020` | TLS 1.3 e 1.2; só de servidor na entrada, mútuo depois | C3; TS 33.122 §6.1 e §6.3.1.1 |
| NEF → CCF (CAPIF-3/4) | `http://127.0.0.1:8020` | sem TLS até o MVP-11 | C3; DV-1 |
| AF → NEF (CAPIF-2e) | `http://192.168.16.130:8000`, descoberto no CCF | sem TLS até o MVP-11 | C3; DV-9 |

As duas portas 8020 ficam em endereços diferentes. O servidor externo só serve as APIs do *invoker*; o interno, só as do provedor (publicação e dados do *invoker*). Uma API pedida no servidor errado recebe `404`.

## 4. Identificadores, chaves e certificados

### 4.1 Identificadores (C17)

| Identificador | Quem atribui | Forma | Norma |
|---|---|---|---|
| `apiInvokerId` | o CCF, na entrada | UUID v4 em minúsculas. É também o `onboardingId` do recurso e o `securityId` do *token* | TS 29.222 tab. 8.4.4.2.2-1; §5.6.2.3.2: `{securityId}` *"is the API invoker identifier"* |
| `aefId`, `apfId` | a configuração do CCF e a do NEF (C6) | `nef-aef`, `nef-apf` | tab. 8.2.4.2.4-1; sem os delimitadores do `scope` (tab. 8.5.4.2.6-1) |
| `apiId` (`serviceApiId`) | o CCF, na publicação | UUID v4 | tab. 8.2.4.2.2-1 |

`[LACUNA]` TS 29.222 não fixa a forma de nenhum deles: são cadeias.

### 4.2 Chaves e certificados (C18, revisto em 30/10/2026)

Feitos por `free5gc/scripts/make-capif-certs.sh`, com o OpenSSL, na pasta `cert/` do free5GC, ao lado de `nrf.pem` e `aanf.pem`. Só para teste: as chaves privadas não são cifradas.

`[3GPP]` TS 33.122 §6.3.1.1, em CAPIF-1e: *"Certificate based authentication shall follow the profiles given in 3GPP TS 33.310 [2], subclauses 6.1.3a and 6.1.4a."* As duas seções somam as regras comuns de TS 33.310 §6.1.1.

| Arquivos | O quê | Quem usa |
|---|---|---|
| `ccf-ca.pem`, `ccf-ca.key` | CA do CCF, ECDSA P-384, assinatura com SHA-384. Perfil de TS 33.310 §6.1.4a: `basicConstraints` crítico, com `CA:TRUE` e `pathlen:0`; `keyUsage` crítico, com `keyCertSign` e `cRLSign` | o CCF assina com ela o próprio certificado e os dos *invokers* (TS 33.122 §6.1 passo 4); o AF a recebe na informação de entrada |
| `ccf.pem`, `ccf.key` | certificado de servidor do CCF, ECDSA P-256, com `subjectAltName` `IP:192.168.16.130`. Perfil de TS 33.310 §6.1.3a: `keyUsage` crítico, com `digitalSignature`; `extendedKeyUsage` `serverAuth`, não crítico | CAPIF-1e |
| `ccf-token.key`, `ccf-token.pub.pem` | ECDSA P-256, para o `ES256` do *token* (C9) | o CCF assina; o NEF confere com a chave pública |
| `ccf-credential.key` | ECDSA P-256, para o `ES256` da credencial de entrada (C19) | só o CCF |

Os nomes seguem TS 33.310 §6.1.1: `O=AKMA test, CN=<nome>`, em `UTF8String` (*"Organization and CN shall be in UTF8 format"*). O CN do servidor é `ccf`.

O par de chaves do *invoker* é gerado pelo AF (ECDSA P-256), e só a chave pública vai ao CCF. O CCF aceita ECDSA P-256 ou P-384, ou RSA de 2048 a 4096 bits com expoente público de pelo menos 65537. `[3GPP]` TS 33.310 §6.1.1: *"A public key length of at least 384-bit shall be supported"*, *"RSA public exponent shall be no less than 65537"* e *"The security level of the public key used to sign the certificate shall be at least the same as the public keys in the certificate"*. Por esta última regra, a CA é P-384. `[INTERPRETAÇÃO]` O RSA de 4096 bits fica abaixo dos 192 bits do P-384 (NIST SP 800-57 Parte 1, tab. 2); TS 33.310 não compara os níveis.

O certificado que o CCF emite na entrada segue TS 33.310 §6.1.3a:
- nome `O=<o da CA>, CN=<identificador do invoker>` (TS 29.222 tab. 8.4.4.2.5-1);
- `keyUsage` crítico, com `digitalSignature` (e `keyEncipherment`, para chave RSA);
- `extendedKeyUsage` `clientAuth`, não crítico;
- nenhuma outra extensão crítica (§6.1.1), e por isso sem `basicConstraints`;
- validade configurável, um ano por padrão.

`[MVP]` Os certificados não têm o ponto de distribuição de CRL que TS 33.310 §6.1.3a torna obrigatório, e não há CRL nem OCSP: DV-16.

`[LACUNA]` `apiInvokerPublicKey` e `apiInvokerCertificate` são cadeias sem formato (tab. 8.4.4.2.5-1). O projeto usa PEM: `PUBLIC KEY` (*SubjectPublicKeyInfo*) e `CERTIFICATE`.

## 5. O CCF

NF nova em `free5gc/ccf/`, na estrutura das NFs do free5GC, como o AAnF: `cmd`, `pkg/factory`, `pkg/service`, `internal/{context,logger,sbi}`. Não se registra no NRF (C1).

### 5.1 Configuração

`free5gc/config/ccfcfg.yaml`, instalado na pasta `config/` do free5GC:

```yaml
info:
  version: 1.0.0
  description: CCF initial local configuration

configuration:
  ccfName: ccf # issuer of the onboarding credentials
  capif1e: # interface of the API invokers (TS 23.222 cl. 6.4.3), HTTPS (TS 33.122 cl. 6.1, 6.3.1.1)
    bindingIPv4: 192.168.16.130
    port: 8020
    certificate: cert/ccf.pem
    key: cert/ccf.key
  internal: # interface of the API provider, CAPIF-3/4 (TS 23.222 cl. 6.4.6, 6.4.7); http until TLS on the SBI
    bindingIPv4: 127.0.0.1
    port: 8020
  ca: # signs the certificates of the API invokers (TS 33.122 cl. 6.1 step 4)
    certificate: cert/ccf-ca.pem
    key: cert/ccf-ca.key # ECDSA P-384 (TS 33.310 cl. 6.1.1)
    invokerCertLifetime: 8760h
  token: # access tokens (TS 29.222 cl. 5.6.2.3.2; TS 33.122 Annex C)
    key: cert/ccf-token.key # ECDSA P-256, signed with ES256 (TS 33.210 cl. 6.3.3.1)
    lifetime: 10m
  credentialKey: cert/ccf-credential.key # onboarding credentials (TS 33.122 cl. 6.1), ECDSA P-256, ES256
  providers: # API provider domain functions, fixed by configuration
    - name: nef
      apfId: nef-apf
      aefIds:
        - nef-aef
  invokerProfiles: # AFs that may onboard; the onboarding credential names the profile
    - name: af1-external
      apis:
        - 3gpp-akma
      afIds: # AF_IDs the NEF accepts from invokers of this profile
        - af1.akma.example.0100011301
        - af1.akma.example.0100011302
        - af1.akma.example.0100011303
        - af1.akma.example.010001C02B
        - af1.akma.example.010001C02C
        - af1.akma.example.0100000007
        - af1.akma.example.010001D001
        - af1.akma.example.0100000009
        - af1.akma.example.0100000008
        - af1.akma.example.FF00000001
        - af4.akma.example.FF00000001 # test only: outside the policy of the NEF

logger:
  enable: true
  level: info
  reportCaller: false
```

### 5.2 Entrada (`Onboard_API_Invoker`)

`POST https://192.168.16.130:8020/api-invoker-management/v1/onboardedInvokers`, com `Content-Type: application/json` e `Authorization: Bearer <credencial>` (TS 29.222 §5.5.2.2.2, NOTA 4). TLS só de servidor: o CCF não pede certificado de cliente nesta operação (TS 33.122 §6.1 passo 2).

| Atributo de `APIInvokerEnrolmentDetails` | P | No MVP-9 |
|---|---|---|
| `onboardingInformation.apiInvokerPublicKey` | M | PEM; ECDSA P-256 ou P-384, ou RSA de 2048 a 4096 bits com expoente de pelo menos 65537 (seção 4.2) |
| `notificationDestination` | M | guardado; o CCF não manda notificações (DV-14) |
| `apiList` | O | ignorado; a resposta lista as APIs do perfil (C20) |
| `apiInvokerInformation` | O | só para o log |
| `supportedFeatures` | C | obrigatório no pedido: *"This attribute shall be provided in the HTTP POST request"* (tab. 8.4.4.2.2-1); sem ele, `400`. O CCF não tem nenhuma *feature* da tab. 8.4.6-1 e responde `"0"` |
| `apiInvokerId` | — | *"shall not be present"* no pedido: com ele, `400` |
| os demais | O | ignorados (são de *features* que o CCF não tem) |

Na ordem:

1. **Credencial** (C19). Ausente, malformada, com assinatura errada, vencida, com `aud` errado ou de perfil desconhecido: `401`, com `ProblemDetails` e `WWW-Authenticate: Bearer`, mais `error="invalid_token"` quando havia credencial.
2. **Corpo.** `Content-Type` que não é JSON: `415`. JSON inválido, atributo M ausente ou chave pública ilegível ou de tipo não aceito: `400`, com `invalidParams`.
3. **Entrada automática** (TS 29.222 §5.5.2.2.2, passo 1 e NOTA 1). O CCF cria o perfil do *invoker* com o identificador novo, o nome do perfil de AF, as APIs e os AF_ID do perfil; assina o certificado; gera o `onboardingSecret`, 32 octetos aleatórios em base64url sem `=` (TS 33.122 §6.1 passo 4).
4. **Resposta** `201 Created`, com `Location: https://192.168.16.130:8020/api-invoker-management/v1/onboardedInvokers/{apiInvokerId}` e o corpo `APIInvokerEnrolmentDetails`: `apiInvokerId`; `onboardingInformation` com a chave pública do pedido, `apiInvokerCertificate` e `onboardingSecret`; `notificationDestination`; `apiList` com as `ServiceAPIDescription` publicadas das APIs do perfil, ou sem `apiList` quando nenhuma está publicada; `supportedFeatures` `"0"`.

A mesma credencial serve a mais de uma entrada enquanto vale; cada entrada cria um *invoker* novo (C19).

### 5.3 Saída (`Offboard_API_Invoker`)

`DELETE https://192.168.16.130:8020/api-invoker-management/v1/onboardedInvokers/{onboardingId}`, com TLS mútuo:

- sem certificado de cliente, com certificado que não é da CA do CCF ou de outro *invoker* (CN diferente de `{onboardingId}`): `403`;
- `{onboardingId}` desconhecido: `404`;
- senão, o CCF apaga o perfil, o contexto de segurança e o segredo (TS 33.122 §6.8 passo 3) e responde `204`. Os *tokens* já emitidos valem até `exp` (DV-14);
- depois do `204`, o CCF encerra a sessão TLS (TS 33.122 §6.8 passo 6: *"The CAPIF core function shall tear down the TLS session with the API invoker"*): com `Connection: close` no HTTP/1.1, e com `GOAWAY` no HTTP/2.

### 5.4 Publicação (`Publish_Service_API`, `Unpublish_Service_API`, `Get_Service_API`)

No servidor interno:

- `POST http://127.0.0.1:8020/published-apis/v1/{apfId}/service-apis`, com `ServiceAPIDescription`;
- `GET .../{apfId}/service-apis` e `GET .../{apfId}/service-apis/{serviceApiId}`;
- `DELETE .../{apfId}/service-apis/{serviceApiId}`.

Regras:

- `{apfId}` fora da configuração: `403`;
- no POST: `apiName` presente; `apiId` ausente; `supportedFeatures` presente (NOTA 1 da tab. 8.2.4.2.2-1); `aefProfiles` com pelo menos um perfil; cada `aefId` do mesmo provedor do `{apfId}`; cada perfil com pelo menos uma versão e com `securityMethods` no perfil ou em cada interface (NOTA 4 da tab. 8.2.4.2.4-1). Senão, `400` com `invalidParams`;
- a mesma API (`apiName` e `aefId`) publicada de novo pelo mesmo APF substitui a anterior e recebe um `apiId` novo, porque o NEF publica a cada partida (C24);
- respostas: POST `201`, com `Location` e o `apiId` no corpo; GET `200`; DELETE `204`; `{serviceApiId}` desconhecido, `404`;
- o CCF não avisa os *invokers* da publicação (§5.3.2.2.2 passo 2d; DV-14).

### 5.5 Descoberta (`Discover_Service_API`)

`GET https://192.168.16.130:8020/service-apis/v1/allServiceAPIs?api-invoker-id=<id>`, com TLS mútuo e, opcionais, os filtros sem *feature* da tab. 8.1.2.2.3.1-1: `api-name`, `api-version`, `comm-type`, `protocol`, `aef-id`, `data-format`, `api-cat` e `preferred-aef-loc` (este último ignorado: o CCF não tem a localização dos AEF).

- sem certificado de cliente, ou com CN diferente de `api-invoker-id`: `403`; *invoker* desconhecido: `403` (§5.2.2.2.2 passo 1);
- filtros das *features* que o CCF não tem (`req-api-prov-name`, `api-supported-features`, `ue-ip-addr`, `service-kpis`, `net-slice-info`, `grant-types`, `api-ids`, `res-ops`, `supported-features`): `400`, com `invalidParams`;
- resultado: as APIs publicadas do perfil do *invoker* que passam nos filtros (a política de descoberta, §5.2.2.2.2 passo 2b), sem `shareableInfo`. Resposta `200` com `DiscoveredAPIs`; sem nenhuma API, `200` sem o atributo `serviceAPIDescriptions` (cardinalidade 1..N, tab. 8.1.4.2.2-1).

### 5.6 Método de segurança (`Obtain_Security_Method`)

`PUT https://192.168.16.130:8020/capif-security/v1/trustedInvokers/{apiInvokerId}`, com TLS mútuo e o corpo `ServiceSecurity` (tab. 8.5.4.2.2-1):

- sem certificado de cliente, ou com CN diferente de `{apiInvokerId}`: `403`; *invoker* desconhecido: `404`;
- cada item de `securityInfo` com `aefId` ou `interfaceDetails`, nunca os dois (NOTA da tab. 8.5.4.2.3-1), e com `prefSecurityMethods`; `notificationDestination` presente. Senão, `400`;
- para cada item, o CCF escolhe o primeiro método da preferência do *invoker* que o AEF publicou e que a política aceita: só `OAUTH` (C3; DV-14). Sem método comum, `selSecurityMethod` fica ausente (tab. 8.5.4.2.3-1);
- o CCF guarda o contexto de segurança e responde `201 Created`, com `Location` e o corpo com `selSecurityMethod` e `supportedFeatures` `"0"`. Um segundo PUT substitui o contexto.

### 5.7 *Token* (`Obtain_Authorization`)

`POST https://192.168.16.130:8020/capif-security/v1/securities/{securityId}/token`, com TLS mútuo e `Content-Type: application/x-www-form-urlencoded` (NOTA 1 da tab. 8.5.4.2.6-1). Campos: `grant_type=client_credentials`, `client_id`, `client_secret` (o `onboardingSecret`), ou `Authorization: Basic` com os dois (§5.6.2.3.2), e `scope`, opcional.

Na ordem:

1. **Cliente.** Certificado de cliente do *invoker* `{securityId}` e o segredo de entrada certo; `client_id` igual a `{securityId}`. Senão, `invalid_client`: `401`, com `WWW-Authenticate: Basic`, quando o pedido usou `Authorization: Basic`; `400` nos outros casos (tab. 8.5.5.3-1; RFC 6749 §5.2).
2. **Pedido.** `grant_type` diferente de `client_credentials`: `400` `unsupported_grant_type`. Campo obrigatório ausente: `400` `invalid_request`.
3. **Contexto.** Sem o contexto de segurança da seção 5.6, ou sem `OAUTH` escolhido para o AEF do escopo: `400` `unauthorized_client`.
4. **Escopo** (C22). No formato sem `CAPIF_Ext1` e sem `CAPIF_Ext2` (tab. 8.5.4.2.6-1), `3gpp#<aefId>:<apiName>[,<apiName>];...`. Cada par tem de estar no perfil do *invoker* e publicado por aquele AEF; senão, `400` `invalid_scope`. Sem `scope`, o CCF concede todas as APIs do perfil publicadas nos AEF com `OAUTH` escolhido (RFC 6749 §3.3).
5. **Resposta** `200`, com `Cache-Control: no-store` (RFC 6749 §5.1) e `AccessTokenRsp`: `access_token` (seção 5.8), `token_type` `Bearer`, `expires_in` (a vida configurada) e `scope` (o concedido).

Os erros vão no corpo `AccessTokenErr`, sem `ProblemDetails` (tab. 8.5.2.3.4.4.2-3).

### 5.8 Forma do *token* (C9; DV-13)

Cabeçalho JWS `{"alg":"ES256","kid":"<kid>","typ":"JWT"}`, serialização compacta (RFC 7515 §7.1). O `kid` é a impressão digital da chave do *token* (RFC 7638).

`[3GPP]` TS 33.210 §6.3.3.1, o perfil JWS: *"All entities and functions that support JWS according to RFC 7515 [46] shall follow the following restrictions and extensions"*. Entre elas: *"The "alg" header parameter value ES256 (ECDSA using P-256 and SHA-256) shall be supported"*; *"The "kid" header parameter shall be supported"*; *"If the "kid" header parameter is used, the end point shall check that the indicated "alg" matches the one specified by the parameters"*. `[INTERPRETAÇÃO]` TS 33.122 Anexo C só cita a RFC 7515, e o perfil vale para o CCF pela aplicação geral de TS 33.210 §6.3.1. O `x5t#S256`, que §6.3.3.4 recomenda (*should*) no *token*, não é usado: DV-16.

*Claims*:

| *Claim* | Valor | Norma |
|---|---|---|
| `iss` | o identificador do *invoker* | TS 29.222 tab. 8.5.4.2.8-1 (DV-13) |
| `client_id` | o identificador do *invoker* | TS 33.122 C.2.2 (DV-13) |
| `scope` | o escopo concedido, por exemplo `3gpp#nef-aef:3gpp-akma` | tab. 8.5.4.2.8-1; C.2.2 |
| `exp` | instante em segundos desde 1970-01-01 (UTC) | RFC 7519 §4.1.4; TS 33.122 C.6 (DV-13) |

Nenhuma outra *claim*. Nenhum log do CCF, do NEF ou do AF mostra o *token*.

### 5.9 Dados do *invoker* (`Obtain_API_Invoker_Info`)

`GET http://127.0.0.1:8020/capif-security/v1/trustedInvokers/{apiInvokerId}?authenticationInfo=true&authorizationInfo=true`, no servidor interno (tab. 8.5.2.3.3.1-1):

- *invoker* desconhecido, ou sem contexto de segurança: `404`;
- resposta `200` com `ServiceSecurity`, os itens do contexto, cada um com `aefId`, `prefSecurityMethods` e `selSecurityMethod` e, quando pedidos (C21):
  - `authenticationInfo`: o certificado da CA do CCF, em PEM, para o AEF validar o certificado do *invoker* (TS 33.122 §6.5.2.3 passo 5). O NEF o guarda e só o usa no MVP-11;
  - `authorizationInfo`: o texto JSON `{"akmaAfIds":["<AF_ID>", ...]}`, com os AF_ID do perfil (C11).
- `[MVP]` Sem TLS em CAPIF-3 (DV-1), o CCF não sabe qual AEF pergunta e devolve todos os itens; o NEF usa o do seu `aefId`.

### 5.10 Credencial de entrada (C19)

`ccf credential --config ccfcfg.yaml --profile af1-external --validity 24h` escreve a credencial na saída: um JWT `ES256`, com `kid` (seção 5.8), assinado com `ccf-credential.key`. As *claims* são `iss` (o `ccfName`), `sub` (o perfil), `aud` (`api-invoker-management`), `iat` e `exp`. Na entrada, o CCF confere:
- a assinatura e o `kid`;
- `iss`;
- `aud`;
- `exp`, com até 30 s de tolerância;
- se `sub` é um perfil da configuração.

`[3GPP]` TS 33.122 §6.1: *"When the OAuth 2.0 token based mechanism is used as the onboarding credential, the access token shall be encoded as JSON web token as specified in IETF RFC 7519 [6], shall include the JSON web signature as specified in IETF RFC 7515 [7], and shall be validated per OAuth 2.0 [4], IETF RFC 7519 [6] and IETF RFC 7515 [7]."* NOTA 1: como o *invoker* recebe a informação de entrada fica fora da norma. No laboratório, o `e2e.sh` a copia para a máquina do AF.

### 5.11 Erros, *features* e versões do HTTP

- Todos os erros, menos os da seção 5.7, levam `ProblemDetails` em `application/problem+json`, com `status` e `detail`. `cause` só quando a norma define uma causa para o caso: as causas de TS 29.500 tab. 5.2.7.2-1 valem para o CAPIF, por TS 29.222 §7.7 e TS 29.122 §5.2.6 (*"The protocol and application errors in clause 5.2.7.2 of 3GPP TS 29.500 [44] are applicable for the above status codes for the APIs"*). O CCF usa estas: `INVALID_MSG_FORMAT`, `MANDATORY_IE_MISSING`, `MANDATORY_IE_INCORRECT`, `OPTIONAL_IE_INCORRECT`, `INVALID_QUERY_PARAM`, `MANDATORY_QUERY_PARAM_MISSING`, `INVALID_API`, `RESOURCE_URI_STRUCTURE_NOT_FOUND`, `INCORRECT_LENGTH`, `MAX_JSON_SIZE_EXCEEDED` e `SYSTEM_FAILURE`. Os `401`, os `403` e os `404` de recurso desconhecido não levam causa. `API_INVOKER_NOT_ALLOWED` é de `CAPIF_Ext1`. (Corrigido na implementação, seção 15.1.)
- Códigos da tab. 5.2.6-1 de TS 29.122 que se aplicam: `400`, `401`, `403`, `404`, `411` (`Content-Length` inválido ou diferente do corpo), `413` (corpo acima de 64 KiB), `415`, `500` e `503`. Pela NOTA 1 da tabela, também `405`, com o cabeçalho `Allow` (RFC 9110 §15.5.6).
- `supportedFeatures` `"0"` em todas as respostas que o pedem: o CCF não tem nenhuma das *features* das tabelas 8.1.6-1, 8.2.6-1, 8.4.6-1 e 8.5.6-1.
- CAPIF-1e com HTTP/1.1 e HTTP/2 negociados por ALPN (RFC 9113 §3.2; C29). O servidor interno, com HTTP/1.1.

### 5.12 Estado

Os perfis dos *invokers*, os contextos de segurança e as APIs publicadas ficam na memória (C1). Depois de reiniciar, o CCF não conhece mais nenhum: o NEF publica de novo (seção 6.2), e o AF entra de novo (seção 7.2).

## 6. O NEF

*Patch* `free5gc/patches/nef/0002`, aplicado depois do `0001`, e a configuração em `free5gc/patches/config/0003`.

### 6.1 Configuração

Seção nova no `nefcfg.yaml`. Sem ela, o NEF fica como no MVP-6:

```yaml
configuration:
  capif: # CAPIF (TS 29.522 cl. 7): the NEF as API publishing and exposing function
    ccfUri: http://127.0.0.1:8020 # internal interface of the CCF, CAPIF-3/4 (http until TLS on the SBI)
    apfId: nef-apf
    aefId: nef-aef
    tokenPublicKey: cert/ccf-token.pub.pem # checks the access tokens of the CCF (TS 33.122 cl. 6.5.2.3 step 7)
    interface: # published for 3gpp-akma (CAPIF-2e)
      ipv4Addr: 192.168.16.130
      port: 8000
    invokerInfoLifetime: 10m # how long the NEF keeps the data of an API invoker
    allowNrfTokens: true # 3gpp-akma also without CAPIF, with access tokens of the NRF (TS 29.522 cl. 6)
```

### 6.2 Publicação (C24)

Na partida, o NEF publica `3gpp-akma` (seção 5.4):

- `apiName` `3gpp-akma`; `description`; `supportedFeatures` `"0"`;
- `aefProfiles`, com um perfil:
  - `aefId` `nef-aef`;
  - `versions`: `apiVersion` `v1`, com `custOperations` = `commType` `REQUEST_RESPONSE`, `custOpName` `retrieve`, `operations` `[POST]` (TS 29.522 §5.14.1 e tab. 5.14.3.1-1);
  - `protocol` `HTTP_1_1`; `dataFormat` `JSON`; `securityMethods` `[OAUTH]`;
  - `interfaceDescriptions`: `ipv4Addr` 192.168.16.130, `port` 8000, `securityMethods` `[OAUTH]`.

Se o CCF não responde, o NEF tenta de novo a cada 5 s, com um aviso no log só na primeira falha. Enquanto isso, o modo do MVP-6 continua servindo. Ao encerrar, o NEF retira a publicação.

### 6.3 `check-authentication` (AEF_Security_API)

`POST http://192.168.16.130:8000/aef-security/v1/check-authentication`, com `CheckAuthenticationReq` (`apiInvokerId` e `supportedFeatures`, os dois obrigatórios, tab. 9.1.4.2.2-1). O grupo `aef-security` sobe quando há a seção `capif`; ele não é um serviço do NRF e não entra em `serviceList`.

- `[MVP]` O pedido tem de vir de um endereço de `akma.authorizedAfs`, a mesma política de endereço do MVP-6 (DV-9); senão, `403`. TS 29.222 §6.2.2.2 não pede autorização nessa operação.
- Se não tem os dados do *invoker*, ou se passaram de `invokerInfoLifetime`, o NEF os pede ao CCF (seção 5.9) e os guarda.
- Respostas: `200` com `CheckAuthenticationRsp` (`supportedFeatures` `"0"`); *invoker* desconhecido no CCF, `403`; CCF fora do ar, `503`.

### 6.4 Pedido `retrieve` com o *token* do CCF (C25)

Antes do *middleware* do MVP-6, no grupo `3gpp-akma`:

1. **Que *token* é.** `Authorization: Bearer <jwt>` com assinatura que confere com a chave pública do CCF: modo CAPIF, passos 2 a 4. Um `kid` presente tem de ser o da chave (seção 5.8). Qualquer outro pedido, com `allowNrfTokens: true`: o modo do MVP-6, sem mudança, com o *middleware* do NRF (C13). Com `allowNrfTokens: false`: `401`.
2. **Forma.** `alg` `ES256` (TS 33.210 §6.3.3.1); `exp` presente e não vencido, com até 30 s de tolerância (TS 33.122 C.2.2); `iss` e `client_id` presentes e iguais. Senão, `401` `invalid_token`.
3. **Escopo.** `scope` começa com `3gpp#` e tem um item com o `aefId` do NEF e `3gpp-akma` entre as APIs desse item (TS 29.522 §7.2). Senão, `403` `insufficient_scope`.
4. **Dados do *invoker*.** Da memória ou do CCF (seção 5.9; TS 23.222 §8.16.3 passo 2a). O CCF responde `404`: `401` `invalid_token`, porque o *invoker* saiu ou o CCF reiniciou. O CCF fora do ar: `503`.

Depois, o processamento do MVP-6 (`nef.md` §3.2, passos 2 a 6), com uma conferência a mais no passo 3: o AF_ID do pedido tem de estar nos `akmaAfIds` do *invoker* (C11) e na política de endereço. Senão, `403` sem `cause`; o `detail` diz qual política recusou.

Respostas do modo CAPIF, com `ProblemDetails` em `application/problem+json`:

| Caso | Resposta |
|---|---|
| *token* inválido ou vencido (passos 2 e 4) | `401`, `WWW-Authenticate: Bearer error="invalid_token"` |
| sem *token*, com `allowNrfTokens: false` | `401`, `WWW-Authenticate: Bearer` |
| escopo de outro AEF ou de outra API | `403`, `WWW-Authenticate: Bearer error="insufficient_scope", scope="3gpp#nef-aef:3gpp-akma"` |
| AF_ID fora dos `akmaAfIds` ou da política de endereço | `403` sem `cause` |
| CCF fora do ar | `503` |

`[3GPP]` RFC 6750 §3 e §3.1; TS 29.122 tab. 5.2.6-1 e NOTA 2. Com `allowNrfTokens: true`, um pedido sem *token*, ou com um *token* que não é do CCF, cai no modo do MVP-6 e recebe o `401` do *middleware* do free5GC (F12).

## 7. O AF como *API invoker*

### 7.1 Configuração

`af/config/afcfg-capif.yaml`: o `afcfg-external.yaml` sem `nefUri` e `nrfUri`, com a seção `capif`:

```yaml
  capif: # the AF as API invoker (TS 29.522 cl. 7.2); the NEF is discovered at the CCF
    ccfUri: https://192.168.16.130:8020
    ccfCa: ../build/capif/ccf-ca.pem # enrolment information (TS 33.122 cl. 6.1 step 1)
    credential: ../build/capif/af1-external.jwt
    state: ../build/capif/af1-external.json # API invoker ID, certificate, private key and secret
    notificationDestination: https://192.168.16.129/capif/notify # mandatory; not served: the CCF sends no notification
    apiScheme: http # scheme of the discovered apiRoot (TS 29.501 cl. 4.4.1); http while the NEF has no TLS
    offboardOnExit: false
```

### 7.2 Partida e uso (C26)

1. **Entrada.** Sem o arquivo de estado, o AF gera o par ECDSA P-256 e entra no CCF (seção 5.2). Grava o estado com permissão `0600`.
2. **Descoberta** (seção 5.5), com `api-name=3gpp-akma&api-version=v1`. O AF escolhe o primeiro perfil com `OAUTH` e monta o `apiRoot` com `apiScheme`, `ipv4Addr` (ou `fqdn`), `port` e `apiPrefix`, se houver (TS 29.222 §5.2.2.2.2 NOTA).
3. **Método** (seção 5.6), para esse `aefId`, com `[OAUTH]`. Sem `selSecurityMethod` `OAUTH`, o AF não sobe.
4. **`check-authentication`** no NEF (seção 6.3).
5. **Ua\*.** O AF serve os perfis. O *token* (seção 5.7, `scope=3gpp#<aefId>:3gpp-akma`) é pedido no primeiro pedido de K_AF, fica só na memória e é renovado 30 s antes do fim de `expires_in` (TS 33.122 C.4.1 e C.5).

Se o CCF responde que o *invoker* não existe (`403` na descoberta, `404` no método, `invalid_client` no *token*), o AF apaga o estado e volta ao passo 1, uma vez. Com `offboardOnExit: true`, ao encerrar, o AF sai do CCF (seção 5.3) e apaga o estado (TS 33.122 §6.8 passo 5).

Quando o CCF ou o NEF não respondem, o pedido de K_AF falha como no MVP-6 (`aanf.ErrUnavailable`), e o Ua\* responde ao UE o que já responde hoje.

## 8. AAnF

Sem mudança no código, na configuração e nos testes (C14). O CCF não fala com o AAnF. O NEF continua a pedir a K_AF pela SBI, com o *token* `naanf-akma` do NRF, e o AAnF continua a conferir o NEF e os AF_ID externos (`authorizedNefs`, P11 b).

## 9. Logs

No formato do free5GC. Nunca aparecem *tokens*, credenciais de entrada, segredos de entrada, chaves privadas nem K_AF. Os identificadores do *invoker* e do AF ficam no nível *info* (D16).

| Onde | Nível | Mensagens |
|---|---|---|
| CCF | info | `API invoker <id> onboarded, profile <perfil>`; `API invoker <id> offboarded`; `<apiName> published by <apfId> as <apiId>, AEF <aefId>`; `<apiName> <apiId> unpublished`; `security method OAUTH for API invoker <id> at AEF <aefId>`; `access token for API invoker <id>, scope <escopo>, expires in <n> s`; `security information of API invoker <id> given to the internal interface` |
| CCF | warn | cada recusa, com o motivo (credencial, certificado, segredo, escopo, perfil) |
| NEF, categoria `CAPIF` | info | `3gpp-akma published at the CCF as <apiId>`; `check-authentication of API invoker <id>: <n> AF_IDs from the CCF`; `access token of API invoker <id> accepted for AF_ID <afId>` |
| NEF, categoria `CAPIF` | warn | *token* recusado e o motivo; AF_ID fora do perfil do *invoker*; CCF fora do ar |
| AF | info | `onboarded at the CCF as API invoker <id>`; `3gpp-akma discovered at <apiRoot>, AEF <aefId>`; `security method OAUTH at AEF <aefId>`; `the NEF has the authentication data of API invoker <id>`; `access token obtained, expires in <n> s`; `offboarded from the CCF` |

## 10. Verificação prevista

### 10.1 Testes Go

| Onde | Casos |
|---|---|
| CCF | cada operação das seções 5.2 a 5.10, com sucesso e recusas; as *claims* do *token* (`iss` = `client_id` = *invoker*; `exp` absoluto; `scope`); a credencial (assinatura, `aud`, `exp`, perfil); o CN do certificado; as regras da publicação; os filtros e a política da descoberta; os atributos obrigatórios das tabelas em cada corpo |
| NEF (*patch*) | a seção 6.4: *token* válido; de outra chave; `alg` `none` e `HS256`; vencido; sem `exp`; `iss` diferente de `client_id`; escopo de outro AEF e de outra API; *invoker* que saiu (CCF `404`); CCF fora do ar (`503`); AF_ID fora do perfil e fora da política de endereço; `allowNrfTokens` `true` e `false`. A seção 6.3. A publicação e a nova tentativa. Os testes do MVP-6 continuam passando |
| AF | a entrada e o estado; a montagem do `apiRoot`; o método; a renovação do *token*; a nova entrada depois de um CCF reiniciado; a saída |

### 10.2 `e2e.sh`: grupo CP e passo R2c

O `e2e-remote.sh` sobe o CCF com o free5GC (`run.sh -aanf -ccf`) e emite a credencial (`ccf credential`); o `e2e.sh` copia para a máquina do AF a CA do CCF e a credencial. O grupo CP roda depois do grupo do NEF do MVP-6. O `af/scripts/negative/capif_api.py` faz o papel de um segundo *invoker*, em Python (`ssl` e o OpenSSL da linha de comando para o par de chaves), independente do AF em Go.

| ID | Verificação | Esperado |
|---|---|---|
| CP00 | o AF do modo CAPIF sobe | no log do AF: entrada, descoberta, `OAUTH` e `check-authentication` |
| CP01 | entrada do `capif_api.py` | `201` e `Location`; certificado com CN igual ao identificador, emitido pela CA do CCF (`openssl verify`) |
| CP02 | descoberta | `apiRoot` `http://192.168.16.130:8000`, `aefId` `nef-aef`, `securityMethods` `[OAUTH]` |
| CP03 | método de segurança | `201`, `selSecurityMethod` `OAUTH` |
| CP04 | `check-authentication` | `200` |
| CP05 | *token* | `200`, `token_type` `Bearer`, `Cache-Control: no-store`; `iss` = `client_id` = identificador; `scope` `3gpp#nef-aef:3gpp-akma`; `exp` igual ao instante do pedido mais a vida, com tolerância de poucos segundos |
| CP06 | K_AF via NEF com o *token* do CCF | igual à do AAnF para o mesmo A-KID e AF_ID; sem `supi` (TS 33.535 §6.3 passo 5) |
| CP07 | o mesmo, anônimo | sem `supi` nem `gpsi` |
| CP08 | entrada sem credencial; com credencial adulterada | `401`, com `WWW-Authenticate: Bearer` |
| CP09 | descoberta sem certificado de cliente | `403` |
| CP10 | descoberta com o `api-invoker-id` do AF em Go | `403` |
| CP11 | *token* com segredo errado | `400` `invalid_client` |
| CP12 | *token* com escopo fora do perfil | `400` `invalid_scope` |
| CP13 | `grant_type=password` | `400` `unsupported_grant_type` |
| CP14 | NEF sem *token* | `401` (o modo do MVP-6 está ligado; F12) |
| CP15 | NEF com o *token* adulterado | `401` |
| CP16 | AF_ID `af3.akma.example.FF00000001`: na política do NEF, fora do perfil | `403` sem `cause` |
| CP17 | AF_ID `af4.akma.example.FF00000001`: no perfil, fora da política do NEF | `403` sem `cause` |
| CP18 | saída do `capif_api.py`; *token* pedido depois | `204`; `400` `invalid_client` |
| R2c | Ua\* do UE com o AF do modo CAPIF, pelo plano de usuário | como o R2n com o AF via NEF |
| CP19 | saída do AF no fim (grupo B) | `204`; linha no log do AF |
| CP20 | logs do CCF, do NEF e do AF | nenhum *token*, credencial, segredo de entrada ou chave privada |

As 161 verificações de hoje continuam, com o grupo do NEF do MVP-6 (C13).

## 11. Limitações

- CAPIF parcial (DV-14): sem eventos, revogação, *logging*, auditoria, política de acesso, roteamento, registro do provedor pela API, descoberta aberta, interconexão, RNAA, ocultação de topologia e tarifação.
- Sem TLS entre o AF e o NEF e entre o NEF e o CCF até o MVP-11 (DV-9, DV-1).
- O estado do CCF fica em memória (seção 5.12).
- Um CCF, um provedor (o NEF), um AEF.
- Sem TLS em CAPIF-3, o CCF não sabe qual AEF pede os dados do *invoker* (seção 5.9).
- Os AF_ID de um AF externo ficam em três configurações: no perfil do CCF, na política do NEF e na lista do AAnF.
- Sem revogação de certificados em CAPIF-1e: sem CRL, sem OCSP e sem o ponto de distribuição de CRL nos certificados (DV-16).

## 12. Achados

Registrados no registro de defeitos em 29/10/2026, com os itens 4 e 5 do plano (MVP-9, "Achados nas normas"):

1. `[DIVERGÊNCIA]` TS 29.222 §7.3 contra a RFC 9113 §3.1 e §3.2: o HTTP/2 pelo *upgrade* do HTTP (seção 1). Igual na V20.1.0. D-AE.
2. `[LACUNA]` TS 33.122 A.1: a AEFPSK do método 1 sai do *Session ID* e do *Master Secret* da sessão TLS de CAPIF-1e, que são do TLS 1.2. Conferido em 29/10/2026: TS 33.122 A.1 remete o perfil TLS a TS 33.310 Anexo E, que o remete a TS 33.210, e TS 33.210 §6.2 exige o TLS 1.3 além do 1.2. A.1 não diz como derivar a chave com o TLS 1.3. D-AD.
3. `[LACUNA]` O formato de `apiInvokerPublicKey`, `apiInvokerCertificate`, `authenticationInfo` e `authorizationInfo`, e o esquema do `apiRoot` em `InterfaceDescription` (seções 4.2, 5.9 e 7.2). Sugestões D-AF e D-AG.

## 13. Pontos para aprovação

| Nº | Ponto | Proposta | Base e alternativas |
|---|---|---|---|
| C17 | Identificadores | `apiInvokerId` em UUID v4, que é também o `onboardingId` e o `securityId`; `aefId` `nef-aef` e `apfId` `nef-apf`, fixos; `apiId` em UUID v4 (seção 4.1) | `[LACUNA]` TS 29.222 só diz que são cadeias; §5.6.2.3.2 iguala o `securityId` ao identificador do *invoker*. Alternativa: identificadores com prefixo legível, como `inv-<n>` |
| C18 | Chaves e certificados | a tabela da seção 4.2: CA do CCF ECDSA P-256; certificado do servidor com o IP; chave RSA do *token* e chave RSA da credencial, separadas; o par do *invoker* gerado pelo AF; certificado do *invoker* com CN igual ao identificador, `clientAuth`, um ano; PEM; `make-capif-certs.sh`, com o OpenSSL, na pasta `cert/` do free5GC | `[3GPP]` TS 33.122 §6.1 passo 4; TS 29.222 tab. 8.4.4.2.5-1 (CN). `[LACUNA]` a codificação das chaves e certificados. Alternativa: o CCF gera as próprias chaves na primeira partida, em Go |
| C19 | Credencial de entrada | JWT `RS512` com a chave da credencial; *claims* `iss`, `sub` (o perfil), `aud` `api-invoker-management`, `iat` e `exp`; emitida por `ccf credential`; enviada como `Authorization: Bearer`; vale para várias entradas enquanto não vence; 24 h por padrão, 1 h no `e2e.sh` (seção 5.10) | `[3GPP]` TS 33.122 §6.1 (JWT e JWS; NOTA 1: como o *invoker* a recebe, fora da norma); TS 29.222 §5.5.2.2.2 NOTA 4. Alternativas: um segredo por AF (*"Other credentials may also be used"*); credencial de uso único, com `jti` |
| C20 | Perfis de AF no CCF | `invokerProfiles`, com o nome, as APIs e os AF_ID; entrada automática com a credencial válida; `apiList` da resposta com as APIs publicadas do perfil, sem considerar a `apiList` do pedido | `[3GPP]` TS 29.222 §5.5.2.2.2 passo 1 e NOTA 1; tab. 8.4.4.2.2-1; TS 23.222 §8.1.3 NOTA 2. Alternativa: recusar a entrada que peça API fora do perfil; pede `failureReasons`, da *feature* `CAPIF_Ext1` |
| C21 | `authenticationInfo` e `authorizationInfo` | `authenticationInfo`: o certificado da CA do CCF, em PEM; `authorizationInfo`: o texto JSON `{"akmaAfIds":[...]}`. Sem TLS em CAPIF-3, a resposta tem todos os itens do contexto (seção 5.9) | `[LACUNA]` tab. 8.5.4.2.3-1: cadeias sem formato. `[3GPP]` TS 33.122 §6.5.2.3 passo 5 (o CCF pode devolver a CA do *invoker*). Alternativas: os AF_ID como lista separada por espaço; os AF_ID no `scope` do *token* (opção c do C11) |
| C22 | Escopo e vida do *token* | escopo no formato sem `CAPIF_Ext1` e `CAPIF_Ext2`; pedido sem `scope` recebe todas as APIs permitidas; vida de 10 min, configurável; tolerância de 30 s no NEF; o AF renova 30 s antes do fim; `Cache-Control: no-store` | `[3GPP]` TS 29.222 tab. 8.5.4.2.6-1 e 8.5.4.2.7-1; TS 33.122 C.2.2, C.4.1 e C.5; RFC 6749 §3.3 e §5.1. Alternativa: `scope` obrigatório |
| C23 | Recusas de autenticação no CCF | credencial ausente ou inválida na entrada: `401` com `WWW-Authenticate: Bearer`; certificado de cliente ausente, inválido ou de outro *invoker* nas demais operações de CAPIF-1e: `403`; *invoker* desconhecido: `403` na descoberta e `404` nos recursos dele; no *token*, `AccessTokenErr` com `400`, ou `401` `invalid_client` com `Authorization: Basic` | `[3GPP]` TS 29.122 tab. 5.2.6-1; TS 29.222 tab. 8.5.5.3-1 e 8.5.2.3.4.4.2-3; RFC 6749 §5.2. RFC 9110 §15.5.2: o `401` *"MUST send a WWW-Authenticate header field"*, e o TLS mútuo não é um esquema de autenticação do HTTP. Alternativa: `401` sem certificado de cliente |
| C24 | Publicação pelo NEF | na partida, com nova tentativa a cada 5 s; a mesma API publicada de novo pelo mesmo APF substitui a anterior; retirada ao encerrar; `HTTP_1_1`, `JSON` e `OAUTH` (seção 6.2) | `[3GPP]` TS 29.222 §5.3.2.2 e tab. 8.2.4.2.2-1, 8.2.4.2.4-1 e 8.2.4.2.7-1; TS 29.522 §5.14.1. Alternativa: o NEF consulta as publicações (`Get_Service_API`) e só publica o que falta, sem a substituição no CCF |
| C25 | Verificação no NEF | o modo vem da chave que assina o *token*; `allowNrfTokens: true` no laboratório (C13); dados do *invoker* por 10 min; `check-authentication` só dos endereços de `akma.authorizedAfs` (DV-9); as respostas da seção 6.4 | `[3GPP]` TS 29.522 §6 e §7.2; TS 29.222 §6.2.2.2; TS 23.222 §8.16.3 passo 2a; RFC 6750 §3. Alternativas: o modo pelo prefixo `3gpp#` do `scope`; `check-authentication` de qualquer endereço |
| C26 | O AF | estado num arquivo JSON `0600`; uma nova entrada quando o CCF não conhece mais o *invoker*; `offboardOnExit` desligado por padrão e ligado no `e2e.sh`; `notificationDestination` configurado e não servido; `apiScheme: http`; *token* só na memória (seção 7) | `[3GPP]` TS 23.222 §8.1.1 (entrada única); TS 33.122 §6.8 passo 5; TS 29.501 §4.4.1. Alternativa: entrar a cada partida e sair a cada parada |
| C27 | Endereços e partida | CAPIF-1e em `192.168.16.130:8020`; o servidor interno em `127.0.0.1:8020`; `run.sh -ccf` põe o CCF antes do NEF; `force_kill.sh` para o CCF; `build-ccf.sh`; `ccfcfg.yaml` na pasta `config/` do free5GC | C1 e C3; como o `-aanf` (D15). Alternativa: o CCF no fim da lista, contando só com a nova tentativa do NEF |
| C28 | Testes | o grupo CP (CP00 a CP20) e o passo R2c da seção 10.2; o `capif_api.py` como segundo *invoker*, em Python; os AF_ID de teste `af3` (já no `nefcfg.yaml`) e `af4` (só no perfil do CCF); os casos que pedem outra chave ou outro relógio nos testes Go | critérios 1 a 7 do plano. Alternativa: só o caminho positivo no `e2e.sh` |
| C29 | HTTP/2 em CAPIF-1e | HTTP/1.1 e HTTP/2 negociados por ALPN, como manda a RFC 9113 §3.2, e não pelo *upgrade* do HTTP de TS 29.222 §7.3. Registrar como DV-15, escolha entre textos | `[DIVERGÊNCIA]` TS 29.222 §7.3 contra a RFC 9113 §3.1 (*"deprecated by this document"*) e §3.2. Alternativa: só HTTP/1.1, que §7.3 torna obrigatório e que dispensa a escolha |
| C30 | Defeitos e sugestões novos | registrar D-AD (TS 33.122 A.1 com o TLS 1.3) e D-AE (TS 29.222 §7.3 contra a RFC 9113), e as sugestões D-AF (formatos das cadeias de segurança) e D-AG (esquema do `apiRoot`) | seção 12 |

**Revisão de 30/10/2026.** Na implementação do CCF, o responsável decidiu três pontos em que esta especificação não seguia TS 33.310 e TS 33.210:
- o C18 foi revisto: CA em ECDSA P-384; chaves do *invoker* em P-256, P-384 ou RSA de 2048 a 4096 bits;
- o C19 e o C9 do plano foram revistos: `ES256` com `kid`;
- foi criado o DV-16.

As seções 4.2, 5.8 e 5.10 já trazem a revisão; os motivos estão na seção 15.1.

## 14. Arquivos previstos

- `free5gc/ccf/` (NF nova) e `free5gc/config/ccfcfg.yaml`;
- `free5gc/patches/openapi/0004` (modelos e clientes do CAPIF, C4);
- `free5gc/patches/nef/0002` e `free5gc/patches/config/0003`;
- `free5gc/patches/run/0002` (`run.sh -ccf` e `force_kill.sh`);
- `free5gc/scripts/`: `build-ccf.sh` e `make-capif-certs.sh`;
- `af/af/internal/capif` (novo), `af/af/internal/config` e `af/config/afcfg-capif.yaml` (novo);
- `af/scripts/e2e.sh`, `af/scripts/e2e-remote.sh` e `af/scripts/negative/capif_api.py` (novo);
- este documento, `nef.md`, `correcoes-free5gc.md` (F13) e o `README.md` deste lado; `docs/af-side/af.md` e `e2e.md`;
- `docs/conformidade/desvios.md` (DV-15 e DV-16, e DV-1, DV-9, DV-13 e DV-14 com a implementação); o registro de defeitos (D-AD a D-AG) e a cópia em `~/akma`;
- o plano, a matriz de rastreabilidade (revisão 10) e o tutorial (parte 11).

## 15. Implementação

Um componente de cada vez, cada um com os seus *commits*:

| Componente | Estado | Onde |
|---|---|---|
| modelos e clientes do CAPIF (C4) | feito em 29/10/2026 | `free5gc/patches/openapi/0004-capif-models-and-clients.patch`; `openapi.md` |
| CCF | feito em 30/10/2026 (seção 15.1) | `free5gc/ccf/`, `free5gc/config/ccfcfg.yaml`, `free5gc/scripts/build-ccf.sh` e `make-capif-certs.sh` |
| NEF como APF e AEF | feito em 30/10/2026 (seção 15.2) | `free5gc/patches/nef/0002-capif.patch`, `free5gc/patches/config/0003-nef-capif.patch` |
| AF como *API invoker* | feito em 30/10/2026 (seção 15.3) | `af/af/internal/capif`, `af/config/afcfg-capif.yaml` |
| partida, VM e `e2e.sh` | feito em 02/11/2026 (seção 15.4) | `free5gc/patches/run/0002-ccf-option.patch`, `af/scripts/` |
| documentos | feito em 02/11/2026 | seção 14 |

### 15.1 O CCF (30/10/2026)

**Decisões do responsável (30/10/2026).** Três pontos achados na implementação, em que a especificação aprovada não seguia TS 33.310 e TS 33.210:

1. **Chaves (C18 revisto).** TS 33.310 §6.1.1 pede suporte a chaves ECDSA de pelo menos 384 bits. Pede também que a chave que assina um certificado tenha pelo menos o nível de segurança da chave certificada. Com a CA em P-256 e só chaves P-256 aceitas, o CCF não atendia à primeira regra; se aceitasse P-384, quebraria a segunda. Decisão: a CA passa a P-384, e as chaves do *invoker* podem ser P-256, P-384 ou RSA de 2048 a 4096 bits com expoente de pelo menos 65537 (seção 4.2). O AF continua com P-256.
2. **JWS (C9 e C19 revistos).** TS 33.210 §6.3.3.1 pede `ES256` e `kid` a toda função que suporte JWS; o C9 e o C19 tinham só `RS512`. Decisão: o *token* e a credencial passam a `ES256`, com `kid` (seções 5.8 e 5.10), e o NEF vai conferir `ES256` (seção 6.4).
3. **Revogação: DV-16.** O ponto de distribuição de CRL é obrigatório nos certificados de entidade TLS (TS 33.310 §6.1.3a); TS 33.310 §6.1.1 pede suporte a CRL; o OCSP e o OCSP *stapling* são recomendações. O `x5t#S256` também fica de fora. Decisão: registrar o desvio, como o DV-10 fez com o OCSP no Ua\*.

**Ajustes à norma, sem decisão.** A especificação não tratava destes pontos:

- a CA com `pathlen:0` (TS 33.310 §6.1.4a);
- os certificados de entidade sem `basicConstraints`. §6.1.3a não pede essa extensão, e §6.1.1 manda marcar como não crítica toda extensão que a norma não exige. O Go e o OpenSSL a punham crítica;
- os nomes com `O`, em `UTF8String` (§6.1.1);
- o expoente RSA de pelo menos 65537 (§6.1.1);
- a sessão TLS encerrada depois da saída do *invoker* (TS 33.122 §6.8 passo 6; seção 5.3).

**Correções do texto desta especificação.**

- Seção 5.11: as causas de TS 29.500 tab. 5.2.7.2-1 valem para o CAPIF (TS 29.222 §7.7; TS 29.122 §5.2.6), e o CCF as usa. O texto aprovado dizia que nenhum caso tinha causa. Os códigos `405`, com `Allow`, e `411` também se aplicam.
- `supportedFeatures` é obrigatório no POST da entrada (tab. 8.4.4.2.2-1) e no da publicação (NOTA 1 da tab. 8.2.4.2.2-1). Sem ele, a resposta é `400` com `MANDATORY_IE_MISSING`; um valor que não é hexadecimal recebe `400` com `MANDATORY_IE_INCORRECT`.
- A credencial: o CCF confere também `iss`, que tem de ser o `ccfName`.
- TLS em CAPIF-1e: o TLS não sabe que operação vem na sessão. Por isso o CCF pede o certificado de cliente em toda sessão, sem exigi-lo. Na entrada, o AF ainda não tem certificado e não manda nenhum, e só o servidor se autentica (TS 33.122 §6.1 passo 2). As outras operações recusam o pedido sem o certificado do *invoker*.

**Outros comportamentos.**

- A configuração é lida em modo estrito: um atributo desconhecido impede a partida. Também impedem a partida um endereço `0.0.0.0`, um identificador com os delimitadores do `scope` (`#`, `:`, `,`, `;` ou espaço), um nome repetido e um AF_ID inválido.
- Na partida, o CCF confere que a chave da CA é a do certificado da CA e que o certificado tem `O`.
- O CCF não grava o registro de chaves TLS que as NFs do free5GC gravam. Com ele, as sessões de CAPIF-1e poderiam ser lidas, com os segredos e os *tokens* (seção 9).

**Verificação (30/10/2026).**

- `go test -race ./...` no módulo `free5gc/ccf`, sem falhas:
  - `jws`: a impressão digital da chave do RFC 7517 A.1 confere com a calculada à parte, em Python; e as recusas de outra chave, de outro `kid`, de `crit`, de `RS512`, de `HS256` e de `none`;
  - `credential`, `context` (as chaves e o nome em DER) e `factory`;
  - `sbi`: as operações das seções 5.2 a 5.10, com as recusas; o perfil dos certificados emitidos; as chaves aceitas e recusadas do *invoker*; o fim da sessão TLS em HTTP/1.1 e em HTTP/2; as rotas; o TLS.
- O teste do fim da sessão falha sem o `Connection: close`: conferido com a linha retirada.
- Teste local nesta máquina, sem nenhuma VM: o binário de `build-ccf.sh`, com os certificados de `make-capif-certs.sh` para 127.0.0.1, e um cliente em Python. O teste cobriu a publicação, a entrada com chave P-384, o `openssl verify` do certificado emitido, a descoberta com e sem certificado, o método de segurança, o *token* `ES256` com `kid`, o segredo errado, os dados do *invoker*, a recusa de uma chave P-521 e a saída com `Connection: close`. Nenhum *token*, credencial ou chave privada apareceu nos logs.

**Achados na implementação.** Vão ao registro de defeitos com os documentos:

1. `TS29222_CAPIF_Security_API.yaml`, `ServiceSecurity.securityInfo`: um *array* com `minimum: 1`, palavra de números, em vez de `minItems: 1`; a tabela 8.5.4.2.2-1 dá a cardinalidade 1..N. Igual no YAML da V20.1.0. Editorial.
2. `[LACUNA]` TS 33.122 Anexo C pede o JWS da RFC 7515 e não remete ao perfil de TS 33.210 §6.3.3, que se declara aplicável a toda função 3GPP com JWS (§6.3.1). Sugestão: citar o perfil no Anexo C, como §6.1 e §6.3.1.1 citam o perfil TLS e os de certificado.

### 15.2 O NEF (30/10/2026)

São dois *patches*: `free5gc/patches/nef/0002-capif.patch`, aplicado depois do `0001`, e `free5gc/patches/config/0003-nef-capif.patch`, depois do `0002`. Foram feitos na cópia do NEF v1.2.5 em `free5gc/build/src`, sem nenhuma VM.

| Parte | Onde |
|---|---|
| a configuração `capif` (seção 6.1) | `pkg/factory/config.go` |
| a chave do *token* e o `kid`, os dados dos *invokers* e o `apiId` publicado | `internal/context/capif.go` |
| o cliente do CCF em CAPIF-3/4: publicação, retirada e `Obtain_API_Invoker_Info` | `internal/sbi/consumer/ccf_service.go` |
| a publicação com nova tentativa, o `check-authentication`, a conferência do *token* e o perfil do *invoker* | `internal/sbi/processor/capif.go` e `akma.go` |
| a rota de `aef-security` e o *middleware* do grupo `3gpp-akma` | `internal/sbi/api_capif.go` e `server.go` |
| a publicação na partida e a retirada no encerramento | `pkg/service/init.go` |
| a categoria de log `CAPIF` | `internal/logger/logger.go` |

**Detalhes que a especificação não fixava.**

- **Que *token* é (seção 6.4, passo 1).** Um *token* é do CCF quando é `ES256`, a assinatura confere com a chave do CCF e o `kid`, se houver, é o dessa chave. Um *token* com `crit` não é tratado como do CCF (RFC 7515 §4.1.11).
  - Com `allowNrfTokens: true`, os outros pedidos vão ao *middleware* do NRF.
  - Com `allowNrfTokens: false`, recebem `401`. Sem *token*, com `WWW-Authenticate: Bearer`; com um *token* que não é do CCF, com `Bearer error="invalid_token"`.
- **Dados do *invoker* (seção 5.9).** O NEF usa só o item do seu `aefId` com `selSecurityMethod` `OAUTH` (TS 33.122 §6.5.2.3). Os AF_ID passam à forma de transporte. Um AF_ID inválido é ignorado, com aviso no log.
- **Ordem da conferência do AF_ID (seção 6.4).** Primeiro o perfil do *invoker*, depois a política de endereço. O `detail` do `403` diz qual das duas recusou: `AF_ID not authorized for the API invoker` ou `AF not authorized by local policy`.
- **`check-authentication` (seção 6.3).** As causas seguem TS 29.500 tab. 5.2.7.2-1, por TS 29.222 §7.7, como no CCF:
  - sem `apiInvokerId` ou sem `supportedFeatures`: `400` com `MANDATORY_IE_MISSING`;
  - `supportedFeatures` que não é hexadecimal: `400` com `MANDATORY_IE_INCORRECT`.

  O grupo `aef-security` não passa pelo *middleware* do NRF, porque TS 29.222 §6.2.2.2 não pede autorização nessa operação.
- **Publicação (seção 6.2).** O aviso no log sai só na primeira falha; as seguintes vão para o nível *debug*. Ao encerrar, o NEF espera acabar a tentativa em curso e retira a publicação. O CCF tem até 2 s para responder.
- **`go.mod`.** `github.com/golang-jwt/jwt/v5`, que já era dependência indireta, passa a direta.

**Verificação (30/10/2026).**

- `go test -race ./...` na cópia do NEF passou sem falhas. Os testes novos estão em quatro pacotes:
  - `context`: a impressão digital da chave do RFC 7517 A.1, a chave do *token* e a vida dos dados do *invoker*;
  - `factory`: a seção `capif`;
  - `sbi`: as rotas e o *middleware*;
  - `processor`: os casos da seção 10.1, além do `check-authentication` e da publicação com nova tentativa e retirada.

  Os casos da seção 10.1 cobertos em `processor`:
  - *token* válido;
  - *token* de outra chave;
  - `none`, `HS256` e `RS512`;
  - *token* vencido, e *token* dentro da tolerância;
  - sem `exp`;
  - `iss` diferente de `client_id`;
  - escopo de outro AEF, e de outra API;
  - *invoker* que saiu (`404` do CCF);
  - CCF fora do ar (`503`);
  - AF_ID fora do perfil, e fora da política de endereço;
  - `allowNrfTokens` `true` e `false`.

  Os testes do MVP-6 e os do free5GC continuam passando.
- Interoperação com o CCF real: o binário do CCF rodou nesta máquina e emitiu um *token*. O *token* e a resposta do CCF a `Obtain_API_Invoker_Info` passaram na conferência do NEF, com o mesmo `kid` e os 11 AF_ID do perfil `af1-external`. O teste foi temporário e não está no *patch*.
- Os *patches* reproduzem a cópia: `src-orig` com `nef/0001` e `0002` dá a cópia do NEF, e `src-orig/config` com `config/0001` a `0003` dá a configuração. O NEF lê e valida o `nefcfg.yaml` resultante.

**Na VM** (próximo componente). Com o `nef/0001` já aplicado, `patch-nf.sh nef apply` recusa o conjunto novo. Antes de aplicar, é preciso reverter com uma cópia do diretório antigo de *patches*, como no MVP-12. O mesmo vale para `config`. Feito em 02/11/2026 (seção 15.4).

### 15.3 O AF (30/10/2026)

Código em `af/af/internal/capif` (novo), `af/af/internal/config`, `af/af/internal/nef`, `af/af/internal/logger` e `af/af/cmd/af`; configuração em `af/config/afcfg-capif.yaml`. Nenhuma VM foi tocada. Os detalhes que a especificação não fixava:

- **Modos.** `aanfUri`, `nefUri`, `nrfUri` e `capif` não se combinam. `ccfUri` só aceita HTTPS. Sem `apiScheme`, vale `https`.
- **Estado.** É um JSON `0600`, escrito por arquivo temporário e renomeação, com `ccfUri`, o identificador, o certificado, a chave privada PKCS #8 e o segredo. Um estado de outro `ccfUri` é posto de lado, e o AF faz nova entrada.
- **Entrada.** O AF confere o certificado emitido: a CA do CCF, `clientAuth`, a própria chave e o CN igual ao identificador.
- **`apiRoot`.** Sai da primeira interface com `OAUTH`, em que os métodos da interface valem antes dos do perfil (tab. 8.2.4.2.3-1). Usa `ipv4Addr`, `ipv6Addr` ou `fqdn`, mais `port` e `apiPrefix`, ou o `domainName` quando o perfil não tem interface.
- **TLS de CAPIF-1e.** TLS 1.3, ou 1.2 só com ECDHE e AEAD (TS 33.210 §6.2), e HTTP/2 por ALPN.
- **Token recusado pelo NEF.** Com um `401`, o AF descarta o *token*, pede outro ao CCF e repete o pedido uma vez. `[3GPP]` TS 33.122 C.4.1: *"the API Invoker shall return to the CAPIF Core Function"*. Isso só vale no modo CAPIF; o modo do NRF não muda.
- **Saída.** Um `404` na saída conta como saída feita, e o estado é apagado.

**Verificação (30/10/2026).**

- `go test -race ./...` no módulo `af`, sem falhas. Os testes novos usam um CCF falso, com CA própria, e um NEF falso, e cobrem:
  - a entrada e o estado;
  - o `apiRoot`;
  - a falta de `OAUTH` e a credencial errada;
  - o *token*: guardado, renovado 30 s antes do fim e renovado depois da recusa;
  - o CCF reiniciado, na partida e no pedido de *token*;
  - a saída;
  - a seção `capif` e o `afcfg-capif.yaml`;
  - nenhum segredo no log.
- Teste local, sem VM: o binário do CCF, um NEF falso em Go que confere o *token* `ES256` com a chave real do CCF e recusa o primeiro pedido, o binário do AF com `digest-http` e o `uaclient` com o vetor `testbed-ue`. O teste passou por entrada, descoberta, método, `check-authentication` e dois *tokens* (o primeiro recusado). As duas sessões Ua\* foram estabelecidas. Na parada, o AF saiu do CCF e apagou o estado. Nenhum segredo apareceu no log do AF.

### 15.4 Partida, VM e `e2e.sh` (02/11/2026)

- **`run/0002-ccf-option.patch`.** Com `-ccf`, o `run.sh` põe o CCF logo antes do NEF (C27). O `force_kill.sh` também para o CCF.
- **`e2e.sh` e `e2e-remote.sh`.**
  - O R1 sobe o free5GC com `-aanf -ccf`, confere o CCF e a publicação do NEF e emite a credencial de entrada.
  - No ponto A, a máquina do AF copia a CA do CCF e a credencial, sobe o AF do modo CAPIF na 8086 (CP00) e roda o `af/scripts/negative/capif_api.py` (CP01 a CP18).
  - O R2c faz o Ua\* do UE com esse AF.
  - O grupo B confere a saída do AF (CP19). O CP20 confere os logs dos dois lados.
- **Resultado.** `185 passed, 0 failed, 3 skipped` na execução `20261102-093804` (`docs/af-side/e2e.md` §5.9). As três puladas (D03, O04, O05) dependem de ferramentas de conferência que não estavam compiladas e não tocam no CAPIF.
- **Dois *scripts* corrigidos na implantação na máquina do free5GC.**
  - `patch-nf.sh` agora confere uma série de *patches* em ordem, numa cópia dos arquivos. Os *patches* novos dependem dos anteriores, e a conferência de cada um sozinho os recusava.
  - `make-capif-certs.sh` faz a CA com uma configuração própria. O OpenSSL 1.1.1 da máquina acrescentava o `basicConstraints` do `-addext` ao da seção `v3_ca`, e o Go recusa uma extensão repetida (RFC 5280 §4.2).
- **Tutorial.** A parte 11 e o passo das chaves do CCF na seção 9.1. O NEF com o `nef/0002` lê a chave pública dos *tokens* na partida, por causa da seção `capif` que o `config/0003` põe no `nefcfg.yaml`.
