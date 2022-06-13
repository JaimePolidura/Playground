cmake_minimum_required(VERSION 3.22)
project(programacpp)

set(CMAKE_CXX_STANDARD 17)

add_executable(c++ main.cpp bank/bank.cpp bank/bank.h bank/accounts/domain/account.cpp bank/accounts/domain/account.h list/linkedlist.h bank/accounts/domain/accountrespostiory.cpp bank/accounts/domain/accountrepository.h bank/accounts/infrastructure/inmemaccountrepository.cpp bank/accounts/infrastructure/inmemaccountrepository.h bank/identifiers/domain/indentifiergenerator.cpp bank/identifiers/domain/indentifiergenerator.h bank/loggers/domain/logger.cpp bank/loggers/domain/logger.h bank/loggers/domain/loglevel.cpp bank/loggers/domain/loglevel.h bank/loggers/infrastructure/inconsoleloggerinfo.cpp bank/loggers/infrastructure/inconsoleloggerinfo.h)
