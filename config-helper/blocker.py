# -*- coding: utf-8 -*-
import argparse
import yaml
import MySQLdb
from datetime import datetime
from collections import defaultdict

# utils
import utils


def getBlocker(i, cur):
    d = {}
    # get from blocker table
    cur.execute("select type, capacity from blocker where id=%s", (i,))
    t, c = cur.fetchone()
    d['id'] = i
    d['type'] = t
    d['capacity'] = c

    # get from blocker_parameter table
    cur.execute(
        "select `key`, value from blocker_parameter where blocker_id=%s", (i,))
    for row in cur:
        k, v = row
        d[k] = v
    return d


def setBlocker(d, cur):
    # insert into blocker table
    t = datetime.now()
    cur.execute("insert into blocker values(%s, %s, %s, %s, %s, %s)",
                (d['id'], d['type'], d['capacity'], 0, t, t))

    # insert into blocker_parameter table
    for k in ['nextHop', 'host', 'port']:
        cur.execute("insert into blocker_parameter values(%s, %s, %s, %s, %s, %s)",
                    (0, d['id'], k, d[k], t, t))


def delBlocker(i, cur):
    # delete from blocker table
    cur.execute("delete from blocker where id=%d" % i)
    # delete from blocker_parameter table
    cur.execute("delete from blocker_parameter where blocker_id=%d" % i)


if __name__ == "__main__":
    # get command line parameter
    parser = argparse.ArgumentParser(description='argparse blocker')
    parser.add_argument('-m', dest='method',
                        help='method get/set/del', default='get')
    parser.add_argument('-i', dest='id', help='id int', default=0)
    parser.add_argument('-f', dest='filename', help='filename str', default='')
    parser.add_argument('--host', dest='host',
                        help='db ip/hostname', default='127.0.0.1')
    parser.add_argument('--port', dest='port', help='db port', default=3306)
    parser.add_argument('--user', dest='user', help='db user', default='')
    parser.add_argument('--passwd', dest='passwd',
                        help='db passwd', default='')
    parser.add_argument('--db', dest='db', help='db database', default='dots')
    args = parser.parse_args()

    try:
        conn = MySQLdb.connect(
            host=args.host, port=args.port, user=args.user, passwd=args.passwd, db=args.db)
        cur = conn.cursor()

        ids = utils.getRows(cur, col='id', table='blocker')
        i = int(args.id)

        if args.method == "get":
            if i in ids:
                b = [getBlocker(i, cur)]
            elif ids:
                b = [getBlocker(i, cur) for i in ids]
            blockers = {'blocker': b}
            utils.printYamlDump(blockers)

        elif args.method == "set":
            blockers = utils.fileYamlLoad(args.filename)
            for d in blockers['blocker']:
                if d['id'] not in ids:
                    setBlocker(d, cur)
                else:
                    print("Duplicate id entry.")
                    utils.printYamlDump(d)

        elif args.method == "del":
            blockers = {'blocker': []}
            if i in ids:
                 b = getBlocker(i, cur)
                 blockers['blocker'].append(b)
                 delBlocker(i, cur)
            print("Deleted.")
            utils.printYamlDump(blockers)

        conn.commit()

    finally:
        cur.close()
        conn.close()
