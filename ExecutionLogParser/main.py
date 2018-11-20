#! /usr/local/bin/python3
# coding: utf-8


start = False
cache = []
with open("/tmp/exec.log.txt", "r") as f:
    for row in f.readlines():
        if start:
            for k in ["path", "hash", "remote_cache_hit"]:
                row = row.strip()
                if row.startswith(k+":") and row not in cache:
                    print(row.strip())
                    cache.append(row.strip())

        if row.startswith("actual_outputs"):
            start = True
        elif row.startswith("remote_cache_hit"):
            cache.pop()
            print("\n")
            start = False
