# ADR-0002 — Codificação de identificadores e material de chave

**Estado:** Proposto
**Data:** 02 de setembro de 2026
**Substitui:** as recomendações preliminares registradas no baseline v2.0 §§3.4 e 4.3

---

## 1. Questão

As especificações do baseline não fixam a codificação textual de três elementos. É preciso decidir, registrar e parametrizar.

---

## 2. Decisões

| Elemento | Decisão | Fundamento |
|---|---|---|
| `kAkma` e `kaf` sobre SBI | **Hexadecimal minúsculo, 64 caracteres**, padrão ancorado `^[A-Fa-f0-9]{64}$` | Precedente direto de TS 29.509 V19.6.0 |
| Valor do A-TID na parte `username` do A-KID | **Base64 conforme IETF RFC 4648**, 44 caracteres | Analogia com o B-TID do GBA |
| Rótulo do A-TID no A-KID | **`atid`** na escrita; `atid` e `akid` aceitos na leitura | Texto normativo de TS 29.522 prevalece sobre o exemplo |
| K_AF como senha do HTTP Digest no Ua* | **Base64 conforme IETF RFC 4648** | TS 24.109 K.2 c) 2), normativo |
| `AF_ID` como P0 da KDF | `UTF-8_NFKC(FQDN) ‖ 5 octetos brutos`, distinto da forma de transporte | TS 33.535 A.4; analogia com `NAF_Id` |

Todas as três primeiras permanecem **parametrizáveis em configuração**, com os valores acima como padrão.

---

## 3. Fundamentação

### 3.1 Material de chave sobre SBI — hexadecimal

TS 29.509 V19.6.0 declara material de chave de comprimento fixo com padrão hexadecimal explícito:

| Tipo | Padrão | Bits |
|---|---|---|
| `Kseaf` | `[A-Fa-f0-9]{64}` | 256 |
| `KnrProSe` | `[A-Fa-f0-9]{64}` | 256 |
| `Msk` | `[A-Fa-f0-9]{128}` | 512 |
| `ResStar`, `HxresStar` | `[A-Fa-f0-9]{32}` | 128 |
| `SorMac`, `UpuMac` | `^[A-Fa-f0-9]{32}$` | 128 |

As mesmas especificações reservam base64 — via `Bytes` de TS 29.571 ou `format: byte` — para contêineres opacos de comprimento variável: `SorHeader`, `SorTransparentInfo`, `UpuTransparentInfo`, `SecuredPacket`, `EapPayload`.

A convenção da série 29 é discriminante: **chave de comprimento fixo → hexadecimal com padrão; contêiner opaco → base64**. K_AKMA e K_AF são chaves de 256 bits produzidas pela mesma KDF que K_SEAF. Seguem a mesma convenção.

Argumento adicional por ausência: `TS29522_AKMA.yaml` já importa `TS29122_CommonData.yaml`, de onde obtém `Uri` e `DateTime`. O tipo `Bytes` estava disponível e não foi referenciado por `kaf`.

Ver defeito D-D no registro de defeitos.

### 3.2 Valor do A-TID — base64

Esta decisão **reverte** a recomendação preliminar de hexadecimal registrada no baseline v2.0.

TS 33.220 V19.1.0, cláusula 4.5.2, passo 7: o B-TID do GBA é gerado em formato NAI como `base64encode(RAND)@BSF_servers_domain_name`, com base64 conforme IETF RFC 4648. A mesma especificação usa base64 no TMPI.

Argumentos:

1. **Analogia estrutural.** O A-KID e o B-TID têm a mesma estrutura — NAI com `username` derivado de material criptográfico e `realm` com o domínio da rede — e a mesma função.
2. **Herança declarada.** TS 24.109 Anexo K constrói os perfis Ua* do AKMA por modificação explícita dos perfis do GBA, com o AF do AKMA assumindo o papel do NAF.
3. **Precedente interno.** TS 24.109 K.2 c) 2) codifica K_AF em base64 no perfil Digest do próprio AKMA.
4. **Ausência de precedente contrário.** Nenhuma fonte do baseline usa hexadecimal em parte `username` de identificador NAI.

Contra-argumento avaliado e descartado: o alfabeto base64 padrão inclui `+`, `/` e `=`. O B-TID do GBA contém esses caracteres — RAND de 16 octetos produz 24 caracteres com preenchimento — e a cláusula 2.2 da IETF RFC 7542 os admite na parte `utf8-username`.

Um A-TID de 32 octetos produz 44 caracteres em base64, dos quais o último é de preenchimento.

**Esta decisão é do projeto, fundamentada em analogia, e não requisito 3GPP.** Deve ser declarada como tal no trabalho e submetida à discussão na *issue* free5GC #1146 e, se possível, ao 3GPP CT1.

### 3.3 As três codificações coexistem por camada

Não há inconsistência entre as decisões acima. Elas pertencem a camadas distintas:

| Camada | Objeto | Codificação |
|---|---|---|
| SBI, JSON | `kAkma`, `kaf` | hexadecimal |
| Identificador NAI | valor do A-TID no A-KID | base64 |
| Ua*, HTTP Digest | K_AF como senha | base64 |

A implementação mantém K_AKMA, K_AF e A-TID internamente como 32 octetos e aplica a codificação apropriada em cada fronteira. Confundi-las produz falha silenciosa: divergência de K_AF sem erro no SBI.

---

## 4. Parametrização

```
akma.codec.key_encoding   = hex        # hex | base64
akma.codec.atid_encoding  = base64     # base64 | base64url | hex
akma.akid.atid_label      = atid       # atid | akid   (escrita)
akma.akid.accept_labels   = atid,akid  #               (leitura)
```

Os valores padrão são os da seção 2. A parametrização existe para permitir interoperação com implementações que tenham decidido de outro modo, e para acomodar eventual esclarecimento normativo.

---

## 5. Testes derivados

| Teste | Verifica |
|---|---|
| `codec_test.go` | Ida e volta de cada codificação; rejeição de comprimento incorreto |
| `akid_test.go` | Construção e análise sintática do A-KID; RID `"0"`, `"012"`, `"0000"` |
| `akid_test.go` | Leitura tolerante de `akid`, escrita estrita de `atid` |
| `afid_kdf_test.go` | Que a forma de transporte do `AfId` e a forma de P0 produzem resultados distintos, e que só a segunda é usada na KDF |
| `supi_test.go` | Remoção do prefixo `imsi-` antes da KDF; rejeição de SUPI com prefixo em P1 |
