import 'package:flutter/material.dart';
import 'package:digital_scrum_assistant/screen/featureList/components/body.dart';
import '../../../constant.dart';
import '../../../size_config.dart';
import 'home_header.dart';

class Body extends StatelessWidget {
  const Body({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: SingleChildScrollView(
        child: Column(
          children: [
            // Logo/Branding
            Image.asset(
              'assets/images/logo.png', // Replace with your actual logo image path
              width: getProportionateScreenWidth(
                  100), // Set the desired width for your logo
              height: getProportionateScreenWidth(
                  100), // Set the desired height for your logo
            ),
            SizedBox(height: getProportionateScreenHeight(20)),
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: getProportionateScreenWidth(20)),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  // Navigation/Menu
                  Row(
                    children: [
                      TextButton(
                        onPressed: () {
                          // Handle menu item 1 press
                        },
                        child: Text(
                          'Menu Item 1',
                          style: TextStyle(
                            color: kPrimaryColor,
                            fontWeight: FontWeight.bold,
                            fontSize: getProportionateScreenWidth(14),
                          ),
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          // Handle menu item 2 press
                        },
                        child: Text(
                          'Menu Item 2',
                          style: TextStyle(
                            color: kPrimaryColor,
                            fontWeight: FontWeight.bold,
                            fontSize: getProportionateScreenWidth(14),
                          ),
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          // Handle menu item 3 press
                        },
                        child: Text(
                          'Menu Item 3',
                          style: TextStyle(
                            color: kPrimaryColor,
                            fontWeight: FontWeight.bold,
                            fontSize: getProportionateScreenWidth(14),
                          ),
                        ),
                      ),
                    ],
                  ),
                  // CTA (Call-to-Action) Button
                  ElevatedButton(
                    onPressed: () {
                      // Handle CTA button press
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
            SizedBox(height: getProportionateScreenHeight(20)),
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
            // Feature Highlights
            // Add your feature highlights widgets or content here
            SizedBox(height: getProportionateScreenHeight(20)),
            // Testimonials
            // Add your testimonials widgets or content here
            SizedBox(height: getProportionateScreenHeight(20)),
            // Social Proof
            // Add your social proof widgets or content here
            SizedBox(height: getProportionateScreenHeight(20)),
            // your footer widgets or content here
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: getProportionateScreenWidth(20)),
              child: Column(
                children: [
                  // Add your footer content here, such as contact information, links, etc.
                  Text(
                    'Contact us: contact@example.com',
                    style: TextStyle(
                      color: kTextColor,
                      fontSize: getProportionateScreenWidth(14),
                    ),
                  ),
                  SizedBox(height: getProportionateScreenHeight(10)),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      IconButton(
                        onPressed: () {
                          // Handle social media icon 1 press
                        },
                        icon: const Icon(
                          Icons.facebook,
                          color: kPrimaryColor,
                        ),
                      ),
                      // IconButton(
                      //   onPressed: () {
                      //     // Handle social media icon 2 press
                      //   },
                      //   icon: Icon(
                      //     Icons.twitter,
                      //     color: kPrimaryColor,
                      //   ),
                      // ),
                      // IconButton(
                      //   onPressed: () {
                      //     // Handle social media icon 3 press
                      //   },
                      //   icon: Icon(
                      //     Icons.instagram,
                      //     color: kPrimaryColor,
                      //   ),
                      // ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
