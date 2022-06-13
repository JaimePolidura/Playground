#pragma once

#include <iostream>

using String = std::string;

class IdentifierGenerator {
public:
    virtual String& generate();
};
