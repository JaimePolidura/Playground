//
// Created by polid on 10/06/2022.
//

#include "./accounts/domain/accountrepository.h"

class Bank {
private:
    AccountRepository& accountRepository;

public:
    Bank(AccountRepository& accountRepositoryToAdd): accountRepository{accountRepositoryToAdd} {}

    void transfer(String& fromAccountId, String& toAccountId, unsigned long money){
        this->ensureNotTheSame(fromAccountId, toAccountId);
        Account * fromAccount = this->accountRepository.findById(fromAccountId);
        Account * toAccount = this->accountRepository.findById(toAccountId);
        this->ensureEnoguhBalance(fromAccount, money);

        Account * newFromAccount =  fromAccount->decrementBalanceBy(money);
        Account * newToAccount =  toAccount->incrementBalanceBy(money);


    }
private:

    void ensureNotTheSame(String &from, String &to) {
        if(from.compare(to) == 0)
            throw std::logic_error("You cannot be the same");
    }

    void ensureEnoguhBalance(Account * fromAccount, unsigned long money) {
        if(fromAccount->getBalance() < money)
            throw std::logic_error("You dont have enough money to commit that transfer");
    }
};