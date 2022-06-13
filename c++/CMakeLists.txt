cmake_minimum_required(VERSION 3.22)
project(programacpp)

set(CMAKE_CXX_STANDARD 17)

add_executable(c++ main.cpp bank/bank.cpp bank/bank.h bank/accounts/domain/account.cpp bank/accounts/domain/account.h list/linkedlist.h bank/accounts/domain/accountrespostiory.cpp bank/accounts/domain/accountrepository.h bank/accounts/infrastructure/inmemaccountrepository.cpp bank/accounts/infrastructure/inmemaccountrepository.h)
