#include <vector>
#include <thread>
#include <atomic>
#include <limits>

#define THREADS 8

class Decrementer {
public:
    std::atomic_int32_t * counter;

    Decrementer(std::atomic_int32_t * counter): counter(counter) {}

    void decrement() {
        bool success = this->atomicDecrementRefcount();

        if(success){
            printf("EXITO\n");
        }
    }

private:
    bool atomicDecrementRefcount() const {
        int32_t expected;
        bool success = false;

        do {
            expected = this->counter->load();
        } while (expected > 0 && !(success = this->counter->compare_exchange_strong(expected, expected - 1)));

        return expected - 1 == 0 && success;
    }
};

int main() {
    std::atomic_int32_t * atomic = new std::atomic_int32_t(std::numeric_limits<int32_t>::max() / 4);

    std::vector<std::thread> threads;
    std::vector<Decrementer> decrementers;

    for (int i = 0; i < THREADS; ++i)
        decrementers.push_back(Decrementer{atomic});

    for(int i = 0; i < THREADS; i++)
        threads.emplace_back([decrementers, i]{
            for(int j = 0; j < 100'000'000; j++) {
                auto decrementer = decrementers.at(i);
                decrementer.decrement();
            }
        });

    for(auto & thread : threads)
        thread.join();

    printf("The final counter is %i\n", atomic->load());
}