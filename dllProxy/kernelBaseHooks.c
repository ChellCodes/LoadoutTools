#include <Windows.h>
#include <stdio.h>

typedef HANDLE(__stdcall* tCreateFileA) (
	LPCSTR                lpFileName,
	DWORD                 dwDesiredAccess,
	DWORD                 dwShareMode,
	LPSECURITY_ATTRIBUTES lpSecurityAttributes,
	DWORD                 dwCreationDisposition,
	DWORD                 dwFlagsAndAttributes,
	HANDLE                hTemplateFile
);

HANDLE LogHandle;
tCreateFileA gwCreateFileA;
HANDLE __stdcall myCreateFileA(
	LPCSTR                lpFileName,
	DWORD                 dwDesiredAccess,
	DWORD                 dwShareMode,
	LPSECURITY_ATTRIBUTES lpSecurityAttributes,
	DWORD                 dwCreationDisposition,
	DWORD                 dwFlagsAndAttributes,
	HANDLE                hTemplateFile
){
	HANDLE a = gwCreateFileA(
		lpFileName,
		dwDesiredAccess,
		dwShareMode,
		lpSecurityAttributes,
		dwCreationDisposition,
		dwFlagsAndAttributes,
		hTemplateFile
	);

	if (strstr(lpFileName, "game_log.txt") != NULL){
		printf("Found Log File %p\n", a);
		LogHandle = a;
	}
	if (strstr(lpFileName, ".ini") != NULL){
		printf("Config File Loading: %s\n", lpFileName);
	}
	if (strstr(lpFileName, ".ind") != NULL){
		printf("Ind File Loading: %s\n", lpFileName);
	}
	if (strstr(lpFileName, ".ARC") != NULL){
		printf("Opening Data File: %s\n", lpFileName);
	}
	return a;
}

typedef BOOL(__stdcall* tWriteFile)(
  HANDLE       hFile,
  LPCVOID      lpBuffer,
  DWORD        nNumberOfBytesToWrite,
  LPDWORD      lpNumberOfBytesWritten,
  LPOVERLAPPED lpOverlapped
);

tWriteFile gwWriteFile;
BOOL __stdcall myWriteFile(
  HANDLE       hFile,
  LPCVOID      lpBuffer,
  DWORD        nNumberOfBytesToWrite,
  LPDWORD      lpNumberOfBytesWritten,
  LPOVERLAPPED lpOverlapped
){
	if (hFile == LogHandle){
		hFile = GetStdHandle(STD_OUTPUT_HANDLE);
	}
	return gwWriteFile(
		hFile,
		lpBuffer,
		nNumberOfBytesToWrite,
		lpNumberOfBytesWritten,
		lpOverlapped
	);
}
