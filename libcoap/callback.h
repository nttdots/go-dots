#pragma once

#include <coap/coap.h>
#include <openssl/ssl.h>

void response_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_tid_t);

void method_handler(coap_context_t *,
                    coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    str *,
                    str *,
                    coap_pdu_t *);

void nack_handler(struct coap_context_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_nack_reason_t,
                    const coap_tid_t);

int coap_dtls_get_peer_common_name(coap_session_t *,
                                    char *,
                                    size_t);

int verify_certificate(coap_context_t *, coap_dtls_pki_t *);


typedef struct coap_dtls_context_t {
  SSL_CTX *ctx;
  SSL *ssl;	/* OpenSSL object for listening to connection requests */
  HMAC_CTX *cookie_hmac;
  BIO_METHOD *meth;
  BIO_ADDR *bio_addr;
} coap_dtls_context_t;

typedef struct coap_tls_context_t {
  SSL_CTX *ctx;
  BIO_METHOD *meth;
} coap_tls_context_t;

typedef struct coap_openssl_context_t {
  coap_dtls_context_t dtls;
  coap_tls_context_t tls;
  int psk_pki_enabled;
} coap_openssl_context_t;