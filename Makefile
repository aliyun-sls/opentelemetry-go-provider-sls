.PHONY: upgrade_main_project_dependencies_version
upgrade_main_project_dependencies_version:
	go get -u ./...

.PHONY: upgrade_example_projects_dependencies_version
upgrade_example_projects_dependencies_version:
	cd examples && \
	for d in */ ; do \
		cd $$d && \
		go get -u ./... && \
		go mod tidy && \
		cd .. ; \
	done

.PHONY: remove_all_example_projects_build_stuff
remove_all_example_projects_build_stuff:
	cd examples && \
	for d in * ; do \
		cd $$d && \
		echo $$d && \
		cd .. ; \
	done

.PHONY: test_all_example_projects_build
test_all_example_projects_build: remove_all_example_projects_build_stuff
	cd examples && \
	for d in * ; do \
		cd $$d && \
		go build -o $$d && \
		cd .. ; \
	done

.PHONY: all
all: upgrade_main_project_dependencies_version upgrade_example_projects_dependencies_version test_all_example_projects_build