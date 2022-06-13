#include "../domain/logger.h"

class InConsoleLoggerInfo: Logger {
    public: void log(LogLevel level, const String& message) override {
        printf("[%s] %s", level, message.data());
    }
};