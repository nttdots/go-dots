# -*- coding: utf-8 -*-
import yaml
import MySQLdb
from datetime import datetime


def getRows(cur, col, table):
    cur.execute("select %s from %s" % (col,table))
    rows = [row[0] for row in cur]
    return rows


