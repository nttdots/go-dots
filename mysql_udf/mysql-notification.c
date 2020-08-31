#ifdef STANDARD
/* STANDARD is defined. Don't use any MySQL functions */
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#ifdef __WIN__
typedef unsigned __int64 ulonglong;     /* Microsoft's 64 bit types */
typedef __int64 longlong;
#else
typedef unsigned long long ulonglong;
typedef long long longlong;
#endif /*__WIN__*/
#else
#include <string.h>
#include <my_global.h>
#include <my_sys.h>
#endif
#include <mysql.h>
#include <ctype.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

static int _server = -1;

enum TRIGGER_TYPE {
    TRIGGER_CLOSE = 1,
    TRIGGER_UPDATE = 2,
    TRIGGER_INSERT = 3,
    TRIGGER_DELETE = 4
};

#define PORT 9999
#define MITIGATION_SCOPE         "mitigation_scope"
#define SESSION_CONFIGURATION    "signal_session_configuration"
#define PREFIX_ADDRESS_RANGE     "prefix"
#define DATA_ACLS                "data_acls"
#define	TELEMETRY_PRE_MITIGATION "telemetry_pre_mitigation"

#define	TELEMETRY_TRAFFIC                 "telemetry_traffic"
#define	TELEMETRY_TOTAL_ATTACK_CONNECTION "telemetry_total_attack_connection"
#define	TELEMETRY_ATTACK_DETAIL           "telemetry_attack_detail"
#define	TELEMETRY_SOURCE_COUNT            "telemetry_source_count"
#define	TELEMETRY_TOP_TALKER              "telemetry_top_talker"
#define	TELEMETRY_SOURCE_PREFIX           "telemetry_source_prefix"
#define	TELEMETRY_SOURCE_PORT_RANGE       "telemetry_source_port_range"
#define	TELEMETRY_SOURCE_ICMP_TYPE_RANGE  "telemetry_source_icmp_type_range"

#define URI_FILTERING_TRAFFIC                      "uri_filtering_traffic"
#define URI_FILTERING_TRAFFIC_PROTOCOL             "uri_filtering_traffic_per_protocol"
#define URI_FILTERING_TRAFFIC_PORT                 "uri_filtering_traffic_per_port"
#define URI_FILTERING_TOTAL_ATTACK_CONNECTION      "uri_filtering_total_attack_connection"
#define URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT "uri_filtering_total_attack_connection_port"
#define URI_FILTERING_ATTACK_DETAIL                "uri_filtering_attack_detail"
#define URI_FILTERING_SOURCE_COUNT                 "uri_filtering_source_count"
#define URI_FILTERING_TOP_TALKER                   "uri_filtering_top_talker"
#define URI_FILTERING_SOURCE_PREFIX                "uri_filtering_source_prefix"
#define URI_FILTERING_ICMP_TYPE_RANGE              "uri_filtering_icmp_type_range"
#define URI_FILTERING_SOURCE_PORT_RANGE            "uri_filtering_source_port_range"

my_bool MySQLNotification_init(UDF_INIT *initid, 
                                          UDF_ARGS *args,
                                          char *message) {
    // allocate memory here
    // longlong* i = malloc(sizeof(*i));
    //initid->ptr = (char*)i;
    
    struct sockaddr_in remote, saddr;

    if (strcmp((char*)args->args[0], MITIGATION_SCOPE) == 0) {

        // check the arguments format
        if(args->arg_count != 7) {
            strcpy(message, "MySQLNotification() requires exactly seven arguments");
            return 1;
        }

        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT || args->arg_type[2] != INT_RESULT ||
            args->arg_type[3] != STRING_RESULT || args->arg_type[4] != INT_RESULT || args->arg_type[5] != STRING_RESULT || args->arg_type[6] != INT_RESULT) {
            strcpy(message, "MySQLNotification() requires four integers, two strings, and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], PREFIX_ADDRESS_RANGE) == 0) {

        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }

        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT) {
            strcpy(message, "MySQLNotification() requires two integers, and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], DATA_ACLS) == 0) {

        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }

        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT) {
            strcpy(message, "MySQLNotification() requires one integer, and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], SESSION_CONFIGURATION) == 0) {

        // check the arguments format
        if(args->arg_count != 3) {
            strcpy(message, "MySQLNotification() requires exactly three arguments");
            return 1;
        }

        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT || args->arg_type[2] != INT_RESULT) {
            strcpy(message, "MySQLNotification() requires two integers, and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_PRE_MITIGATION) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_TRAFFIC) == 0) {
        // check the arguments format
        if(args->arg_count != 3) {
            strcpy(message, "MySQLNotification() requires exactly three arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != STRING_RESULT || args->arg_type[2] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer, one string and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_TOTAL_ATTACK_CONNECTION) == 0) {
        // check the arguments format
        if(args->arg_count != 3) {
            strcpy(message, "MySQLNotification() requires exactly three arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != STRING_RESULT || args->arg_type[2] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer, one string and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_ATTACK_DETAIL) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_COUNT) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_TOP_TALKER) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_PREFIX) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_PORT_RANGE) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_ICMP_TYPE_RANGE) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC) == 0) {
        // check the arguments format
        if(args->arg_count != 3) {
            strcpy(message, "MySQLNotification() requires exactly three arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != STRING_RESULT || args->arg_type[2] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer, one string and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC_PROTOCOL) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC_PORT) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOTAL_ATTACK_CONNECTION) == 0) {
        // check the arguments format
        if(args->arg_count != 3) {
            strcpy(message, "MySQLNotification() requires exactly three arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != STRING_RESULT || args->arg_type[2] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer, one string and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_ATTACK_DETAIL) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_COUNT) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOP_TALKER) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_PREFIX) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_PORT_RANGE) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else if (strcmp((char*)args->args[0], URI_FILTERING_ICMP_TYPE_RANGE) == 0) {
        // check the arguments format
        if(args->arg_count != 2) {
            strcpy(message, "MySQLNotification() requires exactly two arguments");
            return 1;
        }
        if(args->arg_type[0] != STRING_RESULT || args->arg_type[1] != INT_RESULT ) {
            strcpy(message, "MySQLNotification() requires one integer and table name");
            return 1;
        }
    } else {
        strcpy(message, "MySQLNotification() unknown trigger");
        return 1;
    }
   
    // create a socket that will talk to our node server
    _server = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if(_server == -1) {
       return -1;
    }
    
    // bind to local address
    memset(&saddr, 0, sizeof(saddr));
    saddr.sin_family = AF_INET;
    saddr.sin_port = htons(0);
    saddr.sin_addr.s_addr = inet_addr("127.0.0.1");
    if(bind(_server, (struct sockaddr*)&saddr, sizeof(saddr)) != 0) {
        return -1;
    }

    // connect to server
    memset(&remote, 0, sizeof(remote));
    remote.sin_family = AF_INET;
    remote.sin_port = htons(PORT);
    remote.sin_addr.s_addr = inet_addr("127.0.0.1");
    if(connect(_server, (struct sockaddr*)&remote, sizeof(remote)) != 0) {
        sprintf(message, "Failed to connect to server on port: %d", PORT);
        return -1;
    }  

    return 0;
}


void MySQLNotification_deinit(UDF_INIT *initid) {
    // free any allocated memory here
    //free((longlong*)initid->ptr);
    // close server socket
    if(_server != -1) {
        close(_server);
    }
}

longlong MySQLNotification(UDF_INIT *initid, UDF_ARGS *args,
                           char *is_null, char *error) {
    
    char packet[512];

    if(strcmp((char*)args->args[0], MITIGATION_SCOPE) == 0){

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"id\":\"%lld\", \"cid\":\"%lld\", \"cuid\":\"%s\", \"mid\":\"%lld\", \"cdid\":\"%s\", \"status\":\"%lld\"}", ((char*)args->args[0]),
            *((longlong*)args->args[1]), *((longlong*)args->args[2]), ((char*)args->args[3]), *((longlong*)args->args[4]), ((char*)args->args[5]), *((longlong*)args->args[6]));
    } else if(strcmp((char*)args->args[0], SESSION_CONFIGURATION) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"cid\":\"%lld\", \"sid\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]), *((longlong*)args->args[2]));
    } else if(strcmp((char*)args->args[0], PREFIX_ADDRESS_RANGE) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"cid\":\"%lld\" }", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if(strcmp((char*)args->args[0], DATA_ACLS) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"aclId\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_PRE_MITIGATION) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_TRAFFIC) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"prefix_type\":\"%s\", \"prefix_type_id\":\"%lld\"}", ((char*)args->args[0]), ((char*)args->args[1]), *((longlong*)args->args[2]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_TOTAL_ATTACK_CONNECTION) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"prefix_type\":\"%s\", \"prefix_type_id\":\"%lld\"}", ((char*)args->args[0]), ((char*)args->args[1]), *((longlong*)args->args[2]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_ATTACK_DETAIL) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"mitigation_scope_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_COUNT) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_attack_detail_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_TOP_TALKER) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_attack_detail_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_PREFIX) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_PORT_RANGE) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], TELEMETRY_SOURCE_ICMP_TYPE_RANGE) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"prefix_type\":\"%s\", \"prefix_type_id\":\"%lld\"}", ((char*)args->args[0]), ((char*)args->args[1]), *((longlong*)args->args[2]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC_PROTOCOL) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_pre_mitigation_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TRAFFIC_PORT) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_pre_mitigation_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOTAL_ATTACK_CONNECTION) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"prefix_type\":\"%s\", \"prefix_type_id\":\"%lld\"}", ((char*)args->args[0]), ((char*)args->args[1]), *((longlong*)args->args[2]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOTAL_ATTACK_CONNECTION_PORT) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_pre_mitigation_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_ATTACK_DETAIL) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_pre_mitigation_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_COUNT) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_attack_detail_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_TOP_TALKER) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_attack_detail_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_PREFIX) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_SOURCE_PORT_RANGE) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    } else if (strcmp((char*)args->args[0], URI_FILTERING_ICMP_TYPE_RANGE) == 0) {

        // format a message containing id of row and type of change
        sprintf(packet, "{\"table_trigger\":\"%s\", \"tele_top_talker_id\":\"%lld\"}", ((char*)args->args[0]), *((longlong*)args->args[1]));
    }

    send(_server, packet, strlen(packet), 0);

    return 0;
}


