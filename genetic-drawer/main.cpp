#include <SFML/Graphics.hpp>
#include <opencv2/core.hpp>
#include <opencv2/imgcodecs.hpp>
#include <opencv2/highgui.hpp>
#include <iostream>
#include <vector>

#include "GeneticArt.h"

using namespace cv;
using namespace std;

void drawProgressBar(int progress, int limit) {
	string progressBar = "   [";
	for (char i = 0; i < 20; i++) {
		if (i < progress * 20 / limit) {
			progressBar += "#";
		} else {
			progressBar += " ";
		}
	}
	progressBar += "]";
	cout << progressBar << " (" << progress << "/" << limit << ")\r";
}

int main(int argc, char** argv) {
	if (argc != 4) {
		cout << "Using default parameters: 50 50 50" << endl;
	}

	const int MAX_GENERATIONS = 100;

	vector<Mat> population = generateBeautyFromNoise(50, 50, 50);

	vector<int> aptitudes = evaluateArtisticAppeal(population);

	int avgAptitude = 0;

	for (int i : aptitudes) {
		avgAptitude += i;
	}

	cout << "Initial average aptitude: " << (avgAptitude / aptitudes.size()) << endl;

	for (int i = 0; i <= MAX_GENERATIONS; i++) {
		vector<int> aptitudes = evaluateArtisticAppeal(population);
		population = naturallySelectArt(population, aptitudes);

		int avgAptitude = 0;

		for (int i : aptitudes) {
			avgAptitude += i;
		}

		cout << "average aptitude: " << (avgAptitude / aptitudes.size()) << endl;

		//drawProgressBar(i, MAX_GENERATIONS);
	}


	aptitudes = evaluateArtisticAppeal(population);


	avgAptitude = 0;

	for (int i : aptitudes) {
		avgAptitude += i;
	}

	cout << "\nFinal average aptitude: " << (avgAptitude / aptitudes.size()) << endl;

	Mat relic = extractOpusMagnum(population);

	namedWindow("Opus Magnum", WINDOW_AUTOSIZE); // Create a window for display.
	imshow("Opus Magnum", relic); // Show our image inside it.
	waitKey(0); // Wait for a keystroke in the window
}

/*

int main(int argc, char** argv)
{
	if (argc != 2)
	{
		cout << " Usage: display_image ImageToLoadAndDisplay" << endl;
		return -1;
	}
	Mat image;
	image = imread(argv[1], IMREAD_COLOR); // Read the file
	if (image.empty()) // Check for invalid input
	{
		cout << "Could not open or find the image" << std::endl;
		return -1;
	}
	namedWindow("Display window", WINDOW_AUTOSIZE); // Create a window for display.
	imshow("Display window", image); // Show our image inside it.
	waitKey(0); // Wait for a keystroke in the window
	return 0;
}

int main()
{
	sf::RenderWindow window(sf::VideoMode(200, 200), "SFML works!");
	sf::CircleShape shape(100.f);
	shape.setFillColor(sf::Color::Green);

	while (window.isOpen())
	{
		sf::Event event;
		while (window.pollEvent(event))
		{
			if (event.type == sf::Event::Closed)
				window.close();
		}

		window.clear();
		window.draw(shape);
		window.display();
	}

	return 0;
}
*/