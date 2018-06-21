"""Python Module Checks"""

import enum
import subprocess

from packaging import version

from .exc import CheckFailure


class ComparisonOps(enum.Enum):
    """An enumeration of all the available comparison operators for
    PythonModuleCheck."""
    EQUALS = 'eq'
    GTE = 'gte'
    LTE = 'lte'
    GT = 'gt'
    LT = 'lt'


class PythonModuleCheck:
    """Checks if a given python module is installed. Optionally
    checks the module version.

    Args:
        name (required): String name of the module to check for
        version: String version specifier. Remember that for yaml you should
            wrap this in quotes or 1.0 will be considered a float and
            cause an error. Example: '1.0' not 1.0
        python: Which python to run, if not provided uses the first 'python'
            executable found in $PATH. Should be an absolute file path.
        comparison: How to compare the installed version. Defaults to gte
            (greater than or equal to). Available values are:
                - gte: Greater than or equals to (module_version >= specified_version)
                - lte: Less than or equals to (module_version <= specified_version)
                - gt: Greater than (module_version > specified_version)
                - lt: Less than (module_version < specified_version)
                - eq: Equal to (module_version == specified_version)
        min_version: String version specifier if specified then the python
            module version will be checked to be between version and
            min_version
        statement: Python statement used to get the module version. Defaults to
            {module_name}.__version__ but not all python modules use this convention.
            Should be a valid expression to put inside of a print call like: 'print(%s)'
    """

    def __init__(
            self,
            module,
            version=None,  #pylint: disable-msg=redefined-outer-name
            python='python',
            comparison='gte',
            statement=None,
            min_version=None):
        self.module = module
        self.version = version
        self.python = python
        self.comparison = comparison
        self.min_version = min_version
        if min_version is not None and version is None:
            raise CheckFailure(
                'Invalid arguments: min_version requires that version is set')

        self.statement = statement
        if statement is None:
            self.statement = '{}.__version__'.format(self.module)

    def version_check(self, installed_ver, expected_ver, minimum_ver=None):
        """Compare installed_ver with expected_ver according to self.comparison
        or verify that installed_ver is between minimum and expected_ver"""
        if minimum_ver:
            return minimum_ver <= installed_ver <= expected_ver
        elif self.comparison == ComparisonOps.EQUALS.value:
            return installed_ver == expected_ver
        elif self.comparison == ComparisonOps.LTE.value:
            return installed_ver <= expected_ver
        elif self.comparison == ComparisonOps.GT.value:
            return installed_ver > expected_ver
        elif self.comparison == ComparisonOps.LT.value:
            return installed_ver < expected_ver
        # Default to greater than or equal to
        return installed_ver >= expected_ver

    def check(self):
        """Check that the module is installed and optionally the
        correct version."""
        proc = subprocess.run(
            [
                self.python,
                '-c',
                'import {module}; print({statement})'.format(
                    module=self.module, statement=self.statement),
            ],
            stdout=subprocess.PIPE)

        if proc.returncode != 0:
            raise CheckFailure('{} module is not installed')

        # End early if version check not required
        if self.version is None:
            return

        output = proc.stdout.decode('utf-8').strip()
        # Sometimes Python 2 will print the parens. This will strip them out
        if output[0] == '(':
            output = output[1:-1]

        installed_ver = version.parse(output)
        expected_ver = version.parse(self.version)
        minimum_ver = None
        if self.min_version:
            minimum_ver = version.parse(self.min_version)

        if not self.version_check(
                installed_ver, expected_ver, minimum_ver=minimum_ver):
            raise CheckFailure('{} is version: {}'.format(self.module, output))
