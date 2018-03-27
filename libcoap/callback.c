#include <coap/coap.h>
#include "callback.h"

#include "_cgo_export.h"

void response_handler(coap_context_t *context,
                      coap_session_t *session,
                      coap_pdu_t *sent,
                      coap_pdu_t *received,
                      const coap_tid_t id) {

    export_response_handler(context, session, sent, received, id);
}

void method_handler(coap_context_t *context,
                    coap_resource_t *resource,
                    coap_session_t *session,
                    coap_pdu_t *request,
                    str *token,
                    str *queryString,
                    coap_pdu_t *response) {

    export_method_handler(context, resource, session, request, token, queryString, response);
}

void pong_handler(coap_context_t *context,
                      coap_session_t *session,
                      coap_pdu_t *received,
                      const coap_tid_t id) {

    export_pong_handler(context, session, received, id);
}