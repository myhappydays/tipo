
import requests
import unittest
import os

# Load environment variables from .env file
from dotenv import load_dotenv
load_dotenv()

class TestTrendingKeywordsAPI(unittest.TestCase):

    BASE_URL = "http://localhost:8080/api/trending_keywords"

    def test_get_default_trending_keywords(self):
        """
        Test if the API returns a list of default trending keywords (5 items).
        """
        response = requests.get(self.BASE_URL)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIsInstance(data, list)
        self.assertEqual(len(data), 5)
        self.assertIsInstance(data[0], str)

    def test_get_specific_count_trending_keywords(self):
        """
        Test if the API returns a specific count of trending keywords.
        """
        params = {"count": 3}
        response = requests.get(self.BASE_URL, params=params)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIsInstance(data, list)
        self.assertEqual(len(data), 3)

    def test_get_random_trending_keyword(self):
        """
        Test if the API returns a single random trending keyword.
        """
        params = {"random": "true"}
        response = requests.get(self.BASE_URL, params=params)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIsInstance(data, dict)
        self.assertIn("keyword", data)
        self.assertIsInstance(data["keyword"], str)


if __name__ == '__main__':
    # To run this test, you need to have the backend server running.
    # Also, make sure to install the required libraries:
    # pip install requests python-dotenv
    unittest.main()
