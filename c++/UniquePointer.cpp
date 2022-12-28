template <typename T>
class UniquePointer {
public:
    UniquePointer(T * pointer) : pointer(pointer) {}

    UniquePointer(const UniquePointer& other) = delete;

    UniquePointer& operator=(const UniquePointer&) = delete;

    UniquePointer& operator=(UniquePointer&& other) noexcept {
        delete this->pointer;
        this->pointer = other.pointer;
        other.pointer = nullptr;

        return * this;
    }

    UniquePointer(UniquePointer&& other) noexcept {
        this->pointer = other.pointer;
        other.pointer = nullptr;
    }

    ~UniquePointer() {
        delete this->pointer;
    }

private:
    T * pointer;
};