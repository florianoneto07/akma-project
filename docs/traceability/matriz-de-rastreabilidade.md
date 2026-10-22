# Matriz de rastreabilidade

**Revisão:** 7 (MVP-8, fase (c): perfis DTLS do Ua\*). Revisão 6 (MVP-8, fases (a) e (b)): serviço de chaves do ME e perfis TLS do Ua\*. Revisão 5 (MVP-6): K_AF via NEF, acesso anônimo, remoção por OAM e UERANSIM em outra máquina. Revisão 4 (MVP-5): etapas mínimas do fluxo ligadas a evidências do roteiro ponta a ponta
**Data:** 02 de setembro de 2026; revisão 3.1 em 16 de outubro de 2026; revisões 3.2 a 3.4 em 17 de outubro de 2026; revisões 3.5 a 3.7 e 4 em 18 de outubro de 2026; revisão 5 em 19 de outubro de 2026; revisões 6 e 7 em 21 de outubro de 2026

*Nota (16/10/2026):* os caminhos da coluna "Implementação" seguem a estrutura anterior ao plano. A estrutura aprovada (plano §3) é: `akma/` → `shared/akma/`; `NFs/aanf/` → `free5gc/aanf/`; `patches/` → `free5gc/patches/`; `af/` → `af/af/`; `ue/` → `af/ue/`. Os caminhos serão atualizados linha a linha à medida que cada MVP os criar. Os desvios citados como DV-n estão em `docs/conformidade/desvios.md`.

Legenda de estado:

- **Especificado** — parâmetros integralmente determinados por cláusula; pronto para implementação
- **Parcial** — parcialmente determinado; ver questão associada
- **Decisão** — determinado por decisão do projeto, registrada em ADR ou no registro de defeitos
- **Pendente** — fora do escopo do MVP ou aguardando decisão
- **Implementado** — código e testes no repositório, com o MVP indicado

## Etapas mínimas do fluxo e evidências (revisão 4)

`[LACUNA]` O plano (MVP-5, critério 3) liga cada linha da "tabela de etapas mínimas" do material de partida a uma evidência. Esse material nunca entrou no repositório (plano §8, I17 e I23). `[INTERPRETAÇÃO]` A tabela abaixo o reconstrói a partir do texto de TS 33.535 V19.0.0, passo a passo: §6.1 (registro após a autenticação primária) e §6.2.1 (K_AF para um AF na rede do operador), com as cláusulas que eles invocam.

Evidências: os identificadores entre colchetes são verificações da execução `20261018-171510` de `af/scripts/e2e.sh` (`docs/af-side/e2e.md`): R na máquina do free5GC; S, N, B, C e L na do AF. Os demais são testes do repositório.

| Etapa | Cláusula | Implementação | Evidência | Estado |
|---|---|---|---|---|
| AUSF obtém do UDM os dados de autenticação (`Nudm_UEAuthentication_Get`) | 6.1 passo 1 | free5GC (sem alteração) | [R1.ue], [R2.registration] | **Verificado** |
| UDM inclui a indicação AKMA e o RID | 6.1 passo 2; TS 33.501 6.1.3.x passo 2; TS 29.503 Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` | log `[UDM][AKMA] AKMA indication sent … Routing Indicator 0000` em [R2]; ausente sem assinatura [R5.no-indication] | **Verificado** |
| AUSF gera K_AKMA e A-KID após a autenticação bem-sucedida, com o RID do UDM | 6.1 passo 3; A.2, A.3 | `free5gc/patches/ausf/` | [R2.akid-ausf]; `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` | **Verificado** |
| UE gera K_AKMA e A-KID antes de comunicar com um AF | 6.1 passo 3; 6.2.1 passo 1 | `af/ue/ueransim/patches/0001` | A-KID do UE igual ao do AUSF e do AAnF [R2.akid-ausf], [R2.akid-aanf]; [S1] | **Verificado** |
| AUSF seleciona o AAnF | 6.1 passo 4; 6.7 | `free5gc/patches/ausf/`, `free5gc/patches/nrf/0002` | log `AAnF for Routing Indicator 0000 … (discovered at the NRF)` em [R2]; sem AAnF registrado, `no uri for naanf-akma` [R4.ausf-error] | **Verificado** |
| AUSF envia A-KID, K_AKMA e SUPI (`Naanf_AKMA_AnchorKey_Register`); AAnF guarda o mais recente | 6.1 passo 4 e NOTA 1a | `free5gc/patches/ausf/`, `free5gc/aanf/` | [R2.akid-aanf]; substituição [R3.replaced]; A-KID anterior recusado [B1] | **Verificado** |
| AAnF responde ao AUSF | 6.1 passo 5 | `free5gc/aanf/` | log `AKMA key material … registered at the AAnF` em [R2] | **Verificado** |
| Autenticação primária não depende do AAnF | 6.1 passo 3 (D9) | `free5gc/patches/ausf/` (`startAkma`) | [R4.registration], [R4.ausf-error] | **Verificado** — D9 |
| UE inclui o A-KID no pedido de sessão; K_AF derivada pelo ME | 6.2.1 passo 1; A.4 | `0003-akma-kaf-nr-cli-selftest.patch`; `af/ue/uaclient/` | [R2.ua], [R3.ua-new], [R5.ua-restored] | **Verificado** |
| AF sem contexto ativo pede a K_AF ao AAnF com A-KID e AF_ID | 6.2.1 passo 2; 6.7 (configuração local) | `af/af/internal/aanf/` | log do AF `K_AF … obtained from the AAnF`; `TestApplicationKey` | **Verificado** |
| AAnF verifica a política local para o AF | 6.2.1 passo 2 | `free5gc/aanf/` (`AfPolicy`) | [N03] (403); o AF fora da política recebe 503 (MVP-3) | **Verificado** — DV-1 |
| AAnF verifica a presença da K_AKMA | 6.2.1 passo 2 | `free5gc/aanf/` | [N06], [B1], [C1], [N21] (403 `K_AKMA_NOT_PRESENT`) | **Verificado** — DV-7 |
| GPSI pelo UDM | 6.2.1 passos 3 e 4 | — | — | Não aplicável: a política local do MVP não pede GPSI |
| Subscrição de `ROAMING_STATUS` | 6.2.1 passos 5 e 6 | — | — | **Desvio DV-2**: etapa desabilitada |
| AAnF deriva a K_AF se ainda não a tem e guarda o `expiry` | 6.2.1 passo 7; A.4; 5.2 | `free5gc/aanf/` | [N01b]; a mesma `expiry` em pedidos seguintes (log do AAnF em [R2] e [N01]) | **Verificado** |
| AAnF responde com SUPI, K_AF e `expiry` | 6.2.1 passo 8 | `free5gc/aanf/` | [N01b] | **Verificado** |
| AF responde ao UE; com falha, recusa com causa | 6.2.1 passo 9 | `af/af/internal/server/` | 200 [R2.ua]; `403 AKMA_KEY_NOT_AVAILABLE` [R5.ua]; `503 AANF_UNAVAILABLE` [R4.ua] | **Verificado** |
| Iniciação do AKMA pelo AF | 6.5 | `af/af/internal/server/` | `AKMA initiation: challenge sent` no log do AF em cada Ua\*; `TestInitiation` | **Verificado** |
| K_AF vencida no AF; K_AF anterior válida até o `expiry` | 6.4.2; 5.2; 4.4.0 | `af/af/internal/server/keys.go` | [R3.ua-old]; [N22], [N23]; `TestKafExpiry*` | **Verificado** |
| Remoção do contexto AKMA | 6.6 | `free5gc/aanf/` | [R6.removed], [R6.not-found], [C1] | **Verificado** |
| Aplicações do UE sem acesso à K_AKMA | 4.4.3 | `0003-…` (`akma-kaf`) | `docs/af-side/teste-ue.md` §4 | **Verificado** |
| Nenhuma chave em log | plano, todos os MVPs | todos | [R7.no-key-*], [L1]; `TestNoKeyInLogs` | **Verificado** |

## Procedimentos do MVP-6 e evidências (revisão 5)

Evidências: verificações da execução `20261019-213720` de `af/scripts/e2e.sh` (`docs/af-side/e2e.md` §4.1 e §5.2). E e N01c são da máquina do AF; R2n, da do free5GC. "3H" é o teste com o UERANSIM em outra máquina (`docs/af-side/teste-nef.md`).

| Etapa | Cláusula | Implementação | Evidência | Estado |
|---|---|---|---|---|
| AF de fora da rede pede a K_AF ao NEF com A-KID e AF_ID; sem CAPIF, ponto de terminação configurado | 33.535 6.3 passo 1 e NOTA; 29.522 4.4.23.2, 5.14 | `af/af/internal/nef`; `af/config/afcfg-external.yaml` | [E01]; [R2n.ua-nef]; 3H | **Verificado** |
| NEF autoriza o AF, inclusive pelo AF_ID | 33.535 6.3 passo 2 | `free5gc/patches/nef/0001` (`AkmaAfPolicy`) | [E06]; [E04], [E05] (*token*) | **Verificado** — DV-9 |
| NEF seleciona o AAnF pelo RID | 33.535 6.3 passo 2; 6.7 | `nef/0001` (`aanf_service.go`) | log do NEF e [R2n.aanf-via-nef]; teste `no AAnF at the NRF` | **Verificado** |
| NEF pede a K_AF ao AAnF; o AAnF processa como em 6.2.1 | 33.535 6.3 passo 3 | `nef/0001`; `free5gc/aanf/` (`NefPolicy`) | [R2n.aanf-via-nef]; [E07] (fora da política do AAnF); [E08] (sem K_AKMA) | **Verificado** — P11 |
| AAnF responde ao NEF com K_AF, validade e SUPI | 33.535 6.3 passo 4 | `free5gc/aanf/` | `TestRetrieveViaNef`; [E02] (a mesma K_AF entregue direto) | **Verificado** |
| NEF responde ao AF sem o SUPI; GPSI opcional não enviado | 33.535 6.3 passo 5; 29.522 tab. 5.14.5.3.3-1 | `nef/0001` (`processor/akma.go`) | [E01]; [R2n.nef-no-supi] | **Verificado** — GPSI fora do escopo (F10) |
| NEF repassa `403 K_AKMA_NOT_PRESENT` | 29.522 4.4.23.2 | `nef/0001` | [E08] | **Verificado** |
| 204 quando o dado não existe | 29.522 tab. 5.14.3.2.2-2 | `nef/0001`: repassa o 204 do AAnF | teste `204 of the AAnF relayed` | **Verificado** — `[INTERPRETAÇÃO]`; o AAnF nunca emite 204 (DV-7) |
| Acesso anônimo: K_AF sem identidade do UE; o AAnF guarda a validade | 33.535 6.2.2; 7.1.5; 29.535 4.2.2.3.2 | `free5gc/aanf/` (MVP-1); AF `anonymous`; NEF repassa `anonInd` | [N01c], [E03], [R2n.ua-anon], [R2n.ua-nef-anon]; 3H | **Verificado** |
| Remoção do contexto AKMA pelo OAM | 33.535 6.6, tab. 7.1.1-1; 29.535 4.2.1.3.2, tab. 4.2.2.1-1 | `seed remove` no papel do OAM | [R6.removed], [R6.not-found] | **Verificado** |
| Segurança NEF–AF: TLS mútuo e OAuth | 33.501 12.2 a 12.4 | `http`, *token* do NRF, política AF_ID ↔ endereço | [E04], [E06] | **Desvio DV-9**; CAPIF no MVP-9 |
| UE em outra máquina: N2 e N3 entre máquinas | 23.501 (N2, N3); guia oficial do free5GC | configuração do AMF, do SMF, do UPF e dos gNB | 3H: NG Setup a partir de 192.168.16.131, mesmo A-KID no UE, no AUSF e no AAnF | **Verificado** |
| Nenhuma chave em log, inclusive no NEF | plano, todos os MVPs | todos | [R7.no-key-*], [L1]; `TestAkmaNoKeyInLogs` do NEF; 3H | **Verificado** |

## Procedimentos do MVP-8 e evidências (revisão 6)

Evidências: verificações da execução `20261021-164401` de `af/scripts/e2e.sh`, com 115 verificações e nenhuma falha (`docs/af-side/e2e.md` §5.3 e §5.4). R2k, R2t e R2n são da máquina do free5GC; T e L1, da do AF. Especificação: `docs/af-side/ua-star-norma.md`.

| Etapa | Cláusula | Implementação | Evidência | Estado |
|---|---|---|---|---|
| As aplicações do UE não acessam a K_AKMA | 33.535 4.4.3 item 1 | `af/ue/ueransim/patches/0005` (serviço de chaves do ME) | respostas só com A-KID, AF_ID e K_AF; [R7] | **Verificado** |
| Cada aplicação só recebe a K_AF dos AF_IDs autorizados | 33.535 4.4.3 item 2 | `0005`: política `akmaKeyService`, *uid* pelo `SO_PEERCRED` | [R2k.given], [R2k.af-id], [R2k.uid], [R2t.me] | **Verificado** — `[INTERPRETAÇÃO]` *uid* como fronteira (U3) |
| Uma aplicação não recebe a K_AF de outra | 33.535 4.4.3 item 3 | `0005`: AF_ID em duas aplicações recusado | [R2k.duplicate] | **Verificado** (U5) |
| `nr-cli akma-kaf` só para teste | 33.535 4.4.3 | `0005` (`akmaCliKaf`) | [R2k.nr-cli] | **Verificado** |
| K.2: TLS com certificado do AF conferido contra o FQDN; Digest dentro do TLS | 24.109 K.2 a) a d); 33.535 B.1.2.2; 33.222 5.3.0, 5.3.1.3 | `af/akmatls`, `af/af/internal/{profile,server}`, `af/ue/uaclient` | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09]; `TestK2RealmOfAnotherAF` | **Verificado** — DV-5, DV-6 |
| Identificador Ua\* `01 00 01 yy zz` do *cipher suite* (K.2 e TLS 1.2) e `01 00 00 00 07` (TLS 1.3) | 33.220 H.3; 33.222 5.4.0.2 NOTA 0 | `af/akmatls` (`UaStarTLS`); AF e UE | AF_IDs no AF, no AAnF e no ME [T01] a [T04], [R2t.me] | **Verificado** — ADR-0003 revisto |
| PSK TLS 1.3: `psk_dhe_ke`, `server_name`, identidade `3GPP-AKMA;<A-KID>`, PSK da K_AF | 24.109 K.3.3; 33.222 5.4.0.2; 33.210 6.2.2 | `af/wolftls` (GPLv3), perfil `psk-tls13` | [R2t.psk13], [T03], [T05] (OpenSSL), [T08] (`psk_ke` recusado); testes do `ClientHello` | **Verificado** — D-E; U1 `[LACUNA]` |
| PSK TLS 1.2: suíte `D0 01`, *hint* `3GPP-AKMA`, identidade com `;`, *premaster* da RFC 5489 com a K_AF | 24.109 K.3.2; 33.222 5.4.0.1; 33.210 6.2.3 | `af/wolftls`, perfil `psk-tls12` | [R2t.psk12], [T04], [T06]; testes do `ClientHello` e do `ServerHello` | **Verificado** — U1, U2 `[LACUNA]`; sem `D0 02` (DV-10) |
| AF recusa a identidade sem K_AF; *handshake* falha | 33.222 5.4.0.1 passo 4; 33.535 6.2.1 passo 2 | `af/af/internal/server` (`PSKKey`) | [T07], [T07b], [T10], [T10b] | **Verificado** |
| Perfil TLS de TS 33.210: TLS 1.2 e 1.3; AEAD com ECDHE; EMS, SNI, RFC 5746; grupos | 33.210 6.2; 33.310 Anexo E | `af/akmatls`; `af/scripts/build-wolfssl.sh` (macros); `af/wolftls` | testes de `af/wolftls` e `af/akmatls`; `build-wolfssl.sh` | **Verificado** — sem OCSP *stapling* (DV-10) |
| Retomada de sessão: PSK até a expiração da K_AF e 24 h; K.2 até 24 h | 33.222 5.3.1.5, 5.4.1.4 | *tickets* do `af/wolftls`; `boundTickets` do `af/akmatls` | [R2t.k2-resumed], [R2t.psk12-resumed], [T03], [T04]; `TestResumption`, `TestExpiry`, `TestTicketLifetime` | **Verificado** (U8) |
| K_AF vencida numa sessão TLS: sessão encerrada; `refetch` ou `reject` | 33.535 6.4.2, 4.4.0 | `af/wolftls` (`Expired`); `keys.go` | [T11], [T12], [T12b], [T13], [T13b]; `TestPSKExpiry` | **Verificado** (U9) |
| AF externo com perfil PSK, K_AF pelo NEF | 33.535 6.3; 24.109 K.3.3 | AF com `nefUri` e perfil `psk-tls13` | [R2n.ua-nef-psk], [R2n.nef-psk-relayed] | **Verificado** |
| Políticas do AAnF, do NEF e do UE com um AF_ID por perfil e *cipher suite* | 33.535 6.2.1 passo 2, 6.3 passo 2, 4.4.3 | `free5gc/config/aanfcfg.yaml`, `free5gc/patches/config/0002`, `0005` | [T10] (fora da política); demais [R2t] e [T] | **Verificado** (U7) |
| Nenhuma chave em log | plano, todos os MVPs | todos; wolfSSL sem *debug* | [R7.no-key-*], [L1] | **Verificado** |

## Perfis DTLS do MVP-8 e evidências (revisão 7)

Evidências: verificações R2d (máquina do free5GC) e D (máquina do AF) da execução `20261022-124514` de `af/scripts/e2e.sh`, com 126 verificações e nenhuma falha (`docs/af-side/e2e.md` §5.5), e testes de `af/wolftls`, `af/coap` e `af/af/internal/profile`. Especificação: `docs/af-side/ua-star-norma.md` §4.

| Etapa | Cláusula | Implementação | Evidência | Estado |
|---|---|---|---|---|
| UE como cliente DTLS; identidade PSK de B.1 levada em mensagem DTLS | 33.535 C.1.1, C.2.2 | `af/wolftls` (DTLS), `af/ue/uaclient` (`psk-dtls13`, `psk-dtls12`) | [R2d.dtls13], [R2d.dtls12]; `TestDTLSCookieOnTheWire` (identidade no `ClientHello`) | **Verificado** |
| DTLS 1.3 com o procedimento do TLS 1.3 (B.1.3.2.2): `psk_dhe_ke`, PSK externa da K_AF | 33.535 C.2.2; 33.210 6.2.2 | `af/wolftls`; perfil `psk-dtls13` | [R2d.dtls13], [D01], [R2d.psk-ke]; `TestDTLSRefused` | **Verificado** — U1 |
| Identificador Ua\* do DTLS 1.3 `01 00 00 00 09` | 33.220 H.3; 33.535 A.4 | `af/akmatls` (`UaStarDTLS13`); políticas | [D01], [R2d.me] | **Verificado** — `[INTERPRETAÇÃO]` U18 |
| DTLS 1.2, que TS 33.210 exige de quem tem DTLS | 33.210 6.2.1; 33.535 C.1.1 | perfil `psk-dtls12`, procedimento e identificador do TLS 1.2 | [R2d.dtls12], [R2d.dtls12-resumed], [D02] | **Verificado** — `[LACUNA]` U13 (b) |
| Perfil TLS de TS 33.210 aplicado ao DTLS | 33.535 C.2.1; 33.210 6.2 | as mesmas listas de suítes e grupos do TLS | testes de `af/wolftls` | **Verificado** — DV-11 (RFC 7252 §9.1.3.1) |
| *Cookie* contra negação de serviço; nenhuma consulta ao AAnF antes dele | RFC 9147 5.1; RFC 6347 4.2.1 | `af/wolftls` (`wolftls_dtls_peer`, *callback* de *cookie*, fase sem estado) | `TestDTLSCookieOnTheWire`, `TestDTLSNoKeyBeforeCookie` | **Verificado** (U16) |
| CoAP dentro do DTLS: `GET` da sessão, `2.05` em JSON | RFC 7252 3, 4.2, 5.2.1 | `af/coap`; `server.CoAP` | [R2d], [D03]; `TestCoAPHandler`, testes de `af/coap` | **Verificado** — `[LACUNA]` U14 |
| Conferência independente do DTLS 1.3 | — | OpenSSL 4.1 *beta* 2, fora do repositório | [D03]; `TestOpenSSLDTLS13Client` | **Verificado** (U21) |
| Retomada; K_AF vencida | 33.222 5.4.1.4; 33.535 6.4.2 | *tickets* do `af/wolftls`; `Expired` | [R2d.dtls12-resumed], [D01]; `TestDTLSResumption`, `TestDTLSExpiry`, `TestPSKExpiry` | **Verificado** (U19) |

## Derivação criptográfica

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| KDF genérica HMAC-SHA-256 | TS 33.220 V19.1.0 | B.1, B.2.0, B.2.1 | `shared/akma/kdf/kdf.go` | `kdf_test.go`, `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Codificação de cadeias UTF-8 + NFKC | TS 33.220 V19.1.0 | B.2.1.2 | `shared/akma/kdf/kdf.go` (`String`) | `kdf_test.go`; vetores `nfkc-*` | **Implementado** (MVP-0) |
| Codificação de comprimentos em 2 octetos | TS 33.220 V19.1.0 | B.2.1.3 | `shared/akma/kdf/kdf.go` | `kdf_test.go` | **Implementado** (MVP-0) |
| Derivação de K_AKMA — FC `0x80` | TS 33.535 V19.0.0 | A.2 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Derivação de A-TID — FC `0x81` | TS 33.535 V19.0.0 | A.3 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| Derivação de K_AF — FC `0x82` | TS 33.535 V19.0.0 | A.4 | `shared/akma/derive/derive.go` | `vectors_test.go`, `shared/crosscheck` | **Implementado** (MVP-0) |
| SUPI como P0/P1, sem prefixo `imsi-` | TS 33.501 V19.7.0 | A.7.0 | `shared/akma/supi/supi.go` | `supi_test.go`, `derive_test.go`; contraexemplo nos vetores | **Implementado** (MVP-0); só IMSI (DV-3) |

## Identificadores

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Identificador Ua* — 5 octetos | TS 33.220 V19.1.0 | H.1–H.3 | `shared/akma/afid/uastar.go` | `afid_test.go`, `derive_test.go` | **Implementado** (MVP-0) |
| `AfId` de transporte — `<FQDN>.<10 hex>` | TS 29.522 V19.7.0 | Tab. 5.14.5.4.2-1 | `shared/akma/afid/afid.go` | `afid_test.go` | **Implementado** (MVP-0) — decisão P9: maiúsculas na escrita, as duas caixas na leitura |
| Conversão `AfId` → P0 da KDF | TS 33.535 A.4 + TS 29.522 | — | `shared/akma/afid/afid.go` | `afid_test.go`; contraexemplo nos vetores | **Implementado** (MVP-0) — decisão ADR-0002 |
| Home Network Domain | TS 23.003 V19.7.0 | 2.2, 28.2 | `shared/akma/akid/realm.go` | `akid_test.go` | **Implementado** (MVP-0) |
| Routing Indicator, cadeia, valor padrão `"0"` | TS 23.003 V19.7.0 | 2.2B | `shared/akma/akid/rid.go` | `akid_test.go` — casos `"0"`, `"012"`, `"12"`, `"0000"` | **Implementado** (MVP-0) |
| A-KID — formato NAI `rid<v>.atid<v>@<realm>` | TS 33.535 V19.0.0; TS 29.522 V19.7.0 | 6.1; Tab. 5.14.5.4.2-1 | `shared/akma/akid/akid.go`, `shared/akma/derive/derive.go` | `akid_test.go`, `vectors_test.go` | **Implementado** (MVP-0) — defeito D-A: `atid` na escrita, `atid`/`akid` na leitura |
| Codificação do valor do A-TID | — | — | `shared/akma/akid/akid.go` | `akid_test.go`, `vectors_test.go` | **Implementado** (MVP-0) — decisão ADR-0002 (base64; base64url e hexadecimal configuráveis) |
| Codificação de `kAkma` e `kaf` | TS 29.503 (`Kausf`) e TS 29.509 (precedentes) | — | `shared/akma/codec/key.go` | `key_test.go` | **Implementado** (MVP-0) — decisão ADR-0002 |

## API Naanf_AKMA

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `POST /naanf-akma/v1/register-anchorkey` | TS 29.535 V19.5.0 | 4.2.2.2.2, 5.1.4.2 | `free5gc/aanf/internal/sbi/processor/register_anchorkey.go` | `TestRegisterAnchorKey`, `TestRegisterBadRequests`; VM | **Implementado** (MVP-1) |
| `POST /naanf-akma/v1/retrieve-applicationkey` | TS 29.535 + TS 29.522 | 4.2.2.3.2, 5.14.5.3 | `free5gc/aanf/internal/sbi/processor/retrieve_applicationkey.go` | `TestRetrieveApplicationKey`, `TestRetrieveAnonymous`; VM (K_AF = vetor) | **Implementado** (MVP-1); `notifUri` ignorado (sem `RoamingRestriction`) |
| `POST /naanf-akma/v1/remove-context` | TS 29.535 V19.5.0 | 4.2.2.4.2, 5.1.4.4 | `free5gc/aanf/internal/sbi/processor/remove_context.go` | `TestRemoveContext`; VM | **Implementado** (MVP-1) |
| Erro `K_AKMA_NOT_PRESENT` (403) | TS 29.535 V19.5.0 | 5.1.7.3 | `free5gc/aanf/internal/sbi/processor/` | `TestRetrieveKAkmaNotPresent`; VM | **Implementado** (MVP-1) |
| Erro `AKMA_CONTEXT_NOT_FOUND` (404) | TS 29.535 V19.5.0 | 5.1.7.3 | `free5gc/aanf/internal/sbi/processor/` | `TestRemoveContext`; VM | **Implementado** (MVP-1) |
| `retrieve-applicationkey` nunca responde 204 | TS 29.535 V19.5.0 | 4.2.2.3.2 | `free5gc/aanf/internal/sbi/processor/` | conferido em todo teste de API | **Implementado** (MVP-1) — DV-7, defeito D-G |
| AF_ID autorizado associado ao AF que o pede | TS 33.535 V19.0.0 | 4.2.2, 6.2.1 passo 2 | `free5gc/aanf/internal/context/context.go` (`AfPolicy`) | `TestRetrieveAfNotAuthorized`; VM | **Implementado** (MVP-1) — risco residual sem TLS (DV-1) |
| Códigos HTTP obrigatórios para POST | TS 29.500 V19.7.0 | 5.2.7.1, 5.2.7.2 | `free5gc/aanf/internal/sbi/api_akma.go` | 8 códigos testados; 429, 502, 503 justificados (`docs/free5gc-side/aanf.md` §6) | **Implementado** (MVP-1) |
| Método não suportado → 501; versão de API desconhecida → 400 `INVALID_API` | TS 29.500 V19.7.0 | 5.2.7.2, Tab. 5.2.7.2-1 | `free5gc/aanf/internal/sbi/api_akma.go` | `TestRoutingErrors` | **Implementado** (MVP-1) |
| OAuth2 em `register-anchorkey` e `remove-context`: assinatura, escopo, `exp`, `aud` | TS 29.535 V19.5.0; TS 33.501 V19.7.0 | 5.1.9; 13.4.1.1.2 | `free5gc/aanf/internal/context/context.go` | `TestOAuth2`; VM (401 sem token; 200 com token do NRF corrigido) | **Implementado** (MVP-1; `aud` estrito desde o MVP-1b, DV-8 encerrado) |
| Negociação de *features* (nenhuma suportada) | TS 29.535 V19.5.0; TS 29.500 | 5.1.8; 6.6 | `free5gc/aanf/internal/sbi/processor/processor.go` | `TestRegisterFeatureNegotiation` | **Implementado** (MVP-1) |
| `oneOf` supi/gpsi em `AkmaKeyInfo` | `TS29535_Naanf_AKMA.yaml` v1.2.0; TS 29.535 | Anexo A.2; Tab. 5.1.6.2.2-1 | `free5gc/aanf/internal/sbi/processor/register_anchorkey.go` | `TestRegisterBadRequests` | **Implementado** (MVP-1) — sem `AKMA_GPSI_Support`, só `supi` |
| Modelos e cliente do Naanf_AKMA no `free5gc/openapi` | TS 29.535 Anexo A.2; TS 29.522 | 5.14.5.3 | `free5gc/patches/openapi/` | `TestOpenapiClient` | **Implementado** (MVP-1) — D6 |

## Contexto do AAnF

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Contexto AKMA por SUPI | TS 33.535 V19.0.0 | 6.6.1, 6.2.1 | `free5gc/aanf/internal/context/akma_context.go` | `akma_context_test.go` | **Implementado** (MVP-1) — em memória |
| Substituição de (A-KID, K_AKMA) em reautenticação | TS 33.535 V19.0.0 | 6.1 NOTA 1a | `free5gc/aanf/internal/context/akma_context.go` | `TestReRegistrationReplacesPair`, `TestReRegistrationOverAPI` | **Implementado** (MVP-1) |
| K_AF reaproveitada até expirar; K_AF emitidas mantidas após novo K_AKMA | TS 33.535 V19.0.0 | 6.2.1 passo 7; 5.2 | `free5gc/aanf/internal/context/akma_context.go` | `TestKafReusedUntilExpiry` | **Implementado** (MVP-1) |
| A-KID transmitido por AF | TS 33.535 §6.8.1 NOTA 1, §6.2.1 passo 8 (exigido com subscrição de *roaming*); TS 29.522 `ServiceDisableNotif` | — | `free5gc/aanf/internal/context/akma_context.go` (`AfKeyRecord.AKIDSent`) | unitário | **Implementado** (MVP-1) — guardado; notificação fora do MVP |
| Estado de *roaming* por tipo de acesso | TS 33.535 §6.2.1 passo 7; `RoamingStatusReport` | — | `free5gc/aanf/internal/context/` | unitário | Pendente — DV-2 |
| `expiry` como instante absoluto | TS 29.522 + TS 29.122 | — | `free5gc/aanf/internal/context/akma_context.go` | `TestRetrieveApplicationKey` (RFC 3339, UTC); desvio de relógio no AF fica para o MVP-3 | **Implementado** no AAnF (MVP-1) |

## Registro e descoberta

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Registro do AAnF no NRF | TS 29.510 V19.7.0 | — | `free5gc/aanf/internal/sbi/consumer/nrf_service.go` | `TestRegisterNFInstance*`; VM | **Implementado** (MVP-1) |
| `nfType = "AANF"`, `serviceName = "naanf-akma"` | TS 29.510 V19.7.0 | — | `free5gc/aanf/internal/context/context.go` | `TestRegisterNFInstanceProfile`; VM (MongoDB) | **Implementado** (MVP-1) |
| `AanfInfo.routingIndicators` no perfil | TS 29.510 V19.7.0 | Tab. 6.1.6.2.73-1 | `free5gc/aanf/internal/sbi/consumer/nrf_service.go` | `TestRegisterNFInstanceProfile`; VM (MongoDB) | **Implementado** (MVP-1 no AAnF; guardado pelo NRF desde o MVP-1b) |
| NRF guarda `aanfInfoList` | TS 29.510 V19.7.0 | 5.2.2.2.2; Tab. 6.1.6.2.2-1 | `free5gc/patches/nrf/0001-store-aanf-info-list.patch` | `TestNnrfNFManagementDataModelKeepsAanfInfoList`; VM | **Implementado** (MVP-1b) |
| NRF descobre `AANF` por `routing-indicator` | TS 29.510 V19.7.0 | Tab. 6.2.3.2.3.1-1, Tab. 6.1.6.2.73-1 | `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` | `TestValidateQueryParametersNfTypes`, `TestBuildFilterAanfRoutingIndicator`; VM (MongoDB 5.0, quatro perfis sintéticos) | **Implementado** (MVP-1b); `complexQuery` fora do escopo |
| NRF preenche `aud` do token | TS 29.510 V19.7.0; TS 33.501 V19.7.0 | Tab. 6.3.5.2.4-1; 13.4.1.1.2 | `free5gc/patches/nrf/0003-access-token-audience.patch` | `TestAccessTokenAudience`; VM | **Implementado** (MVP-1b) |
| Descoberta do AAnF pelo AUSF, por `routing-indicator` | TS 29.510 V19.7.0; TS 33.535 §6.7 | — | `free5gc/patches/ausf/` (`aanfUri`) | `TestRegisterAkmaDiscoversAanf`; UERANSIM (`routing-indicator=0000`, 200) | **Implementado** (MVP-2) — `akma.aanfUri` como alternativa |
| Endereço do AAnF alcançável a partir da VM do AF | ADR-0001-R2 §5.2 | — | `free5gc/config/aanfcfg.yaml` | VM, 17/10/2026 | **Implementado** (MVP-1) |

## Integração AUSF e UDM

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| `akmaInd` e `routingId` na resposta do UDM | TS 33.501 V19.7.0; TS 29.503 V19.7.0 | 6.1.3.1 passo 2, 6.1.3.2.0 passo 2; Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` (`internal/sbi/processor/akma.go`) | `TestSetAkmaIndication`, `TestGenerateAuthDataAkma`, `TestGenerateAuthDataWithoutAkma`; banco de testes | **Implementado** (MVP-2) |
| Fonte do RID no UDM: `routingId` do `AuthenticationSubscription` no UDR; comparação com o RID do SUCI, quando houver, com divergência em log | TS 29.505 V19.6.0; TS 29.503 V19.7.0; TS 33.535 §6.1 | Tab. 5.4.2.2-1; Tab. 6.3.6.2.3-1 | `free5gc/patches/udm/` (`suci.RoutingIndicator`, `setAkmaIndication`) | `TestRoutingIndicator`, `TestSetAkmaIndication`; UERANSIM | **Implementado** (MVP-2) — D10 |
| Assinatura AKMA no perfil do assinante (`akmaAllowed`, `routingId`) | TS 33.535 V19.0.0; TS 29.505 V19.6.0 | 3.1, 4.2.5; Tab. 5.4.2.2-1 | `free5gc/scripts/provision-akma-subscriber.sh` (a *webconsole* não tem os campos) | VM (MongoDB) | **Implementado** (MVP-2) |
| Acionamento do AKMA no AUSF — 5G-AKA | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | 6.1 passo 3; 6.1.3.2.0 passo 11 | `free5gc/patches/ausf/` — `Auth5gAkaComfirmRequestProcedure` | `TestAkma5GAKA`; UERANSIM | **Implementado** (MVP-2) |
| Acionamento do AKMA no AUSF — EAP-AKA' | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | 6.1 passo 3; 6.1.3.1 passo 10 | `free5gc/patches/ausf/` — `EapAuthComfirmRequestProcedure` | `TestAkmaEAPAKAPrime` | **Implementado** (MVP-2) |
| Domínio de origem do A-KID: MCC e MNC do SUCI; com SUPI, da PLMN única de `plmnSupportList` | TS 33.535 V19.0.0; TS 23.003 V19.7.0 | 6.1; 2.2, 28.2 | `free5gc/patches/ausf/` (`homeNetwork`) | `TestHomeNetworkFromSuci`, `TestBuildAkmaKeyInfoMatchesOracle` | **Implementado** (MVP-2) — `[INTERPRETAÇÃO]` |
| Registro no AAnF sem bloquear a resposta ao AMF; falha no AAnF não altera a autenticação | TS 33.535 V19.0.0 | 6.1 passo 3 | `free5gc/patches/ausf/` (`startAkma`) | `TestStartAkmaDoesNotBlock`; UERANSIM com AAnF parado | **Implementado** (MVP-2) — D9 |
| Consumidor `Naanf_AKMA` no AUSF | TS 29.535 V19.5.0 | 4.2.2.2 | `free5gc/patches/ausf/` (`internal/sbi/consumer/aanf_service.go`) | `TestRegisterAkmaDiscoversAanf`, `TestRegisterAkmaConfiguredUri`; banco de testes; UERANSIM | **Implementado** (MVP-2) |
| GPSI via `GET /{ueId}/id-translation-result` | TS 29.503 V19.7.0 | 5.2.2.2.10 | — | — | Pendente — política local |
| Subscrição `ROAMING_STATUS` com `immediateFlag` | TS 29.503 V19.7.0 | 5.5.2.2.2, 6.4.6.2.8 | — | — | Decisão — DV-2 (configurável, desabilitada) |
| `POST /{ueId}/registrations/trigger-auth` | TS 33.501 §14.2.6; TS 29.503 §5.3.2.16 | — | — | — | Pendente — defeito D-B |

## Lado UE (ME)

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| ME deriva K_AKMA e A-KID da K_AUSF antes de comunicar com um AF | TS 33.535 V19.0.0 | 6.1; 6.2.1 passo 1 | `af/ue/ueransim/patches/0001-akma-key-derivation.patch` (`NasMm::deriveAkmaKeys`) | VM: A-KID do UE igual ao do AUSF e do AAnF | **Implementado** (MVP-3a, D14); instante da derivação `[INTERPRETAÇÃO]` (`docs/af-side/ue-ueransim.md` §2) |
| K_AKMA (FC `0x80`) e A-TID (FC `0x81`) no ME; SUPI = IMSI | TS 33.535 V19.0.0; TS 33.501 V19.7.0 | A.2, A.3; A.7.0 | `af/ue/ueransim/patches/0001-akma-key-derivation.patch` (`keys::CalculateKAkma`, `keys::CalculateATid`) | `nr-akma-selftest` (8 vetores do MVP-0); VM (A-KID igual ao da rede) | **Implementado** (MVP-3a; conferido no MVP-4) |
| A-KID no ME: RID do SUCI, A-TID em base64, domínio com MNC de três dígitos | TS 33.535 V19.0.0; TS 29.522 V19.7.0; TS 23.003 V19.7.0 | 6.1; Tab. 5.14.5.4.2-1; 2.2B, 28.2 | idem (`keys::ConstructAKid`) | VM, RID `"0000"` | **Implementado** (MVP-3a) |
| Opção `akma` na configuração do UE | — | — | idem; `0002-free5gc-ue-akma.patch` | VM: sem a opção, nenhuma derivação | **Implementado** (MVP-3a) — decisão do plano |
| K_AF no ME, sob demanda, para um AF_ID | TS 33.535 V19.0.0; TS 33.220 V19.1.0 | 6.2.1 passo 1; A.4; Anexo H | `0003-akma-kaf-nr-cli-selftest.patch` (`keys::CalculateKAf`) | `nr-akma-selftest`; Ua\* com a K_AF do UE (`docs/af-side/teste-ue.md`) | **Implementado** (MVP-4) |
| Aplicações do UE sem acesso à K_AKMA; cada uma só com a K_AF dos seus AF_IDs | TS 33.535 V19.0.0 | 4.4.3 | serviço de chaves do ME (`0005-akma-key-service.patch`); `akma-kaf` do `nr-cli` só com `akmaCliKaf` | [R2k], [R2t.me] (revisão 6) | **Verificado** (MVP-8) — antes, no MVP-4, o `nr-cli` entregava a K_AF a qualquer processo local |
| FQDN do AF em UTF-8 com NFKC na K_AF | TS 33.220 V19.1.0 | B.2.1.2 | ME recusa FQDN não ASCII | `nr-akma-selftest` (vetores de NFKC) | **Implementado** (MVP-4) — DV-4 |
| Cliente Ua\* do UE pelo plano de usuário | TS 33.535 V19.0.0 | 6.2.1 passo 1 | `af/ue/uaclient/` (`-me-socket`, `-profile`, `-bind`), `af/scripts/ue-netns.sh` | `teste-ue.md` §2 e §3; [R2t] | **Verificado** (MVP-4; perfis TLS no MVP-8) — D5 |

## Logs

| Requisito | Origem | Implementação | Teste | Estado |
|---|---|---|---|---|
| Nenhuma chave em log, em nível algum | plano, critério de todos os MVPs | UDM, AUSF, AAnF, UERANSIM | `TestNoKeyInLogs` (AAnF); VM: nenhuma cadeia de 64 hexadecimais nos logs e na tela do UE | **Verificado** em 18/10/2026 |
| Eventos AKMA em `info`, categoria `AKMA`, com SUPI, A-KID e AF_ID | D16 | `free5gc/patches/{udm,ausf}/`, `free5gc/aanf/` | VM (`docs/free5gc-side/verificacao.md` §4) | **Implementado** (MVP-3a) |
| AAnF na tela do `run.sh` (`-aanf`) | D15 | `free5gc/patches/run/0001-aanf-option.patch` | VM | **Implementado** (MVP-3a) |

## Ponto de referência Ua*

| Requisito | Norma | Cláusula | Implementação | Teste | Estado |
|---|---|---|---|---|---|
| Ua* simplificado com id de faixa privada | TS 33.220 H.2 NOTA 2 | — | `af/uastar/`, `af/af/`, `af/ue/uaclient/`; `docs/af-side/ua-star.md` | `docs/af-side/teste-af.md` | **Implementado** (MVP-3); perfil de teste `digest-http`, desligado por padrão, desde o MVP-8 — ADR-0003 revisto |
| Perfil K.2: Digest dentro de TLS, AF autenticado por certificado | TS 24.109 V19.0.0; TS 33.535 V19.0.0 | K.2; B.1.2.2 | `af/akmatls`, perfil `digest-tls` | revisão 6 | **Verificado** (MVP-8) |
| Perfil HTTP Digest — username, senha, realm | TS 24.109 V19.0.0; TS 33.535 V19.0.0; RFC 7616 | K.2; B.1.2.2; 3.3–3.5 | `af/uastar/` | `TestDigestVectors` (RFC 7616 §3.9.1 e Python independente); integração | **Implementado** (MVP-3) — DV-5 (RFC 7616, SHA-256), DV-6 (`realm = "3GPP-bootstrapping-akma@<FQDN>"`) |
| Confirmação de chave nos dois sentidos (`rspauth`) | TS 33.535 V19.0.0; RFC 7616 | 4.4.1; 3.5 | `af/af/internal/server/`, `af/ue/uaclient/` | `TestAFWithoutKaf`, `TestWrongKafRejected`; integração | **Implementado** (MVP-3) |
| *Product token* `3gpp-akma`; o UE verifica o FQDN no `realm` | TS 24.109 V19.0.0; TS 33.535 V19.0.0 | K.2 b) 1), d); B.1.2.2 passos 2 e 4 | `af/ue/uaclient/` | `TestChallengeChecks`, `TestHostAndProductToken` | **Implementado** (MVP-3) |
| FQDN do AF levado no `Host` e conferido pelo AF | TS 33.220 V19.1.0; TS 33.222 V19.0.0 | 4.5.2 pré-requisito (3); 5.3.0 passo 2 | `af/af/internal/server/` | `TestHostPathMethod`; integração (421) | **Implementado** (MVP-3) — `[INTERPRETAÇÃO]` 421 para outro nome |
| AF pede K_AF ao AAnF sem contexto ativo; `afId` e `aKId` | TS 33.535 V19.0.0; TS 29.535 V19.5.0 | 6.2.1 passos 2 e 9; 4.2.2.3 | `af/af/internal/aanf/` | `TestApplicationKey`; integração | **Implementado** (MVP-3) |
| Erro sem o `realm` do desafio quando o AF não obtém K_AF | TS 33.535 V19.0.0 | B.1.2.2 passo 6 | `af/af/internal/server/` | `TestKeyNotAvailable`, `TestAanfUnavailable` | **Implementado** (MVP-3) |
| Iniciação do AKMA pelo AF | TS 33.535 V19.0.0 | 6.5 | `af/af/internal/server/` (o `401` com desafio) | `TestInitiation`; integração | **Implementado** (MVP-3) |
| Expiração da K_AF no AF: nunca usada após `expiry`; renovar no AAnF ou recusar | TS 33.535 V19.0.0 | 4.4.0, 4.4.1, 6.4.2 | `af/af/internal/server/keys.go` (`kafExpiry`) | `TestKafExpiryRefetch`, `TestKafExpiryReject`, `TestExpiredOnArrival` | **Implementado** (MVP-3) |
| Perfil TLS 1.2 PSK | TS 24.109 V19.0.0 | K.3.2 | `af/wolftls`, perfil `psk-tls12` | revisão 6 | **Verificado** (MVP-8) |
| Perfil TLS 1.3 PSK — `"3GPP-AKMA" ";" A-KID` | TS 24.109 V19.0.0 | K.3.3 | `af/wolftls`, perfil `psk-tls13` | revisão 6 | **Verificado** (MVP-8) — defeito D-E |
| Ausência de renovação de K_AF | TS 33.535 §6.4.3; TS 24.109 K.2, K.3 | — | — | — | Especificado |

## Requisitos de ambiente

| Requisito | Origem | Verificação | Estado |
|---|---|---|---|
| Sincronização NTP entre a VM do AF e a do free5GC (par AAnF↔AF) | ADR-0001-R2 §5.3 | `timedatectl` nas duas VMs | Especificado |
| AF alcança a porta do AAnF | ADR-0001-R2 §5.1, §5.2 | requisição a partir de `192.168.16.129` | **Verificado** em 17/10/2026 (MVP-1) |
| UE alcança o AF pelo plano de usuário (`uesimtun0` → UPF → N6) | D5; `docs/free5gc-side/plano-de-usuario.md` | `ping` e Ua\* a partir do `uesimtun0` no *namespace* do UE; contadores da `upfgtp` | **Verificado** em 18/10/2026 (MVP-4): UE no próprio *namespace*, NAT no N6 |
| Proteção das interfaces SBI | TS 33.535 §4.4.0; TS 33.501 §13.1.0 | TLS habilitado | Decisão — DV-1; TLS no MVP-6+ |
