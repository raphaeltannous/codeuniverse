FROM gcc:15.2-trixie AS build

WORKDIR /codeuniverse

RUN mkdir -p /codeuniverse/pch /codeuniverse/run

RUN wget -O /codeuniverse/pch/json.hpp https://github.com/nlohmann/json/releases/latest/download/json.hpp

RUN cat > /codeuniverse/pch/precompiled.h << 'EOF'
// C++ Standard Library Headers
#include <algorithm>
#include <any>
#include <array>
#include <atomic>
#include <bit>
#include <bitset>
#include <charconv>
#include <chrono>
#include <climits>
#include <cmath>
#include <compare>
#include <complex>
#include <concepts>
#include <condition_variable>
#include <coroutine>
#include <cstddef>
#include <cstdint>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <deque>
#include <exception>
#include <expected>
#include <filesystem>
#include <format>
#include <forward_list>
#include <fstream>
#include <functional>
#include <initializer_list>
#include <iomanip>
#include <ios>
#include <iosfwd>
#include <iostream>
#include <istream>
#include <iterator>
#include <limits>
#include <list>
#include <map>
#include <memory>
#include <mutex>
#include <numeric>
#include <optional>
#include <ostream>
#include <print>
#include <queue>
#include <random>
#include <ranges>
#include <ratio>
#include <regex>
#include <set>
#include <shared_mutex>
#include <source_location>
#include <span>
#include <spanstream>
#include <sstream>
#include <stack>
#include <stdexcept>
#include <string>
#include <string_view>
#include <system_error>
#include <thread>
#include <tuple>
#include <type_traits>
#include <typeindex>
#include <typeinfo>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <valarray>
#include <variant>
#include <vector>
#include <version>
#include "/codeuniverse/pch/json.hpp"

using namespace std;
EOF

RUN g++ -std=c++23 -O2 -x c++-header \
    /codeuniverse/pch/precompiled.h \
    -o /codeuniverse/pch/precompiled.h.gch

RUN ls -lh /codeuniverse/pch/precompiled.h.gch

WORKDIR /codeuniverse/run
