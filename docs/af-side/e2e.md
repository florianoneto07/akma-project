# Ponta a ponta e testes negativos (MVP-5 a MVP-12)

**Data:** 18 de outubro de 2026; MVP-6 em 19 de outubro de 2026; MVP-8 de 20 a 22 de outubro de 2026; MVP-12 em 28 de outubro de 2026 · **Resultado do MVP-12:** 161 verificações, nenhuma falha (execução `20261028-162316`, seção 5.7) · **MVP-8, fase (d):** 141 verificações, nenhuma falha (execução `20261022-201824`, seção 5.6) · **Fase (c):** 126 verificações, nenhuma falha (execução `20261022-124514`, seção 5.5) · **Fase (b):** 115 verificações, nenhuma falha (execução `20261021-164401`, seção 5.4) · **Resultado do MVP-6:** 83 verificações, nenhuma falha (execução `20261019-213720`, seção 5.2) · **MVP-5:** 63 verificações, nenhuma falha (execução `20261018-171510`) · **Instalação limpa:** 63 verificações, nenhuma falha, em 19 de outubro de 2026 (execução `20261019-160723`, seção 5.1)

## 1. O que o roteiro faz

`af/scripts/e2e.sh` executa o fluxo AKMA inteiro no laboratório e os testes negativos do plano, de forma reprodutível:

- autenticação primária real;
- K_AKMA e A-KID no AUSF e no UE;
- registro no AAnF;
- K_AF entregue ao AF;
- autenticação Ua\* do UE pelo plano de usuário.

Ele roda na máquina do AF e conduz a do free5GC por SSH. Lá, pede a senha de `sudo` **uma vez**: a parte remota (`af/scripts/e2e-remote.sh`) roda numa única sessão `ssh -t`, porque o free5GC (UPF), o *namespace* e o UE precisam de `root`.

As duas máquinas se sincronizam por arquivos-marcadores em três pontos (A, B, C). Em cada um, a parte remota espera enquanto a máquina do AF faz as verificações que precisam partir dela: o AAnF autoriza a K_AF por AF_ID e endereço de origem, e só a máquina do AF é o AF autorizado.

| Grupo | Onde | O que verifica |
|---|---|---|
| S1 | AF, antes da execução | autoteste do ME (`nr-akma-selftest`) com os vetores do MVP-0 |
| R1 | free5GC | NAT do N6, *namespace* do UE, `linkIp` do gNB, assinatura AKMA; free5GC com `-aanf` e, no MVP-9, `-ccf`, gNB e UE no ar. MVP-9: o CCF no ar, `3gpp-akma` publicado pelo NEF no CCF e a credencial de entrada do AF emitida. MVP-11: a CRL assinada de novo e o serviço da CRL no ar, com `-crl` |
| R1s | free5GC | MVP-11: o TLS mútuo da SBI (seção 4.3): cada porta, o perfil dos certificados e o `seed` contra um `openssl s_server` |
| R2 | free5GC | registro; A-KID do UE = do AUSF = do AAnF; Ua\* do UE pelo plano de usuário, 3 requisições |
| N (ponto A) | AF | testes negativos da API do AAnF, do `register-anchorkey`, do RID e do AF (seção 4); acesso anônimo direto ao AAnF (N01c, MVP-6) |
| S (ponto A) | AF | MVP-11: o TLS mútuo da SBI da máquina do AF, no AAnF e no NEF (seção 4.3) |
| E (ponto A) | AF | MVP-6: três AF extras (externo na 8083, externo anônimo na 8084, interno anônimo na 8085) e, no MVP-8, um externo com PSK TLS 1.3 na 8456; API do NEF (seção 4.1). MVP-11: os externos são *invokers* do CCF |
| R2k | free5GC | MVP-8, fase (a): serviço de chaves do ME (TS 33.535 §4.4.3). O `uaclient` recebeu a K_AF do seu AF_ID; outro AF_ID, outro usuário e o `akma-kaf` do `nr-cli` são recusados; a configuração com um AF_ID em duas aplicações não sobe |
| R2t | free5GC | MVP-8, fase (b): Ua\* do UE pelos perfis TLS, com a K_AF do ME para o AF_ID de cada perfil e *cipher suite*. K.2 em TLS 1.3 (duas conexões, a segunda retomada) e em TLS 1.2; PSK TLS 1.3 e PSK TLS 1.2, duas conexões cada; `ClientHello` só com `psk_ke` recusado |
| R2d | free5GC | MVP-8, fase (c): Ua\* do UE pelos perfis DTLS, com CoAP: DTLS 1.3 com o AF_ID `…0100000009` e DTLS 1.2 com `…010001D001`, duas associações cada, a segunda retomada; `ClientHello` só com `psk_ke` recusado |
| R2o | free5GC | MVP-8, fase (d): Ua\* do UE pelo OSCORE: `POST /akma` em CBOR e dois pedidos protegidos no mesmo contexto, AF_ID `…0100000008`; contexto renovado pelo Apêndice B.2 da RFC 8613 antes de cada pedido; um `POST` novo para cada pedido |
| D (ponto A) | AF | MVP-8, fase (c): o que o AF viu do R2d; DTLS 1.3 com o `openssl s_client` do OpenSSL 4.1 no papel do UE, mandando o GET em CoAP (com `OPENSSL_DTLS13`; sem ela, pulada); A-KID desconhecido e K_AF de outro AF_ID recusados |
| O (ponto A) | AF | MVP-8, fase (d): o que o AF viu do R2o (quatro contextos, uma consulta ao AAnF, seis pedidos atendidos, duas renovações); o `coap-client` da libcoap no papel do UE, com o pedido e o contexto calculados por `af/scripts/oscore-check.py` (com `LIBCOAP_CLIENT`; sem ela, puladas); OSC-INP com outro algoritmo, A-KID desconhecido e K_AF de outro AF_ID recusados; vencimento da K_AF com `refetch` e `reject`, nos AF de relógio adiantado |
| T (ponto A) | AF | MVP-8, fase (b): o que o AF viu do R2t (AF_ID de cada *cipher suite*; retomadas); PSK TLS 1.3 com o `openssl s_client` no papel do UE; recusas (K_AF de outro AF_ID, A-KID desconhecido, `psk_ke`, certificado de outro FQDN, AF_ID fora da política do AAnF); expiração da K_AF numa sessão TLS, com as políticas `refetch` e `reject` |
| CP (ponto A) | AF | MVP-9: o AF como *API invoker* do CAPIF na 8086, com a CA do CCF e a credencial copiadas da máquina do free5GC (CP00), e o segundo *invoker*, `capif_api.py` (seção 4.2) |
| R2n | free5GC | MVP-6: Ua\* do UE com os três AF extras; o AAnF entrega via NEF; o NEF repassa a K_AF e não registra o SUPI. MVP-8: Ua\* por PSK TLS 1.3 com um quarto AF extra, externo, na porta 8456 |
| R2c | free5GC | MVP-9: Ua\* do UE com o AF do modo CAPIF; o NEF aceita o *token* do CCF para o AF_ID `af1.akma.example.FF00000001` |
| R3 | free5GC | nova autenticação: A-KID novo, contexto substituído, Ua\* com o novo; Ua\* com o A-KID e a K_AF anteriores, ainda aceitos pelo AF até `expiry` |
| B (ponto B) | AF | os AF extras se descadastram do NRF (E12); o AF do modo CAPIF sai do CCF (CP19); o AAnF recusa o A-KID anterior e aceita o novo |
| R4 | free5GC | AAnF parado durante uma autenticação (D9) |
| R5 | free5GC | assinante sem assinatura AKMA; depois restaurada |
| R6 | free5GC | `remove-context`, e de novo (404) |
| R6r | free5GC | MVP-12: renovação da K_AKMA pelo AAnF (TS 33.535 §6.4.4). O AAnF sobe com `kakmaRefresh: 20s` (cópia da configuração); K_AF de uma K_AKMA nova não pede renovação; depois de 20 s, K_AF para outro AF_ID, que o AF não tem guardada: o AAnF pede o `trigger-auth`, o UDM notifica o AMF, o AMF autentica o UE registrado e faz o NAS SMC, o UE deriva um A-KID novo, igual ao do AUSF e ao do AAnF, e a sessão PDU continua. Depois, o `seed trigger-auth`: 403 dentro de 60 s, 405 com GET, 400 sem `supi`, 404 para SUPI desconhecido. A configuração instalada volta no fim |
| C (ponto C) | AF | `retrieve` depois do `remove-context` |
| R7, L1 | os dois | nenhuma cadeia de 64 hexadecimais nos logs do free5GC, do UE, do gNB e do AF; no MVP-9, nenhum *token*, credencial, segredo de entrada ou chave privada nos logs do free5GC e do AF do modo CAPIF (CP20); no MVP-11, nenhum arquivo de segredos do TLS escrito na execução e nenhum segredo do TLS nos logs (S19) |

## 2. Pré-requisitos

Na máquina do free5GC, as partes 3 a 7 e 10 do tutorial, isto é:

- o repositório em `~/akma-free5gc`;
- os *patches* do NRF, do UDM (os dois), do AUSF, do AMF (parte 10, MVP-12), da configuração e do `run.sh`;
- o AAnF implantado (`bin/aanf`, `config/aanfcfg.yaml` com a linha `kakmaRefresh`, que o R6r troca numa cópia) e o `seed` em `~/akma-tools/`;
- o CCF e os *patches* do NEF da parte 11 (MVP-9): `bin/ccf`, `config/ccfcfg.yaml` e os certificados de `make-capif-certs.sh` em `cert/`; `nef/0001` e `0002`, `config/0001` a `0003` e `run/0001` e `0002`;
- o TLS da SBI da parte 12 (MVP-11): os *patches* `*-sbi-tls`, `config/0004`, `run/0003` e `test/0002`, as cópias de `openapi` e `util`, `bin/crl`, `config/crlcfg.yaml`, os certificados de `make-sba-certs.sh` em `cert/sba/` e os nomes da SBI no `/etc/hosts`;
- o UERANSIM com os *patches*, o `linkIp: 10.100.0.1` do gNB e o `uaclient` em `~/akma-free5gc/af/build/bin/`.

Na máquina do AF: `af/scripts/build-wolfssl.sh`, `af/scripts/make-test-ca.sh`, `af/scripts/build-af.sh` e `free5gc/scripts/build-aanf.sh`; o `openssl` (o grupo T o usa como UE independente); opcionalmente, um OpenSSL 4.1 compilado à parte para o D03 e a libcoap para O04 e O05 (seção 3); acesso SSH por chave à máquina do free5GC; o `af/config/afcfg.yaml` com o `aanfUri` do AAnF da máquina do free5GC; e, no `/etc/hosts`, os nomes do AAnF e do NEF.

O roteiro sobe o AF com os perfis de `afcfg.yaml` e acrescenta o perfil de teste `digest-http` na porta 8080, que os grupos R2, N, R3 a R5 usam. Ele copia a CA de teste para a máquina do free5GC e, de lá, para `af/build/sba/`, os certificados da SBI desta máquina, a cada execução. A chave da CA do operador fica na máquina do free5GC. As políticas do AAnF, do NEF e do UE precisam dos AF_IDs dos perfis (`docs/af-side/af.md` §3).

O roteiro confere esses itens no início (*preflight*) e mostra o estado dos *patches*. O NAT do N6 e o *namespace* do UE ele mesmo aplica.

## 3. Como rodar

```bash
af/scripts/build-af.sh
```

```bash
af/scripts/e2e.sh woca@192.168.16.130
```

O argumento é a máquina do free5GC, `usuario@endereço`. Sem ele, vale a do laboratório, `woca@192.168.16.130`.

A verificação D03 confere o DTLS 1.3 com o `openssl s_client` no papel do UE. O OpenSSL do Ubuntu 24.04, 3.0, não tem DTLS 1.3 no `s_client`; o 4.1 tem. Ele fica fora do repositório, compilado à parte; a verificação foi feita com o 4.1.0 *beta* 2:

```bash
mkdir -p ~/openssl41 && cd ~/openssl41 && wget https://github.com/openssl/openssl/releases/download/openssl-4.1.0-beta2/openssl-4.1.0-beta2.tar.gz
```

```bash
echo "457df0ec01a148954f843cb3589ff7465bdacc50250c7aba37e20f6c30c4e8b5  openssl-4.1.0-beta2.tar.gz" | sha256sum -c -
```

```bash
tar -xzf openssl-4.1.0-beta2.tar.gz && cd openssl-4.1.0-beta2 && ./Configure --prefix="$HOME/openssl41/install" --libdir=lib no-shared no-tests no-docs && make -j"$(nproc)" && make install_sw
```

O roteiro recebe o binário em `OPENSSL_DTLS13`. Sem ele, o D03 aparece como pulado (`SKIP`):

```bash
OPENSSL_DTLS13=$HOME/openssl41/install/bin/openssl af/scripts/e2e.sh woca@192.168.16.130
```

As verificações O04 e O05 usam o `coap-client` da libcoap, com OSCORE, no papel do UE. A versão conferida é a 4.3.5b, ligada ao OpenSSL 3.0 do Ubuntu 24.04. Sem CMake, sem compilador C++ e sem os cabeçalhos do OpenSSL na máquina, os passos são estes (`ua-star-norma.md` §7.11):

```bash
mkdir -p ~/libcoap && cd ~/libcoap && wget https://github.com/Kitware/CMake/releases/download/v3.31.6/cmake-3.31.6-linux-x86_64.tar.gz https://github.com/Kitware/CMake/releases/download/v3.31.6/cmake-3.31.6-SHA-256.txt
```

```bash
grep " cmake-3.31.6-linux-x86_64.tar.gz$" cmake-3.31.6-SHA-256.txt | sha256sum -c - && tar -xzf cmake-3.31.6-linux-x86_64.tar.gz
```

```bash
wget https://github.com/openssl/openssl/releases/download/openssl-3.0.13/openssl-3.0.13.tar.gz && echo "88525753f79d3bec27d2fa7c66aa0b92b3aa9498dafd93d7cfa4b3780cdae313  openssl-3.0.13.tar.gz" | sha256sum -c -
```

```bash
tar -xzf openssl-3.0.13.tar.gz && (cd openssl-3.0.13 && ./Configure --prefix=/usr --libdir=lib/x86_64-linux-gnu linux-x86_64 && make build_generated)
```

```bash
wget -O libcoap-4.3.5b.tar.gz https://github.com/obgm/libcoap/archive/refs/tags/v4.3.5b.tar.gz && echo "383a17d8466cee7c1cb1d4dfbffad2651004850b29eb590e9591c7bedd46741d  libcoap-4.3.5b.tar.gz" | sha256sum -c -
```

```bash
tar -xzf libcoap-4.3.5b.tar.gz && cd libcoap-4.3.5b && sed -i 's/  LANGUAGES CXX C)/  LANGUAGES C)/; s/CXX_COMPILER_ID/C_COMPILER_ID/g' CMakeLists.txt
```

```bash
../cmake-3.31.6-linux-x86_64/bin/cmake -S . -B build -DDTLS_BACKEND=openssl -DOPENSSL_INCLUDE_DIR=$HOME/libcoap/openssl-3.0.13/include -DOPENSSL_SSL_LIBRARY=/usr/lib/x86_64-linux-gnu/libssl.so.3 -DOPENSSL_CRYPTO_LIBRARY=/usr/lib/x86_64-linux-gnu/libcrypto.so.3 -DENABLE_DOCS=OFF -DENABLE_OSCORE=ON -DBUILD_SHARED_LIBS=OFF -DCMAKE_INSTALL_PREFIX=$HOME/libcoap/install
```

```bash
../cmake-3.31.6-linux-x86_64/bin/cmake --build build -j"$(nproc)" && ../cmake-3.31.6-linux-x86_64/bin/cmake --install build
```

O roteiro recebe o `coap-client` em `LIBCOAP_CLIENT`, ao lado do `OPENSSL_DTLS13`. Sem ele, O04 e O05 aparecem como pulados:

```bash
OPENSSL_DTLS13=$HOME/openssl41/install/bin/openssl LIBCOAP_CLIENT=$HOME/libcoap/install/bin/coap-client af/scripts/e2e.sh woca@192.168.16.130
```

Os endereços vêm de duas configurações:

- o do AAnF, do `aanfUri` de `af/config/afcfg.yaml`. No *preflight*, o roteiro confere que ele é o `sbi.bindingIPv4` do `~/free5gc/config/aanfcfg.yaml` instalado e para se não for. A parte remota toma o endereço desse mesmo `aanfcfg.yaml`. O `seed` recebe o endereço em `--aanf` nas duas máquinas;
- o dos AF extras, de relógio deslocado, do `listen` de `af/config/afcfg.yaml`, nas portas 8081 e 8082.

Com `--keep`, o free5GC, o gNB e o UE continuam no ar no fim. Sem ele, o roteiro os encerra. Os arquivos ficam em `af/build/e2e/<data-hora>/`:

| Arquivo | Conteúdo |
|---|---|
| `report.md` | resumo, A-KID de cada autenticação, tabela de verificações e as linhas AKMA do free5GC, do UE e do AF |
| `local-results.txt`, `remote/remote-results.txt` | uma linha `PASS`, `FAIL` ou `SKIP` por verificação |
| `remote/` | `free5gc.log`, `ue.log`, `gnb.log`, saídas do cliente Ua\* e do `seed`, `state` (os A-KID) |
| `af.log`, `af-*.log` | o AF, e os AF de relógio deslocado |

O código de saída é 0 só se nenhuma verificação falhar; uma verificação pulada não conta como falha. Leva cerca de quatro minutos.

## 4. Testes negativos

A lista do plano (MVP-5), com o resultado esperado de cada um:

| Caso do plano | Verificações | Esperado |
|---|---|---|
| A-KID desconhecido | N06; R5.ua | `403 K_AKMA_NOT_PRESENT` no AAnF; `403 AKMA_KEY_NOT_AVAILABLE` no AF |
| rótulo `akid` aceito | N02 | 200 (ADR-0002) |
| RID `"012"` × `"12"` | N19, N20, N21 | o mesmo A-TID com `rid12` é outro A-KID: 403 |
| AF não autorizado | N03 | 403 sem causa (P3) |
| `afId` malformado (9 dígitos, sem ponto) | N04 | `400 MANDATORY_IE_INCORRECT` |
| `AkmaKeyInfo` com `supi` e `gpsi` | N17 | `400 OPTIONAL_IE_INCORRECT` |
| `kAkma` fora do padrão | N18 | `400 MANDATORY_IE_INCORRECT` |
| `Content-Type` incorreto | N09 | 415 |
| corpo acima do limite | N10 | `413 MAX_JSON_SIZE_EXCEEDED` |
| demais códigos do MVP-1 | N05, N07, N08, N11 a N16 | `400 MANDATORY_IE_INCORRECT`, `MANDATORY_IE_MISSING`, `INVALID_MSG_FORMAT`; `411 INCORRECT_LENGTH`; 501; `404 RESOURCE_URI_STRUCTURE_NOT_FOUND`; `400 INVALID_API`; 401 sem *token*; nunca 204 (DV-7) |
| relógio do AF deslocado | N22, N23 | 2 h adiantado: K_AF válida tratada como vencida, `503 AANF_UNAVAILABLE`; 2 h atrasado: K_AF aceita, que ficaria em uso 2 h além do `expiry` |
| reautenticação | R3, B1, B2 | A-KID antigo: 403 no AAnF; a K_AF antiga continua aceita no AF até `expiry` (TS 33.535 §6.4.2, §5.2) |
| assinante sem AKMA | R5 | sem indicação do UDM, sem contexto no AAnF, `403 AKMA_KEY_NOT_AVAILABLE` no AF |
| AAnF indisponível na autenticação | R4 | registro aceito, erro só no log do AUSF (D9), `503 AANF_UNAVAILABLE` no AF |
| `remove-context` seguido de `retrieve` | R6, C1 | 204, depois 404 `AKMA_CONTEXT_NOT_FOUND`; `retrieve`: 403 |
| — | N24, N25 | Ua\* com K_AF errada: 401; `Host` de outro AF: 421 |

Os testes da API do AAnF (N01 a N16) estão em `af/scripts/negative/aanf_api.py`. São requisições HTTP/1.1 cruas, sobre TLS mútuo com o certificado `af1`, sem o código Go do projeto. Os do `register-anchorkey` precisam de *token* e usam o `seed` na máquina do free5GC, com as opções de teste `--gpsi` e `--k-akma`.

`[MVP]` Os casos de RID, relógio e Ua\* usam contextos semeados com os vetores `rid-012` e `mnc-3-digits` do MVP-0, de outros SUPI, para não tocar o contexto do UE real; o roteiro os remove depois.

`[MVP]` O relógio deslocado é simulado pela opção de teste `clockOffset` do AF (`af.md`). O AF não muda o relógio do sistema.

### 4.1 API do NEF e acesso anônimo (MVP-6)

`af/scripts/negative/nef_api.py` faz requisições HTTP/1.1 cruas sobre TLS mútuo, sem o código Go do projeto. Desde o MVP-11, ele entra no CCF como *invoker*, com as funções do `capif_api.py`, pede o *token* ao CCF e sai no fim (`docs/free5gc-side/tls-sbi.md` §9.3).

| Verificação | Caso | Esperado |
|---|---|---|
| E00 | os três AF externos entram no CCF e descobrem o NEF | 3 entradas |
| E01 | K_AF via NEF | 200 com `kaf` e `expiry`, sem `supi` e sem `gpsi` (TS 33.535 §6.3 passo 5) |
| E02 | a mesma K_AF via NEF e direto do AAnF, para o mesmo A-KID e AF_ID | iguais; nenhuma das duas é impressa |
| E03 | acesso anônimo via NEF | 200 sem `supi` e sem `gpsi` (§6.2.2; TS 29.522 tab. 5.14.5.3.3-1 NOTA 3) |
| E04 | sem *token* | 401 (*middleware* do free5GC, sem `ProblemDetails`: F12) |
| E05 | *token* do NRF para outro serviço (`naanf-akma`), pedido pelo `seed` na máquina do free5GC | 401 |
| E06 | AF_ID fora da política do NEF | 403 sem causa (§6.3 passo 2) |
| E07 | AF_ID aceito pelo NEF e pelo CCF e não pelo AAnF (`af5`, entrada de teste do `nefcfg.yaml` e do `ccfcfg.yaml`) | 403 sem causa, repassado (P11) |
| E08 | A-KID desconhecido | 403 `K_AKMA_NOT_PRESENT`, repassado (TS 29.522 §4.4.23.2); nunca 204 |
| E09, E10 | `afId` malformado; `aKId` ausente | 400 com `invalidParams` |
| E11 | `Content-Type` `text/plain` | 415 (TS 29.122 tab. 5.2.6-1) |
| E12 | os AF externos saem do CCF ao encerrar | 3 saídas |
| N01c | acesso anônimo direto ao AAnF (`anonInd`) | 200 sem `supi` |

### 4.2 CAPIF (MVP-9)

`af/scripts/negative/capif_api.py` é um segundo *invoker*, independente do código Go do projeto. Ele faz requisições HTTP/1.1 cruas sobre TLS mútuo, com o módulo `ssl` do Python, no CCF e no NEF, e usa o OpenSSL da linha de comando para o par de chaves. Os casos e o esperado estão em `docs/free5gc-side/capif.md` §10.2.

| Verificação | Caso | Esperado |
|---|---|---|
| CP00 | o AF do modo CAPIF sobe | entrada, descoberta, `OAUTH` e `check-authentication` no log do AF |
| CP01 | entrada do `capif_api.py` | `201` e `Location`; o certificado tem CN igual ao identificador e é da CA do CCF (`openssl verify`) |
| CP02 | descoberta | `apiRoot` `https://nef.5gc.mnc093.mcc208.3gppnetwork.org:8000`, `aefId` `nef-aef`, `[OAUTH]` |
| CP03, CP04 | método de segurança; `check-authentication` no NEF | `201` com `OAUTH`; `200` |
| CP05 | *token* | `Bearer` e `no-store`; `iss` = `client_id` = identificador; `scope`; `exp` = agora + `expires_in` |
| CP06, CP07 | K_AF via NEF com o *token* do CCF, normal e anônima | a K_AF do AAnF, sem `supi` (e sem `gpsi`) |
| CP08 | entrada sem credencial; com credencial adulterada | `401` com `WWW-Authenticate: Bearer` |
| CP09, CP10 | descoberta sem certificado; com o `api-invoker-id` do AF em Go | `403` |
| CP11 a CP13 | segredo errado; escopo fora do perfil; `grant_type=password` | `400` `invalid_client`, `invalid_scope`, `unsupported_grant_type` |
| CP14, CP15 | NEF sem *token*; com o *token* adulterado | `401` |
| CP16, CP17 | AF_ID `af3`, fora do perfil; AF_ID `af4`, fora da política do NEF | `403` sem `cause` |
| CP18 | saída; *token* pedido depois | `204`; `400` `invalid_client` |
| CP19 | saída do AF do modo CAPIF no fim (grupo B) | linha no log do AF |
| CP20 | logs do free5GC (CCF e NEF) e do AF | nenhum *token*, credencial, segredo de entrada ou chave privada |

### 4.3 TLS da SBI (MVP-11)

`af/scripts/negative/sbi_tls.py` usa só o módulo `ssl` do Python, sem o código Go do projeto. A máquina do free5GC o roda no passo R1s, em cada porta da SBI: NRF, AMF, SMF, UDR, PCF, UDM, NSSF, AUSF, CHF, NEF e AAnF. A do AF o roda no ponto A, no AAnF e no NEF. Os identificadores dizem a porta e a máquina: `S02.nrf.free5gc`, `S02.nef.af`. Especificação: `docs/free5gc-side/tls-sbi.md` §9.2 e §15.2.

| Verificação | Caso | Esperado |
|---|---|---|
| S01 | HTTP/1.1 e h2c em texto claro | recusados no TLS: o `400` do Go a um pedido HTTP/1.1, conexão fechada ao h2c |
| S02 | TLS sem certificado de cliente | alerta, sem resposta HTTP |
| S03 | certificado `other-ca` | alerta |
| S04 | o certificado do servidor, com o SNI do nome | o DNS-ID, o URI do NF Instance ID e o IP da entidade; onde o arquivo existe, o certificado instalado |
| S05 | `make-sba-certs.sh check` | nenhuma recusa |
| S06, S07 | certificado `revoked`; `no-crldp` | alerta |
| S08.name | o `seed` contra um `openssl s_server` com o certificado do UDM, no nome do AAnF | o cliente recusa pelo nome |
| S08.ca | o mesmo, com um certificado do nome do AAnF emitido por outra CA | o cliente recusa pela CA |
| S09 | o mesmo, com o certificado do AAnF | o `s_server` registra o SNI com o nome |
| S19 | arquivos `*sslkey*` escritos na execução; logs das duas máquinas | nenhum arquivo com conteúdo; nenhuma chave privada ou segredo do TLS |

## 5. Resultado em 18/10/2026

`63 passed, 0 failed`: 30 verificações na máquina do free5GC e 33 na do AF. Os cinco A-KID da execução, um por autenticação primária:

| Autenticação | A-KID | Situação |
|---|---|---|
| 1 | `rid0000.atidmq7LS8IJTQYlRh7lrP2XD1SpvrLD1uIU5AOoM7RJzYo=@…` | registrado; Ua\* aceito |
| 2 (nova autenticação) | `rid0000.atidAdq+4U+8HPV+Bw1n5F+EfWvjyoEqAE9rE6nMU8nP3Ic=@…` | substituiu o 1; o 1 passou a 403 no AAnF |
| 3 (AAnF parado) | `rid0000.atidqU2g296EbEun+yEPMArMbIpkKogSw5ZRmhaz8eTLNOI=@…` | derivado pelo AUSF; registro falhou (D9) |
| 4 (sem assinatura) | `rid0000.atidATIugkCuUJTyczuW/cDgcdldY0jOravQOMEMcmEUShk=@…` | derivado só pelo UE; desconhecido do AAnF |
| 5 (restaurada) | `rid0000.atidy6lXvgC9uQe8p/NNUb0B+rQtNOJTI2VLGBtqfV/bDeQ=@…` | registrado; removido no R6 |

No console do `run.sh` aparece uma cadeia de 64 hexadecimais: o nome do arquivo de CDR que o CHF do free5GC apaga ao encerrar. Não é material AKMA, e a verificação R7 não a alcança, porque copia o `free5gc.log` antes do encerramento.

### 5.1 Instalação limpa, em 19/10/2026 (critério 4)

`63 passed, 0 failed` (30 e 33), na execução `20261019-160723`. Ambiente:

- a VM 192.168.16.131, com o free5GC v4.2.3 e o UERANSIM instalados pelo guia oficial e sem nenhum *patch* do projeto;
- o tutorial seguido do início, partes 1 a 8, a partir de um clone limpo do repositório na máquina do AF;
- `woca@192.168.16.131` como argumento do roteiro.

A primeira passagem achou três problemas que impediam o roteiro fora da VM 192.168.16.130:

- o `e2e.sh` tinha o endereço do AAnF fixo;
- o `e2e-remote.sh` esperava o AAnF em `192.168.16.130:8000`;
- os dois chamavam o `seed` sem `--aanf`, e o padrão dele é o AAnF do laboratório.

Os endereços agora vêm das configurações (seção 3). As divergências do tutorial e a correção de cada uma estão no plano, MVP-5, critério 4.

Os cinco A-KID da execução:

| Autenticação | A-KID |
|---|---|
| 1 | `rid0000.atidvfVWKPf69ZWBMVQGrvOkLaDQFdGba3KyQ5U0gCgNJFg=@…` |
| 2 (nova autenticação) | `rid0000.atidEnHVYP6TrlaX6wXXUX9Y7QNGnhfI3zcLH/mb/s6INFA=@…` |
| 3 (AAnF parado) | `rid0000.atidSvpTFXFecorSiZI+PGTCEXoMHM9XMgnQAoi3LAJiebw=@…` |
| 4 (sem assinatura) | `rid0000.atidHGpmuka9FkKJw7ogzBXHnbejSfqB1XfXUHVpT4iHchc=@…` |
| 5 (restaurada) | `rid0000.atidSaWLZ7zwKidV8wciRplt7HcrhqFnt2YAV+Mu5eBlhs4=@…` |

### 5.2 MVP-6, em 19/10/2026

`83 passed, 0 failed` na execução `20261019-213720`: 36 verificações na máquina do free5GC e 47 na do AF. A configuração é a do MVP-6:

- AAnF em 192.168.16.130:8010;
- NEF corrigido em 192.168.16.130:8000;
- N2 e N3 em 192.168.16.130;
- gNB local com `ngapIp` 192.168.16.130 e `gtpIp` 10.100.0.1.

As 20 verificações novas são E00 a E12, N01c e R2n.ua-nef, R2n.ua-nef-anon, R2n.ua-anon, R2n.aanf-via-nef, R2n.nef-relayed e R2n.nef-no-supi. As 63 do MVP-5 passaram com a configuração nova.

Antes dela, duas execuções no mesmo dia:

- **`20261019-203849` (63/0):** a regressão do MVP-5 com a configuração nova.
- **`20261019-213431` (80/3):** as três falhas eram do roteiro, não do sistema. A função que gera as configurações dos AF extras falhava com `set -u` quando chamada com três argumentos, e o AF externo da porta 8083 não subiu. Corrigido.

O teste com o UERANSIM em outra máquina (critério 6) está em `docs/af-side/teste-nef.md`. O roteiro supõe o gNB e o UE na máquina do free5GC.

### 5.3 MVP-8, fase (a), em 20/10/2026

`88 passed, 0 failed` na execução `20261020-215257`: 41 verificações na máquina do free5GC e 47 na do AF. Mesma montagem do MVP-6, com o UERANSIM nos *patches* 0001 a 0005 (serviço de chaves do ME).

- **Verificações novas.** As cinco do R2k: K_AF ao `uaclient` para o seu AF_ID; outro AF_ID recusado; outro usuário (`nobody`) recusado; `akma-kaf` desligado; configuração com um AF_ID em duas aplicações recusada.
- **Execução anterior, `20261020-203951` (87/1).** A falha foi do roteiro: o R3.ua-old guardava a K_AF anterior pelo `akma-kaf`, agora desligado. O roteiro passou a pedir a K_AF ao serviço de chaves, como o usuário do `uaclient`.
- **Conferência à mão do executável** (em 21/10/2026, na VM do free5GC, com um `nr-ue` avulso e uma política com `executable`):
  - o `uaclient` passou pela política e recebeu `NO_KEY_MATERIAL`, porque não havia registro;
  - o `python3` do mesmo usuário recebeu `ERR NOT_AUTHORIZED`;
  - ao encerrar, o `nr-ue` apagou o *socket*.

### 5.4 MVP-8, fase (b), em 21/10/2026

`115 passed, 0 failed` na execução `20261021-164401`: 51 verificações na máquina do free5GC e 64 na do AF, em cerca de três minutos. O AF e o `uaclient` foram compilados com o wolfSSL 5.9.4. A máquina do free5GC tem a política nova do AAnF, os *patches* de configuração com a política nova do NEF e o *patch* 0005 com a política nova do UE.

- **Verificações novas.** As oito do R2t, as dezessete do grupo T (seção 1) e duas no R2n: um AF fora da rede do operador com o perfil PSK TLS 1.3 (porta 8456), cuja K_AF vem pelo NEF. Os casos e o resultado de cada um estão em `ua-star-norma.md` §3.9.
- **Execuções anteriores.** As falhas foram do roteiro, não do sistema:
  - **`20261021-115934` (106/7).** O `uaclient` da máquina do free5GC era de antes da opção `-tls12` (R2t.k2-tls12, R2t.me, T02). O R2k.af-id pedia ao cliente um AF_ID de outro FQDN; o cliente agora só pergunta ao ME depois do desafio, e o AF recusa o FQDN antes (421). O R2k.af-id passou a perguntar ao serviço de chaves direto, como o usuário do `uaclient`. T12, T12b e T13b: o AAnF devolvia a K_AF que ainda valia para ele, e o AF, com o relógio adiantado, a achava vencida (`ua-star-norma.md` §3.9).
  - **`20261021-120446`.** O registro de novo pelo `seed`, sem remover o contexto, não renova a K_AF no AAnF: o AAnF guarda as K_AF de um contexto enquanto o A-KID não muda. O roteiro passou a remover o contexto antes de registrar de novo. Além disso, o `e2e.sh` foi editado enquanto rodava e voltou a executar a parte remota; ela foi interrompida com Ctrl-C.
  - **`20261021-132202` (113/0).** Execução válida, antes das duas verificações do AF externo com PSK no R2n.
  - **`20261021-120909` (66/36).** Não vale. A parte remota interrompida na execução anterior não parou: o `e2e-remote.sh` ligava o controle de *jobs* (`set -m`), e o Ctrl-C e o fim da sessão SSH só alcançavam o comando em curso. Ela seguiu sozinha até o fim e parou o free5GC desta execução no meio do R2. Agora o controle de *jobs* só vale para iniciar o `run.sh`, e o Ctrl-C para o roteiro remoto.
- **AUSF.** Antes da execução, o AUSF da máquina do free5GC foi revertido com o texto anterior do *patch* e reaplicado com o atual (`124b2d5`, que só mudou comentários). Assim o *preflight* volta a mostrar `ausf: patches applied`.

### 5.5 MVP-8, fase (c), em 22/10/2026

`126 passed, 0 failed, 0 skipped` na execução `20261022-124514`: 56 verificações na máquina do free5GC e 70 na do AF, em cerca de dois minutos e meio. O AF e o `uaclient` foram compilados com o wolfSSL 5.9.4 com DTLS; o D03 usou o OpenSSL 4.1.0 *beta* 2, compilado à parte (`OPENSSL_DTLS13`). A máquina do free5GC tem a política do AAnF com os nove AF_IDs, os *patches* de configuração com a política nova do NEF e o *patch* 0005 com o AF_ID do DTLS 1.3 na política do UE.

- **Verificações novas.** As cinco do R2d e as seis do grupo D (seção 1). Os casos e o resultado de cada um estão em `ua-star-norma.md` §4.8.
- **Execuções anteriores.**
  - **`20261021-183854` (124/2).** T11: o texto do log do AF no vencimento da K_AF passou a dizer o protocolo (`TLS 1.3 session`), e o roteiro ainda procurava o texto antigo. D03: o segundo `ClientHello` do OpenSSL 4.1, com 1702 octetos, passava do que o wolfSSL lê por padrão com o MTU de 1200 (`ua-star-norma.md` §4.7). O `build-wolfssl.sh` passou a compilar a biblioteca com uma leitura maior.
  - **`20261022-120810` (125/1).** D03 de novo: o AF ainda tinha a biblioteca antiga, porque o Go reaproveitou o binário já ligado (`ua-star-norma.md` §4.7). O `build-af.sh` passou a compilar o *hash* da biblioteca dentro de `af/wolftls`, e o `uaclient` da máquina do free5GC foi trocado pelo compilado de novo.

### 5.6 MVP-8, fase (d), em 22/10/2026

`141 passed, 0 failed, 0 skipped` na execução `20261022-201824`: 60 verificações na máquina do free5GC e 81 na do AF, em cerca de três minutos. Com o `OPENSSL_DTLS13` da seção 5.5 e o `coap-client` da libcoap 4.3.5b em `LIBCOAP_CLIENT` (seção 3). Antes da execução, a máquina do free5GC passou a ter o AF_ID `…0100000008` na política do AAnF, na do NEF (*patch* de configuração `0002`) e na do UE (*patch* `0005`); as cópias anteriores ficaram em `~/akma-mvp8c-*` e `aanfcfg.yaml.mvp8c`.

- **Verificações novas.** As quatro do R2o e as onze do grupo O (seção 1). Os casos e o resultado de cada um estão em `ua-star-norma.md` §7.12.
- **Execução anterior, `20261022-201455` (140/1).** O O10 falhou por causa do roteiro: a saída do cliente tem duas linhas, e o `grep` da verificação procura linha por linha. O AF respondeu o esperado, `4.03 KAF_EXPIRED`. A verificação passou a juntar as linhas.

### 5.7 MVP-12, em 28/10/2026

`161 passed, 0 failed, 0 skipped` na execução `20261028-162316`: 80 verificações na máquina do free5GC e 81 na do AF, em cerca de quatro minutos. Com o `OPENSSL_DTLS13` e o `LIBCOAP_CLIENT` da seção 3, recompilados nesta data pelos passos daquela seção. Antes da execução, a máquina do free5GC recebeu as duas cópias do `openapi` com o *patch* `0003`, o UDM com o `0002`, o AMF com o `amf/0001`, o AAnF e o `seed` novos e o `aanfcfg.yaml` com a linha `kakmaRefresh: 0`. As cópias anteriores ficaram em `~/akma-mvp12-*` e `aanfcfg.yaml.mvp8d`.

- **Verificações novas.** As vinte do R6r (seção 1). Os casos e o resultado de cada um estão em `docs/free5gc-side/renovacao-kakma.md` §11.
- **Execução `20261028-151849` (151/10).** As dez falhas, no R6r: o roteiro pedia, depois de 20 s, a K_AF do mesmo AF_ID, e o AF, que a tinha guardada, não consultou o AAnF; o `seed trigger-auth` repetia o prefixo da API no caminho. Corrigidos o roteiro e o `seed`.
- **Execução `20261028-154345` (160/1).** O R6r.smc falhou: o UERANSIM pôs o último `Registration Request` no `Security Mode Complete`, e o AMF o tratou como um registro. Corrigido no AMF (`docs/free5gc-side/amf.md` §2).

### 5.8 MVP-13, instalação limpa da V2, em 29/10/2026

`161 passed, 0 failed, 0 skipped` na execução `20261029-124631`: 80 verificações na máquina do free5GC e 81 na do AF. A máquina do free5GC é o *snapshot* limpo do free5GC v4.2.3 e do UERANSIM `85a0fbf` (o da .132, que voltou como 192.168.16.133). O repositório veio de um `git bundle` do *commit* `03c0606` do repositório público, instalado só pelo `docs/install.md` público, partes 1 a 7, 9.1 e 11.1; o AF, nesta máquina, de um clone do mesmo *bundle*. Os endereços mudaram nos lugares da tabela 1.2 do guia, e por isso o relatório diz `03c0606-dirty`. O `OPENSSL_DTLS13` e o `LIBCOAP_CLIENT` foram recompilados pelos passos de `public/docs/testing.md` §5.1. O relatório ficou em `af/build/e2e/20261029-124631-clean-v2`, fora do controle de versão.

- **Lacuna do guia.** A parte 8 e `testing.md` §5.1 pediam as partes 3 a 7 e a seção 9.1, sem a 11.1 (o *patch* do AMF), sem a qual o R6r falha, e sem a linha `kakmaRefresh: 0` no `aanfcfg.yaml` instalado, que o R6r troca numa cópia. Corrigido no guia público antes da execução.

### 5.9 MVP-9, em 02/11/2026

`185 passed, 0 failed, 3 skipped` na execução `20261102-093804`: 86 verificações na máquina do free5GC e 99 na do AF. As três puladas são D03, O04 e O05: o `OPENSSL_DTLS13` e o `LIBCOAP_CLIENT` da seção 3 não estavam compilados nesta máquina. Nada no MVP-9 depende delas.

- **Verificações novas.** São 27: R1.ccf, R1.nef-published e R1.credential; R2c.ua-capif e R2c.nef-capif; CP00 a CP19; CP20.af e CP20.free5gc.
- **Antes da execução, na máquina do free5GC.** Os conjuntos anteriores de *patches* do NEF, da configuração e do `run.sh` foram guardados em `~/akma-mvp9-old-patches` e revertidos. O repositório foi sincronizado. As duas cópias do `openapi` foram refeitas com o `0004`. Foram instalados o CCF (`bin/ccf`, `config/ccfcfg.yaml` e os certificados) e os conjuntos novos: `nef/0001` e `0002`, `config/0001` a `0003` e `run/0001` e `0002`. O `bin/nef.orig` continua lá.
- **Dois *scripts* corrigidos na implantação.**
  - `patch-nf.sh` conferia cada *patch* sozinho, contra a árvore original. Os *patches* novos dependem dos anteriores e eram recusados. Agora a série é conferida em ordem, numa cópia dos arquivos que ela toca.
  - `make-capif-certs.sh`: com o OpenSSL 1.1.1 da máquina do free5GC, o `-addext` *acrescentava* o `basicConstraints` ao da seção `v3_ca`, e o certificado da CA saía com a extensão duas vezes, o que o Go recusa. Agora a CA é feita com uma configuração própria.

### 5.10 MVP-11, fase (a), em 04/11/2026

`288 passed, 0 failed, 0 skipped` na execução `20261104-190824`, em quatro minutos: 171 verificações na máquina do free5GC e 117 na do AF. É o primeiro `e2e.sh` com TLS mútuo em toda a SBI.

- **Verificações novas.** São 100: R1.crl-signed e R1.crl; S01 a S09 e S19, 83 na máquina do free5GC (onze portas, sete conferências em cada uma, mais S05, S08, S09 e S19) e 15 na do AF (AAnF e NEF, mais S19.af).
- **As três puladas voltaram.** D03, O04 e O05 passaram, com o OpenSSL 4.1.0 *beta* 2 e a libcoap 4.3.5b em `~/openssl41` e `~/libcoap`, fora do repositório (seção 3).
- **Grupo do NEF.** Os três AF externos do MVP-6 são *invokers* do CCF. O E07 é recusado pelo AAnF, como o log mostra: `NEF at 192.168.16.130 not authorized for AF_ID af5.akma.example.FF00000001`.
- **Antes da execução, na máquina do free5GC.** `config/` e `cert/` ficaram em `~/akma-mvp11-backup`, e os *patches* antigos em `~/akma-mvp11-old-patches`, depois do *snapshot* da VM. Os conjuntos antigos foram revertidos, o repositório sincronizado, as duas cópias do `openapi` e a do `util` refeitas, as onze NFs compiladas com os *patches* novos e os certificados feitos com o OpenSSL 1.1.1. O `/etc/hosts` das duas máquinas ganhou os nomes.
- **O banco de testes da máquina estava defasado.** O `akma_test.go` instalado era anterior a uma revisão dos comentários do `test/0001`, e o `patch-nf.sh` o dava como não aplicado. Os outros arquivos foram revertidos com o *patch*, e o `akma_test.go` foi apagado, com cópia no *backup*, antes de aplicar a série.
- **Um defeito do teste corrigido.** O `tamper` do `capif_api.py` trocava o último caractere da assinatura. Numa assinatura `ES256`, em *base64url*, os quatro bits baixos desse caractere são enchimento, e o decodificador os ignora. Com assinaturas terminadas em `A`, a credencial "adulterada" era a mesma, e o CP08 falhava. A chance é de uma em quatro. O MVP-9 passou por sorte. Agora muda o caractere anterior.
- **Compatibilidade com a máquina do free5GC.** Ela tem o Python 3.6.9 e o OpenSSL 1.1.1. O `sbi_tls.py` trata a recusa do cliente nos dois jeitos, porque `ssl.SSLCertVerificationError` só existe a partir do Python 3.7.

## 6. Limitações

- **Senha de `sudo`.** É pedida uma vez por execução; o roteiro não altera o `sudoers`.
- **Endereços.** O roteiro não lê a máquina do free5GC das configurações. Ela é o argumento, e o padrão é a do laboratório.
- **Ordem fixa.** Uma falha no meio não interrompe o roteiro: as verificações seguintes rodam no estado que houver e podem falhar em cascata. A primeira `FAIL` do relatório é a que importa.
