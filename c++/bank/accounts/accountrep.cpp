#include "account.h"
#include <iostream>
#include <cstddef>

using String = std::string;

class AccountRepository {
    private: Account * accounts;
    private: size_t max_size;
    private: int actual_size;

    public: explicit AccountRepository(size_t max_size): actual_size{0}, accounts{}, max_size{max_size}{
//        this->accounts = new Account[max_size];
    }

    public: void addAccount(Account& account){

    }

    public: Account * findAccountById(String& accountId){
        for(int i = 0; i < this->actual_size; i++){
            Account * actual_account = this->accounts + i;

            if(actual_account->getAccountId()->compare(accountId))
                return actual_account;
        }

        return nullptr;
    }

    private: void ensureNotMaxSizeReached() const{
        if(this->actual_size + 1 >= this->max_size)
            throw std::runtime_error("Max account number has been reached");
    }

    private: void ensureAccountNotAlreadyIncluded(Account& account){
//        std::string * accountId = account.getAccountId();
    }
};