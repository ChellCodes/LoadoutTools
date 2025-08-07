#pragma once
#include <cstdint>

#define CHARACTERBASE 0x00B54270

char* UberentAddress = (char*)0x1015434;
char* UesAddress = (char*)0x0f438b8;
char* MatchmakingAddress = (char*)0x1015540;
// char* MapAddress = (char*)0x0cc94d0;

struct CharSlot{
    uint32_t start;
    uint8_t u1;
    char Name[256];
    uint8_t SelectedChar;
    uint8_t u2[26];
    uint32_t Cosmetic[120];
};
