from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import random
import logging

chrome_options = Options()
chrome_options.add_argument('--headless')
chrome_options.add_argument('--disable-logging') 
logging.getLogger('selenium').setLevel(logging.CRITICAL)
login_count = 0
passwords = ['1234','1234','1234','1234','1235']

def execute_script():
    driver = webdriver.Chrome(options=chrome_options)

    try:
        driver.get("http://localhost:8080/auth")

        password_field = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.NAME, "password"))
        )

        random_password = random.choice(passwords)

        if random_password == '1234':
            print("Sending correct password")
        else:
            print("Sending wrong password")

        password_field.send_keys(random_password)

        login_button = driver.find_element("css selector", "input[type='button']")
        login_button.click()

        time.sleep(2)

        print("Completed")

    finally:
        driver.quit()

try:
    print("Starting Script")
    while True:
        login_count += 1
        print("Trying login number {}".format(login_count))
        execute_script()
        time.sleep(5)

except KeyboardInterrupt:
    print("Stopped Script")
