#pragma once

#include <cstdint>
#include <cstddef>
#include <vector>
#include <string>
#include <map>
#include <atomic>
#include <memory>
#include <algorithm>
#include <queue>
#include <functional>
#include <set>
#include <cmath>
#include <cstring>

#include "params.hpp"

template<typename T>
T roundLessTo8(T input) {
    return input & (~7);
}