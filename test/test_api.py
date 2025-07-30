
import requests
import unittest
import os

# Load environment variables from .env file
from dotenv import load_dotenv
load_dotenv()

class TestContentAPI(unittest.TestCase):

    BASE_URL = "http://localhost:8080/api/content"

    def test_get_news_content_successfully(self):
        """
        Test if the API returns a well-formed list of news content successfully.
        """
        params = {"query": "IT", "start": 1, "display": 10}
        response = requests.get(self.BASE_URL, params=params)

        # Check if the request was successful
        self.assertEqual(response.status_code, 200)

        # Check the content of the response - expect a list of content objects
        data = response.json()
        self.assertIsInstance(data, list)
        self.assertGreater(len(data), 0) # Expect at least one item

        # Check the structure of the first item
        first_item = data[0]
        self.assertIn("title", first_item)
        self.assertIn("source", first_item)
        self.assertIn("keyword", first_item) # Changed from category to keyword
        self.assertIn("pubDate", first_item)
        self.assertIn("sentences", first_item)

        # Check data types and values for the first item
        self.assertIsInstance(first_item["title"], str)
        self.assertIsInstance(first_item["source"], str)
        self.assertNotEqual(first_item["source"], "") # Source should not be empty
        self.assertEqual(first_item["keyword"], "IT") # Assuming query is 'IT'
        self.assertIsInstance(first_item["pubDate"], str)
        self.assertIsInstance(first_item["sentences"], list)

        # Save the response to a file for verification
        import json
        with open("content_response.json", "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print('''
--- Content API Response saved to content_response.json ---''')

    def test_missing_parameters(self):
        """
        Test the API's response when required parameters are missing.
        """
        # Case 1: Missing 'query'
        params = {"start": 1, "display": 10}
        response = requests.get(self.BASE_URL, params=params)
        self.assertEqual(response.status_code, 400)

    def test_multiple_keywords_and_interleaving(self):
        """
        Test if the API handles multiple keywords and interleaves results.
        """
        params = {"query": "IT,경제", "start": 1, "display": 10}
        response = requests.get(self.BASE_URL, params=params)

        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIsInstance(data, list)
        self.assertEqual(len(data), 10) # Expect 10 items (5 from each keyword)

        # Basic check for interleaving (e.g., first item from IT, second from Economy, etc.)
        # This is a simplified check, a more robust test would analyze content.
        # For now, just check if titles are different, implying mixed sources.
        self.assertNotEqual(data[0]["title"], data[1]["title"])

        # Save the response to a file for verification
        import json
        with open("content_interleaved_response.json", "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print('''
--- Content Interleaved API Response saved to content_interleaved_response.json ---''')


if __name__ == '__main__':
    # To run this test, you need to have the backend server running.
    # Also, make sure to install the required libraries:
    # pip install requests python-dotenv
    unittest.main()

