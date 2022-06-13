#include "../domain/accountrepository.h"
#include "../../../list/linkedlist.h"
#include <vector>

class InMemoryAccountRepository: AccountRepository {
    private: Linkedlist<Account> accounts;

    public: void save(const Account * account) override{
        this->accounts.add(* account);
    }

    public: Account * findById(const String& accountId) override {
        for(int i = 0; i < this->accounts.getSize(); i++){
            Account account = this->accounts.get(i);

            if(account.getAccountId().compare(accountId) == 0)
                return new Account(accountId, account.getNombre(), account.getBalance());
        }

        throw std::out_of_range("Account ID not found in array");
    };
};
