#include "modloader.c"
#include "kernelBaseHooks.c"
#include "mempatch.c"
#include "threadSus.cpp"

#include "dx11hook.cpp"

// 007AD330 log address
int attachGameLog = FALSE;
void myCode() {
	AllocConsole();
	freopen_s((FILE**)stdout, "CONOUT$", "w", stdout);
	freopen_s((FILE**)stdin, "CONIN$", "r", stdin);

	Sleep(50);
	printf("Wello Horld~\n");
	HMODULE kernelBase = GetModuleHandleA("kernelbase.dll");
	ModInit();

	tCreateFileA orgCreateFileA = (tCreateFileA)GetProcAddress(kernelBase, "CreateFileA");
	gwCreateFileA = (tCreateFileA)TrampolineHook((BYTE*)orgCreateFileA, (BYTE*)myCreateFileA, 5);

	tReadFile orgReadFile = (tReadFile)GetProcAddress(kernelBase, "ReadFile");
	gwReadFile = (tReadFile)TrampolineHook((BYTE*)orgReadFile, (BYTE*)myReadFile, 5);

	if (attachGameLog){
		tWriteFile orgWriteFile = (tWriteFile)GetProcAddress(kernelBase, "WriteFile");
		gwWriteFile = (tWriteFile)TrampolineHook((BYTE*)orgWriteFile, (BYTE*)myWriteFile, 7);
	}

	HMODULE graphics = NULL;
	while (!graphics) {
		graphics = GetModuleHandleA("d3d11.dll");
		if (graphics != NULL) {
			printf("d3d11 Addr: %p\n", graphics);
			FARPROC createDevicePtr = GetProcAddress(graphics, "D3D11CreateDeviceAndSwapChain");
			printf("D3D11CreateDeviceAndSwapChain Addr: %p\n", createDevicePtr);

			printf("myPresnent %p\n", myPresent);
			// tPresent orgPresent = getPresentProcAddr(myPresent);
			gwPresent = getPresentProcAddr(myPresent);
			printf("Got Present: %p\n", gwPresent);
			// gwPresent = (tPresent)TrampolineHook((BYTE*)orgPresent+5, (BYTE*)myPresent, 6);
		}
		Sleep(25);
	}
}

void xdbg() {
	char xdbgString[64];
	sprintf_s(xdbgString, sizeof(xdbgString), "N:\\RE\\xdbg\\release\\x32\\x32dbg.exe -p %lu", GetCurrentProcessId());
	// system(xdbgString);
	STARTUPINFO si;
	PROCESS_INFORMATION pi;

	CreateProcessA(NULL, xdbgString, NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi);
}

BOOL WINAPI DllMain(HMODULE hModule, DWORD  ul_reason_for_call, LPVOID lpReserved)
{
	switch (ul_reason_for_call)
	{
	case DLL_PROCESS_ATTACH:
	{
		HANDLE handle = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)myCode, NULL, 0, NULL);
		CloseHandle(handle);
		
		int returnCode = MessageBox(NULL, "Game Launch Suspended\nAttach game log to console?", "Loadout Tools", MB_YESNO | MB_ICONQUESTION | MB_DEFBUTTON2);
		if (returnCode == IDYES){
			attachGameLog = TRUE;
		}
		break;
	}
	case DLL_PROCESS_DETACH:
		break;
	case DLL_THREAD_ATTACH:
		break;
	case DLL_THREAD_DETACH:
		break;
	}
	return TRUE;
}
