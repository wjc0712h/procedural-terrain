from manim import *

class Intro(Scene):
    def construct(self):
        text = Text('Hello world').scale(3)
        self.add(text)