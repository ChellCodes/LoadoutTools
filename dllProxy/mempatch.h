#pragma once
#include <Windows.h>

void writeJumpBytes(BYTE* src, BYTE* dst, size_t size);
BYTE* TrampolineHook(BYTE* src, BYTE* dst, size_t size);
