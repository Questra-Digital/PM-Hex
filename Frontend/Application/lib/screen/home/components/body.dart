import 'package:digital_scrum_assistant/screen/About.dart';
import 'package:digital_scrum_assistant/screen/ContactUsScreen.dart';
import 'package:flutter/material.dart';
import 'package:digital_scrum_assistant/screen/featureList/components/body.dart';
import '../../../size_config.dart';
import 'home_header.dart';

class Body extends StatelessWidget {
  const Body({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Logo/Branding
            SizedBox(height: getProportionateScreenHeight(20)),
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: getProportionateScreenWidth(20)),
              child: Text(
                'Digital Scrum Assistant',
                style: TextStyle(
                  color: Colors.black87,
                  fontWeight: FontWeight.bold,
                  fontSize: getProportionateScreenWidth(18),
                ),
              ),
            ),
            SizedBox(height: getProportionateScreenHeight(20)),
            const HomeHeader(),
            SizedBox(height: getProportionateScreenWidth(10)),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                ElevatedButton(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => FeaturePlayListScreen(),
                      ),
                    );
                  },
                  child: Text(
                    'Feature List',
                    style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: getProportionateScreenWidth(14),
                    ),
                  ),
                ),
                SizedBox(width: getProportionateScreenWidth(10)),
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
                SizedBox(width: getProportionateScreenWidth(10)),
                ElevatedButton(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => const AboutScreen(),
                      ),
                    );
                  },
                  child: Text(
                    'About',
                    style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: getProportionateScreenWidth(14),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
                height: getProportionateScreenHeight(
                    10)), // Add SizedBox with desired height
            // Hero Section
            Container(
              height: getProportionateScreenHeight(
                  200), // Set the desired height for your Hero Section
              decoration: const BoxDecoration(
                image: DecorationImage(
                  image: AssetImage(
                      'assets/images/hero_image.png'), // Replace with your actual hero image path
                  fit: BoxFit.cover,
                ),
              ),
              // Add any additional widgets or content for your Hero Section
            ),
            SizedBox(height: getProportionateScreenHeight(20)),
            // CTA (Call-to-Action) Button
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => FeaturePlayListScreen(),
                  ),
                );
              },
              child: Text(
                'Get Started',
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: getProportionateScreenWidth(14),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
