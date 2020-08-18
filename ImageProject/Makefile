OBJ = imageSurface.c

CC = gcc

COMPILER_FLAGS = -w

LINKER_FLAGS = -lSDL2 -g -lzmq -lm

OBJ_NAME = imageSurface

all : $(OBJS)
	$(CC) $(OBJ) $(COMPILER_FLAGS) $(LINKER_FLAGS) -o $(OBJ_NAME)
	
clean:
	rm $(OBJ_NAME)

