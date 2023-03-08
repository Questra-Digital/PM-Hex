import 'package:flutter/widgets.dart';
import 'package:digital_scrum_assistant/screen/forgot_password/forgot_password_screen.dart';
import 'package:digital_scrum_assistant/screen/home/home_screen.dart';
import 'package:digital_scrum_assistant/screen/login_success/login_success_screen.dart';
import 'package:digital_scrum_assistant/screen/otp/otp_screen.dart';
import 'package:digital_scrum_assistant/screen/complete_profile/complete_profile.dart';
import 'package:digital_scrum_assistant/screen/signin/signin_screen.dart';
import 'package:digital_scrum_assistant/screen/splash/splash_screen.dart';

import 'screen/signup/signup_page.dart';

// We use name route
// All our routes will be available here
final Map<String, WidgetBuilder> routes = {
  SplashScreen.routeName: (context) => const SplashScreen(),
  SignInScreen.routeName: (context) => const SignInScreen(),
  ForgotPasswordScreen.routeName: (context) => const ForgotPasswordScreen(),
  LoginSuccessScreen.routeName: (context) => const LoginSuccessScreen(),
  SignUpScreen.routeName: (context) => const SignUpScreen(),
  CompleteProfileScreen.routeName: (context) => const CompleteProfileScreen(),
  OtpScreen.routeName: (context) => const OtpScreen(),
  HomeScreen.routeName: (context) => const HomeScreen(),
};
