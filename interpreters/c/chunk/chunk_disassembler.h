#pragma once

#include "../shared.h"
#include "chunk.h"
#include "../bytecode.h"

void disassemble_chunk(const struct chunk * chunk, char * name);
int disassemble_chunk_instruction(const struct chunk * chunk, int offset);