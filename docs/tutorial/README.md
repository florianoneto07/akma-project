# Tutorial — AKMA sobre o free5GC

Este tutorial descreve o que fazer **depois de instalar o free5GC** pelo guia oficial para usar AKMA (3GPP TS 33.535). Ele cresce a cada MVP; esta versão corresponde ao MVP-0.

| Parte | Conteúdo | Disponível |
|---|---|---|
| 1 | Pré-requisitos | sim |
| 2 | Conferir as derivações de chave | sim |
| 3 | Instalar o AAnF no free5GC | MVP-1 |
| 4 | Habilitar AKMA no UDM e no AUSF; provisionar um assinante | MVP-2 |
| 5 | Escrever um AF que usa AKMA | MVP-3 |
| 6 | Preparar o UE (UERANSIM) | MVP-4 |
| 7 | Executar o fluxo de ponta a ponta | MVP-5 |

---

## 1. Pré-requisitos

| Item | Versão usada no projeto |
|---|---|
| free5GC | v4.2.3, instalado pelo guia oficial, com `TestRegistration` passando |
| Go | 1.26.2 — a mesma do free5GC |
| Python | 3.10 ou mais recente (para o oráculo de referência) |
| UERANSIM | v3.2.6 (a partir da parte 6) |

Topologia de referência (ADR-0001-R2): o free5GC e o AAnF numa máquina; o AF e este repositório em outra. Uma única máquina também serve para as partes 1 e 2.

Se o Go não estiver instalado, instale-o num diretório de usuário, sem privilégios:

```bash
mkdir -p ~/sdk/go1.26.2 && curl -fsSL https://go.dev/dl/go1.26.2.linux-amd64.tar.gz | tar -xz -C ~/sdk/go1.26.2 --strip-components=1
```

```bash
export PATH=$HOME/sdk/go1.26.2/bin:$PATH
```

Confira o SHA-256 do arquivo com o publicado em https://go.dev/dl/ antes de usar. O projeto usou `990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282` para `go1.26.2.linux-amd64.tar.gz`.

---

## 2. Conferir as derivações de chave

O AKMA deriva três valores da chave K_AUSF que a autenticação primária do 5G deixa no AUSF e no UE (TS 33.535 Anexo A):

- **K_AKMA**, a chave âncora, guardada no AAnF;
- **A-TID**, parte do identificador A-KID que o UE apresenta à aplicação;
- **K_AF**, a chave que a aplicação (AF) recebe do AAnF.

Todas usam a KDF de TS 33.220 Anexo B.2 — a mesma que o free5GC já usa para K_SEAF. Detalhes em `docs/shared/kdf.md`.

### 2.1 Rodar o oráculo de referência

```bash
python3 reference/kdf/akma_kdf.py
```

A saída mostra a cadeia `S` e o resultado de cada derivação para o exemplo do baseline. O valor de K_AKMA deve ser `d6aa856e127fb8de3eb1a1d0cebf52a7df6b69b4a6c68c21c18c1c27d466dbba`.

### 2.2 Testar a biblioteca contra os vetores

```bash
cd shared/akma && go test ./...
```

```bash
cd shared/crosscheck && go test ./...
```

O primeiro compara a biblioteca Go com os vetores do oráculo; o segundo, com a KDF do próprio free5GC (`free5gc/util`). Os vetores são do projeto, não do 3GPP — ver `docs/shared/vetores.md`.

### 2.3 Dois erros que não dão mensagem de erro

1. **O SUPI vai à KDF sem o prefixo `imsi-`.** O free5GC trafega `imsi-208930000000001`; a KDF usa `208930000000001` (TS 33.501 A.7.0). Com o prefixo, todas as chaves saem erradas.
2. **O AF_ID da KDF não é o `afId` da API.** A API transporta `af1.akma.example.FF00000001`; a KDF usa o FQDN seguido de cinco octetos brutos `FF 00 00 00 01` (TS 33.535 A.4).

Os dois casos estão em `shared/akma/testdata/vectors.json`, seção `counterexamples`, com o resultado errado que produzem.

### 2.4 Um erro que só aparece para alguns assinantes

O Routing Indicator é uma cadeia de 1 a 4 dígitos em que cada dígito conta: `"0000"`, `"012"` e `"12"` são valores diferentes (TS 23.003 §2.2B). Guardá-lo como número apaga os zeros à esquerda e quebra o A-KID só dos assinantes cujo RID começa com zero. O UE do testbed usa `"0000"`.
