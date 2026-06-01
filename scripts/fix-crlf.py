#!/usr/bin/env python3
import pathlib
import sys

root = pathlib.Path(sys.argv[1] if len(sys.argv) > 1 else "/opt/ai-forum")
count = 0
for p in root.rglob("*"):
    if not p.is_file():
        continue
    if p.suffix not in {".sh", ".yml", ".sql", ".example"} and p.name not in {
        "docker-compose.yml",
        "docker-compose.server.yml",
    }:
        continue
    try:
        data = p.read_bytes()
    except OSError:
        continue
    if b"\r" not in data:
        continue
    text = data.decode("utf-8", errors="replace").replace("\r\n", "\n").replace("\r", "\n")
    p.write_text(text, encoding="utf-8", newline="\n")
    count += 1
print(f"fixed {count} files under {root}")
