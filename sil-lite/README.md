# SIL-lite (Sysrepo Subscriber)

This is a minimal "SIL" implementation: it listens to Sysrepo changes and applies them to the OS.
It is intentionally small and easy to read.

What it does:
- `/interfaces/interface/enabled` -> `ip link set dev <if> up|down`
- `/interfaces/interface/ipv4/address/prefix-length` -> `ip addr add|del <ip>/<len> dev <if>`

Requirements:
- Sysrepo development libraries
- `ip` command available (iproute2)
- CAP_NET_ADMIN if not in dry-run

Build (example):
```bash
gcc -Wall -Wextra -O2 -o sil-lite sil_lite.c -lsysrepo
```

Run (default is dry-run):
```bash
./sil-lite
```

Apply to system (requires privileges):
```bash
SIL_LITE_APPLY=1 ./sil-lite
```

Then push config from another terminal:
```bash
go run ./cmd/yanglab
```

Notes:
- Only applies changes on `SR_EV_APPLY`.
- It ignores deletes for `enabled` (you can extend if needed).
- This is a teaching example, not production-safe.
