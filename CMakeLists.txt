cmake_minimum_required(VERSION 3.18.0)

project(
	aur-setup
	LANGUAGES CXX
	VERSION 0.6.0
	DESCRIPTION "AUR git hooks helper"
	HOMEPAGE_URL "https://github.com/AlfredoRamos/snippets/tree/master/cpp/aur-setup"
)

if ("${CMAKE_BINARY_DIR}" STREQUAL "${CMAKE_SOURCE_DIR}")
	message(FATAL_ERROR "In-source builds are not allowed")
endif()

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
add_compile_options(-Wall)

find_package(Threads REQUIRED)
add_subdirectory(libs/inipp)
add_subdirectory(libs/spdlog)
add_subdirectory(libs/fmt)

add_executable("${PROJECT_NAME}")

target_sources(
	"${PROJECT_NAME}"
	PRIVATE
	main.cpp
)

target_link_libraries(
	"${PROJECT_NAME}"
	"${CMAKE_THREAD_LIBS_INIT}"
	inipp::inipp
	spdlog::spdlog
	fmt::fmt
)
