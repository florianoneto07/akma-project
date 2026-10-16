# Implementação de referência da KDF do AKMA

Oráculo para validação cruzada com a implementação em Go.

Escrito diretamente a partir de **TS 33.220 V19.1.0, Anexo B.2**, sem consultar outra implementação. O propósito é que as duas implementações sejam independentes, de modo que a concordância entre elas seja evidência de correção.

## Uso

```
python3 akma_kdf.py
```

Mostra a cadeia `S` e o resultado de K_AKMA, A-TID e K_AF para o exemplo de `docs/00-baseline-normativo.md` §2.5.

## Gerador de vetores

```
python3 reference/kdf/gen_vectors.py > shared/akma/testdata/vectors.json
```

`gen_vectors.py` usa as funções de `akma_kdf.py` e acrescenta, escritos também a partir das normas e sem consultar o código Go: SUPI como parâmetro (TS 33.501 A.7.0), Home Network Domain (TS 23.003 §28.2), A-KID (TS 29.522) nas três codificações do ADR-0002 e as duas formas do AF_ID. Produz 8 vetores, as igualdades esperadas e dois contraexemplos. Ver `docs/shared/vetores.md`.

## O que implementa

| Derivação | FC | Norma |
|---|---|---|
| K_AKMA | `0x80` | TS 33.535 Anexo A.2 |
| A-TID | `0x81` | TS 33.535 Anexo A.3 |
| K_AF | `0x82` | TS 33.535 Anexo A.4 |

A KDF genérica é `HMAC-SHA-256(Key, S)` com `S = FC ‖ P0 ‖ L0 ‖ … ‖ Pn ‖ Ln`, `Li` em dois octetos *big-endian*, cadeias de caracteres em UTF-8 com normalização NFKC (TS 33.220 B.2.1.2).

## Aviso sobre os valores

**TS 33.535 V19.0.0 não fornece vetores de teste.** Os valores no arquivo são escolhidos pelo projeto e servem apenas para exercitar a construção da cadeia `S`. Não têm autoridade normativa e não devem ser citados como vetores 3GPP.

## Advertências reproduzidas no código

- O SUPI é passado sem o prefixo `imsi-`. Ver ADR-0002 e `docs/00-baseline-normativo.md` §2.3.
- O `AF_ID` é montado como `UTF-8(FQDN) ‖ 5 octetos brutos`, não como a cadeia de transporte com ponto e hexadecimal. Ver ADR-0002 §2.
- O identificador Ua* é da faixa privada. Ver ADR-0003.
