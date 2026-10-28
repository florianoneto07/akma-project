# Cliente Ua\* do UE

**MVP:** MVP-3 (cliente com vetores), MVP-4 (cliente no UE, com a K_AF do ME, pelo plano de usuário) e MVP-8 (serviço de chaves do ME e perfis TLS, DTLS e OSCORE da norma) · **Data:** 18 de outubro de 2026; perfis de 21 a 22/10/2026

## 1. Papel

`[3GPP]` TS 33.535 §6.2.1 passo 1: quando o UE inicia a comunicação com o AF, inclui o A-KID no pedido de estabelecimento de sessão; a K_AF pode ser derivada antes ou depois. §4.4.3: as aplicações do UE não têm acesso à K_AKMA.

O cliente `uaclient` (`af/ue/uaclient/`) faz o papel da aplicação do UE nos perfis Ua\* do AF (`ua-star-norma.md`), escolhidos com `-profile`:

| Perfil | O que o cliente faz | AF_ID da K_AF |
|---|---|---|
| `digest-tls` (TS 24.109 K.2) | abre TLS 1.3 ou 1.2, confere o certificado do AF contra o FQDN e a CA de `-ca`, e faz o Digest dentro do túnel | `<FQDN>.010001<cipher suite>`, pedido ao ME depois do desafio |
| `psk-tls13` (K.3.3) | TLS 1.3 com a K_AF como PSK externa, só `psk_dhe_ke`, identidade `3GPP-AKMA;<A-KID>` | `<FQDN>.0100000007`, pedido ao ME antes do `ClientHello` |
| `psk-tls12` (K.3.2) | TLS 1.2 com `ECDHE-PSK`; confere o *hint* `3GPP-AKMA` | `<FQDN>.010001<cipher suite>`, pedido ao ME depois do `ServerHello` |
| `psk-dtls13` (TS 33.535 Anexo C) | DTLS 1.3 com a K_AF como PSK; `GET` em CoAP confirmável, com retransmissão | `<FQDN>.0100000009`, pedido ao ME antes do `ClientHello` |
| `psk-dtls12` | DTLS 1.2 com `ECDHE-PSK`; CoAP | `<FQDN>.010001<cipher suite>`, depois do `ServerHello` |
| `oscore` (TS 33.535 Anexo D) | `POST /akma` em CBOR, deriva o contexto OSCORE da K_AF e dos dois *payloads*, e manda o `GET` protegido; com `4.01` sem proteção, estabelece de novo | `<FQDN>.0100000008`, pedido ao ME antes do `POST` |
| `digest-http` (padrão; teste) | o Ua\* simplificado de `ua-star.md`, sem TLS | `<FQDN>.<-ua-star-id>` |

Nos perfis Digest:

- **Pedido de sessão:** envia o *product token* `3gpp-akma`.
- **Desafio:** confere se o `realm` escolhe AKMA para o FQDN que ele usou; se não, desiste sem pedir a K_AF ao ME e sem enviar o A-KID (TS 33.535 B.1.2.2 passo 4).
- **Credenciais:** responde com o A-KID e a K_AF.
- **Prova do AF:** só aceita a resposta se o AF provar, pelo `rspauth`, que tem a mesma K_AF.

Nos perfis PSK, o *handshake* autentica os dois lados, e o `GET` do recurso vai sem autenticação HTTP. Os perfis PSK só existem num cliente compilado com o wolfSSL (`build-af.sh` depois de `build-wolfssl.sh`), que fica sob GPLv3. No perfil OSCORE, o AF prova ter a K_AF ao proteger a resposta com o contexto; o perfil não precisa do wolfSSL.

## 2. De onde vêm o A-KID e a K_AF

| Opção | Fonte | MVP |
|---|---|---|
| `-me-socket <caminho>` | o serviço de chaves do ME (`ue-ueransim.md` §3), que só entrega a K_AF dos AF_IDs da aplicação (TS 33.535 §4.4.3) | MVP-8 |
| `-nr-cli <nr-cli> -ue <imsi-…>` | o ME do UERANSIM, pelo comando `akma-kaf`; só para teste, com `akmaCliKaf: true` no UE | MVP-4 |
| `-vector <nome> [-vectors <arquivo>]` | um vetor de `shared/akma` | MVP-3 |
| `-akid <A-KID> -kaf <hex>` ou `-kaf-base64 <base64>` | linha de comando, para testes negativos | MVP-3 |

Com `-me-socket`, a K_AF chega pelo *socket* do ME; com `-nr-cli`, pela saída do `nr-cli`. Ela nunca aparece numa linha de comando nem na tela: o cliente nunca imprime a K_AF. Se o ME recusar o AF_ID, o cliente para com *"application not authorized for this AF_ID by the ME"*.

Com o ME, o cliente pede a K_AF quando o perfil fixa o identificador Ua\*, e uma vez por identificador. Com uma chave ou um vetor explícitos, usa a mesma K_AF em qualquer perfil: serve aos testes negativos, como uma K_AF de outro AF_ID.

## 3. Pelo plano de usuário

`[MVP]` Decisão D5: o cliente roda no *namespace* do UE e vai ao AF pela sessão PDU. A opção `-bind <endereço>` fixa o endereço de origem das conexões no endereço da sessão (10.60.0.x, no `uesimtun0`). Assim, a regra que o `nr-ue` instala (`from 10.60.0.x lookup rt_uesimtun0`) leva o tráfego pelo `uesimtun0`, e não pelo `veth` do *namespace*.

O *script* `af/scripts/ue-netns.sh ua` junta tudo:

- pega o SUPI na configuração do UE e o endereço do `uesimtun0` no *namespace*;
- roda o cliente no *namespace*, com o seu usuário, e não como `root`, com `-me-socket /run/ueransim/akma-<SUPI>.sock` (ou `AKMA_SOCKET`) e `-bind`. É o *uid* desse usuário que a política do ME confere.

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -n 3
```

```
AF af1.akma.example at 192.168.16.129:8080 from 10.60.0.1, profile digest-http
ME: K_AF derived for AF_ID af1.akma.example.FF00000001
session established (nc 00000001): AF authenticated by rspauth; A-KID rid0000.atid…; AF_ID af1.akma.example.FF00000001; K_AF valid until …
```

Com os perfis TLS (a CA de teste vem de `af/scripts/make-test-ca.sh`):

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls13 -reconnect -n 2
```

```
AF af1.akma.example at 192.168.16.129:8444 from 10.60.0.1, profile psk-tls13
ME: K_AF derived for AF_ID af1.akma.example.0100000007
session established: AF authenticated by the TLS handshake with K_AF; A-KID rid0000.atid…; AF_ID af1.akma.example.0100000007; TLS 1.3, cipher suite 1301; K_AF valid until …
session established: AF authenticated by the TLS handshake with K_AF; A-KID rid0000.atid…; AF_ID af1.akma.example.0100000007; TLS 1.3, cipher suite 1301, session ticket offered; K_AF valid until …
```

No TLS 1.3, o wolfSSL informa ao cliente que ofereceu o *ticket*, e não se o AF o aceitou; quem diz se a sessão foi retomada é o log do AF (`af.md` §4). No TLS 1.2 e no K.2, o cliente diz `resumed true` ou `false`.

O AF vê a requisição vindo de 192.168.16.130, o endereço da VM do free5GC, por causa do NAT do N6 (`docs/free5gc-side/plano-de-usuario.md`).

## 4. Opções

| Opção | Padrão | |
|---|---|---|
| `-profile` | `digest-http` | `digest-tls`, `psk-tls13`, `psk-tls12`, `psk-dtls13`, `psk-dtls12` ou `digest-http` |
| `-fqdn` | `af1.akma.example` | FQDN do AF: vai no `Host` e no `server_name`, é conferido no `realm` e no certificado, e entra na K_AF |
| `-connect` | `192.168.16.129` | endereço do AF; sem porta, a do perfil em `afcfg.yaml`: 8443, 8444, 8445, 5684 (UDP), 5685 (UDP), 5683 (UDP) ou 8080 |
| `-ca` | — | certificado da CA do AF, PEM; obrigatório com `digest-tls` |
| `-tls12` | — | `digest-tls` só com TLS 1.2 |
| `-reconnect` | — | uma conexão TLS ou associação DTLS nova por requisição, retomando a sessão anterior; no `oscore`, um `POST /akma` novo por requisição |
| `-rekey` | — | `oscore`: renova o contexto antes de cada requisição (RFC 8613 Apêndice B.2) |
| `-psk-ke-only` | — | `psk-tls13` ou `psk-dtls13` oferecendo só `psk_ke`, que TS 33.210 §6.2.2 proíbe; teste negativo |
| `-path` | `/ua/v1/session` | |
| `-ua-star-id` | `FF00000001` | identificador Ua\* de `digest-http` (ADR-0003) |
| `-qop` | `auth` | ou `auth-int` |
| `-n` | 1 | requisições; nos perfis Digest, na mesma sessão Digest, com `nc` crescente |
| `-bind` | — | endereço de origem |
| `-kaf-file` | — | arquivo com a K_AF em hexadecimal (teste: K_AF anterior a uma nova autenticação, MVP-5) |
| `-user-agent` | `akma-uaclient/0.1 3gpp-akma` | |

Uma K_AF dada explicitamente (`-kaf`, `-kaf-base64`, `-kaf-file`) ou um vetor têm precedência sobre `-nr-cli`. Isso permite, no *namespace* do UE, apresentar um A-KID e uma K_AF anteriores.

## 5. Arquivos

| Caminho | Conteúdo |
|---|---|
| `af/ue/uaclient/client.go` | sessão Digest, com ou sem TLS: desafio, K_AF do identificador Ua\* da conexão, credenciais, `rspauth`, `stale`; endereço de origem |
| `af/ue/uaclient/psk.go`, `psk_wolfssl.go`, `psk_nowolfssl.go` | sessão PSK pelo `af/wolftls`; sem a *tag* `wolfssl`, um erro que manda compilar com o wolfSSL |
| `af/ue/uaclient/dtls_wolfssl.go` | sessão DTLS: associação, CoAP pelo `af/coap`, nova associação quando o AF encerra a anterior |
| `af/ue/uaclient/me.go` | `FromKeyService`: pede a K_AF ao serviço de chaves do ME; `FromME`: chama o `nr-cli` (teste) |
| `af/akmatls/` | identificadores Ua\* de TS 33.220 H.3, identidade PSK e configuração TLS do K.2 |
| `af/ue/uaclient/cmd/uaclient/` | o comando |
| `af/scripts/ue-netns.sh` | `ua`: o cliente no *namespace* do UE |

Testes: `TestHonestAF`, `TestAFWithoutKaf`, `TestChallengeChecks`, `TestHostAndProductToken`, `TestParseAkmaKaf`, `TestFromME` (com um `nr-cli` falso), `TestFromMERefused`, `TestParseKeyService`, `TestFromKeyService` (com um serviço de chaves falso), `TestLocalAddr`, `TestK2KeyOfTheCipherSuite` e `TestK2RealmOfAnotherAF`. Os perfis PSK são testados com o AF em `af/af/internal/profile` (`TestPSKProfiles`, `TestPSKRefused`, `TestPSKExpiry`, com a *tag* `wolfssl`).

## 6. Limitações

- **Identidade da aplicação.** O ME separa as aplicações pelo usuário; o executável é conferência de melhor esforço (`ua-star-norma.md` §2.2). Até o MVP-7, o `nr-cli` entregava a K_AF de qualquer AF a qualquer processo local; desde o MVP-8 isso só acontece com `akmaCliKaf: true`.
- **Perfis PSK só com PSK.** O `ClientHello` do TLS 1.2 oferece também suítes sem PSK (TS 33.222 §5.4.0.1 passo 1), e o do TLS 1.3 leva os algoritmos de assinatura; mas o cliente não confia em certificado nenhum nesses perfis. Se o AF escolher a autenticação por certificado, o *handshake* falha. Cada perfil tem a sua porta (U6 de `ua-star-norma.md`); o K.2 é o perfil `digest-tls`.
- **`digest-http` sem TLS.** É o perfil de teste; ver `ua-star.md` §10.
