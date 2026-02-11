#include <sysrepo.h>
#include <signal.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static volatile sig_atomic_t exit_app = 0;
static int dry_run = 1;

static void sigint_handler(int signum) {
  (void)signum;
  exit_app = 1;
}

static void run_cmd(const char *fmt, ...) {
  char cmd[512];
  va_list ap;
  va_start(ap, fmt);
  vsnprintf(cmd, sizeof(cmd), fmt, ap);
  va_end(ap);

  if (dry_run) {
    printf("[DRY-RUN] %s\n", cmd);
    return;
  }

  int rc = system(cmd);
  if (rc != 0) {
    fprintf(stderr, "command failed (%d): %s\n", rc, cmd);
  }
}

static int extract_key_value(const char *xpath, const char *key, char *out, size_t out_len) {
  char pattern[64];
  snprintf(pattern, sizeof(pattern), "[%s='", key);
  const char *start = strstr(xpath, pattern);
  if (!start) {
    return 0;
  }
  start += strlen(pattern);
  const char *end = strchr(start, '\'');
  if (!end) {
    return 0;
  }
  size_t n = (size_t)(end - start);
  if (n >= out_len) {
    n = out_len - 1;
  }
  memcpy(out, start, n);
  out[n] = '\0';
  return 1;
}

static void apply_enabled(const char *ifname, int enabled) {
  run_cmd("ip link set dev %s %s", ifname, enabled ? "up" : "down");
}

static void apply_ipv4(const char *ifname, const char *ip, uint8_t prefix, int is_delete) {
  run_cmd("ip addr %s %s/%u dev %s", is_delete ? "del" : "add", ip, prefix, ifname);
}

static int module_change_cb(sr_session_ctx_t *session, uint32_t sub_id, const char *module_name,
                            const char *xpath, sr_event_t event, uint32_t request_id,
                            void *private_data) {
  (void)sub_id;
  (void)module_name;
  (void)xpath;
  (void)request_id;
  (void)private_data;

  if (event != SR_EV_APPLY) {
    return SR_ERR_OK;
  }

  sr_change_iter_t *it = NULL;
  int rc = sr_get_changes_iter(session, "/lab-net-device:interfaces/interface//*", &it);
  if (rc != SR_ERR_OK) {
    fprintf(stderr, "sr_get_changes_iter failed: %s\n", sr_strerror(rc));
    return rc;
  }

  sr_change_oper_t op;
  const sr_val_t *old_val = NULL;
  const sr_val_t *new_val = NULL;

  while ((rc = sr_get_change_next(session, it, &op, &old_val, &new_val)) == SR_ERR_OK) {
    const sr_val_t *val = (op == SR_OP_DELETED) ? old_val : new_val;
    if (!val || !val->xpath) {
      continue;
    }

    char ifname[128];
    if (!extract_key_value(val->xpath, "name", ifname, sizeof(ifname))) {
      continue;
    }

    if (strstr(val->xpath, "/enabled")) {
      if (op == SR_OP_DELETED) {
        continue;
      }
      apply_enabled(ifname, val->data.bool_val ? 1 : 0);
      continue;
    }

    if (strstr(val->xpath, "/ipv4/address") && strstr(val->xpath, "/prefix-length")) {
      char ip[64];
      if (!extract_key_value(val->xpath, "ip", ip, sizeof(ip))) {
        continue;
      }
      apply_ipv4(ifname, ip, val->data.uint8_val, op == SR_OP_DELETED);
      continue;
    }
  }

  if (rc != SR_ERR_NOT_FOUND) {
    fprintf(stderr, "sr_get_change_next failed: %s\n", sr_strerror(rc));
  }

  sr_free_change_iter(it);
  return SR_ERR_OK;
}

int main(void) {
  sr_conn_ctx_t *conn = NULL;
  sr_session_ctx_t *sess = NULL;
  sr_subscription_ctx_t *sub = NULL;

  if (getenv("SIL_LITE_APPLY") != NULL) {
    dry_run = 0;
  }

  signal(SIGINT, sigint_handler);
  signal(SIGTERM, sigint_handler);

  int rc = sr_connect(0, &conn);
  if (rc != SR_ERR_OK) {
    fprintf(stderr, "sr_connect failed: %s\n", sr_strerror(rc));
    return 1;
  }

  rc = sr_session_start(conn, SR_DS_RUNNING, &sess);
  if (rc != SR_ERR_OK) {
    fprintf(stderr, "sr_session_start failed: %s\n", sr_strerror(rc));
    sr_disconnect(conn);
    return 1;
  }

  rc = sr_module_change_subscribe(
      sess,
      "lab-net-device",
      "/lab-net-device:interfaces/interface",
      module_change_cb,
      NULL,
      0,
      SR_SUBSCR_DEFAULT,
      &sub);
  if (rc != SR_ERR_OK) {
    fprintf(stderr, "sr_module_change_subscribe failed: %s\n", sr_strerror(rc));
    sr_session_stop(sess);
    sr_disconnect(conn);
    return 1;
  }

  printf("SIL-lite running. DRY_RUN=%s\n", dry_run ? "1" : "0");
  while (!exit_app) {
    sleep(1);
  }

  sr_unsubscribe(sub);
  sr_session_stop(sess);
  sr_disconnect(conn);
  return 0;
}
