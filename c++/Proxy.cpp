#include <algorithm>

template<typename T>
struct InProxy {
    const T& value;
};

struct InTag{
    template<typename T>
    InProxy<T> operator < (const T& value) {
        return InProxy<T>(value);
    }

//    template <typename T, typename Range>
//    bool operator > (const InProxy<T>& p, const Range& r) {
//        return std::find(r.begin(), r.end(), p.value) != r.end();
//    }
};
static constexpr InTag in{};