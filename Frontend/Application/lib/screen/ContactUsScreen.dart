import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class ContactUsScreen extends StatelessWidget {
  final TextEditingController _feedbackController = TextEditingController();

  ContactUsScreen({super.key});

  void _submitFeedback(BuildContext context) async {
    String feedback = _feedbackController.text;
    String subject = 'Feedback from Contact Us Form';
    String body = 'Feedback: $feedback';

    final Uri emailLaunchUri = Uri(
      scheme: 'mailto',
      path: 'disgitalscrumassistant.1440@gmail.com',
      queryParameters: {
        'subject': subject,
        'body': body,
      },
    );

    String emailUri = emailLaunchUri.toString();
    // ignore: deprecated_member_use
    if (await canLaunch(emailUri)) {
      // ignore: deprecated_member_use
      await launch(emailUri);
    } else {
      // ignore: use_build_context_synchronously
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to send email')),
      );
    }

    _feedbackController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Contact Us'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Feedback Form',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16.0),
            TextField(
              controller: _feedbackController,
              maxLines: 4,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Your Feedback',
              ),
            ),
            const SizedBox(height: 16.0),
            ElevatedButton(
              onPressed: () => _submitFeedback(context),
              child: const Text('Submit'),
            ),
            const Spacer(),
            InkWell(
              onTap: () => _submitFeedback(context),
              child: const Text(
                'disgitalscrumassistant.1440@gmail.com',
                style: TextStyle(
                  color: Colors.blue,
                  decoration: TextDecoration.underline,
                ),
                textAlign: TextAlign.center,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

void main() {
  runApp(MaterialApp(
    home: ContactUsScreen(),
  ));
}
