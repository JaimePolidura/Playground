#include <iostream>
#include <cstring>

class String {
public:
    String() = default;

    explicit String(const char * source):
            size(strlen(source)), content(new char[strlen(source)])
    {
        printf("created\n");
        memcpy(this->content, source, this->size);
    }

    String(String&& other) noexcept {
        content = other.data();
        size = other.size;

        other.content = nullptr;
        other.size = 0;
    }

    String(String& other):
        size(other.size), content(new char[other.size])
    {
        printf("copied\n");
        memcpy(this->content, other.content, this->size);
    }

    ~ String(){
        delete[] content;
    }

    char * data (){
        return this->content;
    }

private:
    char * content;
    size_t size;
};

void printString(String& string){
    printf("El string es %s\n", string.data());
}

int main(){
    std::string xd = "hola";

    String hola(xd.data());
    printString(hola);

    return 0;
}