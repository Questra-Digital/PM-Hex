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
            SizedBox(height: getProportionateScreenHeight(20)),
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: getProportionateScreenWidth(20)),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Digital Scrum Assistant',
                    style: TextStyle(
                      color: kPrimaryColor,
                      fontWeight: FontWeight.bold,
                      fontSize: getProportionateScreenWidth(18),
                    ),
                  ),
                  Row(
                    children: [
                      IconButton(
                        onPressed: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                                builder: (context) => FeaturePlayListScreen()),
                          );
                        },
                        icon: const Icon(Icons.featured_play_list),
                      ),
                      IconButton(
                        onPressed: () {},
                        icon: const Icon(Icons.person),
                      ),
                      IconButton(
                        onPressed: () {},
                        icon: const Icon(Icons.contact_mail),
                      ),
                      IconButton(
                        onPressed: () {},
                        icon: const Icon(Icons.info_outline),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const HomeHeader(),
            SizedBox(height: getProportionateScreenWidth(10)),
          ],
        ),
      ),
    );
  }
}
