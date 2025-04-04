#include <stdio.h>
#include "mempatch.h"

void writeJumpBytes(BYTE* src, BYTE* dst, size_t size)
{
	DWORD orgProtect;
	VirtualProtect(src, size, PAGE_EXECUTE_READWRITE, &orgProtect);
	
	// set jmp op
	*src = 0xe9;
	uintptr_t relAddress = dst - src - 5;
	*(uintptr_t*)(src + 1) = relAddress;

	if(size > 5){
		memset(src+5, 0x90, size-5);
	}

	VirtualProtect(src, size, orgProtect, &orgProtect);
}

// Gets bytes that will be over written, copy them to
// our allocated memory.
BYTE* TrampolineHook(BYTE* src, BYTE* dst, size_t size)
{
	if (size < 5) {
		printf("Bad JMP Size\n");
		return 0;
	}

	BYTE* gateway = (BYTE*)VirtualAlloc(0, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);

	// Copy org bytes
	memcpy_s(gateway, size, src, size);
	//printf_s("Org bytes: %.*x\n", -size, *src);

	uintptr_t gatewayRel = src - gateway - 5;

	// Jump back to src
	*(gateway+size) = 0xe9;
	*(uintptr_t*)((uintptr_t)gateway + size + 1) = gatewayRel;

	writeJumpBytes(src, dst, size);

	return gateway;
}

