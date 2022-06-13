#pragma once

#include "account.h"

class AccountRepository {
    public: virtual void save(const Account * account);
    public: virtual Account * findById(const String& accountId);
    public: virtual ~AccountRepository() = default;
};