template<typename T>
class Iterator {
public:
    virtual bool hasNext() = 0;
    virtual const T& next() = 0;
    virtual size_t size() = 0;
};