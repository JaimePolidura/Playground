#pragma once

#include "shared.hpp"
#include "types//types.hpp"

namespace VM {
    enum State {
        RUNNING,
        WAITING,
        FINISHED,
    };

    class Thread {
    public:
        uint16_t id;
        State state;

        Types::Object * stack[256];
        uint8_t esp;

        void * gc;
    };
}