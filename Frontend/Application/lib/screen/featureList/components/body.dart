import 'package:flutter/material.dart';
import '../../../size_config.dart';
import 'package:digital_scrum_assistant/screen/featureList/StandupRecordForm.dart';

class FeaturePlayListScreen extends StatelessWidget {
  FeaturePlayListScreen({Key? key}) : super(key: key);

  // Define a list of features
  final List<String> features = [
    'Classic Text Based Standups',
    'Remote Based Standups',
    'Mood Report',
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Feature Play List'),
      ),
      body: Column(
        children: [
          const Padding(
            padding: EdgeInsets.all(10),
            child: Text(
              'Feature List',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 24,
              ),
            ),
          ),
          Expanded(
            child: Center(
              child: SizedBox(
                width: SizeConfig.screenWidth * 0.7,
                child: ListView.builder(
                  itemCount: features.length,
                  itemBuilder: (BuildContext context, int index) {
                    return Padding(
                      padding: const EdgeInsets.symmetric(vertical: 10),
                      child: GestureDetector(
                        onTap: () {
                          if (index == 0) {
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                  builder: (context) => StandupRecordForm()),
                            );
                          }
                        },
                        child: Container(
                          height: SizeConfig.screenHeight * 0.3,
                          decoration: BoxDecoration(
                            color: const Color.fromARGB(255, 245, 238, 227),
                            borderRadius: BorderRadius.zero,
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withOpacity(0.5),
                                spreadRadius: 2,
                                blurRadius: 5,
                                offset: const Offset(0, 3),
                              ),
                            ],
                          ),
                          child: Column(
                            children: [
                              Image.asset(
                                'assets/images/feature${index + 1}.jpeg', // Assuming images are named as feature1.jpeg, feature2.jpeg, etc.
                                width:
                                    310, // Adjust the width of the image as needed
                                height:
                                    140, // Adjust the height of the image as needed
                                fit: BoxFit
                                    .cover, // Adjust the fit of the image as needed
                              ),
                              const SizedBox(
                                  height:
                                      10), // Add spacing between image and text
                              Center(
                                child: Text(
                                  features[index],
                                  style: const TextStyle(
                                    fontWeight: FontWeight.bold,
                                    fontSize: 20,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
