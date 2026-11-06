# Registro de alterações

## [Não lançado]

### MVP-11 — O NEF com o AF de fora e com o CCF, fase (c) (06/11/2026)
- **CCF:** o servidor interno em `https`, com TLS mútuo e a CA do operador; o provedor reconhecido pelo NF Instance ID do certificado; `Obtain_API_Invoker_Info` só com os itens do AEF que pergunta; a CRL da CA do CCF, assinada na partida, a cada saída e antes do `nextUpdate`; o certificado do *invoker* revogado na saída e recusado em CAPIF-1e.
- **NEF:** *patch* `nef/0005`: o CN do certificado igual ao `apiInvokerId` e ao `client_id`, a CA conferida contra `authenticationInfo`, `sourceAddresses` opcional, `invokerCrlMaxAge`.
- **`openapi`:** *patch* `0007`: a idade máxima da CRL, `sbitls.CrlChecker` e `CloseIdleConnections`.
- **AF:** a CRL no certificado do CCF; as conexões fechadas ao trocar de certificado.
- **Configuração:** `config/0006`, com `ccfUri` por nome, `invokerCrl: true` e `allowNrfTokens: false`; `ccfcfg.yaml`; `make-capif-certs.sh`, com o ponto de distribuição de CRL.
- **Roteiro:** S16 a S18, o servidor interno do CCF no R1s, E05 e CP18 revistos, `R1.profiles`. Execução de 06/11/2026: 317 verificações, nenhuma falha, nenhuma pulada.
- **Documentos:** `tls-sbi.md` §17, `capif.md`, `nef.md`, `af.md`, `e2e.md`, o tutorial, F13 e F27, DV-1 e DV-9 encerrados, DV-16 reduzido, D-AJ, a matriz na revisão 11, o plano.

### MVP-11 — Identidade nos *tokens* e nas políticas, fase (b) (05/11/2026)
- **`openapi`:** *patch* `0006`: a verificação do produtor em todas as NFs, só `ES256`, com o certificado de assinatura da CA do operador, `iss`, `aud`, `sub`, `exp` e escopo, e as respostas de TS 29.500 §6.7.3 (F5, F12, F23); `oauth.Signer`.
- **NRF:** *patch* `nrf/0005`: o certificado da sessão TLS confere o pedido de *token*, o registro, a atualização, o descadastro e a descoberta; *tokens* em `ES256`, com `kid`, `x5t#S256` e `iss` (F23, F26); `allowedOperationsPerNfType` (F4); saem a geração de certificados e o `rootcert` (F25).
- **AAnF:** a política pelo NF Instance ID do certificado, sem endereço de origem; escopos por operação; sem `supi-access`, a resposta vem sem `supi`.
- **AUSF e NEF:** *patches* `ausf/0003` e `nef/0004`, com os escopos novos; o segundo fecha o F14.
- **Configuração:** `config/0005`, com o `nrf-token.pem`; `test/0003`.
- **Roteiro:** `af/scripts/negative/sbi_token.py`, com S10 a S15, no passo R1s. Execução de 05/11/2026: 301 verificações, nenhuma falha, nenhuma pulada.
- **Documentos:** `tls-sbi.md` §16, `aanf.md`, `nrf.md`, `nef.md`, `openapi.md`, `logs-akma.md`, `e2e.md`, o tutorial, F4, F5, F12, F14, F23, F25 e F26, DV-1.

### MVP-11 — TLS mútuo nas SBI, fase (a) (04/11/2026)
- **PKI:** `free5gc/scripts/make-sba-certs.sh` e `free5gc/config/sba-entities.conf`: a CA do operador, um certificado por NF no perfil de TS 33.310 §6.1.3c.3, os certificados de teste e a conferência dos perfis.
- **CRL:** `free5gc/crl/`, um serviço HTTP só de leitura, com `run.sh -crl` (*patch* `run/0003`). A CRL é conferida nas duas pontas, com falha fechada (DV-18).
- **`openapi`:** *patch* `0005`, com o pacote `sbitls`: o perfil TLS, a CA, a CRL e a identidade do certificado. O cliente passa a conferir o servidor e a apresentar o certificado da NF (F21). `GetNFServiceUri` usa o nome (F24).
- **`util`:** *patch* `0001`, sem o registro dos segredos do TLS (F22).
- **NFs:** um *patch* `*-sbi-tls` por NF, com `sbi.fqdn` e `sbi.tls.ca`; `config/0004` põe o laboratório em `https`, pelos nomes. O modo por IP continua opcional (DV-17).
- **NEF:** a CA de cada grupo de rotas: a do operador na SBI e a do CCF no CAPIF-2e.
- **AAnF e `seed`:** TLS mútuo; `seed --cert`, `--key`, `--ca` e `seed token`.
- **AF:** `sbiTls` e `aefCa`; o AF de fora é só *invoker* do CAPIF, e o `afcfg-external.yaml` e o túnel SSH saíram.
- **Roteiro:** passo R1s e grupo S, com `af/scripts/negative/sbi_tls.py`; o grupo do NEF no modo CAPIF; o E07 com o AF_ID `af5`. O `tamper` dos testes do CAPIF trocava um caractere que só tem enchimento, e o CP08 falhava em uma de cada quatro credenciais.
- **Execução de 04/11/2026:** 288 verificações, nenhuma falha, nenhuma pulada; D03, O04 e O05 voltaram.
- **Documentos:** `tls-sbi.md` §15, `nef.md`, `aanf.md`, `capif.md`, `nrf.md`, `openapi.md`, `af.md`, `e2e.md`, o tutorial (parte 12 e os comandos das outras partes), F21 a F25 de `correcoes-free5gc.md`, DV-1, DV-9, DV-17 e DV-18.

### MVP-12 — Renovação da K_AKMA pelo AAnF (28/10/2026)
- **AAnF:** `akma.kakmaRefresh` (desligado por padrão). Ao entregar uma K_AF de uma K_AKMA mais velha que isso, o AAnF pede ao UDM, em segundo plano, uma nova autenticação primária (`Nudm_UECM_AuthTrigger`, TS 33.535 §6.4.4), uma vez por K_AKMA. Descobre o UDM no NRF e pede os *tokens* quando o OAuth2 está ligado. Cada contexto guarda quando a K_AKMA foi registrada.
- **`openapi`:** *patch* `0003`, com `AuthTriggerInfo`, `ReauthNotificationInfo`, `reauthNotifyCallbackUri` nos registros do AMF, a operação `AuthTrigger` (POST, DV-12) e a notificação ao AMF.
- **UDM:** *patch* `0002`, com `trigger-auth` (POST; GET recebe 405), a política de 60 s por SUPI e a notificação ao AMF sem *token* (TS 29.500 §6.7.3).
- **AMF:** *patch* novo `amf/0001`. O registro no UDM leva o `reauthNotifyCallbackUri`. O AMF recebe a notificação e responde pela tabela 6.2.5.6-3 de TS 29.503; com o UE conectado, autentica-o de novo com um ngKSI novo, faz o NAS SMC e volta a `Registered`; com o UE em CM-IDLE, marca a autenticação como pendente.
- **Ferramenta:** `seed trigger-auth`.
- **Roteiro:** passo R6r, com o AAnF numa cópia da configuração com `kakmaRefresh: 20s`. Execução de 28/10/2026: 161 verificações, nenhuma falha.
- **Documentos:** `renovacao-kakma.md` (especificação aprovada, implementação e verificação), `amf.md` (novo), `udm.md` §6, `aanf.md` §11, `openapi.md`, `logs-akma.md`, F16 a F20 de `correcoes-free5gc.md`, DV-12, D-O e D-P, parte 10 do tutorial; no guia público, a parte 11 e F16 a F20.

### MVP-8 — Ua\* da norma e autorização por aplicação no UE, fases (a) a (d) (20 a 22/10/2026)
- **Fase (a), serviço de chaves do ME:** `af/ue/ueransim/patches/0005-akma-key-service.patch`. O `nr-ue` atende num *socket* Unix e dá a cada aplicação, identificada pelo *uid* (`SO_PEERCRED`), só a K_AF dos seus AF_IDs (TS 33.535 §4.4.3). O `akma-kaf` do `nr-cli` só funciona com `akmaCliKaf`. O `uaclient` usa o serviço (`-me-socket`). Execução de 20/10/2026: 88 verificações, nenhuma falha.
- **Fase (b), perfis TLS** (TS 24.109 Anexo K):
  - K.2, Digest dentro de TLS 1.3 e 1.2, com o `crypto/tls` (`af/akmatls`);
  - PSK TLS 1.3 e PSK TLS 1.2 com a K_AF, pelo wolfSSL 5.9.4 e uma ligação cgo própria, `af/wolftls`, sob GPLv3;
  - o AF serve um perfil por porta; o Ua\* simplificado vira o perfil de teste `digest-http`, desligado por padrão;
  - identificadores Ua\* de TS 33.220 H.3, um AF_ID por perfil e *cipher suite* nas políticas do AAnF, do NEF (`free5gc/patches/config/0002`) e do UE (`0005`);
  - retomada por *tickets* presos à expiração da K_AF; sessão TLS encerrada quando a K_AF vence.
- **Fase (c), perfis DTLS** (TS 33.535 Anexo C):
  - DTLS 1.3 com a K_AF como PSK externa, AF_ID `…0100000009`; DTLS 1.2 com o procedimento e o AF_ID do TLS 1.2 (`…010001D001`);
  - CoAP dentro do DTLS (`af/coap`, RFC 7252): `GET coaps://<FQDN>/ua/v1/session`, resposta `2.05` em JSON; UDP 5684 e 5685;
  - *cookies* nas duas versões, e nenhuma consulta ao AAnF antes do *cookie*; MTU 1200; sem *connection ID*;
  - `af/wolftls` com DTLS: um *socket* UDP no AF, uma associação por endereço de origem, retransmissão pelos temporizadores do wolfSSL;
  - AF_ID `…0100000009` nas políticas do AAnF, do NEF e do UE;
  - sem `TLS_PSK_WITH_AES_128_CCM_8`, que RFC 7252 exige do CoAP e TS 33.210 não admite (DV-11).
- **Fase (d), perfil OSCORE** (TS 33.535 Anexo D):
  - `POST coap://<FQDN>/akma` em CBOR, `2.01` com N2 e UE-SID, e pedidos CoAP protegidos pelo OSCORE com o contexto da K_AF; UDP 5683; AF_ID `…0100000008`;
  - OSCORE (RFC 8613) com AES-CCM sobre o `crypto/aes` e o Apêndice B.2 (`af/oscore`), o CBOR necessário (`af/cbor`) e os *payloads* e derivações do Anexo D (`af/akmaoscore`), todos em Go e sem wolfSSL;
  - servidor CoAP sobre UDP em `af/coap`; `uaclient -profile oscore`, com `-rekey`;
  - as lacunas do Anexo D resolvidas pelos pontos U24 a U30, e cinco defeitos do texto no registro de defeitos (D-J a D-N);
  - AF_ID `…0100000008` nas políticas do AAnF, do NEF e do UE.
- **Scripts:**
  - `af/scripts/build-wolfssl.sh` (SHA-256 fixado, conferência das macros), `af/scripts/make-test-ca.sh`; `build-af.sh` compila com o wolfSSL quando a biblioteca existe;
  - fase (c): o wolfSSL com DTLS, MTU ajustável e leitura de datagramas maiores que o MTU; o `build-af.sh` compila o *hash* da biblioteca dentro de `af/wolftls`, para o Go ligar os binários de novo quando ela muda;
  - fase (d): `af/scripts/oscore-check.py`, a conferência independente do OSCORE com a libcoap.
- **Roteiro:** passo R2t e grupo T; um AF externo com PSK TLS 1.3 no R2n; `e2e-remote.sh` sem controle de *jobs*, para o Ctrl-C parar a parte remota. Execução de 21/10/2026: 115 verificações, nenhuma falha. Fase (c): passo R2d e grupo D, com o OpenSSL 4.1 no papel do UE (`OPENSSL_DTLS13`). Execução de 22/10/2026: 126 verificações, nenhuma falha. Fase (d): passo R2o e grupo O, com a libcoap no papel do UE (`LIBCOAP_CLIENT`). Execução de 22/10/2026: 141 verificações, nenhuma falha.
- **Documentos:** `docs/af-side/ua-star-norma.md` (especificação, implementação e verificação); ADR-0003 revisto; DV-10 (sem `D0 02` e sem OCSP *stapling*, recomendados por TS 33.210); DV-11; matriz revisões 6 a 8; registro de defeitos conferido de novo contra as normas (D-A a D-I são todos das normas), com D-J a D-N e o D-B como MVP-12; tutorial, partes 5 a 9; `af.md`, `ua-client.md`, `ue-ueransim.md`, `e2e.md`.

### MVP-6 — AF externo via NEF, acesso anônimo e UERANSIM em outra máquina (19/10/2026)
- **NEF:** `free5gc/patches/nef/0001-akma-api.patch` acrescenta ao NEF do free5GC a API `3gpp-akma` (TS 29.522 §5.14; TS 33.535 §6.3). O NEF:
  - autoriza o AF por AF_ID e endereço de origem;
  - descobre o AAnF pelo RID;
  - nunca repassa o SUPI;
  - repassa `403 K_AKMA_NOT_PRESENT` e o 204 do AAnF.
- **Configuração:** `free5gc/patches/config/0002-nef-akma.patch` põe o NEF em 192.168.16.130:8000, com `3gpp-akma` e a política de AF_ID.
- **AAnF:** passa para a porta 8010, fora das portas do free5GC. Aceita pedidos vindos do NEF pela lista `akma.authorizedNefs`, com endereço de origem e *token* (P11).
- **AF:**
  - modo externo (`nefUri`, `nrfUri`, `nfInstanceId`): registro no NRF e *token* `3gpp-akma`;
  - `anonymous` (TS 33.535 §6.2.2);
  - `af/config/afcfg-external.yaml`.
- **Roteiro:** `e2e.sh` com o grupo do NEF (`af/scripts/negative/nef_api.py`, E00 a E12) e o passo R2n. Execução de 19/10/2026: 83 verificações, nenhuma falha.
- **UERANSIM em outra máquina:** gNB e UE na VM 192.168.16.131, com o N2 e o N3 do free5GC em 192.168.16.130 (`docs/af-side/teste-nef.md`).
- **Documentos:**
  - `docs/free5gc-side/nef.md` e `docs/free5gc-side/correcoes-free5gc.md`, este com 15 pontos do free5GC com a cláusula e a forma da norma;
  - DV-9;
  - tutorial, parte 9;
  - matriz revisão 5;
  - o CAPIF planejado para o MVP-8 e o repositório público para o MVP-7.


### MVP-5 — ponta a ponta e testes negativos (18/10/2026)
- `af/scripts/e2e.sh` e `af/scripts/e2e-remote.sh`: o fluxo AKMA inteiro, nas duas máquinas, com uma senha de `sudo`. A sequência:
  - autoteste do ME;
  - free5GC com o AAnF, gNB e UE no *namespace*;
  - A-KID igual no UE, no AUSF e no AAnF; Ua\* pelo plano de usuário;
  - nova autenticação; AAnF parado (D9); assinante sem AKMA; `remove-context`.
  O relatório vai para `af/build/e2e/`.
- `af/scripts/negative/aanf_api.py`: testes negativos da API do AAnF em HTTP/1.1 cru (400, 401, 403, 404, 411, 413, 415, 501; nunca 204).
- `seed register --gpsi` e `--k-akma`; AF com `clockOffset` (teste de relógio); `uaclient -kaf-file`, com as chaves explícitas tendo precedência sobre o ME.
- Execução de 18/10/2026: 63 verificações, nenhuma falha.
- Matriz de rastreabilidade revisão 4, com as etapas de TS 33.535 §6.1 e §6.2.1 ligadas a evidências; desvios consolidados; tutorial, parte 8; `docs/af-side/e2e.md`. A pergunta-guia está respondida.
- Critério 4, verificado em 19/10/2026: o tutorial, seguido do início numa VM nova com o free5GC v4.2.3 do guia oficial, reproduz o resultado (63 verificações, nenhuma falha). A reprodução levou a duas correções:
  - `e2e.sh` e `e2e-remote.sh` tiram o endereço do AAnF do `afcfg.yaml` e do `aanfcfg.yaml`, em vez de 192.168.16.130 fixo; o `seed` recebe `--aanf`; o *preflight* confere que os dois endereços coincidem;
  - tutorial: nova seção 1.1 (repositório nas duas máquinas, SSH, gNB numa máquina só, assinante OPc, Ubuntu 18.04), caminho e `--aanf` do `seed`, qual `aanfcfg.yaml` ajustar, comandos completos da 6.3, aviso sobre o SUPI do vetor `testbed-ue` e `e2e.sh` com a máquina do free5GC.

### MVP-4 — lado UE: K_AF no ME e Ua\* pelo plano de usuário (18/10/2026)
- `af/ue/ueransim/patches/0003-akma-kaf-nr-cli-selftest.patch`:
  - o `nr-ue` deriva a K_AF sob demanda (TS 33.535 A.4, §6.2.1 passo 1) e a entrega, com o A-KID, pelo comando `akma-kaf` do `nr-cli`; a K_AKMA nunca sai do ME (§4.4.3);
  - FQDN não ASCII é recusado (DV-4);
  - `nr-akma-selftest` confere K_AKMA, A-KID e K_AF contra os vetores do MVP-0.
- `0004-free5gc-ue-gnb-on-veth.patch`: o UE procura o gNB no `veth` do seu *namespace*.
- `af/scripts/ue-netns.sh`: *namespace* do UE (D5), com `up`, `ue`, `cli`, `ua` e `down`.
- `af/scripts/patch-ueransim.sh`: confere a série numa cópia temporária, compila um alvo por vez com o `cmake` que configurou a árvore e também compila o autoteste.
- Cliente Ua\*: obtém A-KID e K_AF do ME pelo `nr-cli` (`-nr-cli`, `-ue`) e liga as conexões ao endereço da sessão PDU (`-bind`).
- Plano de usuário corrigido. O UE passou a rodar no próprio *namespace*, porque no do host o endereço da sessão é local e o kernel descartava os pacotes que saíam da `upfgtp`. `free5gc/scripts/n6-nat.sh` aplica o encaminhamento e o NAT do N6.
- Verificação na VM:
  - autoteste com 5 vetores aprovados e 3 de NFKC recusados (DV-4);
  - UE alcança 8.8.8.8 e o AF pelo `uesimtun0`;
  - A-KID igual no UE, no AUSF, no AAnF e no AF; Ua\* com a K_AF do UE, pelo plano de usuário;
  - nova autenticação troca o A-KID, e o antigo recebe 403 do AAnF.
- DV-4 implementado. Documentação:
  - `docs/af-side/{ue-ueransim,ua-client,teste-ue}.md` e `docs/free5gc-side/plano-de-usuario.md`;
  - tutorial, partes 5 (*namespace*) e 7; matriz, revisão 3.7.

### MVP-3 — AF e Ua\* simplificado (18/10/2026)
- `docs/af-side/ua-star.md`: especificação do Ua\* simplificado, aprovada antes do código. É o HTTP Digest de TS 24.109 K.2 sem TLS, com:
  - `username` = A-KID e senha = base64(K_AF);
  - `realm = "3GPP-bootstrapping-akma@<FQDN>"` (DV-6), SHA-256 (DV-5) e `rspauth` obrigatório;
  - FQDN conferido no `Host` (TS 33.220 §4.5.2) e iniciação do AKMA pelo desafio (§6.5);
  - política `refetch` ou `reject` para a K_AF vencida (§6.4.2).
- `af/`: módulo Go com:
  - `uastar`: Digest de RFC 7616, compartilhado;
  - o AF (`af/af`): servidor Ua\*, `Naanf_AKMA_ApplicationKey_Get` pelo cliente do `openapi` corrigido e K_AF guardada até `expiry`;
  - o cliente Ua\* (`af/ue/uaclient`), que recusa AF sem `rspauth` válido;
  - `af/config/afcfg.yaml` e `af/scripts/build-af.sh`.
- `af/uastar/testdata/gen_digest_vectors.py`: implementação Python independente do Digest; reproduz o exemplo de RFC 7616 §3.9.1 e gera os vetores do projeto.
- Verificação com o AAnF da VM do free5GC e um contexto semeado:
  - a K_AF do oráculo autentica;
  - K_AF errada → 401; A-KID desconhecido → 403; outro `Host` → 421; AF_ID fora da política do AAnF → 503;
  - nenhuma chave nos logs.
- DV-5 e DV-6 implementados. Documentação: `docs/af-side/{ua-star,af,teste-af}.md`; tutorial, parte 6; matriz, revisão 3.6.

### MVP-3a — mensagens AKMA na tela e derivação no ME (18/10/2026)
- Pedido do responsável: acompanhar o AKMA pelas telas do `run.sh` e do UE. Decisões D14 (escopo; derivação no ME antecipada do MVP-4), D15 (`run.sh -aanf`) e D16 (SUPI, A-KID e AF_ID em `info`).
- `free5gc/patches/udm/0001-*.patch` e `free5gc/patches/ausf/0001-*.patch` regenerados:
  - categoria `AKMA` no UDM;
  - no UDM, a indicação enviada em `info`;
  - no AUSF, em `info`: indicação recebida, K_AKMA e A-KID derivados (com o A-KID), AAnF selecionado e registro.
- AAnF: identificadores em `info` e uma linha de partida com RIDs, tempo de vida da K_AF e AF_ID autorizados.
- `free5gc/patches/run/0001-aanf-option.patch`: opção `-aanf` no `run.sh` e `aanf` no `force_kill.sh`. `patch-nf.sh` ganha o alvo `run` e passa a checar alterações locais só nos arquivos que os *patches* tocam. `run-aanf.sh` exporta `GIN_MODE=release`.
- `af/ue/ueransim/patches/`: o `nr-ue` deriva K_AKMA, A-TID e A-KID (TS 33.535 §6.1, A.2, A.3) ao ativar o contexto de segurança da autenticação primária, com a opção `akma`, e mostra o A-KID em `[nas] [info]`. `af/scripts/patch-ueransim.sh` aplica, reverte e recompila só o `nr-ue`.
- Verificação na VM: o A-KID do UE é igual ao do AUSF e do AAnF em duas autenticações. Sem `akma`, o UE não deriva nada. Não há chave nos logs nem nas telas. Os testes `TestAkma*` continuam passando.
- Documentação: `docs/free5gc-side/logs-akma.md`, `docs/af-side/ue-ueransim.md`, `verificacao.md` §4; tutorial, parte 5 (as seguintes renumeradas); matriz, revisão 3.5.

### MVP-2 — UDM e AUSF (17/10/2026)
- `free5gc/patches/udm/0001-akma-indication-and-routing-id.patch`: o UDM preenche `akmaInd` e `routingId` (TS 33.501 §6.1.3.1 e §6.1.3.2.0; TS 29.503 Tab. 6.3.6.2.3-1) a partir de `akmaAllowed` e `routingId` do UDR (TS 29.505 Tab. 5.4.2.2-1). Compara com o RID do SUCI e registra divergências (D10).
- `free5gc/patches/ausf/0001-akma-anchor-key-register.patch`: depois do 5G AKA ou do EAP-AKA' bem-sucedido, o AUSF deriva K_AKMA e o A-KID e os registra no AAnF (TS 33.535 §6.1), sem bloquear a autenticação (D9). O AAnF é descoberto no NRF pelo Routing Indicator ou configurado (`akma.aanfUri`).
- `free5gc/patches/config/0001-ausf-akma.patch`: `akma.enable` no `ausfcfg.yaml`.
- `free5gc/patches/test/0001-akma-tests.patch`: AAnF no banco de testes do free5GC; testes `TestAkma5GAKA`, `TestAkmaEAPAKAPrime` e `TestAkmaWithoutSubscription`, que rodam sem `sudo`.
- `free5gc/scripts/provision-akma-subscriber.sh`; `prepare-openapi.sh` com argumento de versão; `patch-nf.sh` com os alvos `test` e `config`.
- Verificação: K_AF do AAnF igual à derivada pelo UE de teste nos dois métodos; A-KID `rid0000…@5gc.mnc093.mcc208.3gppnetwork.org` com UERANSIM; reautenticação substitui o contexto; AAnF parado não impede o registro; `TestRegistration`, `TestGUTIRegistration`, `TestReSynchronization` e `TestEAPAKAPrimeAuthentication` continuam passando.
- Documentação: `docs/free5gc-side/{udm,ausf,assinatura-akma,verificacao}.md`; tutorial, parte 4; matriz, revisão 3.4.

### MVP-1b — correção do NRF do free5GC (17/10/2026)
- `free5gc/patches/nrf/`: três *patches* com testes sobre o `free5gc/nrf` v1.4.5:
  - 0001: guarda `aanfInfoList` (TS 29.510 §5.2.2.2.2, Tab. 6.1.6.2.2-1);
  - 0002: descoberta de `AANF` com filtro por `routing-indicator` sobre `aanfInfoList` (Tab. 6.2.3.2.3.1-1, Tab. 6.1.6.2.73-1);
  - 0003: `aud` do token com a instância ou o tipo do produtor (Tab. 6.3.5.2.4-1).
- `free5gc/scripts/patch-nf.sh`: aplica, reverte e informa *patches* de uma NF do free5GC e a recompila, guardando o binário original.
- AAnF: `aud` passa a ser exigido; DV-8 encerrado. `seed discover` aceita `--routing-indicator`.
- Aplicado na VM do free5GC. Verificados `aanfInfoList` guardado, descoberta por RID (com o filtro conferido no MongoDB 5.0) e token com `aud`. Regressão parcial com nove NFs; a regressão com UE fica com o responsável (exige `sudo`).
- Decisões: D13 revista (NRF corrigido agora); D8 (TS 29.505 V19.6.0 incorporada ao baseline); D10 (RID vem do `routingId` do UDR, com comparação ao RID do SUCI em log).
- Documentação: `docs/free5gc-side/nrf.md`; tutorial, seção 3.3; matriz, revisão 3.3.

### MVP-1 — AAnF autônomo (17/10/2026)
- `free5gc/patches/openapi/`: `0001-naanf-akma-models.patch` (`AkmaKeyInfo`, `CtxRemove`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif`) e `0002-naanf-akma-client.patch` (cliente `aanf/AKMA`) sobre `free5gc/openapi` v1.2.4; `free5gc/scripts/prepare-openapi.sh`.
- `free5gc/aanf/`: NF AAnF com a estrutura das NFs do free5GC. Três operações do Naanf_AKMA (TS 29.535); contexto em memória com substituição do par (A-KID, K_AKMA) e K_AF reaproveitada até expirar (TS 33.535 §6.1, §6.2.1, §5.2); política local que associa AF_ID a endereço de origem; OAuth2 com verificação de `exp` e `aud`; erros de TS 29.500; registro e cancelamento no NRF.
- `free5gc/aanf/tools/seed/`: ferramenta de teste nos papéis do AUSF (`register`) e do AF (`retrieve`), mais `remove` e `discover`.
- `free5gc/config/aanfcfg.yaml`; `free5gc/scripts/{build-aanf,run-aanf,deploy}.sh`.
- Verificação na VM do free5GC e a partir da VM do AF: K_AF igual ao vetor `testbed-ue`; 403 para AF fora da política e para AF_ID pedido de outra máquina.
- Decisões: DV-8 (`aud` vazio tolerado no token, por limitação do NRF do free5GC); D13 (limitações do NRF v1.4.5 — `aanfInfoList` descartado, descoberta de `AANF` recusada, `aud` vazio — a corrigir com *patch* futuro). DV-1, DV-2 e DV-7 implementados.
- Documentação: `docs/free5gc-side/{aanf,openapi,registro-nrf}.md`; tutorial, parte 3; matriz de rastreabilidade, revisão 3.2.
- MVP-0: `go test -race` passa em `shared/akma` (critério 1 atendido por inteiro).

### MVP-0 — biblioteca AKMA comum (16/10/2026)
- `shared/akma`: KDF de TS 33.220 B.2 com NFKC; SUPI como parâmetro (TS 33.501 A.7.0, só IMSI — DV-3); Routing Indicator, Home Network Domain e A-KID (TS 23.003, TS 29.522); AF_ID nas formas da KDF e de transporte; K_AKMA, A-TID e K_AF (TS 33.535 Anexo A); codecs por camada (ADR-0002).
- `shared/crosscheck`: comparação com `ueauth.GetKDFValue` do free5GC (`free5gc/util` v1.3.2).
- `reference/kdf/gen_vectors.py` e `shared/akma/testdata/vectors.json`: 8 vetores do projeto, igualdades esperadas e contraexemplos; não são vetores 3GPP.
- Decisão P9: identificador Ua\* no `AfId` em maiúsculas na escrita (errata no `00` §6.4 e no ADR-0003).
- Documentação: `docs/shared/` (KDF, identificadores, codificação, vetores), `docs/tutorial/` (pré-requisitos e conferência das derivações), `docs/free5gc-side/` e `docs/af-side/` criados.

## [2.1] — 2026-10-16

Verificação do ambiente, conferência contra o texto das normas, revisão da topologia e plano de desenvolvimento aprovado.

### Novo
- `docs/baseline/ambiente-verificado.md`: free5GC v4.2.3 e submódulos fixados; `free5gc/openapi` v1.2.4 já contém `AANF`, `naanf-akma`, `AanfInfo`, `AkmaInd` e `RoutingId` (questão 5 encerrada).
- `docs/baseline/conferencia-normativa.md`: conferência da documentação contra TS 33.535, TS 29.535 e as demais normas do baseline.
- ADR-0001-R2: AAnF na VM do free5GC.
- Adendo E: TS 23.502 do acervo na V19.9.0; baseline mantido em V19.8.0.
- `docs/literatura/analise-dos-artigos.md`: 14 trabalhos sobre AKMA.
- `docs/conformidade/desvios.md`: desvios DV-1 a DV-7.
- `docs/plano/plano-de-desenvolvimento.md`: MVP-0 a MVP-5, aprovado.

### Revisto
- `realm` do Digest do Ua*: passa a seguir TS 33.535 B.1.2.2, `3GPP-bootstrapping-akma@<FQDN>` (DV-6; defeito D-C).
- Registro de defeitos: correspondência entre as letras A–E e D-A a D-I; novos D-F (RFC 3310 em TS 24.109 K.2), D-G (204 inalcançável em TS 29.535), D-H e D-I (editoriais).
- ADR-0001 marcado como revisto; ADR-0002 e ADR-0003 aceitos; ADR-0003 registra por que não usa `(0x01,0x00,0x00,0x00,0x02)`.
- Matriz de rastreabilidade, revisão 3.1.
- Normas e artigos ficam fora do controle de versão (`/specs/`, `/papers/`).

### Corrigido (errata marcada no texto)
- `00` §2.5: espaço espúrio na cadeia `S` do A-TID.
- `00` §5.1: guarda de K_AUSF citada em TS 33.501 §6.1.1.1 e §6.2.2.1 (não §6.2.2).
- `00` §2.4: forma hexadecimal do identificador Ua\* no `AfId` reclassificada como interpretação; tabela de H.3 completada.
- Adendo A §A.5: `immediateFlag` em TS 29.503 §5.5.2.2.2, com relatório imediato não garantido.
- Adendo B: marcação `[HIPÓTESE]` substituída; lista de cabeçalhos obrigatórios de TS 29.500 completada.
- Adendo D: precedente ancorado de `Kausf` (TS 29.503); conflito RFC 3310 × RFC 7616 (DV-5); decisão do `realm` revista.

### Removido
- `scripts/bootstrap-repo.sh`: fazia commits e criava repositório público, o que contraria as regras do projeto; o repositório já está inicializado.

## [2.0] — 2026-09-02

Consolidação da fase F0 (aquisição documental) e da fase F1 (decisões de codificação).

### Baseline
- Release 19 fixado. Dezesseis especificações incorporadas.
- Par corpo/OpenAPI da TS 29.535 fixado em V19.5.0 + `TS29535_Naanf_AKMA.yaml` v1.2.0.

### Encerrado
- Definição da KDF — TS 33.220 Anexo B.2.0.
- Codificação do SUPI como parâmetro de entrada — TS 33.501 Anexo A.7.0.
- Identificador de protocolo de segurança Ua* — TS 33.220 Anexo H.
- Tipos `AKId`, `AfId`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` — TS 29.522.
- Existência e Stage 3 de `Nudm_UECM_AuthTrigger` — TS 33.501 §14.2.6, TS 29.503 §5.3.2.16.
- Reconhecimento do AAnF como consumidor do UDM — TS 29.503 §5.2.2.1.
- `NFType` `AANF`, `AanfInfo` e descoberta por `routing-indicator` — TS 29.510.
- Codificação de `kAkma` e `kaf` — por precedente de TS 29.509.
- Equivalência de `DateTime` entre TS 29.122 e TS 29.571.

### Revisto
- Codificação do valor do A-TID: recomendação alterada de hexadecimal para base64, por analogia com o B-TID do GBA (TS 33.220 §4.5.2). Registrado em ADR-0002.

### Novo
- Stage 3 do ponto de referência Ua* — TS 24.109 Anexo K, três perfis.
- Registro de cinco defeitos em especificações.
- ADR-0001, ADR-0002 e ADR-0003.
- Matriz de rastreabilidade, revisão 3.
