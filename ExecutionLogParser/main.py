#! /usr/local/bin/python3
# coding: utf-8


start = False
cache = []
with open("/tmp/exec.log.txt", "r") as f:
    for row in f.readlines():
        if start:
            row = row.strip()
            if row.startswith("path:"):
                if row not in cache:
                    cache.append(row)
                else:
                    start = False
                    continue
            if row.startswith(("path:", "hash:", "remote_cache_hit:")):
                print(row)

        if row.startswith("actual_outputs"):
            start = True
        elif row.startswith("remote_cache_hit"):
            print()
            start = False
