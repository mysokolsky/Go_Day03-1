# Константа имени операционной системы
OS:=$(shell uname -s)

removecash:
ifeq ($(OS), Darwin)
	@echo "\033[33m\n — Свободное место до очистки : $(shell df -h "$$HOME" | grep "$$HOME" | awk '{print($$4)}' | tr 'i' 'B') —\033[0m"
	@echo "\033[31m\n — Очистка...\n\033[0m "
	@make zerocash > /dev/null
	@sleep 1
	@echo "\033[32m — Свободное место после очистки : $(shell df -h "$$HOME" | grep "$$HOME" | awk '{print($$4)}' | tr 'i' 'B') —\n\033[0m"
endif

zerocash:
	rm -rf ~/Library/Developer/CoreSimulator/Devices/*
	rm -f ~/.zsh_history
	rm -rf ~/Library/42_cache/
	rm -rf ~/Library/Caches/*
	rm -rf ~/.Trash/*
	rm -rf ~/Library/Safari/*
	rm -rf ~/.kube/cache/*
	rm -rf ~/Library/Application\ Support/Code/CachedData/*
	rm -rf ~/Library/Application\ Support/Code/User/workspaceStoratge
	rm -rf ~/Library/Containers/com.docker.docker/Data/vms/*
	rm -rf ~/Library/Containers/com.apple.Safari/Data/Library/Caches/*
	rm -rf ~/Library/Caches/*
	rm -rf ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage/
	rm -rf ~/Library/Application\ Support/Slack/Cache/
	rm -rf ~/Library/Application\ Support/Slack/Code\ Cache/
	rm -rf ~/Library/Application\ Support/Code/Cache
	rm -rf ~/Library/Application\ Support/Code/CachedData
	rm -rf ~/Library/Application\ Support/Code/CachedExtension
	rm -rf ~/Library/Application\ Support/Code/CachedExtensions
	rm -rf ~/Library/Application\ Support/Code/CachedExtensionVSIXs
	rm -rf ~/Library/Application\ Support/Code/Code\ Cache
	rm -rf ~/Library/Application\ Support/Code/Service\ Worker/ScriptCache
	rm -rf ~/Library/Application\ Support/Slack/Service\ Worker/ScriptCache
	rm -rf ~/Library/Application\ Support/Telegram\ Desktop/tdata/user_data
	rm -rf ~/Library/Application\ Support/Google/Chrome/Default/Service\ Worker/CacheStorage
	rm -rf ~/Library/Application\ Support/Google/Chrome/Default/Service\ Worker/ScriptCache
	rm -rf ~/Library/Developer/CoreSimulator/Caches
	rm -rf ~/Library/Developer/CoreSimulator/Devices
	rm -rf ~/Library/Developer/Xcode/DerivedData
	# rm -rf ~/Library/Application\ Support/Code/Cache
	rm -rf ~/Library/Containers/com. apple. Safari/Data/Library/Caches
	# rm -rf ~/.Trash
	rm -rf ~/Desktop/Screen*
	rm -rf ~/Desktop/Pre*
	rm -rf ~/Documents/*
	# rm -rf ~/Downloads/*
	rm -rf ~/Movies/*
	rm -rf ~/Music/*
	rm -rf ~/Pictures/*
	rm -rf ~/Library/Application\ Support/Slack/Code\ Cache/
	rm -rf ~/Library/Application\ Support/Slack/Cache/
	rm -rf ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage/
	#rm -rf ~/Library/ApplicationSupport/CrashReporter/*
	#rm -rf ~/Library/Application\ Support/Code/*
	#rm -rf ~/Library/Group\ Containers/*
	rm -rf ~/Library/42_cache/
	rm -rf ~/Library/Caches/CloudKit
	rm -rf ~/Library/Caches/com.apple.akd
	rm -rf ~/Library/Caches/com.apple.ap.adprivacyd
	rm -rf ~/Library/Caches/com.apple.appstore
	rm -rf ~/Library/Caches/com.apple.appstoreagent
	rm -rf ~/Library/Caches/com.apple.cache_delete
	rm -rf ~/Library/Caches/com.apple.commerce
	rm -rf ~/Library/Caches/com.apple.iCloudHelper
	rm -rf ~/Library/Caches/com.apple.imfoundation.IMRemoteURLConnectionAgent
	rm -rf ~/Library/Caches/com.apple.keyboardservicesd
	rm -rf ~/Library/Caches/com.apple.nbagent
	rm -rf ~/Library/Caches/com.apple.nsservicescache.plist
	rm -rf ~/Library/Caches/com.apple.nsurlsessiond
	rm -rf ~/Library/Caches/storeassetd
	rm -rf ~/Library/Caches/com.microsoft.VSCode.ShipIt
	rm -rf ~/Library/Caches/com.microsoft.VSCode
	rm -rf ~/Library/Caches/com.google.SoftwareUpdate
	rm -rf ~/Library/Caches/com.google.Keystone
	rm -rf ~/Library/Caches/com.apple.touristd
	rm -rf ~/Library/Caches/com.apple.tiswitcher.cache
	rm -rf ~/Library/Caches/com.apple.preferencepanes.usercache
	rm -rf ~/Library/Caches/com.apple.preferencepanes.searchindexcache
	rm -rf ~/Library/Caches/com.apple.parsecd
	# rm -rf ~/Library/Caches/
	# rm -rf ~/Library/Caches/
	# rm -rf ~/Library/Caches/
	# rm -rf ~/Library/Caches/
	# rm -rf ~/Library/Caches/
	# rm -rf ~/Library/Caches/
	rm -rf ~/.Trash/*
	rm -rf ~/.kube/cache/*
	rm -rf ~/Library/Containers/com.docker.docker/Data/vms/*
	rm -rf ~/Library/Application\ Support/Firefox/Profiles/hdsrd79k.default-release/storage
	rm -rf ~/Library/42_cache
	rm -rf ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage/
	rm -rf ~/Library/Application\ Support/Slack/Cache/
	rm -rf ~/Library/Application\ Support/Slack/Code\ Cache/
	rm -rf ~/Library/Application\ Support/Code/User/workspaceStorage
	rm -rf ~/Library/Application\ Support/Spotify/PersistentCache
	rm -rf ~/Library/Application\ Support/Telegram\ Desktop/tdata/user_data
	rm -rf ~/Library/Application\ Support/Telegram\ Desktop/tdata/emoji
	rm -rf ~/Library/Application\ Support/Code/Cache/Library/Application\ Support/Code/Cachei
	rm -rf ~/Library/Application\ Support/Code/CacheData
	rm -rf ~/Library/Application\ Support/Code/Cache
	rm -rf ~/Library/Application\ Support/Code/CacheData
	rm -rf ~/Library/Application\ Support/Code/Crashpad/completed
	rm -rf ~/Library/Application\ Support/Code/Service\ Worker/CacheStorage
	rm -rf ~/Library/Application\ Support/BraveSoftware/Brave-Browser/Default/Service\ Worker/CacheStorage/
	rm -rf ~/Library/Application\ Support/Slack/Code\ Cache/* -y
	rm -rf ~/Library/Application\ Support/Slack/Cache/* -y
	rm -rf ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage/* -y
	rm -rf ~/Library/Application\ Support/Google/Chrome/Default/Service\ Worker/CacheStorage/*
	rm -rf ~/Library/Application\ Support/Google/Chrome/Crashpad/completed/*
	# rm -rf ~/Library/Caches/* -y
	rm -rf ~/.Trash/* -y
	rm -rf ~/Library/Safari/* -y
	rm -rf ~/.kube/cache/* -y
	rm -rf ~/Library/Application\ Support/Code/CachedData/* -y
	rm -rf ~/Library/Application\ Support/Code/Crashpad/completed/* -y
	rm -rf ~/Library/Application\ Support/Code/User/workspaceStoratge/* -y
	rm -rf ~/Library/Containers/com.docker.docker/Data/vms/* -y
	rm -rf ~/Library/Containers/com.apple.Safari/Data/Library/Caches/* -y
	rm -rf ~/Library/Containers/org.telegram.desktop/Data/Library/Application\ Support/Telegram\ Desktop/tdata/emoji/* -y
	# rm -rfv ~/Library/Caches/*
	rm -rfv ~/Library/Application\ Support/Slack/Cache/*
	rm -rfv ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage/*
	rm -rfv ~/Library/Group\ Containers/6N38VWS5BX.ru.keepcoder.Telegram/account-570841890615083515/postbox/*
	# rm -rfv ~/Library/Caches
	rm -rfv ~/Library/Application\ Support/Code/Cache
	rm -rfv ~/Library/Application\ Support/Code/CachedData
	rm -rfv ~/Library/Application\ Support/Code/CachedExtension
	rm -rfv ~/Library/Application\ Support/Code/CachedExtensions
	rm -rfv ~/Library/Application\ Support/Code/CachedExtensionVSIXs
	rm -rfv ~/Library/Application\ Support/Code/Code\ Cache
	rm -rfv ~/Library/Application\ Support/Slack/Cache
	rm -rfv ~/Library/Application\ Support/Slack/Code\ Cache
	rm -rfv ~/Library/Application\ Support/Slack/Service\ Worker/CacheStorage
	rm -rfv ~/Library/Application\ Support/Code/User/workspaceStorage