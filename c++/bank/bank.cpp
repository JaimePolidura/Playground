#include "./accounts/domain/accountrepository.h"
#include "./identifiers/domain/indentifiergenerator.h"
#include "./loggers/domain/logger.h"

class Bank {
private:
    AccountRepository& accountRepository;
    IdentifierGenerator& identifierGenerator;
    Logger& loggerService;

public:
    Bank(AccountRepository& accountRepository, IdentifierGenerator& identifierGenerator, Logger& loggerService):
        accountRepository{accountRepository}, identifierGenerator{identifierGenerator}, loggerService{loggerService}
    {}

    void transfer(String& fromAccountId, String& toAccountId, unsigned long money){
        this->ensureNotTheSame(fromAccountId, toAccountId);
        Account * fromAccount = this->accountRepository.findById(fromAccountId);
        Account * toAccount = this->accountRepository.findById(toAccountId);
        this->ensureEnoguhBalance(fromAccount, money);

        Account * newFromAccountWithBalanceChanged =  fromAccount->decrementBalanceBy(money);
        Account * newToAccountWithBalanceChanged =  toAccount->incrementBalanceBy(money);

        this->accountRepository.save(newFromAccountWithBalanceChanged);
        this->accountRepository.save(newToAccountWithBalanceChanged);

        this->loggerService.log(LogLevel::SUCCESS, "Transfer made");
    }

    String& addUser(String& name, unsigned long initialBalance){
        String& userAccountId = this->identifierGenerator.generate();
        Account * newAccount = new Account(userAccountId, name, initialBalance);

        this->accountRepository.save(newAccount);

        this->loggerService.log(LogLevel::SUCCESS, "User created");

        return userAccountId;
    }

    void withdraw(String& toAccountId, unsigned long money){
        Account * accountWithdraw = this->accountRepository.findById(toAccountId);
        this->ensureEnoguhBalance(accountWithdraw, money);

        Account * newAccountWithChange = accountWithdraw->incrementBalanceBy(money);
        this->accountRepository.save(newAccountWithChange);

        this->loggerService.log(LogLevel::SUCCESS, "User withdrawn");
    }

    void deposit(String& toAccountId, unsigned long money){
        Account * accountToDeposit = this->accountRepository.findById(toAccountId);
        Account * newAccountWithDeposit = accountToDeposit->incrementBalanceBy(money);

        this->accountRepository.save(newAccountWithDeposit);

        this->loggerService.log(LogLevel::SUCCESS, "User deposited");
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