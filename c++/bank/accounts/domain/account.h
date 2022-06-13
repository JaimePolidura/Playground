#pragma once

#include <iostream>

using String = std::string;

class Account {
    public: Account(const String& account_id_c, const String& nombre_c, const double balance_c);

    const String& getAccountId() const;
    const String& getNombre() const;
    double getBalance() const;

    Account * incrementBalanceBy(unsigned long toIncrement);
    Account * decrementBalanceBy(unsigned long toDecrement);
};