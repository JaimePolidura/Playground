#include "account.h"

class AccountRepository {
    virtual void save(const Account& account);
    virtual Account& findById(const String& accountId);
    virtual ~AccountRepository() = default;
};