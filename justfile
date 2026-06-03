BINARY := "transitlab"
SRC_DIR := "src"

build:
  @echo "Building {{BINARY}}..."
  @go build -C {{SRC_DIR}} -o ../{{BINARY}} .
  @echo "Build complete: ./{{BINARY}}"

run: build
  @echo "Running ./{{BINARY}}..."
  @./{{BINARY}}

clean:
  @rm -f {{BINARY}}
  @echo "Cleaned: ./{{BINARY}}"
