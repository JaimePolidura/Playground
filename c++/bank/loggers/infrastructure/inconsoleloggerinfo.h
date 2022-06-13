#pragma once

#include "../domain/logger.h"

class inconsoleloggerinfo: Logger {
    public: void log(LogLevel level, const String& message) override;
};