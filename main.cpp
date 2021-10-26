#include <cstdlib>
#include <iostream>
#include <string>
#include <fstream>
#include <filesystem>
#include <vector>

#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>

#include "inipp.h"
#include "spdlog/spdlog.h"
#include "fmt/core.h"

struct App {
	std::string homePath;
	std::string configFile;
	std::vector<std::string> repoList;
};

struct Aur {
	std::string rootPath;
	std::string hooksPath;
	std::vector<std::string> hookList;
};

struct Git {
	std::string name;
	std::string email;
};

void setupGitConfig(const std::string &, const Git &);
void setupGitHooks(const std::string &, const std::vector<std::string> &);

int main() {
	App app;
	Aur aur;
	Git git;

	app.homePath = std::string(std::getenv("HOME"));

	if (app.homePath.empty()) {
		struct passwd *pw = getpwuid(getuid());
		app.homePath = std::string(pw->pw_dir);
	}

	app.configFile = (app.homePath + "/.config/aur-setup.conf");

	inipp::Ini<char> ini;
	std::ifstream input{app.configFile};
	ini.parse(input);

	if (!std::filesystem::exists(app.configFile)) {
		spdlog::warn("Configuration file does not exist");

		ini.sections["aur"]["root_path"];
		ini.sections["aur"]["hooks_path"];
		ini.sections["git"]["name"];
		ini.sections["git"]["email"];

		std::ofstream output{app.configFile};
		ini.generate(output);

		if (std::filesystem::exists(app.configFile)) {
			spdlog::info("Configuration file created at {}", app.configFile);
		}

		return 1;
	}

	if (ini.sections.contains("aur")) {
		inipp::get_value(ini.sections.at("aur"), "root_path", aur.rootPath);
		inipp::get_value(ini.sections.at("aur"), "hooks_path", aur.hooksPath);

		if (aur.rootPath.empty() || aur.hooksPath.empty()) {
			spdlog::error("You need to set the AUR root and hooks path");
			return 1;
		}

		if (aur.rootPath.starts_with('~')) {
			aur.rootPath.replace(0, 1, app.homePath);
		}

		if (aur.hooksPath.starts_with('~')) {
			aur.hooksPath.replace(0, 1, app.homePath);
		}

		for (const std::filesystem::directory_entry &e : std::filesystem::directory_iterator{aur.rootPath}) {
			if (e.is_directory()) {
				app.repoList.push_back(e.path().string());
			}
		}

		for (const std::filesystem::directory_entry &e : std::filesystem::directory_iterator{aur.hooksPath}) {
			if (e.is_regular_file()) {
				aur.hookList.push_back(e.path().string());
			}
		}
	}

	if (ini.sections.contains("git")) {
		inipp::get_value(ini.sections.at("git"), "name", git.name);
		inipp::get_value(ini.sections.at("git"), "email", git.email);

		if (git.name.empty() || git.email.empty()) {
			spdlog::error("You need to set the git user name and email");
			return 1;
		}
	}

	for (const std::string &repo : app.repoList) {
		std::string package{std::filesystem::path(repo).filename()};
		spdlog::info("Processing package {}", package);

		setupGitConfig(repo, git);
		setupGitHooks(repo, aur.hookList);
	}

	return 0;
}

void setupGitConfig(const std::string &repo, const Git &git) {
	if (repo.empty() || git.name.empty() || git.email.empty()) {
		return;
	}

	std::filesystem::current_path(repo);
	std::system(fmt::format("git config --local --replace-all user.name '{}'\n", git.name).c_str());
	std::system(fmt::format("git config --local --replace-all user.email '{}'\n", git.email).c_str());

	spdlog::info(" -> Git configuration [{:s}]", true);
}

void setupGitHooks(const std::string &repo, const std::vector<std::string> &hooks) {
	if (repo.empty() || hooks.empty()) {
		return;
	}

	std::filesystem::current_path(repo);

	for (const std::string &h : hooks) {
		std::filesystem::path hook{h};
		hook.replace_extension("");
		bool success = false;

		try {
			std::string file{std::filesystem::current_path().string() + "/.git/hooks/" + hook.filename().string()};

			// std::filesystem::copy_options::{update_existing,overwrite_existing} still throws error
			if (std::filesystem::exists(file)) {
				std::filesystem::remove(file);
			}

			success = std::filesystem::copy_file(h, file, std::filesystem::copy_options::overwrite_existing);
		} catch (const std::filesystem::filesystem_error &ex) {
			spdlog::error("[{1}] {0}", ex.code().message(), ex.code().value());
		}

		spdlog::info(" -> Git hooks [{:s}]", success);
	}
}
