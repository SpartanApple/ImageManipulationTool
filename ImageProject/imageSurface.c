#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#include <SDL2/SDL.h>
#include <endian.h>
#include <math.h>

#include <arpa/inet.h> //Hopefully can be removed later, used to ntohl
#include <netinet/in.h>

#include <zmq.h>
#include <string.h>
#include <unistd.h> // Used for access as well as other stuff
#include <assert.h>




//#define WIDTH 500
//#define HEIGHT 400


struct pixel {
	uint8_t red;
	uint8_t green;
	uint8_t blue;
};

// http://www.ece.ualberta.ca/~elliott/ee552/studentAppNotes/2003_w/misc/bmp_file_format/bmp_file_format.htm
struct fileHeader {
	char sig[2];
	uint32_t fileSize;
	uint32_t reserved;
	uint32_t offSet;

	uint32_t size;
	uint32_t width;
	uint32_t height;
	uint16_t planes;
	uint16_t bitsPerPixel;
	uint32_t compression;
	uint32_t imageSize;
	uint32_t xPixelsPerM;
	uint32_t yPixelsPerM;
	uint32_t coloursUsed;
	uint32_t importantColours;
};


int drawCircle(struct pixel* pixelArr, struct pixel penColor, int mouseX, int mouseY, int penWidth, uint32_t width, uint32_t height)
{
	for(int x = mouseX-penWidth; x < mouseX + (2 * penWidth); x++)
	{
		for(int y = mouseY-penWidth; y < mouseY + (2 * penWidth); y++)
		{
			if(((y-mouseY)*(y-mouseY)) + ((x-mouseX)*(x-mouseX)) < penWidth*penWidth)
			{
				if(x + (y*width) >= 0 && (x+(y*width) < width*height)) // Makes sure pen is within boundary of array (Avoids seg-fault)
				{
					// Says which pixel colour to set each of the RGB values
					pixelArr[x + (y*width)].red = penColor.red;
					pixelArr[x + (y*width)].green = penColor.green;
					pixelArr[x + (y*width)].blue = penColor.blue;
				}
			}
		}
	}


}


//https://www.tutorialspoint.com/computer_graphics/line_generation_algorithm.htm (Derived from)
int drawLine(struct pixel* pixelArr, struct pixel penColor, int mouseX, int mouseY, int previousX, int previousY, uint32_t width)
{
	int dx = mouseX - previousX;
	int dy = mouseY - previousY;

	float steps = abs(dx) > abs(dy) ? abs(dx) : abs(dy);

	float xIncr = dx / steps;
	float yIncr = dy / steps;

	float x = mouseX;
	float y = mouseY;

	for(int v = 0; v < steps; v++)
	{
		x += xIncr;
		y += yIncr;
		pixelArr[(int)floorf(y * width + x)].red = penColor.red;
		pixelArr[(int)floorf(y * width + x)].green = penColor.green;
		pixelArr[(int)floorf(y * width + x)].blue = penColor.blue;
	}
	return;
}

void headerSave(struct pixel* pixelArr, struct fileHeader header, uint32_t  height, uint32_t width)
{
	FILE *fp;
	//fp = fopen("scr/bmpReader/forGoHeader.CMP", "wb+");
	fp = fopen("forGoHeader.CMP", "wb+");

	fwrite(&header.sig, sizeof(header.sig), 1, fp); 
	fwrite(&header.fileSize, sizeof(header.fileSize), 1, fp); 
	fwrite(&header.reserved, sizeof(header.reserved), 1, fp); 
	fwrite(&header.offSet, sizeof(header.offSet), 1, fp); 
	fwrite(&header.size, sizeof(header.size), 1, fp); 
	fwrite(&header.width, sizeof(header.width), 1, fp); 
	fwrite(&header.height, sizeof(header.height), 1, fp); 
	fwrite(&header.planes, sizeof(header.planes), 1, fp); 
	fwrite(&header.bitsPerPixel, sizeof(header.bitsPerPixel), 1, fp); 
	fwrite(&header.compression, sizeof(header.compression), 1, fp); 
	fwrite(&header.imageSize, sizeof(header.imageSize), 1, fp); 
	fwrite(&header.xPixelsPerM, sizeof(header.xPixelsPerM), 1, fp); 
	fwrite(&header.yPixelsPerM, sizeof(header.yPixelsPerM), 1, fp); 
	fwrite(&header.coloursUsed, sizeof(header.coloursUsed), 1, fp); 
	fwrite(&header.importantColours, sizeof(header.importantColours), 1, fp); 
	fclose(fp);
	return;
}


void headerSaveNoLoad(struct pixel *pixelArr, uint32_t height, uint32_t width)
{
	FILE *fp;
	fp = fopen("forGoHeader.CMP", "wb+");

	fwrite("BM", 2, 1, fp);
	uint32_t fs = 600122; // All values gotten from header from custom made bitmap
	fwrite(&fs, sizeof(fs), 1, fp); 
	uint32_t res = 0;
	fwrite(&res, sizeof(res), 1, fp); 
	uint32_t off = 122;
	fwrite(&off, sizeof(off), 1, fp);
	uint32_t sz = 108;
	fwrite(&sz, sizeof(sz), 1, fp); 
	fwrite(&width, sizeof(width), 1, fp); 
	fwrite(&height, sizeof(height), 1, fp); 
	uint16_t pl = 1;
	fwrite(&pl, sizeof(pl), 1, fp); 
	uint16_t bPP = 24;
	fwrite(&bPP, sizeof(bPP), 1, fp); 
	uint32_t comp = 0;
	fwrite(&comp, sizeof(comp), 1, fp); 
	uint32_t is = 500;
	fwrite(&is, sizeof(is), 1, fp); 
	uint32_t x = 11811;
	fwrite(&x, sizeof(x), 1, fp); 
	uint32_t y = 11811;
	fwrite(&y, sizeof(y), 1, fp); 
	uint32_t coloursU = 0;
	fwrite(&coloursU, sizeof(coloursU), 1, fp); 
	uint32_t iC = 0;
	fwrite(&iC, sizeof(iC), 1, fp); 

	fclose(fp);
	return;
}
// https://stackoverflow.com/questions/8465006/how-do-i-concatenate-two-strings-in-c (taken from)
char * concat(const char *s1, const char *s2)
{
	char *result = malloc(strlen(s1) + strlen(s2) + 1);

	strcpy(result, s1);
	strcat(result, s2);
	return result;
}

void dataSave(struct pixel* pixelArr, struct fileHeader header, uint32_t height, uint32_t width)
{
	FILE *fp;
	fp = fopen("forGoData.CMP", "wb+");

	for(int y = height-1; y >= 0; y--)
	{ 
		for(int x = 0; x < width; x++)
		{
			fwrite(&pixelArr[(y*width) + x].red, sizeof(uint8_t), 1, fp);
			fwrite(&pixelArr[(y*width) + x].green, sizeof(uint8_t), 1, fp);
			fwrite(&pixelArr[(y*width) + x].blue, sizeof(uint8_t), 1, fp);
		}
	}

	
	fclose(fp);
	return;
	

}

struct fileHeader readHeader(struct fileHeader header, char *fileNameHeader)
{
	
	FILE *fp;
	fp = fopen(fileNameHeader, "r");

	fread(header.sig, sizeof(header.sig), 1, fp);
	fread(&header.fileSize, sizeof(header.fileSize), 1, fp);
	fread(&header.reserved, sizeof(header.reserved), 1, fp);
	fread(&header.offSet, sizeof(header.offSet), 1, fp);
	fread(&header.size, sizeof(header.size), 1, fp);
	fread(&header.width, sizeof(header.width), 1, fp);
	fread(&header.height, sizeof(header.height), 1, fp);
	fread(&header.planes, sizeof(header.planes), 1, fp);
	fread(&header.bitsPerPixel, sizeof(header.bitsPerPixel), 1, fp);
	fread(&header.compression, sizeof(header.compression), 1, fp);
	fread(&header.imageSize, sizeof(header.imageSize), 1, fp);
	fread(&header.xPixelsPerM, sizeof(header.xPixelsPerM), 1, fp);
	fread(&header.yPixelsPerM, sizeof(header.yPixelsPerM), 1, fp);
	fread(&header.coloursUsed, sizeof(header.coloursUsed), 1, fp);
	fread(&header.importantColours, sizeof(header.importantColours), 1, fp);

	fclose(fp);
	return header;
}

void readData(struct pixel *pixelArr, struct fileHeader header, char *fileNameData)
{
	FILE *fp;
	fp = fopen(fileNameData, "r");

	
	for(int y = header.height-1; y >= 0; y--)
	{ 
		for(int x = 0; x < header.width; x++)
		{
			fread(&pixelArr[(y*header.width) + x].red, sizeof(uint8_t), 1, fp);
			fread(&pixelArr[(y*header.width) + x].green, sizeof(uint8_t), 1, fp);
			fread(&pixelArr[(y*header.width) + x].blue, sizeof(uint8_t), 1, fp);
		}
	}
	return;
}


int main(int argc, char **argv)
{
	// Default size
    int width = 500;
    int height = 400;
	
    char *fileNameData = "forCData.CMP";
    char *fileNameHeader = "forCHeader.CMP";

    SDL_Window *window;
    SDL_Renderer *renderer;
    SDL_Texture *texture;
    SDL_Event event;

	// Communication sockets set-up
	//https://zeromq.org/get-started/?language=c# (Reference)
    void *context = zmq_ctx_new();
    void *responder = zmq_socket(context, ZMQ_REP);
    int rc = zmq_bind(responder, "tcp://*:5555");
    assert(rc == 0);

    bool leftMouseButtonDown = false;

    SDL_Init(SDL_INIT_VIDEO);

    // https://wiki.libsdl.org/SDL_CreateWindow

    window = SDL_CreateWindow("Image Manipulation Tool",
							  SDL_WINDOWPOS_UNDEFINED,
							  SDL_WINDOWPOS_UNDEFINED,
							  width, height,
							  SDL_WINDOW_OPENGL | SDL_WINDOW_RESIZABLE);

    // https://wiki.libsdl.org/SDL_CreateRenderer
    renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC);

    // https://wiki.libsdl.org/SDL_CreateTexture
    texture = SDL_CreateTexture(renderer,
				SDL_PIXELFORMAT_RGBA8888,
				SDL_TEXTUREACCESS_TARGET,
				width, height);
    bool running = true;
	int size = 300;
	char example[size]; // Used for buffer read in

	memset(example, 0, size);

	struct pixel* pixelArr = (struct pixel*)malloc(width * height * sizeof(struct pixel));

	// Fills array with only white pixels
	for(long pos = 0; pos < width*height; pos++)
	{
		pixelArr[pos].red   = 255;
		pixelArr[pos].green = 255;
		pixelArr[pos].blue  = 255;
	}

	struct pixel penColor = {0, 0, 0};
	penColor.red = 0;
	penColor.green = 0;
	penColor.blue = 0;

	int previousY = -1;
	int previousX = -1;

	int penWidth = 1;

	char buffer[size];
		
	int whichButton = 0;
	
	int loop = 0;

	int hasLoaded = 0; // Is from a loaded file flag

	struct fileHeader header;

    while(running)
    {
		
		// Process events
        while(SDL_PollEvent(&event))
        {
            if(event.type == SDL_QUIT)
            {
                running = false;
            }
        }
        // Clear screen
        SDL_RenderClear(renderer);


		memset(buffer, 0, size);
		

		zmq_recv(responder, buffer, size, ZMQ_NOBLOCK);
		// Handles python Input
		if(strcmp(buffer, example) != 0)
		{
			//http://www.cplusplus.com/reference/cstdlib/strtol/
			char *pEnd;
			whichButton = strtol(buffer, &pEnd, 10);
			if(whichButton != 2) // Handles the fact we need to return a value to python
			{
				zmq_send(responder, "World", 5, 0);
			}
			if(whichButton == 0)
			{
				penColor.red = strtol(pEnd, &pEnd, 10);
				penColor.green = strtol(pEnd, &pEnd, 10);
				penColor.blue = strtol(pEnd, &pEnd, 10);
			}
			else if(whichButton == 1) // Pen Width
			{
				penWidth = strtol(pEnd, &pEnd, 10);
	
			}
			else if(whichButton == 2) // Eye dropper
			{
				loop = 0; // Value set so eye dropper only functions once
			}
			else if(whichButton == 3) // Eraser
			{
				penColor.red   = 255;
				penColor.green = 255;
				penColor.blue  = 255;
			}
			else if(whichButton == 4)
			{
				
				if(hasLoaded == 1)
				{
					headerSave(pixelArr, header, height, width);
				}
				else
				{
					headerSaveNoLoad(pixelArr, height, width);
				}
				dataSave(pixelArr, header, height, width);
				

				int status = system("scr/bmpReader/saveBMPreadCMP output.bmp");
			}
			
			else if(whichButton == 5) // Reading in
			{
				char * fileLocation = buffer+2;
				char * fileWithArgs = "scr/bmpReader/readBMPSaveCMP ";
				char *fileLocationandArgs = concat(fileWithArgs, fileLocation);
				int status = system(fileLocationandArgs);
			
				//struct fileHeader* header = (struct fileHeader*)malloc(sizeof(fileHeader));
				header = readHeader(header, fileNameHeader);

				// https://www.tutorialspoint.com/c_standard_library/c_function_realloc.htm
				pixelArr = (struct pixel*)realloc(pixelArr, sizeof(struct pixel) * header.width * header.height);
				
				readData(pixelArr, header, fileNameData);			
				// https://wiki.libsdl.org/SDL_SetWindowSize
				SDL_SetWindowSize(window, header.width, header.height);
				width = header.width;
				height = header.height;
				hasLoaded = 1;
			}


		}

		// https://dzone.com/articles/sdl2-pixel-drawing (Derived from, only for switch)
        switch (event.type)
        {
            case SDL_QUIT:
                running = false;
                break;
            case SDL_MOUSEBUTTONUP:
                if (event.button.button == SDL_BUTTON_LEFT)
                    leftMouseButtonDown = false;
					previousY = -1;
					previousX = -1;
                break;
            case SDL_MOUSEBUTTONDOWN:
                if (event.button.button == SDL_BUTTON_LEFT)
                    leftMouseButtonDown = true;
            case SDL_MOUSEMOTION:
                if (leftMouseButtonDown)
                {
                    int mouseX = event.motion.x;
                    int mouseY = event.motion.y;

		    if(whichButton == 0 || whichButton == 1|| whichButton == 3 )
		    {
		         drawCircle(pixelArr, penColor, mouseX, mouseY, penWidth, width, height);
		    }
		    else if(whichButton == 2 && loop == 0)
		    {
			int const finalStringSize = 20;
			char finalString[finalStringSize];
			memset(finalString, 0, finalStringSize*sizeof(char));
			sprintf(finalString, "%d %d %d", pixelArr[mouseX + (mouseY * width)].red, pixelArr[mouseX + (mouseY * width)].green, pixelArr[mouseX + (mouseY * width)].blue);

			zmq_send(responder, finalString, finalStringSize, 0);
			loop = 1;
                    }

		    previousY = mouseY;
		    previousX = mouseX;
					
                }
                break;
        }
        int num  = height*width;
	for(int i = 0; i < num; i++)
	{	
    	    SDL_SetRenderDrawColor(renderer, pixelArr[i].red, pixelArr[i].green, pixelArr[i].blue, 255);
	    SDL_RenderDrawPoint(renderer, i%width, i/width);
	}
			
        SDL_RenderPresent(renderer);
    }

    // Release resources
    free(pixelArr);
    free(penColor.red);
    free(penColor.green);
    free(penColor.blue);
    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();

    return 0;
}


