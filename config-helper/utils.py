# -*- coding: utf-8 -*-
import yaml


def getRows(cur, col, table):
    cur.execute("select %s from %s" % (col, table))
    rows = [row[0] for row in cur]
    return rows


def printYamlDump(d):
    try:
        print(yaml.dump(d, default_flow_style=False))
    except yaml.YAMLError as err:
        print(err)


def fileYamlLoad(filename):
    with open(filename, 'r') as f:
        try:
            a = yaml.load(f)
        except yaml.YAMLError as err:
            print(err)
    return a
