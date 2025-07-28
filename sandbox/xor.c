#include <windows.h>

int xor_encrypt(char* data, size_t len, char key) {
    for (size_t i = 0; i < len; i++){
        data[i] ^= key;
    }
    return 0;
}

int xor_decrypt(char* data, size_t len, char key) {
    for (size_t i = 0; i < len; i++){
        data[i] ^= key;
    }
    return 0;
}

size_t my_strlen(const char* s) {
    size_t i = 0;
    while(s[i]) ++i;
    return i;
}

int main(int argc, char* argv[]) {

    if(argc != 3) {
        printf("Usage: %s <string> 0x<key>\n", argv[0]);
        printf("Example: %s Hello 0xAA\n", argv[0]);
        return 1;
    }

    char* input = argv[1];
    
    

    // Convertir argv[2] (ex: "0xAA") en nombre entier
    char* endptr;
    long key_long = strtol(argv[2], &endptr, 0); // base 0 => détecte automatiquement 0x pour hex
    if (*endptr != '\0' || key_long < 0 || key_long > 0xFF) {
        printf("Invalid hex key: %s\n", argv[2]);
        return 1;
    }


    printf("Xoring: \"%s\" with key: %d\n", input, key_long);


    xor_encrypt(input, my_strlen(input), key_long);
    
    
    printf("%s\n", input);


    xor_decrypt(input, my_strlen(input), key_long);
    printf("Decrypted: %s\n", input);

    return 0;
}

// 0xe2