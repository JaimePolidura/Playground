#pragma once

#include "../shared.h"
#include "../vm/vm.h"
#include "../scanner/scanner.h"
#include "../bytecode.h"
#include "../types/object.h"

bool compile(char * source_code, struct chunk * output_chunk);
