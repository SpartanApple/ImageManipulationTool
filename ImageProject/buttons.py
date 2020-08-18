from tkinter import *
# https://zeromq.org/get-started/?language=python&library=pyzmq# (Derived From)
import zmq
from tkinter import colorchooser
from tkinter import filedialog
from tkinter import Scale
import os

top = Tk()

def saveFile(): # 4
	context = zmq.Context()

	#  Socket to talk to server
	# https://zeromq.org/get-started/?language=c# (Reference)
	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")
	
	socket.send(b"4")




def fileExplorerOpen():
	global socket

	context = zmq.Context()

	#  Socket to talk to server
	# https://zeromq.org/get-started/?language=c# (Reference)
	print("Connecting to hello world server…")
	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")


	# https://www.geeksforgeeks.org/file-explorer-in-python-using-tkinter/ (How to use the file explorer tool)
	# https://stackoverflow.com/questions/46784609/python-tkinter-browse-directory-and-save-to-new-directory
	curr_directory = os.getcwd()
	fileName = filedialog.askopenfilename(initialdir = curr_directory, title = "Select a bitmap file", filetypes = (("BMP files", "*.bmp"),("all files","*.*")))

	fileName5 = "5 "+ fileName	

	fileByte = fileName5.encode("ASCII")

	socket.send(fileByte)



def eyeDropper():
	global socket
	context = zmq.Context()

	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")
	
	socket.send(b"2")
	colors = socket.recv()


def colourSelect():
	global socket
	#https://docs.python.org/3.9/library/tkinter.colorchooser.html
	color_code = colorchooser.askcolor(title ="Choose color")
	if(color_code == ""):
		return

	context = zmq.Context()

	#  Socket to talk to server
	# https://zeromq.org/get-started/?language=c# (Reference)
	print("Connecting to hello world server…")
	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")
	
	red = round(color_code[0][0])
	red = 255 if red>255 else red
	red = 0 if red<0 else red
	green = round(color_code[0][1])
	green = 255 if green>255 else green
	green = 0 if green<0 else green
	blue = round(color_code[0][2])
	blue = 255 if blue>255 else blue
	blue = 0 if blue<0 else blue

	colorString = "0 " + str(red) +" "+ str(green) +" "+ str(blue)
	print(colorString)
	byteColor = colorString.encode("ASCII")

	socket.send(byteColor)

def updatePenSize():
	global socket	
	context = zmq.Context()

	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")
	
	# https://www.tutorialspoint.com/python/tk_scale.htm (Derived from)
	penSizeString = "1 " + str(penSizeSliderVar.get())
	print(penSizeString)
	penSize = penSizeString.encode("ASCII")

	socket.send(penSize)

def eraserTool():
	global socket
	context = zmq.Context()

	socket = context.socket(zmq.REQ)
	socket.connect("tcp://localhost:5555")
	
	eraserString = "3"
	
	eraserEncoded = eraserString.encode("ASCII")

	socket.send(eraserEncoded)

#Handles tkinter button making

B1 = Button(top, text ="Open File Explorer", command = fileExplorerOpen)

B1.pack()

saveButton = Button(top, text="Save File", command = saveFile)
saveButton.pack()


B2 = Button(top, text ="Select Colour", command = colourSelect)
B2.pack()

B3 = Button(top, text ="Update Pen Size", command = updatePenSize)
B3.pack()

# https://www.tutorialspoint.com/python/tk_scale.htm (Derived from)
penSizeSliderVar = DoubleVar()
penSizeSlider = Scale(top, from_=0, to=100, var = penSizeSliderVar, orient=HORIZONTAL)
penSizeSlider.pack()



eyeDropper = Button(top, text ="Eye Dropper", command = eyeDropper)
eyeDropper.pack()

eraser = Button(top, text ="Eraser", command = eraserTool)
eraser.pack()


top.mainloop()

