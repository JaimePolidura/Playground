#pragma once

#include "shared.hpp"
#include "thread.hpp"
#include "package.hpp"

namespace VM {
    struct VM {
        std::vector<Thread> threads;
        std::map<std::string, std::shared_ptr<Package>> packages;

        void * gc;

        void stopThreadsGC();

        void awakeThreadsGC();
    };
}