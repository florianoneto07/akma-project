# Ua\* simplificado do MVP: especificação

**MVP:** MVP-3 · **Data:** 18 de outubro de 2026 · **Estado:** aprovada pelo responsável em 18/10/2026; implementada em `af/` (`af.md`) e verificada (`teste-af.md`)

Este documento define, mensagem a mensagem, o protocolo entre o UE e o AF de teste. O AF (`af.md`) e o cliente Ua\* (`ua-client.md`) são implementados a partir dele.

## 1. Situação normativa

`[3GPP]` TS 33.535 §4.4.1: o Ua\* é específico de cada aplicação. As exigências genéricas são três:

- carregar o A-KID;
- proteger o ponto de referência com a K_AF;
- tratar a expiração da K_AF.

Os perfis normativos estão em TS 33.535 Anexo B.1 e TS 24.109 Anexo K:

- **K.2:** HTTP Digest dentro de um túnel TLS, com o AF autenticado por certificado;
- **K.3:** TLS com PSK.

`[MVP]` O Ua\* do MVP é o Digest de K.2 **sem o túnel TLS**. Ele segue as regras de K.2 b), c) e d), de TS 24.109 §5.2 e de TS 33.535 B.1.2.2 em tudo o que não depende do TLS. Como falta o TLS, não é o perfil K.2 e usa o identificador Ua\* da faixa privada, `FF 00 00 00 01` (ADR-0003, TS 33.220 H.2 NOTA 2). É um protocolo de laboratório; as limitações estão na seção 10.

## 2. Identificadores e segredo

| Item | Valor no MVP | Fonte |
|---|---|---|
| FQDN do AF | `af1.akma.example` (configurável) | — |
| AF_ID | FQDN ‖ `FF 00 00 00 01`; forma de transporte `af1.akma.example.FF00000001` | TS 33.535 A.4; ADR-0003; plano P9 |
| *product token* do UE | `3gpp-akma` no `User-Agent`, em **toda** requisição | TS 24.109 K.2 b) 1) e §5.2.1; TS 33.535 B.1.2.2 passo 2 |
| `realm` | `3GPP-bootstrapping-akma@af1.akma.example` | TS 33.535 B.1.2.2 passo 3; TS 24.109 §5.2.2.1 3); DV-6 |
| `username` | o A-KID, exatamente como o UE o derivou | TS 24.109 K.2 c) 1); TS 33.535 B.1.2.2 passo 5 |
| senha | base64(K_AF), RFC 4648 §4, 44 caracteres | TS 24.109 K.2 c) 2) |
| algoritmo | `SHA-256` | RFC 7616; DV-5; TS 33.222 §5.3.0 NOTA 1b (MD5 desaconselhado) |

## 3. Transporte e FQDN

`[MVP]` HTTP/1.1 sobre TCP, sem TLS. No laboratório, o AF escuta em `192.168.16.129:8080` (configurável).

`[3GPP]` TS 33.220 §4.5.2: a K_AF só coincide se o UE e a rede usarem o mesmo FQDN na derivação. Cada protocolo sobre o Ua precisa dizer como garante isso. O Ua\* do MVP usa a hipótese (3) daquela cláusula: o protocolo leva ao AF o nome usado pelo UE, no cabeçalho `Host` do HTTP/1.1 (como TS 33.222 §5.3.0 passo 2).

- O cliente envia `Host: af1.akma.example:8080`. Sem DNS no laboratório, ele se conecta ao endereço IP configurado, mas o `Host` leva o FQDN.
- `[INTERPRETAÇÃO]` O AF compara o nome do `Host`, sem a porta e sem diferenciar maiúsculas, com o seu único FQDN configurado. Se diferir, responde `421 Misdirected Request` (RFC 9110 §15.5.20), sem desafio. O mesmo FQDN vai no `realm` e no AF_ID enviado ao AAnF.

## 4. Recurso protegido

`[MVP]` Um recurso só: `GET /ua/v1/session`.

- A requisição é o *"Application Session Establishment Request"* de TS 33.535 §6.2.1 passo 1.
- A resposta 200 é o *"Application Session Establishment Response"* do passo 9, com corpo JSON:

  ```json
  {"aKid": "<A-KID>", "kafExpiry": "<RFC 3339>"}
  ```

  `kafExpiry` deixa o UE saber até quando a sessão vale (§4.4.1). O SUPI recebido do AAnF nunca vai para o UE.
- Outro caminho dá `404`; outro método, `405`.

## 5. Troca de mensagens

### 5.1 Primeira requisição e desafio

UE → AF, sem `Authorization`:

```
GET /ua/v1/session HTTP/1.1
Host: af1.akma.example:8080
User-Agent: akma-uaclient/0.1 3gpp-akma
```

AF → UE:

```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Digest realm="3GPP-bootstrapping-akma@af1.akma.example", qop="auth", algorithm=SHA-256, nonce="<base64 de 32 octetos aleatórios>", opaque="<base64 de 16 octetos aleatórios>", charset=UTF-8
```

- `realm`, `qop`, `nonce` e `opaque` vão entre aspas; `algorithm` e `charset`, sem aspas (RFC 7616 §3.3, último parágrafo).
- `[MVP]` O AF oferece `qop="auth"`. TS 24.109 §5.2.6 permite `auth-int` (*"may"*); o AF o oferece só se configurado, e o cálculo da seção 5.3 cobre os dois.
- `[INTERPRETAÇÃO]` Sem `userhash`, embora RFC 7616 §3.4.4 diga que o servidor *SHOULD* ativá-lo. O AF precisa do A-KID em claro para pedir a K_AF ao AAnF, e não conseguiria recuperá-lo de H(A-KID:realm). TS 24.109 K.2 c) 1) manda pôr o A-KID no `username`.
- O AF guarda o desafio: nonce, opaque, instante de emissão e maior `nc` aceito. O nonce vale 300 s (configurável).

### 5.2 Verificação no UE

Antes de responder, o cliente confere o desafio:

1. esquema `Digest` e `algorithm=SHA-256`;
2. `realm` com duas partes separadas pelo último `@`. A primeira é `3GPP-bootstrapping-akma`: o AF escolheu AKMA (B.1.2.2 passo 3). A segunda é igual ao FQDN que o cliente usou (B.1.2.2 passo 4; K.2 d));
3. `qop` contém `auth`.

Se alguma conferência falhar, o cliente encerra e não envia o A-KID.

### 5.3 Requisição autenticada

UE → AF:

```
GET /ua/v1/session HTTP/1.1
Host: af1.akma.example:8080
User-Agent: akma-uaclient/0.1 3gpp-akma
Authorization: Digest username="<A-KID>", realm="3GPP-bootstrapping-akma@af1.akma.example", uri="/ua/v1/session", algorithm=SHA-256, nonce="<nonce>", nc=00000001, cnonce="<base64 de 16 octetos aleatórios>", qop=auth, response="<64 hexadecimais>", opaque="<opaque>"
```

`[3GPP]` Cálculo (RFC 7616 §3.4.1 a §3.4.3), com H = SHA-256 em hexadecimal minúsculo e KD(s, d) = H(s ":" d):

```
A1       = A-KID ":" realm ":" base64(K_AF)
A2       = "GET" ":" uri                         (qop=auth)
A2       = "GET" ":" uri ":" H(corpo)            (qop=auth-int)
response = KD(H(A1), nonce ":" nc ":" cnonce ":" qop ":" H(A2))
```

`nc` tem exatamente oito dígitos hexadecimais e começa em `00000001` para cada nonce (RFC 7616 §3.4 e §3.5).

### 5.4 Verificação no AF

O AF verifica nesta ordem, e para no primeiro erro:

| # | Verificação | Falha | Referência |
|---|---|---|---|
| 1 | `Host` com o FQDN do AF | `421` | seção 3 |
| 2 | `Authorization` presente, esquema `Digest` | `401` com desafio novo (iniciação do AKMA) | TS 33.535 §6.5; RFC 7616 §3.3 |
| 3 | parâmetros obrigatórios presentes e bem formados; `algorithm=SHA-256`; `qop` entre os oferecidos; sem `username*` nem `userhash=true` | `400` | RFC 7616 §3.4 (*"the proper response is a 4xx error code"*) |
| 4 | `uri` igual ao alvo da requisição | `400` | RFC 7616 §3.4.6 |
| 5 | `realm` igual ao do AF | `401` com desafio novo | TS 24.109 K.2 d) |
| 6 | `nonce` emitido pelo AF e dentro da validade; `opaque` igual ao emitido com ele | `401` com desafio novo e `stale=true` se o nonce venceu | RFC 7616 §3.3 (`stale`) |
| 7 | `nc` maior que o último aceito para o nonce | `401` com desafio novo (repetição) | RFC 7616 §3.4 (`nc`), §5.5 |
| 8 | `username` é um A-KID sintaticamente válido | `400` | TS 29.522 Tab. 5.14.5.4.2-1 |
| 9 | K_AF disponível (seção 6) | `403` ou `503` (seção 8) | TS 33.535 B.1.2.2 passo 6; §6.2.1 passos 2 e 9 |
| 10 | `response` igual ao calculado, em comparação de tempo constante | `401` com desafio novo; tentativa registrada em `warn`, sem a senha | RFC 7616 §3.4, §3.6 |

### 5.5 Resposta e confirmação mútua

AF → UE, se tudo confere:

```
HTTP/1.1 200 OK
Authentication-Info: rspauth="<64 hexadecimais>", qop=auth, nc=00000001, cnonce="<cnonce>"
Content-Type: application/json

{"aKid": "<A-KID>", "kafExpiry": "2026-10-18T17:35:31Z"}
```

`[3GPP]` `rspauth` é calculado como `response`, mas com `A2 = ":" uri` para `qop=auth` e `A2 = ":" uri ":" H(corpo)` para `qop=auth-int`. `rspauth`, `cnonce` e `nc` são obrigatórios quando há `qop` (RFC 7616 §3.5).

`[MVP]` O cliente **recusa** qualquer resposta 200 sem `Authentication-Info` válido: falta de `rspauth`, valor errado, ou `cnonce` e `nc` diferentes dos enviados. É a confirmação de chave do AF para o UE, critério 4 do MVP-3; sem ela, um AF falso poderia aceitar qualquer cliente.

### 5.6 Requisições seguintes

`[3GPP]` RFC 7616 §3.6 e TS 33.222 §5.3.0 NOTA 2: toda requisição seguinte leva `Authorization`.

- O cliente reaproveita o nonce, com `nc` incrementado e `cnonce` novo.
- Se o nonce vencer, o AF responde `401` com `stale=true`. O cliente refaz o cálculo com o nonce novo, com o mesmo A-KID e a mesma K_AF.
- Cada requisição passa de novo por todas as verificações da seção 5.4.

## 6. K_AF no AF

`[3GPP]` TS 33.535 §6.2.1 passo 2: se não tem contexto ativo para o A-KID, o AF pede a K_AF ao AAnF com `Naanf_AKMA_ApplicationKey_Get` (TS 29.535 §4.2.2.3), pelo `POST {apiRoot}/naanf-akma/v1/retrieve-applicationkey`. O corpo é `AkmaAfKeyRequest`:

- `afId` na forma de transporte;
- `aKId` exatamente o `username` recebido.

| Resposta do AAnF | O AF |
|---|---|
| `200` com `AkmaAfKeyData` (`kaf` hexadecimal, `expiry`, `supi`) | guarda `{A-KID, K_AF, expiry, SUPI}` e segue a verificação |
| `403` com `K_AKMA_NOT_PRESENT` | `403`, causa `AKMA_KEY_NOT_AVAILABLE` |
| `403` sem causa (AF fora da política do AAnF), outro erro HTTP, falha de rede | `503`, causa `AANF_UNAVAILABLE` |

- **Seleção do AAnF:** `[MVP]` o endereço do AAnF é configurado no AF (TS 33.535 §6.7, *"locally configured"*). Sem *token* OAuth2: o AAnF autoriza o AF pelo NF Instance ID do certificado de cliente, desde o MVP-11 (`docs/free5gc-side/tls-sbi.md` §7.5).
- **Contexto ativo:** a K_AF guardada vale até `expiry`. `[3GPP]` O AF nunca a usa depois disso (§4.4.0: tempo de vida máximo). A K_AF fica em memória e nunca aparece em log.

## 7. Expiração da K_AF

`[3GPP]` TS 33.535 §6.4.2: quando a K_AF vence, o AF pode recusar o acesso do UE ou renovar a K_AF, *"based on its policy"*. Se recusar, pode informar ao UE que a K_AF expirou, por meio próprio do Ua\*. O §6.4.3 diz que não há renovação explícita de K_AF; sem nova autenticação primária, a K_AF continua a mesma.

`[MVP]` Política configurável `kafExpiry`:

| Valor | Comportamento |
|---|---|
| `refetch` (padrão) | o AF descarta a K_AF vencida e pede de novo ao AAnF com o mesmo A-KID. Enquanto a K_AKMA não muda, o AAnF devolve a mesma K_AF com `expiry` novo. Se o A-KID foi substituído por uma reautenticação, o AAnF responde 403, e o AF devolve `403 AKMA_KEY_NOT_AVAILABLE` |
| `reject` | o AF responde `403` com causa `KAF_EXPIRED`, o aviso de expiração do §6.4.2. O UE pode tentar de novo com o A-KID de uma K_AUSF nova |

A renovação de chaves de sessão derivadas da K_AF (§6.4.3; TS 24.109 K.2 NOTA 2) fica fora do MVP.

## 8. Respostas de erro

`[MVP]` O Ua\* é específico da aplicação (§4.4.1), então o formato dos erros é do projeto. Erro sem desafio leva corpo JSON `{"cause": "<CAUSA>"}`:

| Código | Causa | Quando |
|---|---|---|
| `400` | `INVALID_REQUEST` | parâmetro do Digest ausente, malformado ou não suportado; `uri` diferente do alvo; A-KID inválido |
| `401` | — (desafio) | sem `Authorization`; `realm`, nonce, `nc` ou `response` não conferem |
| `403` | `AKMA_KEY_NOT_AVAILABLE` | o AAnF não tem K_AKMA para o A-KID. TS 33.535 B.1.2.2 passo 6: erro **sem** o `realm` do desafio. O UE pode tentar com o A-KID mais recente (§6.2.1 passo 9) |
| `403` | `KAF_EXPIRED` | K_AF vencida, com a política `reject` |
| `404`, `405` | `NOT_FOUND`, `METHOD_NOT_ALLOWED` | outro caminho ou método |
| `421` | `MISDIRECTED_REQUEST` | `Host` com outro nome |
| `503` | `AANF_UNAVAILABLE` | o AF não conseguiu a K_AF por falha do AAnF ou da própria configuração |

## 9. Iniciação do AKMA

`[3GPP]` TS 33.535 §6.5: se o AF exige AKMA e a requisição não traz parâmetros AKMA, o AF responde com uma mensagem de iniciação do AKMA, cuja forma depende do Ua\*.

`[MVP]` A mensagem de iniciação é o `401` da seção 5.1. Toda requisição sem `Authorization` recebe o desafio, com ou sem o *product token* `3gpp-akma`. O AF só oferece AKMA; o *product token* é registrado em log e não muda a resposta (critério 6 do MVP-3).

## 10. Limitações do protocolo sem TLS

| Limitação | Efeito | Mitigação no MVP |
|---|---|---|
| sem confidencialidade | o A-KID e o conteúdo trafegam em claro; quem escuta associa o A-KID à sessão | nenhuma; o A-KID é pseudônimo e muda a cada autenticação primária |
| AF não autenticado pelo transporte | um AF falso pode receber o pedido | o UE só aceita resposta com `rspauth` válido, e só quem tem a K_AF o calcula |
| integridade só da autenticação | com `qop=auth`, o corpo das mensagens pode ser alterado | `qop=auth-int` disponível por configuração |
| `Host` pode ser alterado no caminho | — | a K_AF depende do FQDN: com outro nome, o `response` não confere |
| ataque de dicionário *offline* | — | a senha é base64 de 256 bits aleatórios |

O perfil normativo (K.2 com TLS, ou K.3.3) fica para o MVP-6+; a migração está descrita no ADR-0003.

## 11. Configuração do AF

| Campo | Padrão no laboratório | Observação |
|---|---|---|
| `fqdn` | `af1.akma.example` | único nome aceito no `Host` |
| `uaStarId` | `FF00000001` | ADR-0003 |
| `listen` | `192.168.16.129:8080` | |
| `realmPrefix` | `3GPP-bootstrapping-akma` | DV-6: primeira parte configurável |
| `qop` | `[auth]` | `auth-int` opcional |
| `nonceLifetime` | `300s` | |
| `kafExpiry` | `refetch` | ou `reject` |
| `aanfUri` | `http://192.168.16.130:8010` | §6.7, configuração local |

## 12. Vetores de teste

- **RFC 7616 §3.9.1:** `username="Mufasa"`, senha `Circle of Life`, `realm="http-auth@example.org"`, `GET /dir/index.html`, SHA-256, os `nonce` e `cnonce` do exemplo e `nc=00000001` dão `response=753927fa0e85d155564e2e272a28d1802ca10daf4496794697cf8db5856cb6c1`.
- **Vetor do projeto:** A-KID e base64(K_AF) do vetor `testbed-ue` do MVP-0 (`af1.akma.example.FF00000001`), com `nonce`, `cnonce` e `nc` fixos. O `response` e o `rspauth` são calculados também por um *script* Python independente, guardado com os vetores do Ua\* em `af/`, para não depender só da implementação Go. O `reference/` continua só com o oráculo da KDF (D1).
