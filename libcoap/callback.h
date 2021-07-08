#pragma once

#include <coap3/coap.h>
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <openssl/x509v3.h>

coap_response_t response_handler(coap_session_t *,
                      coap_pdu_t *,
                      coap_pdu_t *,
                      const coap_mid_t);

void method_handler(coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_string_t *,
                    coap_pdu_t *);
void method_from_server_handler(coap_resource_t *,
                    coap_session_t *,
                    coap_pdu_t *,
                    coap_string_t *,
                    coap_pdu_t *);

void nack_handler(coap_session_t *,
                  coap_pdu_t *,
                  coap_nack_reason_t,
                  const coap_mid_t);

void event_handler(void *data, coap_event_t event);

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

int coap_set_dirty(coap_resource_t *resource, char *query, int length);

int coap_check_subscribers(coap_resource_t *resource);
int coap_check_dirty(coap_resource_t *resource);
char* coap_get_token_subscribers(coap_resource_t *resource);
int coap_get_size_block2_subscribers(coap_resource_t *resource);
coap_block_t coap_create_block(unsigned int num, unsigned int m, unsigned int size);

coap_strlist_t* coap_common_name(coap_strlist_t* head, coap_strlist_t* tail, char* str);

void coap_session_handle_release(coap_session_t *session);
size_t coap_handle_add_option(coap_pdu_t *pdu, uint16_t type, unsigned int val);
coap_string_t * coap_get_token_from_request_pdu (coap_pdu_t *pdu);


/* Copied from the internal file in the libcoap */
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

  struct UT_hash_handle *hh;

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

typedef struct UT_hash_bucket {
   struct UT_hash_handle *hh_head;
   unsigned count;

   /* expand_mult is normally set to 0. In this situation, the max chain length
    * threshold is enforced at its default value, HASH_BKT_CAPACITY_THRESH. (If
    * the bucket's chain exceeds this length, bucket expansion is triggered).
    * However, setting expand_mult to a non-zero value delays bucket expansion
    * (that would be triggered by additions to this particular bucket)
    * until its chain length reaches a *multiple* of HASH_BKT_CAPACITY_THRESH.
    * (The multiplier is simply expand_mult+1). The whole idea of this
    * multiplier is to reduce bucket expansions, since they are expensive, in
    * situations where we know that a particular bucket tends to be overused.
    * It is better to let its chain length grow to a longer yet-still-bounded
    * value, than to do an O(n) bucket expansion too often.
    */
   unsigned expand_mult;

} UT_hash_bucket;

/* random signature used only to find hash tables in external analysis */
#define HASH_SIGNATURE 0xa0111fe1u
#define HASH_BLOOM_SIGNATURE 0xb12220f2u

typedef struct UT_hash_table {
   UT_hash_bucket *buckets;
   unsigned num_buckets, log2_num_buckets;
   unsigned num_items;
   struct UT_hash_handle *tail; /* tail hh in app order, for fast append    */
   ptrdiff_t hho; /* hash handle offset (byte pos of hash handle in element */

   /* in an ideal situation (all buckets used equally), no bucket would have
    * more than ceil(#items/#buckets) items. that's the ideal chain length. */
   unsigned ideal_chain_maxlen;

   /* nonideal_items is the number of items in the hash whose chain position
    * exceeds the ideal chain maxlen. these items pay the penalty for an uneven
    * hash distribution; reaching them in a chain traversal takes >ideal steps */
   unsigned nonideal_items;

   /* ineffective expands occur when a bucket doubling was performed, but
    * afterward, more than half the items in the hash had nonideal chain
    * positions. If this happens on two consecutive expansions we inhibit any
    * further expansion, as it's not helping; this happens when the hash
    * function isn't a good fit for the key domain. When expansion is inhibited
    * the hash will still work, albeit no longer in constant time. */
   unsigned ineff_expands, noexpand;

   uint32_t signature; /* used only to find hash tables in external analysis */
#ifdef HASH_BLOOM
   uint32_t bloom_sig; /* used only to test bloom exists in external analysis */
   uint8_t *bloom_bv;
   uint8_t bloom_nbits;
#endif

} UT_hash_table;

typedef struct UT_hash_handle {
   struct UT_hash_table *tbl;
   void *prev;                       /* prev element in app order      */
   void *next;                       /* next element in app order      */
   struct UT_hash_handle *hh_prev;   /* previous hh in bucket order    */
   struct UT_hash_handle *hh_next;   /* next hh in bucket order        */
   const void *key;                  /* ptr to enclosing struct's key  */
   unsigned keylen;                  /* enclosing struct's key len     */
   unsigned hashv;                   /* result of hash-fcn(key)        */
} UT_hash_handle;

struct coap_session_t {
  coap_proto_t proto;               /**< protocol used */
  coap_session_type_t type;         /**< client or server side socket */
  coap_session_state_t state;       /**< current state of relationaship with
                                         peer */
  unsigned ref;                     /**< reference count from queues */
  size_t tls_overhead;              /**< overhead of TLS layer */
  size_t mtu;                       /**< path or CSM mtu */
  struct coap_addr_hash_t *addr_hash;  /**< Address hash for server incoming packets */
  UT_hash_handle hh;
  coap_addr_tuple_t addr_info;      /**< key: remote/local address info */
  int ifindex;                      /**< interface index */
  struct coap_socket_t *sock;               /**< socket object for the session, if
                                         any */
  coap_endpoint_t *endpoint;        /**< session's endpoint */
  coap_context_t *context;          /**< session's context */
  void *tls;                        /**< security parameters */
  uint16_t tx_mid;                  /**< the last message id that was used in
                                         this session */
  uint8_t con_active;               /**< Active CON request sent */
  uint8_t csm_block_supported;      /**< CSM TCP blocks supported */
  coap_mid_t last_ping_mid;         /**< the last keepalive message id that was
                                         used in this session */
  coap_queue_t *delayqueue;         /**< list of delayed messages waiting to
                                         be sent */
  coap_lg_xmit_t *lg_xmit;          /**< list of large transmissions */
  coap_lg_crcv_t *lg_crcv;       /**< Client list of expected large receives */
  coap_lg_srcv_t *lg_srcv;       /**< Server list of expected large receives */
  size_t partial_write;             /**< if > 0 indicates number of bytes
                                         already written from the pdu at the
                                         head of sendqueue */
  uint8_t read_header[8];           /**< storage space for header of incoming
                                         message header */
  size_t partial_read;              /**< if > 0 indicates number of bytes
                                         already read for an incoming message */
  coap_pdu_t *partial_pdu;          /**< incomplete incoming pdu */
  coap_tick_t last_rx_tx;
  coap_tick_t last_tx_rst;
  coap_tick_t last_ping;
  coap_tick_t last_pong;
  coap_tick_t csm_tx;
  coap_dtls_cpsk_t cpsk_setup_data; /**< client provided PSK initial setup
                                         data */
  coap_bin_const_t *psk_identity;   /**< If client, this field contains the
                                      current identity for server; When this
                                      field is NULL, the current identity is
                                      contained in cpsk_setup_data

                                      If server, this field contains the client
                                      provided identity.

                                      Value maintained internally */
  coap_bin_const_t *psk_key;        /**< If client, this field contains the
                                      current pre-shared key for server;
                                      When this field is NULL, the current
                                      key is contained in cpsk_setup_data

                                      If server, this field contains the
                                      client's current key.

                                      Value maintained internally */
  coap_bin_const_t *psk_hint;       /**< If client, this field contains the
                                      server provided identity hint.

                                      If server, this field contains the
                                      current hint for the client; When this
                                      field is NULL, the current hint is
                                      contained in context->spsk_setup_data

                                      Value maintained internally */
  void *app;                        /**< application-specific data */
  uint8_t block_mode;             /**< Zero or more COAP_BLOCK_ or'd options */
  coap_pdu_t *saved_pdu;          /**< Saved PDU when testing for remote
                                       feature */
  coap_fixed_point_t ack_timeout;   /**< timeout waiting for ack
                                         (default 2.0 secs) */
  coap_fixed_point_t ack_random_factor; /**< ack random factor backoff (default
                                             1.5) */
  uint16_t max_retransmit;          /**< maximum re-transmit count
                                         (default 4) */
  uint16_t nstart;                  /**< maximum concurrent confirmable xmits
                                         (default 1) */
  coap_fixed_point_t default_leisure; /**< Mcast leisure time
                                           (default 5.0 secs) */
  uint32_t probing_rate;            /**< Max transfer wait when remote is not
                                         respoding (default 1 byte/sec) */
  uint16_t max_payloads;            /**< maximum Q-BlockX payloads before delay
                                         (default 10) */
  uint16_t non_max_retransmit;      /**< maximum Q-BlockX non re-transmit count
                                         (default 4) */
  coap_fixed_point_t non_timeout;   /**< Q-BlockX timeout waiting for response
                                         (default 2.0 secs) */
  coap_fixed_point_t non_receive_timeout;  /**< Q-BlockX receive timeout before
                                         requesting missing packets.
                                         (default 4.0 secs) */
  coap_fixed_point_t non_probing_wait_base; /**< Q-BlockX max wait time base
                                              while probing
                                             (default 247.0 secs) */
  coap_fixed_point_t non_partial_timeout; /**< Q-BlockX time to wait before
                                           discarding partial data for a body.
                                           (default 247.0 secs) */
  unsigned int dtls_timeout_count;      /**< dtls setup retry counter */
  int dtls_event;                       /**< Tracking any (D)TLS events on this
                                             sesison */
  uint64_t tx_token;              /**< Next token number to use */
  uint64_t tx_rtag;               /**< Next Request-Tag number to use */
};
struct coap_addr_hash_t {
  coap_address_t remote;       /**< remote address and port */
  uint16_t lport;              /**< local port */
  coap_proto_t proto;          /**< CoAP protocol */
};

struct coap_socket_t {
#if defined(WITH_LWIP)
  struct udp_pcb *pcb;
#elif defined(WITH_CONTIKI)
  void *conn;
#else
  coap_fd_t fd;
#endif /* WITH_LWIP */
#if defined(RIOT_VERSION)
  gnrc_pktsnip_t *pkt; /* pointer to received packet for processing */
#endif /* RIOT_VERSION */
  coap_socket_flags_t flags;
  coap_session_t *session; /* Used by the epoll logic for an active session. */
  coap_endpoint_t *endpoint; /* Used by the epoll logic for a listening
                                endpoint. */
};

struct coap_context_t {
  coap_opt_filter_t known_options;
  coap_resource_t *resources; /**< hash table or list of known
                                   resources */
  coap_resource_t *unknown_resource; /**< can be used for handling
                                          unknown resources */
  coap_resource_t *proxy_uri_resource; /**< can be used for handling
                                            proxy URI resources */
  coap_resource_release_userdata_handler_t release_userdata;
                                        /**< function to  release user_data
                                             when resource is deleted */

#ifndef WITHOUT_ASYNC
  /**
   * list of asynchronous message ids */
  coap_async_t *async_state;
#endif /* WITHOUT_ASYNC */

  /**
   * The time stamp in the first element of the sendqeue is relative
   * to sendqueue_basetime. */
  coap_tick_t sendqueue_basetime;
  coap_queue_t *sendqueue;
  coap_endpoint_t *endpoint;      /**< the endpoints used for listening  */
  coap_session_t *sessions;       /**< client sessions */

#ifdef WITH_CONTIKI
  struct uip_udp_conn *conn;      /**< uIP connection object */
  struct etimer retransmit_timer; /**< fires when the next packet must be
                                       sent */
  struct etimer notify_timer;     /**< used to check resources periodically */
#endif /* WITH_CONTIKI */

#ifdef WITH_LWIP
  uint8_t timer_configured;       /**< Set to 1 when a retransmission is
                                   *   scheduled using lwIP timers for this
                                   *   context, otherwise 0. */
#endif /* WITH_LWIP */

  coap_response_handler_t response_handler;
  coap_nack_handler_t nack_handler;
  coap_ping_handler_t ping_handler;
  coap_pong_handler_t pong_handler;

  /**
   * Callback function that is used to signal events to the
   * application.  This field is set by coap_set_event_handler().
   */
  coap_event_handler_t handle_event;

  ssize_t (*network_send)(coap_socket_t *sock, const coap_session_t *session,
                          const uint8_t *data, size_t datalen);

  ssize_t (*network_read)(coap_socket_t *sock, coap_packet_t *packet);

  size_t(*get_client_psk)(const coap_session_t *session, const uint8_t *hint,
                          size_t hint_len, uint8_t *identity,
                          size_t *identity_len, size_t max_identity_len,
                          uint8_t *psk, size_t max_psk_len);
  size_t(*get_server_psk)(const coap_session_t *session,
                          const uint8_t *identity, size_t identity_len,
                          uint8_t *psk, size_t max_psk_len);
  size_t(*get_server_hint)(const coap_session_t *session, uint8_t *hint,
                          size_t max_hint_len);

  void *dtls_context;

  coap_dtls_spsk_t spsk_setup_data;  /**< Contains the initial PSK server setup
                                          data */

  unsigned int session_timeout;    /**< Number of seconds of inactivity after
                                        which an unused session will be closed.
                                        0 means use default. */
  unsigned int max_idle_sessions;  /**< Maximum number of simultaneous unused
                                        sessions per endpoint. 0 means no
                                        maximum. */
  unsigned int max_handshake_sessions; /**< Maximum number of simultaneous
                                            negotating sessions per endpoint. 0
                                            means use default. */
  unsigned int ping_timeout;           /**< Minimum inactivity time before
                                            sending a ping message. 0 means
                                            disabled. */
  unsigned int csm_timeout;           /**< Timeout for waiting for a CSM from
                                           the remote side. 0 means disabled. */
  uint8_t observe_pending;         /**< Observe response pending */
  uint8_t block_mode;              /**< Zero or more COAP_BLOCK_ or'd options */
  uint64_t etag;                   /**< Next ETag to use */

  coap_cache_entry_t *cache;       /**< CoAP cache-entry cache */
  uint16_t *cache_ignore_options;  /**< CoAP options to ignore when creating a
                                        cache-key */
  size_t cache_ignore_count;       /**< The number of CoAP options to ignore
                                        when creating a cache-key */
  void *app;                       /**< application-specific data */
#ifdef COAP_EPOLL_SUPPORT
  int epfd;                        /**< External FD for epoll */
  int eptimerfd;                   /**< Internal FD for timeout */
  coap_tick_t next_timeout;        /**< When the next timeout is to occur */
#endif /* COAP_EPOLL_SUPPORT */
};

struct coap_pdu_t {
  coap_pdu_type_t type;     /**< message type */
  coap_pdu_code_t code;     /**< request method (value 1--31) or response code
                                 (value 64-255) */
  coap_mid_t mid;           /**< message id, if any, in regular host byte
                                 order */
  uint8_t max_hdr_size;     /**< space reserved for protocol-specific header */
  uint8_t hdr_size;         /**< actual size used for protocol-specific
                                 header */
  uint8_t token_length;     /**< length of Token */
  uint16_t max_opt;         /**< highest option number in PDU */
  size_t alloc_size;        /**< allocated storage for token, options and
                                 payload */
  size_t used_size;         /**< used bytes of storage for token, options and
                                 payload */
  size_t max_size;          /**< maximum size for token, options and payload,
                                 or zero for variable size pdu */
  uint8_t *token;           /**< first byte of token, if any, or options */
  uint8_t *data;            /**< first byte of payload, if any */
#ifdef WITH_LWIP
  struct pbuf *pbuf;        /**< lwIP PBUF. The package data will always reside
                             *   inside the pbuf's payload, but this pointer
                             *   has to be kept because no exact offset can be
                             *   given. This field must not be accessed from
                             *   outside, because the pbuf's reference count
                             *   is checked to be 1 when the pbuf is assigned
                             *   to the pdu, and the pbuf stays exclusive to
                             *   this pdu. */
#endif
  const uint8_t *body_data; /**< Holds ptr to re-assembled data or NULL */
  size_t body_length;       /**< Holds body data length */
  size_t body_offset;       /**< Holds body data offset */
  size_t body_total;        /**< Holds body data total size */
  coap_lg_xmit_t *lg_xmit;  /**< Holds ptr to lg_xmit if sending a set of
                                 blocks */
};

struct coap_endpoint_t {
  struct coap_endpoint_t *next;
  coap_context_t *context;        /**< endpoint's context */
  coap_proto_t proto;             /**< protocol used on this interface */
  uint16_t default_mtu;           /**< default mtu for this interface */
  coap_socket_t sock;             /**< socket object for the interface, if
                                       any */
  coap_address_t bind_addr;       /**< local interface address */
  coap_session_t *sessions;       /**< hash table or list of active sessions */
};

struct coap_queue_t {
  struct coap_queue_t *next;
  coap_tick_t t;                /**< when to send PDU for the next time */
  unsigned char retransmit_cnt; /**< retransmission counter, will be removed
                                 *    when zero */
  unsigned int timeout;         /**< the randomized timeout value */
  coap_session_t *session;      /**< the CoAP session */
  coap_mid_t id;                /**< CoAP message id */
  coap_pdu_t *pdu;              /**< the CoAP PDU to send */
};

struct coap_lg_xmit_t {
  struct coap_lg_xmit_t *next;
  uint8_t blk_size;      /**< large block transmission size */
  uint16_t option;       /**< large block transmisson CoAP option */
  int last_block;        /**< last acknowledged block number */
  const uint8_t *data;   /**< large data ptr */
  size_t length;         /**< large data length */
  size_t offset;         /**< large data next offset to transmit */
  union {
    struct coap_l_block1_t *b1;
    struct coap_l_block2_t *b2;
  } b;
  coap_pdu_t pdu;        /**< skeletal PDU */
  coap_tick_t last_payload; /**< Last time MAX_PAYLOAD was sent or 0 */
  coap_tick_t last_used; /**< Last time all data sent or 0 */
  coap_release_large_data_t release_func; /**< large data de-alloc function */
  void *app_ptr;         /**< applicaton provided ptr for de-alloc function */
};

/**
 * Structure to hold large body (many blocks) client receive information
 */
struct coap_lg_crcv_t {
  struct coap_lg_crcv_t *next;
  uint8_t observe[3];    /**< Observe data (if set) (only 24 bits) */
  uint8_t observe_length;/**< Length of observe data */
  uint8_t observe_set;   /**< Set if this is an observe receive PDU */
  uint8_t etag_set;      /**< Set if ETag is in receive PDU */
  uint8_t etag_length;   /**< ETag length */
  uint8_t etag[8];       /**< ETag for block checking */
  uint16_t content_format; /**< Content format for the set of blocks */
  uint8_t last_type;     /**< Last request type (CON/NON) */
  uint8_t initial;       /**< If set, has not been used yet */
  uint8_t szx;           /**< size of individual blocks */
  size_t total_len;      /**< Length as indicated by SIZE2 option */
  coap_binary_t *body_data; /**< Used for re-assembling entire body */
  coap_binary_t *app_token; /**< app requesting PDU token */
  uint8_t base_token[8]; /**< established base PDU token */
  size_t base_token_length; /**< length of token */
  uint8_t token[8];      /**< last used token */
  size_t token_length;   /**< length of token */
  coap_pdu_t pdu;        /**< skeletal PDU */
  struct coap_rblock_t *rec_blocks; /** < list of received blocks */
  coap_tick_t last_used; /**< Last time all data sent or 0 */
  uint16_t block_option; /**< Block option in use */
};

/**
 * Structure to hold large body (many blocks) server receive information
 */
struct coap_lg_srcv_t {
  struct coap_lg_srcv_t *next;
  uint8_t observe[3];    /**< Observe data (if set) (only 24 bits) */
  uint8_t observe_length;/**< Length of observe data */
  uint8_t observe_set;   /**< Set if this is an observe receive PDU */
  uint8_t rtag_set;      /**< Set if RTag is in receive PDU */
  uint8_t rtag_length;   /**< RTag length */
  uint8_t rtag[8];       /**< RTag for block checking */
  uint16_t content_format; /**< Content format for the set of blocks */
  uint8_t last_type;     /**< Last request type (CON/NON) */
  uint8_t szx;           /**< size of individual blocks */
  size_t total_len;      /**< Length as indicated by SIZE1 option */
  coap_binary_t *body_data; /**< Used for re-assembling entire body */
  size_t amount_so_far;  /**< Amount of data seen so far */
  coap_resource_t *resource; /**< associated resource */
  coap_str_const_t *uri_path; /** set to uri_path if unknown resource */
  struct coap_rblock_t *rec_blocks; /** < list of received blocks */
  uint8_t last_token[8]; /**< last used token */
  size_t last_token_length; /**< length of token */
  coap_mid_t last_mid;   /**< Last received mid for this set of packets */
  coap_tick_t last_used; /**< Last time data sent or 0 */
  uint16_t block_option; /**< Block option in use */
};

struct coap_lg_range {
  uint32_t begin;
  uint32_t end;
};

#define COAP_RBLOCK_CNT 4
typedef struct coap_rblock_t {
  uint32_t used;
  uint32_t retry;
  struct coap_lg_range range[COAP_RBLOCK_CNT];
  coap_tick_t last_seen;
} coap_rblock_t;

/**
 * Structure to keep track of block1 specific information
 * (Requests)
 */
typedef struct coap_l_block1_t {
  coap_binary_t *app_token; /**< original PDU token */
  uint8_t token[8];      /**< last used token */
  size_t token_length;   /**< length of token */
  uint32_t count;        /**< the number of packets sent for payload */
} coap_l_block1_t;

/**
 * Structure to keep track of block2 specific information
 * (Responses)
 */
typedef struct coap_l_block2_t {
  coap_resource_t *resource; /**< associated resource */
  coap_string_t *query;  /**< Associated query for the resource */
  uint64_t etag;         /**< ETag value */
  coap_time_t maxage_expire; /**< When this entry expires */
} coap_l_block2_t;

struct coap_async_t {
  struct coap_async_t *next; /**< internally used for linking */
  coap_tick_t delay;    /**< When to delay to before triggering the response
                             0 indicates never trigger */
  coap_session_t *session;         /**< transaction session */
  coap_pdu_t *pdu;                 /**< copy of request pdu */
  void* appdata;                   /** User definable data pointer */
};

struct coap_cache_entry_t {
  UT_hash_handle hh;
  coap_cache_key_t *cache_key;
  coap_session_t *session;
  coap_pdu_t *pdu;
  void* app_data;
  coap_tick_t expire_ticks;
  unsigned int idle_timeout;
  coap_cache_app_data_free_callback_t callback;
};
/* End copied */