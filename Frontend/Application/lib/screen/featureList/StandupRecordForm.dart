import 'package:flutter/material.dart';

class StandupRecordForm extends StatefulWidget {
  const StandupRecordForm({Key? key}) : super(key: key);

  @override
  _StandupRecordFormState createState() => _StandupRecordFormState();
}

class _StandupRecordFormState extends State<StandupRecordForm> {
  final _formKey = GlobalKey<FormState>();
  String _title = '';
  String _participants = '';
  String _updates = '';
  String _email = '';
  List<String> _daysOfWeek = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Column(children: [
      // Header
      Container(
        decoration: const BoxDecoration(
          image: DecorationImage(
            image: AssetImage('assets/f1.png'),
            fit: BoxFit.cover,
          ),
        ),
        alignment: Alignment.topCenter,
        padding: const EdgeInsets.only(top: 16),
        child: const Text(
          'Classic Text Based Standups',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      Row(
        children: [
          // Navigation Panel
          Container(
            width: MediaQuery.of(context).size.width * 0.2,
            color: Color.fromARGB(153, 238, 130, 7),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Updates:',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  '${DateTime.now().toString()}',
                  style: const TextStyle(fontSize: 16),
                ),
              ],
            ),
          ),

          // Form Content
          Expanded(
            child: SingleChildScrollView(
              child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 16),
                    Form(
                      key: _formKey,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          TextFormField(
                            decoration: const InputDecoration(
                              hintText: 'Enter meeting title',
                              labelText: 'Title',
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter title';
                              }
                              return null;
                            },
                            onSaved: (value) {
                              _title = value!;
                            },
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            decoration: const InputDecoration(
                              hintText: 'Enter participants names',
                              labelText: 'Participants',
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter participants';
                              }
                              return null;
                            },
                            onSaved: (value) {
                              _participants = value!;
                            },
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            decoration: const InputDecoration(
                              hintText: 'Enter your email',
                              labelText: 'Email',
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter email';
                              }
                              return null;
                            },
                            onSaved: (value) {
                              _email = value!;
                            },
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'Select Days of Week:',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Wrap(
                            spacing: 8,
                            children: [
                              _buildDayOfWeekChip('Mon'),
                              _buildDayOfWeekChip('Tue'),
                              _buildDayOfWeekChip('Wed'),
                              _buildDayOfWeekChip('Thu'),
                              _buildDayOfWeekChip('Fri'),
                              _buildDayOfWeekChip('Sat'),
                              _buildDayOfWeekChip('Sun'),
                            ],
                          ),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () {
                              if (_formKey.currentState!.validate()) {
                                _formKey.currentState!.save();
                                _submitForm();
                              }
                            },
                            child: const Text('Submit'),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    ]));
  }

  Widget _buildDayOfWeekChip(String label) {
    final isSelected = _daysOfWeek.contains(label);
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        setState(() {
          if (selected) {
            _daysOfWeek.add(label);
          } else {
            _daysOfWeek.remove(label);
          }
        });
      },
    );
  }

  void _submitForm() {
// TODO: Implement submitting form data
    print('Title: $_title');
    print('Participants: $_participants');
    print('Updates: $_updates');
    print('Email: $_email');
    print('Days of Week: $_daysOfWeek');
  }
}
