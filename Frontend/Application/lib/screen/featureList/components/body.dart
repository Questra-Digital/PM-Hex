import 'package:flutter/material.dart';
import '../../../size_config.dart';

class FeaturePlayListScreen extends StatelessWidget {
  FeaturePlayListScreen({Key? key}) : super(key: key);

  // Define a list of features
  final List<String> features = [
    'Classic Text Based Standup Meetings',
    'Remote Based Standup Meetings',
    'Mood Report',
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Feature Play List'),
      ),
      body: Center(
        child: SizedBox(
          width: SizeConfig.screenWidth * 0.4,
          child: ListView.builder(
            itemCount: features.length,
            itemBuilder: (BuildContext context, int index) {
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 10),
                child: Container(
                  height: SizeConfig.screenHeight * 1.2,
                  decoration: BoxDecoration(
                    color: Colors.orange,
                    borderRadius: BorderRadius.circular(10),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.5),
                        spreadRadius: 2,
                        blurRadius: 5,
                        offset: const Offset(0, 3),
                      ),
                    ],
                  ),
                  child: Center(
                    child: Text(
                      features[index],
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 20,
                      ),
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}
