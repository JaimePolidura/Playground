#include <iostream>

using String = std::string;

class IdentifierGenerator {
    virtual String& generate();
};