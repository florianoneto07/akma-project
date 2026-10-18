# `free5gc/openapi`: modelos e cliente do Naanf_AKMA

**MVP:** MVP-1 · **Decisão:** D6 do plano · **Data:** 17 de outubro de 2026

## 1. O que falta no `openapi` do free5GC

`[FREE5GC]` O `free5gc/openapi` v1.2.4, fixado pelo free5GC v4.2.3, já traz o tipo de NF `AANF`, o nome de serviço `naanf-akma` e o `AanfInfo` (`docs/baseline/ambiente-verificado.md` §3). Não traz os tipos de dados do Naanf_AKMA nem um cliente para ele. Sem os tipos não há produtor (AAnF) nem consumidores (AUSF no MVP-2, AF no MVP-3).

`[3GPP]` Os tipos vêm de dois lugares. `AkmaKeyInfo` e `CtxRemove` estão em TS 29.535 §5.1.6.2. `AkmaAfKeyRequest`, `AkmaAfKeyData` e `ServiceDisableNotif` estão em TS 29.522 §5.14.5.3. Os nomes JSON seguem `TS29535_Naanf_AKMA.yaml` v1.2.0 (TS 29.535 Anexo A.2), que prevalece sobre o corpo do texto na codificação (Anexo A.1): `aKId`, `kAkma`, `afId`, `kaf`, `expiry`, `anonInd`, `notifUri`, `suppFeat`.

## 2. Arquivos

Nenhum arquivo existente do `openapi` é alterado; os *patches* só acrescentam arquivos.

| *Patch* (repositório) | Arquivos acrescentados ao `openapi` | Conteúdo |
|---|---|---|
| `free5gc/patches/openapi/0001-naanf-akma-models.patch` | `models/model_akma_key_info.go`, `model_ctx_remove.go`, `model_akma_af_key_request.go`, `model_akma_af_key_data.go`, `model_service_disable_notif.go` | os cinco tipos de dados |
| `free5gc/patches/openapi/0002-naanf-akma-client.patch` | `aanf/AKMA/{client,configuration,api_common}.go` e um arquivo por operação | cliente das três operações, no estilo dos clientes gerados do free5GC |

O cliente segue a Tab. 5.1.4.1-1 de TS 29.535: as três operações são `POST` em `{apiRoot}/naanf-akma/v1/` com os sufixos `register-anchorkey`, `retrieve-applicationkey` e `remove-context`. Ele trata:

- os 11 códigos de erro obrigatórios para `POST` (TS 29.500 Tab. 5.2.7.1-1), devolvendo o `ProblemDetails` recebido;
- 307 e 308 com `Location` e `3gpp-Sbi-Target-Nf-Id` (TS 29.535 Tab. 5.1.4.x.2-3/-4);
- o 204 de `retrieve-applicationkey` (campo `NoContent`), embora o AAnF deste projeto nunca o emita (DV-7).

## 3. Como o código usa a cópia com *patch*

O módulo do Go em `~/go/pkg/mod` é somente leitura e compartilhado por todas as NFs; ele não é alterado. `free5gc/scripts/prepare-openapi.sh`:

1. baixa `github.com/free5gc/openapi@v1.2.4` com `go mod download`;
2. copia o módulo para `free5gc/build/openapi`, fora do controle de versão;
3. aplica `free5gc/patches/openapi/*.patch` em ordem.

O `go.mod` do AAnF aponta para essa cópia com `replace github.com/free5gc/openapi => ../build/openapi`. O AUSF corrigido no MVP-2 usa a mesma cópia, por `../../../akma-free5gc/free5gc/build/openapi`.

O banco de testes do free5GC (`~/free5gc/test`) fixa outra versão, `v1.2.5-0.20260527003827-02dc71b4d94f`. Os mesmos *patches* aplicam sobre ela sem mudança. O segundo argumento do *script* escolhe a versão:

```bash
free5gc/scripts/prepare-openapi.sh free5gc/build/openapi-test v1.2.5-0.20260527003827-02dc71b4d94f
```

`[MVP]` No trabalho *upstream* os arquivos entram no próprio `free5gc/openapi`, e o `replace` desaparece.

## 4. Como reproduzir numa instalação limpa

Na máquina em que o AAnF será compilado, com Go 1.26.2 no `PATH` e este repositório clonado:

```bash
free5gc/scripts/prepare-openapi.sh
```

A saída lista os dois *patches* aplicados e termina com `patched free5gc/openapi v1.2.4 ready at …/free5gc/build/openapi`. O *script* só substitui um diretório que já seja uma cópia do `free5gc/openapi`.

## 5. Como verificar

```bash
cd free5gc/aanf && go test -run TestOpenapiClient ./internal/sbi/
```

O teste usa o cliente do *patch* contra o servidor do AAnF sobre HTTP/2 sem TLS (h2c). Ele registra o par do vetor `testbed-ue`, recupera a K_AF esperada pelo oráculo e remove o contexto. Depois confere que os erros chegam ao cliente com o `ProblemDetails` do AAnF: 403 `K_AKMA_NOT_PRESENT` e 404 `AKMA_CONTEXT_NOT_FOUND`.
