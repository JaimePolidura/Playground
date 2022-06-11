cmake_minimum_required(VERSION 3.22)
project(programacpp)

set(CMAKE_CXX_STANDARD 17)

add_executable(c++ main.cpp bank/bank.cpp bank/bank.h bank/accounts/account.cpp bank/accounts/account.h bank/accounts/accountrep.cpp bank/accounts/accountrep.h list/linkedlist.cpp list/linkedlist.h)
