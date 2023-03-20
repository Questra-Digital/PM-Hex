import 'package:flutter/material.dart';
import '../../../size_config.dart';

class FeaturePlayListScreen extends StatelessWidget {
  const FeaturePlayListScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Feature Play List'),
      ),
      body: const Center(
        child: Text(
          'This is the feature play list screen.\n'
          'Add your feature list here.',
          textAlign: TextAlign.center,
        ),
      ),
    );
  }
}
