import pytest
import uuid
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from urllib.parse import urlparse

BASE_URL = 'http://127.0.0.1:8080'

@pytest.fixture(scope="module")
def driver():
    geckodriver_path = "/snap/bin/geckodriver"
    driver_service = webdriver.FirefoxService(executable_path=geckodriver_path)
    drv = webdriver.Firefox(service=driver_service)
    yield drv
    drv.quit()


# Get the path from a URL
def get_path(url):
    p = urlparse(url)
    return p.path

# Click on button that is beyound the screen
def click_workaround(driver, element):
    driver.execute_script("arguments[0].click();", element)
    # TODO: Get rid of sleep
    time.sleep(1)

def test_login_page_failure(driver):
    driver.get(BASE_URL + "/login")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    login_input.send_keys("admin")
    password_input.send_keys("wrong_password")
    submit_button.click()

    assert get_path(driver.current_url) == '/login'

    body = driver.find_element(By.TAG_NAME, "body")
    assert 'Invalid username or password' in body.text


def test_login_page_success(driver):
    driver.get(BASE_URL + "/login")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    login_input.send_keys("admin")
    password_input.send_keys("admin")
    submit_button.click()

    # Find the h1 element with the text "Users"
    users_header = driver.find_element(By.TAG_NAME, "h1")
    assert users_header.text == "Users"

    assert get_path(driver.current_url) == '/admin/users'


def test_user_management_page(driver):
    driver.get(BASE_URL + "/admin/users")

    # Find the h1 element with the text "Users"
    users_header = driver.find_element(By.TAG_NAME, "h1")
    assert users_header.text == "Users"

    # Verify table headers
    headers = driver.find_elements(By.TAG_NAME, "th")
    assert headers[0].text == "ID"
    assert headers[1].text == "Login"
    assert headers[2].text == "Actions"

    # Verify the presence of user rows in the table
    rows = driver.find_elements(By.CSS_SELECTOR, "table tr")
    assert len(rows) > 1  # More than one row indicates users are present

    # Verify first user actions
    first_user_actions = rows[1].find_elements(By.TAG_NAME, "a")
    assert "Edit" in first_user_actions[1].text
    delete_form = rows[1].find_element(By.TAG_NAME, "form")
    assert delete_form is not None


def test_user_create_failure(driver):
    driver.get(BASE_URL + "/admin/users/new")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    assert login_input is not None
    assert password_input is not None
    assert submit_button is not None

    # Optionally, you can fill out the form and submit it to test the form submission
    login_input.send_keys("1")
    password_input.send_keys("1")
    submit_button.click()

    body = driver.find_element(By.TAG_NAME, "body")
    assert "Validation error" in body.text
    assert "Validation error" in body.text
    assert "Field is too short" in body.text


def test_user_create_uniqueness_validation(driver):
    driver.get(BASE_URL + "/admin/users/new")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    assert login_input is not None
    assert password_input is not None
    assert submit_button is not None

    unique_name = str(uuid.uuid4())

    # Optionally, you can fill out the form and submit it to test the form submission
    login_input.send_keys(unique_name)
    password_input.send_keys("newpassword")
    submit_button.click()

    # Find the h1 element with the text "Users"
    users_header = driver.find_element(By.TAG_NAME, "h1")
    assert users_header.text == "Users"

    body = driver.find_element(By.TAG_NAME, "body")
    assert "User was added." in body.text

    body = driver.find_element(By.TAG_NAME, "body")
    assert unique_name in body.text

    driver.get(BASE_URL + "/admin/users/new")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    assert login_input is not None
    assert password_input is not None
    assert submit_button is not None

    # Optionally, you can fill out the form and submit it to test the form submission
    login_input.send_keys(unique_name)
    password_input.send_keys("newpassword")
    submit_button.click()

    body = driver.find_element(By.TAG_NAME, "body")
    assert 'ERROR: duplicate key value violates unique constraint' in body.text


def test_user_create_and_delete_success(driver):
    driver.get(BASE_URL + "/admin/users/new")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    assert login_input is not None
    assert password_input is not None
    assert submit_button is not None

    unique_name = str(uuid.uuid4())

    # Optionally, you can fill out the form and submit it to test the form submission
    login_input.send_keys(unique_name)
    password_input.send_keys("newpassword")
    submit_button.click()

    # Find the h1 element with the text "Users"
    users_header = driver.find_element(By.TAG_NAME, "h1")
    assert users_header.text == "Users"

    body = driver.find_element(By.TAG_NAME, "body")
    assert "User was added." in body.text

    body = driver.find_element(By.TAG_NAME, "body")
    assert unique_name in body.text

    # Find the delete button for the user with ID 7
    delete_button = driver.find_element(
        By.CSS_SELECTOR, f'[data-selenium="delete-{unique_name}"]')
    click_workaround(driver, delete_button)

    body = driver.find_element(By.TAG_NAME, "body")
    assert "User was deleted" in body.text


def test_user_create_and_edit_success(driver):
    driver.get(BASE_URL + "/admin/users/new")

    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")

    assert login_input is not None
    assert password_input is not None
    assert submit_button is not None

    unique_name = str(uuid.uuid4())

    # Optionally, you can fill out the form and submit it to test the form submission
    login_input.send_keys(unique_name)
    password_input.send_keys("newpassword")
    submit_button.click()

    # Find the h1 element with the text "Users"
    users_header = driver.find_element(By.TAG_NAME, "h1")
    assert users_header.text == "Users"

    body = driver.find_element(By.TAG_NAME, "body")
    assert "User was added." in body.text

    body = driver.find_element(By.TAG_NAME, "body")
    assert unique_name in body.text

    edit_button = driver.find_element(
        By.CSS_SELECTOR, "a[data-selenium='edit-" + unique_name + "']")
    click_workaround(driver, edit_button)

    header = driver.find_element(By.TAG_NAME, "h1")
    assert header.text == "Edit User"
    login_input = driver.find_element(By.ID, "login")
    password_input = driver.find_element(By.ID, "password")
    submit_button = driver.find_element(
        By.CSS_SELECTOR, "button[type='submit']")
    login_input.send_keys(unique_name + "_")
    submit_button.click()

    assert get_path(driver.current_url) == '/admin/users'
    body = driver.find_element(By.TAG_NAME, "body")
    assert unique_name + "_" in body.text


if __name__ == "__main__":
    pytest.main(["test.py"])
