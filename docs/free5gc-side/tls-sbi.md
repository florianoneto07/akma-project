# TLS nas SBI e entre o NEF e o AF — especificação

**MVP:** MVP-11 · **Data:** 2 de novembro de 2026 · **Estado:** **aprovada** pelo responsável em 04/11/2026, com os pontos T18 a T39 da seção 13 como propostos; fases (a), (b) e (c) implementadas e verificadas (seções 15 a 17); o DV-18, o F26 e os defeitos D-AJ e D-AK registrados na mesma data. O escopo foi aprovado no plano em 02/11/2026, com os pontos T1 a T17 (plano, seção 4): o T4 revisto, com os dois modos e o padrão por nome; o T14 na opção (b). **Fases (a) e (b) implementadas e verificadas** em 04 e 05/11/2026 (seções 15 e 16); a fase (c) está por fazer.

Este documento detalha o que o plano aprovou: a PKI de teste do operador, os nomes das NFs, o TLS mútuo em toda a SBI, a identidade dos certificados nos *tokens* e nas políticas, e o TLS mútuo entre o NEF e o AF de fora e entre o NEF e o CCF. Os pontos abertos continuam a numeração do plano, do T18 em diante. Neste texto, "CA do operador" é a CA de teste da SBI (T3), e "CA do CCF" é a do MVP-9, que assina os certificados dos *invokers*.

## 1. Situação normativa

O plano (MVP-11, "Situação normativa") resume as normas. Esta seção registra as versões conferidas e o que a leitura acrescentou.

### 1.1 Versões conferidas (02/11/2026)

Cada cláusula citada foi lida na baseline (`specs/`) e na versão mais recente do arquivo da 3GPP, baixada em 02/11/2026, e as duas foram comparadas.

| Norma | Baseline | Mais recente | Cláusulas comparadas | Resultado |
|---|---|---|---|---|
| TS 33.501 | V19.7.0 | V20.3.0 (09/2026) | §12.1 a §12.5; §13.1.0; §13.3.0 a §13.3.2; §13.3.8.3; §13.4.1.0 a §13.4.1.1.2 | iguais. A NOTA 1 de §13.1.0 passa a dizer *"IPsec or any other lower layer security protection"* no lugar de *"NDS/IP"*; §13.3.8.3 (CCA, fora do escopo) ganha uma frase |
| TS 33.310 | V19.5.0 | V19.5.0 | §6.1.1, §6.1.3a, §6.1.3c, §6.1.4a, §6.1a, §6.3.1 e §7.1 | a mesma versão. Fonte: o portal da 3GPP (`specificationId=2293`); a lista do arquivo não abriu |
| TS 33.210 | V19.3.0 | V20.0.0 (09/2026) | §6.2 e §6.3 | §6.2 igual. §6.3.3.1 ganha duas frases: *"Usage of "alg" header parameter values RS256, RS384 and RS512 ... is not recommended."* e *"For the validation of header parameter values, the validation guidelines described in section 3.10 of RFC 8725 should be followed."* (CR 0103 e 0104) |
| TS 29.510 | V19.7.0 | V20.1.0 (09/2026) | §5.4.2.2.1; tab. 6.1.6.2.2-1, 6.1.6.2.3-1, 6.1.6.2.5-1, 6.1.6.3.3-1, 6.3.4.2.2-2, 6.3.5.2.2-1 e 6.3.5.2.4-1 | iguais; as tabelas só ganham atributos que o MVP não usa (`energyPriority`, `senfInfo`) |
| TS 29.500 | V19.7.0 | V20.0.0 (09/2026) | §6.7.1 a §6.7.3 | iguais, salvo um parágrafo de §6.7.3 sobre o SCP |
| TS 29.501 | V19.5.0 | V20.0.0 | §4.4.1 e §4.4.3 | iguais |
| TS 23.003 | V19.7.0 | V20.1.0 (09/2026) | §28.2, §28.3.2.3, §28.3.2.4 e §28.3.2.5 | iguais |
| TS 33.122 | V19.4.0 | V19.4.0 | §6.3.1.1, §6.5.2.2, §6.5.2.3, §6.6 e §6.8 | a mesma versão |
| TS 29.222 | V19.8.0 | V20.1.0 | §10 e tab. 8.2.4.3.6-1 e 8.5.4.2.3-1 | iguais |
| TS 29.535 | V19.5.0 | V20.1.0 (09/2026) | §5.1.9 e tab. 4.2.2.1-1 | iguais |

As RFC que TS 33.310 cita para as extensões novas, lidas no rfc-editor.org:

- RFC 9310: `id-pe-nftype` = 1.3.6.1.5.5.7.1.34; `NFTypes ::= SEQUENCE SIZE (1..MAX) OF NFType`; `NFType ::= IA5String (SIZE (1..32))`; *"This extension MUST NOT be marked critical."* Os valores são os de TS 29.510 tab. 6.1.6.3.3-1; §5: *"an operator MAY assign its own NF Types"*, sem colidir com os da norma.
- RFC 9509: `id-kp-oauthAccessTokenSigning` = 1.3.6.1.5.5.7.3.39 e `id-kp-jwt` = 1.3.6.1.5.5.7.3.37. §3: quem confere a assinatura de um *token* *"SHOULD require that corresponding KeyPurposeIds be specified by the EKU extension"*.

O texto segue a baseline; onde a versão nova acrescenta algo que pesa, o texto diz.

### 1.2 O que a leitura acrescentou

Pontos que o plano não trazia. Os que pediam decisão estão na seção 13.

1. `[3GPP]` TS 29.510 §5.4.2.2.1: *"The NRF may verify that the input attributes (e.g. NF type) in the access token request match with the corresponding ones in the public key certificate of the NF service consumer. If the verification is successful, other authorization check shall be performed, otherwise, the request shall be rejected immediately with "400 Bad Request" status code, and "error" attribute set to "invalid_client"."* O T9 aprovado dizia 403 (T18).
2. `[3GPP]` TS 33.501 §13.4.1.1.2 passo 2: além de `exp`, `aud` e `sub`, o produtor *"shall check that the identity in the issuer claim in the access token matches the identity of the OAuth 2.0 authorization server(s)"*. NOTA 2a: o produtor precisa da raiz para validar *"the NRF's OAuth 2.0 Access Token certificate (chain)"*. TS 33.310 §6.1.3c.3, NOTA 1b: sem o URI com o NF Instance ID no certificado de assinatura, *"the identity of the NRF in JWT access token "iss" Issuer claim ... cannot be securely validated"* (T29).
3. `[DIVERGÊNCIA]` TS 33.310 §6.1a: *"CRL retrieval with LDAPv3 (RFC 4510, RFC 4517, RFC 4523, RFC 4512) shall be supported as the primary method. HTTP may be used for checking the revocation status of TLS and NE certificates."* O T7 prevê só HTTP (T27).
4. `[3GPP]` TS 33.310 §6.3.1: *"Any NE, SEG shall not validate received certificates from a peer entity whose CRL distribution point field is empty, but end the path validation with a negative result."* A regra cita NE e SEG; para a entidade TLS, a frase é *"Validity of certificates received from the TLS entity shall be verified by CRLs or OCSP responses ... based on the CRL Distribution Point or Authority Information Access extensions in the certificates."* `[INTERPRETAÇÃO]` Sem o ponto de distribuição, a conferência que §6.3.1 manda fazer é impossível: o T7 (recusar) segue valendo (T18). A mesma cláusula: *"Certificate validity calculation results shall not be cached in TLS entities for longer than the TLS connection lifetime."*
5. `[3GPP]` TS 23.003 dá o formato de mais dois nomes, além do NRF: o NSSF, `nssf.5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org` (§28.3.2.4.2), e o AMF, `<AMF-id>.amf.5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org` (§28.3.2.5). TS 29.510 tab. 6.1.6.2.2-1: *"For AMF, the FQDN registered with the NRF shall be that of the AMF Name"*; TS 33.310 §6.1.3c.3: o DNS-ID do AMF é o AMF Name (T20).
6. `[3GPP]` TS 29.510 tab. 6.1.6.2.5-1, NOTA 4: com `https`, o consumidor *"shall use such port number along with the FQDN parameter present in the NFService or NFProfile data types to construct the target URI"*; NOTA 2: sem porta, 443. Tab. 6.1.6.2.3-1, NOTA 1: com TLS, o FQDN do serviço, se houver, senão o do perfil; NOTA 14: o consumidor pode resolver o FQDN ou usar os endereços de `ipEndPoints`. `apiPrefix` é *"optional deployment-specific string that starts with a "/" character"* (tab. 6.2.6.2.4-1; TS 29.501 §4.4.1) (T22).
7. `[3GPP]` TS 29.501 §4.4.3: no URI de *callback*, *"'host' is either an FQDN or an IP address"* (T23).
8. `[3GPP]` TS 29.500 §6.7.3: as respostas do produtor que recusa o *token*: 401 com `WWW-Authenticate: Bearer`, `realm` igual ao URI da API e `error="invalid_token"` quando havia *token*; `ACCESS_TOKEN_CLAIM_MISSING` (401) e `ACCESS_TOKEN_CLAIM_MISMATCH` (403) em `ProblemDetails`, com `invalidParams`; 403 com `error="insufficient_scope"` e `scope` quando falta escopo (T29).
9. `[3GPP]` TS 29.510 tab. 6.3.5.2.4-1: `iss` é obrigatório e *"shall contain NF instance id of the NRF"*. O NRF do free5GC o manda vazio (achado 6, F26, T30).
10. `[LACUNA]` TS 33.122 §6.5.2.3 passo 5: no método 3, o TLS de CAPIF-2e é só de servidor ou mútuo *"as indicated by CAPIF core function"*. TS 29.222 (V19.8.0 e V20.1.0) não tem atributo para essa indicação: `SecurityMethod` só tem `PSK`, `PKI` e `OAUTH` (tab. 8.2.4.3.6-1), e `SecurityInformation` só `authenticationInfo` e `authorizationInfo`, cadeias sem formato (tab. 8.5.4.2.3-1). O mesmo passo 5 põe o `Authentication Initiation Request` *"prior to establishment of TLS session"* sem dizer com que proteção. Candidato a defeito D-AJ (T34).
11. `[3GPP]` TS 33.122 §6.6: em CAPIF-3/4, certificados nos perfis de TS 33.310 §6.1.3a e §6.1.4a. O CCF não é um tipo de NF de TS 29.510 tab. 6.1.6.3.3-1: o certificado do seu servidor interno segue §6.1.3a, sem `nfTypes` (T24).
12. `[3GPP]` TS 33.210 §6.2.2: com HTTP/2 sobre TLS 1.3, *"HTTP/2 servers shall not send post-handshake TLS 1.3 CertificateRequest messages"*. O certificado de cliente é pedido no *handshake*.
13. `[3GPP]` TS 33.501 §13.1.0, NOTA 3: *"It is a vendor implementation decision how the SNI extension is being used in TLS servers."* RFC 9113 §9.2: *"If the server is identified by a domain name [DNS-TERMS], clients MUST send the server_name TLS extension"*; RFC 6066 §3: *"Literal IPv4 and IPv6 addresses are not permitted in "HostName"."*
14. `[DIVERGÊNCIA]` editorial: TS 33.310 §6.1.3c.3, antes da tab. 6.1.3c.3-2, diz *"the rest of extensions remains the same as in table 6.1.3c.3-2"*, a própria tabela; o sentido pede a 6.1.3c.3-1. Candidato a defeito D-AK.

## 2. Situação no free5GC v4.2.3

Os achados 1 a 5 do plano, conferidos nas fontes (`free5gc/build/src-orig/NFs`, `openapi` v1.2.4 e v1.2.5-0.20260527003827-02dc71b4d94f, `util` v1.3.2), e um sexto. Estão em `correcoes-free5gc.md` como F21 a F25 (aprovados com o T17) e F26 (T30).

| Achado | Onde | O que o free5GC faz | Registro |
|---|---|---|---|
| 1 | `openapi`, `client.go` (`innerHTTP2Client`) e `udr/ImsDataRepository/client.go` | o cliente HTTP/2 usa `InsecureSkipVerify: true` e não apresenta certificado. O NRF, com `oauth` e `https`, exige e confere certificado de cliente contra `cert/root.pem` (`internal/sbi/server.go`, `configureOAuthMutualTLS`): nenhuma NF consegue falar com ele | F21 |
| 2 | `util`, `httpwrapper.NewHttp2Server` | com um caminho de log, grava os segredos de cada sessão TLS (`KeyLogWriter`). Os servidores de métricas usam a mesma função. O arquivo é `log/<data>/key/<nf>sslkey.log` (o plano dizia `<nf>_sslkey.log`); o do AAnF, `aanf_sslkey.log`. Na .130, lido em 02/11/2026: os doze arquivos existem, vazios, porque a SBI está em `http` | F22 |
| 3 | NRF, `processor/access_token.go` (`AccessTokenProcedure`) e `context.GetTokenCtx`; `openapi`, `oauth.VerifyOAuth` | os *tokens* saem em `RS512`, e `VerifyOAuth` só aceita `RS512` | F23 |
| 4 | NRF e todas as NFs | os perfis não têm `fqdn`; `apiPrefix` leva o URI inteiro, `<esquema>://<IP>:<porta>`. Os consumidores montam o URI de quatro jeitos: `SearchNFServiceUri` (AMF, SMF, UDM, PCF) devolve o `fqdn` sem esquema nem porta; o do NEF, sem porta; `openapi.GetNFServiceUri` (AUSF, e o AAnF e o *patch* `ausf/0001` do projeto) prefere `apiPrefix`; o SMF e o PCF usam `service.ApiPrefix` direto. Pôr o `fqdn` no perfil, sem mais nada, quebra o AMF, o SMF, o UDM e o PCF | F24 |
| 5 | `openapi`, `oauth.GenerateCertificate`; NRF, `context.SignNFCert` e `AccessTokenScopeCheck` | o NRF gera para cada NF um certificado RSA 2048, de 10 anos, `C=TW, ST=Taiwan, L=Hsinchu, O=free5gc, OU=free5gc`, CN e DNS com o tipo da NF, sem `keyUsage`, `extendedKeyUsage`, ponto de distribuição de CRL nem `nfTypes`. Regera o próprio certificado a cada partida (na .130: `CN=NRF`, `DNS:NRF`, sem URI). No pedido de *token*, confere o arquivo que ele mesmo gerou, não o certificado da sessão TLS, e só avisa quando a cadeia não confere (*"unknown authority"*) | F25 |
| 6 | NRF, `context.NRFContext.Nrf_NfInstanceID` | o campo nunca é preenchido: todo *token* sai com `iss` vazio, e o certificado do NRF, sem o URI | F26 (T30) |

Outros fatos do free5GC que o projeto usa:

- `[FREE5GC]` O próprio URI de cada NF (o `apiPrefix` do perfil, os *callbacks* e o `Location`) sai de uma função só na maioria das NFs: `GetIPv4Uri()` (AMF, PCF, UDM, UDR, NSSF), `context.Url` (AUSF, CHF), `factory.SbiUri()` (NEF) e `GetSbiUri()` (NRF). O SMF e o BSF montam o URI no lugar, com `RegisterIPv4`.
- `[FREE5GC]` As NFs já leem `sbi.tls.pem` e `sbi.tls.key` (o NRF, `sbi.cert` e `sbi.rootcert`) e servem `https` com `ListenAndServeTLS`. Na .130, todas em `http`, cada uma num endereço de *loopback* na porta 8000 (seção 3.1).
- `[FREE5GC]` O banco de testes monta as configurações em Go, com `http`, e roda todas as NFs num processo só, com o cliente global do `openapi` (`test/nfs_setup.go`).
- `[FREE5GC]` O `openapi` não depende do `util`, e o `util` não depende do `openapi`. O AF, o AAnF e o CCF já compilam com a cópia do `openapi` por `replace`.
- `[FREE5GC]` O OpenSSL da .130 é o 1.1.1.

## 3. Nomes, endereços e os dois modos (T4, DV-17)

### 3.1 Nomes e endereços

Domínio de origem (TS 23.003 §28.2): `5gc.mnc093.mcc208.3gppnetwork.org`, com o MNC em três dígitos.

| Entidade | FQDN (modo por nome) | Endereço e porta | `nfTypes` | Base |
|---|---|---|---|---|
| NRF | `nrf.5gc.mnc093.mcc208.3gppnetwork.org` | 127.0.0.10:8000 | `NRF` | `[3GPP]` TS 23.003 §28.3.2.3.2 |
| AMF | `amf1.amf.5gc.mnc093.mcc208.3gppnetwork.org` | 127.0.0.18:8000 | `AMF` | `[3GPP]` §28.3.2.5; o AMF Name (T20) |
| NSSF | `nssf.5gc.mnc093.mcc208.3gppnetwork.org` | 127.0.0.31:8000 | `NSSF` | `[3GPP]` §28.3.2.4.2 |
| SMF, UDR, PCF, UDM, AUSF, CHF | `<tipo>.5gc.mnc093.mcc208.3gppnetwork.org`, em minúsculas | 127.0.0.2, .4, .7, .3, .9, .113, porta 8000 | o tipo | `[INTERPRETAÇÃO]` o modelo do NRF |
| BSF | `bsf.5gc.mnc093.mcc208.3gppnetwork.org` | 127.0.0.15:8000 | `BSF` | idem; só com `run.sh -bsf`, que o laboratório não usa |
| NEF | `nef.5gc.mnc093.mcc208.3gppnetwork.org` | 192.168.16.130:8000 | `NEF` | idem |
| AAnF | `aanf.5gc.mnc093.mcc208.3gppnetwork.org` | 192.168.16.130:8010 | `AANF` | idem |
| CCF, CAPIF-3/4 | `ccf.5gc.mnc093.mcc208.3gppnetwork.org` | 127.0.0.1:8020 | — | idem; TS 33.122 §6.6 (seção 4.4) |
| CCF, CAPIF-1e | sem nome: por IP | 192.168.16.130:8020 | — | T34 |
| AF interno | `af1.akma.example` | cliente, em 192.168.16.129 | `AF` | o FQDN do AF_ID; T13 |
| CRL | sem nome: por IP | `http://192.168.16.130:8030` | — | T26 |

Os nomes resolvem pelo `/etc/hosts`, como decidiu o T4:

- na .130: os treze nomes acima, cada um no seu endereço (o do CCF em 127.0.0.1);
- na .129: `aanf.5gc...` e `nef.5gc...`, em 192.168.16.130. São os dois servidores da SBI que a máquina do AF alcança: o AAnF, para o AF interno, e o NEF, para o AF de fora.

### 3.2 A escolha do modo (T21)

- **NFs, AAnF e CCF.** Uma chave nova, `sbi.fqdn` (no CCF, `internal.fqdn`). Com ela, modo por nome: o perfil leva o `fqdn`, e os URIs próprios (o `apiPrefix`, os *callbacks* e o `Location`) saem com o nome. Sem ela, modo por IP, como hoje (DV-17). Com `scheme: https` e sem `fqdn`, a NF avisa no log, uma vez: `SBI over https without FQDN: URIs by IP address (DV-17; TS 29.510 cl. 6.1.6.2.2 NOTE 1)`.
- **Padrão.** No código, a ausência da chave é o modo por IP: as configurações do free5GC e o banco de testes continuam valendo. Nas configurações do projeto (`config/0004`, `aanfcfg.yaml`, `ccfcfg.yaml`, os `afcfg*.yaml`), o modo por nome, que é o padrão do laboratório pedido no T4.
- **Clientes.** O modo de quem só é cliente (o AF, o `seed`, os testes em Python) é a forma do URI configurado. Com um nome, o cliente manda o SNI e confere o DNS-ID; com um IP, confere o IP do certificado e não manda o SNI (RFC 6066 §3). É o que o `crypto/tls` do Go faz sozinho, a partir do *host* do URI. Os consumidores que acham o produtor no NRF seguem o que o produtor registrou (seção 3.3).

Os certificados são os mesmos nos dois modos: o `subjectAltName` leva o nome e o IP (seção 4.2).

**MVP-15, fase (a) (12/11/2026).** Decisão do responsável no MVP-14 (Q6, opção b): o modo por IP sai com `https`. As onze NFs, o AAnF e o servidor interno do CCF não partem sem `sbi.fqdn` (`internal.fqdn` no CCF), com a mensagem `sbi.fqdn is required with https (TS 29.510 table 6.1.6.2.2-1 NOTE 1)`. Com `http`, nada muda: a NOTA 1 só vale com `https`, e o banco de testes do free5GC continua por IP. O aviso no log descrito acima não chegou a existir no código, e a chave ausente agora é um erro. O AF, cliente, muda na fase (c). Os *patches* são `*/000N-fqdn-required.patch` e `config/0007`, que acerta os comentários das configurações (DV-17 em `desvios.md`).

### 3.3 Como os URIs se montam (T22, T23)

- **Produtor.** No modo por nome, o perfil leva `fqdn`, e cada serviço leva `ipEndPoints` (o IP e a porta) e `apiPrefix`. O `apiPrefix` continua com o URI inteiro, como o free5GC faz, agora com o nome: `https://<FQDN>:<porta>`. `[FREE5GC]` É o que o SMF e o PCF usam direto; corrigir o `apiPrefix` (TS 29.501 §4.4.1) fica de fora (F24, T22).
- **Consumidor.** Um *helper* só, `openapi.GetNFServiceUri`, corrigido no *patch* `openapi/0005` (TS 29.510 tab. 6.1.6.2.3-1 NOTA 1 e tab. 6.1.6.2.5-1 NOTAS 2 a 4):
  1. com o `fqdn` do serviço ou, sem ele, o do perfil: `<esquema>://<FQDN>:<porta>`, com a porta da primeira entrada de `ipEndPoints` que tiver porta, ou 443 e 80 sem porta; mais o `apiPrefix`, só quando ele começa com `/`;
  2. sem `fqdn`: o `apiPrefix`, quando é um URI absoluto (o jeito do free5GC);
  3. senão: o IP e a porta de `ipEndPoints`, como hoje.

  Os `SearchNFServiceUri` do AMF, do SMF, do UDM e do PCF e o `searchNFServiceUri` do NEF passam a chamar esse *helper*.
- **URIs próprios.** No modo por nome, a função única de cada NF (seção 2) devolve `https://<FQDN>:<porta>`. Com isso, o perfil, os *callbacks* (`deregCallbackUri`, `reauthNotifyCallbackUri`, `smContextStatusUri`, as notificações de política) e o `Location` saem com o nome, sem mexer em cada lugar. No SMF e no BSF, o *patch* cria essa função e troca as montagens no lugar.
- **AMF.** O AMF Name (`amfName`) passa a ser o FQDN `amf1.amf.5gc.mnc093.mcc208.3gppnetwork.org`, e o AMF recusa partir com `sbi.fqdn` diferente de `amfName` (TS 29.510 tab. 6.1.6.2.2-1). O mesmo nome vai no `NG Setup Response` (TS 38.413, *AMF Name*), que o UERANSIM só mostra.

### 3.4 Topologia

```
                     VM do AF (.129)                                VM do free5GC (.130)
  AF interno ──── SBI, TLS mútuo (CA do operador) ───────────────► AAnF aanf.5gc...:8010
  AF de fora ──── CAPIF-1e, TLS mútuo (CA do CCF), por IP ───────► CCF 192.168.16.130:8020
  AF de fora ──── CAPIF-2e, TLS mútuo (cliente: CA do CCF;  ─────► NEF nef.5gc...:8000
                  servidor: CA do operador)
                                                  NEF ── CAPIF-3/4, TLS mútuo (CA do operador) ──► CCF ccf.5gc...:8020
                                                  NFs, NEF, AAnF ── SBI, TLS mútuo (CA do operador) ──► entre si
  todas as pontas ── HTTP, a CRL de cada CA ────────────────────► CRL http://192.168.16.130:8030
```

O túnel SSH do MVP-6 sai: o AF de fora não fala mais com o NRF (T14).

## 4. A PKI do laboratório (T3)

Toda a PKI sai de `free5gc/scripts/make-sba-certs.sh` (seção 4.6), na pasta `cert/sba/` do free5GC, separada dos certificados que vêm com o free5GC e dos do CCF. A chave da CA não sai da .130. Só para teste: as chaves privadas não são cifradas.

### 4.1 A CA do operador

`sba-ca.pem` e `sba-ca.key`: ECDSA P-384, assinatura com SHA-384, validade de 5 anos. Perfil de TS 33.310 §6.1.4a (`[INTERPRETAÇÃO]` TS 33.310 não tem um perfil próprio de CA da SBA; §6.1.3c.3 só pede que a CA seja *"the CA in the operator domain"*):

- *Subject* `C=FR, O=5gc.mnc093.mcc208.3gppnetwork.org, CN=SBA test CA`, em `UTF8String` (§6.1.1);
- `basicConstraints` crítico, `CA:TRUE`, `pathlen:0`;
- `keyUsage` crítico, `keyCertSign` e `cRLSign`;
- `subjectKeyIdentifier` pelo método (1) da RFC 5280 §4.2.1.2.

Pela regra de §6.1.1 (*"The security level of the public key used to sign the certificate shall be at least the same as the public keys in the certificate"*), a CA é P-384, como a do CCF.

### 4.2 Certificados das NFs (TS 33.310 tab. 6.1.3c.3-1)

Um por NF, pelo AAnF e pelo AF interno, cada um com a sua chave ECDSA P-256 (§6.1.3c.2: ECDSA recomendado):

| Campo ou extensão | Valor | Norma |
|---|---|---|
| versão, número de série | v3; 16 octetos aleatórios, positivo | §6.1.1; tab. 6.1.3c.3-1 (até 20 octetos) |
| assinatura | `ecdsa-with-SHA384`, da CA P-384 | §6.1.1 |
| *Subject* | `C=FR, O=5gc.mnc093.mcc208.3gppnetwork.org, CN=<FQDN>`; para o AF interno, `CN=af1.akma.example` | tab. 6.1.3c.3-1 (C e O); §6.1.3c.3 NOTA (o CN, quando existe, é um dos valores do `subjectAltName`); RFC 5280 §4.1.2.6 (*Subject* único por entidade) |
| validade | 1 ano | tab. 6.1.3c.3-1: *"3 years or less"*; T3 |
| `keyUsage` | crítico; `digitalSignature` | tab. 6.1.3c.3-1 |
| `extendedKeyUsage` | não crítico; `clientAuth` e `serverAuth` (o AF interno, só `clientAuth`) | idem: *"NF that may be both client and server shall have both OIDs set"* |
| `authorityKeyIdentifier`, `subjectKeyIdentifier` | não críticos; método (1) | idem |
| `cRLDistributionPoints` | não crítico; `URI:http://192.168.16.130:8030/sba-ca.crl` | idem; T26 |
| `subjectAltName` | **crítico**; `URI:urn:uuid:<nfInstanceId>`; `DNS:<FQDN>`; `IP:<endereço>`; nos servidores, também `URI:https://<FQDN>:<porta>`, o `apiRoot` | idem: o URI-ID do NF Instance ID (*shall*), o DNS-ID (*shall* no servidor), o URI do `apiRoot` (*should*), *"shall not contain only IP address in TLS server certificates"*; o IP, para o modo por IP (DV-17) |
| `nfTypes` | não crítico; uma `IA5String` com o tipo da NF (seção 3.1) | idem; RFC 9310 §3 |

Nenhuma outra extensão: sem `basicConstraints` (§6.1.1 manda marcar como não crítica toda extensão que a norma não exige, e o perfil não a pede, como nos certificados do CCF no MVP-9), sem `authorityInfoAccess` e sem a extensão *TLS feature* (não há OCSP; DV-10).

O NF Instance ID de cada certificado é o `nfInstanceId` fixo da configuração (T8). A lista das entidades, com o arquivo, o tipo, o FQDN, o endereço, a porta e o UUID, fica num arquivo só, `free5gc/config/sba-entities.conf`, que o *script* lê e a conferência da seção 4.7 cruza com as configurações.

### 4.3 Certificado de assinatura dos *tokens* do NRF (TS 33.310 tab. 6.1.3c.3-2)

`nrf-token.pem` e `nrf-token.key`, ECDSA P-256, para o `ES256` (T11):

- `keyUsage` crítico, `digitalSignature`;
- `extendedKeyUsage` não crítico, `id-kp-oauthAccessTokenSigning` (1.3.6.1.5.5.7.3.39, RFC 9509);
- `subjectAltName` crítico, só `URI:urn:uuid:<nfInstanceId do NRF>` (NOTA 1b: o `iss`);
- o resto como na seção 4.2: AKI, SKI, ponto de distribuição de CRL e `nfTypes` `NRF` (*"the rest of extensions remains the same"*; achado 14);
- *Subject* `C=FR, O=5gc.mnc093.mcc208.3gppnetwork.org, CN=NRF access token signing`. `[INTERPRETAÇÃO]` Sem DNS-ID: o certificado não é de TLS, e a regra do DNS-ID é dos certificados TLS.

O URI opcional do *Token Endpoint* (*"subjectAltName may ... contain URI-ID with the HTTPS URI for the "Token Endpoint""*) não entra.

### 4.4 Certificado do servidor interno do CCF

`ccf-internal.pem` e `ccf-internal.key`, ECDSA P-256, da CA do operador. Perfil de TS 33.310 §6.1.3a (TS 33.122 §6.6): `keyUsage` crítico `digitalSignature`; `extendedKeyUsage` `serverAuth`; `subjectAltName` com `DNS:ccf.5gc.mnc093.mcc208.3gppnetwork.org` e `IP:127.0.0.1`; ponto de distribuição de CRL; AKI e SKI. Sem `nfTypes`: o CCF não é um tipo de NF de TS 29.510 (achado 11).

### 4.5 Certificados de teste (T25)

| Arquivo | Para quê | Diferença |
|---|---|---|
| `seed-ausf` | o `seed` no papel do AUSF (`register`, `remove`) | `nfTypes` `AUSF`; UUID próprio; só `clientAuth` |
| `seed-af` | o `seed` no papel do AF no NRF (`trigger-auth`) | `nfTypes` `AF`; UUID próprio |
| `af2` | um AF interno fora da política do AF_ID de `af1` (T13) | `nfTypes` `AF`; `DNS:af2.akma.example` |
| `revoked` | conferência da CRL (T7) | `nfTypes` `AF`; revogado na CRL |
| `no-crldp` | certificado sem ponto de distribuição | sem `cRLDistributionPoints` |
| `other-ca`, `other-ca-ca` | certificado de outra CA | uma CA descartável, que só assina esse certificado |

Os testes de servidor com outro nome usam o certificado de outra NF (seção 9.2).

### 4.6 `make-sba-certs.sh`

```bash
free5gc/scripts/make-sba-certs.sh <pasta> <arquivo de entidades> [all|crl|revoke <nome>|check]
```

- `all` (padrão): a CA, se não existe ou se a chave não é P-384; cada certificado que falta, que vence em menos de 30 dias ou cuja linha do arquivo de entidades mudou; os de teste; a CRL.
- `crl`: assina de novo a CRL, com o `nextUpdate` uma semana à frente (seção 5.1).
- `revoke <nome>`: revoga e assina de novo a CRL.
- `check`: a conferência da seção 4.7.

Com o OpenSSL 1.1.1 da .130:

- `openssl ca`, com o banco da CA (`index.txt`, `serial`, `crlnumber`) em `cert/sba/ca/`: emite, revoga e assina a CRL;
- a configuração da CA é própria, sem a seção `v3_ca` padrão (o mesmo cuidado do `make-capif-certs.sh` no MVP-9);
- o `nfTypes`, que o OpenSSL 1.1.1 não conhece, entra pela sintaxe genérica: `1.3.6.1.5.5.7.1.34 = ASN1:SEQUENCE:nftypes`, com `[nftypes]` e `t1 = IA5STRING:<tipo>`;
- o `extendedKeyUsage` do certificado de assinatura entra pelo OID;
- arquivos com `umask 077`; certificados e CRL legíveis por todos.

Na máquina do AF, `af/build/sba/` recebe por `scp`, no `e2e.sh`, só o que o AF usa: `sba-ca.pem`, `af1.pem` e `af1.key`, `af2`, `revoked`, `no-crldp` e `other-ca`.

### 4.7 Conferência dos perfis (critério 2)

`make-sba-certs.sh check` lê cada certificado (`openssl x509 -text`, `openssl asn1parse` para o `nfTypes`) e recusa, com a linha do motivo:

- versão diferente de 3, série com mais de 20 octetos ou não positiva;
- *Subject* sem `C` ou com `O` diferente do domínio de origem; `CN` que não está no `subjectAltName`;
- validade acima de 3 anos ou vencida; assinatura que não é ECDSA com SHA-256 ou SHA-384;
- `keyUsage` ausente, não crítico ou sem `digitalSignature`; `extendedKeyUsage` diferente do papel;
- AKI diferente do SKI da CA; SKI ausente;
- ponto de distribuição ausente ou diferente do da CA;
- `subjectAltName` não crítico; sem `urn:uuid:` de um UUID v4 igual ao da configuração; sem o DNS-ID do servidor; só com IP;
- `nfTypes` ausente, crítico ou com tipo diferente;
- outra extensão crítica.

O `e2e.sh` roda a conferência no grupo S (S05).

## 5. Revogação (T7)

### 5.1 As CRL

- `sba-ca.crl`: a CRL da CA do operador, v2, assinada com SHA-384 pela CA P-384 (TS 33.310 §6.1a), feita por `make-sba-certs.sh`. `nextUpdate` uma semana depois de `thisUpdate`. O `e2e-remote.sh` a assina de novo no R1 (`make-sba-certs.sh crl`), e o tutorial manda fazer o mesmo antes de subir o free5GC.
- `ccf-ca.crl`: a CRL da CA do CCF, feita pelo CCF (seção 8.4).

### 5.2 O serviço da CRL (T26)

Um programa novo e pequeno, `free5gc/crl/`, compilado por `free5gc/scripts/build-crl.sh`:

- serve `GET` e `HEAD` de `/<nome>.crl`, os arquivos de `cert/crl/`, com `Content-Type: application/pkix-crl` (RFC 5280 §4.2.1.13; RFC 2585 §4);
- escuta em `192.168.16.130:8030`, porta que o free5GC não usa (`nef.md` §2), só HTTP/1.1, sem TLS (TS 33.310 §7.1: o repositório *"shall be accessible without relying on a secure connection"*);
- não lista a pasta e responde 404 a qualquer outro caminho;
- sobe com `run.sh -crl`, antes do NRF (*patch* `run/0003`), e o `force_kill.sh` o para.

O ponto de distribuição dos certificados é por IP, `http://192.168.16.130:8030/<ca>.crl`: vale nos dois modos, inclusive onde os nomes não resolvem.

### 5.3 A conferência em cada ponta

No pacote novo `sbitls` (seção 6.2), para o certificado do outro lado, no cliente e no servidor, depois de o `crypto/tls` conferir a cadeia e a validade:

1. O certificado tem de ter um ponto de distribuição com URI `http`. Sem ele: recusa (T7; achado 4).
2. A CRL vem da memória, pela URI. Quando não está na memória, ou depois de `nextUpdate`, a ponta baixa de novo: tempo máximo de 2 s, até 1 MiB, sem redirecionamento, uma busca por URI de cada vez.
3. A CRL tem de ser v2, assinada pela CA que emitiu o certificado (`RevocationList.CheckSignatureFrom`), com o mesmo emissor e `thisUpdate ≤ agora < nextUpdate`.
4. Sem uma CRL válida, depois da tentativa de baixar: recusa. A falha fica fechada (T26).
5. O número de série na lista: recusa.

Na partida, cada ponta já baixa a CRL do ponto de distribuição do próprio certificado. A conferência roda em todo *handshake*, com a CRL da memória; o resultado não fica guardado além da conexão (TS 33.310 §6.3.1).

`[FREE5GC]` O `crypto/tls` do Go não confere CRL: o código é do projeto, com testes próprios (seção 9.1).

## 6. Fase (a): a PKI, os nomes e o TLS mútuo em toda a SBI

Pela revisão do T19, a fase (a) também leva o transporte de CAPIF-2e e a troca do grupo do NEF do MVP-6 para o CAPIF: com o NRF e o NEF em TLS mútuo, o AF de fora do MVP-6, sem certificado do operador, deixa de funcionar.

### 6.1 Onde cada mudança entra (T28)

| Mudança | Onde | Por quê |
|---|---|---|
| a configuração TLS, o cliente, a conferência da CRL e a identidade do outro lado | `free5gc/patches/openapi/0005-sbi-tls.patch`: pacote novo `sbitls`, `client.go`, `udr/ImsDataRepository/client.go` | um lugar só para as onze NFs, o AAnF, o CCF e o AF, que já compilam com a cópia do `openapi`; corrige F21 |
| a montagem dos URIs | o mesmo *patch*: `util.go` (`GetNFServiceUri`) | um *helper* só (seção 3.3); F24 |
| nenhum segredo do TLS em arquivo | `free5gc/patches/util/0001-no-tls-key-log.patch`, numa cópia do `util` v1.3.2 em `free5gc/build/util`, preparada por `free5gc/scripts/prepare-util.sh` | um lugar só para os servidores SBI e de métricas de todas as NFs; corrige F22 |
| as chaves da configuração, a partida do TLS, o servidor, o `fqdn` no perfil, os URIs próprios, o *helper* dos URIs e a passagem do estado TLS à verificação do *token* | um *patch* por NF: `amf/0002`, `ausf/0002`, `bsf/0001`, `chf/0001`, `nef/0003`, `nrf/0004`, `nssf/0001`, `pcf/0001`, `smf/0001`, `udm/0003`, `udr/0001`, todos `-sbi-tls.patch` | cada NF tem o seu `go.mod`, a sua configuração e o seu servidor; F24 |
| as configurações | `free5gc/patches/config/0004-sbi-tls.patch` | T8 |
| o serviço da CRL | `free5gc/crl/`, `build-crl.sh`, `free5gc/patches/run/0003-crl-option.patch` | seção 5.2 |
| o banco de testes do free5GC | `free5gc/patches/test/0002-sbi-tls.patch`: o `replace` do `util` | seção 6.9 |
| AAnF, AF, `seed`, testes em Python | código do projeto | seções 6.6 a 6.8 |

O *patch* de cada NF já passa o estado TLS da conexão (`c.Request.TLS`) à verificação do *token*, que na fase (a) continua a do free5GC. A fase (b) só mexe no `openapi`, no NRF, no AAnF e nas configurações: nenhum *patch* novo por NF.

### 6.2 `openapi/0005`: o pacote `sbitls` e o cliente

**Configuração.** `sbitls.Init(Config)`, chamada uma vez na partida de cada NF, com o certificado, a chave e as CA (`sbi.tls.pem`, `sbi.tls.key`, `sbi.tls.ca`). Ela:

- confere que a chave é a do certificado, que o certificado é da CA e tem o URI do NF Instance ID igual ao `nfInstanceId` da configuração;
- baixa a CRL do próprio ponto de distribuição;
- monta as configurações TLS de cliente e de servidor;
- escreve no log a identidade, sem chave nenhuma (seção 6.10).

**Perfil TLS** (TS 33.210 §6.2; RFC 9113 §9.2), nas duas pontas:

- TLS 1.3 e 1.2; nada abaixo;
- no TLS 1.2, só ECDHE com AEAD: `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256` (obrigatória), `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256`, e as mesmas com RSA;
- grupos X25519, P-256 e P-384;
- ALPN `h2`;
- sem renegociação, sem OCSP *stapling* (DV-10).

**Cliente.** `innerHTTP2Client` passa a usar a configuração do `sbitls`: as CA do operador, o certificado da NF, o SNI pelo *host* do URI e a conferência da CRL. Sem `InsecureSkipVerify`. Um pedido `https` antes de `Init` falha (`SBI TLS not configured`), sem cair num cliente inseguro. O cliente `http` (h2c) continua, para o banco de testes e o modo sem TLS.

**Servidor.** `sbitls.ServerTLSConfig()`: `RequireAndVerifyClientCert`, com as CA do operador em `ClientCAs` e a conferência da CRL em `VerifyConnection`. O certificado de cliente é pedido no *handshake*, nunca depois (achado 12).

**Identidade.** `sbitls.PeerNfInstanceId` e `sbitls.PeerNfTypes` leem do certificado do outro lado o URI `urn:uuid:` e o `nfTypes`. Um certificado com mais de um `urn:uuid:` é recusado.

### 6.3 `util/0001`: nenhum segredo do TLS em arquivo (F22)

`httpwrapper.NewHttp2Server` não abre nem cria o arquivo e não põe `KeyLogWriter`; o parâmetro fica, ignorado, para não mudar a assinatura. A pasta `log/<data>/key/`, que cada NF cria na partida, fica vazia. O AAnF deixa de passar o caminho e não grava mais `aanf_sslkey.log`.

### 6.4 As NFs do free5GC

O *patch* de cada NF:

| Parte | O quê |
|---|---|
| `go.mod` | `replace` do `openapi` (a cópia do projeto; o AMF, a de `openapi-test`) e do `util` |
| `pkg/factory/config.go` | `sbi.fqdn` (opcional; FQDN válido) e `sbi.tls.ca` (obrigatória com `https`) |
| partida (`pkg/service/init.go` ou o contexto) | `sbitls.Init`, com `https` |
| `internal/sbi/server.go` | `TLSConfig` do `sbitls` antes de `ListenAndServeTLS` |
| contexto | `fqdn` no perfil e a função do URI próprio pelo nome (seção 3.3) |
| `router_auth_check.go` e `AuthorizationCheck` | o estado TLS da conexão passado à verificação do *token* |

E o que muda em cada uma:

| NF | URI próprio | Consumidores |
|---|---|---|
| AMF | `GetIPv4Uri()`; `amfName` igual a `sbi.fqdn` | `util/search_nf_service.go` → *helper* |
| SMF | função nova no contexto, no lugar de `nf_profile.go` e das montagens nos consumidores | `util/search_nf_service.go` → *helper* |
| UDM, PCF | `GetIPv4Uri()` | `util/search_nf_service.go` → *helper*; no PCF, também o recurso de `bsf_service.go` |
| UDR, NSSF | `GetIPv4Uri()` (`GetIpv4Uri()` no NSSF) | — |
| AUSF, CHF | `context.Url` | o AUSF já usa o *helper* |
| NEF | `factory.SbiUri()` | `consumer/nrf_service.go` → *helper* |
| BSF | função nova no contexto | — |
| NRF | `GetSbiUri()` | — (seção 6.5) |

### 6.5 O NRF na fase (a)

*Patch* `nrf/0004-sbi-tls.patch`:

- o TLS da SBI passa a uma seção nova, `sbi.tls` (`pem`, `key`, `ca`), igual à das outras NFs. `configureOAuthMutualTLS` sai: com `https`, o servidor sempre exige o certificado de cliente da CA do operador;
- `sbi.fqdn` e o `fqdn` no perfil próprio;
- `sbi.cert` continua a ser o par que assina os *tokens* (`RS512`, com a chave RSA do free5GC), e `sbi.rootcert` continua a gerar os arquivos de certificado das NFs para a conferência atual do pedido de *token*. Os dois mudam na fase (b) (seção 7.1).

### 6.6 O AAnF

Código do projeto, sem *patch*: `sbi.fqdn`, `sbi.tls.ca` e `scheme: https` no `aanfcfg.yaml`; `sbitls.Init`; o servidor com o `TLSConfig` do `sbitls`; `NewHttp2Server` sem o caminho do registro de segredos; o `fqdn` no perfil; o `nrfUri` pelo nome. A política continua por endereço até a fase (b).

### 6.7 O AF e as ferramentas

- **AF interno.** Seção nova no `afcfg.yaml`, `sbiTls` (`certificate`, `key`, `ca`), com o certificado `af1` (seção 4.2); `aanfUri: https://aanf.5gc.mnc093.mcc208.3gppnetwork.org:8010`. O cliente do AAnF usa o `sbitls`.
- **AF de fora** (modo CAPIF). `apiScheme: https` e `aefCa`, a CA do operador, para conferir o certificado do NEF: TS 33.122 §6.5.2.3 supõe as pontas *"pre-provisioned with the appropriate credentials"*. No CAPIF-2e, o AF apresenta o certificado que o CCF emitiu na entrada.
- **Modo do NRF do AF de fora.** O código (`af/af/internal/nrf`) fica, com os testes Go, como o modo sem CAPIF no NEF; `afcfg-external.yaml` sai do repositório, e a documentação diz que o modo não é usado no laboratório (T34).
- **`seed`.** Opções `--cert`, `--key` e `--ca`; padrões por papel: `seed-ausf` em `register` e `remove`, `seed-af` em `trigger-auth`; na máquina do AF, `retrieve` usa o certificado `af1`. `--aanf` e `--nrf` passam a `https` pelo nome.
- **Testes em Python** (`aanf_api.py`, `nef_api.py`, `capif_api.py`). `ssl.SSLContext` com o certificado de cliente e as CA, HTTP/1.1. O Go serve HTTP/1.1 e HTTP/2 por ALPN sobre TLS.

### 6.8 Configuração (`config/0004`; T8)

Exemplo, o AUSF:

```yaml
configuration:
  sbi:
    scheme: https
    registerIPv4: 127.0.0.9
    bindingIPv4: 127.0.0.9
    port: 8000
    fqdn: ausf.5gc.mnc093.mcc208.3gppnetwork.org # by name (TS 29.510 cl. 6.1.6.2.2 NOTE 1); without it, by IP address (DV-17)
    tls:
      pem: cert/sba/ausf.pem # client and server certificate (TS 33.310 cl. 6.1.3c.3)
      key: cert/sba/ausf.key
      ca: cert/sba/sba-ca.pem # CA of the operator: the certificates of the peers (TS 33.310 cl. 6.3.1)
  nfInstanceId: <UUID fixo, o do certificado>
  nrfUri: https://nrf.5gc.mnc093.mcc208.3gppnetwork.org:8000
  nrfCertPem: cert/nrf.pem # phase (b): cert/sba/nrf-token.pem
```

O mesmo em cada NF, com o seu nome, endereço e UUID, e mais:

- o NRF: `sbi.tls` (seção 6.5);
- o AMF: `amfName` com o FQDN;
- o NEF: em `capif.interface`, `fqdn` no lugar de `ipv4Addr` (o NEF publica o nome; seção 8.1).

Na .130, *backup* de `config/` e um *snapshot* da VM antes da fase (a) (T8).

### 6.9 O banco de testes do free5GC (T37)

Continua em `http`. Ele monta as configurações em Go, sem `sbi.fqdn` nem `sbi.tls`, e roda todas as NFs num processo só, com o cliente global do `openapi`: um certificado de cliente só serviria a todas, e a identidade de cada NF não teria sentido. O `test/0002` acrescenta o `replace` do `util`. O banco tem de passar com os *patches* novos no modo `http`: `TestAkma` e o registro do UE com sessão PDU. O TLS fica coberto pelos testes Go dos pacotes (seção 9.1) e pelo laboratório (seção 9.2).

### 6.10 Logs

No formato do free5GC. Nunca aparecem chaves privadas, segredos de sessão, *tokens*, K_AKMA nem K_AF, e nenhum arquivo de segredos é criado.

| Onde | Nível | Mensagens |
|---|---|---|
| todas as pontas, categoria `SBI` (ou a da NF) | info | `SBI TLS: certificate of <FQDN>, NF instance <UUID>, type <tipo>, issued by <CA>, valid until <data>`; `CRL <URI> loaded: <n> revoked, next update <data>` |
| todas as pontas | warn | `TLS handshake with <endereço> refused: <motivo>` (sem certificado, outra CA, vencido, sem ponto de distribuição, revogado, CRL indisponível, nome diferente); `CRL <URI> not available: <erro>`; o aviso do DV-17 (seção 3.2) |

### 6.11 Erros

Uma recusa do TLS é um alerta no *handshake*, sem resposta HTTP. O `crypto/tls` do Go manda `bad_certificate` para os certificados recusados pela conferência do projeto, e `certificate_required` no TLS 1.3 sem certificado de cliente.

## 7. Fase (b): a identidade nos *tokens* e nas políticas

### 7.1 O NRF: o certificado e o pedido (T9, T32)

*Patch* `nrf/0005-token-identity.patch`. O NRF lê o certificado da sessão TLS (`c.Request.TLS.PeerCertificates[0]`), e não o arquivo que gerava:

| Operação | Conferência | Recusa |
|---|---|---|
| `POST /oauth2/token` | `nfInstanceId` igual ao URI do certificado; `nfType` presente no `nfTypes`; `requesterPlmn`, se vier, igual à PLMN do `O` do certificado | `400` `invalid_client`, com `Cache-Control: no-store` e `Pragma: no-cache` (TS 29.510 §5.4.2.2.1; tab. 6.3.4.2.2-4) |
| `PUT`, `PATCH` e `DELETE /nnrf-nfm/v1/nf-instances/{nfInstanceID}` | `{nfInstanceID}` igual ao URI; no `PUT`, `nfType` do perfil presente no `nfTypes` | `403` com `ProblemDetails`, sem `cause` (TS 29.500 tab. 5.2.7.1-1) |
| `GET /nnrf-disc/v1/nf-instances` | `requester-nf-type` presente no `nfTypes`; `requester-nf-instance-id`, se vier, igual ao URI | `403` com `ProblemDetails`, sem `cause`. `[3GPP]` TS 33.501 §13.3.1.3 NOTA 1a: *"This verification is implementation specific"* |

Saem `SignNFCert`, a conferência do arquivo, a tolerância a *"unknown authority"* e o `rootcert` (F25).

### 7.2 O NRF: o perfil JWS e o `iss` (T11, T29, T30)

- O NRF assina com a chave de `nrf-token.key` (`sbi.cert` passa a apontar para `nrf-token.pem` e `nrf-token.key`), em `ES256`.
- Cabeçalho `{"alg":"ES256","kid":"<kid>","typ":"JWT","x5t#S256":"<impressão>"}`. O `kid` é a impressão digital da chave pela RFC 7638, como no CCF. O `x5t#S256` é a impressão SHA-256 do certificado de assinatura (TS 33.210 §6.3.3.4: *"should be supported"*). TS 33.210 §6.3.3.1: `ES256` e `kid`, *"shall be supported"*; V20.0.0: `RS512` *"not recommended"*.
- `iss` passa a ser o `nfInstanceId` do NRF, o do URI do certificado de assinatura (F26). O NRF não parte se os dois forem diferentes.
- O *token* que o NRF emite para si mesmo (`GetTokenCtx`) segue a mesma forma.

### 7.3 Escopos por tipo (T12, T31; F4)

O AAnF registra no serviço `naanf-akma`:

```json
"allowedOperationsPerNfType": {
  "AUSF": ["naanf-akma", "naanf-akma:anchorkey"],
  "NEF":  ["naanf-akma", "naanf-akma:applicationkeyget", "naanf-akma:applicationkeyget:supi-access"],
  "AF":   ["naanf-akma", "naanf-akma:applicationkeyget", "naanf-akma:applicationkeyget:supi-access"]
}
```

`[3GPP]` TS 29.535 tab. 4.2.2.1-1: o AUSF inicia o `AnchorKey_Register`; o NEF e o AF, o `ApplicationKey_Get`. TS 33.535 §6.2.1 passo 8 e §6.3 passo 4: o AAnF manda o SUPI ao AF interno e ao NEF. O `ContextRemove` é iniciado *"e.g. OAM"*, que não é tipo de NF; no laboratório, o `seed` o faz no papel do AUSF, com `naanf-akma:anchorkey` (*"store and remove"*, tab. 5.1.9-1).

O NRF, pela NOTA 11 da tab. 6.1.6.2.3-1:

- com `allowedOperationsPerNfType`, só concede o escopo que está na lista do tipo do consumidor; um tipo fora da lista não recebe nenhum escopo do serviço (*"any NF type ... not listed in these IEs is disallowed"*);
- sem a lista, como hoje: só o escopo do serviço e `allowedNfTypes`;
- escopo pedido fora da lista: `400` `invalid_scope`.

`allowedOperationsPerNfInstance` não entra.

O AAnF, com *token*:

- `register-anchorkey` e `remove-context` exigem `naanf-akma` e `naanf-akma:anchorkey`;
- `retrieve-applicationkey` exige `naanf-akma` e `naanf-akma:applicationkeyget`;
- falta de escopo: `403`, `WWW-Authenticate: Bearer error="insufficient_scope", scope="<o que falta>"` (TS 29.500 §6.7.3);
- sem `naanf-akma:applicationkeyget:supi-access`, a resposta vem sem `supi`, como no acesso anônimo. `[INTERPRETAÇÃO]` A tab. 5.1.9-1 liga esse escopo à devolução do SUPI, não ao acesso.

Os pedidos sem *token* do AF interno seguem a política local (seção 7.5).

### 7.4 O produtor: a verificação do *token* (T10, T29; F5, F12)

*Patch* `openapi/0006-access-token.patch`. A verificação recebe o *token*, o nome do serviço, o estado TLS e o certificado de assinatura do NRF (`nrfCertPem`), e confere, na ordem:

1. **Forma.** `Authorization: Bearer <jwt>`. Sem *token*: `401`, `WWW-Authenticate: Bearer realm="<URI da API>"`.
2. **Assinatura.** `alg` `ES256` (os outros são recusados, inclusive `RS512` e `none`); `kid`, se vier, igual ao da chave do certificado; `x5t#S256`, se vier, igual ao certificado; o certificado de assinatura da CA do operador, com `id-kp-oauthAccessTokenSigning` (RFC 9509 §3) e não revogado. Falha: `401`, `error="invalid_token"`.
3. ***Claims* presentes.** `iss`, `sub`, `aud`, `scope` e `exp`. Falta alguma: `401`, `error="invalid_token"`, com `ProblemDetails` `ACCESS_TOKEN_CLAIM_MISSING` e `invalidParams`.
4. **Validade.** `exp` no futuro (TS 33.501 §13.4.1.1.2), sem tolerância: as duas máquinas usam NTP (ADR-0001-R2 §5.3). Vencido: `401`, `error="invalid_token"`.
5. **Valores.** `iss` igual ao URI do certificado de assinatura; `aud` igual ao `nfInstanceId` do produtor ou ao seu tipo; `sub` igual ao URI do certificado de cliente da sessão (comunicação direta). Diferença: `403`, `WWW-Authenticate: Bearer error="invalid_token"`, `ProblemDetails` `ACCESS_TOKEN_CLAIM_MISMATCH` e `invalidParams` com a *claim*.
6. **Escopo.** O nome do serviço em `scope`. Falta: `403` `insufficient_scope` (passo do item 7.3).

As respostas seguem TS 29.500 §6.7.3, com `ProblemDetails` em `application/problem+json`, em todas as NFs, pela mesma função do *patch*. O 401 do NEF, que era `{"error": ...}` (F12), passa a seguir a norma (T29).

O AAnF, código do projeto, faz as mesmas conferências com a mesma função; `VerifyAccessToken` sai.

### 7.5 O AAnF: o AF e o NEF pelo certificado (T13, T33)

No `aanfcfg.yaml`, `sourceAddresses` sai das duas políticas:

```yaml
    authorizedAfs: # TS 33.535 cl. 6.2.1 step 2; the AF is the NF instance of its client certificate (TS 33.501 cl. 13.1.0)
      - afId: af1.akma.example.0100011301
        nfInstanceIds:
          - <UUID do af1>
    authorizedNefs: # TS 33.535 cl. 6.3; the NEF is the NF instance of its client certificate
      - nfInstanceIds:
          - <UUID do NEF>
        afIds:
          - af1.akma.example.0100011301
```

No `retrieve-applicationkey`:

1. A identidade é o URI do certificado de cliente; o tipo, o `nfTypes`.
2. Com *token*: as conferências da seção 7.4, inclusive `sub` igual à identidade. Se a identidade está em `authorizedNefs` e o tipo é `NEF`, o pedido é do NEF: o `afId` tem de estar nos `afIds` da entrada.
3. Sem *token*: a identidade tem de estar em `nfInstanceIds` do `afId` pedido e o tipo tem de ser `AF`.
4. Senão: `403` sem `cause`, como hoje.

O endereço de origem sai da política e fica só no log. O risco do DV-1, um AF apresentar o AF_ID de outro, fica fechado pelo certificado.

## 8. Fase (c): o NEF com o AF de fora e com o CCF

### 8.1 CAPIF-2e com TLS mútuo (T14, T34)

O transporte entra na fase (a) (T19); a ligação do certificado ao *token*, aqui.

- **Confiança do NEF.** O servidor do NEF aceita certificados de cliente de duas CA: a do operador e a do CCF (`capif.invokerCa: cert/ccf-ca.pem`). `[MVP]` Cada grupo de rotas confere de qual CA veio o certificado: os grupos de CAPIF-2e (`3gpp-akma` e `aef-security`) só aceitam a CA do CCF, e os da SBI só a do operador. Outra: `403` sem `cause`.
- **`check-authentication`.** O CN do certificado tem de ser o `apiInvokerId` do corpo; senão, `403`. O NEF pede os dados do *invoker* ao CCF (`capif.md` §5.9) e confere que o certificado do *invoker* é da CA de `authenticationInfo`.
- **`retrieve`.** Além de `capif.md` §6.4: o CN do certificado de cliente tem de ser igual ao `client_id` do *token*; senão, `401`, `WWW-Authenticate: Bearer error="invalid_token"`. `[INTERPRETAÇÃO]` O *token* não vale para esse cliente.
- **Revogação.** O NEF confere o certificado do *invoker* contra a CRL da CA do CCF (seção 8.4).
- **Fim do modo do NRF.** `allowNrfTokens: false`: um *token* do NRF em `3gpp-akma` recebe `401` (`capif.md` §15.2). O modo sem CAPIF fica no código, com os testes Go.
- **A política de endereço.** No modo CAPIF, `akma.authorizedAfs` do NEF deixa de conferir o endereço: `sourceAddresses` fica opcional e só vale no modo sem CAPIF. A lista de AF_ID continua (o CP17).
- **A indicação do TLS mútuo** (achado 10). `[LACUNA]` TS 29.222 não tem atributo para ela. A política do CCF é sempre o TLS mútuo, e a indicação é a CA do *invoker* em `authenticationInfo` (TS 33.122 §6.5.2.3 passo 5). O AF apresenta o certificado quando o NEF o pede, no *handshake*. O `check-authentication` também vai sobre o TLS mútuo.
- **Publicação.** No modo por nome, o NEF publica `fqdn` em `InterfaceDescription` (tab. 8.2.4.2.3-1). O AF monta o `apiRoot` com `apiScheme: https` (`capif.md` C26; D-AG).
- **CAPIF-1e** continua por IP, com o `ccf.pem` da CA do CCF, que ganha o ponto de distribuição de CRL. `[INTERPRETAÇÃO]` A norma não dá nome ao CCF, e o URI por IP é válido no `https` (RFC 9110 §4.3.4; RFC 9525) (T34).

### 8.2 CAPIF-3/4: o servidor interno do CCF em `https` (T15, T35)

- `internal` no `ccfcfg.yaml` ganha `fqdn`, `certificate`, `key` e `ca`: `https://ccf.5gc.mnc093.mcc208.3gppnetwork.org:8020`, em 127.0.0.1, com TLS mútuo e a CA do operador, HTTP/1.1 e HTTP/2 por ALPN.
- O CCF reconhece o provedor pelo certificado: cada entrada de `providers` ganha `nfInstanceId`, o UUID do NEF. Só esse cliente publica e retira APIs desse `apfId` e pede dados de *invokers* para os seus `aefIds` (C6). Outro certificado: `403`.
- `Obtain_API_Invoker_Info` passa a devolver só os itens do AEF que pergunta. Fecha a limitação de `capif.md` §5.9 e §11.
- O NEF passa a `ccfUri: https://ccf.5gc.mnc093.mcc208.3gppnetwork.org:8020`.

### 8.3 O AF de fora

A conta do AF no CCF não muda. Ele confere o certificado do NEF pela CA do operador e pelo nome publicado, apresenta o certificado da entrada e confere a CRL do CCF no `ccf.pem` e a do operador no certificado do NEF.

### 8.4 A CRL da CA do CCF e a revisão do DV-16 (T36)

- Os certificados que o CCF emite na entrada e o `ccf.pem` ganham `cRLDistributionPoints` `URI:http://192.168.16.130:8030/ccf-ca.crl` (TS 33.310 §6.1.3a).
- O CCF assina a CRL com a CA do CCF e a grava em `cert/crl/ccf-ca.crl`: na partida, a cada mudança e um dia antes de `nextUpdate`, com `nextUpdate` de uma semana.
- Na saída de um *invoker* (TS 33.122 §6.8), o CCF revoga o certificado dele.
- O CCF confere a CRL no CAPIF-1e; o NEF e o AF, nas suas pontas.
- **DV-16.** A parte da revogação se encerra. Ficam o OCSP *stapling*, como no DV-10, e o `x5t#S256` do *token* do CCF, cuja chave continua sem certificado.

## 9. Verificação

### 9.1 Testes Go

| Onde | Casos |
|---|---|
| `openapi/sbitls` | o perfil TLS (versões, suítes, grupos); cliente e servidor com TLS mútuo; outra CA; nome diferente; IP sem SNI; SNI com nome; sem ponto de distribuição; revogado; CRL vencida, de outra CA, inacessível e grande demais; a CRL baixada uma vez por URI; o `nfTypes` e o URI do certificado; dois `urn:uuid:` recusados |
| `openapi` (`GetNFServiceUri`) | `fqdn` do serviço, do perfil, com e sem porta, com `apiPrefix` caminho e absoluto, só IP |
| `openapi/oauth` | cada passo da seção 7.4, com o código, os cabeçalhos e a causa |
| `util` | nenhum arquivo criado com o caminho do registro de segredos |
| NRF | a seção 7.1 (cada recusa e o caso certo), o `ES256` com `kid`, `x5t#S256` e `iss`, a NOTA 11 da seção 7.3 |
| AAnF | as seções 7.3 a 7.5; a política sem endereço; o TLS mútuo no servidor |
| NEF | a seção 8.1: a CA por grupo, o CN e o `client_id`, o `check-authentication`, `allowNrfTokens: false` |
| CCF | a seção 8.2 e a CRL da seção 8.4 |
| AF | o `sbiTls`, o `aefCa`, o `apiScheme: https` |

### 9.2 `e2e.sh`: grupo S (T38)

Da máquina do AF e, por SSH, da do free5GC. As portas são as onze NFs, o AAnF, o NEF e o servidor interno do CCF.

| ID | Verificação | Esperado |
|---|---|---|
| S01 | cada porta da SBI com texto claro (h2c) | conexão recusada no TLS |
| S02 | cada porta com TLS sem certificado de cliente | alerta, sem resposta HTTP |
| S03 | cada porta com o certificado `other-ca` | alerta |
| S04 | o certificado de servidor de cada porta, com o SNI do nome | o DNS-ID e o URI do UUID da configuração; o perfil da seção 4.7 |
| S05 | `make-sba-certs.sh check` em todos os certificados | nenhuma recusa (critério 2) |
| S06 | o certificado `revoked` no AAnF, no NRF e no NEF | alerta (T7) |
| S07 | o certificado `no-crldp` no AAnF | alerta |
| S08 | o `seed` contra um `openssl s_server` com o certificado do UDM no nome do AAnF, e com o `other-ca` | o cliente recusa os dois (critério 3) |
| S09 | o `seed` contra um `openssl s_server` com o nome certo | o `s_server` registra o SNI com o nome |
| S10 | pedido de *token* com o NF Instance ID de outro certificado; com o tipo de outro certificado; registro com o UUID de outro certificado | `400` `invalid_client`; `400` `invalid_client`; `403` (T9) |
| S11 | *token* de `seed-ausf` usado com o certificado `seed-af` no AAnF | `403` `ACCESS_TOKEN_CLAIM_MISMATCH` (`sub`, T10) |
| S12 | *token* para outro tipo de produtor no AAnF | `403` `ACCESS_TOKEN_CLAIM_MISMATCH` (`aud`) |
| S13 | o cabeçalho e as *claims* de um *token* do NRF | `ES256`, `kid`, `x5t#S256`; `iss` igual ao UUID do NRF (T11, F26) |
| S14 | o AUSF pede `naanf-akma:applicationkeyget`; `retrieve` com *token* só de `naanf-akma` | `400` `invalid_scope`; `403` `insufficient_scope` (T12) |
| S15 | `retrieve` com o certificado `af2` e o AF_ID de `af1` | `403` sem `cause` (T13) |
| S16 | no NEF: *token* de um *invoker* com o certificado de outro; *token* do NRF; o certificado de *invoker* num grupo da SBI | `401` `invalid_token`; `401`; `403` (T14) |
| S17 | no servidor interno do CCF: sem certificado; com o certificado de outra NF | alerta; `403` (T15) |
| S18 | o certificado de um *invoker* que saiu, no NEF | alerta (revogado; DV-16) |
| S19 | nenhum `*sslkey.log` com conteúdo e nenhuma chave nos logs das duas máquinas | nenhum (critério 8) |

As verificações entram por fase: S01 a S09 e S19 na (a); S10 a S15 na (b); S16 a S18 na (c).

### 9.3 O grupo do NEF do MVP-6 no modo CAPIF (fase a)

- Os AF extras do grupo (`nef`, `nef-anon` e `nef-psk`) passam a *invokers* do CCF, cada um com a sua entrada e o seu estado (`capif_variant`). O `anon`, AF interno anônimo, só ganha o TLS.
- O E00 passa a conferir a entrada no CCF, e não o registro no NRF.
- O `nef_api.py` obtém o *token* do CCF (as funções de entrada do `capif_api.py`) e faz o E01 a E12 sobre TLS mútuo. O E05 passa a ser o *token* do NRF para `3gpp-akma`, obtido pelo `seed-af` na .130, recusado com `401` (na fase c, com `allowNrfTokens: false`; antes, o E05 confere o 401 de *token* para outro serviço, como hoje).
- O túnel SSH e o `stop_extra` dele saem.

### 9.4 As três verificações puladas

O D03, o O04 e o O05 voltam, com o OpenSSL 4.1 e o `coap-client` da libcoap compilados na máquina do AF fora do *scratchpad*, pelos passos de `docs/af-side/e2e.md` (`OPENSSL_DTLS13` e `LIBCOAP_CLIENT`). Não dependem do TLS da SBI.

### 9.5 Em cada fase

O `e2e.sh` completo, com código 0, inclusive o registro do UE com sessão PDU (R1) e o Ua\* com o AF interno, o de fora e o anônimo. O `TestAkma` e o registro do banco de testes do free5GC, com os *patches* da fase. `go vet` e `go test -race` nos módulos do projeto e nas NFs com *patches* (critério 9).

## 10. Registros

- **Desvios.** O DV-18 (T27, opção a) está registrado desde 04/11/2026. Na implementação: DV-1 e DV-9 encerrados no fim da fase (c); DV-16 reduzido ao OCSP e ao `x5t#S256` (T36); DV-17 e DV-18 implementados.
- **Correções do free5GC.** F21 a F26 estão em `correcoes-free5gc.md`. Na implementação: F4, F5 e F12 corrigidos no projeto (T12, T29).
- **Defeitos.** D-AJ (a indicação do TLS mútuo no método 3, achado 10) e D-AK (a referência da tab. 6.1.3c.3-2, achado 14) estão no registro de defeitos e na cópia em `~/akma`, desde 04/11/2026.

## 11. Limitações

- O banco de testes do free5GC continua em `http` (seção 6.9).
- `apiPrefix` continua com o URI inteiro, fora de TS 29.501 §4.4.1 (F24; T22).
- A CRL só por HTTP, sem LDAPv3 (DV-18); sem OCSP nem OCSP *stapling* (DV-10, DV-16).
- Sem a inscrição automática de certificados (TS 33.310 §10.2 e §10.3): os certificados saem do *script*, e a chave da CA fica na .130.
- A *webconsole* do free5GC, que o laboratório não usa (a assinatura sai do `mongosh`), não fala com o NRF em TLS mútuo.
- O CHF tem outros servidores TLS que não são SBI (o CGF e o ABMF, `chfcfg.yaml`), com os certificados do free5GC: ficam como estão.
- Um CCF, um provedor e um AEF, como no MVP-9.

## 12. Fora do escopo

O que o plano já tirou: NDS/IP, o SEPP e o N32; o SCP, a comunicação indireta e o CCA; o OCSP; a inscrição automática; os métodos 1 e 2 de CAPIF-2e (DV-14); o Ua\*; o GPSI; o *roaming*; o repositório público. E mais: `allowedOperationsPerNfInstance`; as interfaces de NGAP, GTP-U e PFCP; os servidores de métricas, que continuam em `http` (o `util/0001` só tira deles o registro de segredos).

## 13. Pontos para aprovação

| Nº | Ponto | Proposta | Base e alternativas |
|---|---|---|---|
| T18 | Correções do texto aprovado | (a) T9: o pedido de *token* com dados diferentes do certificado recebe `400` `invalid_client`, e não 403; o registro, a atualização e o descadastro recebem `403`; (b) T7: a regra da recusa sem ponto de distribuição cita NE e SEG; para a entidade TLS, a recusa vem da frase anterior de §6.3.1, por interpretação; o comportamento não muda; (c) T4: TS 23.003 também dá o nome do NSSF e o do AMF; (d) o arquivo de segredos é `<nf>sslkey.log` | `[3GPP]` TS 29.510 §5.4.2.2.1 e tab. 6.3.4.2.2-2; TS 29.500 tab. 5.2.7.1-1; TS 33.310 §6.3.1; TS 23.003 §28.3.2.4.2 e §28.3.2.5. Alternativa para (a): manter 403, contra TS 29.510 |
| T19 | Fases revistas | a fase (a) também leva o transporte de CAPIF-2e (o NEF aceita a CA do CCF nos grupos de CAPIF-2e; o AF usa `https`) e a troca do grupo do NEF do MVP-6 para o CAPIF, sem túnel; a fase (c) fica com a ligação do certificado ao *token*, `allowNrfTokens: false`, a CRL do CCF e CAPIF-3/4 | `[MVP]` com o NRF e o NEF em TLS mútuo na fase (a), o AF de fora do MVP-6, sem certificado do operador, para de funcionar, e o `e2e.sh` da fase falharia. Alternativas: (b) juntar as fases (a) e (c); (c) deixar o NEF em `http` até a fase (c), contra o T2 |
| T20 | Nomes das NFs | a tabela da seção 3.1: o NRF e o NSSF pelo formato de TS 23.003; o AMF com o AMF Name `amf1.amf.5gc...` e `amfName` igual; as outras `<tipo>.5gc...`; o BSF configurado, mas sem subir | `[3GPP]` TS 23.003 §28.2, §28.3.2.3.2, §28.3.2.4.2 e §28.3.2.5; TS 29.510 tab. 6.1.6.2.2-1. `[INTERPRETAÇÃO]` as outras NFs. Alternativa: manter `amfName: AMF` e registrar o AMF por outro nome, contra TS 29.510 |
| T21 | Chave do modo | `sbi.fqdn` (`internal.fqdn` no CCF): com ela, por nome; sem ela, por IP (DV-17), com aviso quando há `https`. No código, o padrão é por IP; nas configurações do projeto, por nome. Os clientes seguem a forma do URI | `[3GPP]` TS 29.510 §6.1.6.2.2 NOTA 1; RFC 6066 §3. Alternativa: uma chave explícita, `sbi.addressing: fqdn|ip`, além do nome |
| T22 | Montagem dos URIs | um *helper* só, `openapi.GetNFServiceUri`, corrigido pela norma e chamado pelos outros cinco; `apiPrefix` continua com o URI inteiro, agora com o nome, para o SMF e o PCF que o usam direto; a correção do `apiPrefix` fica aberta em F24 | `[3GPP]` TS 29.510 tab. 6.1.6.2.3-1 NOTA 1 e tab. 6.1.6.2.5-1 NOTAS 2 a 4; TS 29.501 §4.4.1. Alternativa: `apiPrefix` só com o caminho, como a norma, trocando os usos diretos no SMF e no PCF |
| T23 | URIs próprios | no modo por nome, o perfil, os *callbacks* e o `Location` saem com o nome, pela função única de cada NF | `[3GPP]` TS 29.501 §4.4.3 (o *callback* pode ter IP); TS 33.310 §6.1.3c.3 (DNS-ID do *callback*, *should*). Alternativa: *callbacks* por IP, com mais trocas por NF |
| T24 | Perfil dos certificados | as seções 4.1 a 4.4: *Subject* com `CN=<FQDN>`; `subjectAltName` crítico com o URI do UUID, o DNS-ID, o IP e, no servidor, o URI do `apiRoot`; sem `basicConstraints`; série aleatória de 16 octetos; 1 ano; P-256 assinado com SHA-384; `cert/sba/`; `sba-entities.conf`; o certificado de assinatura do NRF sem DNS-ID; o do CCF interno sem `nfTypes` | `[3GPP]` TS 33.310 §6.1.1, §6.1.3a, §6.1.4a e tab. 6.1.3c.3-1 e -2; RFC 9310; RFC 9509; RFC 5280 §4.1.2.6. Alternativas: sem o IP no `subjectAltName` (só possível sem o modo por IP); `nfTypes` `CCF` como tipo do operador (RFC 9310 §5) |
| T25 | Certificados de teste | os da seção 4.5; o `seed` com duas identidades, AUSF e AF; na máquina do AF, o `retrieve` com o `af1` | os testes dos critérios 1 a 6. Alternativa: um certificado de teste por verificação |
| T26 | CRL e serviço | ponto de distribuição por IP, `http://192.168.16.130:8030/<ca>.crl`; `nextUpdate` de uma semana, assinada de novo no R1; o programa `free5gc/crl/`, com `run.sh -crl`; a falha fechada; a CRL baixada na partida e conferida em todo *handshake* | `[3GPP]` TS 33.310 §6.1a, §6.3.1 e §7.1. Alternativas: o ponto de distribuição por nome; `python3 -m http.server`; o CCF servindo as CRL; uma tolerância à CRL vencida |
| T27 | **`[DIVERGÊNCIA]` CRL sem LDAPv3** | (a) só HTTP, registrado como DV-18: TS 33.310 §6.1a manda suportar o LDAPv3 como método principal, e o laboratório não tem diretório; o mesmo §6.1a permite o HTTP para certificados TLS | `[3GPP]` TS 33.310 §6.1a. Impacto: um operador com a CRL só em LDAP não serve a este código. Alternativas: (b) ler o *"HTTP may be used"* como suficiente para entidades TLS, sem desvio; (c) o LDAPv3 também, com um servidor de diretório na .130 e um cliente LDAP em Go |
| T28 | Onde as mudanças entram | a tabela da seção 6.1: `openapi/0005` e `0006`; `util/0001` com `prepare-util.sh`; um *patch* por NF, só na fase (a), já com o estado TLS levado à verificação do *token*; o NRF com `0004` (a) e `0005` (b); `config/0004` a `0006`; `run/0003`; `test/0002` | o menor número de *patches*. Alternativas: sem a cópia do `util`, com a troca em cada NF (o `main.go`, o servidor SBI e o de métricas); *patches* por NF em cada fase |
| T29 | Verificação do *token* | a seção 7.4: só `ES256`; `kid` e `x5t#S256` no cabeçalho do NRF; a cadeia, o EKU e a revogação do certificado de assinatura; `iss`, `aud`, `sub` e `exp`, sem tolerância; as respostas de TS 29.500 §6.7.3 em todas as NFs, o que corrige F12 no NEF | `[3GPP]` TS 33.501 §13.4.1.1.2 (o `iss`, que o T10 não trazia); TS 33.210 §6.3.3.1 e §6.3.3.4; RFC 9509 §3; TS 29.500 §6.7.3. Alternativas: manter o corpo `{"error": ...}` do free5GC (F12 aberto); sem `x5t#S256` |
| T30 | F26 | registrar em `correcoes-free5gc.md`: o NRF manda `iss` vazio e o próprio certificado sem o URI; corrigido no `nrf/0005` | `[3GPP]` TS 29.510 tab. 6.3.5.2.4-1; TS 33.501 §13.4.1.1.2; TS 33.310 NOTA 1b |
| T31 | Escopos por tipo | a lista da seção 7.3; o AAnF exige `anchorkey` no registro e na remoção, e `applicationkeyget` no `retrieve` com *token*; sem `supi-access`, `200` sem `supi` | `[3GPP]` TS 29.535 tab. 4.2.2.1-1 e 5.1.9-1; TS 29.510 tab. 6.1.6.2.3-1 NOTA 11; TS 33.535 §6.2.1 e §6.3. Alternativa: `403` sem `supi-access` |
| T32 | O NRF e o certificado | a tabela da seção 7.1: o pedido de *token*, o registro, a atualização, o descadastro e a descoberta; saem a geração de certificados e a conferência dos arquivos | `[3GPP]` TS 33.501 §13.4.1.1.2 passo 2 e §13.3.1.3 NOTA 1a; TS 29.510 §5.4.2.2.1. Alternativa: só o *token* e o registro, como no T9 |
| T33 | Política do AAnF | `nfInstanceIds` no lugar de `sourceAddresses`, em `authorizedAfs` e `authorizedNefs`; o tipo pelo `nfTypes` (seção 7.5) | `[3GPP]` TS 33.501 §13.1.0; TS 33.535 §6.2.1 passo 2. Alternativa: o certificado e o endereço juntos |
| T34 | CAPIF-2e e CAPIF-1e | a seção 8.1: a CA por grupo; o CN e o `client_id`; o `check-authentication` sobre TLS mútuo; a indicação pela CA em `authenticationInfo`; `sourceAddresses` opcional no NEF; `fqdn` publicado; CAPIF-1e por IP; `afcfg-external.yaml` fora | `[3GPP]` TS 33.122 §6.5.2.3; TS 33.501 §12.2 e §12.3. `[LACUNA]` achado 10 (D-AJ). Alternativas: um servidor separado para CAPIF-2e no NEF; CAPIF-1e com DNS-ID e `ccfUri` pelo nome |
| T35 | CAPIF-3/4 | a seção 8.2: `https` pelo nome em 127.0.0.1; o provedor pelo UUID do certificado; os dados do *invoker* só do AEF que pergunta | `[3GPP]` TS 33.122 §6.6; C6 do MVP-9 |
| T36 | DV-16 | a seção 8.4: a CRL da CA do CCF, a revogação na saída e a conferência nas pontas; o DV-16 fica só com o OCSP e o `x5t#S256` | `[3GPP]` TS 33.310 §6.1.3a e §6.3.1; TS 33.122 §6.8. Alternativa: fechar também o `x5t#S256`, com um certificado para a chave do *token* do CCF |
| T37 | Banco de testes do free5GC | continua em `http` (seção 6.9) | `[FREE5GC]` um processo, um cliente global. Alternativa: um cliente por NF no `openapi`, uma mudança grande no free5GC |
| T38 | `e2e.sh` | o grupo S da seção 9.2, por fase; o grupo do NEF do MVP-6 no modo CAPIF (seção 9.3); as três puladas de volta | critérios 1 a 8. Alternativa: os negativos de certificado só nos testes Go |
| T39 | Registros | a seção 10, com F26 (T30), DV-18 (T27), D-AJ e D-AK | a regra do registro de desvios |

**Aprovação de 04/11/2026.** O responsável aprovou os pontos T18 a T39 como propostos. No T27, ficou a opção (a): só HTTP, com o desvio DV-18. Registrados na mesma data: o DV-18, as notas do DV-9, do DV-16 e do DV-17, o F26, os defeitos D-AJ e D-AK e, no plano, os pontos T1, T4, T7, T9, T14 e T17 revistos.

## 14. Arquivos e *commits* por fase

Cada fase termina com o `e2e.sh` completo. O responsável faz os *commits*; as mensagens são propostas na implementação, em inglês, com a cláusula no corpo.

**Fase (a).**

1. `free5gc/scripts/make-sba-certs.sh`, `free5gc/config/sba-entities.conf`;
2. `free5gc/crl/`, `free5gc/scripts/build-crl.sh`, `free5gc/patches/run/0003-crl-option.patch`;
3. `free5gc/patches/openapi/0005-sbi-tls.patch`;
4. `free5gc/patches/util/0001-no-tls-key-log.patch`, `free5gc/scripts/prepare-util.sh`;
5. os onze `free5gc/patches/<nf>/…-sbi-tls.patch`, `free5gc/patches/config/0004-sbi-tls.patch`, `free5gc/patches/test/0002-sbi-tls.patch`;
6. `free5gc/aanf/`, `free5gc/config/aanfcfg.yaml`, `free5gc/aanf/tools/seed`;
7. `af/af/internal/{aanf,nef,capif,config}`, `af/config/*.yaml` (sem `afcfg-external.yaml`);
8. `af/scripts/e2e.sh`, `e2e-remote.sh`, `negative/*.py`;
9. documentos: este, `aanf.md`, `nef.md`, `capif.md`, `nrf.md`, `README.md`, `docs/af-side/af.md` e `e2e.md`, o tutorial (parte 12).

**Fase (b).** `free5gc/patches/openapi/0006-access-token.patch`; `free5gc/patches/nrf/0005-token-identity.patch`; `free5gc/patches/config/0005-access-token.patch`; `free5gc/aanf/`; o `e2e.sh`; os documentos.

**Fase (c).** `free5gc/ccf/` e `ccfcfg.yaml`; `free5gc/patches/openapi/0007-capif-tls.patch`, `free5gc/patches/nef/0005-capif-tls.patch` e `free5gc/patches/config/0006-capif-tls.patch`; o AF; o `e2e.sh`; `desvios.md` (DV-1, DV-9, DV-16, DV-17 e DV-18), `correcoes-free5gc.md`, o registro de defeitos, a matriz de rastreabilidade (revisão 11) e o plano.

Na .130, antes da fase (a): *backup* de `config/` e `cert/` e um *snapshot* da VM, anunciados antes. O `/etc/hosts` das duas máquinas muda na fase (a), também anunciado.

## 15. Implementação da fase (a)

Esta seção registra o que a especificação não fixava. Nada aqui muda um ponto aprovado nem a leitura da norma.

### 15.1 O código

- **O par de quem fala com o produtor.** O `AuthorizationCheck` de cada NF continua com a assinatura do free5GC. O roteador chama, depois dele, `oauth.CheckTokenPeer(token, c.Request.TLS)` e, na recusa, `oauth.WriteTokenError`. Assim, os *mocks* dos testes das NFs não mudam. Na fase (a), as duas funções só guardam o lugar: a primeira aceita, a segunda escreve o `401` de hoje. A fase (b) as preenche (seção 7.4).
- **`http.DefaultTransport`.** O `sbitls.Init` também configura esse transporte. O NEF e os clientes do BSF no SMF e no PCF usam o cliente padrão do Go, que sem isso não teria o TLS da SBI.
- **Grupos do NEF.** A conferência da CA por grupo (seção 8.1) já entra na fase (a), com o transporte do CAPIF-2e: um *middleware* novo, `checkClientCA`, em `internal/sbi/client_ca.go`. O `capif.invokerCa` diz a CA dos *invokers*, e o `capif.invokerCrl: false` dispensa a CRL deles até a fase (c), quando os certificados de *invoker* passam a ter ponto de distribuição (seção 8.4).
- **`interface` do CAPIF.** O `capif.interface` do `nefcfg.yaml` aceita `fqdn` ou `ipv4Addr`, um dos dois. O CCF já aceitava: TS 29.222 pede exatamente um dos três endereços.
- **NSSF.** O `nrfId` da lista de NSI passa ao URI do NRF pelo nome.
- **BSF.** Chave nova, `nfInstanceId`, porque o BSF não tinha como fixar o identificador que o certificado carrega.
- **AMF.** A recusa da seção 3.3, `sbi.fqdn` diferente de `amfName`, fica no `NewApp`, junto do `sbitls.Init`.
- **`config/0004`.** O *patch* é gerado com uma linha de contexto, para valer sobre as configurações do AMF, do SMF e da UPF com as edições manuais do N2 e do N3 (`docs/tutorial` parte 9).
- **AAnF.** O `NewApp` recebe o contexto e a configuração, porque o `sbitls.Init` precisa da configuração antes de criar o servidor. O banco de testes acompanha (`test/0002`).

### 15.2 As ferramentas

- **`seed`.** O NF Instance ID e o tipo do perfil temporário no NRF saem do certificado (`sbitls.NfInstanceId` e `NfTypes`), para que a fase (b) não mude nada na ferramenta. O comando novo `token` imprime um *token* do NRF, para o E05.
- **`sbi_tls.py`.** O S01 a S04, o S06 e o S07 estão em `af/scripts/negative/sbi_tls.py`, só com o módulo `ssl` do Python. Roda da máquina do free5GC em todas as portas (passo R1s) e da máquina do AF no AAnF e no NEF. É mais do que a tabela da seção 9.2 pede: o S06 e o S07 valem em todas as portas.
- **S01.** O servidor do Go, numa porta TLS, responde a um pedido HTTP/1.1 com `HTTP/1.0 400 Bad Request` e o texto `Client sent an HTTP request to an HTTPS server`, e fecha. A um h2c, fecha sem responder. Os dois contam como recusa no TLS: nenhum dos dois chega à aplicação.
- **S03.** O `sbitls.New` recusa partir com um certificado próprio de outra CA. Por isso, o S03 é do Python, e não do `seed`.
- **S04.** Além do nome, do URI e do IP, onde o arquivo existe, o certificado apresentado é conferido contra o arquivo instalado. Com o S05, isso vale pelo perfil da seção 4.7.
- **S08.** O Go confere o nome antes da cadeia. O `other-ca` da seção 4.5, com o nome dele, seria recusado pelo nome, e não pela CA. O passo emite na hora, com a CA de `other-ca`, um certificado com o nome do AAnF, e o `seed` só pode recusar pela CA.
- **S09.** O `s_server` mostra o SNI com `-servername` e `-cert2`.
- **S19.** Os `*sslkey.log` que execuções antigas deixaram em `log/` não contam: a conferência olha só os arquivos escritos na execução.

### 15.3 O E07

O AF_ID `af3` do MVP-6 estava na política do NEF e fora do perfil do CCF, de propósito (CP16). Com o *token* do CCF, o NEF o recusa antes de chegar ao AAnF, e o E07 passaria sem provar o que diz. O AF_ID novo `af5.akma.example.FF00000001`, de teste, está na política do NEF e no perfil `af1-external` do CCF, e fora da lista do AAnF. O E07 passa a usá-lo; o `af3` fica para o CP16.

### 15.4 A verificação

- **`e2e.sh`.** `288 passed, 0 failed, 0 skipped`, na execução `20261104-190824` (`docs/af-side/e2e.md` §5.10). O D03, o O04 e o O05 voltaram.
- **Banco de testes do free5GC.** Na .130, com os *patches* da fase: os três `TestAkma` e o `TestRegistration` passam.
- **`go vet` e `go test -race`.** Nenhum aviso e nenhuma corrida, nos módulos do projeto e nas onze NFs. Os testes do AF também passam com o wolfSSL. Duas falhas não vêm do projeto:
  - no SMF, `internal/pfcp/handler`: `TestHandlePfcpAssociationSetupRequest`, `TestHandlePfcpAssociationUpdateRequest` e `TestHandlePfcpManagementRequest` falham do mesmo jeito no SMF sem *patch*;
  - no UDR, `TestUDR_InfluData_CreateThenGet` precisa do MongoDB, que a máquina do AF não tem.
- **Os *patches*.** A árvore original com os *patches* reconstrói as cópias testadas das onze NFs.

### 15.5 Para a fase (c)

- O NEF guarda a CRL da CA do CCF até o `nextUpdate`. O S18 precisa de uma idade máxima menor, senão o certificado de um *invoker* que saiu ainda passa.
- Depois de uma nova entrada, o AF continua com a conexão HTTP/2 antiga, com o certificado anterior. O AF precisa fechar as conexões ociosas ao trocar de certificado.

## 16. Implementação da fase (b)

Como a seção 15: o que a especificação não fixava. Nada aqui muda a leitura da norma.

### 16.1 Dois *patches* a mais

O T28 punha um *patch* por NF só na fase (a). Dois saíram agora, porque o consumidor tem de pedir o escopo que o AAnF passou a exigir:

- `ausf/0003-token-scopes.patch`: o AUSF pede `naanf-akma naanf-akma:anchorkey`;
- `nef/0004-token-scopes.patch`: o NEF pede `naanf-akma`, `naanf-akma:applicationkeyget` e `naanf-akma:applicationkeyget:supi-access`. Como o *patch* já mexe no consumidor do NEF, ele também corrige o F14: o descadastro pede o *token* de `nnrf-nfm` com o tipo `NRF`.

Com isso, o *patch* do NEF da fase (c) passa a `nef/0005-capif-tls.patch`. Os escopos são constantes do `openapi`, em `oauth/scopes.go`. `[3GPP]` TS 29.500 §6.7.3 deixa o consumidor descobri-los no perfil do produtor; o projeto não faz essa descoberta, porque os dois consumidores têm uma operação cada.

### 16.2 A verificação do *token*

- **A ordem nas NFs.** O `AuthorizationCheck` de cada NF manteve a assinatura (seção 15.1): ele faz os passos 1 a 4, o `iss`, o `aud` e o escopo, e o roteador confere o `sub` depois, em `oauth.CheckTokenPeer`. Um *token* com `sub` errado e sem o escopo recebe, nas NFs, o `403` de escopo. O AAnF chama `oauth.Verify` com o estado TLS e segue a ordem exata da seção 7.4.
- **Sem TLS.** Com `scheme: http`, como no banco de testes, `sbitls.Current()` é nulo: não há cadeia nem CRL do certificado de assinatura, nem identidade para o `aud`, nem cliente para o `sub`. Valem a assinatura, as *claims* presentes, o `exp`, o `iss` e o escopo.
- ***Claims* vazias.** Contam como ausentes: `ACCESS_TOKEN_CLAIM_MISSING`, em `401`.
- **O `aud` no NEF.** Um *token* do NRF para outro produtor é recusado no `aud`, com `403`, antes do escopo. Por isso, o E05 passa de `401` para `403` `ACCESS_TOKEN_CLAIM_MISMATCH`. Na fase (c), com `allowNrfTokens: false`, ele volta a `401`, como a seção 9.3 diz.

### 16.3 O AAnF

- **O AF com *token*.** A seção 7.5 fala do AF sem *token*. Como `allowedOperationsPerNfType` dá `applicationkeyget` ao tipo `AF`, um AF com *token* válido segue a mesma regra depois das conferências da seção 7.4: a identidade em `nfInstanceIds` do `afId`, com o tipo `AF`.
- **Sem TLS.** Sem certificado não há identidade: em `http`, a política é só a lista de AF_ID. É o modo do banco de testes.
- **O NEF sem *token*.** Com OAuth2 ligado, a identidade do NEF sem *token* cai na regra do AF e recebe `403`.

### 16.4 O banco de testes

`test/0003-access-token.patch`, que a seção 14 não trazia. O campo `rootcert` saiu do NRF, e o banco o preenchia. O *patch* tira esse campo, aponta o `cert` do NRF e o `nrfCertPem` das NFs para `cert/sba/nrf-token.*`, fixa o NF Instance ID do NRF, o do certificado, e troca a política do AAnF.

### 16.5 As ferramentas

- **`sbi_token.py`.** S10 a S15 estão em `af/scripts/negative/sbi_token.py`, só com a biblioteca padrão do Python 3.6. Roda na .130, no passo R1s, porque o NRF está em *loopback*. Além da tabela da seção 9.2, confere o PLMN do pedido de *token*, o tipo no registro e o tipo na descoberta. O S13 também compara o `x5t#S256` com o arquivo `nrf-token.pem`. O S15 tem um caso de controle: o mesmo pedido com `af1` passa da política e recebe `K_AKMA_NOT_PRESENT`.
- **`seed`.** `register` e `remove` pedem `naanf-akma naanf-akma:anchorkey`.

### 16.6 A verificação

- **`e2e.sh`.** `301 passed, 0 failed, 0 skipped`, na execução `20261105-115322` (`docs/af-side/e2e.md` §5.11). Depois dela, só um comentário do `nrf/0005` mudou, com a citação da cláusula corrigida.
- **Banco de testes do free5GC.** Na .130: os três `TestAkma`; `TestRegistration`, com e sem `oauth`; `TestOAuth2Callback oauth`. O modo `oauth` do banco usa `cert/sba/nrf-token.*` e roda em `http`: exercita o `ES256`, sem cadeia, `aud` nem `sub` (seção 16.2).
- **`go vet` e `go test -race`.** Nenhum aviso e nenhuma corrida, nos módulos do projeto e nas onze NFs. Ficam as duas falhas da seção 15.4, que não vêm do projeto: no SMF, `internal/pfcp/handler`; no UDR, `TestUDR_InfluData_CreateThenGet`.
- **Os *patches*.** A árvore original com os *patches* reconstrói as cópias testadas das onze NFs.

## 17. Implementação da fase (c)

Como as seções 15 e 16: o que a especificação não fixava. Nada aqui muda a leitura da norma.

### 17.1 Os *patches*

- `nef/0005-capif-tls.patch`, e não `nef/0004`, como a seção 16.1 já dizia. A seção 14 foi corrigida.
- `openapi/0007-capif-tls.patch`, que a seção 14 não trazia. Ele traz a idade máxima da CRL (seção 17.3), `sbitls.CrlChecker`, para conferir a CRL de um par cuja CA não é a do processo, e `openapi.CloseIdleConnections` (seção 17.4).
- `config/0006-capif-tls.patch`: o `ccfUri` por nome, `invokerCrl: true`, `invokerCrlMaxAge`, `allowNrfTokens: false` e `authorizedAfs` sem `sourceAddresses`.

### 17.2 O CCF

- **O servidor interno.** Com `certificate`, `key` e `ca`, serve `https`, com `sbitls.Init` e `ConfigureServer`, como as NFs. Sem os três, continua em `http`, para os testes Go. O `nfInstanceId` do provedor é obrigatório com TLS.
- **A CRL depois de reiniciar.** O CCF guarda o estado em memória (`capif.md` §5.12), mas uma revogação não pode se perder. Na partida, ele lê o `ca.crlFile` anterior, confere a assinatura com a sua CA, mantém as entradas e continua a numeração. Uma CRL de outra CA é posta de lado, com um aviso.
- **A assinatura da CRL.** Versão 2, com `cRLNumber` crescente, `nextUpdate` de uma semana e a razão `cessationOfOperation` (RFC 5280 §5.3.1) em cada entrada. O arquivo é trocado por `rename`, e o serviço da CRL nunca lê um arquivo pela metade.
- **CAPIF-1e.** O CCF confere o número de série do certificado de cliente em memória, no `VerifyConnection`, sem baixar a própria CRL. A recusa vale a partir da resposta da saída.
- **`ccf credential`.** O subcomando usa o contexto do CCF e não toca a CRL: a CRL só é assinada por quem serve.

### 17.3 O NEF

- **A idade máxima da CRL.** A seção 5.3 confere a CRL até o `nextUpdate`. Na CRL do CCF, de uma semana, isso deixaria um *invoker* que saiu aceito no NEF por dias. `[MVP]` `capif.invokerCrlMaxAge`, 5 s no laboratório, faz o NEF buscar a CRL de novo depois desse tempo, no *handshake* seguinte. Se a busca falha, a CRL em memória continua valendo até o `nextUpdate`, com um aviso. A CRL do operador não muda de regra. `[3GPP]` TS 33.310 §6.3.1 pede a conferência da CRL e não fixa com que frequência.
- **A ordem no `check-authentication`.** O CN é conferido antes de pedir os dados ao CCF. A CA é conferida depois, porque vem em `authenticationInfo`: o NEF compara o emissor da cadeia verificada com o PEM do CCF. Os dois casos dão `403`, com `detail` diferente.
- **A ordem no `retrieve`.** O CN é conferido depois do `iss` e do `client_id`, e antes do escopo e dos dados do *invoker*.
- **Sem TLS.** Nos testes Go, sem certificado, o `check-authentication` volta à política de endereço do MVP-6, e o `retrieve` não confere o CN.
- **A sessão aberta.** A CRL e o CN são conferidos no *handshake* e em cada pedido, respectivamente. Numa sessão TLS já aberta, um certificado revogado continua até a sessão cair, e o `client_id` do *token* continua a valer até o `exp`. É o limite de TS 33.122 C.6, com a vida curta do *token* (DV-14).

### 17.4 O AF

- **O CCF.** `sbitls.CrlChecker`, no `VerifyConnection` do cliente de CAPIF-1e. O `ccf.pem` não é da CA do processo, que é a do operador, e por isso a conferência não vem do `sbitls.Init`.
- **O certificado de cada sessão.** Com a conferência do CN, uma sessão aberta com o certificado de uma entrada anterior seria de outro *invoker*. Numa nova entrada, o AF fecha as conexões ociosas (`openapi.CloseIdleConnections`), as dos dois clientes do processo.
- **Um *invoker* revogado.** Só o próprio *invoker* pede a saída, e o AF apaga o estado logo depois (`capif.md` §7). Por isso, o AF não trata o caso de um estado com certificado revogado: a partida falharia no TLS, até o estado ser apagado à mão.

### 17.5 As ferramentas

- **S16.** Em `capif_api.py`, com um terceiro *invoker*, que entra só para isso e sai no fim. Além da tabela da seção 9.2, confere o `check-authentication` de um *invoker* com o certificado do outro (`S16.check`). O caso do *token* do NRF é o E05, que `nef_api.py` também registra como `S16.nrf`.
- **E05.** Volta a `401`, como a seção 16.2 previa: `seed token --target NEF --scope 3gpp-akma`.
- **S17.** Em `sbi_token.py`, na .130, porque o servidor interno está em *loopback*. Tem um caso de controle: os dados de um *invoker* desconhecido, com o certificado do NEF, dão `404`, e não `403`.
- **S01 a S07 no CCF.** `ccf-internal` entrou na lista do R1s. O S04 não confere o URI do NF Instance ID nele, porque o servidor interno não é uma NF e o certificado não o tem (seção 4.4).
- **CP18 e S18.** O CP18 confere a recusa no TLS, em vez do `400` `invalid_client`: o certificado revogado não abre mais sessão em CAPIF-1e. O `invalid_client` continua nos testes Go, nos outros casos de recusa. O S18 espera 7 s, o `invokerCrlMaxAge` e uma margem.

### 17.6 A verificação

- **`e2e.sh`.** `317 passed, 0 failed, 0 skipped`, na execução `20261106-170311` (`docs/af-side/e2e.md` §5.12).
- **Banco de testes do free5GC.** Na .130: os três `TestAkma`; `TestRegistration`, com e sem `oauth`; `TestOAuth2Callback oauth`.
- **`go vet` e `go test -race`.** Nenhum aviso e nenhuma corrida, nos módulos do projeto e nas onze NFs. Ficam as duas falhas da seção 15.4, que não vêm do projeto.
- **Os *patches*.** A árvore original com os *patches* reconstrói as cópias testadas das onze NFs e a do `openapi`.
- **Um achado no caminho.** O E05 falhou no primeiro ensaio com `invalid_scope`: o banco de testes deixou perfis no MongoDB, e o NRF confere o escopo em um só perfil do tipo pedido (F27). O passo R1 passou a apagar os perfis antigos (`R1.profiles`).
