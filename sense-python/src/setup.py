from setuptools import setup, find_packages

setup(
    name='bubblesnet',
    version='0.1.0',
    packages=find_packages(include=['', '.*']),
    setup_requires=['pytest-runner'],
    tests_require=['pytest']
)