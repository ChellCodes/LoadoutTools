#include "imgui/backends/imgui_impl_dx11.h"
#include "imgui/backends/imgui_impl_win32.h"
#include "imgui/imgui.h"
#include "mempatch.h"
#include <cstdint>
#include <d3d11.h>
#include <libloaderapi.h>
#include <stdio.h>
#include <winuser.h>
#include "loadoutDefs.hpp"
#pragma comment(lib, "d3d11.lib")

extern IMGUI_IMPL_API LRESULT ImGui_ImplWin32_WndProcHandler(HWND hWnd,
                                                             UINT msg,
                                                             WPARAM wParam,
                                                             LPARAM lParam);

typedef LRESULT(__stdcall *tWndProc)(HWND hWnd, UINT msg, WPARAM wParam,
                                     LPARAM lParam);

tWndProc gwWndProc;
LRESULT WINAPI myWndProc(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam) {
  if (ImGui_ImplWin32_WndProcHandler(hWnd, msg, wParam, lParam))
    return true;
  return gwWndProc(hWnd, msg, wParam, lParam);
}

typedef HRESULT(__stdcall *tPresent)(IDXGISwapChain *, UINT, UINT);
tPresent gwPresent;

ID3D11Device *device = nullptr;
IDXGISwapChain *swap_chain;
ID3D11DeviceContext *context;
ID3D11RenderTargetView *target = nullptr;
ImVec4 clear_color = ImVec4(0.45f, 0.55f, 0.60f, 1.00f);

bool init = false;
bool showMenu = true;
bool showDemoMenu = false;
char* baseAddr = nullptr;
HRESULT __stdcall myPresent(IDXGISwapChain *thisptr, UINT sync, UINT flags) {
  if (!init) {
    baseAddr = (char*)GetModuleHandle(NULL);
    DXGI_SWAP_CHAIN_DESC sd;
    thisptr->GetDesc(&sd);
    thisptr->GetDevice(__uuidof(ID3D11Device), (void **)&device);
    device->GetImmediateContext(&context);

    WNDPROC wndProcPtr =
        (WNDPROC)GetWindowLongPtr(sd.OutputWindow, GWLP_WNDPROC);
    gwWndProc =
        (tWndProc)TrampolineHook((BYTE *)wndProcPtr, (BYTE *)myWndProc, 7);

    context->OMSetRenderTargets(1, &target, nullptr);
    if (!target) {
      ID3D11Texture2D *pBackBuffer = nullptr;
      thisptr->GetBuffer(0, __uuidof(ID3D11Texture2D),
                         reinterpret_cast<void **>(&pBackBuffer));

      device->CreateRenderTargetView(pBackBuffer, nullptr, &target);
      pBackBuffer->Release();
      // Make sure our render target is set, only needed if creating our own, if
      // already exist use original
      context->OMSetRenderTargets(1, &target, nullptr);
    }

    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    ImGuiIO &io = ImGui::GetIO();
    (void)io;
    io.IniFilename = NULL;
    io.ConfigFlags |=
        ImGuiConfigFlags_NavEnableKeyboard; // Enable Keyboard Controls
    io.ConfigFlags |= ImGuiConfigFlags_NavEnableGamepad;
    io.ConfigFlags |= ImGuiConfigFlags_DockingEnable;
    // io.DisplaySize = ImVec2(1280, 720);

    ImGuiStyle* style = &ImGui::GetStyle();
    style->WindowRounding = 3.0f;

    ImGui::StyleColorsDark();
    ImGui_ImplWin32_Init(sd.OutputWindow);
    ImGui_ImplDX11_Init(device, context);
    init = true;
    printf("Imgui Initialized\n");
  }

  if(GetAsyncKeyState(VK_F1) & 1)
    showMenu = !showMenu;
  if(GetAsyncKeyState(VK_F5) & 1)
    showDemoMenu = !showDemoMenu;

  ImGui_ImplDX11_NewFrame();
  ImGui_ImplWin32_NewFrame();
  ImGui::NewFrame();
  if (showDemoMenu)
    ImGui::ShowDemoWindow(&showDemoMenu);

  if (showMenu){
    ImGui::Begin("LoadoutTools", &showMenu, ImGuiWindowFlags_NoCollapse);
    if (ImGui::Button("Patch Endpoints")) {
        strncpy(UberentAddress, "api.loadout.rip", 15 * sizeof(char));
        strncpy(UesAddress, "api.loadout.rip", 15 * sizeof(char));
        strncpy(MatchmakingAddress, "api.loadout.rip", 15 * sizeof(char));
    }
    if (ImGui::CollapsingHeader("Endpoints")) {
      // ImGui::PushStyleVar(Imgui)
      ImGui::Indent();
      ImGui::Text("%s: %p", UberentAddress, (void *)UberentAddress);
      ImGui::Text("%s: %p", UesAddress, (void *)UesAddress);
      ImGui::Text("%s: %p", MatchmakingAddress, (void *)MatchmakingAddress);
      ImGui::Unindent();
    }
    if (ImGui::CollapsingHeader("Mods", ImGuiTreeNodeFlags_DefaultOpen)) {
      if (filesCount > 0) {
      ImGui::Indent();
        for (int i=0; i <= filesCount; i++) {
          ImGui::Text(modFiles[i]);
        }
      ImGui::Unindent();
      }
    }

    uintptr_t* base = (uintptr_t *)((char*)baseAddr+CHARACTERBASE);
    if (*base){
      uintptr_t customCharData = *((uintptr_t*)((*base)+0x8));
      // uint32_t TrueDataSize = *((uintptr_t*)((*base)+0xc));
      uint32_t numOfSlots = *((uintptr_t*)((*base)+0x10));
      if (customCharData){
        for (uint8_t i = 0; i<numOfSlots; i++) {
          CharSlot* slot = (((CharSlot**)customCharData)[i]);
          ImGui::Text("%s: %u",
                  slot->Name,
                  slot->SelectedChar);
        }
      }else {
        ImGui::Text("Custom Characters Not Loaded");
      }
    }

    ImGui::End();
  }
  ImGui::EndFrame();
  ImGui::Render();
  ImGui_ImplDX11_RenderDrawData(ImGui::GetDrawData());

  return gwPresent(thisptr, sync, flags);
}

tPresent getPresentProcAddr(tPresent funcOverride) {
  DXGI_SWAP_CHAIN_DESC sd;
  ZeroMemory(&sd, sizeof(sd));
  sd.BufferCount = 2;
  sd.BufferDesc.Format = DXGI_FORMAT_R8G8B8A8_UNORM;
  sd.BufferUsage = DXGI_USAGE_RENDER_TARGET_OUTPUT;
  sd.OutputWindow = GetForegroundWindow();
  sd.SampleDesc.Count = 1;
  sd.Windowed = TRUE;
  sd.SwapEffect = DXGI_SWAP_EFFECT_DISCARD;

  const D3D_FEATURE_LEVEL feature_levels[] = {
      D3D_FEATURE_LEVEL_11_0,
      D3D_FEATURE_LEVEL_10_0,
  };
  if (D3D11CreateDeviceAndSwapChain(NULL, D3D_DRIVER_TYPE_HARDWARE, NULL, 0,
                                    feature_levels, 2, D3D11_SDK_VERSION, &sd,
                                    &swap_chain, &device, NULL,
                                    &context) == S_OK) {
    tPresent ptr = (tPresent)(*(void ***)swap_chain)[8];

    BYTE *swVTable = (BYTE *)(*(void ***)swap_chain);
    printf("VTable: %p\n", (void *)swVTable);

    // swap_chain->Release();
    // device->Release();
    // context->Release();

    Sleep(1500);
    printf("address of PresentPtr: %p\n", (void *)(swVTable + 0x20));
    DWORD tmpPro;
    VirtualProtect(swVTable + 0x20, sizeof(DWORD), PAGE_EXECUTE_READWRITE,
                   &tmpPro);

    *(uintptr_t *)(swVTable + 0x20) = (uintptr_t)funcOverride;
    // (*reinterpret_cast<void ***>(swap_chain))[8] = (void*)funcOverride;

    VirtualProtect(swVTable + 0x20, sizeof(DWORD), tmpPro, &tmpPro);
    return ptr;
  }
  return NULL;
}

/*
typedef HRESULT(__stdcall *fn_D3D11CreateDeviceAndSwapChain)(
    IDXGIAdapter *, D3D_DRIVER_TYPE, HMODULE, UINT, const D3D_FEATURE_LEVEL *,
    UINT, UINT, const DXGI_SWAP_CHAIN_DESC *, IDXGISwapChain **,
    ID3D11Device **, D3D_FEATURE_LEVEL *, ID3D11DeviceContext **);
*/
