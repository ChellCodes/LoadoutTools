#pragma once
#include <Windows.h>

char* UberentAddress = (char*)0x1015434;
char* UesAddress = (char*)0x0f438b8;
char* MatchmakingAddress = (char*)0x1015540;
// char* MapAddress = (char*)0x0cc94d0;

void writeJumpBytes(BYTE* src, BYTE* dst, size_t size);
BYTE* TrampolineHook(BYTE* src, BYTE* dst, size_t size);
