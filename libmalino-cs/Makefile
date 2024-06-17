all: clean build

clean:
	@rm -f libmalino-cs.dll libmalino-cs.deps.json libmalino-cs.pdb
	@rm -rf bin
	@rm -rf obj

build:
	@echo " C# libmalino-cs"
	@/home/$(shell whoami)/.dotnet/dotnet publish -c Release -r linux-x64 --verbosity quiet -o . | grep -v "MSBuild version" | cat
	@rm -f libmalino-cs.deps.json libmalino-cs.pdb