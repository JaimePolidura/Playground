#include "utils.hpp"

static void traverseStruct(Types::StructObject * structObject, std::queue<Types::Object *>& pending);
static void traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending);

std::size_t Types::sizeofObject(Types::ObjectType type) {
    switch (type) {
        case ARRAY:
            return sizeof(Types::ArrayObject);
        case STRUCT:
            return sizeof(Types::StructObject);
        case STRING:
            return sizeof(Types::StringObject);
        default:
            return 0;
    }
}

void traverseObjectDeep(Types::Object * object, std::function<bool(Types::Object *)> callback) {
    std::queue<Types::Object *> pending;
    pending.push(object);

    while(!pending.empty()) {
        Types::Object * currentObject = pending.front();
        pending.pop();

        if(!callback(currentObject)){
            continue;
        }

        switch (currentObject->type) {
            case Types::ObjectType::ARRAY: {
                traverseArray(reinterpret_cast<Types::ArrayObject *>(currentObject), pending);
            }
            case Types::ObjectType::STRUCT: {
                traverseStruct(reinterpret_cast<Types::StructObject *>(currentObject), pending);
            }
            default:
                break;
        }

    }
}

void traverseStruct(Types::StructObject * structObject, std::queue<Types::Object *>& pending) {
    for(auto currentField = structObject->fields;
        currentField < (structObject->fields + structObject->n_fields);
        currentField++) {

        pending.push(currentField);
    }
}

void traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending) {
    for(auto currentElement = arrayObject->content;
        currentElement < (arrayObject->content + arrayObject->size);
        currentElement++) {

        pending.push(currentElement);
    }
}