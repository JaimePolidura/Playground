#include <iostream>

using String = std::string;

class Account {

private:
    const String& account_id;
    const String& nombre;
    const double balance;

public:
    Account(const String& account_id_c, const String& nombre_c, const double balance_c):
        balance{balance_c}, account_id{account_id_c}, nombre{nombre_c}
    {}

    const String& getAccountId() const {
        return account_id;
    }

    const String& getNombre() const {
        return nombre;
    }

    double getBalance() const {
        return balance;
    }

    Account * incrementBalanceBy(unsigned long toIncrement){
        Account * newAccount = new Account(account_id, nombre, this->balance + toIncrement);
        return newAccount;
    }

    Account * decrementBalanceBy(unsigned long toDecrement){
        Account * newAccount = new Account(account_id, nombre, this->balance - toDecrement);
        return newAccount;
    }

};