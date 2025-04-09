#include <stdio.h>
#include <windows.h>

char modFiles[18][1024] = {""};
int filesCount = -1;

void ModInit() {
  {
    WIN32_FIND_DATA fdFile;
    HANDLE hFind = NULL;

    if ((hFind = FindFirstFile("mods\\*.*", &fdFile)) == INVALID_HANDLE_VALUE) {
      printf("Path not found: mods\n");
      return;
    }

    do {
      if (strcmp(fdFile.cFileName, ".") != 0 &&
          strcmp(fdFile.cFileName, "..") != 0) {

        // Is the entity a File or Folder?
        if (fdFile.dwFileAttributes & !FILE_ATTRIBUTE_DIRECTORY) {
        } else {
          filesCount++;
          snprintf(modFiles[filesCount], sizeof(modFiles[filesCount]), "%s",
                   fdFile.cFileName);
        }
      }
    } while (FindNextFile(hFind, &fdFile)); // Find the next file.

    FindClose(hFind); // Always, Always, clean things up!
  }
  printf("Mods Initialized\n");
}

char *GetModFile(LPCSTR fileName) {
  if (filesCount == -1)
    return NULL;

  for (int i = 0; i <= filesCount; i++) {
    if (strstr(fileName, modFiles[i]) != NULL) {
      return modFiles[i];
    }
  }
  return NULL;
}
