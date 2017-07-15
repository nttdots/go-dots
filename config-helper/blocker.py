# -*- coding: utf-8 -*-
import argparse
import yaml
import MySQLdb
from datetime import datetime
from collections import defaultdict

def getBlockerIds(cur):
  # get existing blocker_id
  cur.execute("select id from blocker")
  ids = [row[0] for row in cur]
  return ids

def getBlocker(i, cur):
  d={}
  # get blocker
  cur.execute("select type, capacity from blocker where id=%d" % i)
  t, c = cur.fetchone()
  d['id'] = i
  d['type'] = t
  d['capacity'] = c

  # get blocker_parameter
  cur.execute("select `key`, value from blocker_parameter where blocker_id=%d" % i)
  for row in cur:
    k, v = row
    d[k] = v
  return d

def setBlocker(d, cur):
  # set blocker
  t = datetime.now()
  cur.execute("insert into blocker values(%s, %s, %s, %s, %s, %s)", (d['id'], d['type'], d['capacity'], 0, t, t))

  # set blocker_parameter
  for k in ['nextHop', 'host', 'port']:
    cur.execute("insert into blocker_parameter values(%s, %s, %s, %s, %s, %s)", (0, d['id'], k, d[k], t, t))
  conn.commit()


def printYamlDump(d):
  print(yaml.dump(d, default_flow_style=False))

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

    ids = getBlockerIds(cur)
    i = int(args.id)

    if args.method == "get":
      if i in ids: 
        b = [ getBlocker(i, cur) ]
      elif ids:
        b = [ getBlocker(i, cur) for i in ids]
      blockers = {'blocker': b}
      printYamlDump(blockers)

    elif args.method == "set":
      with open(args.filename, 'r') as f:
        try:
          blockers = yaml.load(f)
        except yaml.YAMLError as err:
          print(err)
      for d in blockers['blocker']:
        if d['id'] not in ids:
          setBlocker(d, cur)
        else:
          print("Duplicate id entry.")
          printYamlDump(d)

  finally:
    conn.close()
    cur.close()
