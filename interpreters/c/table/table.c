#include "table.h"

#define TABLE_MAX_LOAD 0.75

static struct hash_table_entry * find_entry(struct hash_table_entry * entries, int capacity, struct string_object * key);
static void adjust_hash_table_capacity(struct hash_table * table, int new_capcity);

void init_hash_table(struct hash_table * table) {
    table->size = 0;
    table->capacity = 0;
    table->entries = NULL;
}

void add_all_hash_table(struct hash_table * from, struct hash_table * to) {
    for (int i = 0; i < from->capacity; i++) {
        struct hash_table_entry* entry = &from->entries[i];
        if (entry->key != NULL) {
            put_hash_table(to, entry->key, entry->value);
        }
    }
}

bool get_hash_table(struct hash_table * table, struct string_object * key, lox_value_t *value){
    if (table->size == 0) {
        return false;
    }

    struct hash_table_entry * entry = find_entry(table->entries, table->capacity, key);
    if (entry->key == NULL) {
        return false;
    }

    *value = entry->value;

    return true;
}

bool remove_hash_table(struct hash_table * table, struct string_object * key) {
    if (table->size == 0) {
        return false;
    }

    struct hash_table_entry * entry = find_entry(table->entries, table->capacity, key);
    if (entry->key == NULL)  {
        return false;
    }

    entry->key = NULL;
    entry->value = FROM_BOOL(true); //Tombstone, marked as deleted

    return true;
}

bool put_hash_table(struct hash_table * table, struct string_object * key, lox_value_t value) {
    if (table->size + 1 > table->capacity * TABLE_MAX_LOAD) {
        adjust_hash_table_capacity(table, GROW_CAPACITY(table->capacity));
    }

    struct hash_table_entry * entry = find_entry(table->entries, table->capacity, key);
    entry->key = key;
    entry->value = value;

    bool is_new_key = entry->key == NULL;
    bool is_tombstone = entry->key == NULL && IS_BOOL(entry->value) && entry->value.as.boolean;

    if (is_new_key && !is_tombstone) {
        table->size++;
    }

    return is_new_key;
}

static void adjust_hash_table_capacity(struct hash_table * table, int new_capcity) {
    struct hash_table_entry * new_entries = ALLOCATE(struct hash_table_entry, new_capcity);

    for (int i = 0; i < new_capcity; i++) {
        new_entries[i].key = NULL;
    }

    table->size = 0;

    for (int i = 0; i < table->capacity; i++) {
        struct hash_table_entry * entry = &table->entries[i];
        if (entry->key != NULL) {
            struct hash_table_entry * dest = find_entry(new_entries, new_capcity, entry->key);
            dest->key = entry->key;
            dest->value = entry->value;
            table->size++;
        }
    }

    FREE_ARRAY(struct hash_table_entry, table->entries, table->capacity);
    table->entries = new_entries;
    table->capacity = new_capcity;
}

static struct hash_table_entry * find_entry(struct hash_table_entry * entries, int capacity, struct string_object * key) {
    struct hash_table_entry * first_tombstone_found = NULL;
    uint32_t index = key->hash & (capacity - 1); //Optimized %
    for (;;) {
        struct hash_table_entry * entry = &entries[index];
        bool is_tombstone = entry->key == NULL && IS_BOOL(entry->value) && entry->value.as.boolean;

        if(is_tombstone && first_tombstone_found == NULL){
            first_tombstone_found = entry;
            continue;
        }
        if (entry->key == key) {
            return entry;
        }
        if(entry->key == NULL && first_tombstone_found != NULL){
            return first_tombstone_found;
        }
        if(entry->key == NULL && first_tombstone_found == NULL){
            return entry;
        }

        index = (index + 1) & (capacity - 1); //Optimized %
    }
}