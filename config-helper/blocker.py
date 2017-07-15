# -*- coding: utf-8 -*-
import argparse
import yaml
import MySQLdb
import datetime
from collections import defaultdict


def getBlocker(blocker_id, cur):
  blockers = defaultdict(list)

  # get existing blocker_id
  cur.execute("select id from blocker")
  ids = [row[0] for row in cur]
  if blocker_id and blocker_id in ids:
    ids = [blocker_id]
  else:
    print("existing blocker_id:", ids)

  # get each blocker infomation from blocker and blocker_parameter
  for i in ids:
    d={}
    # get blocker
    command = "select type, capacity from blocker where id=%d" % i
    cur.execute(command)
    t, c = cur.fetchone()
    d['id'] = i
    d['type'] = t
    d['capacity'] = c

    # get blocker_parameter
    command = "select `key`, value from blocker_parameter where blocker_id=%d" % i
    cur.execute(command)
    for row in cur:
      k, v = row
      d[k] = v
    blockers['blocker'].append(d)
  return blockers
 



if __name__ == "__main__":
  # get command line parameter
  parser = argparse.ArgumentParser(description = 'argparse blocker')
  parser.add_argument('-m', dest = 'method', help = 'method get/set/del', default = 'get')
  parser.add_argument('-i', dest = 'id', help = 'id int', default = 0)
  parser.add_argument('-f', dest = 'filename', help = 'filename str', default = '')
  parser.add_argument('--host', dest = 'host', help = 'db ip/hostname', default = '127.0.0.1')
  parser.add_argument('--port', dest = 'port', help = 'db port', default = 3306)
  parser.add_argument('--user', dest = 'user', help = 'db user', default = '')
  parser.add_argument('--passwd', dest = 'passwd', help = 'db passwd', default = '')
  parser.add_argument('--db', dest = 'db', help = 'db database', default = 'dots')
  args = parser.parse_args()

  try:
    conn = MySQLdb.connect(host=args.host, port=args.port, user=args.user, passwd=args.passwd, db=args.db)
    cur = conn.cursor()

    print(args)

    if args.method == "get":
      b = getBlocker(int(args.id), cur)
      print(yaml.dump(dict(b), default_flow_style=False))

  finally:
    conn.close()
    cur.close()
