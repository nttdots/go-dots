#pragma once

#include <coap2/coap.h>
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <openssl/x509v3.h>

void response_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_tid_t);

void method_handler(coap_context_t *,
                    coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_string_t *,
                    coap_string_t *,
                    coap_pdu_t *);
void method_from_server_handler(coap_context_t *,
                    coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_string_t *,
                    coap_string_t *,
                    coap_pdu_t *);

void nack_handler(struct coap_context_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_nack_reason_t,
                    const coap_tid_t);

void event_handler(coap_context_t *context,
                    coap_event_t event,
                    void *data);

/**
 * Validate common name call back
 *
 * @param cn                The determined CN from the certificate
 * @param asn1_public_cert  The ASN.1 DER encoded X.509 certificate
 * @param asn1_length       The ASN.1 length
 * @param coap_session      The CoAP session associated with the certificate update
 * @param depth             Depth in cert chain.  If 0, then client cert, else a CA
 * @param validated         TLS layer can find no issues if 1
 * @param arg               The same as was passed into coap_context_set_pki()
 *                          in setup_data->cn_call_back_arg
 *
 * @return @c 1 if accepted, else @c 0 if to be rejected.
 */
int validate_cn_call_back(const char *cn,
                          const uint8_t *asn1_public_cert,
                          size_t asn1_length,
                          coap_session_t *coap_session,
                          unsigned depth,
                          int validated,
                          void *arg);
/**
 * Get peer common name (from certificate issuer names)
 * @param session   The CoAP session
 * @param buf       The return value
 * @param buflen    The max length of the return value 
 * @return          Return the index of the next matching entry -1 if not found
 */
int coap_dtls_get_peer_common_name(coap_session_t *session,
                                    char *buf,
                                    size_t buf_len);

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

typedef struct sni_entry {
  char *sni;
  SSL_CTX *ctx;
} sni_entry;

typedef struct coap_openssl_context_t {
  coap_dtls_context_t dtls;
  coap_tls_context_t tls;
  coap_dtls_pki_t setup_data;
  int psk_pki_enabled;
  size_t sni_count;
  sni_entry *sni_entry_list;
} coap_openssl_context_t;

typedef struct coap_strlist_t {
    struct coap_strlist_t *next;
    coap_string_t* str;
} coap_strlist_t;

// Coppied from the coap_session_internal.h file in the libcoap
struct coap_subscription_t {
  struct coap_subscription_t *next; /**< next element in linked list */
  struct coap_session_t *session;   /**< subscriber session */

  unsigned int non_cnt:4;  /**< up to 15 non-confirmable notifies allowed */
  unsigned int fail_cnt:2; /**< up to 3 confirmable notifies can fail */
  unsigned int dirty:1;    /**< set if the notification temporarily could not be
                            *   sent (in that case, the resource's partially
                            *   dirty flag is set too) */
  unsigned int has_block2:1; /**< GET request had Block2 definition */
  uint16_t tid;             /**< transaction id, if any, in regular host byte order */
  coap_block_t block2;     /**< GET request Block2 definition */
  size_t token_length;     /**< actual length of token */
  unsigned char token[8];  /**< token used for subscription */
  struct coap_string_t *query; /**< query string used for subscription, if any */
};

void coap_set_dirty(coap_resource_t *resource, char *query, int length);

int coap_check_subscribers(coap_resource_t *resource);
int coap_check_dirty(coap_resource_t *resource);
char* coap_get_token_subscribers(coap_resource_t *resource);
int coap_get_size_block2_subscribers(coap_resource_t *resource);
coap_block_t coap_create_block(unsigned int num, unsigned int m, unsigned int size);

coap_strlist_t* coap_common_name(coap_strlist_t* head, coap_strlist_t* tail, char* str);

void coap_session_handle_release(coap_session_t *session);
coap_session_t* coap_get_session_from_resource(coap_resource_t *resource);