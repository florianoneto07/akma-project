# Cliente Ua\* do UE

**MVP:** MVP-3 (cliente com vetores) e MVP-4 (cliente no UE, com a K_AF do ME, pelo plano de usuário) · **Data:** 18 de outubro de 2026

## 1. Papel

`[3GPP]` TS 33.535 §6.2.1 passo 1: quando o UE inicia a comunicação com o AF, inclui o A-KID no pedido de estabelecimento de sessão; a K_AF pode ser derivada antes ou depois. §4.4.3: as aplicações do UE não têm acesso à K_AKMA.

O cliente `uaclient` (`af/ue/uaclient/`) faz o papel da aplicação do UE no Ua\* simplificado de `ua-star.md`:

- **Pedido de sessão:** envia o *product token* `3gpp-akma`.
- **Desafio:** confere se o `realm` escolhe AKMA para o FQDN que ele usou; se não, desiste sem enviar o A-KID.
- **Credenciais:** responde com o A-KID e a K_AF.
- **Prova do AF:** só aceita a resposta se o AF provar, pelo `rspauth`, que tem a mesma K_AF.

## 2. De onde vêm o A-KID e a K_AF

| Opção | Fonte | MVP |
|---|---|---|
| `-me-socket <caminho>` | o serviço de chaves do ME (`ue-ueransim.md` §3), que só entrega a K_AF dos AF_IDs da aplicação (TS 33.535 §4.4.3) | MVP-8 |
| `-nr-cli <nr-cli> -ue <imsi-…>` | o ME do UERANSIM, pelo comando `akma-kaf`; só para teste, com `akmaCliKaf: true` no UE | MVP-4 |
| `-vector <nome> [-vectors <arquivo>]` | um vetor de `shared/akma` | MVP-3 |
| `-akid <A-KID> -kaf <hex>` ou `-kaf-base64 <base64>` | linha de comando, para testes negativos | MVP-3 |

Com `-me-socket`, a K_AF chega pelo *socket* do ME; com `-nr-cli`, pela saída do `nr-cli`. Ela nunca aparece numa linha de comando nem na tela: o cliente nunca imprime a K_AF. Se o ME recusar o AF_ID, o cliente para com *"application not authorized for this AF_ID by the ME"*.

## 3. Pelo plano de usuário

`[MVP]` Decisão D5: o cliente roda no *namespace* do UE e vai ao AF pela sessão PDU. A opção `-bind <endereço>` fixa o endereço de origem das conexões no endereço da sessão (10.60.0.x, no `uesimtun0`). Assim, a regra que o `nr-ue` instala (`from 10.60.0.x lookup rt_uesimtun0`) leva o tráfego pelo `uesimtun0`, e não pelo `veth` do *namespace*.

O *script* `af/scripts/ue-netns.sh ua` junta tudo:

- pega o SUPI na configuração do UE e o endereço do `uesimtun0` no *namespace*;
- roda o cliente no *namespace*, com o seu usuário, e não como `root`, com `-me-socket /run/ueransim/akma-<SUPI>.sock` (ou `AKMA_SOCKET`) e `-bind`. É o *uid* desse usuário que a política do ME confere.

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -n 3
```

```
ME: K_AF derived for AF_ID af1.akma.example.FF00000001
AF af1.akma.example at 192.168.16.129:8080 from 10.60.0.1, A-KID rid0000.atid…
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until …
```

O AF vê a requisição vindo de 192.168.16.130, o endereço da VM do free5GC, por causa do NAT do N6 (`docs/free5gc-side/plano-de-usuario.md`).

## 4. Opções

| Opção | Padrão | |
|---|---|---|
| `-fqdn` | `af1.akma.example` | FQDN do AF: vai no `Host`, é conferido no `realm` e entra na K_AF |
| `-connect` | `192.168.16.129:8080` | endereço do AF |
| `-path` | `/ua/v1/session` | |
| `-ua-star-id` | `FF00000001` | identificador Ua\* (ADR-0003) |
| `-qop` | `auth` | ou `auth-int` |
| `-n` | 1 | requisições na mesma sessão Digest, com `nc` crescente |
| `-bind` | — | endereço de origem |
| `-kaf-file` | — | arquivo com a K_AF em hexadecimal (teste: K_AF anterior a uma nova autenticação, MVP-5) |
| `-user-agent` | `akma-uaclient/0.1 3gpp-akma` | |

Uma K_AF dada explicitamente (`-kaf`, `-kaf-base64`, `-kaf-file`) ou um vetor têm precedência sobre `-nr-cli`. Isso permite, no *namespace* do UE, apresentar um A-KID e uma K_AF anteriores.

## 5. Arquivos

| Caminho | Conteúdo |
|---|---|
| `af/ue/uaclient/client.go` | sessão Digest: desafio, credenciais, `rspauth`, `stale`; endereço de origem |
| `af/ue/uaclient/me.go` | `FromKeyService`: pede a K_AF ao serviço de chaves do ME; `FromME`: chama o `nr-cli` (teste) |
| `af/ue/uaclient/cmd/uaclient/` | o comando |
| `af/scripts/ue-netns.sh` | `ua`: o cliente no *namespace* do UE |

Testes: `TestHonestAF`, `TestAFWithoutKaf`, `TestChallengeChecks`, `TestHostAndProductToken`, `TestParseAkmaKaf`, `TestFromME` (com um `nr-cli` falso), `TestFromMERefused`, `TestParseKeyService`, `TestFromKeyService` (com um serviço de chaves falso) e `TestLocalAddr`.

## 6. Limitações

- **Identidade da aplicação.** O ME separa as aplicações pelo usuário; o executável é conferência de melhor esforço (`ua-star-norma.md` §2.2). Até o MVP-7, o `nr-cli` entregava a K_AF de qualquer AF a qualquer processo local; desde o MVP-8 isso só acontece com `akmaCliKaf: true`.
- **Sem TLS.** Ver `ua-star.md` §10.
