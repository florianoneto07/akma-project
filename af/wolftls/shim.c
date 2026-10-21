//go:build wolfssl

/*
 * Copyright 2026 Floriano Neto.
 * Licensed under the GNU General Public License v3.0 or later; see COPYING.
 */

#include "shim.h"
#include "_cgo_export.h"

/* The hash of an external PSK defaults to SHA-256 (RFC 8446 cl. 4.2.11). */
static const char tls13PSKSuite[] = "TLS13-AES128-GCM-SHA256";

static uintptr_t handle(WOLFSSL *ssl)
{
    return (uintptr_t)wolfSSL_get_psk_callback_ctx(ssl);
}

static int recv_cb(WOLFSSL *ssl, char *buf, int sz, void *ctx)
{
    (void)ssl;
    return wolftlsRecv((uintptr_t)ctx, buf, sz);
}

static int send_cb(WOLFSSL *ssl, char *buf, int sz, void *ctx)
{
    (void)ssl;
    return wolftlsSend((uintptr_t)ctx, buf, sz);
}

static unsigned int server_key(WOLFSSL *ssl, const char *identity,
                               unsigned char *key, unsigned int max)
{
    unsigned int lifetime = 0;
    unsigned int n = wolftlsServerKey(handle(ssl), (char *)identity,
                                      wolfSSL_get_current_cipher_suite(ssl),
                                      key, max, &lifetime);
    if (n > 0 && lifetime > 0)
        wolfSSL_set_timeout(ssl, lifetime);
    return n;
}

static unsigned int server_psk12(WOLFSSL *ssl, const char *identity,
                                 unsigned char *key, unsigned int max)
{
    return server_key(ssl, identity, key, max);
}

static unsigned int server_psk13(WOLFSSL *ssl, const char *identity,
                                 unsigned char *key, unsigned int max,
                                 const char **suite)
{
    *suite = tls13PSKSuite;
    return server_key(ssl, identity, key, max);
}

static unsigned int client_psk12(WOLFSSL *ssl, const char *hint, char *identity,
                                 unsigned int idMax, unsigned char *key,
                                 unsigned int keyMax)
{
    return wolftlsClientKey(handle(ssl), (char *)hint,
                            wolfSSL_get_current_cipher_suite(ssl),
                            identity, idMax, key, keyMax);
}

static unsigned int client_psk13(WOLFSSL *ssl, const char *hint, char *identity,
                                 unsigned int idMax, unsigned char *key,
                                 unsigned int keyMax, const char **suite)
{
    *suite = tls13PSKSuite;
    return wolftlsClientKey(handle(ssl), (char *)hint, 0x1301,
                            identity, idMax, key, keyMax);
}

static int ticket_cb(WOLFSSL *ssl, unsigned char name[WOLFSSL_TICKET_NAME_SZ],
                     unsigned char iv[WOLFSSL_TICKET_IV_SZ],
                     unsigned char mac[WOLFSSL_TICKET_MAC_SZ], int enc,
                     unsigned char *ticket, int inLen, int *outLen, void *ctx)
{
    (void)ctx;
    return wolftlsTicket(handle(ssl), name, iv, mac, enc, ticket, inLen, outLen);
}

void wolftls_setup_server(WOLFSSL_CTX *ctx, int tls13, int tickets)
{
    wolfSSL_CTX_SetIORecv(ctx, recv_cb);
    wolfSSL_CTX_SetIOSend(ctx, send_cb);
    if (tls13)
        wolfSSL_CTX_set_psk_server_tls13_callback(ctx, server_psk13);
    else
        wolfSSL_CTX_set_psk_server_callback(ctx, server_psk12);
    /* No resumption by session ID: only tickets, which ticket_cb binds to
     * the expiry of the key. */
    wolfSSL_CTX_set_session_cache_mode(ctx, WOLFSSL_SESS_CACHE_OFF);
    if (tickets) {
        wolfSSL_CTX_set_TicketEncCb(ctx, ticket_cb);
    } else if (tls13) {
        wolfSSL_CTX_no_ticket_TLSv13(ctx);
    } else {
        wolfSSL_CTX_NoTicketTLSv12(ctx);
    }
}

void wolftls_setup_client(WOLFSSL_CTX *ctx, int tls13)
{
    wolfSSL_CTX_SetIORecv(ctx, recv_cb);
    wolfSSL_CTX_SetIOSend(ctx, send_cb);
    if (tls13)
        wolfSSL_CTX_set_psk_client_tls13_callback(ctx, client_psk13);
    else
        wolfSSL_CTX_set_psk_client_callback(ctx, client_psk12);
}

void wolftls_attach(WOLFSSL *ssl, uintptr_t h)
{
    wolfSSL_SetIOReadCtx(ssl, (void *)h);
    wolfSSL_SetIOWriteCtx(ssl, (void *)h);
    wolfSSL_set_psk_callback_ctx(ssl, (void *)h);
}

void wolftls_error_string(int err, char *buf, int len)
{
    wolfSSL_ERR_error_string_n((unsigned long)err, buf, (unsigned long)len);
}
