#include <iostream>

#ifndef PROGRAMACPP_ACCOUNT_H
#define PROGRAMACPP_ACCOUNT_H

using String = std::string;

class Account {
    private: const String account_id[36];
    private: const String nombre[16];
    private: const double balance;

    public: Account(const String& account_id_c, const String& nombre_c, const double balance_c):
                balance{balance_c}, account_id{account_id_c}, nombre{nombre_c}
    {}

    const std::string *getAccountId() const {
        return account_id;
    }

    const std::string *getNombre() const {
        return nombre;
    }

    double getBalance() const {
        return balance;
    }

};


#endif
