#pragma once

#include <coap2/coap.h>
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <openssl/x509v3.h>

coap_response_t response_handler(struct coap_context_t *,
                      coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_mid_t);

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
                    const coap_mid_t);

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

// Coppied from the internal file in the libcoap
struct coap_subscription_t {
  struct coap_subscription_t *next; /**< next element in linked list */
  struct coap_session_t *session;   /**< subscriber session */

  unsigned int non_cnt:4;  /**< up to 15 non-confirmable notifies allowed */
  unsigned int fail_cnt:2; /**< up to 3 confirmable notifies can fail */
  unsigned int dirty:1;    /**< set if the notification temporarily could not be
                            *   sent (in that case, the resource's partially
                            *   dirty flag is set too) */
  unsigned int has_block2:1; /**< GET request had Block2 definition */
  uint8_t code;            /** request type code (GET/FETCH)*/
  uint16_t mid;             /**< message id, if any, in regular host byte order */
  coap_block_t block;      /**< GET/FETCH request Block definition */
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
size_t coap_handle_add_option(coap_pdu_t *pdu, uint16_t type, unsigned int val);


/* Copied from the internal file in the libcoap */
coap_subscription_t *
coap_add_observer(coap_resource_t *resource,
                  coap_session_t *session,
                  const coap_binary_t *token,
                  coap_string_t *query,
                  int has_block2,
                  coap_block_t block,
                  uint8_t code);
int coap_delete_observer(coap_resource_t *resource,
                         coap_session_t *session,
                         const coap_binary_t *token);

struct coap_resource_t {
  unsigned int dirty:1;          /**< set to 1 if resource has changed */
  unsigned int partiallydirty:1; /**< set to 1 if some subscribers have not yet
                                  *   been notified of the last change */
  unsigned int observable:1;     /**< can be observed */
  unsigned int cacheable:1;      /**< can be cached */
  unsigned int is_unknown:1;     /**< resource created for unknown handler */
  unsigned int is_proxy_uri:1;   /**< resource created for proxy URI handler */

  /**
   * Used to store handlers for the seven coap methods @c GET, @c POST, @c PUT,
   * @c DELETE, @c FETCH, @c PATCH and @c IPATCH.
   * coap_dispatch() will pass incoming requests to handle_request() and then
   * to the handler that corresponds to its request method or generate a 4.05
   * response if no handler is available.
   */
  coap_method_handler_t handler[7];

  UT_hash_handle hh;

  coap_attr_t *link_attr; /**< attributes to be included with the link format */
  coap_subscription_t *subscribers;  /**< list of observers for this resource */

  /**
   * Request URI Path for this resource. This field will point into static
   * or allocated memory which must remain there for the duration of the
   * resource.
   */
  coap_str_const_t *uri_path;  /**< the key used for hash lookup for this
                                    resource */
  int flags; /**< zero or more COAP_RESOURCE_FLAGS_* or'd together */

  /**
  * The next value for the Observe option. This field must be increased each
  * time the resource changes. Only the lower 24 bits are sent.
  */
  unsigned int observe;

  /**
   * Pointer back to the context that 'owns' this resource.
   */
  coap_context_t *context;

  /**
   * Count of valid names this host is known by (proxy support)
   */
  size_t proxy_name_count;

  /**
   * Array valid names this host is known by (proxy support)
   */
  coap_str_const_t ** proxy_name_list;

  /**
   * This pointer is under user control. It can be used to store context for
   * the coap handler.
   */
  void *user_data;

};
coap_subscription_t * coap_find_observer(coap_resource_t *resource, coap_session_t *session, const coap_binary_t *token);
coap_subscription_t * coap_find_observer_query(coap_resource_t *resource, coap_session_t *session, const coap_string_t *query);
/* End copied */