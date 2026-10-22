# Ua\* da norma e autorização por aplicação no UE: especificação

**MVP:** MVP-8 · **Data:** 20 de outubro de 2026 · **Estado:** fases (a) e (b) aprovadas pelo responsável em 20/10/2026, com os pontos U1 a U12 como propostos. A fase (a) está implementada e verificada (20/10/2026). A fase (b) está implementada e verificada (21/10/2026): a implementação e as diferenças em relação a este texto estão na seção 3.8, e a verificação, com 115 verificações e nenhuma falha, na seção 3.9. A fase (c), DTLS, foi aprovada pelo responsável em 21/10/2026, com os pontos U13 a U21 como propostos, e está implementada e verificada (22/10/2026): a implementação está na seção 4.7, e a verificação, com 126 verificações e nenhuma falha, na seção 4.8. A fase (d), OSCORE, entra depois (seção 7).

Decisões já tomadas pelo responsável em 20/10/2026 (plano, MVP-8):

- um MVP em quatro fases: (a) serviço de chaves do ME; (b) TLS; (c) DTLS 1.3; (d) OSCORE;
- o `akma-kaf` do `nr-cli` fica atrás de uma opção de teste, desligada por padrão;
- o Ua\* simplificado (`FF00000001`, `ua-star.md`) fica como perfil de teste, desligado por padrão;
- o PSK TLS e o DTLS usam o wolfSSL, e os componentes que o usam ficam sob GPLv3, como os *patches* do UERANSIM.

## 1. Situação normativa

### 1.1 O UE

`[3GPP]` TS 33.535 §4.4.3, *"Requirements on the UE"*:

1. *"Applications on the UE shall not be able to get access to KAKMA."*
2. *"An application on the UE shall only get the KAF keys related to specific AF Identifiers (AF_IDs) that the application is authorized to get."*
3. *"An application on the UE shall not be able to get access to the KAF keys that belong to other applications."*
4. NOTA: *"How these requirements are satisfied is out of scope of 3GPP."*

`[FREE5GC]` O UERANSIM com os *patches* do MVP-4 atende o item 1: nenhum comando devolve a K_AKMA. Não atende os itens 2 e 3:

- o comando `akma-kaf` deriva e devolve a K_AF de qualquer FQDN, com qualquer identificador Ua\* (`src/ue/app/cmd_handler.cpp`);
- o `nr-cli` fala com o `nr-ue` por UDP em 127.0.0.1, numa porta publicada em `/tmp/UERANSIM.proc-table/` (`src/lib/app/cli_base.hpp`, `src/utils/constants.hpp`), sem autenticar quem pergunta.

Qualquer processo no *namespace* do UE obtém qualquer K_AF.

`[LACUNA]` A NOTA deixa a forma para a implementação. O GBA trata do mesmo problema só de modo informativo (TS 33.222 Anexo C, controle das credenciais GBA no UE).

### 1.2 Os perfis TLS

`[3GPP]` TS 33.535 Anexo B adapta ao AKMA os perfis de TS 33.222, com o AF no papel do NAF. O Stage 3 está em TS 24.109 Anexo K. Os identificadores Ua\* vêm de TS 33.220 H.3.

| Perfil | Procedimento | Identificador Ua\* | Como a K_AF entra |
|---|---|---|---|
| **K.2**: autenticação do UE por chave compartilhada e do AF por certificado | TS 33.535 B.1.2.2; TS 33.222 §5.3.0; TS 24.109 K.2 | `01 00 01 yy zz`, com `yy zz` o *cipher suite* do TLS (TS 33.220 H.3) | HTTP Digest dentro do TLS: `username` = A-KID; senha = base64(K_AF) |
| **K.3, TLS 1.2**: autenticação mútua por chave compartilhada | TS 33.535 B.1.3.2.1; TS 33.222 §5.4.0.1; TS 24.109 K.3.2 | `01 00 01 yy zz` | *hint* `3GPP-AKMA`; identidade `3GPP-AKMA` + A-KID; *premaster secret* derivado da K_AF |
| **K.3, TLS 1.3** | TS 33.535 B.1.3.2.2; TS 33.222 §5.4.0.2; TS 24.109 K.3.3 | `01 00 00 00 07` | identidade `3GPP-AKMA;<A-KID>`; PSK externa derivada da K_AF |

`[3GPP]` TS 33.222 §5.4.0.2, NOTA 0: no TLS 1.3 o UE precisa do identificador Ua\* antes do `ClientHello`, por isso ele é fixo. No TLS 1.2 e no K.2 o identificador depende do *cipher suite*. O UE só o conhece depois do `ServerHello`, e é então que deriva a K_AF.

`[3GPP]` Perfil TLS. TS 33.222 §5.3.1 e §5.4.1 e TS 24.109 K.2 e K.3.1 remetem ao TS 33.310 Anexo E. Na V19.5.0, esse anexo só diz *"The TLS protocol profiles are located in TS 33.210"*. TS 33.210 V19.3.0 §6.2 manda:

- TLS 1.2 e TLS 1.3 *"shall be supported"*; SSL, TLS 1.0 e TLS 1.1 não;
- TLS 1.3: *"psk_ke shall not be supported"*; o PSK vai sempre com ECDHE (`psk_dhe_ke`);
- TLS 1.2: só suítes AEAD com ECDHE. Obrigatórias: `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256` e `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`. Recomendadas: as duas `AES_256_GCM_SHA384`;
- TLS 1.2 com PSK: RFC 5489. `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256` (RFC 8442, `D0 01`) é obrigatória, e `TLS_ECDHE_PSK_WITH_AES_256_GCM_SHA384` (`D0 02`) recomendada;
- TLS 1.2: SNI, *Extended Master Secret* (RFC 7627) e renegociação segura (RFC 5746) *"shall be supported"*.

`[3GPP]` TS 33.222 §5.3.1.5 e §5.4.1.4: no TLS 1.2, cliente e servidor *"shall allow for resuming a session"*. No K.2, a sessão vale no máximo 24 horas (§5.3.1.5). No PSK, vale a vida do segredo compartilhado, no máximo 24 horas (§5.4.1.4).

`[3GPP]` TS 33.222 §5.3.0 passo 1 e §5.3.1.3: no K.2, o UE confere que o certificado do AF corresponde ao FQDN do AF (RFC 2818 §3.1) e não manda certificado de cliente.

### 1.3 A pilha

`[FREE5GC]` O `crypto/tls` do Go faz o K.2 (TLS 1.2 e 1.3 com certificado), mas não faz PSK nem DTLS.

Verificado em 20/10/2026:

- o OpenSSL 3.0.13 da VM do AF não tem `ECDHE-PSK-AES128-GCM-SHA256`;
- o OpenSSL só traz DTLS 1.3 a partir da 4.1;
- o wolfSSL tem TLS 1.2 e 1.3 com PSK e DTLS 1.3, sob GPLv3 ou licença comercial. Conferido compilando em 21/10/2026 (seção 3.8): a suíte `ECDHE-PSK-AES128-GCM-SHA256` (`D0 01`) está lá; a `ECDHE-PSK-AES256-GCM-SHA384` (`D0 02`) não existe no wolfSSL;
- o wolfSSL 5.8.4-stable (20/11/2025) corrige o rebaixamento de PSK com PFS para PSK sem PFS no TLS 1.3 (CVE-2025-11935) e a verificação do *binder* fora de tempo constante (CVE-2025-11932);
- não há ligação Go do wolfSSL que o projeto possa adotar com segurança.

## 2. Fase (a): serviço de chaves AKMA do ME

### 2.1 O que muda

| Parte | Mudança |
|---|---|
| UERANSIM, *patch* novo `0005-akma-key-service.patch` | o `nr-ue` atende pedidos de K_AF num *socket* Unix e aplica a política por aplicação |
| configuração do UE | `akmaKeyService` (política) e `akmaCliKaf` (opção de teste) |
| `af/ue/uaclient` | obtém o A-KID e a K_AF pelo *socket* (`-me-socket`); o caminho pelo `nr-cli` só com a opção de teste |
| `af/scripts/ue-netns.sh` | o subcomando `ua` passa o *socket* ao cliente |

### 2.2 Identidade da aplicação

`[INTERPRETAÇÃO]` Fica a cargo do projeto, pela NOTA de §4.4.3.

- **O *uid* é a fronteira.** O ME lê do kernel as credenciais de quem conectou (`SO_PEERCRED`: `pid`, `uid`, `gid`). No Linux, a separação entre processos é por usuário. Um processo do mesmo usuário pode, por exemplo, trocar o executável da aplicação ou injetar código nela ao iniciá-la (`LD_PRELOAD`). O Yama do Ubuntu (`ptrace_scope` 1) restringe o `ptrace`, mas não fecha esses caminhos. Então o usuário é o que separa uma aplicação da outra de verdade.
- **O executável é conferência adicional.** O ME também confere o executável do processo (`/proc/<pid>/exe`) contra a política. É uma conferência de melhor esforço: um processo pode conectar e depois trocar de imagem com `exec`. Fica declarada como tal.
- **Aplicações separadas por usuário.** Para separar de verdade duas aplicações, cada uma roda com um usuário próprio, como nos sistemas móveis.

### 2.3 Política

Na configuração do UE (`free5gc-ue.yaml`):

```yaml
akma: true
akmaKeyService:
  socket: /run/ueransim/akma-imsi-208930000000001.sock   # opcional; padrão /run/ueransim/akma-<SUPI>.sock
  applications:
    - name: uaclient
      uid: 1000
      executable: /home/woca/akma-free5gc/af/build/bin/uaclient   # opcional
      afIds:                                                       # forma de transporte
        - af1.akma.example.0100000007
        - af1.akma.example.010001D001
        - af1.akma.example.0100011301
akmaCliKaf: false                      # teste: o comando akma-kaf do nr-cli
```

Regras:

- **Item 2 de §4.4.3.** A aplicação só recebe a K_AF dos AF_IDs da sua lista. Pedido de outro AF_ID: recusa e log `warn`, com o nome da aplicação e o AF_ID, sem chave.
- **Item 3 de §4.4.3.** A K_AF depende só do AF_ID, então duas aplicações com o mesmo AF_ID teriam a mesma chave. A validação da configuração recusa um AF_ID que apareça em duas aplicações, e o `nr-ue` não sobe.
- **Item 1 de §4.4.3.** Nenhuma resposta leva a K_AKMA.
- **Sem política.** Sem `akmaKeyService.applications`, nenhuma aplicação recebe K_AF, e o `nr-ue` registra isso na partida.

### 2.4 Protocolo do *socket*

Texto ASCII, uma linha por pedido e uma por resposta, vários pedidos por conexão:

```
KAF <FQDN do AF> <identificador Ua* em 10 hexadecimais>
OK a-kid=<A-KID> af-id=<AF_ID> k-af=<64 hexadecimais>
ERR <causa>
```

| Causa | Quando |
|---|---|
| `AKMA_DISABLED` | `akma: false` |
| `NO_KEY_MATERIAL` | ainda não houve autenticação primária bem-sucedida |
| `BAD_REQUEST` | FQDN que não é nome DNS ASCII (DV-4), ou identificador que não tem 10 hexadecimais |
| `NOT_AUTHORIZED` | *uid* ou executável fora da política, ou AF_ID fora da lista da aplicação |
| `UNAVAILABLE` | o ME não conseguiu pausar as tarefas a tempo para ler o contexto de segurança NAS, como acontece com o `nr-cli` |

O *socket* fica em caminho do sistema de arquivos, então vale entre *namespaces* de rede: o cliente no *namespace* do UE o alcança. Permissão `0666`: a autorização é a da política, não a do arquivo.

### 2.5 O `akma-kaf` como teste

Com `akmaCliKaf: false`, o padrão, o comando responde *"akma-kaf is disabled; use the AKMA key service"*. Com `true`, ele funciona como hoje, e o `nr-ue` avisa na partida que a K_AF está exposta a qualquer processo local.

### 2.6 Verificação prevista

| Caso | Esperado |
|---|---|
| `uaclient` do usuário e do executável da política, AF_ID da lista | `OK`; Ua\* com o AF passa |
| a mesma aplicação, AF_ID fora da lista | `ERR NOT_AUTHORIZED`; `warn` no log do UE |
| outro usuário (`sudo -u nobody`) | `ERR NOT_AUTHORIZED` |
| outro executável do mesmo usuário | `ERR NOT_AUTHORIZED` (melhor esforço, seção 2.2) |
| AF_ID em duas aplicações | o `nr-ue` não sobe |
| `akma-kaf` sem `akmaCliKaf` | erro |
| todas as respostas e logs | nenhuma K_AKMA; nenhuma K_AF em log |

## 3. Fase (b): TLS

### 3.1 Componentes e licença

| Componente | Licença | Conteúdo |
|---|---|---|
| `af/akmatls` | Apache-2.0 | identificadores Ua\* por perfil, identidade PSK (formato e leitura), K.2 com o `crypto/tls` |
| `af/wolftls` | **GPLv3** | ligação cgo própria ao wolfSSL: `net.Conn` sobre uma sessão TLS 1.2/1.3 com PSK, *callbacks* de PSK do servidor e do cliente |
| `af/scripts/build-wolfssl.sh` | GPLv3 | baixa a versão fixada do wolfSSL (5.8.4-stable ou mais nova), confere o SHA-256 e compila a biblioteca estática |
| AF e `uaclient` | Apache-2.0 sem a *tag* `wolfssl`; **GPLv3** com ela | sem a *tag*: Ua\* simplificado e K.2; com ela: também PSK TLS 1.2 e 1.3 |

O `COPYING` da GPLv3 vai em `af/wolftls/`, e o `NOTICE` do repositório público diz que os binários compilados com a *tag* `wolfssl` são GPLv3.

### 3.2 K.2: Digest dentro de TLS, AF com certificado

`[3GPP]` TS 33.535 B.1.2.2 sobre TS 33.222 §5.3.0:

1. O UE abre TLS com o AF, confere o certificado contra o FQDN do AF e não manda certificado.
2. O pedido HTTP vai dentro do TLS, com `3gpp-akma` no `User-Agent`.
3. O AF responde 401 com o desafio Digest e o `realm` AKMA.
4. O UE confere que o FQDN do `realm` é o do AF com quem abriu o TLS. Se não for, encerra o TLS.
5. O UE responde com `username` = A-KID e senha base64(K_AF).
6. O AF confere com a K_AF pedida ao AAnF. Se não conseguir a chave, responde com erro sem o `realm`.

O protocolo HTTP e as conferências são os do `ua-star.md` (seções 5 a 9), com três mudanças:

- **Transporte.** HTTPS em vez de HTTP.
- **Identificador.** Passa a ser `01 00 01 yy zz`, com o *cipher suite* da conexão. O AF o lê da conexão TLS (`r.TLS.CipherSuite`), e o UE o lê da resposta (`resp.TLS.CipherSuite`). O AF_ID enviado ao AAnF é `<FQDN>.010001YYZZ`, por exemplo `af1.akma.example.0100011301` para `TLS_AES_128_GCM_SHA256`.
- **Certificado.** O do AF é emitido por uma CA de teste, gerada por *script*. O UE recebe o certificado da CA por configuração (`-ca`).

O Digest continua o de RFC 7616 (DV-5), e o `realm` continua `3GPP-bootstrapping-akma@<FQDN>` (DV-6).

`[3GPP]` *Cipher suites* (TS 33.210 §6.2):

- TLS 1.3: os do `crypto/tls`, `13 01`, `13 02` e `13 03`;
- TLS 1.2: `C0 2B` e `C0 2F` (obrigatórios) e `C0 2C` e `C0 30` (recomendados).

Cada *cipher suite* aceito dá um AF_ID, e cada AF_ID precisa estar na política do AAnF (ou na do NEF, para o AF externo) e na do UE.

### 3.3 K.3 com TLS 1.3

`[3GPP]` TS 24.109 K.3.3 e TS 33.222 §5.4.0.2:

- **UE:**
  - deriva a K_AF com o identificador `01 00 00 00 07` antes do `ClientHello`, pelo serviço de chaves da fase (a);
  - no `ClientHello`: a extensão `psk_key_exchange_modes` só com `psk_dhe_ke` (TS 33.210 §6.2.2), o `server_name` com o FQDN do AF e a identidade `3GPP-AKMA;<A-KID>` em `psk_identities`. No wolfSSL, o modo é forçado com `wolfSSL_CTX_only_dhe_psk()`, como no repositório `adaptive-security` do responsável; a conferir na compilação;
  - sem 0-RTT (`early_data`): a literatura do projeto (artigo 13, Ko et al., 2024) mostra que o 0-RTT não tem sigilo futuro e admite repetição;
  - oferece também autenticação por certificado (TS 24.109 K.3.3 1)c)).
- **AF:**
  - no *callback* de PSK, lê a identidade e confere o prefixo `3GPP-AKMA`;
  - pede a K_AF ao AAnF (ou ao NEF) com o AF_ID `<FQDN>.0100000007`, e a K_AF vira a PSK externa;
  - se não conseguir a K_AF, recusa a identidade, e o *handshake* falha.
- **A PSK externa:** ver o ponto U1 da seção 5. O *hash* associado é o SHA-256: RFC 8446 §4.2.11, *"default to SHA-256 if no such algorithm is defined"*.

Depois do *handshake*, o UE faz `GET /ua/v1/session` sem Digest: o TLS já autenticou os dois lados. A resposta 200 é a do `ua-star.md` §4.

### 3.4 K.3 com TLS 1.2

`[3GPP]` TS 24.109 K.3.2 e TS 33.222 §5.4.0.1, com RFC 5489:

- o `ClientHello` oferece `ECDHE-PSK-AES128-GCM-SHA256` e `ECDHE-PSK-AES256-GCM-SHA384`, além de suítes sem PSK (TS 33.222 §5.4.0.1 passo 1), e leva o `server_name`;
- o AF escolhe uma suíte PSK e manda `ServerKeyExchange` com o *hint* `3GPP-AKMA`;
- o UE deriva a K_AF com o identificador `01 00 01` mais o código da suíte escolhida (`D0 01` ou `D0 02`) e manda `ClientKeyExchange` com a identidade `3GPP-AKMA;<A-KID>` (ponto U2);
- os dois formam o *premaster secret* da RFC 5489 §2, com a K_AF como PSK (ponto U1). O *Extended Master Secret* vai ligado.

### 3.5 Portas do AF

O AF passa a ter uma lista de perfis, cada um na sua porta:

```yaml
configuration:
  profiles:
    - type: digest-tls           # K.2
      listen: 192.168.16.129:8443
      certificate: cert/af1.akma.example.pem
      key: cert/af1.akma.example.key
    - type: psk-tls13            # K.3, TLS 1.3 (precisa da tag wolfssl)
      listen: 192.168.16.129:8444
    - type: psk-tls12            # K.3, TLS 1.2 (precisa da tag wolfssl)
      listen: 192.168.16.129:8445
    - type: digest-http          # Ua* simplificado, teste; desligado por padrão
      listen: 192.168.16.129:8080
```

Sem `profiles`, vale o `listen` de hoje, com o Ua\* simplificado, para que as configurações atuais continuem funcionando.

### 3.6 Expiração da K_AF e retomada de sessão

`[3GPP]` TS 33.535 §6.4.2 e §4.4.0 (tempo máximo de vida); TS 33.222 §5.4.1.4.

- **Sessões abertas.** O AF encerra (`close_notify`) toda sessão TLS cuja K_AF venceu.
- **Retomada.** É suportada:
  - nos perfis PSK, a validade do *ticket* ou da sessão vai até o menor entre a expiração da K_AF e 24 horas (TS 33.222 §5.4.1.4);
  - no K.2, até 24 horas (§5.3.1.5). A K_AF não entra no TLS do K.2: o Digest de cada pedido é que a usa, com as regras do `ua-star.md` §7.
- **Política `refetch`.** No *handshake* seguinte, o AF pede a K_AF de novo ao AAnF.
- **Política `reject`.**
  - Nos perfis PSK, o AF recusa a identidade e o *handshake* falha. O TLS não tem como levar a causa `KAF_EXPIRED`, que o §6.4.2 deixa opcional (*"may provide a cause"*).
  - No K.2, a resposta é a do `ua-star.md` §7.

### 3.7 Verificação prevista

| Caso | Esperado |
|---|---|
| UE → AF por K.2, TLS 1.3 e TLS 1.2 | 200 com `rspauth`; AF_ID `…0100011301` (ou o da suíte) no log do AAnF |
| K.2 com certificado de outro FQDN | o UE encerra antes de mandar o A-KID |
| K.2 com `realm` de outro FQDN | o UE encerra o TLS (B.1.2.2 passo 4) |
| UE → AF por PSK TLS 1.3 | *handshake* com `psk_dhe_ke`; 200; AF_ID `…0100000007` |
| UE → AF por PSK TLS 1.2 | suíte `D0 01`; *hint* `3GPP-AKMA`; 200; AF_ID `…010001D001` |
| cliente que só oferece `psk_ke` | *handshake* recusado (TS 33.210 §6.2.2) |
| PSK com K_AF errada | *handshake* falha no `Finished` ou no *binder* |
| A-KID desconhecido do AAnF | *handshake* recusado; 403 `K_AKMA_NOT_PRESENT` no log do AF |
| AF_ID de um perfil fora da política do AAnF | *handshake* recusado |
| conferência independente | `openssl s_client -tls1_3 -psk_identity "3GPP-AKMA;<A-KID>" -psk <K_AF>` (OpenSSL 3.0) contra o AF, com a K_AF de um vetor |
| K_AF vencida, `refetch` e `reject` | sessão encerrada; novo *handshake* conforme a política |
| logs do AF, do UE e do wolfSSL | nenhuma chave |

O `e2e.sh` ganha um grupo T, com o UE no seu *namespace* passando pelos três perfis TLS, e os negativos acima.

### 3.8 Implementação (21/10/2026)

**Componentes.**

| Componente | Conteúdo |
|---|---|
| `af/scripts/build-wolfssl.sh` (GPLv3) | wolfSSL 5.9.4-stable (25/09/2026), `wolfssl-5.9.4.zip` de wolfssl.com com SHA-256 fixado; os fontes comuns ao *tag* `v5.9.4-stable` do GitHub foram comparados e são iguais. Biblioteca estática em `af/build/wolfssl`. Confere em `options.h`: TLS 1.3, DTLS 1.2 e 1.3, AES-GCM e AES-CCM, ECC e x25519, SNI, EMS, `renegotiation_info`, *supported groups*, *tickets*, `NO_DH` e `NO_OLD_TLS`; e que não há `NO_PSK`, *early data*, TLS 1.0, SSLv3 nem *truncated HMAC* |
| `af/wolftls` (GPLv3) | ligação cgo: TLS 1.2 ou 1.3 com PSK sobre um `net.Conn`. O wolfSSL não bloqueia: lê de um *buffer* que uma *goroutine* enche, e por isso `Read` e `Write` podem correr juntos, como faz o `net/http` |
| `af/akmatls` | identificadores Ua\* de H.3, identidade `3GPP-AKMA;<A-KID>`, configuração `crypto/tls` do K.2 |
| `af/af/internal/profile` | um servidor por perfil; os perfis PSK só com a *tag* `wolfssl` |
| `af/scripts/make-test-ca.sh` | CA de teste e certificado ECDSA P-256 do AF, com o FQDN em `subjectAltName` |
| `af/scripts/build-af.sh` | compila com o wolfSSL quando a biblioteca existe: binários estáticos, GPLv3 |

**O que a ligação impõe**, qualquer que seja a configuração:

- uma versão por contexto: TLS 1.2 ou TLS 1.3;
- TLS 1.3: `only_dhe_psk` no AF e no UE, e nenhum *early data* (U12);
- TLS 1.2: só `ECDHE-PSK-AES128-GCM-SHA256` no AF, com EMS e `renegotiation_info`;
- grupos P-256, x25519 e P-384;
- o AF exige o `server_name` igual ao FQDN e recusa um `ClientHello` sem ele;
- retomada só por *tickets*. O cache de sessões por *session ID* fica desligado no AF. Os *tickets* são cifrados por um *callback* próprio, com AES-256-GCM, e o nome de cada *ticket* aponta para a identidade, o A-KID e a expiração da K_AF do *handshake* completo. Um *ticket* só vale até o menor entre a expiração da K_AF e 24 horas (U8). Depois disso, o *handshake* é completo, e a política de expiração se aplica.

**Conferido nos testes de `af/wolftls`**, com o `ClientHello` e o `ServerHello` lidos da rede:

- TLS 1.2: o `ClientHello` oferece `D0 01` primeiro e depois `C0 2B`, `C0 2F`, `C0 2C` e `C0 30` (TS 33.222 §5.4.0.1 passo 1); tem SNI, EMS, `renegotiation_info`, *ticket* e os grupos `00 17`, `00 1D`, `00 18`. O `ServerHello` escolhe `D0 01` e tem EMS e `renegotiation_info`. O *hint* `3GPP-AKMA` chega ao UE;
- TLS 1.3: `psk_key_exchange_modes` só com `psk_dhe_ke`; a identidade em `pre_shared_key`; sem `early_data`;
- recusas: K_AF errada; identidade desconhecida (o erro do AF guarda a causa, como `K_AKMA_NOT_PRESENT`); `ClientHello` só com `psk_ke`; outro `server_name`;
- retomada nas duas versões, com a identidade e o A-KID do *handshake* completo; *ticket* recusado depois da expiração da K_AF;
- o `net/http` sobre a ligação, com cinco requisições numa conexão;
- interoperabilidade: o `openssl s_client` 3.0.13 completa o TLS 1.3 com PSK contra o AF do wolfSSL.

**Diferenças em relação ao texto aprovado.**

1. `[MVP]` **wolfSSL 5.9.4**, e não a 5.8.4. A U10 pedia 5.8.4 ou mais nova; a 5.9.4 é a mais nova estável em 21/10/2026 e corrige CVEs de 2026.
2. `[DIVERGÊNCIA]` **Sem `D0 02`** (DV-10). A §3.4 dizia que o `ClientHello` oferece `ECDHE-PSK-AES128-GCM-SHA256` e `ECDHE-PSK-AES256-GCM-SHA384`. O wolfSSL não tem a segunda. TS 33.210 §6.2.3 só a recomenda (*"Support of the following cipher suite is recommended"*); a obrigatória, `D0 01`, está. Registrado como DV-10 por decisão do responsável em 21/10/2026.
3. `[MVP]` **Sem os *handlers* de `fork` do RNG do wolfSSL** (`--disable-rng-autofork`). Eles usam `dlopen`, que um binário estático não tem, e os processos Go não usam o wolfSSL depois de um `fork()`.
4. `[INTERPRETAÇÃO]` **Identidade PSK só como texto.** Depois de recusar um *ticket*, o wolfSSL oferece os octetos dele ao *callback* de PSK externa como se fossem uma identidade. A ligação recusa, sem consultar o AF, toda identidade que não seja UTF-8 sem caracteres de controle (RFC 4279 §5.1: a identidade é UTF-8).
5. `[INTERPRETAÇÃO]` **O UE dos perfis PSK só aceita PSK.** TS 24.109 K.3.3 1)c) e TS 33.222 §5.4.0.1 passo 1 mandam o UE indicar outros métodos de autenticação no `ClientHello`. O `ClientHello` do TLS 1.2 oferece suítes sem PSK, e o do TLS 1.3 leva os algoritmos de assinatura. Mas o cliente não confia em certificado nenhum nesses perfis: se o AF escolher certificado, o *handshake* falha. É a U6: o K.2 está noutra porta.
6. `[MVP]` **Retomada vista pelo UE no TLS 1.3.** O wolfSSL informa ao cliente que ele ofereceu o *ticket*, e não se o AF o aceitou. O `uaclient` escreve *"session ticket offered"*; quem diz se houve retomada é o log do AF.
7. `[DIVERGÊNCIA]` **OCSP *stapling*** (DV-10). TS 33.210 §6.2.2 e §6.2.3 recomendam (*should*) a extensão *certificate status request*. A CA de teste não tem OCSP, e o K.2 não a usa. Registrado como DV-10 por decisão do responsável em 21/10/2026.
8. `[MVP]` **HTTP/1.1** nos perfis TLS. TS 33.210 §6.2.1 só impõe regras se houver HTTP/2.
9. **Renegociação.** O AF e o UE mandam `renegotiation_info` (RFC 5746), e nenhum dos dois renegocia: o wolfSSL foi compilado sem renegociação, e o Go não a aceita no servidor. TS 33.210 §6.2.3: *"The server shall accept client-initiated renegotiation only if secured"*; aqui não aceita nunca.
10. **No K.2, *tickets* de 24 horas.** O `crypto/tls` aceita *tickets* por sete dias. O AF marca cada *ticket* com a hora de emissão e recusa os de mais de 24 horas (TS 33.222 §5.3.1.5). Quando o Go reempacota um *ticket* do TLS 1.2 retomado, perde a marca. O AF então recusa esse *ticket*, e a sessão não passa das 24 horas.
11. **O UE pede a K_AF ao ME só depois do desafio**, também no `digest-http`. Antes, o `uaclient` pedia a K_AF antes de conectar. Agora pede quando o identificador Ua\* é conhecido: depois do desafio no Digest, quando o *cipher suite* já foi escolhido, e no *callback* de PSK. Um desafio de outro AF já não leva a pedido nenhum ao ME.
12. **Documentos.** Os perfis entraram no tutorial nas partes 5 a 9, onde o AF, o UE e o roteiro já estavam, e não numa parte 10 (seção 6).

### 3.9 Verificação (21/10/2026)

`e2e.sh` com **115 verificações e nenhuma falha** (execução `20261021-164401`; `e2e.md` §5.4). Os casos da seção 3.7:

| Caso | Resultado |
|---|---|
| UE → AF por K.2, TLS 1.3 e TLS 1.2 | atendido: `nc 00000002` em TLS 1.3 com AF_ID `…0100011301` [R2t.k2-tls13, T01]; TLS 1.2 com `…010001C02B` [R2t.k2-tls12, T02]; retomada na segunda conexão [R2t.k2-resumed] |
| K.2 com certificado de outro FQDN | atendido: o UE encerra o TLS antes de mandar o A-KID [T09]; teste de unidade `TestDigestTLSWrongCertificate` |
| K.2 com `realm` de outro FQDN | atendido no teste de unidade `TestK2RealmOfAnotherAF`: o UE não pede a K_AF ao ME nem manda o A-KID. O AF do laboratório tem um FQDN só, e o roteiro não monta esse caso |
| UE → AF por PSK TLS 1.3 | atendido: `psk_dhe_ke` (teste de `af/wolftls`), 200, AF_ID `…0100000007`, segunda conexão retomada [R2t.psk13, T03] |
| UE → AF por PSK TLS 1.2 | atendido: suíte `D001`, *hint* `3GPP-AKMA`, 200, AF_ID `…010001D001`, retomada [R2t.psk12, R2t.psk12-resumed, T04] |
| cliente que só oferece `psk_ke` | atendido: recusado [R2t.psk-ke, T08] |
| PSK com K_AF errada | atendido: o *handshake* falha [T06] |
| A-KID desconhecido do AAnF | atendido: *handshake* recusado; `No K_AKMA at the AAnF … handshake refused` no log do AF [T07, T07b] |
| AF_ID de um perfil fora da política do AAnF | atendido: *handshake* recusado; o AAnF responde 403 [T10, T10b] |
| conferência independente | atendido: `openssl s_client` 3.0.13 com `-psk_identity "3GPP-AKMA;<A-KID>"` e a K_AF calculada pelo oráculo Python recebe 200 com o A-KID [T05] |
| K_AF vencida, `refetch` e `reject` | atendido: sessão encerrada no vencimento [T11]; `refetch` pede a K_AF de novo, e o *handshake* seguinte passa [T12, T12b]; `reject` recusa [T13, T13b] |
| logs do AF, do UE e do wolfSSL | atendido: nenhuma cadeia de 64 hexadecimais [L1, R7]; o wolfSSL foi compilado sem *debug* e não escreve log |
| o ME entrega a K_AF de cada AF_ID dos perfis | atendido [R2t.me] |

**Como a expiração foi testada.** O AAnF dá a K_AF com validade de 1 hora e devolve a mesma enquanto ela vale. O roteiro sobe AF extras com o relógio adiantado em 3590 s: o AF vê a K_AF vencer em 10 s. Para o `refetch` ter uma K_AF nova, o roteiro remove o contexto no AAnF e o registra de novo com o `seed`, como faria uma nova autenticação primária. Sem isso, o AAnF devolve a mesma K_AF, que o AF, com o relógio adiantado, já acha vencida. É o mesmo efeito que o N22 mostra de propósito.

## 4. Fase (c): DTLS

**Estado (22/10/2026):** **aprovada** pelo responsável em 21/10/2026, com os pontos U13 a U21 como propostos (seção 5): DTLS 1.3 e DTLS 1.2, CoAP dentro do DTLS. **Implementada e verificada** (seções 4.7 e 4.8).

### 4.1 Situação normativa

`[3GPP]` TS 33.535 Anexo C, *"AKMA Ua\* protocol based on DTLS"*:

- C.1: o DTLS é **opcional** no UE e no AF; se suportado, valem as cláusulas seguintes;
- C.1.1: o UE é o cliente DTLS. *"The PSK identity specified in B.1 for TLS is also applicable for DTLS"*;
- C.1.2: o AF *"should be able to require that a certain key (i.e., KAF) used to secure the Ua reference point based on DTLS"*;
- C.2.1: *"The TLS profile specified in TS 33.210 clause 6.2 apply to DTLS 1.3"*;
- C.2.2: *"The procedures given in B.1.3.2.2 for TLS 1.3 is also applicable for DTLS 1.3"*. *"AKMA PSK identity should be delivered via DTLS message."*

`[3GPP]` TS 33.220 H.3: `( 0x01,0x00,0x00,0x00,0x09 )` *"Ua security protocol for DTLS according to Annex O"*. O Anexo O é o DTLS do GBA, que o Anexo C do AKMA reproduz: só DTLS 1.3, com o procedimento do TLS 1.3 de TS 33.222 §5.4.0.2; e *"It does not have any dependency on the application layer protocol"* (O.2.2). TS 33.535 A.4 manda tirar o identificador Ua\* do Anexo H.

`[3GPP]` TS 33.210 §6.2.1: *"If DTLS is supported then DTLS 1.2 as specified in RFC 6347 shall be supported and DTLS 1.3 as specified in RFC 9147 should be supported"*; NOTA 2: os requisitos do perfil valem para o DTLS.

`[LACUNA]` **DTLS 1.2.** TS 33.210 o exige de quem tem DTLS; TS 33.535 C.1.1 aceita a identidade de B.1, que inclui o TLS 1.2; mas C.2 só descreve o DTLS 1.3, e H.3 não dá identificador Ua\* para o DTLS 1.2: o `…09` é do Anexo O, que só tem DTLS 1.3, e o `01 00 01 yy zz` é dito *"for TLS 1.2"*. Ponto U13.

`[3GPP]` RFC 9147 §5.1: contra negação de serviço, o servidor *"MAY respond with a HelloRetryRequest"* com um *cookie* sem estado, e o cliente *"MUST send a new ClientHello with the cookie"*. No DTLS 1.2, a mesma função é do `HelloVerifyRequest` (RFC 6347 §4.2.1).

`[3GPP]` RFC 8446 Apêndice E.7: *"Implementations can ensure safety from cross-protocol related output by not reusing PSKs between TLS 1.3 and TLS 1.2"*.

### 4.2 Protocolo de aplicação

`[LACUNA]` Nem o Anexo C nem o Anexo O dizem o que vai dentro do DTLS. Proposta (U14): **CoAP** (RFC 7252) sobre DTLS, o esquema `coaps`, porta UDP 5684. É o protocolo de aplicação usual do DTLS, e a fase (d), OSCORE, precisa dele de qualquer modo.

O UE faz `GET coaps://af1.akma.example/ua/v1/session`:

- mensagem confirmável (CON) com *token* aleatório e as opções `Uri-Host` (o FQDN) e `Uri-Path` (`ua`, `v1`, `session`);
- o AF responde no ACK (*piggybacked*) com `2.05 Content`, `Content-Format` 50 (`application/json`) e o mesmo corpo dos perfis TLS, `{"aKid": …, "kafExpiry": …}`; o *handshake* já autenticou os dois lados, como no PSK TLS;
- erros: `4.04 Not Found` para outro recurso ou outro `Uri-Host`, `4.05 Method Not Allowed`, `4.03 Forbidden` com `{"cause": "KAF_EXPIRED"}`;
- retransmissão do CON com os parâmetros padrão de RFC 7252 §4.8 (`ACK_TIMEOUT` 2 s, `MAX_RETRANSMIT` 4); o DTLS retransmite o *handshake* com os temporizadores do wolfSSL.

O CoAP fica num pacote Go próprio, `af/coap` (Apache-2.0), só com o necessário: codificação da mensagem (RFC 7252 §3), opções, CON e ACK.

`[DIVERGÊNCIA]` RFC 7252 §9.1.3.1: no modo PSK, as implementações de CoAP *"MUST support the mandatory-to-implement cipher suite TLS_PSK_WITH_AES_128_CCM_8"*, uma suíte de DTLS 1.2 sem ECDHE. TS 33.210 §6.2.3 só admite suítes AEAD com ECDHE, e TS 33.535 C.2.1 aplica TS 33.210 ao DTLS. Proposta (U15): seguir TS 33.210 e não ter `TLS_PSK_WITH_AES_128_CCM_8`, registrando a escolha em `desvios.md`.

### 4.3 DTLS 1.3

`[3GPP]` C.2.2, com B.1.3.2.2 e TS 24.109 K.3.3, como o PSK TLS 1.3 da seção 3.3:

- **UE:** deriva a K_AF com o identificador `01 00 00 00 09` antes do `ClientHello`, pelo serviço de chaves do ME; oferece `3GPP-AKMA;<A-KID>` em `psk_identities` (a identidade vai *"via DTLS message"*), só `psk_dhe_ke`, o `server_name` com o FQDN e nenhum *early data*;
- **AF:** responde ao primeiro `ClientHello` com `HelloRetryRequest` e *cookie* (U16); no segundo, lê o A-KID, pede ao AAnF a K_AF de `<FQDN>.0100000009` e a usa como PSK externa, com SHA-256 e `TLS_AES_128_GCM_SHA256`. Sem K_AF, recusa a identidade;
- AF_ID: `af1.akma.example.0100000009`, nas políticas do AAnF, do NEF e do UE.

### 4.4 DTLS 1.2

Se o ponto U13 ficar na proposta (b): o procedimento do PSK TLS 1.2 da seção 3.4, sobre DTLS 1.2 (RFC 6347):

- o AF manda `HelloVerifyRequest` com *cookie* antes de qualquer estado;
- a suíte é `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256`, que RFC 8442 define para o TLS 1.2 e o DTLS 1.2; *hint* `3GPP-AKMA`; identidade `3GPP-AKMA;<A-KID>`;
- o identificador Ua\* é `01 00 01 D0 01`, o da suíte, como no TLS 1.2. A K_AF é a mesma do perfil `psk-tls12`: o AF_ID é o mesmo. Cada sessão (D)TLS 1.2 tira as suas chaves dos números aleatórios das duas pontas, e a K_AF não é repartida com o DTLS 1.3 (RFC 8446 E.7).

### 4.5 Retomada, expiração e o que mais muda

- **Portas:** `psk-dtls13` em UDP 5684, a porta padrão do `coaps`; `psk-dtls12` em UDP 5685. Um perfil por porta, como na U6.
- **Retomada:** como nos perfis PSK TLS (U8): *tickets* até a expiração da K_AF e no máximo 24 horas.
- **K_AF vencida:** como na U9: o AF encerra a associação DTLS com o alerta `close_notify`; o *handshake* seguinte segue a política `refetch` ou `reject`.
- **MTU:** 1200 octetos, abaixo do MTU do túnel GTP; o *handshake* com PSK cabe com folga.
- **Connection ID** (RFC 9147 §9): não (U20). O endereço do UE é fixo enquanto dura a sessão PDU. Atrás do NAT do N6, o mapeamento UDP pode expirar com a associação parada; o AF passa a ver outra porta de origem, e o UE faz um novo *handshake*, retomado pelo *ticket*.
- **Ligação `af/wolftls`:** passa a ter DTLS. No servidor, um *socket* UDP recebe os datagramas de todos os UEs e os separa pelo endereço de origem, uma associação por endereço; o wolfSSL lê um datagrama inteiro por vez, e a ligação trata os temporizadores de retransmissão.
- **Cliente:** `uaclient -profile psk-dtls13` ou `psk-dtls12`, com o CoAP do `af/coap`.

### 4.6 Verificação prevista

| Caso | Esperado |
|---|---|
| UE → AF por DTLS 1.3, CoAP GET | `HelloRetryRequest` com *cookie*; `psk_dhe_ke`; `2.05` com o A-KID; AF_ID `…0100000009` no log do AAnF |
| UE → AF por DTLS 1.2 (se U13 b) | `HelloVerifyRequest`; suíte `D0 01`; *hint* `3GPP-AKMA`; `2.05`; AF_ID `…010001D001` |
| segunda associação | retomada pelo *ticket* |
| `ClientHello` que só oferece `psk_ke` | recusado (TS 33.210 §6.2.2) |
| K_AF errada; A-KID desconhecido; AF_ID fora da política | *handshake* recusado |
| `ClientHello` sem o *cookie* | o AF não pede nada ao AAnF antes do segundo `ClientHello` |
| K_AF vencida, `refetch` e `reject` | associação encerrada; novo *handshake* conforme a política |
| conferência independente (U21) | o `openssl s_client -dtls1_3` do OpenSSL 4.1 com `-psk_identity` e a K_AF de um vetor, mandando o GET em CoAP, recebe `2.05` |
| logs | nenhuma chave |

O `e2e.sh` ganha o passo R2d, com o UE no seu *namespace*, e um grupo D na máquina do AF.

### 4.7 Implementação (21/10/2026)

**Componentes.**

| Componente | Mudança |
|---|---|
| `af/scripts/build-wolfssl.sh` | `--enable-dtls-mtu`, para fixar o MTU; `WOLFSSL_DTLS_MTU_ADDITIONAL_READ_BUFFER` de 16384 (abaixo); a conferência das macros passa a exigir `WOLFSSL_DTLS_MTU` |
| `af/wolftls` | DTLS 1.2 e 1.3. A fila de entrada guarda datagramas inteiros, e cada `Read` devolve um registro. Durante o *handshake*, a espera segue o temporizador de retransmissão do wolfSSL (`wolfSSL_dtls_get_current_timeout`, `wolfSSL_dtls_got_timeout`); depois dele, só no DTLS 1.3 com mensagem não confirmada. No servidor, `ListenDTLS`: um *socket* UDP, uma associação por endereço de origem, até 1024 |
| `af/scripts/build-af.sh` | o *hash* da biblioteca entra numa constante do *shim* de `af/wolftls` (abaixo) |
| `af/coap` (novo) | mensagem e opções (RFC 7252 §3, §5.4, §5.10); cliente com CON, retransmissão (§4.2, §4.8) e resposta no ACK; respondedor que repete a resposta a um CON repetido (§4.5), responde `Reset` a um CON malformado e ao *ping* |
| `af/af` | perfis `psk-dtls13` e `psk-dtls12`; recurso CoAP (`server.CoAP`); *handshake* de no máximo 15 s; associação sem pedido por 5 minutos encerrada |
| `af/ue/uaclient` | `-profile psk-dtls13` e `psk-dtls12`; nova associação quando o AF encerra a anterior |
| políticas | AF_ID `af1.akma.example.0100000009` no AAnF, no NEF e no UE; o `…010001D001` já estava |

**Cookies (U16).**

- **DTLS 1.2.** O `HelloVerifyRequest` usa um *callback* próprio: HMAC-SHA-256, com um segredo do servidor, sobre o endereço do cliente. O padrão do wolfSSL lê o endereço do *socket*, que a ligação não lhe dá.
- **DTLS 1.3.** O *cookie* do `HelloRetryRequest` fica preso ao endereço do cliente com `wolfSSL_dtls_set_peer`.
- **Achado.** No DTLS 1.3, ao responder sem estado ao primeiro `ClientHello`, o wolfSSL chama o *callback* de PSK para saber a suíte da PSK (*"Ask the user for the ciphersuite matching this identity"*, `src/dtls.c`); não confere *binder* nessa hora. Do jeito ingênuo, o AF pediria a K_AF ao AAnF antes do *cookie*, contra a U16. A ligação responde a essa fase sem perguntar nada (`wolfSSL_is_stateful`): devolve só a suíte. A K_AF é pedida no segundo `ClientHello`, depois do *cookie*, quando o wolfSSL confere o *binder*. Teste: `TestDTLSNoKeyBeforeCookie`.

**Datagramas maiores que o MTU.** O wolfSSL lê no máximo o MTU mais `WOLFSSL_DTLS_MTU_ADDITIONAL_READ_BUFFER` (500 por padrão), 1700 octetos com o MTU de 1200, e o resto do datagrama se perde. O OpenSSL 4.1 oferece uma partilha de chave híbrida com ML-KEM: com uma identidade do tamanho de um A-KID, o seu segundo `ClientHello` tinha 1702 octetos, chegava cortado ao AF e o *handshake* não terminava. O MTU de 1200 vale para o que o AF manda; para receber, a biblioteca passa a ser compilada com 16384 octetos a mais. O teste `TestOpenSSLDTLS13Client` cobre as duas identidades, a curta e a do tamanho de um A-KID.

**Compilação depois de trocar a biblioteca.** O Go não vê o `libwolfssl.a`. Desde a fase (b), o `build-af.sh` punha o *hash* da biblioteca numa macro do compilador C, para que uma biblioteca nova recompilasse `af/wolftls`. O pacote era recompilado, mas nenhum código usava a macro: o objeto saía igual, e o Go reaproveitava os binários já ligados. Com isso, o AF da execução `20261022-120810` ainda tinha o wolfSSL de antes do ajuste da leitura, e o D03 falhou do mesmo jeito. Agora o *hash* vai para uma constante do *shim* (`wolftls_lib`): uma biblioteca nova muda o objeto, e o Go liga os binários de novo.

**Testes de `af/wolftls`** (DTLS): `psk_dhe_ke`, identidade e `server_name` no segundo `ClientHello`; `HelloRetryRequest` com *cookie* (1.3) e `HelloVerifyRequest` (1.2) como primeira resposta; nenhuma chave pedida antes do *cookie*; K_AF errada, identidade desconhecida, outro `server_name` e `psk_ke` recusados; retomada nas duas versões com uma consulta só; associação encerrada no vencimento da K_AF; um registro por `Read`; e, com `OPENSSL_DTLS13`, o `openssl s_client -dtls1_3` do OpenSSL 4.1.

**Diferenças em relação ao texto aprovado.**

1. `[MVP]` **Tempos.** *Handshake* de no máximo 15 s e associação encerrada depois de 5 minutos sem pedido. O UDP não tem fim de conexão, e as associações abandonadas ocupariam o AF.
2. `[MVP]` **CoAP só com resposta no ACK.** O AF responde sempre no ACK, e o cliente não trata a resposta separada (RFC 7252 §5.2.2).
3. **Recusas no DTLS 1.2.** Com uma K_AF errada, o registro que não decifra é descartado em silêncio (RFC 6347 §4.1.2.7). O cliente só desiste no fim do seu prazo.
4. **OpenSSL 4.1 *beta* 2** (09/10/2026) na conferência independente: a versão final ainda não saiu. Ele fica fora do repositório e é dado ao roteiro por `OPENSSL_DTLS13`; sem ele, a verificação D03 aparece como pulada.

### 4.8 Verificação (22/10/2026)

`e2e.sh` com **126 verificações e nenhuma falha** (execução `20261022-124514`; `e2e.md` §5.5). Os casos da seção 4.6:

| Caso | Resultado |
|---|---|
| UE → AF por DTLS 1.3, CoAP GET | atendido: `2.05` com o A-KID, AF_ID `…0100000009` e suíte `1301` no log do AF [R2d.dtls13, D01]; `HelloRetryRequest` com *cookie* como primeira resposta e `psk_dhe_ke` no segundo `ClientHello`, nos testes de `af/wolftls` (`TestDTLSCookieOnTheWire`) |
| UE → AF por DTLS 1.2 | atendido: suíte `D001`, `2.05`, AF_ID `…010001D001` [R2d.dtls12, D02]; `HelloVerifyRequest` e *hint* `3GPP-AKMA` nos testes de `af/wolftls` |
| segunda associação | atendido: retomada nas duas versões [D01, R2d.dtls12-resumed, D02] |
| `ClientHello` que só oferece `psk_ke` | atendido: recusado [R2d.psk-ke] |
| K_AF errada; A-KID desconhecido; AF_ID fora da política | atendido: K_AF de outro AF_ID [D05]; A-KID desconhecido, com `K_AKMA_NOT_PRESENT` do AAnF no log [D04, D04b]. O AF_ID fora da política é recusado pelo AAnF no mesmo caminho do PSK TLS [T10, T10b]; o roteiro não o repete no DTLS |
| `ClientHello` sem o *cookie* | atendido nos testes de `af/wolftls` (`TestDTLSNoKeyBeforeCookie`): nenhuma chave pedida nas duas versões. O roteiro não monta esse caso |
| K_AF vencida, `refetch` e `reject` | atendido nos testes: associação encerrada no vencimento (`TestDTLSExpiry`; teste do perfil `psk-dtls13` do AF). A política do *handshake* seguinte é a mesma do PSK TLS [T11 a T13b] |
| conferência independente (U21) | atendido: `openssl s_client -dtls1_3` do OpenSSL 4.1.0 *beta* 2, com `-psk_identity "3GPP-AKMA;<A-KID>"` e a K_AF calculada em Python a partir do vetor, manda o GET em CoAP e recebe `2.05` com o A-KID [D03] |
| logs | atendido: nenhuma cadeia de 64 hexadecimais [L1, R7] |

## 5. Pontos para aprovação

| Nº | Ponto | Proposta | Base |
|---|---|---|---|
| U1 | Como derivar da K_AF a PSK do TLS 1.3 e o *premaster* do TLS 1.2 | A K_AF, os 32 octetos, é a PSK. No TLS 1.2, a K_AF entra como o PSK da RFC 5489 §2; no TLS 1.3, como a PSK externa, com SHA-256. Nenhuma KDF a mais | `[LACUNA]` TS 33.222 §5.4, TS 24.109 K.3 e TS 33.535 B.1.3 dizem só "derivar da K_AF", e TS 33.210 não define nada. A separação de chaves por protocolo já está na K_AF, pelo identificador Ua\* no AF_ID (TS 33.220 H.2). Uma KDF a mais seria invenção. A análise formal de Ko et al. (2024, artigo 13 da literatura) usa a própria K_AF como PSK (*"compute KAF as the pre-shared key PSK"*), e o `adaptive-security` do responsável também |
| U2 | Separador da identidade PSK no TLS 1.2 | `3GPP-AKMA;<A-KID>`, em UTF-8 (RFC 4279 §5.1) | `[LACUNA]` TS 24.109 K.3.2 e) não dá o separador; o GBA usa `;` (TS 24.109 §5.3.3.1.2), e o TLS 1.3 do AKMA também (K.3.3) |
| U3 | Identidade da aplicação no ME | *uid* como fronteira; executável como conferência de melhor esforço (seção 2.2) | `[INTERPRETAÇÃO]` NOTA de TS 33.535 §4.4.3 |
| U4 | Canal e protocolo do serviço de chaves | *socket* Unix em `/run/ueransim/akma-<SUPI>.sock`, protocolo de linha (seção 2.4). Na implementação, o nome leva o SUPI, porque o `nr-ue -n` roda vários UEs num processo; o `socket` da configuração continua valendo para um UE só | — |
| U5 | Um AF_ID por aplicação | a configuração do UE recusa AF_ID repetido | TS 33.535 §4.4.3 item 3 |
| U6 | Uma porta por perfil no AF | sim (seção 3.5). Uma porta que negociasse K.2 ou K.3 no mesmo `ClientHello` (TS 33.222 §5.4.0.1 NOTA 3) fica fora: o K.2 é do `crypto/tls`, e o K.3 do wolfSSL | `[MVP]` |
| U7 | *Cipher suites* | os de TS 33.210 §6.2 (seções 3.2 a 3.4); cada um vira um AF_ID nas políticas do AAnF, do NEF e do UE | `[3GPP]` TS 33.220 H.3; TS 33.210 §6.2 |
| U8 | Retomada de sessão | suportada; nos perfis PSK, até a expiração da K_AF e no máximo 24 horas; no K.2, até 24 horas | `[3GPP]` TS 33.222 §5.3.1.5 e §5.4.1.4 |
| U9 | K_AF vencida nos perfis PSK | `refetch`: novo pedido ao AAnF no *handshake* seguinte; `reject`: *handshake* recusado; sessões abertas encerradas | `[3GPP]` TS 33.535 §6.4.2; §4.4.0 |
| U10 | Licença e *build* | `af/wolftls` e `build-wolfssl.sh` sob GPLv3; *tag* `wolfssl`; wolfSSL 5.8.4-stable ou mais nova, fixada com SHA-256. Depois de compilar, o *script* confere as macros de `options.h` (`WOLFSSL_TLS13`, `WOLFSSL_DTLS13`, ausência de `NO_PSK`, `HAVE_AESCCM`), como o `check_wolfssl.sh` do `adaptive-security`. Primeiro passo da fase (b): compilar e conferir a suíte `D0 01` e o `psk_dhe_ke` | decisão do responsável de 20/10/2026; CVE-2025-11935 e CVE-2025-11932 |
| U11 | ADR | ADR-0003 revisto: a faixa privada fica para o Ua\* simplificado; os perfis da norma usam os identificadores de TS 33.220 H.3 | — |
| U12 | 0-RTT | desligado nos perfis PSK | `[MVP]` sem sigilo futuro e sujeito a repetição (artigo 13 da literatura); TS 33.210 §6.2 não o exige |

**Fase (c), DTLS** (aprovados pelo responsável em 21/10/2026, como propostos):

| Nº | Ponto | Proposta | Base |
|---|---|---|---|
| U13 | DTLS 1.2 | **(b)** DTLS 1.2 com o procedimento do TLS 1.2 e o identificador `01 00 01 yy zz` da suíte (`…010001D001`), no perfil `psk-dtls12`. Alternativas: (a) só DTLS 1.3, registrando em `desvios.md` a falta do DTLS 1.2 que TS 33.210 exige; (c) DTLS 1.2 com o `…09` do DTLS 1.3, o que repartiria a mesma K_AF entre o DTLS 1.2 e o 1.3, contra RFC 8446 E.7 | `[LACUNA]` TS 33.210 §6.2.1 (*shall*); TS 33.535 C.1.1 e C.2; TS 33.220 H.3. `[INTERPRETAÇÃO]` na (b): H.3 dá `01 00 01 yy zz` ao TLS 1.2, e o DTLS 1.2 é o TLS 1.2 sobre datagramas (TS 33.210 §6.2.1 NOTA 2) |
| U14 | Protocolo dentro do DTLS | CoAP (RFC 7252), `coaps`, `GET /ua/v1/session` com resposta `2.05` em JSON (seção 4.2); pacote `af/coap` | `[LACUNA]` Anexos C e O não o fixam (*"no dependency on the application layer protocol"*). O CoAP é também a base da fase (d) |
| U15 | Suíte obrigatória do CoAP | não ter `TLS_PSK_WITH_AES_128_CCM_8`; registrar em `desvios.md` a escolha de TS 33.210 sobre RFC 7252 | `[DIVERGÊNCIA]` RFC 7252 §9.1.3.1 contra TS 33.210 §6.2.3 e TS 33.535 C.2.1 |
| U16 | *Cookie* contra negação de serviço | sempre: `HelloRetryRequest` no DTLS 1.3, `HelloVerifyRequest` no DTLS 1.2; o AF só pede a K_AF ao AAnF depois do *cookie* | RFC 9147 §5.1 (*MAY*); RFC 6347 §4.2.1. Sem ele, um endereço forjado faria o AF consultar o AAnF |
| U17 | Portas | UDP 5684 (`psk-dtls13`, porta do `coaps`) e 5685 (`psk-dtls12`) | RFC 7252 §6.2; U6 |
| U18 | Identificador do DTLS 1.3 | `01 00 00 00 09`, AF_ID `<FQDN>.0100000009` | `[INTERPRETAÇÃO]` TS 33.220 H.3 dá o `…09` ao DTLS do Anexo O, que o Anexo C do AKMA reproduz; TS 33.535 A.4 remete ao Anexo H. O TLS 1.3 do AKMA usa o `…07` do GBA da mesma forma |
| U19 | Retomada, expiração, MTU | como U8 e U9; MTU de 1200 octetos | TS 33.222 §5.4.1.4; TS 33.535 §6.4.2 |
| U20 | *Connection ID* | não | `[MVP]` RFC 9147 §9 o deixa opcional. O endereço do UE não muda durante a sessão PDU; se o mapeamento do NAT do N6 expirar, um novo *handshake* retomado pelo *ticket* resolve |
| U21 | Conferência independente | compilar à parte, só para o teste, o OpenSSL 4.1 (a versão final, se já tiver saído; senão a *beta* mais nova), que traz o DTLS 1.3, e usar o `s_client -dtls1_3 -psk` contra o AF, como o T05 fez com o TLS 1.3. Não vai para o repositório nem para o guia | o OpenSSL 3.0 da VM do AF não tem DTLS 1.3, e nenhum deles tem `D0 01` no DTLS 1.2: o DTLS 1.2 fica sem conferência independente |

## 6. Arquivos previstos

- **Fase (a):**
  - `af/ue/ueransim/patches/0005-akma-key-service.patch`;
  - `af/ue/uaclient` (`-me-socket`);
  - `af/scripts/ue-netns.sh`;
  - `docs/af-side/ue-ueransim.md` e `ua-client.md`.
- **Fase (b):**
  - `af/akmatls`, `af/wolftls`;
  - `af/scripts/build-wolfssl.sh` e o *script* da CA de teste;
  - `af/af` (perfis) e `af/ue/uaclient` (perfis);
  - `af/config/afcfg*.yaml`, `free5gc/config/aanfcfg.yaml` e o `nefcfg.yaml` do *patch* de configuração (AF_IDs por perfil);
  - `af/scripts/e2e.sh` (grupo T);
  - `docs/decisions/` (ADR-0003 revisto) e `desvios.md`, se algum ponto virar desvio.
- **Nas duas:** matriz de rastreabilidade revisão 6; tutorial, parte 10.
- **Fase (c):**
  - `af/coap` (novo, Apache-2.0) e `af/wolftls` (DTLS);
  - `af/af` (perfis `psk-dtls13` e `psk-dtls12`) e `af/ue/uaclient`;
  - `af/config/afcfg*.yaml`, `free5gc/config/aanfcfg.yaml`, o `nefcfg.yaml` do *patch* de configuração e o *patch* 0005 (AF_ID `…0100000009`);
  - `af/scripts/e2e.sh` e `e2e-remote.sh` (R2d e grupo D);
  - `desvios.md` (U15 e, se for a escolha, U13);
  - tutorial, guia público e matriz revisão 7.

## 7. Fase (d): OSCORE

Fica para a próxima revisão deste documento, depois da fase (c):

- **OSCORE**, TS 33.535 Anexo D: CoAP, OSCORE (RFC 8613) e CBOR. Duas opções:
  - em Go, sobre o `af/coap` da fase (c), conferido com os vetores do RFC 8613 Apêndice C;
  - a libcoap (`--enable-oscore`, sobre o wolfSSL), como no `adaptive-security` do responsável. Lá o contexto OSCORE é estático; o estabelecimento AKMA de D.3 (`POST /akma`, CBOR e HKDF) teria de ser escrito nas duas opções.

  Pontos a decidir:
  - os parâmetros do HKDF;
  - a forma do HPLMN-ID no A-KID em CBOR;
  - a remissão de D.3.2 ao "Anexo Y";
  - o campo que D.3.2 lista e D.3.5 não tem.
