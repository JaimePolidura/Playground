#include "utils.hpp"

static void traverseStruct(Types::StructObject * structObject, std::queue<Types::Object *>& pending);
static void traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending);

void Types::copy(Types::Object * dst, Types::Object * src) {
    *dst = *src;

    switch (src->type) {
        case ARRAY: {
            auto srcArray = AS_ARRAY(src);
            std::memcpy(AS_ARRAY(dst)->elements, srcArray->elements, srcArray->nElements);
            break;
        }
        case STRUCT: {
            auto srcStruct = AS_STRUCT(src);
            std::memcpy(AS_STRUCT(dst)->fields, srcStruct->fields, srcStruct->nFields);
            break;
        }
        default:
            break;
    }
}

std::size_t Types::sizeofObject(Types::Object * object) {
    switch (object->type) {
        case ARRAY:
            return sizeof(Types::ArrayObject) + (sizeof(Types::Object *) * AS_ARRAY(object)->nElements);
        case STRUCT:
            return sizeof(Types::StructObject) + (sizeof(Types::Object *) * AS_STRUCT(object)->nFields) ;
        case STRING:
            return sizeof(Types::StringObject) + std::ceil(AS_STRING(object)->nChars / sizeof(Types::Object));
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
    for(int i = 0; i < structObject->nFields; i++){
        pending.push(structObject->fields[i]);
    }
}

void traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending) {
    for(int i = 0; i < arrayObject->nElements; i++) {
        pending.push(arrayObject->elements[i]);
    }
}