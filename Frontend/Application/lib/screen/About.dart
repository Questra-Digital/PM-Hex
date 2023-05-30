import 'package:flutter/material.dart';
import '../size_config.dart';
import 'ContactUsScreen.dart';

class AboutScreen extends StatelessWidget {
  const AboutScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('About Us'),
      ),
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header
            Container(
              // Company logo and navigation menu can be added here
              alignment: Alignment.center,
              child: Image.asset(
                'assets/AboutPic1.jpg',
                width: 200,
                height: 200,
              ),
            ),
            const SizedBox(height: 16),

            // Introduction
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Welcome to Digital Scrum Assistant!',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.0),
              child: Text(
                'We are dedicated to automating agile processes and enhancing project management efficiency.',
                style: TextStyle(fontSize: 16),
              ),
            ),
            SizedBox(height: 16),

            // Team
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Meet Our Team',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            // Team members' information and photos can be added here

            // Services
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Our Services',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.0),
              child: Text(
                'At Digital Scrum Assistant, we offer cutting-edge solutions for automating agile processes, optimizing workflow, and enhancing collaboration within teams.',
                style: TextStyle(fontSize: 16),
              ),
            ),

            // Client Testimonials
            // Testimonials from satisfied clients can be added here

            // Achievements or Milestones
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Our Achievements',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            // Overview of achievements or notable projects can be added here

            // Company Culture
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Our Company Culture',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.0),
              child: Text(
                'At Digital Scrum Assistant, we foster a culture of innovation, collaboration, and continuous learning. Our team is passionate about delivering exceptional results and exceeding client expectations.',
                style: TextStyle(fontSize: 16),
              ),
            ),

            // Contact Information
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => ContactUsScreen(),
                  ),
                );
              },
              child: Text(
                'Contact Us',
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: getProportionateScreenWidth(14),
                ),
              ),
            ),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.0),
              child: Text(
                'Since 2023\nLahore, Pakistan\nPhone: +92 344 4080878\nEmail: digitalscrumassistant.1440\n\nFollow us on social media:\n- Facebook: facebook.com/digitalscrumassistant\n- Twitter: twitter.com/digitalscrumassistant\n- Instagram: instagram.com/digitalscrumassistant\n\nWe look forward to hearing from you and collaborating on your agile journey!',
                style: TextStyle(fontSize: 16),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
