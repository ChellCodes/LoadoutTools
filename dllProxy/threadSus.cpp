#include <tlhelp32.h>


void SetOtherThreadsSuspended(bool suspend) {
  HANDLE hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPTHREAD, 0);
  if (hSnapshot != INVALID_HANDLE_VALUE) {
    THREADENTRY32 te;
    te.dwSize = sizeof(THREADENTRY32);
    if (Thread32First(hSnapshot, &te)) {
      do {
        if (te.dwSize >= (FIELD_OFFSET(THREADENTRY32, th32OwnerProcessID) +
                          sizeof(DWORD)) &&
            te.th32OwnerProcessID == GetCurrentProcessId() &&
            te.th32ThreadID != GetCurrentThreadId()) {

          HANDLE thread =
              ::OpenThread(THREAD_ALL_ACCESS, FALSE, te.th32ThreadID);
          if (thread != NULL) {
            if (suspend) {
              SuspendThread(thread);
            } else {
              ResumeThread(thread);
            }
            CloseHandle(thread);
          }
        }
      } while (Thread32Next(hSnapshot, &te));
    }
  }
}
