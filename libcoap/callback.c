
#include "callback.h"
#include "string.h"

/* Copied from the internal file in the libcoap */
#define COAP_FREE_TYPE(Type, Object) coap_free(Object)
#define LL_FOREACH(head,el)                                                                    \
    LL_FOREACH2(head,el,next)

#define LL_FOREACH2(head,el,next)                                                              \
    for ((el) = (head); el; (el) = (el)->next)

#define COAP_MALLOC_TYPE(Type) \
  ((coap_##Type##_t *)coap_malloc(sizeof(coap_##Type##_t)))

#define min(a,b) ((a) < (b) ? (a) : (b))

#define LL_PREPEND(head,add)                                                                   \
    LL_PREPEND2(head,add,next)

#define LL_PREPEND2(head,add,next)                                                             \
do {                                                                                           \
  (add)->next = (head);                                                                        \
  (head) = (add);                                                                              \
} while (0)

#define LL_DELETE(head,del)                                                                    \
    LL_DELETE2(head,del,next)

#define LL_DELETE2(head,del,next)                                                              \
do {                                                                                           \
  LDECLTYPE(head) _tmp;                                                                        \
  if ((head) == (del)) {                                                                       \
    (head)=(head)->next;                                                                       \
  } else {                                                                                     \
    _tmp = (head);                                                                             \
    while (_tmp->next && (_tmp->next != (del))) {                                              \
      _tmp = _tmp->next;                                                                       \
    }                                                                                          \
    if (_tmp->next) {                                                                          \
      _tmp->next = (del)->next;                                                                \
    }                                                                                          \
  }                                                                                            \
} while (0)

#define LDECLTYPE(x) __typeof(x)
/* End copied */

extern coap_response_t export_response_handler(coap_context_t *ctx,
                            coap_session_t *sess,
                            coap_pdu_t *sent,
                            coap_pdu_t *received,
                            coap_mid_t id);

extern void export_method_handler(coap_context_t *ctx,
                           coap_resource_t *rsrc,
                           coap_session_t *sess,
                           coap_pdu_t *req,
                           coap_string_t *tok,
                           coap_string_t *query,
                           coap_pdu_t *resp);
extern void export_method_from_server_handler(coap_context_t *ctx,
                           coap_resource_t *rsrc,
                           coap_session_t *sess,
                           coap_pdu_t *req,
                           coap_string_t *tok,
                           coap_string_t *query,
                           coap_pdu_t *resp);

extern void export_nack_handler(coap_context_t *ctx,
                    coap_session_t *sess,
                    coap_pdu_t *sent,
                    coap_nack_reason_t reason,
                    coap_mid_t id);

extern void export_event_handler(coap_context_t *ctx,
                    coap_event_t event,
                    coap_session_t *sess);

extern int export_validate_cn_call_back(const char *cn,
                        unsigned depth,
                        coap_strlist_t *cn_list);

coap_response_t response_handler(coap_context_t *context,
                      coap_session_t *session,
                      coap_pdu_t *sent,
                      coap_pdu_t *received,
                      const coap_mid_t id) {

    export_response_handler(context, session, sent, received, id);
}

void method_handler(coap_context_t *context,
                    coap_resource_t *resource,
                    coap_session_t *session,
                    coap_pdu_t *request,
                    coap_string_t *token,
                    coap_string_t *queryString,
                    coap_pdu_t *response) {

    export_method_handler(context, resource, session, request, token, queryString, response);
}

void method_from_server_handler(coap_context_t *context,
                    coap_resource_t *resource,
                    coap_session_t *session,
                    coap_pdu_t *request,
                    coap_string_t *token,
                    coap_string_t *queryString,
                    coap_pdu_t *response) {

    export_method_from_server_handler(context, resource, session, request, token, queryString, response);
}

void nack_handler(coap_context_t *context,
                    coap_session_t *session,
                    coap_pdu_t *sent,
                    coap_nack_reason_t reason,
                    const coap_mid_t id){

    export_nack_handler(context, session, sent, reason, id);
}

void event_handler(coap_context_t *context,
                      coap_event_t event,
                      void *data) {

    export_event_handler(context, event, (coap_session_t *)data);
}

int validate_cn_call_back(const char *cn,
                        const uint8_t *asn1_public_cert,
                        size_t asn1_length,
                        coap_session_t *coap_session,
                        unsigned depth,
                        int validated,
                        void *arg){

    X509 *x509 = d2i_X509(NULL, &asn1_public_cert, (long) asn1_length);
    int valid = 0;
    char *cnt;
    // If the present identifier is CA (depth = 1), the client doesn't validate this idenifier
    if (depth == 1) return 1;

    if (x509) {
        STACK_OF(GENERAL_NAME) *san_list;
        san_list = X509_get_ext_d2i(x509, NID_subject_alt_name, NULL, NULL);
        // If existed the Subject Alternative Name, the client validate DNS-ID/SRV-ID (Subject Alternative Name)
        // Else the client validate CN-ID (Common Name)
        if (san_list) {
            int san_count = sk_GENERAL_NAME_num(san_list);

            for (int n = 0; n < san_count; n++) {
                const GENERAL_NAME * name = sk_GENERAL_NAME_value(san_list, n);

                if (name->type == GEN_DNS) {
                    const char *dns_name = (const char *)ASN1_STRING_get0_data(name->d.dNSName);

                    /* Make sure that there is not an embedded NUL in the dns_name */
                    if (ASN1_STRING_length(name->d.dNSName) != (int)strlen (dns_name))
                        continue;
                    cnt = OPENSSL_strdup(dns_name);
                }
                valid = export_validate_cn_call_back(cnt, depth, (coap_strlist_t*) arg);
                if (valid == 1) return valid;
            }
        } else {
            valid = export_validate_cn_call_back(cn, depth, (coap_strlist_t*) arg);
        }
    }

    if (valid == 0) {
        coap_log(LOG_ERR, "Terminate the communication attempt with a bad certificate error \n");
        coap_session_release(coap_session);
    }

    return valid;
}

int coap_dtls_get_peer_common_name(coap_session_t *session,
                                    char *buf,
                                    size_t buf_len){
    SSL *ssl;
    X509 *cert;
    X509_NAME *name;
    int cn_len;

    coap_openssl_context_t *ctx = (coap_openssl_context_t *)session->context->dtls_context;
    coap_dtls_pki_t *setup_data = &ctx->setup_data;

    if (session->tls == NULL) {
        return -1;
    }
    ssl = (SSL *)session->tls;

    long verify_result = SSL_get_verify_result(ssl);
    switch (verify_result) {
    case X509_V_ERR_CERT_NOT_YET_VALID:
    case X509_V_ERR_CERT_HAS_EXPIRED:
        if (setup_data->allow_expired_certs)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_SELF_SIGNED_CERT_IN_CHAIN:
        if (setup_data->allow_self_signed)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_UNABLE_TO_GET_CRL:
        if (setup_data->allow_no_crl)
            verify_result = X509_V_OK;
        break;
    case X509_V_ERR_CRL_NOT_YET_VALID:
    case X509_V_ERR_CRL_HAS_EXPIRED:
        if (setup_data->allow_expired_crl)
            verify_result = X509_V_OK;
        break;
    default:
        break;
    }
    if (X509_V_OK != verify_result) {
        coap_log(LOG_WARNING, "    %s\n", X509_verify_cert_error_string(verify_result));
        return -1;
    }
    cert = SSL_get_peer_certificate(ssl);
    if (cert == NULL) {
        return -1;
    }

    name = X509_get_subject_name(cert);
    cn_len = X509_NAME_get_text_by_NID(name, NID_commonName, NULL, 0);
    if (cn_len < 0) {
        return -1;
    }
    if (buf_len < (size_t)cn_len + 1) {
        return -1;
    }
    return X509_NAME_get_text_by_NID(name, NID_commonName, buf, buf_len);

}

void coap_set_dirty(coap_resource_t *resource, char *key, int length) {
    if (*key == '\0' && length == 0) {
        coap_resource_notify_observers(resource, NULL);
    } else {
        coap_string_t *query = coap_new_string(length);
        query->s = (uint8_t*)key;
        query->length = (size_t)length;
        coap_resource_notify_observers(resource, query);
    }
}

int coap_check_subscribers(coap_resource_t *resource) {
    return !(resource->subscribers == NULL);
}

int coap_check_dirty(coap_resource_t *resource) {
    return resource->dirty;
}

// Get token from subcribers
char* coap_get_token_subscribers(coap_resource_t *resource) {
    coap_subscription_t *subscriber = resource->subscribers;
    if (subscriber != NULL) {
        return subscriber->token;
    }
    return (char*)0;
}

// Get size of block 2 from subcribers
int coap_get_size_block2_subscribers(coap_resource_t *resource) {
    coap_subscription_t *subscriber = resource->subscribers;
    if (subscriber != NULL) {
        coap_block_t block2 = subscriber->block;
        return block2.szx;
    }
    return 0;
}

// create coap_block_t
coap_block_t coap_create_block(unsigned int num, unsigned int m, unsigned int size) {
   coap_block_t block = { num, m, size };
   return block;
}

// create coap_strlist_t
coap_strlist_t* coap_common_name(coap_strlist_t* head, coap_strlist_t* tail, char* str) {
    coap_strlist_t* element = malloc(sizeof(coap_strlist_t));
    coap_string_t *cstr = coap_new_string(strlen(str));
    cstr->s = (uint8_t*)str;
    cstr->length = strlen(str);
    element->str = cstr;
    element->next = NULL;

    if (head == NULL) {
        head = tail = element;
    } else {
        tail->next = element;
        tail = element;
    }
    return element;
}

// handle release session
void coap_session_handle_release(coap_session_t *session) {
    session->context->handle_event = NULL;
    coap_session_release(session);
}

// Get session from resource
coap_session_t* coap_get_session_from_resource(coap_resource_t *resource) {
    if (resource->subscribers == NULL) {
        return NULL;
    }
    return resource->subscribers->session;
}

// Handle add option
size_t coap_handle_add_option(coap_pdu_t *pdu, uint16_t type, unsigned int val) {
    unsigned char buf[4];
    size_t t;
    t = coap_add_option(pdu, type, coap_encode_var_safe(buf, sizeof(buf), val), buf);
    return t;
}

/* Copied from the internal file in the libcoap */
coap_subscription_t *
coap_add_observer(coap_resource_t *resource,
                  coap_session_t *session,
                  const coap_binary_t *token,
                  coap_string_t *query,
                  int has_block2,
                  coap_block_t block,
                  uint8_t code) {
    coap_subscription_t *s;
    assert( session );

    /* Check if there is already a subscription for this peer. */
    s = coap_find_observer(resource, session, token);
    if (!s) {
        /*
        * Cannot allow a duplicate to be created for the same query as application
        * may not be cleaning up duplicates.  If duplicate found, then original
        * observer is deleted and a new one created with the new token
        */
        s = coap_find_observer_query(resource, session, query);
        if (s) {
            /* Delete old entry with old token */
            coap_binary_t tmp_token = { s->token_length, s->token };
            coap_delete_observer(resource, session, &tmp_token);
            s = NULL;
        }
    }

    /* We are done if subscription was found. */
    if (s) {
        if (s->query)
            coap_delete_string(s->query);
        s->query = query;
        s->code = code;
        return s;
    }

    /* Create a new subscription */
    s = COAP_MALLOC_TYPE(subscription);

    if (!s) {
        /* query is not deleted so it can be used in the calling function
        * which must give up ownership of query only if this function
        * does not return NULL. */
        return NULL;
    }

    // coap_subscription_init(s);
    assert(s);
    memset(s, 0, sizeof(coap_subscription_t));
    s->session = coap_session_reference( session );

    if (token && token->length) {
        s->token_length = token->length;
        memcpy(s->token, token->s, min(s->token_length, 8));
    }

    s->query = query;
    s->has_block2 = has_block2;
    s->block = block;
    s->code = code;

    /* add subscriber to resource */
    LL_PREPEND(resource->subscribers, s);
    coap_log(LOG_DEBUG, "create new subscription\n");
    return s;
}

int coap_delete_observer(coap_resource_t *resource, coap_session_t *session,
                     const coap_binary_t *token) {
    coap_subscription_t *s;

    s = coap_find_observer(resource, session, token);

    if ( s && coap_get_log_level() >= LOG_DEBUG ) {
        char outbuf[2 * 8 + 1] = "";
        unsigned int i;
        for ( i = 0; i < s->token_length; i++ )
            snprintf( &outbuf[2 * i], 3, "%02x", s->token[i] );
        coap_log(LOG_DEBUG, "removed observer with token '%s'\n", outbuf);
    }

    if (resource->subscribers && s) {
        LL_DELETE(resource->subscribers, s);
        coap_session_release( session );
        if (s->query)
            coap_delete_string(s->query);
        coap_free(s);
    }
    return s != NULL;
}

coap_subscription_t * coap_find_observer(coap_resource_t *resource, coap_session_t *session, const coap_binary_t *token) {
    coap_subscription_t *s;

    assert(resource);
    assert(session);

    LL_FOREACH(resource->subscribers, s) {
        if (s->session == session && (!token || (token->length == s->token_length && memcmp(token->s, s->token, token->length) == 0)))
            return s;
    }
    return NULL;
}

coap_subscription_t * coap_find_observer_query(coap_resource_t *resource, coap_session_t *session, const coap_string_t *query) {
    coap_subscription_t *s;

    assert(resource);
    assert(session);

    LL_FOREACH(resource->subscribers, s) {
        if (s->session == session && ((!query && !s->query) || (query && s->query && coap_string_equal(query, s->query))))
            return s;
    }
    return NULL;
}
/* End copied */