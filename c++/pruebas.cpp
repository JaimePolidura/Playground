#include <iostream>
#include "./list/linkedlist.h"

struct Color {
    Color(){
        printf("[Color] Construido empty: %i\n", this);
    }

    Color(Color * other){
        printf("[Color] Construido copy: %i\n", this);
    }

    Color(const Color& other){
        printf("[Color] Construido copy: %i\n", this);
    }

    ~ Color(){
        printf("[Color] Deleted %i\n", this);
    }

    int rgb{};
};

struct Point {
public:
    Point(){
        printf("[Point] Construido empty: %i\n", this);
    }

    Point(const Point& other){
        printf("[Point] Construido copy: %i\n", this);
    }

    ~ Point(){
        printf("[Point] Deleted %i\n", this);
    }

    int x;
    int y;
    int z;
    Color * color;
};

Point& make_point () {
    Point * point = new Point();
    Color * color = new Color();
    point->color = color;

    return * point;
};

void show_point_ref(Point& point){
    printf("[Point] show_point_ref Direcccion ref: %i\n", &point);
    printf("[Color] show_point_ref Direcccion ref: %i\n", point.color);
}

void show_point(Point point){
    printf("[Point] show_point Direcccion: %i\n", &point);
    printf("[Color] show_point Point: %i\n", point.color);
}

int main(){
    auto * list = new Linkedlist<Point>();

    Point& ref = make_point(); //From heap
    printf("[Point] Direcccion: %i\n", &ref);
    printf("[Color] Direcccion: %i\n", ref.color);
    show_point(ref);
    show_point_ref(ref);

    printf("------ LIST --------\n");
    list->add(ref);

    return 0;
}