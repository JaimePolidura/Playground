#include "loglevel.h"
#include <iostream>

using String = std::string;

class Logger {
public:
    virtual void log(LogLevel level, const String& message);
};